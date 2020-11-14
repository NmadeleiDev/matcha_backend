package model

type SocketMessage struct {
	MessageType int    `json:"messageType"`
	ToChat      string `json:"toChat"`
	Payload     interface{} `json:"payload"`
}
