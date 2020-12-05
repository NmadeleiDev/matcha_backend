package main

import (
	"backend/db/userDataStorage"
	"backend/db/structuredDataStorage"
	"backend/server"
	"os"
)

func main() {
	port := os.Getenv("BACKEND_PORT")

	structuredDataStorage.Init()
	userDataStorage.Init()

	defer structuredDataStorage.Manager.CloseConnection()
	defer userDataStorage.Manager.CloseConnection()

	server.StartServer(port)
}
