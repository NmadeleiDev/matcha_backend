package handlers

import (
	"fmt"
	"net/http"

	"backend/db/structuredDataStorage"
	"backend/utils"
	"backend/wsClient"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func ManagerChatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		session := utils.GetCookieValue(r, "session_id")
		user, err := structuredDataStorage.Manager.GetUserLoginDataBySession(session)
		if err != nil {
			log.Error("Failed to get user data by session")
			utils.SendFailResponse(w, "incorrect user data")
			return
		}
		if chatId, exists := mux.Vars(r)["chat_id"]; exists {
			chat := wsClient.GetManager().FindChat(chatId)
			if chat != nil {
				utils.SendDataResponse(w, chat)
			} else {
				utils.SendFailResponse(w, fmt.Sprintf("chat not found for id = %v", chat))
			}
		} else {
			chats := wsClient.GetManager().GetUserChats(user.Id)
			utils.SendDataResponse(w, chats)
		}
	} else if r.Method == http.MethodPost {
		chat, ok := utils.UnmarshalHttpBodyToChat(w, r)
		if !ok {
			return
		}

		newChatId := wsClient.GetManager().CreateChat(*chat)

		if len(newChatId) > 0 {
			utils.SendDataResponse(w, newChatId)
		} else {
			utils.SendFailResponse(w, "Error creating chat")
		}
	}
}
