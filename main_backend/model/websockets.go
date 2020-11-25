package model

type SocketMessage struct {
	MessageType int    `json:"messageType"`
	Payload     interface{} `json:"payload"`
}
