package userMetaDataStorage

import (
	"net/http"

	"backend/model"
	"backend/utils"

	"github.com/sirupsen/logrus"
)

func (m *ManagerStruct) RefreshRequestSessionKeyCookie(w http.ResponseWriter, user model.LoginData) bool {
	sessionKey, err := Manager.IssueUserSessionKey(user)
	if err != nil {
		if err.Error() == "STATE" {
			utils.SendFailResponse(w, "User account not verified")
			return false
		}
		logrus.Errorf("Error refreshing cookie: %v", err)
		utils.SendFailResponse(w, "incorrect user data")
		return false
	}

	utils.SetCookieForDay(w, "session_id", sessionKey)
	return true
}

func (m *ManagerStruct) AuthUserBySessionId(w http.ResponseWriter, r *http.Request) *model.LoginData {
	session := utils.GetCookieValue(r, "session_id")
	if len(session) != 32 {
		utils.SendFailResponseWithCode(w, "incorrect user data", http.StatusUnauthorized)
		return nil
	}

	user, err := Manager.GetUserLoginDataBySession(session)
	if err != nil {
		logrus.Error("Failed to get user data by session")
		utils.SendFailResponseWithCode(w, "incorrect user data", http.StatusUnauthorized)
		return nil
	}
	return &user
}
