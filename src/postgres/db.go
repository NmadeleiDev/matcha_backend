package postgres

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type DBHandler struct {
	DB *sql.DB
}

func InitTables() {
	dbConn := getConnection()
	defer dbConn.DB.Close()

	query := "CREATE TABLE IF NOT EXISTS users_data (id SERIAL PRIMARY KEY, username VARCHAR(64) NOT NULL, password VARCHAR(255) NOT NULL, unique_key SMALLINT DEFAULT 0)"
	_, err := dbConn.DB.Query(query)
	if err != nil {
		log.Fatal(" InsertUserData error. Err: " + err.Error())
	}
}

func InsertUserData(userName, password, uniqueKey string) {
	dbConn := getConnection()
	defer dbConn.DB.Close()

	query := "INSERT INTO users_data(username, password, unique_key) VALUES(?, ?, ?)"
	_, err := dbConn.DB.Query(query, userName, password, uniqueKey)
	if err != nil {
		log.Fatal(" InsertUserData error. Err: " + err.Error())
	}
}

func AuthUser(email, password string) bool {
	var databasePassword string

	dbConn := getConnection()
	defer dbConn.DB.Close()

	query := "SELECT FROM users_data (password) WHERE email=?"
	resultRow := dbConn.DB.QueryRow(query, email)
	resultRow.Scan(&databasePassword)
	err := bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func getConnection() *DBHandler{
	dsn := getDSN()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Database is down. Err: " + err.Error())
	}
	connection := &DBHandler{db}
	return connection
}

func getDSN() string {
	user, _ := os.LookupEnv("POSTGRES_USER")
	password, _ := os.LookupEnv("POSTGRES_PASSWORD")
	port, _ := os.LookupEnv("POSTGRES_PORT")
	host, _ := os.LookupEnv("POSTGRES_HOST")
	db, _ := os.LookupEnv("POSTGRES_DB")

	dsn := "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + db + "?sslmode=verify-full"
	println("DSN: ", dsn)
	return dsn
}