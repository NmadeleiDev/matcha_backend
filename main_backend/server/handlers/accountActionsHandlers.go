package handlers

import (
	"backend/db/userMetaDataStorage"
	"backend/db/userFullDataStorage"
	"backend/dto"
	"backend/emails"
	"backend/hashing"
	"backend/model"
	"backend/utils"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
	"time"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		userData, ok := utils.UnmarshalHttpBodyToUserData(w, r)
		if !ok {
			return
		}

		authKey, ok := userMetaDataStorage.Manager.CreateUser(userData)
		if !ok {
			utils.SendFailResponse(w, "User already exists")
			return
		}

		if !userFullDataStorage.Manager.CreateUser(*userData) {
			utils.SendFailResponse(w, "failed to create user")
			return
		}

		emails.Manager.SendAccountVerificationKey(userData.Email, authKey)
		utils.SendSuccessResponse(w)
	} else {
		utils.SendFailResponseWithCode(w, "Not allowed", http.StatusMethodNotAllowed)
	}
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		loginData, ok := utils.UnmarshalHttpBodyToLoginData(w, r)
		if !ok {
			return
		}
		if err := userMetaDataStorage.Manager.LoginUser(loginData); err == nil {
			var err error
			var userData model.FullUserData

			if len(r.Header.Get("latitude")) == 0 || len(r.Header.Get("longitude")) == 0 {
				userData, err = userFullDataStorage.Manager.GetUserData(*loginData, false)
			} else {
				position := model.Coordinates{
					Lat: utils.UnsafeAtof(r.Header.Get("latitude"), 0),
					Lon: utils.UnsafeAtof(r.Header.Get("longitude"), 0),
				}
				userData, err = userFullDataStorage.Manager.FindUserAndUpdateGeo(*loginData, position)
			}
			if err != nil {
				log.Error("Failed to get user data")
				utils.SendFailResponse(w,"Failed to get user data")
			} else {
				userData.Id = loginData.Id
				if userMetaDataStorage.Manager.RefreshRequestSessionKeyCookie(w, *loginData) {
					utils.SendDataResponse(w, userData)
				}
			}
		} else {
			if err.Error() == "STATE" {
				utils.SendFailResponseWithCode(w,"account not verified", http.StatusUnauthorized)
			} else {
				utils.SendFailResponseWithCode(w,"incorrect user data", http.StatusUnauthorized)
			}
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

	newSessionKey, ok := userMetaDataStorage.Manager.VerifyUserAccountState(key)
	if ok {
		utils.SetCookieForDay(w, "session_id", newSessionKey)
		_, err := userMetaDataStorage.Manager.GetUserLoginDataBySession(newSessionKey)
		if err != nil {
			utils.SendFailResponse(w,"Failed to get user data")
		} else {
			var url string
			host := os.Getenv("PROJECT_HOST")
			port := os.Getenv("PROJECT_PORT")

			if len(host) == 0 {
				url = fmt.Sprintf("http://localhost:%v", port)
			} else {
				url = fmt.Sprintf("https://%v", host)
			}
			http.Redirect(w, r, url, http.StatusPermanentRedirect)
			//data, err := userFullDataStorage.Manager.GetFullUserData(login, true)
			//if err != nil {
			//	utils.SendFailResponse(w, "Failed to get user data")
			//} else {
			//	utils.SendDataResponse(w, data)
			//}
		}
	} else {
		utils.SendFailResponse(w, "failed to verify user")
	}
}

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var key string
		userEmail, ok := utils.UnmarshalHttpBodyToLoginData(w, r)
		if !ok {
			return
		}

		if randStr, err := utils.GenerateRandomString(10); err != nil {
			log.Errorf("error generating rand str: %v", err)
			utils.SendFailResponse(w, "Internal error")
			return
		} else {
			key = hashing.CalculateSha256(randStr + time.Now().String())
		}

		userId, err := userMetaDataStorage.Manager.GetUserIdByEmail(userEmail.Email)
		if err != nil {
			log.Errorf("Error getting user id by email: %v", err)
			utils.SendSuccessResponse(w) // тут специально отправляется успешный ответ, так как запросившему смену пароля не нужно знать, что этого меила не существует (защита от user enumeration)
			return
		}

		if err := userMetaDataStorage.Manager.CreateResetPasswordRecord(userId, key); err != nil {
			log.Errorf("Error creating reset password lot: %v", err)
			utils.SendFailResponse(w, "internal error")
			return
		}
		emails.Manager.SendPasswordResetEmail(userEmail.Email, key)
		utils.SendSuccessResponse(w)

	} else if r.Method == http.MethodGet {
		key := r.URL.Query().Get("k")
		var newKey string
		if len(key) != 64 {
			utils.SendFailResponse(w, "Invalid key")
			log.Infof("len = %v", len(key))
			return
		}
		if randStr, err := utils.GenerateRandomString(5); err != nil {
			log.Errorf("error generating rand str: %v", err)
			utils.SendFailResponse(w, "Internal error")
			return
		} else {
			newKey = hashing.CalculateSha256(randStr + key + time.Now().String())
		}

		if err := userMetaDataStorage.Manager.SetNextStepResetKey(key, newKey); err != nil {
			utils.SendFailResponse(w, "Key is invalid")
		} else {
			utils.SetHttpOnlyCookie(w, "reset_key", newKey)
			utils.SendSuccessResponse(w)
		}
	} else if r.Method == http.MethodPut {
		newPassword, ok := utils.UnmarshalHttpBodyToLoginData(w, r)
		if !ok {
			return
		}

		resetKey := utils.GetCookieValue(r, "reset_key")

		if id, err := userMetaDataStorage.Manager.GetAccountIdByResetKey(resetKey); err != nil {
			utils.SendFailResponse(w, "Key is invalid")
		} else {
			if err1 := userMetaDataStorage.Manager.SetNewPasswordForAccount(id, newPassword.Password); err1 != nil {
				utils.SendFailResponse(w, "Password update error")
			} else {
				utils.SendSuccessResponse(w)
			}
		}
	}
}

func SignOutHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodDelete {
		session := utils.GetCookieValue(r,"session_id")
		ok := userMetaDataStorage.Manager.UpdateSessionKey(session, "")
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
		utils.SendFailResponseWithCode(w, "incorrect session id", http.StatusUnauthorized)
		return
	}

	loginData, err := userMetaDataStorage.Manager.GetUserLoginDataBySession(session)
	if err != nil {
		utils.SendFailResponse(w, "incorrect session id")
		return
	}
	if r.Method == http.MethodGet {
		userDto := dto.GetUserDTO(&model.FullUserData{Id: loginData.Id}).LoadUserData(false).PrepareUserDataForClient()
		userData := userDto.GetUser()
		if userDto.GetError() != nil {
			utils.SendFailResponseWithCode(w,"Failed to get user data", http.StatusInternalServerError)
		} else {
			filteredLiked := make([]string, 0, len(userData.LikedBy))
			for _, id := range userData.LikedBy {
				if !utils.DoesArrayContain(userData.Matches, id) {
					filteredLiked = append(filteredLiked, id)
				}
			}
			userData.LikedBy = filteredLiked
			utils.SendDataResponse(w, userData)
			return
		}
	} else if r.Method == http.MethodPut {
		userData, ok := utils.UnmarshalHttpBodyToUserData(w, r)
		if !ok {
			return
		}

		loginData, err := userMetaDataStorage.Manager.GetUserLoginDataBySession(utils.GetCookieValue(r, "session_id"))
		if err != nil {
			log.Error("Can't get user is: ", err)
			utils.SendFailResponse(w, "Session cookie not present")
			return
		}

		userData.Id = loginData.Id
		if !userFullDataStorage.Manager.UpdateUser(*userData) {
			utils.SendFailResponse(w, "failed to update user")
			return
		}
		utils.SendSuccessResponse(w)
	} else if r.Method == http.MethodDelete {
		userData, err := userFullDataStorage.Manager.GetFullUserData(loginData, true)
		if err != nil {
			utils.SendFailResponse(w, fmt.Sprintf("Failed to get user data: %v", err))
			return
		}
		for _, tagId := range userData.TagIds {
			_ = userMetaDataStorage.Manager.DecrTagById(tagId)
		}
		if err := userFullDataStorage.Manager.DeleteAccount(loginData); err != nil {
			utils.SendFailResponse(w, fmt.Sprintf("Failed to delete user data: %v", err))
			return
		}
		if err := userMetaDataStorage.Manager.DeleteAccount(loginData); err != nil {
			utils.SendFailResponse(w, fmt.Sprintf("Failed to delete user account metadata: %v", err))
			return
		}
		if err := userFullDataStorage.Manager.DeleteAccountRecordsFromOtherUsers(loginData); err != nil {
			utils.SendFailResponse(w, fmt.Sprintf("Failed to delete records from other users accounts: %v", err))
			return
		}
		emails.Manager.SendGoodbyeMessage(userData.Email)
		log.Infof("Deleted %v", userData.Email)

		utils.SendSuccessResponse(w)
	}
}

