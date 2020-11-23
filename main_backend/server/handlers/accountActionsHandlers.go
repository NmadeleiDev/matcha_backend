package handlers

import (
	"backend/db/structuredDataStorage"
	"backend/db/userDataStorage"
	"backend/emails"
	"backend/model"
	"backend/utils"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		userData, ok := utils.UnmarshalHttpBodyToUserData(w, r)
		if !ok {
			return
		}

		authKey, ok := structuredDataStorage.Manager.CreateUser(userData)
		if !ok {
			utils.SendFailResponse(w, "User already exists")
			return
		}

		if !userDataStorage.Manager.CreateUser(*userData) {
			utils.SendFailResponse(w, "failed to create user")
			return
		}

		emails.Manager.SendVerificationKey(userData.Email, authKey)
		utils.SendSuccessResponse(w)
	}
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		loginData, ok := utils.UnmarshalHttpBodyToLoginData(w, r)
		if !ok {
			return
		}
		if structuredDataStorage.Manager.LoginUser(loginData) {
			userData, err := utils.GetFullUserData(*loginData, false)
			if err != nil {
				log.Error("Failed to get user data")
				utils.SendFailResponse(w,"Failed to get user data")
			} else {
				userData.Id = loginData.Id
				utils.RefreshRequestSessionKeyCookie(w, *loginData)
				utils.SendDataResponse(w, userData)
			}
		} else {
			utils.SendFailResponse(w,"incorrect user data")
		}
	}
}

func VerifyAccountHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	if len(key) == 0 {
		log.Error("Can't read request path for verify: ", r.URL.String())
		utils.SendFailResponse(w, "Not correct key")
		return
	}

	newSessionKey, ok := structuredDataStorage.Manager.VerifyUserAccountState(key)
	if ok {
		utils.SetCookie(w, "session_id", newSessionKey)
		login, err := structuredDataStorage.Manager.GetUserLoginDataBySession(newSessionKey)
		if err != nil {
			utils.SendFailResponse(w,"Failed to get user data")
		} else {
			data, err := userDataStorage.Manager.GetFullUserData(login, "public")
			if err != nil {
				utils.SendFailResponse(w, "Failed to get user data")
			} else {
				utils.SendDataResponse(w, data)
			}
		}
	} else {
		utils.SendFailResponse(w, "failed to verify user")
	}
}

func SignOutHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodDelete {
		session := utils.GetCookieValue(r,"session_id")
		ok := structuredDataStorage.Manager.UpdateSessionKey(session, "")
		if ok {
			utils.SendSuccessResponse(w)
		} else {
			utils.SendFailResponse(w, "incorrect session")
		}
	}
}

func ManageOwnAccountHandler(w http.ResponseWriter, r *http.Request) {
	session := utils.GetCookieValue(r,"session_id")
	if len(session) == 0 {
		utils.SendFailResponse(w, "incorrect session id")
		return
	}

	loginData, err := structuredDataStorage.Manager.GetUserLoginDataBySession(session)
	if err != nil {
		utils.SendFailResponse(w, "incorrect session id")
		return
	}
	if r.Method == http.MethodGet {
		userData, err := utils.GetFullUserData(loginData, false)
		if err != nil {
			utils.SendFailResponse(w,"Failed to get user data")
		} else {
			utils.SendDataResponse(w, userData)
			return
		}
	} else if r.Method == http.MethodPut {
		userData, ok := utils.UnmarshalHttpBodyToUserData(w, r)
		if !ok {
			return
		}

		loginData, err := structuredDataStorage.Manager.GetUserLoginDataBySession(utils.GetCookieValue(r, "session_id"))
		if err != nil {
			log.Error("Can't get user is: ", err)
			utils.SendFailResponse(w, "Session cookie not present")
			return
		}

		userData.Id = loginData.Id
		if !userDataStorage.Manager.UpdateUser(*userData) {
			utils.SendFailResponse(w, "failed to update user")
			return
		}
		utils.SendSuccessResponse(w)
	} else if r.Method == http.MethodDelete {
		userData, err := userDataStorage.Manager.GetFullUserData(loginData, "public")
		if err != nil {
			utils.SendFailResponse(w, fmt.Sprintf("Failed to get user data: %v", err))
			return
		}
		for _, tagId := range userData.TagIds {
			_ = structuredDataStorage.Manager.DecrTagById(tagId)
		}
		if err := userDataStorage.Manager.DeleteAccount(loginData); err != nil {
			utils.SendFailResponse(w, fmt.Sprintf("Failed to delete user data: %v", err))
			return
		}
		if err := structuredDataStorage.Manager.DeleteAccount(loginData); err != nil {
			utils.SendFailResponse(w, fmt.Sprintf("Failed to delete user account metadata: %v", err))
			return
		}
		emails.Manager.SendGoodbyeMessage(userData.Email)
		log.Infof("Deleted %v", userData.Email)

		utils.SendSuccessResponse(w)
	}
}

