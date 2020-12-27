package model

type Notification struct {
	Type string		`json:"type"` // look/like/match
	User	string	`json:"user"` // id fo user, how originated the action (i.e. looked profile)
}