func EmailActionsHandler(w http.ResponseWriter, r *http.Request) {
	const (
		changeEmailAction = "change"
		verifyEmailAction = "verify"
		keyLen = 32
	)

	action := mux.Vars(r)["action"]
	if r.Method == http.MethodPut && action == changeEmailAction {
		loginData := userMetaDataStorage.Manager.AuthUserBySessionId(w, r)
		emailData, ok := utils.UnmarshalHttpBodyToLoginData(w, r)
		if loginData == nil || !ok {
			return
		}
		key, err := utils.GenerateRandomString(keyLen)
		if err != nil {
			log.Errorf("Failed to generate key to change email: %v", err)
			utils.SendFailResponseWithCode(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
			return
		}
		if userMetaDataStorage.Manager.CreateResetEmailRecord(loginData.Id, emailData.Email, key) {
			emails.Manager.SendEmailVerificationKey(emailData.Email, key)
			utils.SendSuccessResponse(w)
		} else {
			utils.SendFailResponseWithCode(w, "Failed to init procedure", http.StatusInternalServerError)
		}
	} else if r.Method == http.MethodGet && action == verifyEmailAction {
		key := r.URL.Query().Get("key")
		if len(key) != keyLen {
			utils.SendFailResponseWithCode(w, fmt.Sprintf("Key len expected %v, got %v", keyLen, len(key)), http.StatusBadRequest)
			return
		}
		userId, email, err := userMetaDataStorage.Manager.GetResetEmailRecord(key)
		if err != nil {
			log.Errorf("Error getting email: %v", err)
			utils.SendFailResponseWithCode(w, "Failed to get email!", http.StatusInternalServerError)
			return
		}
		if err := userMetaDataStorage.Manager.SetNewEmail(userId, email); err != nil {
			utils.SendFailResponseWithCode(w, "Failed to set email! " + err.Error(), http.StatusInternalServerError)
			return
		}
		if err := userFullDataStorage.Manager.SetNewEmail(userId, email); err != nil {
			utils.SendFailResponseWithCode(w, "Failed to set email! " + err.Error(), http.StatusInternalServerError)
			return
		}
		utils.SendDataResponse(w, model.LoginData{Id: userId, Email: email})
	}
}

func UserTagsHandler(w http.ResponseWriter, r *http.Request) {
	loginData := userMetaDataStorage.Manager.AuthUserBySessionId(w, r)
	if loginData == nil {
		return
	}
	if r.Method == http.MethodPost {
		tags, ok := utils.UnmarshalHttpBodyToTags(w, r)
		if !ok {
			return
		}
		failedTags := make([]string, 0, len(tags.Tags))
		for _, tag := range tags.Tags {
			id, err := userMetaDataStorage.Manager.IncOrInsertTag(tag)
			if err != nil {
				failedTags = append(failedTags, tag)
			} else {
				userFullDataStorage.Manager.AddTagToUserTags(*loginData, id)
			}
		}

		if len(failedTags) == 0 {
			utils.SendSuccessResponse(w)
		} else {
			utils.SendFailResponse(w, fmt.Sprintf("Failed to save tags: %v", failedTags))
		}
	} else if r.Method == http.MethodPut {
		userTagIds := userFullDataStorage.Manager.GetUserDataWithCustomProjection(*loginData, []string{"tag_ids"}, true).TagIds
		userTags := userMetaDataStorage.Manager.GetTagsById(userTagIds)
		newTags, ok := utils.UnmarshalHttpBodyToTags(w, r)
		if !ok {
			return
		}
		failedTags := make([]string, 0, len(userTagIds) + len(newTags.Tags))

		for _, tag := range append(userTags, newTags.Tags...) {
			exists := utils.DoesArrayContain(userTags, tag)
			posted := utils.DoesArrayContain(newTags.Tags, tag)
			if exists && !posted {
				id, err := userMetaDataStorage.Manager.DecrTagByValue(tag)
				ok := userFullDataStorage.Manager.DeleteTagFromUserTags(*loginData, id)
				if err != nil || !ok {
					failedTags = append(failedTags, tag)
				}
			} else if posted && !exists {
				id, err := userMetaDataStorage.Manager.IncOrInsertTag(tag)
				if err != nil {
					failedTags = append(failedTags, tag)
				} else {
					userFullDataStorage.Manager.AddTagToUserTags(*loginData, id)
				}
			}
		}
		if len(failedTags) == 0 {
			utils.SendSuccessResponse(w)
		} else {
			utils.SendFailResponse(w, fmt.Sprintf("Failed to manage tags: %v", failedTags))
		}
	} else if r.Method == http.MethodDelete {
		tags, ok := utils.UnmarshalHttpBodyToTags(w, r)
		if !ok {
			return
		}
		failedTags := make([]string, 0, len(tags.Tags))
		for _, tag := range tags.Tags {
			id, err := userMetaDataStorage.Manager.DecrTagByValue(tag)
			ok := userFullDataStorage.Manager.DeleteTagFromUserTags(*loginData, id)
			if err != nil || !ok {
				failedTags = append(failedTags, tag)
			}
		}
		go userMetaDataStorage.Manager.ClearUnmentionedTags()

		if len(failedTags) == 0 {
			utils.SendSuccessResponse(w)
		} else {
			utils.SendFailResponse(w, fmt.Sprintf("Failed to delete tags: %v", failedTags))
		}
	} else if r.Method == http.MethodGet {
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

		utils.SendDataResponse(w, userMetaDataStorage.Manager.GetAllTags(limit, offset))
	}
}

func GetUserDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		var userData interface{}
		var err error

		id := mux.Vars(r)["id"]
		isShortData := r.URL.Query().Get("full") == "false"

		if userMetaDataStorage.Manager.AuthUserBySessionId(w, r) == nil {
			return
		}

		if isShortData {
			userData, err = userFullDataStorage.Manager.GetShortUserData(model.LoginData{Id: id})
		} else {
			userDto := dto.GetUserDTO(&model.FullUserData{Id: id}).LoadUserData(true).PrepareUserDataForClient()
			userData = userDto.GetUser()
			err = userDto.GetError()
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

		data := userMetaDataStorage.Manager.AuthUserBySessionId(w, r)
		if data == nil {
			return
		}

		actions, err := userFullDataStorage.Manager.GetPreviousInteractions(*data, action)
		if err != nil {
			utils.SendFailResponse(w, err.Error())
			return
		}
		utils.SendDataResponse(w, actions)
	}
}

func GetUserOwnImagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data := userMetaDataStorage.Manager.AuthUserBySessionId(w, r)
		if data == nil {
			return
		}
		utils.SendDataResponse(w, userFullDataStorage.Manager.GetUserImages(data.Id))
	}
}
