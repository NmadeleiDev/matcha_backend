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

	log.Info("Listening ", port)
	if err := http.ListenAndServe(":" + port, router); err != nil {
		log.Fatal("Server error: ", err)
	}
}
