package userMetaDataStorage

import (
	"backend/model"
	"crypto/md5"
	"database/sql"
	_ "database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
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
	passwordResetTable = "user_data_schema.password_reset"
	tagsTable = "user_data_schema.tags"
	messagesTable = "message_data_schema.messages"
	hashCost      = 14
)

type ManagerStruct struct {
	Conn	*sqlx.DB
}

func (m *ManagerStruct) MakeConnection() {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	db := os.Getenv("POSTGRES_DB")

	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", user, password, host, port, db)
	conn, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}
	err = conn.Ping()
	if err != nil {
		log.Fatal("Error pinging postgres: ", err)
	}
	m.Conn = conn
}

func (m *ManagerStruct) CloseConnection() {
	if err := m.Conn.Close(); err != nil {
		log.Error("Error closing postgres: ", err)
	}
}

func (m *ManagerStruct) InitTables() {
	query := `create schema if not exists ` + strings.Split(userDataTable, ".")[0]

	if _, err := m.Conn.Exec(query); err != nil {
		log.Fatal("Error creating schema: ", err)
	}

	query = `create schema if not exists ` + strings.Split(messagesTable, ".")[0]

	if _, err := m.Conn.Exec(query); err != nil {
		log.Fatal("Error creating schema: ", err)
	}

	query = `create table if not exists ` + userDataTable + `
(
    id            varchar(256)       not null
        constraint users_pk
            primary key,
    password      varchar(255) not null,
    email varchar(128) unique,
    session_key   varchar(128) default NULL::character varying,
	acc_state		integer default 2
)`
	if _, err := m.Conn.Exec(query); err != nil {
		log.Fatal("Error creating table: ", err)
	}

	query = `create table if not exists ` + passwordResetTable + `
(
    user_id            varchar(256)       not null
        constraint password_reset_pk
            primary key,
    key      varchar(128) not null,
	state		integer default 0,
	created_at		timestamp default now()::timestamp
)`
	if _, err := m.Conn.Exec(query); err != nil {
		log.Fatal("Error creating table: ", err)
	}

	query = `create table if not exists ` + tagsTable + `
(
    id            serial
		constraint tags_pk
			primary key,
    value      varchar(255) not null,
    hashing varchar(32) unique not null,
    times_mentioned   integer default 1
)`
	if _, err := m.Conn.Exec(query); err != nil {
		log.Fatal("Error creating table: ", err)
	}

	query = `create table if not exists ` + messagesTable + `
(
    id            varchar(256)       not null
        constraint users_pk
            primary key,
    sender      varchar(128) not null,
    recipient      varchar(128) not null,
    state integer default 0,
	date integer not null,
    text   varchar(1024) default ''
)`
	if _, err := m.Conn.Exec(query); err != nil {
		log.Fatal("Error creating table: ", err)
	}
}

func (m *ManagerStruct) SaveMessage(message model.Message) bool {
	query := `
INSERT INTO ` + messagesTable + ` (sender, recipient, date, text) 
VALUES (?, ?, ?, ?)`

	if _, err := m.Conn.Exec(query, message.Sender, message.Recipient, message.Date, message.Text); err != nil {
		log.Errorf("Error saving message: %v", err)
		return false
	}
	return true
}

func (m *ManagerStruct) UpdateMessageState(messageId string, state int) bool {
	query := `
UPDATE ` + messagesTable + ` 
SET state=? 
WHERE id=?`

	if _, err := m.Conn.Exec(query, state, messageId); err != nil {
		log.Errorf("Error updating message: %v", err)
		return false
	}
	return true
}

func (m *ManagerStruct) DeleteMessage(id string) bool {
	query := `
DELETE FROM ` + messagesTable + ` 
WHERE id=?`

	if _, err := m.Conn.Exec(query, id); err != nil {
		log.Errorf("Error updating message: %v", err)
		return false
	}
	return true
}

func (m *ManagerStruct) SetSessionKeyById(sessionKey string, id string) bool {
	query := `
UPDATE ` + userDataTable + ` 
SET session_key=$1
WHERE id=$2`

	if _, err := m.Conn.Exec(query, sessionKey, id); err != nil {
		log.Error("Error setting session key: ", err)
		log.Error("Key: ", sessionKey, " ID: ", id)
		return false
	}
	return true
}

func (m *ManagerStruct) GetUserEmailBySession(sessionKey string) (user model.LoginData, err error) {

	query := `
SELECT email
FROM ` + userDataTable + ` 
WHERE session_key=$1`

	row := m.Conn.QueryRow(query, sessionKey)
	err = row.Scan(&user.Email)
	return user, err
}

func (m *ManagerStruct) GetUserIdByEmail(email string) (id string, err error) {
	query := `
SELECT id
FROM ` + userDataTable + ` 
WHERE email=$1`

	row := m.Conn.QueryRow(query, email)
	err = row.Scan(&id)
	return id, err
}

func (m *ManagerStruct) GetUserLoginDataBySession(sessionKey string) (user model.LoginData, err error) {
	query := `
SELECT id
FROM ` + userDataTable + ` 
WHERE session_key=$1`

	row := m.Conn.QueryRow(query, sessionKey)
	err = row.Scan(&user.Id)
	return user, err
}

