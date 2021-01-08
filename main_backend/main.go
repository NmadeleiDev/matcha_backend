package main

import (
	"os"

	"backend/db/realtimeDataDb"
	"backend/db/userMetaDataStorage"
	"backend/db/userFullDataStorage"
	"backend/server"
)

func main() {
	port := os.Getenv("BACKEND_PORT")

	userMetaDataStorage.Init()
	userFullDataStorage.Init()
	realtimeDataDb.GetManager().MakeConnection()

	defer userMetaDataStorage.Manager.CloseConnection()
	defer userFullDataStorage.Manager.CloseConnection()
	defer realtimeDataDb.GetManager().CloseConnection()

	server.StartServer(port)
}
