package postgres

import (
	"backend/structs"
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"os"
)


const (
	userDataTable = "user_data_schema.user_data"
	hashCost = 14
)

var connection *sql.DB

func MakeConnection() {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	db := os.Getenv("POSTGRES_DB")

	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", user, password, host, port, db)
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}
	connection = conn
}

func CloseConnection() {
	if err := connection.Close(); err != nil {
		log.Error("Error closing postgres connection: ", err)
	}
}

func	CreateUser(userData structs.UserData) bool {

	query := `
INSERT INTO ` + userDataTable + ` (email, phone, password, username, born_date, gender, country, city, max_dist, look_for, min_age, max_age)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userData.Password), hashCost)
	if err != nil {
		log.Error("Error hashing password", err)
		return false
	}

	_, err = connection.Exec(query, userData.Email, userData.Phone, passwordHash, userData.BornDate, userData.Gender, userData.Country, userData.City, userData.MaxDist, userData.LookFor, userData.MinAge, userData.MaxAge)
	if err != nil {
		log.Error("Error creating user: ",err)
		return false
	}
	return true
}

func	LoginUser(loginData structs.LoginData) bool {
	var truePassword string

	query := `
SELECT password FROM ` + userDataTable + ` 
WHERE email = $1`

	row := connection.QueryRow(query, loginData.Email)
	if err := row.Scan(&truePassword); err != nil {
		log.Error("Error getting user info: ", err)
		return false
	}
	if err := bcrypt.CompareHashAndPassword([]byte(truePassword), []byte(loginData.Password)); err != nil {
		log.Error("Error verifying password: ", err)
		return false
	}
	return true
}

func	GetFullUserData(loginData structs.LoginData) (userData structs.UserData, err error) {

	query := `
SELECT email, phone, username, born_date, gender, country, city, max_dist, look_for, min_age, max_age
FROM ` + userDataTable + `
WHERE email=$1`

	row := connection.QueryRow(query, loginData.Email)
	if err = row.Scan(&userData.Email, &userData.Phone, &userData.Username, &userData.BornDate, &userData.Gender, &userData.Country, &userData.City, &userData.MaxDist, &userData.LookFor, &userData.MinAge, &userData.MaxAge); err != nil {
		log.Error("Error getting user data: ", err)
		return structs.UserData{}, err
	}
	return userData, nil
}
