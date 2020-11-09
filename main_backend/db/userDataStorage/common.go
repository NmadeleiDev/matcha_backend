package userDataStorage

import (
	"backend/types"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"time"
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
		log.Error("Env is empty", user, password, addr)
	}

	connStr := fmt.Sprintf("mongodb://%v:%v@%v", user, password, addr)
	log.Info("Connecting to mongo: ", connStr)
	opts := options.Client().ApplyURI(connStr)
	m.Conn, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal("Error getting mongo client: ", err)
	}
	if err != nil {
		log.Fatal("Error connecting to mongo: ", err)
	}
	if err = m.Conn.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal("Error pinging: ", err)
	}
	log.Info("Connected.")
}

func (m *ManagerStruct) CloseConnection() {
	if err := m.Conn.Disconnect(context.TODO()); err != nil {
		log.Error("Error closing mongo: ", err)
	}
}

func (m *ManagerStruct) AddTagToUserTags(user types.LoginData, tagId int64) bool {
	database := m.Conn.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	filter := bson.M{"id": user.Id}
	update := bson.D{{"$addToSet", bson.D{{"tag_ids", tagId}}}}

	res, err := userCollection.UpdateOne(context.TODO(), filter, update)
	if  err != nil {
		log.Error("Error updating user document: ", err)
		return false
	}
	if res.MatchedCount != 1 {
		log.Error("Error find user document (res.MatchedCount != 1): ", err)
		return false
	}
	return true
}

func (m *ManagerStruct) DeleteTagFromUserTags(user types.LoginData, tagId int64) bool {
	database := m.Conn.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	filter := bson.M{"id": user.Id}
	update := bson.D{{"$pull", bson.D{{"tag_ids", tagId}}}}

	res, err := userCollection.UpdateOne(context.TODO(), filter, update)
	if  err != nil {
		log.Error("Error updating user document: ", err)
		return false
	}
	if res.MatchedCount != 1 {
		log.Error("Error find user document (res.MatchedCount != 1): ", err)
		return false
	}
	return true
}

func (m *ManagerStruct) GetFittingUsers(user types.FullUserData) (results []types.FullUserData, ok bool) {
	database := m.Conn.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	var filter bson.D

	nowIs := time.Now().Unix() * 1000
	minStamp := nowIs - int64(user.MinAge * yearInMilisecs)
	maxStamp := nowIs - int64(user.MaxAge * yearInMilisecs)

	log.Infof("Now  = %17v", nowIs)
	log.Infof("User = %17v", user.BirthDate)
	log.Infof("Max  = %17v", maxStamp)
	log.Infof("Min  = %17v", minStamp)

	if user.LookFor == "male" || user.LookFor == "female" {
		filter = bson.D{
			{"id", bson.M{"$nin": user.BannedUserIds}},
			{"gender", user.LookFor},
			{"country", user.Country},
			{"city", user.City},
			{"$and",  bson.A{bson.D{{"birth_date", bson.D{{"$gte", maxStamp}}}}, bson.D{{"birth_date", bson.D{{"$lte", minStamp}}}}}}}
	} else {
		filter = bson.D{
			{"id", bson.M{"$nin": user.BannedUserIds}},
			{"country", user.Country},
			{"city", user.City},
			{"$and",  bson.A{bson.D{{"birth_date", bson.D{{"$gte", maxStamp}}}}, bson.D{{"birth_date", bson.D{{"$lte", minStamp}}}}}}}
	}
	cur, err := userCollection.Find(context.Background(), filter)
	if  err != nil {
		log.Error("Error finding user document: ", err)
		return nil, false
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		container := types.FullUserData{}
		err := cur.Decode(&container)
		if err != nil {
			log.Error("Error decoding user: ", err)
		}
		results = append(results, container)
	}
	if err := cur.Err(); err != nil {
		log.Error("Error in mongo cursor: ", err)
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
		log.Errorf("Error pushing looked_by: %v", err)
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
		log.Errorf("Error pushing liked_by: %v", err)
		return false
	}
	return true
}

func (m *ManagerStruct) SaveMatch(matched1Id, matched2Id string) bool {
	return m.makeMatchForAccount(matched1Id, matched2Id) && m.makeMatchForAccount(matched2Id, matched1Id)
}

func (m *ManagerStruct) DeleteInteraction(acc types.LoginData, pairId string) bool {
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
		log.Errorf("Error deleting interactions for pair: %v", err)
		return false
	}
	_, err = userCollection.UpdateOne(context.TODO(), filterUser, updateUser, opts)
	if err != nil {
		log.Errorf("Error deleting interactions for user: %v", err)
		return false
	}
	return true
}

func (m *ManagerStruct) GetPreviousInteractions(acc types.LoginData, actionType string) (result []string, err error) {
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
		log.Errorf("Unknown action type to get interactions: %v", actionType)
		return nil, fmt.Errorf("unknown action type to get interactions: %v", actionType)
	}

	projection := bson.M{"id": 1, "_id": 0}
	opts := options.Find().SetProjection(projection)

	cursor, err := userCollection.Find(context.TODO(), filter, opts)
	if err != nil {
		log.Errorf("Error pushing liked_by: %v", err)
		return nil, fmt.Errorf("error quering interactions: %v", actionType)
	}

	if err := cursor.All(context.TODO(), &container); err != nil {
		log.Errorf("Error getting liked users: %v", err)
		return nil, fmt.Errorf("error reading interactions: %v", actionType)
	}
	log.Infof("got type '%v' container: %v", actionType, container)
	result = make([]string, len(container))
	for i, item := range container {
		result[i] = item.Id
	}
	return result, nil
}

func (m *ManagerStruct) makeMatchForAccount(userId, matchedId string) bool {
	user := types.FullUserData{}
	database := m.Conn.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	filter := bson.D{{"id", userId}}
	update := bson.D{{"$addToSet", bson.D{{"matched", matchedId}}}, {"$pull", bson.D{{"liked_by", matchedId}}}}
	opts := options.FindOneAndUpdate()

	err := userCollection.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&user)
	if err != nil {
		log.Errorf("Error pushing liked_by: %v", err)
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
		log.Error("Error finding user document: ", err)
		return container.Images
	} else {
		log.Info("Got user images: ", container)
	}

	return container.Images
}