package notificationsBroker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"backend/dao"
	"backend/model"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

const (
	LookType         = "NEW_LOOK"
	CreatedLikeType  = "NEW_LIKE"
	DeletedLikeType  = "DELETE_LIKE"
	CreatedMatchType = "NEW_MATCH"
	DeletedMatchType = "DELETE_MATCH"
)

type ManagerStruct struct {
	client	*redis.Client
}

func (m *ManagerStruct) MakeConnection() {
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	password := os.Getenv("REDIS_PASSWORD")
	addr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	
	m.client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	logrus.Infof("Connected to redis! %v %v %v", addr, password, db)
}

func (m *ManagerStruct) CloseConnection() {
	if err := m.client.Close(); err != nil {
		logrus.Errorf("Error closing redis conn: %v", err)
	}
}

func (m *ManagerStruct) PublishMessage(channelId, mType, originId string) {
	message := model.Notification{Type: mType, User: originId}
	body, err := json.Marshal(message)
	if err != nil {
		logrus.Errorf("Error marshal message for notif: %v", err)
		return
	}
	res := m.client.Publish(context.TODO(), channelId, body)
	logrus.Infof("Published message: %v to channel %v; err = %v", message, channelId, res.Err())
}

var manager ManagerStruct

func GetManager() dao.NotificationsBroker {
	return &manager
}

