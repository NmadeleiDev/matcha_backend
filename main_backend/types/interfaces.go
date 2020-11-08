package types

type UserDataStorage interface {
	MakeConnection()
	CloseConnection()

	CreateUser(user FullUserData) bool
	GetFullUserData(user LoginData, variant string) (FullUserData, error) // variant: private/public/full
	GetShortUserData(user LoginData) (ShortUserData, error)
	UpdateUser(user FullUserData) bool
	DeleteAccount(acc LoginData) error

	AddTagToUserTags(user LoginData, tagId int64) bool
	DeleteTagFromUserTags(user LoginData, tagId int64) bool
	GetFittingUsers(user FullUserData) (results []FullUserData, ok bool)
	GetPreviousInteractions(acc LoginData, actionType string) (result []string, err error)
	SaveLiked(likedId, likerId string) bool
	SaveLooked(lookedId, lookerId string) bool
	SaveMatch(matched1Id, matched2Id string) bool
	DeleteInteraction(acc LoginData, pairId string) bool
	GetUserImages(id string) []string

	AddUserIdToBanned(acc LoginData, bannedId string) bool
	GetUserBannedList(acc LoginData) (result []string, err error)
	RemoveUserIdFromBanned(acc LoginData, bannedId string) bool
}

type StructuredDataStorage interface {
	MakeConnection()
	CloseConnection()

	CreateUser(userData *FullUserData) (string, bool)
	LoginUser(loginData *LoginData) bool
	DeleteAccount(loginData LoginData) error

	SetSessionKeyById(sessionKey string, id string) bool
	GetUserEmailBySession(sessionKey string) (user LoginData, err error)
	GetUserLoginDataBySession(sessionKey string) (user LoginData, err error)
	VerifyUserAccountState(key string) (string, bool)
	UpdateSessionKey(old, new string) bool
	IssueUserSessionKey(user LoginData) (string, error)

	IncOrInsertTag(tag string) (id int64, err error)
	GetTagsById(ids []int64) (tags []string)
	DecrTagByValue(tag string) (id int64, err error)
	DecrTagById(tagId int64) (err error)
	GetAllTags(limit, offset int) (tags []string)
	ClearUnmentionedTags()

	SaveMessage(message Message) bool
	UpdateMessageState(messageId string, state int) bool
	DeleteMessage(id string) bool
}

type EmailService interface {
	SendVerificationKey(to, key string)
	SendGoodbyeMessage(to string)
}
