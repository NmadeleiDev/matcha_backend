package main

import (
	"os"

	"backend/db/notificationsBroker"
	"backend/db/userMetaDataStorage"
	"backend/db/userFullDataStorage"
	"backend/server"
)

func main() {
	port := os.Getenv("BACKEND_PORT")

	userMetaDataStorage.Init()
	userFullDataStorage.Init()
	notificationsBroker.GetManager().MakeConnection()

	defer userMetaDataStorage.Manager.CloseConnection()
	defer userFullDataStorage.Manager.CloseConnection()
	defer notificationsBroker.GetManager().CloseConnection()

	server.StartServer(port)
}
