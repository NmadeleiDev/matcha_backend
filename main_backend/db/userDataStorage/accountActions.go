package userDataStorage

import (
	"backend/types"
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
)

func (m *ManagerStruct) CreateUser(user types.FullUserData) bool {
	database := m.Conn.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	position := MongoCoords{Type: "point", Coordinates: []float64{user.GeoPosition.Lon, user.GeoPosition.Lat}}

	userDocument := bson.D{
		{"id", user.Id},
		{"username", user.Username},
		{"email", user.Email},
		{"birth_date", user.BirthDate},
		{"gender", user.Gender},
		{"phone", user.Phone},
		{"country", user.Country},
		{"city", user.City},
		{"max_dist", user.MaxDist},
		{"look_for", user.LookFor},
		{"min_age", user.MinAge},
		{"max_age", user.MaxAge},
		{"looked_by", []string{}},
		{"liked_by", []string{}},
		{"matches", []string{}},
		{"position", position},
	}

	_, err := userCollection.InsertOne(context.TODO(), userDocument)
	if err != nil {
		log.Error("Error creating user in mongo: ", err)
		return false
	}
	return true
}

func (m *ManagerStruct) GetFullUserData(user types.LoginData, variant string) (types.FullUserData, error) {

	var opts *options.FindOneOptions

	database := m.Conn.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	log.Info("UserId: ", user)
	filter := bson.M{"id": user.Id}
	container := types.FullUserData{}

	if variant != "full" {
		projection := bson.M{"banned_user_ids": 0}
		opts = options.FindOne().SetProjection(projection)
	}

	err := userCollection.FindOne(context.Background(), filter, opts).Decode(&container)
	if err != nil {
		log.Error("Error finding user document: ", err)
		return types.FullUserData{}, err
	}

	if variant == "public" {
		container.LikedBy = nil
		container.LookedBy = nil
		container.Matches = nil
	}

	if len(container.Avatar) == 0 && len(container.Images) > 0 {
		container.Avatar = container.Images[rand.Intn(len(container.Images))]
	}

	return container, nil
}

func (m *ManagerStruct) GetShortUserData(user types.LoginData) (types.ShortUserData, error) {

	database := m.Conn.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	log.Info("UserId: ", user)
	filter := bson.M{"id": user.Id}
	container := types.ShortUserData{}
	err := userCollection.FindOne(context.Background(), filter).Decode(&container)
	if err != nil {
		log.Error("Error finding user document: ", err)
		return types.ShortUserData{}, err
	} else {
		log.Infof("Got user document: %v; avatar = %v", container, container.Avatar)
	}

	if len(container.Avatar) == 0 && len(container.Images) > 0 {
		container.Avatar = container.Images[rand.Intn(len(container.Images))]
	}

	return container, nil
}

func (m *ManagerStruct) UpdateUser(user types.FullUserData) bool {
	database := m.Conn.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	position := MongoCoords{Type: "point", Coordinates: []float64{user.GeoPosition.Lon, user.GeoPosition.Lat}}

	filter := bson.M{"id": user.Id}
	update := bson.D{
		{"$set", bson.D{{"username", user.Username}}},
		{"$set", bson.D{{"name", user.Name}}},
		{"$set", bson.D{{"surname", user.Surname}}},
		{"$set", bson.D{{"birth_date", user.BirthDate}}},
		{"$set", bson.D{{"gender", user.Gender}}},
		{"$set", bson.D{{"phone", user.Phone}}},
		{"$set", bson.D{{"country", user.Country}}},
		{"$set", bson.D{{"city", user.City}}},
		{"$set", bson.D{{"bio", user.Bio}}},
		{"$set", bson.D{{"max_dist", user.MaxDist}}},
		{"$set", bson.D{{"look_for", user.LookFor}}},
		{"$set", bson.D{{"min_age", user.MinAge}}},
		{"$set", bson.D{{"max_age", user.MaxAge}}},
		{"$set", bson.D{{"position", position}}}}

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

func (m *ManagerStruct) AddUserIdToBanned(acc types.LoginData, bannedId string) bool {
	database := m.Conn.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	filter := bson.M{"id": acc.Id}
	update := bson.D{{"$addToSet", bson.D{{"banned_user_ids", bannedId}}}}
	_, err := userCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Error("Error updating user document (ban user): ", err)
		return false
	}
	return true
}

func (m *ManagerStruct) GetUserBannedList(acc types.LoginData) (result []string, err error) {
	database := m.Conn.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	filter := bson.M{"id": acc.Id}
	projection := bson.M{"banned_user_ids": 1}
	opts := options.FindOne().SetProjection(projection)

	err = userCollection.FindOne(context.Background(), filter, opts).Decode(&result)
	if err != nil {
		log.Error("Error finding user document: ", err)
		return nil, err
	}

	return result, nil
}

func (m *ManagerStruct) RemoveUserIdFromBanned(acc types.LoginData, bannedId string) bool {
	database := m.Conn.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	filter := bson.M{"id": acc.Id}
	update := bson.D{{"$pull", bson.D{{"banned_user_ids", bannedId}}}}
	_, err := userCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Error("Error updating user document (ban user): ", err)
		return false
	}
	return true
}

func (m *ManagerStruct) DeleteAccount(acc types.LoginData) error {
	database := m.Conn.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	filter := bson.M{"id": acc.Id}

	if _, err := userCollection.DeleteOne(context.TODO(), filter); err != nil {
		return err
	}
	return nil
}
