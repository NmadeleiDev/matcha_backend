package responce

import (
	"Matcha/postgres"
	"fmt"
	"net/http"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/signup/")
}

func Login(w http.ResponseWriter, r *http.Request)  {
	email := r.FormValue("email")
	password := r.FormValue("passwd")

	if postgres.AuthUser(email, password) == true {
		fmt.Println("User logged in")
		http.Redirect(w, r, "/home", 200)
	}
}