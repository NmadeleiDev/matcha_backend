package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"backend/db/notificationsBroker"
	"backend/db/userFullDataStorage"
	"backend/model"
	"backend/utils"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func GetStrangersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		user := utils.AuthUserBySessionId(w, r)
		if user == nil {
			return
		}

		userData, err := userFullDataStorage.Manager.GetFullUserData(*user, "full")
		if err != nil {
			log.Error("Failed to get user data from mongo")
			utils.SendFailResponse(w, "sorry!")
			return
		}

		userData.MinAge = utils.UnsafeAtoi(r.URL.Query().Get("minAge"), userData.MinAge)
		userData.MaxAge = utils.UnsafeAtoi(r.URL.Query().Get("maxAge"), userData.MaxAge)
		userData.MaxDist = utils.UnsafeAtoi(r.URL.Query().Get("maxDist"), userData.MaxDist)

		if r.URL.Query().Get("city") != "" {
			userData.City = r.URL.Query().Get("city")
		}
		if r.URL.Query().Get("country") != "" {
			userData.Country = r.URL.Query().Get("country")
		}
		if r.URL.Query().Get("gender") != "" {
			userData.Gender = r.URL.Query().Get("gender")
		}
		if r.URL.Query().Get("tags") != "" {
			userData.Tags = strings.Split(
				strings.ReplaceAll(r.URL.Query().Get("tags"), " ", ""), ",")
		}

		strangers, ok := userFullDataStorage.Manager.GetFittingUsers(userData)
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

	loginData := utils.AuthUserBySessionId(w, r)
	if loginData == nil {
		return
	}

	if r.Method == http.MethodPost {
		if userFullDataStorage.Manager.SaveLooked(lookedUserId.Id, loginData.Id) {
			notificationsBroker.GetManager().PublishMessage(lookedUserId.Id, notificationsBroker.LookType, loginData.Id)
			utils.SendSuccessResponse(w)
		} else {
			utils.SendFailResponse(w,"failed to save looked to db.")
		}
	} else if r.Method == http.MethodGet {
		looks, err := userFullDataStorage.Manager.GetPreviousInteractions(*loginData, "look")
		if err != nil {
			utils.SendFailResponse(w,"failed to get looks")
		} else {
			utils.SendDataResponse(w, looks)
		}
	}
}

func LikeActionHandler(w http.ResponseWriter, r *http.Request) {
	likedUserId, ok := utils.UnmarshalHttpBodyToLoginData(w, r)
	if !ok {
		return
	}

	loginData := utils.AuthUserBySessionId(w, r)
	if loginData == nil {
		return
	}

	if r.Method == http.MethodPost {
		userData := userFullDataStorage.Manager.GetUserDataWithCustomProjection(*loginData, []string{"liked_by"}, true)
		if utils.DoesArrayContain(userData.LikedBy, likedUserId.Id) {
			if userFullDataStorage.Manager.SaveMatch(loginData.Id, likedUserId.Id) {
				notificationsBroker.GetManager().PublishMessage(loginData.Id, notificationsBroker.MatchType, likedUserId.Id)
				notificationsBroker.GetManager().PublishMessage(likedUserId.Id, notificationsBroker.MatchType, loginData.Id)
				utils.SendSuccessResponse(w)
			} else {
				utils.SendFailResponse(w,"failed to save matched to db.")
			}
		} else {
			if userFullDataStorage.Manager.SaveLiked(likedUserId.Id, loginData.Id) {
				notificationsBroker.GetManager().PublishMessage(likedUserId.Id, notificationsBroker.LikeType, loginData.Id)
				utils.SendSuccessResponse(w)
			} else {
				utils.SendFailResponse(w,"failed to save looked to db.")
			}
		}
	} else if r.Method == http.MethodDelete {
		likedId := r.URL.Query().Get("id")
		if len(likedId) == 0 {
			utils.SendFailResponse(w,"Id to delete is invalid: " + likedId)
			return
		}
		if userFullDataStorage.Manager.DeleteInteraction(*loginData, likedId) {
			utils.SendSuccessResponse(w)
		} else {
			utils.SendFailResponse(w,"failed to delete interactions")
		}
	} else if r.Method == http.MethodGet {
		likes, err := userFullDataStorage.Manager.GetPreviousInteractions(*loginData, "like")
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

		user1Data := utils.AuthUserBySessionId(w, r)
		if user1Data == nil {
			return
		}

		user2Data := &model.LoginData{}
		err = json.Unmarshal(requestData, &user2Data)
		if err != nil {
			log.Error("Can't parse request body for login: ", err)
			utils.SendFailResponse(w, "can't read request body")
			return
		}
		if userFullDataStorage.Manager.SaveMatch(user1Data.Id, user2Data.Id) {
			notificationsBroker.GetManager().PublishMessage(user1Data.Id, notificationsBroker.MatchType, user2Data.Id)
			notificationsBroker.GetManager().PublishMessage(user2Data.Id, notificationsBroker.MatchType, user1Data.Id)
			utils.SendSuccessResponse(w)
		} else {
			utils.SendFailResponse(w,"failed to save matched to db.")
		}
	}
}

func ManageBannedUsersHandler(w http.ResponseWriter, r *http.Request) {
	data := utils.AuthUserBySessionId(w, r)
	if data == nil {
		return
	}

	if r.Method == http.MethodGet {
		bans, err := userFullDataStorage.Manager.GetUserBannedList(*data)
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
		if ok := userFullDataStorage.Manager.AddUserIdToBanned(*data, bannedLoginData.Id); !ok {
			utils.SendFailResponse(w, "Failed to ban user")
		} else {
			utils.SendSuccessResponse(w)
		}
	} else if r.Method == http.MethodDelete {
		idToUnban := mux.Vars(r)["id"]
		if len(idToUnban) < 5 {
			utils.SendFailResponse(w, "id is incorrect")
			return
		}
		if ok := userFullDataStorage.Manager.RemoveUserIdFromBanned(*data, idToUnban); !ok {
			utils.SendFailResponse(w, "Failed to remove user form banned")
		} else {
			utils.SendSuccessResponse(w)
		}
	}
}
