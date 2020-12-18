package userFullDataStorage

import (
	"context"
	"fmt"
	"os"
	"time"

	"backend/model"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const userDataCollection = "users"
const mainDBName = "matcha"
const yearInMilisecs = 31207680000

type ManagerStruct struct {
	Conn *mongo.Client
}

type MongoCoords struct {
	Type	string	`bson:"type"`
	Coordinates		[]float64	`bson:"coordinates"`
}

func (m *ManagerStruct) MakeConnection() {
	var err error
	user := os.Getenv("MONGO_USER")
	password := os.Getenv("MONGO_PASSWORD")
	addr := os.Getenv("MONGO_ADDRESS")

	if user == "" || password == "" || addr == "" {
		logrus.Error("Env is empty", user, password, addr)
	}

	connStr := fmt.Sprintf("mongodb://%v:%v@%v", user, password, addr)
	logrus.Info("Connecting to mongo: ", connStr)
	opts := options.Client().ApplyURI(connStr)
	m.Conn, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		logrus.Fatal("Error getting mongo wsClient: ", err)
	}
	if err != nil {
		logrus.Fatal("Error connecting to mongo: ", err)
	}
	if err = m.Conn.Ping(context.TODO(), readpref.Primary()); err != nil {
		logrus.Fatal("Error pinging: ", err)
	}
	logrus.Info("Connected.")
}

func (m *ManagerStruct) CreateLocationIndex() {
	database := m.Conn.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	key := "position_2dsphere"

	if res, err := userCollection.Indexes().CreateOne(
		context.TODO(), mongo.IndexModel{Keys: bson.D{{"position", "2dsphere"}}, Options: &options.IndexOptions{Name: &key}}); err != nil {
		logrus.Errorf("Error creating index: %v", err)
	} else {
		logrus.Infof("Created index: %v", res)
	}

}

func (m *ManagerStruct) CloseConnection() {
	if err := m.Conn.Disconnect(context.TODO()); err != nil {
		logrus.Error("Error closing mongo: ", err)
	}
}

func (m *ManagerStruct) AddTagToUserTags(user model.LoginData, tagId int64) bool {
	database := m.Conn.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	filter := bson.M{"id": user.Id}
	update := bson.D{{"$addToSet", bson.D{{"tag_ids", tagId}}}}

	res, err := userCollection.UpdateOne(context.TODO(), filter, update)
	if  err != nil {
		logrus.Error("Error updating user document: ", err)
		return false
	}
	if res.MatchedCount != 1 {
		logrus.Error("Error find user document (res.MatchedCount != 1): ", err)
		return false
	}
	return true
}

func (m *ManagerStruct) DeleteTagFromUserTags(user model.LoginData, tagId int64) bool {
	database := m.Conn.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	filter := bson.M{"id": user.Id}
	update := bson.D{{"$pull", bson.D{{"tag_ids", tagId}}}}

	res, err := userCollection.UpdateOne(context.TODO(), filter, update)
	if  err != nil {
		logrus.Error("Error updating user document: ", err)
		return false
	}
	if res.MatchedCount != 1 {
		logrus.Error("Error find user document (res.MatchedCount != 1): ", err)
		return false
	}
	return true
}

func (m *ManagerStruct) GetFittingUsers(user model.FullUserData) (results []model.FullUserData, ok bool) {
	database := m.Conn.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	nowIs := time.Now().Unix() * 1000
	minStamp := nowIs - int64(user.MinAge * yearInMilisecs)
	maxStamp := nowIs - int64(user.MaxAge * yearInMilisecs)

	//logrus.Infof("Now  = %17v", nowIs)
	//logrus.Infof("User = %17v", user.BirthDate)
	//logrus.Infof("Max  = %17v", maxStamp)
	//logrus.Infof("Min  = %17v", minStamp)

	user.BannedUserIds = append(user.BannedUserIds, user.Id)

	filter := bson.M{
		"id": bson.M{"$nin": user.BannedUserIds},
		"$and": bson.A{bson.D{{"birth_date", bson.D{{"$gte", maxStamp}}}}, bson.D{{"birth_date", bson.D{{"$lte", minStamp}}}}},
	}

	if user.LookFor == "male" || user.LookFor == "female" {
		filter["gender"] = user.LookFor
	}

	user.ConvertFromDbCoords()
	logrus.Infof("user before: %v", user)

	if user.UseLocation &&
		user.GeoPosition.Lat != 0 && user.GeoPosition.Lon != 0 && user.MaxDist != 0 {
		filter["position"] = bson.M{
			"$near": bson.M{
				"$geometry": bson.M{
					"type": "Point",
					"coordinates": bson.A{user.GeoPosition.Lon, user.GeoPosition.Lat},
				},
				"$maxDistance": user.MaxDist * 1000,
				"$minDistance": 0,
			},
		}
	} else {
		filter["country"] = user.Country
		filter["city"] = user.City
	}

	logrus.Infof("Full strangers filter: %v", filter)
	cur, err := userCollection.Find(context.Background(), filter)
	if  err != nil {
		logrus.Error("Error finding user document: ", err)
		return nil, false
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		container := model.FullUserData{}
		err := cur.Decode(&container)
		if err != nil {
			logrus.Error("Error decoding user: ", err)
		}
		results = append(results, container)
	}
	if err := cur.Err(); err != nil {
		logrus.Error("Error in mongo cursor: ", err)
	}

	return results, true
}

