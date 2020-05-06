package main

import (
	"backend/db/mongodb"
	"backend/db/postgres"
	"backend/server"
	"os"
)

func main() {
	port := os.Getenv("BACKEND_PORT")

	defer postgres.CloseConnection()
	defer mongodb.CloseConnection()

	postgres.MakeConnection()
	mongodb.MakeConnection()

	postgres.InitTables()

	server.StartServer(port)
}
