package types

type UserDataStorage interface {
	MakeConnection()
	CloseConnection()

	CreateUser(user FullUserData) bool
	GetFullUserData(user LoginData, isPublic bool) (FullUserData, error)
	GetShortUserData(user LoginData) (ShortUserData, error)
	UpdateUser(user FullUserData) bool
	AddTagToUserTags(user LoginData, tagId int64) bool
	DeleteTagFromUserTags(user LoginData, tagId int64) bool
	GetFittingUsers(user FullUserData) (results []FullUserData, ok bool)
	SaveLiked(likedId, likerId string) bool
	SaveLooked(lookedId, lookerId string) bool
	SaveMatch(matched1Id, matched2Id string) bool
	GetUserImages(id string) []string
}

type StructuredDataStorage interface {
	MakeConnection()
	CloseConnection()

	CreateUser(userData *FullUserData) (string, bool)
	LoginUser(loginData *LoginData) bool
	SetSessionKeyById(sessionKey string, id string) bool
	GetUserEmailBySession(sessionKey string) (user LoginData, err error)
	GetUserLoginDataBySession(sessionKey string) (user LoginData, err error)
	VerifyUserAccountState(key string) (string, bool)
	UpdateSessionKey(old, new string) bool
	IssueUserSessionKey(user LoginData) (string, error)

	IncOrInsertTag(tag string) (id int64, err error)
	GetTagsById(ids []int64) (tags []string)
	DecrTag(tag string) (id int64, err error)
	GetAllTags() (tags []string)
	ClearUnmentionedTags()

	SaveMessage(message Message) bool
	UpdateMessageState(messageId string, state int) bool
	DeleteMessage(id string) bool
}
