package main

import (
	"os"

	"backend/db/structuredDataStorage"
	"backend/db/userDataStorage"
	"backend/server"
)

func main() {
	port := os.Getenv("BACKEND_PORT")

	structuredDataStorage.Init()
	userDataStorage.Init()

	defer structuredDataStorage.Manager.CloseConnection()
	defer userDataStorage.Manager.CloseConnection()

	server.StartServer(port)
}
