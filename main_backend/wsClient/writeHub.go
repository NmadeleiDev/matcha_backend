package wsClient

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"time"
)

func	(client *Client) WriteHub() {

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		if err := client.Connection.Close(); err != nil {
			log.Error("Error closing connection in write: ", err)
		}
		close(client.ReadMessageChan)
		Clients[client.Id] = nil
	}()

	for {
		select {
		case message := <- client.ReadMessageChan:
			err := client.Connection.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Error("Error setting write deadline: ", err)
				return
			}
			w, err := client.Connection.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Error("Error getting next writer: ", err)
				return
			}
			messageBytes, err := json.Marshal(message)
			if err != nil {
				log.Error("Error marshal message: ", err)
				return
			}
			_, err = w.Write(messageBytes)
			if err != nil {
				log.Error("Error writing message to ws: ", err)
				return
			}
			if err = w.Close(); err != nil {
				log.Error("Error closing writer: ", err)
			}
			log.Infof("Message %v sent.", message)
		case <-ticker.C:
			err := client.Connection.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Error("Error setting write deadline: ", err)
				return
			}
			if err := client.Connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Error("Error writing ticker message to ws: ", err)
				return
			}
		}
	}
}

