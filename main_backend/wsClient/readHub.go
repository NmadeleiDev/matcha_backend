package wsClient

import (
	"encoding/json"
	"fmt"
	"time"

	"backend/model"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

const (
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
		_, message, err := client.Connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error: %v", err)
			}
			break
		}

		messageType := GetWsMessageType(message)
		if messageType < 1 {
			log.Errorf("Invalid message type: %v", messageType)
			client.ReadMessageChan <- model.SocketMessage{MessageType: 0, Payload: fmt.Sprintf("Invalid message type: %v", messageType)}
			continue
		}

		if messageType >= 0 && messageType < 100 {
			var userMessage struct{
				Payload 	model.Message 	`json:"payload"`
			}

			if err := json.Unmarshal(message, &userMessage); err != nil {
				log.Errorf("Error unmarshal message: %v", err)
			}
			switch messageType {
			case InsertMessage:
				GetManager().SendMessageToChat(userMessage.Payload)
				//structuredDataStorage.Manager.SaveMessage(userMessage.Payload)
			case UpdateMessage:
				GetManager().UpdateMessageToChat(userMessage.Payload)
				//structuredDataStorage.Manager.UpdateMessageState(userMessage.Payload.Id, userMessage.Payload.State)
			case DeleteMessage:
				GetManager().DeleteMessageFromChat(userMessage.Payload)
				//structuredDataStorage.Manager.DeleteMessage(userMessage.Payload.Id)
			default:
				log.Warnf("Unknown message type: %v", messageType)
				continue
			}
		} else if messageType >= 100 && messageType < 200 {
			var newChat struct{
				Payload 	model.Chat 	`json:"payload"`
			}

			if err := json.Unmarshal(message, &newChat); err != nil {
				log.Errorf("Error unmarshal message: %v", err)
			}
			GetManager().CreateChat(newChat.Payload)
		}
	}
}