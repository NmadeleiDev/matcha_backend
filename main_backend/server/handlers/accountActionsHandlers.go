package handlers

import (
	"backend/db/mongodb"
	"backend/db/postgres"
	"backend/structs"
	"backend/utils"
	"encoding/json"
	"github.com/gorilla/mux"
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

		if !postgres.CreateUser(userData) {
			utils.SendFailResponse(w, "failed to create user")
			return
		}

		if !mongodb.CreateUser(*userData) {
			utils.SendFailResponse(w, "failed to create user")
			return
		}

		loginData := structs.LoginData{Email: userData.Email, Password: userData.Password, Id: userData.Id}
		if utils.RefreshRequestSessionKeyCookie(w, loginData) {
			userData.Password = ""
			utils.SendDataResponse(w, userData)
		}
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
		if postgres.LoginUser(loginData) {
			userData, err := mongodb.GetUserData(*loginData)
			if err != nil {
				userData.Id = loginData.Id
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

func UpdateAccountHandler(w http.ResponseWriter, r *http.Request) {
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

		loginData, err := postgres.GetUserIdBySession(utils.GetCookieValue(r, "session_id"))
		if err != nil {
			log.Error("Can't get user is: ", err)
			return
		}

		if userData.Id != loginData.Id {
			log.Warn("request body id does not match session_id id")
			utils.SendFailResponse(w, "id is incorrect")
			postgres.SetSessionKeyById("", loginData.Id)
			return
		}

		if !mongodb.UpdateUser(*userData) {
			utils.SendFailResponse(w, "failed to update user")
			return
		}
		utils.SendSuccessResponse(w)
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

func GetUserDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		id := mux.Vars(r)["id"]
		session := utils.GetCookieValue(r,"session_id")
		_, err := postgres.GetUserIdBySession(session)
		if err != nil {
			utils.SendFailResponse(w, "incorrect session id")
			return
		}
		userData, err := mongodb.GetUserData(structs.LoginData{Id: id})
		userData.LikedBy = []string{}
		userData.LookedBy = []string{}
		userData.Matches = []string{}
		if err != nil {
			userData.Id = id
			log.Error("Failed to get user data")
			utils.SendFailResponse(w,"Failed to get user data")
		} else {
			utils.SendDataResponse(w, userData)
			return
		}
	}
}