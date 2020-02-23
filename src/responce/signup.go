package responce

import (
	"Matcha/postgres"
	"Matcha/utils"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
	"time"
)

func SignupPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/signup/")
}

func Signup(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("passwd1")
	passwordConfirm := r.FormValue("passwd2")
	if strings.Compare(password, passwordConfirm) != 0 {
		w.Write([]byte("User: " + username + "  Password: " + password))
		return
	}

	uniqueKey := utils.GetMD5(username + password + time.Now().String())
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(" InsertUserData error. Err: " + err.Error())
	}
	postgres.InsertUserData(username, string(hashedPassword), uniqueKey)
	w.Write([]byte("Account " + username + " successfully created"))
	fmt.Println("Account created")

}