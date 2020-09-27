package types

type UserDataStorage interface {
	MakeConnection()
	CloseConnection()

	CreateUser(user UserData) bool
	GetUserData(user LoginData) (UserData, error)
	UpdateUser(user UserData) bool
	GetFittingUsers(user UserData) (results []UserData, ok bool)
	SaveLiked(likedId, likerId string) bool
	SaveLooked(lookedId, lookerId string) bool
	SaveMatch(matched1Id, matched2Id string) bool
	GetUserImages(id string) []string
}

type StructuredDataStorage interface {
	MakeConnection()
	CloseConnection()

	CreateUser(userData *UserData) (string, bool)
	LoginUser(loginData *LoginData) bool
	SetSessionKeyById(sessionKey string, id string) bool
	GetUserEmailBySession(sessionKey string) (user LoginData, err error)
	GetUserLoginDataBySession(sessionKey string) (user LoginData, err error)
	VerifyUserAccountState(key string) (string, bool)
	UpdateSessionKey(old, new string) bool
	IssueUserSessionKey(user LoginData) (string, error)

	SaveMessage(message Message) bool
	UpdateMessageState(messageId string, state int) bool
	DeleteMessage(id string) bool
}