func UserTagsHandler(w http.ResponseWriter, r *http.Request) {
	loginData := utils.AuthUserBySessionId(w, r)
	if loginData == nil {
		return
	}
	if r.Method == http.MethodPut {
		tags, ok := utils.UnmarshalHttpBodyToTags(w, r)
		if !ok {
			return
		}
		failedTags := make([]string, 0, len(tags.Tags))
		for _, tag := range tags.Tags {
			id, err := structuredDataStorage.Manager.IncOrInsertTag(tag)
			if err != nil {
				failedTags = append(failedTags, tag)
			} else {
				userDataStorage.Manager.AddTagToUserTags(*loginData, id)
			}
		}

		if len(failedTags) == 0 {
			utils.SendSuccessResponse(w)
		} else {
			utils.SendFailResponse(w, fmt.Sprintf("Failed to save tags: %v", failedTags))
		}
	} else if r.Method == http.MethodDelete {
		tags, ok := utils.UnmarshalHttpBodyToTags(w, r)
		if !ok {
			return
		}
		failedTags := make([]string, 0, len(tags.Tags))
		for _, tag := range tags.Tags {
			id, err := structuredDataStorage.Manager.DecrTagByValue(tag)
			ok := userDataStorage.Manager.DeleteTagFromUserTags(*loginData, id)
			if err != nil || !ok {
				failedTags = append(failedTags, tag)
			}
		}
		go structuredDataStorage.Manager.ClearUnmentionedTags()

		if len(failedTags) == 0 {
			utils.SendSuccessResponse(w)
		} else {
			utils.SendFailResponse(w, fmt.Sprintf("Failed to delete tags: %v", failedTags))
		}
	} else if r.Method == http.MethodGet {
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

		utils.SendDataResponse(w, structuredDataStorage.Manager.GetAllTags(limit, offset))
	}
}

func GetUserDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		var userData interface{}
		var err error

		id := mux.Vars(r)["id"]
		isShortData := r.URL.Query().Get("full") == "false"

		if utils.AuthUserBySessionId(w, r) == nil {
			return
		}

		if isShortData {
			userData, err = userDataStorage.Manager.GetShortUserData(model.LoginData{Id: id})
		} else {
			userData, err = utils.GetFullUserData(model.LoginData{Id: id}, true)
		}
		if err != nil {
			utils.SendFailResponse(w,"Failed to get user data")
		} else {
			utils.SendDataResponse(w, userData)
			return
		}
	}
}

func GetOwnActionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		action := mux.Vars(r)["action"]

		data := utils.AuthUserBySessionId(w, r)
		if data == nil {
			return
		}

		actions, err := userDataStorage.Manager.GetPreviousInteractions(*data, action)
		if err != nil {
			utils.SendFailResponse(w, err.Error())
			return
		}
		utils.SendDataResponse(w, actions)
	}
}

func GetUserOwnImagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data := utils.AuthUserBySessionId(w, r)
		if data == nil {
			return
		}
		utils.SendDataResponse(w, userDataStorage.Manager.GetUserImages(data.Id))
	}
}