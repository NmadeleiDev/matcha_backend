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
	router.Route("/", func(r chi.Router) {
		r.Get("/", responce.Home)
	})
	router.Route("/signup", func(r chi.Router) {
		r.Get("/", responce.SignupPage)
		r.Post("/", responce.Signup)
	})
	router.Route("/login", func(r chi.Router) {
		r.Get("/", responce.LoginPage)
		r.Post("/", responce.Login)
	})
	router.Route("/verify", func(r chi.Router) {
		r.Get("/", responce.VerifyEmail)
	})
	m := &autocert.Manager{
		Cache:      autocert.DirCache("golang-autocert"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("example.org", "www.example.org"),
	}
	server := &http.Server{
		Addr:      ":3333",
		TLSConfig: m.TLSConfig(),
	}
	server.ListenAndServeTLS("", "")
	http.ListenAndServe(":3333", router)
}