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

	router.HandleFunc("/data/{id}", handlers.GetUserDataHandler) // получение данных любого юзера (только для залогиненных юзеров)

	router.HandleFunc("/media", handlers.GetUserOwnImagesHandler) // получить свои фотки

	router.HandleFunc("/strangers", handlers.GetStrangersHandler)

	router.HandleFunc("/look", handlers.SaveAccountLookUpHandler)
	router.HandleFunc("/like", handlers.SaveLikeActionHandler)
	router.HandleFunc("/match", handlers.SaveMatchHandler)

	router.HandleFunc("/ws", handlers.WebSocketHandler)

	log.Info("Listening ", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal("Server error: ", err)
	}
}
