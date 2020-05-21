package postgres

import (
	"backend/structs"
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	userDataTable = "user_data_schema.user_data"
	hashCost      = 14
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
	err = conn.Ping()
	if err != nil {
		log.Fatal("Error pinging postgres: ", err)
	}
	connection = conn
}

func CloseConnection() {
	if err := connection.Close(); err != nil {
		log.Error("Error closing postgres connection: ", err)
	}
}

func InitTables() {
	query := `create schema if not exists ` + strings.Split(userDataTable, ".")[0]

	if _, err := connection.Exec(query); err != nil {
		log.Error("Error creating schema: ", err)
	}

	query = `create table if not exists ` + userDataTable + `
(
    id            varchar(256)       not null
        constraint users_pk
            primary key,
    password      varchar(255) not null,
    email varchar(128) unique,
    session_key   varchar(128) default NULL::character varying
)`
	if _, err := connection.Exec(query); err != nil {
		log.Error("Error creating table: ", err)
	}
}

func CreateUser(userData *structs.UserData) bool {

	query := `
INSERT INTO ` + userDataTable + `(email, password, id)
VALUES ($1, $2, $3)`

	rawId := userData.Email + time.Now().String() + strconv.Itoa(rand.Int())
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userData.Password), hashCost)
	if err != nil {
		log.Error("Error hashing password", err)
		return false
	}

	userData.Id = CalculateSha256(rawId)
	_, err = connection.Exec(query, userData.Email, passwordHash, userData.Id)
	if err != nil {
		log.Error("Error creating user: ", err)
		return false
	}
	return true
}

func LoginUser(loginData *structs.LoginData) bool {
	var truePassword string

	query := `
SELECT id, password FROM ` + userDataTable + ` 
WHERE email = $1`

	row := connection.QueryRow(query, loginData.Email)
	if err := row.Scan(&loginData.Id, &truePassword); err != nil {
		log.Error("Error getting user info: ", err)
		return false
	}
	if err := bcrypt.CompareHashAndPassword([]byte(truePassword), []byte(loginData.Password)); err != nil {
		log.Error("Error verifying password: ", err)
		return false
	}
	return true
}

func SetSessionKeyById(sessionKey string, id string) bool {
	query := `
UPDATE ` + userDataTable + ` 
SET session_key=$1
WHERE id=$2`

	if _, err := connection.Exec(query, sessionKey, id); err != nil {
		log.Error("Error setting session key: ", err)
		log.Error("Key: ", sessionKey, " ID: ", id)
		return false
	}
	return true
}

func GetUserEmailBySession(sessionKey string) (user structs.LoginData, err error) {

	query := `
SELECT email
FROM ` + userDataTable + ` 
WHERE session_key=$1`

	row := connection.QueryRow(query, sessionKey)
	err = row.Scan(&user.Email)
	return user, err
}

func GetUserIdBySession(sessionKey string) (user structs.LoginData, err error) {
	query := `
SELECT id
FROM ` + userDataTable + ` 
WHERE session_key=$1`

	row := connection.QueryRow(query, sessionKey)
	err = row.Scan(&user.Id)
	return user, err
}

func UpdateSessionKey(old, new string) bool {
	query := `
UPDATE ` + userDataTable + ` 
SET session_key=$1
WHERE session_key=$2`

	if _, err := connection.Exec(query, new, old); err != nil {
		log.Error("Error updating session key: ", err)
		return false
	}
	return true
}

func IssueUserSessionKey(user structs.LoginData) (string, error) {
	var truePassword string

	query := `
SELECT password FROM ` + userDataTable + ` 
WHERE id = $1`

	row := connection.QueryRow(query, user.Id)
	if err := row.Scan(&truePassword); err != nil {
		log.Error("Error getting user info: ", err)
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(truePassword), []byte(user.Password)); err != nil {
		log.Error("Error verifying password: ", err)
		return "", err
	}
	sessionKeyBytes := md5.Sum([]byte(time.Now().String() + user.Email + strconv.Itoa(rand.Int())))
	sessionKey := fmt.Sprintf("%x", sessionKeyBytes)

	if SetSessionKeyById(sessionKey, user.Id) {
		return sessionKey, nil
	} else {
		return "", errors.New("error updating session key")
	}
}
