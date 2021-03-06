package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"backend/db/realtimeDataDb"
	"backend/db/userFullDataStorage"
	"backend/db/userMetaDataStorage"
	"backend/dto"
	"backend/utils"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func GetStrangersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		user := userMetaDataStorage.Manager.AuthUserBySessionId(w, r)
		if user == nil {
			return
		}

		userData, err := userFullDataStorage.Manager.GetFullUserData(*user, false)
		if err != nil {
			log.Error("Failed to get user data from mongo")
			utils.SendFailResponse(w, "sorry!")
			return
		}
		mainUserDto := dto.GetUserDTO(&userData).PrepareUserDataForClient()
		userData = *mainUserDto.GetUser()

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
			userData.LookFor = r.URL.Query().Get("gender")
		}
		if r.URL.Query().Get("tags") != "" {
			userData.Tags = strings.Split(
				strings.ReplaceAll(r.URL.Query().Get("tags"), " ", ""), ",")
		}
		//log.Infof("UserData for strangers: %v", userData)

		strangers, ok := userFullDataStorage.Manager.GetFittingUsers(userData)
		for i, acc := range strangers {
			userDto := dto.GetUserDTO(&acc).PrepareUserDataForClient()
			userData := userDto.GetUser()
			if userDto.GetError() != nil {
				log.Errorf("Error preparing user data for strangers: %v", err)
				continue
			}
			strangers[i] = *userData
		}
		if ok {
			utils.SendDataResponse(w, strangers)
		} else {
			utils.SendFailResponse(w, "Failed to load users")
		}
	}
}

func LookActionHandler(w http.ResponseWriter, r *http.Request) {

	loginData := userMetaDataStorage.Manager.AuthUserBySessionId(w, r)
	if loginData == nil {
		return
	}

	if r.Method == http.MethodPost {
		lookedUserId, ok := utils.UnmarshalHttpBodyToLoginData(w, r)
		if !ok {
			return
		}
		if userFullDataStorage.Manager.SaveLooked(lookedUserId.Id, loginData.Id) {
			realtimeDataDb.GetManager().PublishMessage(lookedUserId.Id, realtimeDataDb.LookType, loginData.Id)
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

	loginData := userMetaDataStorage.Manager.AuthUserBySessionId(w, r)
	if loginData == nil {
		return
	}

	if r.Method == http.MethodPost {
		likedUserId, ok := utils.UnmarshalHttpBodyToLoginData(w, r)
		if len(likedUserId.Id) == 0 {
			utils.SendFailResponseWithCode(w, "User is not set!", http.StatusBadRequest)
			return
		}
		if !ok {
			return
		}
		userData := userFullDataStorage.Manager.GetUserDataWithCustomProjection(*loginData, []string{"liked_by"}, true)
		if utils.DoesArrayContain(userData.LikedBy, likedUserId.Id) {
			if userFullDataStorage.Manager.SaveMatch(loginData.Id, likedUserId.Id) &&
				userFullDataStorage.Manager.SaveLiked(likedUserId.Id, loginData.Id) {
				realtimeDataDb.GetManager().PublishMessage(loginData.Id, realtimeDataDb.CreatedMatchType, likedUserId.Id)
				realtimeDataDb.GetManager().PublishMessage(likedUserId.Id, realtimeDataDb.CreatedMatchType, loginData.Id)
				utils.SendSuccessResponse(w)
			} else {
				utils.SendFailResponse(w,"failed to save matched to db.")
			}
		} else {
			if userFullDataStorage.Manager.SaveLiked(likedUserId.Id, loginData.Id) {
				realtimeDataDb.GetManager().PublishMessage(likedUserId.Id, realtimeDataDb.CreatedLikeType, loginData.Id)
				utils.SendSuccessResponse(w)
			} else {
				utils.SendFailResponse(w,"failed to save liked to db.")
			}
		}
	} else if r.Method == http.MethodDelete {
		likedId := mux.Vars(r)["user_id"]
		if len(likedId) != 64 {
			utils.SendFailResponse(w, fmt.Sprintf("User id len is invalid. Must be 64, got %v", len(likedId)))
			return
		}
		if isMatchDelete, ok := userFullDataStorage.Manager.DeleteLikeOrMatch(*loginData, likedId); ok {
			if isMatchDelete {
				realtimeDataDb.GetManager().PublishMessage(likedId, realtimeDataDb.DeletedMatchType, loginData.Id)
				realtimeDataDb.GetManager().PublishMessage(loginData.Id, realtimeDataDb.DeletedMatchType, likedId)
			} else {
				realtimeDataDb.GetManager().PublishMessage(likedId, realtimeDataDb.DeletedLikeType, loginData.Id)
			}
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

func ManageBannedUsersHandler(w http.ResponseWriter, r *http.Request) {
	data := userMetaDataStorage.Manager.AuthUserBySessionId(w, r)
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



func ReportUserHandler(w http.ResponseWriter, r *http.Request) {
	data := userMetaDataStorage.Manager.AuthUserBySessionId(w, r)
	if data == nil {
		return
	}

	if r.Method == http.MethodPost {
		report, ok := utils.UnmarshalHttpBodyToReport(w, r)
		if !ok {
			return
		}
		if userFullDataStorage.Manager.SaveReport(*report) {
			utils.SendSuccessResponse(w)
		} else {
			utils.SendFailResponseWithCode(w, "Failed to save report", http.StatusInternalServerError)
		}
	} else {
		utils.SendFailResponseWithCode(w, "Not allowed", http.StatusMethodNotAllowed)
	}
}
