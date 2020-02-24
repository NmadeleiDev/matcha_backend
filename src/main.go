package main

import (
	"Matcha/postgres"
	"Matcha/responce"
	"github.com/go-chi/chi"
	"net/http"
)

func main()  {
	postgres.InitTables()

	router := chi.NewRouter()
	router.Get("/", responce.Home)
	router.Route("/signup", func(r chi.Router) {
		r.Get("/", responce.SignupPage)
		r.Post("/", responce.Signup)
	})
	router.Route("/login", func(r chi.Router) {
		r.Get("/", responce.LoginPage)
		r.Post("/", responce.Login)
	})
	http.ListenAndServe(":3333", router)
}