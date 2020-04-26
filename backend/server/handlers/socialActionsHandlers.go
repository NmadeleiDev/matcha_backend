package handlers

import (
	"backend/db/mongodb"
	"backend/db/postgres"
	"backend/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func GetStrangersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		session := utils.GetCookieValue(r, "session_id")
		user, err := postgres.GetUserEmailBySession(session)
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
