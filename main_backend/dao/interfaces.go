package dao

import "backend/model"

type UserFullDataStorage interface {
	MakeConnection()
	CloseConnection()

	CreateUser(user model.FullUserData) bool
	FindUserAndUpdateGeo(user model.LoginData, geo model.Coordinates) (model.FullUserData, error)
	GetFullUserData(user model.LoginData, variant string) (model.FullUserData, error) // variant: private/public/full
	GetShortUserData(user model.LoginData) (model.ShortUserData, error)
	GetUserDataWithCustomProjection(user model.LoginData, projectFields []string, doInclude bool) model.FullUserData
	UpdateUser(user model.FullUserData) bool
	DeleteAccount(acc model.LoginData) error
	DeleteAccountRecordsFromOtherUsers(acc model.LoginData) error

	AddTagToUserTags(user model.LoginData, tagId int64) bool
	DeleteTagFromUserTags(user model.LoginData, tagId int64) bool
	GetFittingUsers(user model.FullUserData) (results []model.FullUserData, ok bool)
	GetPreviousInteractions(acc model.LoginData, actionType string) (result []string, err error)
	SaveLiked(likedId, likerId string) bool
	SaveLooked(lookedId, lookerId string) bool
	SaveMatch(matched1Id, matched2Id string) bool
	DeleteInteraction(acc model.LoginData, pairId string) bool
	GetUserImages(id string) []string

	AddUserIdToBanned(acc model.LoginData, bannedId string) bool
	GetUserBannedList(acc model.LoginData) (result []string, err error)
	RemoveUserIdFromBanned(acc model.LoginData, bannedId string) bool

	CreateLocationIndex()
}

type UserMetaDataStorage interface {
	MakeConnection()
	CloseConnection()

	CreateUser(userData *model.FullUserData) (string, bool)
	LoginUser(loginData *model.LoginData) bool
	DeleteAccount(loginData model.LoginData) error

	SetSessionKeyById(sessionKey string, id string) bool
	GetUserEmailBySession(sessionKey string) (user model.LoginData, err error)
	GetUserIdByEmail(email string) (id string, err error)
	GetUserLoginDataBySession(sessionKey string) (user model.LoginData, err error)
	VerifyUserAccountState(key string) (string, bool)
	UpdateSessionKey(old, new string) bool
	IssueUserSessionKey(user model.LoginData) (string, error)

	CreateResetPasswordRecord(userId, key string) error
	SetNextStepResetKey(oldKey, newKey string) error
	GetAccountIdByResetKey(key string) (id string, err error)
	SetNewPasswordForAccount(accountId string, newPassword string) error

	IncOrInsertTag(tag string) (id int64, err error)
	GetTagsById(ids []int64) (tags []string)
	DecrTagByValue(tag string) (id int64, err error)
	DecrTagById(tagId int64) (err error)
	GetAllTags(limit, offset int) (tags []string)
	ClearUnmentionedTags()

	SaveMessage(message model.Message) bool
	UpdateMessageState(messageId string, state int) bool
	DeleteMessage(id string) bool
}

type NotificationsBroker interface {
	MakeConnection()
	CloseConnection()

	PublishMessage(channelId, mType, originId string)
}

type WsDataManager interface {
	CreateChat(chat model.Chat) string
	FindChat(chatId string) *model.Chat
	GetUserChats(userId string) []*model.Chat
	ConnectToChat(chatId string)
	SendMessageToChat(message model.Message)
	UpdateMessageInChat(message model.Message)
	DeleteMessageFromChat(message model.Message)
	AddUserToChat(userId string, chat model.Chat)
}

type EmailService interface {
	SendVerificationKey(to, key string)
	SendGoodbyeMessage(to string)
	SendPasswordResetEmail(to, key string)
}
