package backend

import (
	"backend/db/postgres"
	"backend/server"
	"os"
)

func	main() {
	port := os.Getenv("BACKEND_PORT")

	defer postgres.CloseConnection()
	postgres.MakeConnection()

	server.StartServer(port)
}
