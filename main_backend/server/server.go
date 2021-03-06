package server

import (
	"net/http"

	"backend/server/handlers"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func StartServer(port string) {

	router := mux.NewRouter()

	router.HandleFunc("/signup", handlers.SignUpHandler)
	router.HandleFunc("/signin", handlers.SignInHandler)
	router.HandleFunc("/signout", handlers.SignOutHandler)
	router.HandleFunc("/verify/{key}", handlers.VerifyAccountHandler)
	router.HandleFunc("/reset", handlers.ResetPasswordHandler)

	router.HandleFunc("/account", handlers.ManageOwnAccountHandler)
	router.HandleFunc("/location", handlers.UpdateLocationHandler)
	router.HandleFunc("/email/{action}", handlers.EmailActionsHandler)
	router.HandleFunc("/actions/{action}", handlers.GetOwnActionsHandler)
	router.HandleFunc("/data/{id}", handlers.GetUserDataHandler) // получение данных любого юзера (только для залогиненных юзеров)
	router.HandleFunc("/ban", handlers.ManageBannedUsersHandler)
	router.HandleFunc("/ban/{id}", handlers.ManageBannedUsersHandler)
	router.HandleFunc("/tag", handlers.UserTagsHandler)
	router.HandleFunc("/media", handlers.GetUserOwnImagesHandler) // получить свои фотки

	router.HandleFunc("/strangers", handlers.GetStrangersHandler)
	router.HandleFunc("/look", handlers.LookActionHandler)
	router.HandleFunc("/like", handlers.LikeActionHandler)
	router.HandleFunc("/like/{user_id}", handlers.LikeActionHandler)

	router.HandleFunc("/ws", handlers.WebSocketHandler)
	router.HandleFunc("/chat", handlers.ManagerChatsHandler)
	router.HandleFunc("/chat/{chat_id}", handlers.ManagerChatsHandler)

	router.HandleFunc("/report", handlers.ReportUserHandler)

	log.Info("Listening ", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal("Server error: ", err)
	}
}
