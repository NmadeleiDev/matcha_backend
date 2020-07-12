package main

import (
	"backend/db/userDataStorage"
	"backend/db/structuredDataStorage"
	"backend/server"
	"os"
)

func main() {
	port := os.Getenv("BACKEND_PORT")

	defer structuredDataStorage.Manager.CloseConnection()
	defer userDataStorage.Manager.CloseConnection()

	structuredDataStorage.Init()
	userDataStorage.Init()

	server.StartServer(port)
}
