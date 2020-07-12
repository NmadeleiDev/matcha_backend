package client

import (
	"backend/types"
	"github.com/gorilla/websocket"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 30 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = pongWait * 9 / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024
)

var Clients map[string]*Client

type Client struct {
	Id				string
	IsOnline		bool
	Connection      *websocket.Conn
	ReadMessageChan chan types.SocketMessage
}

func	RegisterNewClient(connection *websocket.Conn, user *types.UserData) (client *Client) {
	client = &Client{Id: user.Id, Connection: connection, ReadMessageChan:make(chan types.SocketMessage), IsOnline: true}
	Clients[user.Id] = client
	return client
}

func	SendMessageToClient(message types.SocketMessage) {
	for _, id := range message.To {
		if Clients[id] != nil && Clients[id].IsOnline {
			Clients[id].ReadMessageChan <- message
		}
	}
}