func (m *ManagerStruct) SaveLooked(lookedId, lookerId string) bool {
	database := m.Conn.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	filter := bson.D{{"id", lookedId}}
	update := bson.D{{"$addToSet", bson.D{{"looked_by", lookerId}}}}
	opts := options.Update()

	_, err := userCollection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		logrus.Errorf("Error pushing looked_by: %v", err)
		return false
	}
	return true
}

func (m *ManagerStruct) SaveLiked(likedId, likerId string) bool {
	database := m.Conn.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	filter := bson.D{{"id", likedId}}
	update := bson.D{{"$addToSet", bson.D{{"liked_by", likerId}}}}
	opts := options.Update()

	_, err := userCollection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		logrus.Errorf("Error pushing liked_by: %v", err)
		return false
	}
	return true
}

func (m *ManagerStruct) SaveMatch(matched1Id, matched2Id string) bool {
	return m.makeMatchForAccount(matched1Id, matched2Id) && m.makeMatchForAccount(matched2Id, matched1Id)
}

func (m *ManagerStruct) DeleteInteraction(acc model.LoginData, pairId string) bool {
	database := m.Conn.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	filterPair := bson.D{{"id", pairId}}
	filterUser := bson.D{{"id", acc.Id	}}
	updatePair := bson.D{{"$pull", bson.D{{"liked_by", acc.Id}}},
		{"$pull", bson.D{{"matches", acc.Id}}}}
	updateUser := bson.D{{"$pull", bson.D{{"matches", pairId}}}}
	opts := options.Update()

	_, err := userCollection.UpdateOne(context.TODO(), filterPair, updatePair, opts)
	if err != nil {
		logrus.Errorf("Error deleting interactions for pair: %v", err)
		return false
	}
	_, err = userCollection.UpdateOne(context.TODO(), filterUser, updateUser, opts)
	if err != nil {
		logrus.Errorf("Error deleting interactions for user: %v", err)
		return false
	}
	return true
}

func (m *ManagerStruct) GetPreviousInteractions(acc model.LoginData, actionType string) (result []string, err error) {
	database := m.Conn.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	var container []struct {
		Id string `bson:"id"`
	}

	var filter bson.D

	if actionType == "like" {
		filter = bson.D{{"liked_by", acc.Id}}
	} else if actionType == "look" {
		filter = bson.D{{"looked_by", acc.Id}}
	} else {
		logrus.Errorf("Unknown action type to get interactions: %v", actionType)
		return nil, fmt.Errorf("unknown action type to get interactions: %v", actionType)
	}

	projection := bson.M{"id": 1, "_id": 0}
	opts := options.Find().SetProjection(projection)

	cursor, err := userCollection.Find(context.TODO(), filter, opts)
	if err != nil {
		logrus.Errorf("Error pushing liked_by: %v", err)
		return nil, fmt.Errorf("error quering interactions: %v", actionType)
	}

	if err := cursor.All(context.TODO(), &container); err != nil {
		logrus.Errorf("Error getting liked users: %v", err)
		return nil, fmt.Errorf("error reading interactions: %v", actionType)
	}
	logrus.Infof("got type '%v' container: %v", actionType, container)
	result = make([]string, len(container))
	for i, item := range container {
		result[i] = item.Id
	}
	return result, nil
}

func (m *ManagerStruct) makeMatchForAccount(userId, matchedId string) bool {
	user := model.FullUserData{}
	database := m.Conn.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	filter := bson.D{{"id", userId}}
	update := bson.D{{"$addToSet", bson.D{{"matched", matchedId}}}, {"$pull", bson.D{{"liked_by", matchedId}}}}
	opts := options.FindOneAndUpdate()

	err := userCollection.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&user)
	if err != nil {
		logrus.Errorf("Error pushing liked_by: %v", err)
		return false
	}
	return true
}

func (m *ManagerStruct) GetUserImages(id string) []string {
	database := m.Conn.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	container := struct {
		Images		[]string	`bson:"images"`
	}{}

	opts := options.FindOne().SetProjection(bson.M{"images": 1})
	filter := bson.M{"id": id}
	err := userCollection.FindOne(context.Background(),filter,opts).Decode(&container)
	if  err != nil {
		logrus.Error("Error finding user document: ", err)
		return container.Images
	} else {
		logrus.Info("Got user images: ", container)
	}

	return container.Images
}
