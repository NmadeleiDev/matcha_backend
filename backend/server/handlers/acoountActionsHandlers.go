package handlers

import (
	"backend/db/postgres"
	"backend/server/utils"
	"backend/structs"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func	SignUpHandler(w http.ResponseWriter, r *http.Request) {

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
			utils.SendFailResponse(w)
			return
		}

		//utils.RefreshRequestSessionKeyCookie(w, *userData)
		utils.SendSuccessResponse(w)
	}
}

func	SignInHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		requestData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error("Can't read request body for login: ", err)
			return
		}

		loginData := &structs.LoginData{}
		err = json.Unmarshal(requestData, loginData)
		if err != nil {
			log.Error("Can't parse request body for login: ", err)
			return
		}
		if postgres.LoginUser(*loginData) {
			userData, err := postgres.GetFullUserData(*loginData)
			if err != nil {
				log.Error("Failed to get user data")
				utils.SendFailResponse(w)
			} else {
				utils.SendDataResponse(w, userData)
			}
		} else {
			utils.SendFailResponse(w)
		}

		//utils.RefreshRequestSessionKeyCookie(w, *userData)
	}
}
