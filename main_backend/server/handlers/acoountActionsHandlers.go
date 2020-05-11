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

func SignUpHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		requestData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error("Can't read request body for login: ", err)
			return
		}

		userData := &structs.UserData{}
		err = json.Unmarshal(requestData, userData)
		if err != nil {
			log.Error("Can't parse request body for login: ", err)
			return
		}

		if !postgres.CreateUser(*userData) {
			utils.SendFailResponse(w, "failed to create user")
			return
		}

		if !mongodb.CreateUser(*userData) {
			utils.SendFailResponse(w, "failed to create user")
			return
		}

		loginData := structs.LoginData{Email: userData.Email, Password: userData.Password}
		utils.RefreshRequestSessionKeyCookie(w, loginData)
		utils.SendSuccessResponse(w)
	}
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		requestData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error("Can't read request body for login: ", err)
			utils.SendFailResponse(w, "can't read request body")
			return
		}

		loginData := &structs.LoginData{}
		err = json.Unmarshal(requestData, loginData)
		if err != nil {
			log.Error("Can't parse request body for login: ", err)
			utils.SendFailResponse(w, "can't read request body")
			return
		}
		if postgres.LoginUser(*loginData) {
			userData, err := mongodb.GetUserData(*loginData)
			if err != nil {
				log.Error("Failed to get user data")
				utils.SendFailResponse(w,"Failed to get user data")
			} else {
				utils.RefreshRequestSessionKeyCookie(w, *loginData)
				utils.SendDataResponse(w, userData)
				return
			}
		} else {
			utils.SendFailResponse(w,"incorrect user data")
		}
	}
}

func SignOutHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodDelete {
		session := utils.GetCookieValue(r,"session_id")
		ok := postgres.UpdateSessionKey(session, "")
		if ok {
			utils.SendSuccessResponse(w)
		} else {
			utils.SendFailResponse(w, "incorrect session")
		}
	}
}