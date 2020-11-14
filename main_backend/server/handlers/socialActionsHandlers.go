package handlers

import (
	"backend/db/structuredDataStorage"
	"backend/db/userDataStorage"
	"backend/model"
	"backend/utils"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func GetStrangersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		session := utils.GetCookieValue(r, "session_id")
		user, err := structuredDataStorage.Manager.GetUserLoginDataBySession(session)
		if err != nil {
			log.Error("Failed to get user data by session")
			utils.SendFailResponse(w, "incorrect user data")
			return
		}
		userData, err := userDataStorage.Manager.GetFullUserData(user, "full")
		if err != nil {
			log.Error("Failed to get user data from mongo")
			utils.SendFailResponse(w, "sorry!")
			return
		}
		strangers, ok := userDataStorage.Manager.GetFittingUsers(userData)
		if ok {
			utils.SendDataResponse(w, strangers)
		} else {
			utils.SendFailResponse(w, "Failed to load users")
		}
	}
}

func LookActionHandler(w http.ResponseWriter, r *http.Request) {
	lookedUserId, ok := utils.UnmarshalHttpBodyToLoginData(w, r)
	if !ok {
		return
	}

	loginData, err := structuredDataStorage.Manager.GetUserLoginDataBySession(utils.GetCookieValue(r, "session_id"))
	if err != nil {
		log.Error("Session Id is invalid: ", err)
		utils.SendFailResponse(w, "Your session ID is invalid, try to refresh the page.")
		return
	}

	if r.Method == http.MethodPost {
		if userDataStorage.Manager.SaveLooked(lookedUserId.Id, loginData.Id) {
			utils.SendSuccessResponse(w)
		} else {
			utils.SendFailResponse(w,"failed to save looked to db.")
		}
	} else if r.Method == http.MethodGet {
		looks, err := userDataStorage.Manager.GetPreviousInteractions(loginData, "look")
		if err != nil {
			utils.SendFailResponse(w,"failed to get looks")
		} else {
			utils.SendDataResponse(w, looks)
		}
	}
}

func LikeActionHandler(w http.ResponseWriter, r *http.Request) {
	lookedUserId, ok := utils.UnmarshalHttpBodyToLoginData(w, r)
	if !ok {
		return
	}

	loginData, err := structuredDataStorage.Manager.GetUserLoginDataBySession(utils.GetCookieValue(r, "session_id"))
	if err != nil {
		log.Error("Session Id is invalid: ", err)
		utils.SendFailResponse(w, "Your session ID is invalid, try to refresh the page.")
		return
	}

	if r.Method == http.MethodPost {
		if userDataStorage.Manager.SaveLiked(lookedUserId.Id, loginData.Id) {
			utils.SendSuccessResponse(w)
		} else {
			utils.SendFailResponse(w,"failed to save looked to db.")
		}
	} else if r.Method == http.MethodDelete {
		if userDataStorage.Manager.DeleteInteraction(loginData, lookedUserId.Id) {
			utils.SendSuccessResponse(w)
		} else {
			utils.SendFailResponse(w,"failed to delete interactions")
		}
	} else if r.Method == http.MethodGet {
		likes, err := userDataStorage.Manager.GetPreviousInteractions(loginData, "like")
		if err != nil {
			utils.SendFailResponse(w,"failed to get likes")
		} else {
			utils.SendDataResponse(w, likes)
		}
	}
}

func MatchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		requestData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error("Can't read request body for login: ", err)
			utils.SendFailResponse(w, "can't read request body")
			return
		}

		user1Data, err := structuredDataStorage.Manager.GetUserLoginDataBySession(utils.GetCookieValue(r, "session_id"))
		if err != nil {
			log.Error("Session Id is invalid: ", err)
			utils.SendFailResponse(w, "Your session ID is invalid, try to refresh the page.")
			return
		}

		user2Data := &model.LoginData{}
		err = json.Unmarshal(requestData, &user2Data)
		if err != nil {
			log.Error("Can't parse request body for login: ", err)
			utils.SendFailResponse(w, "can't read request body")
			return
		}
		if userDataStorage.Manager.SaveMatch(user1Data.Id, user2Data.Id) {
			utils.SendSuccessResponse(w)
		} else {
			utils.SendFailResponse(w,"failed to save matched to db.")
		}
	}
}

func ManageBannedUsersHandler(w http.ResponseWriter, r *http.Request) {
	session := utils.GetCookieValue(r, "session_id")
	data, err := structuredDataStorage.Manager.GetUserLoginDataBySession(session)
	if err != nil {
		utils.SendFailResponse(w, "incorrect session id")
		return
	}

	if r.Method == http.MethodGet {
		bans, err := userDataStorage.Manager.GetUserBannedList(data)
		if err != nil {
			log.Errorf("Error getting banned users: %v", err)
			utils.SendFailResponse(w, "Error getting banned users")
			return
		}
		utils.SendDataResponse(w, bans)
	} else if r.Method == http.MethodPost {
		bannedLoginData, ok := utils.UnmarshalHttpBodyToLoginData(w, r)
		if !ok {
			return
		}
		if ok := userDataStorage.Manager.AddUserIdToBanned(data, bannedLoginData.Id); !ok {
			utils.SendFailResponse(w, "Failed to ban user")
		} else {
			utils.SendSuccessResponse(w)
		}
	} else if r.Method == http.MethodDelete {
		bannedLoginData, ok := utils.UnmarshalHttpBodyToLoginData(w, r)
		if !ok {
			return
		}
		if ok := userDataStorage.Manager.RemoveUserIdFromBanned(data, bannedLoginData.Id); !ok {
			utils.SendFailResponse(w, "Failed to remove user form banned")
		} else {
			utils.SendSuccessResponse(w)
		}
	}
}
