package main

import (
	"backend/client"
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

	client.Clients = make(map[string]*client.Client, 10000)
	server.StartServer(port)
}
