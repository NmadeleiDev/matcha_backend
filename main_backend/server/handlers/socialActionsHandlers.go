package handlers

import (
	"backend/db/mongodb"
	"backend/db/postgres"
	"backend/structs"
	"backend/utils"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func GetStrangersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		session := utils.GetCookieValue(r, "session_id")
		user, err := postgres.GetUserIdBySession(session)
		if err != nil {
			log.Error("Failed to get user data by session")
			utils.SendFailResponse(w, "incorrect user data")
			return
		}
		userData, err := mongodb.GetUserData(user)
		if err != nil {
			log.Error("Failed to get user data from mongo")
			utils.SendFailResponse(w, "sorry!")
			return
		}
		strangers, ok := mongodb.GetFittingUsers(userData)
		if ok {
			utils.SendDataResponse(w, strangers)
		} else {
			utils.SendFailResponse(w, "Failed to load users")
		}
	}
}

func SaveAccountLookUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		requestData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error("Can't read request body for login: ", err)
			utils.SendFailResponse(w, "can't read request body")
			return
		}

		loginData, err := postgres.GetUserIdBySession(utils.GetCookieValue(r, "session_id"))
		if err != nil {
			log.Error("Session Id is invalid: ", err)
			utils.SendFailResponse(w, "Your session ID is invalid, try to refresh the page.")
			return
		}

		lookedUserId := &structs.LoginData{}
		err = json.Unmarshal(requestData, &lookedUserId)
		if err != nil {
			log.Error("Can't parse request body for login: ", err)
			utils.SendFailResponse(w, "can't read request body")
			return
		}
		if mongodb.SaveLooked(lookedUserId.Id, loginData.Id) {
			utils.SendSuccessResponse(w)
		} else {
			utils.SendFailResponse(w,"failed to save looked to db.")
		}
	}
}

func SaveLikeActionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		requestData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error("Can't read request body for login: ", err)
			utils.SendFailResponse(w, "can't read request body")
			return
		}

		loginData, err := postgres.GetUserIdBySession(utils.GetCookieValue(r, "session_id"))
		if err != nil {
			log.Error("Session Id is invalid: ", err)
			utils.SendFailResponse(w, "Your session ID is invalid, try to refresh the page.")
			return
		}

		lookedUserId := &structs.LoginData{}
		err = json.Unmarshal(requestData, &lookedUserId)
		if err != nil {
			log.Error("Can't parse request body for login: ", err)
			utils.SendFailResponse(w, "can't read request body")
			return
		}
		if mongodb.SaveLiked(lookedUserId.Id, loginData.Id) {
			utils.SendSuccessResponse(w)
		} else {
			utils.SendFailResponse(w,"failed to save looked to db.")
		}
	}
}

func SaveMatchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		requestData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error("Can't read request body for login: ", err)
			utils.SendFailResponse(w, "can't read request body")
			return
		}

		user1Data, err := postgres.GetUserIdBySession(utils.GetCookieValue(r, "session_id"))
		if err != nil {
			log.Error("Session Id is invalid: ", err)
			utils.SendFailResponse(w, "Your session ID is invalid, try to refresh the page.")
			return
		}

		user2Data := &structs.LoginData{}
		err = json.Unmarshal(requestData, &user2Data)
		if err != nil {
			log.Error("Can't parse request body for login: ", err)
			utils.SendFailResponse(w, "can't read request body")
			return
		}
		if mongodb.SaveMatch(user1Data.Id, user2Data.Id) {
			utils.SendSuccessResponse(w)
		} else {
			utils.SendFailResponse(w,"failed to save matched to db.")
		}
	}
}