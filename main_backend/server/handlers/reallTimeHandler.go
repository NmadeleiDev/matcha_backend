package handlers

import (
	"net/http"

	"backend/db/structuredDataStorage"
	"backend/utils"
	"backend/wsClient"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	//HandshakeTimeout:  10,
	//ReadBufferSize:    1024,
	//WriteBufferSize:   1024,
	//WriteBufferPool:   nil,
	//Subprotocols:      []string{"chat"},
	//Error:             nil,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	//EnableCompression: false,
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request)  {
	log.Info("Managing websocket")
	session := r.URL.Query().Get("key")
	data, err := structuredDataStorage.Manager.GetUserLoginDataBySession(session)
	if err != nil {
		log.Errorf("Error find user: %v", err)
		return
	}

	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("Error establishing ws connection: ", err)
		utils.SendFailResponse(w, err.Error())
		return
	}


	clientStruct := wsClient.RegisterNewClient(connection, &data)
	go clientStruct.ReadHub()
	go clientStruct.WriteHub()
}



