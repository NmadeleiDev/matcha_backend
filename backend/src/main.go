package main

import (
	"Matcha/postgres"
	"Matcha/api"
	"github.com/go-chi/chi"
	"github.com/rs/cors"
	"net/http"
)

func main()  {
	postgres.InitTables()

	router := chi.NewRouter()

	router.Post("/api/v1/signup", api.Signup)
	router.Post("/api/v1/signin", api.Login)
	//router.Route("/", func(r chi.Router) {
	//	r.Get("/", api.Home)
	//})
	//router.Route("/signup", func(r chi.Router) {
	//	r.Get("/", api.SignupPage)
	//	r.Post("/", api.Signup)
	//})
	//router.Route("/login", func(r chi.Router) {
	//	r.Get("/", api.LoginPage)
	//	r.Post("/", api.Login)
	//})
	//router.Route("/verify", func(r chi.Router) {
	//	r.Get("/", api.VerifyEmail)
	//})
	handler := cors.Default().Handler(router)
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowCredentials: true,
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})
	handler = c.Handler(handler)
	http.ListenAndServe(":2222", handler)
}