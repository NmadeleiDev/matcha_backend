package postgres

import (
	"context"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	pgx "github.com/jackc/pgx/v4"
)

type DBHandler struct {
	DB *pgx.Conn
}

func InitTables() {
	dbConn := getConnection()
	defer dbConn.DB.Close(context.Background())

	query := "CREATE SCHEMA IF NOT EXISTS matcha"
	_, err := dbConn.DB.Exec(context.Background(), query)
	if err != nil {
		log.Println(" InitTables error. Err: " + err.Error())
	}

	query = "CREATE TABLE IF NOT EXISTS matcha.users_data (id SERIAL PRIMARY KEY, username VARCHAR(64) NOT NULL, email VARCHAR(255) NOT NULL UNIQUE, password VARCHAR(255) NOT NULL, unique_key VARCHAR(255) NOT NULL, acc_state SMALLINT DEFAULT 0)"
	_, err = dbConn.DB.Exec(context.Background(), query)
	if err != nil {
		log.Println(" InitTables error. Err: " + err.Error())
	}
	fmt.Println("Tables initialized")
}

func InsertUserData(userName, email, password, uniqueKey string) bool {
	dbConn := getConnection()
	defer dbConn.DB.Close(context.Background())

	query := "INSERT INTO matcha.users_data(username, email, password, unique_key) VALUES($1, $2, $3, $4) ON CONFLICT DO NOTHING"
	result, err := dbConn.DB.Exec(context.Background(), query, userName, email, password, uniqueKey)
	if err != nil {
		log.Println(" InsertUserData error. Err: " + err.Error())
	}
	if result.RowsAffected() < 1 {
		return false
	}
	return true
}

func AuthUser(email, password string) bool {
	var databasePassword string

	dbConn := getConnection()
	defer dbConn.DB.Close(context.Background())

	query := "SELECT password FROM matcha.users_data WHERE email=$1"
	resultRow := dbConn.DB.QueryRow(context.Background(), query, email)
	resultRow.Scan(&databasePassword)
	err := bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func VerifyAccount(vkey string) bool  {
	dbConn := getConnection()
	defer dbConn.DB.Close(context.Background())

	query := "UPDATE matcha.users_data SET acc_state=1 WHERE unique_key=$1 AND acc_state=0"
	result, err := dbConn.DB.Exec(context.Background(), query, vkey)
	if err != nil {
		log.Println(" VerifyAccount error. Err: " + err.Error())
	}
	if result.RowsAffected() < 1 {
		return false
	}
	return true
}

func getConnection() *DBHandler {
	dsn := getDSN()
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		log.Println("Config parse error. Err: " + err.Error())
	}
	db, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Println("Database is down. Err: " + err.Error())
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

	dsn := "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + db + "?sslmode=disable"
	return dsn
}
