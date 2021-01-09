package userMetaDataStorage

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"backend/hashing"
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

	userData.Id = hashing.CalculateSha256(rawId)
	key := hashing.CalculateSha256(userData.Id + strconv.Itoa(rand.Int()))
	_, err = m.Conn.Exec(query, userData.Email, passwordHash, userData.Id, key)
	if err != nil {
		log.Error("Error creating user: ", err)
		return "", false
	}
	return key, true
}

func (m *ManagerStruct) LoginUser(loginData *model.LoginData) error {
	var truePassword string
	var accState int

	query := `
SELECT id, password, acc_state FROM ` + userDataTable + ` 
WHERE email=$1`

	row := m.Conn.QueryRow(query, loginData.Email)
	if err := row.Scan(&loginData.Id, &truePassword, &accState); err != nil {
		log.Error("Error getting user info: ", err)
		return err
	}
	if accState != 2 {
		return fmt.Errorf("STATE")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(truePassword), []byte(loginData.Password)); err != nil {
		log.Error("Error verifying password: ", err)
		return err
	}
	return nil
}

func (m *ManagerStruct) CreateResetEmailRecord(userId, email, key string) bool {
	query := `INSERT INTO ` + emailResetTable + ` (user_id, email, key) VALUES ($1, $2, $3)
	ON CONFLICT (user_id) DO UPDATE SET key=$3, email=$2`
	if _, err := m.Conn.Exec(query, userId, email, key); err != nil {
		log.Errorf("Error creating email reset record: %v; key: %v", err, key)
		return false
	}
	return true
}

func (m *ManagerStruct) GetResetEmailRecord(key string) (userId, email string, err error) {
	query := `DELETE FROM ` + emailResetTable + ` WHERE key=$1 RETURNING user_id, email`
	if err := m.Conn.QueryRow(query, key).Scan(&userId, &email); err != nil {
		return "", "", err
	}
	return userId, email, nil
}

func (m *ManagerStruct) SetNewEmail(userId, email string) error {
	query := `UPDATE ` + userDataTable + ` SET email=$1, acc_state=2 WHERE id=$2`
	if _, err := m.Conn.Exec(query, email, userId); err != nil {
		log.Errorf("Error setting new email: %v", err)
		return err
	}
	return nil
}

func (m *ManagerStruct) DeleteAccount(loginData model.LoginData) error {
	query := `DELETE FROM ` + userDataTable + ` WHERE id=$1`
	if _, err := m.Conn.Exec(query, loginData.Id); err != nil {
		return err
	}
	return nil
}

func (m *ManagerStruct) CreateResetPasswordRecord(userId, key string) error {
	query := `INSERT INTO ` + passwordResetTable + ` (user_id, key) VALUES ($1, $2)
	ON CONFLICT (user_id) DO UPDATE SET key=$2, state=DEFAULT`
	if _, err := m.Conn.Exec(query, userId, key); err != nil {
		return err
	}
	return nil
}

func (m *ManagerStruct) SetNextStepResetKey(oldKey, newKey string) error {
	query := `UPDATE ` + passwordResetTable + ` SET key=$1, state=1 WHERE key=$2 AND state=0`
	if _, err := m.Conn.Exec(query, newKey, oldKey); err != nil {
		log.Errorf("Error setting next step reset password key: %v", err)
		return err
	}
	return nil
}

func (m *ManagerStruct) SetNewPasswordForAccount(accountId string, newPassword string) error {
	query := `UPDATE ` + userDataTable + ` SET password=$1, session_key='' WHERE id=$2`
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), hashCost)
	if err != nil {
		log.Error("Error hashing password", err)
		return err
	}

	if _, err := m.Conn.Exec(query, passwordHash, accountId); err != nil {
		return err
	}

	return nil
}

func (m *ManagerStruct) GetAccountIdByResetKey(key string) (id string, err error) {
	query := `DELETE FROM ` + passwordResetTable + ` WHERE key=$1 AND state=1 RETURNING user_id`
	if err := m.Conn.QueryRow(query, key).Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}


