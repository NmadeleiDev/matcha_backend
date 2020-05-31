package mongodb

import (
	"backend/structs"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
)

const userDataCollection = "users"
const mainDBName = "matcha"
var client *mongo.Client

func MakeConnection() {
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
	client, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal("Error getting client mongo: ", err)
	}
	if err != nil {
		log.Fatal("Error connecting to mongo: ", err)
	}
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal("Error pinging: ", err)
	}
	log.Info("Connected.")
}

func CloseConnection() {
	if err := client.Disconnect(context.TODO()); err != nil {
		log.Error("Error closing mongo: ", err)
	}
}

func CreateUser(user structs.UserData) bool {
	database := client.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	user.Password = ""
	user.LookedBy = []string{}
	user.LikedBy = []string{}
	user.Matches = []string{}

	_, err := userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Error("Error creating user in mongo: ", err)
		return false
	}
	return true
}

func GetUserData(user structs.LoginData) (structs.UserData, error) {

	database := client.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	log.Info("User: ", user)
	filter := bson.M{"id": user.Id}
	container := structs.UserData{}
	err := userCollection.FindOne(context.Background(),filter).Decode(&container)
	if  err != nil {
		log.Error("Error finding user document: ", err)
		return structs.UserData{}, err
	} else {
		log.Info("Got user document: ", container)
	}

	return container, nil
}

func UpdateUser(user structs.UserData) bool {
	database := client.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	filter := bson.M{"id": user.Id}
	update := bson.D{{"$set", bson.D{{"username", user.Username}}},
		{"$set", bson.D{{"age", user.Age}}},
		{"$set", bson.D{{"gender", user.Gender}}},
		{"$set", bson.D{{"phone", user.Phone}}},
		{"$set", bson.D{{"country", user.Country}}},
		{"$set", bson.D{{"city", user.City}}},
		{"$set", bson.D{{"max_dist", user.MaxDist}}},
		{"$set", bson.D{{"look_for", user.LookFor}}},
		{"$set", bson.D{{"min_age", user.MinAge}}},
		{"$set", bson.D{{"max_age", user.MaxAge}}}}

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

func GetFittingUsers(user structs.UserData) (results []structs.UserData, ok bool) {
	database := client.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	filter := bson.D{{"gender", user.LookFor}, {"country", user.Country}, {"city", user.City}, {"$and",  bson.A{bson.D{{"age", bson.D{{"$gte", user.MinAge}}}}, bson.D{{"age", bson.D{{"$lte", user.MaxAge}}}}}}}
	//log.Info("Filter: ", filter)
	//log.Info("User: ", user)
	cur, err := userCollection.Find(context.Background(),filter)
	if  err != nil {
		log.Error("Error finding user document: ", err)
		return nil, false
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		container := structs.UserData{}
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

func SaveLooked(lookedId, lookerId string) bool {
	database := client.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	filter := bson.D{{"id", lookedId}}
	update := bson.D{{"$push", bson.D{{"looked_by", lookerId}}}}
	opts := options.Update()

	_, err := userCollection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		log.Errorf("Error pushing looked_by: %v", err)
		return false
	}
	return true
}

func SaveLiked(likedId, likerId string) bool {
	database := client.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	filter := bson.D{{"id", likedId}}
	update := bson.D{{"$push", bson.D{{"liked_by", likerId}}}}
	opts := options.Update()

	_, err := userCollection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		log.Errorf("Error pushing liked_by: %v", err)
		return false
	}
	return true
}

func SaveMatch(matched1Id, matched2Id string) bool {

	return makeMatchForAccount(matched1Id, matched2Id) && makeMatchForAccount(matched2Id, matched1Id)
}

func makeMatchForAccount(userId, matchedId string) bool {
	user := structs.UserData{}
	database := client.Database(mainDBName)
	userCollection := database.Collection(userDataCollection)

	filter := bson.D{{"id", userId}}
	update := bson.D{{"$push", bson.D{{"matched", matchedId}}}, {"$pull", bson.D{{"liked_by", matchedId}}}}
	opts := options.FindOneAndUpdate()

	err := userCollection.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&user)
	if err != nil {
		log.Errorf("Error pushing liked_by: %v", err)
		return false
	}
	return true
}