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

	router.Post("/api/v1/signup", responce.Signup)
	//router.Route("/", func(r chi.Router) {
	//	r.Get("/", responce.Home)
	//})
	//router.Route("/signup", func(r chi.Router) {
	//	r.Get("/", responce.SignupPage)
	//	r.Post("/", responce.Signup)
	//})
	//router.Route("/login", func(r chi.Router) {
	//	r.Get("/", responce.LoginPage)
	//	r.Post("/", responce.Login)
	//})
	//router.Route("/verify", func(r chi.Router) {
	//	r.Get("/", responce.VerifyEmail)
	//})
	http.ListenAndServe(":8080", router)
}