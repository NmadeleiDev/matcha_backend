package server

import (
	"backend/server/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func StartServer(port string) {

	router := mux.NewRouter()

	router.HandleFunc("/signup", handlers.SignUpHandler)
	router.HandleFunc("/signin", handlers.SignInHandler)
	router.HandleFunc("/signout", handlers.SignOutHandler)
	router.HandleFunc("/user", handlers.UpdateAccountHandler)
	router.HandleFunc("/verify/{key}", handlers.VerifyAccountHandler)

	router.HandleFunc("/account", handlers.ManageOwnAccountHandler)
	router.HandleFunc("/actions/{action}", handlers.GetOwnActionsHandler)
	router.HandleFunc("/data/{id}", handlers.GetUserDataHandler) // получение данных любого юзера (только для залогиненных юзеров)
	router.HandleFunc("/ban", handlers.ManageBannedUsersHandler)

	router.HandleFunc("/tag", handlers.UserTagsHandler)

	router.HandleFunc("/media", handlers.GetUserOwnImagesHandler) // получить свои фотки

	router.HandleFunc("/strangers", handlers.GetStrangersHandler)

	router.HandleFunc("/look", handlers.LookActionHandler)
	router.HandleFunc("/like", handlers.LikeActionHandler)
	router.HandleFunc("/match", handlers.MatchHandler)

	router.HandleFunc("/ws", handlers.WebSocketHandler)

	log.Info("Listening ", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal("Server error: ", err)
	}
}
