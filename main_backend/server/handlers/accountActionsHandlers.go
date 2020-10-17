package handlers

import (
	"backend/db/structuredDataStorage"
	"backend/db/userDataStorage"
	"backend/emails"
	"backend/types"
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

		userData := &types.UserData{}
		err = json.Unmarshal(requestData, userData)
		if err != nil {
			log.Error("Can't parse request body for login: ", err)
			return
		}

		authKey, ok := structuredDataStorage.Manager.CreateUser(userData)
		if !ok {
			utils.SendFailResponse(w, "User with this email already exists")
			return
		}

		if !userDataStorage.Manager.CreateUser(*userData) {
			utils.SendFailResponse(w, "failed to create user")
			return
		}

		emails.Send(userData.Email, authKey)
		utils.SendSuccessResponse(w)
		//loginData := types.LoginData{Email: userData.Email, Password: userData.Password, Id: userData.Id}
		//if utils.RefreshRequestSessionKeyCookie(w, loginData) {
		//	userData.Password = ""
		//	utils.SendDataResponse(w, userData)
		//}
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

		loginData := &types.LoginData{}
		err = json.Unmarshal(requestData, loginData)
		if err != nil {
			log.Error("Can't parse request body for login: ", err)
			utils.SendFailResponse(w, "can't read request body")
			return
		}
		if structuredDataStorage.Manager.LoginUser(loginData) {
			userData, err := userDataStorage.Manager.GetUserData(*loginData)
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

func UpdateAccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		requestData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error("Can't read request body for login: ", err)
			return
		}

		userData := &types.UserData{}
		err = json.Unmarshal(requestData, userData)
		if err != nil {
			log.Error("Can't parse request body: ", err)
			utils.SendFailResponse(w, "Can't parse request body")
			return
		}

		loginData, err := structuredDataStorage.Manager.GetUserLoginDataBySession(utils.GetCookieValue(r, "session_id"))
		if err != nil {
			log.Error("Can't get user is: ", err)
			utils.SendFailResponse(w, "Session cookie not present")
			return
		}

		if userData.Id != loginData.Id {
			log.Warn("request body id does not match session_id id")
			utils.SendFailResponse(w, "id is incorrect")
			structuredDataStorage.Manager.SetSessionKeyById("", loginData.Id)
			return
		}

		if !userDataStorage.Manager.UpdateUser(*userData) {
			utils.SendFailResponse(w, "failed to update user")
			return
		}
		utils.SendSuccessResponse(w)
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
			data, err := userDataStorage.Manager.GetUserData(login)
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

func GetUserDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		id := mux.Vars(r)["id"]
		session := utils.GetCookieValue(r,"session_id")
		_, err := structuredDataStorage.Manager.GetUserLoginDataBySession(session)
		if err != nil {
			utils.SendFailResponse(w, "incorrect session id")
			return
		}
		userData, err := userDataStorage.Manager.GetUserData(types.LoginData{Id: id})
		if err != nil {
			log.Error("Failed to get user data")
			utils.SendFailResponse(w,"Failed to get user data")
		} else {
			userData.LikedBy = []string{}
			userData.LookedBy = []string{}
			userData.Matches = []string{}
			utils.SendDataResponse(w, userData)
			return
		}
	}
}

func GetUserOwnImagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		session := utils.GetCookieValue(r, "session_id")
		data, err := structuredDataStorage.Manager.GetUserLoginDataBySession(session)
		if err != nil {
			utils.SendFailResponse(w, "incorrect session id")
			return
		}
		utils.SendDataResponse(w, userDataStorage.Manager.GetUserImages(data.Id))
	}
}