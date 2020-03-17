package api

import (
	"Matcha/postgres"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type UseLoginData struct {
	Username	string	`json:"username"`
	Password	string	`json:"password"`
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/signup/")
}

func Login(w http.ResponseWriter, r *http.Request)  {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err.Error())
	}
	userData := new(UseLoginData)
	err = json.Unmarshal(data, &userData)
	if err != nil {
		w.Write([]byte("Backend error: " + err.Error()))
		log.Println("Error unmarshalling json: ", err)
	} else {
		log.Println("Parsed json: ", userData)
	}
	if postgres.AuthUser(userData.Username, userData.Password) == true {
		fmt.Println("User logged in")
		w.Write([]byte("OK"))
	} else {
		fmt.Println("Logging failed")
		w.Write([]byte("KO"))
	}
}
