package handlers

import (
	"backend/client"
	"backend/types"
	"backend/utils"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
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
	id, ok := utils.IdentifyUserBySession(r)
	if !ok {
		utils.SendFailResponse(w, "Unauthorized request")
		return
	}

	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("Error establishing ws connection: ", err)
		utils.SendFailResponse(w, err.Error())
		return
	}

	user := &types.FullUserData{Id: id}

	clientStruct := client.RegisterNewClient(connection, user)
	go clientStruct.ReadHub()
	go clientStruct.WriteHub()
}



