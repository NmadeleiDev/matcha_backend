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
	filter := bson.M{"email": user.Email}
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
