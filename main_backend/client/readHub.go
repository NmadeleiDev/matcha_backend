package client

import (
	"backend/db/structuredDataStorage"
	"backend/types"
	"encoding/json"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	NewChatMeta = 1
	InsertMessage = 1
	UpdateMessage = 2
	DeleteMessage = 3

	CreateChat = 100

	MessageSent = 1
	MessageDelivered = 2
	MessageRead = 3
)

func	(client *Client) ReadHub() {
	defer func() {
		if err := client.Connection.Close(); err != nil {
			log.Error("Error closing connection in read: ", err)
		}
	}()

	client.Connection.SetReadLimit(maxMessageSize)
	if err := client.Connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Error("Error setting write deadline: ", err)
		return
	}
	client.Connection.SetPongHandler(func(string) error { client.Connection.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var messageStruct types.SocketMessage
		_, message, err := client.Connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error: %v", err)
			}
			break
		}
		if err := json.Unmarshal(message, &messageStruct); err != nil {
			log.Error("Error unmarshal message: ", err, " Message: ", message)
			continue
		}

		if messageStruct.MessageType >= 0 && messageStruct.MessageType < 100 {
			var userMessage types.Message

			if err := json.Unmarshal(messageStruct.Payload, &userMessage); err != nil {
				log.Errorf("Error unmarshal message: %v", err)
			}
			switch messageStruct.MessageType {
			case InsertMessage:
				structuredDataStorage.Manager.SaveMessage(userMessage)
			case UpdateMessage:
				structuredDataStorage.Manager.UpdateMessageState(userMessage.Id, userMessage.State)
			case DeleteMessage:
				structuredDataStorage.Manager.DeleteMessage(userMessage.Id)
			default:
				log.Warnf("Unknown message type: %v", messageStruct)
				continue
			}
			SendMessageToClient(messageStruct)
		} else if messageStruct.MessageType >= 100 && messageStruct.MessageType < 200 {
			// Create chat
			//SendMessageToClient(messageStruct)
		}
		log.Info("Got message in read hub: ", messageStruct)

	}
}