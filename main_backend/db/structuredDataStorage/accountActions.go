package structuredDataStorage

import (
	"math/rand"
	"strconv"
	"time"

	"backend/hash"
	"backend/model"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func (m *ManagerStruct) CreateUser(userData *model.FullUserData) (string, bool) {

	query := `
INSERT INTO ` + userDataTable + `(email, password, id, session_key)
VALUES ($1, $2, $3, $4)` // здесь session_key создается, чтобы авторизовать почту юзера

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userData.Password), hashCost)
	if err != nil {
		log.Error("Error hashing password", err)
		return "", false
	}
	rawId := userData.Email + time.Now().String() + strconv.Itoa(rand.Int())

	userData.Id = hash.CalculateSha256(rawId)
	key := hash.CalculateSha256(userData.Id + strconv.Itoa(rand.Int()))
	_, err = m.Conn.Exec(query, userData.Email, passwordHash, userData.Id, key)
	if err != nil {
		log.Error("Error creating user: ", err)
		return "", false
	}
	return key, true
}

func (m *ManagerStruct) LoginUser(loginData *model.LoginData) bool {
	var truePassword string

	query := `
SELECT id, password FROM ` + userDataTable + ` 
WHERE email = $1`

	row := m.Conn.QueryRow(query, loginData.Email)
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

func (m *ManagerStruct) DeleteAccount(loginData model.LoginData) error {
	query := `DELETE FROM ` + userDataTable + ` WHERE id=$1`
	if _, err := m.Conn.Exec(query, loginData.Id); err != nil {
		return err
	}
	return nil
}

