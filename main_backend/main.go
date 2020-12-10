package main

import (
	"os"

	"backend/db/userMetaDataStorage"
	"backend/db/userFullDataStorage"
	"backend/server"
)

func main() {
	port := os.Getenv("BACKEND_PORT")

	userMetaDataStorage.Init()
	userFullDataStorage.Init()

	defer userMetaDataStorage.Manager.CloseConnection()
	defer userFullDataStorage.Manager.CloseConnection()

	server.StartServer(port)
}
