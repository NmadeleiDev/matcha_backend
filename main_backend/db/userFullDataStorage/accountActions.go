package userFullDataStorage

import (
	"context"
	"fmt"
	"time"

	"backend/model"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *ManagerStruct) CreateUser(user model.FullUserData) bool {
	database := m.Conn.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	position := user.GeoPosition.ConvertToMongoCoords()

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
		{"use_location", user.UseLocation},
		{"last_online", time.Now().Unix()},
	}

	_, err := userCollection.InsertOne(context.TODO(), userDocument)
	if err != nil {
		log.Error("Error creating user in mongo: ", err)
		return false
	}
	return true
}

func (m *ManagerStruct) GetFullUserData(user model.LoginData, isPublic bool) (model.FullUserData, error) {
	opts := options.FindOne()

	database := m.Conn.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	//log.Info("UserId: ", user)
	filter := bson.M{"id": user.Id}
	container := model.FullUserData{}
	projection := bson.M{}

	if isPublic {
		projection["email"] = 0
		projection["max_dist"] = 0
		projection["look_for"] = 0
		projection["min_age"] = 0
		projection["max_age"] = 0
		projection["use_location"] = 0
		projection["liked_by"] = 0
		projection["looked_by"] = 0
		projection["matches"] = 0
	}

	err := userCollection.FindOne(context.Background(), filter, opts).Decode(&container)
	if err != nil {
		log.Error("Error finding user document: ", err)
		return model.FullUserData{}, err
	}

	return container, nil
}

func (m *ManagerStruct) GetUserDataWithCustomProjection(user model.LoginData, projectFields []string, doInclude bool) model.FullUserData {
	opts := options.FindOne()

	database := m.Conn.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	filter := bson.M{"id": user.Id}
	container := model.FullUserData{}
	projection := bson.M{}

	projVal := 0
	if doInclude {
		projVal = 1
	}
	for _, val := range projectFields {
		projection[val] = projVal
	}

	err := userCollection.FindOne(context.TODO(), filter, opts).Decode(&container)
	if err != nil {
		log.Error("Error finding user document for custom projection: ", err)
		return model.FullUserData{}
	}

	return container
}

func (m *ManagerStruct) FindUserAndUpdateGeo(user model.LoginData, geo model.Coordinates) (model.FullUserData, error) {
	var update bson.M

	opts := options.FindOneAndUpdate()

	database := m.Conn.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	filter := bson.M{"id": user.Id}
	update = bson.M{"$set": bson.D{{"position",geo.ConvertToMongoCoords()}}}
	container := model.FullUserData{}

	projection := bson.M{"banned_user_ids": 0}
	opts.SetProjection(projection).SetReturnDocument(options.After)

	err := userCollection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&container)
	if err != nil {
		log.Error("Error finding user document: ", err)
		return model.FullUserData{}, err
	}

	return container, nil
}

func (m *ManagerStruct) GetShortUserData(user model.LoginData) (model.ShortUserData, error) {

	database := m.Conn.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	log.Info("UserId: ", user)
	filter := bson.M{"id": user.Id}
	container := model.ShortUserData{}
	err := userCollection.FindOne(context.Background(), filter).Decode(&container)
	if err != nil {
		log.Error("Error finding user document: ", err)
		return model.ShortUserData{}, err
	} else {
		log.Infof("Got user document: %v; avatar = %v", container, container.Avatar)
	}

	return container, nil
}

func (m *ManagerStruct) UpdateUser(user model.FullUserData) bool {
	database := m.Conn.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	//user.ConvertToDbCoords()

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
		{"$set", bson.D{{"use_location", user.UseLocation}}},
	}

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

func (m *ManagerStruct) SetNewEmail(userId, email string) error {
	database := m.Conn.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	//user.ConvertToDbCoords()

	filter := bson.M{"id": userId}
	update := bson.D{
		{"$set", bson.D{{"email", email}}},
	}

	res, err := userCollection.UpdateOne(context.TODO(), filter, update)
	if  err != nil {
		log.Error("Error updating user email (mongo): ", err)
		return err
	}
	if res.MatchedCount != 1 {
		log.Error("Error find user document (res.MatchedCount != 1): ", err)
		return fmt.Errorf("user not found")
	}
	return nil
}

func (m *ManagerStruct) AddUserIdToBanned(acc model.LoginData, bannedId string) bool {
	database := m.Conn.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	filter := bson.M{"id": acc.Id}
	update := bson.D{{"$addToSet", bson.M{"banned_user_ids": bannedId}}}
	res, err := userCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Error("Error updating user document (ban user): ", err)
		return false
	}
	log.Infof("Ban res: %v", res)
	return true
}

func (m *ManagerStruct) GetUserBannedList(acc model.LoginData) (result []string, err error) {
	database := m.Conn.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	container := struct{
		BannedUserIds	[]string	`bson:"banned_user_ids"`
	}{}

	filter := bson.M{"id": acc.Id}
	projection := bson.M{"banned_user_ids": 1}
	opts := options.FindOne().SetProjection(projection)

	err = userCollection.FindOne(context.Background(), filter, opts).Decode(&container)
	if err != nil {
		log.Error("Error finding user document: ", err)
		return nil, err
	}

	return container.BannedUserIds, nil
}

func (m *ManagerStruct) RemoveUserIdFromBanned(acc model.LoginData, bannedId string) bool {
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

func (m *ManagerStruct) DeleteAccount(acc model.LoginData) error {
	database := m.Conn.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	filterDelete := bson.M{"id": acc.Id}

	if result, err := userCollection.DeleteOne(context.TODO(), filterDelete); err != nil {
		return err
	} else if result.DeletedCount != 1 {
		return fmt.Errorf("user to delete not found")
	}

	return nil
}

func (m *ManagerStruct) DeleteAccountRecordsFromOtherUsers(acc model.LoginData) error {
	database := m.Conn.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	update := bson.M{"$pull": bson.M{"liked_by": acc.Id, "looked_by": acc.Id, "matches": acc.Id}}

	if _, err := userCollection.UpdateMany(context.TODO(), bson.M{}, update); err != nil {
		return err
	}

	return nil
}

func (m *ManagerStruct) GetUserEmailByUsername(acc model.LoginData) string {
	database := m.Conn.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)
	cont := model.LoginData{}

	filter := bson.M{"username": acc.Username}
	projection := bson.M{"email": 1}
	opts := options.FindOne().SetProjection(projection)
	err := userCollection.FindOne(context.TODO(), filter, opts).Decode(&cont)
	if err != nil {
		log.Errorf("Error finding user document: %v; filter: %v", err, filter)
		return ""
	}
	return cont.Email
}
