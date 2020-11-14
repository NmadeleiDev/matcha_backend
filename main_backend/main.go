package main

import (
	"backend/wsClient"
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

	wsClient.Clients = make(map[string]*wsClient.Client, 10000)
	server.StartServer(port)
}