func (m *ManagerStruct) UpdateSessionKey(old, new string) bool {
	query := `
UPDATE ` + userDataTable + ` 
SET session_key=$1
WHERE session_key=$2`

	if _, err := m.Conn.Exec(query, new, old); err != nil {
		log.Error("Error updating session key: ", err)
		return false
	}
	return true
}

func (m *ManagerStruct) VerifyUserAccountState(key string) (string, bool) {
	sessionKeyBytes := md5.Sum([]byte(time.Now().String() + key + strconv.Itoa(rand.Int())))
	newSessionKey := fmt.Sprintf("%x", sessionKeyBytes)

	query := `
UPDATE ` + userDataTable + ` 
SET acc_state=0, session_key=$2 
WHERE session_key=$1`

	if _, err := m.Conn.Exec(query, key, newSessionKey); err != nil {
		log.Errorf("Error verifying acc state: %v", err)
		return "", false
	}
	return newSessionKey, true
}

func (m *ManagerStruct) IssueUserSessionKey(user model.LoginData) (string, error) {
	var truePassword string
	var state int

	query := `
SELECT password, acc_state FROM ` + userDataTable + ` 
WHERE id = $1`

	row := m.Conn.QueryRow(query, user.Id)
	if err := row.Scan(&truePassword, &state); err != nil {
		log.Error("Error getting user info: ", err)
		return "", err
	}
	//if state == 2 { для разработки закоментил
	//	return "", errors.New("STATE")
	//}
	if err := bcrypt.CompareHashAndPassword([]byte(truePassword), []byte(user.Password)); err != nil {
		log.Error("Error verifying password: ", err)
		return "", err
	}
	sessionKeyBytes := md5.Sum([]byte(time.Now().String() + user.Email + strconv.Itoa(rand.Int())))
	sessionKey := fmt.Sprintf("%x", sessionKeyBytes)

	if m.SetSessionKeyById(sessionKey, user.Id) {
		return sessionKey, nil
	} else {
		return "", errors.New("error updating session key")
	}
}

func (m *ManagerStruct) IncOrInsertTag(tag string) (id int64, err error) {
	tagBytes := md5.Sum([]byte(tag))
	tagHash := fmt.Sprintf("%x", tagBytes)

	query := `INSERT INTO ` + tagsTable + ` AS tt (value, hashing) 
		VALUES ($1, $2)
		ON CONFLICT (hashing) DO UPDATE SET times_mentioned=tt.times_mentioned + 1
		RETURNING id`

	if err := m.Conn.QueryRow(query, tag, tagHash).Scan(&id); err != nil {
		log.Errorf("Error upserting tag: %v", err)
		return 0, err
	}
	return id, nil
}

func (m *ManagerStruct) DecrTagByValue(tag string) (id int64, err error) {
	tagBytes := md5.Sum([]byte(tag))
	tagHash := fmt.Sprintf("%x", tagBytes)

	query := `UPDATE ` + tagsTable + ` AS tt SET times_mentioned=tt.times_mentioned - 1 
		WHERE hashing=$1 RETURNING id`

	if err := m.Conn.QueryRow(query, tagHash).Scan(&id); err != nil {
		log.Errorf("Error decrementing tag: %v", err)
		return 0, err
	}
	return id, nil
}

func (m *ManagerStruct) DecrTagById(tagId int64) (err error) {
	query := `UPDATE ` + tagsTable + ` AS tt SET times_mentioned=tt.times_mentioned - 1 
		WHERE id=$1`

	if _, err := m.Conn.Exec(query, tagId); err != nil {
		log.Errorf("Error decrementing tag: %v", err)
		return err
	}
	return nil
}

func (m *ManagerStruct) ClearUnmentionedTags() {
	query := `DELETE FROM ` + tagsTable + ` WHERE times_mentioned <= 0`

	if _, err := m.Conn.Exec(query); err != nil {
		log.Errorf("Error cleaning tags: %v", err)
	}
}

func (m *ManagerStruct) GetTagsById(ids []int64) (tags []string) {
	query := `SELECT value FROM ` + tagsTable + ` 
		WHERE id=ANY($1)`

	rows, err := m.Conn.Query(query, pq.Array(ids))
	if err != nil {
		log.Errorf("Error finding tags by ids: %v", err)
		return nil
	}
	for rows.Next() {
		cont := ""
		if err := rows.Scan(&cont); err != nil {
			log.Errorf("Error scanning tag: %v", err)
			continue
		}
		tags = append(tags, cont)
	}
	return tags
}

func (m *ManagerStruct) GetAllTags(limit, offset int) (tags []string) {
	query := `SELECT value FROM ` + tagsTable + ` 
		ORDER BY times_mentioned DESC`

	var rows *sql.Rows
	var err error

	if limit > 0 && offset > 0 {
		query += ` LIMIT $1 OFFSET $2`
		rows, err = m.Conn.Query(query, limit, offset)
	} else {
		rows, err = m.Conn.Query(query)
	}

	if err != nil {
		log.Errorf("Error getting all tags: %v", err)
		return nil
	}
	for rows.Next() {
		cont := ""
		if err := rows.Scan(&cont); err != nil {
			log.Errorf("Error scanning tag: %v", err)
			continue
		}
		tags = append(tags, cont)
	}
	return tags
}
