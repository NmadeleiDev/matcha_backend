package types

type SocketMessage struct {
	MessageType		int		`json:"message_type"`
	To				[]string	`json:"to"`
	Payload			[]byte		`json:"payload"`
}

