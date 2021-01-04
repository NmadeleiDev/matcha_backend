package utils

import (
	"encoding/base64"
	"encoding/json"
	"math"
	"math/rand"
	"net/http"
	"strconv"

	"backend/model"

	log "github.com/sirupsen/logrus"
)

const (
	oneDayInSeconds = 86400
)

func SetCookieForDay(w http.ResponseWriter, cookieName, value string) {
	c := http.Cookie{
		Name:  cookieName,
		Value: value,
		Path:  "/",
		//SameSite: http.SameSiteNoneMode,
		//Secure: true,
		MaxAge: oneDayInSeconds * 1}
	http.SetCookie(w, &c)
}

func SetHttpOnlyCookie(w http.ResponseWriter, cookieName, value string) {
	c := http.Cookie{
		Name:  cookieName,
		Value: value,
		Path:  "/",
		HttpOnly: true,
		//Secure: true,
		MaxAge: 360000}
	http.SetCookie(w, &c)
}

func GetCookieValue(r *http.Request, cookieName string) string {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		log.Info("Failed getting cookie: ", err)
		return ""
	} else {
		//log.Info("Got cookie: ", cookie)
	}
	return cookie.Value
}

func SendFailResponse(w http.ResponseWriter, text string) {
	var packet []byte
	var err error

	response := &model.ResponseJson{Status: false, Data: text}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	if packet, err = json.Marshal(response); err != nil {
		log.Error("Error marshalling response: ", err)
	}
	if _, err = w.Write(packet); err != nil {
		log.Error("Error sending response: ", err)
	}
}

func SendFailResponseWithCode(w http.ResponseWriter, text string, code int) {
	var packet []byte
	var err error

	response := &model.ResponseJson{Status: false, Data: text}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)

	if packet, err = json.Marshal(response); err != nil {
		log.Error("Error marshalling response: ", err)
	}
	if _, err = w.Write(packet); err != nil {
		log.Error("Error sending response: ", err)
	}
}

func SendSuccessResponse(w http.ResponseWriter) {
	var packet []byte
	var err error

	response := &model.ResponseJson{Status: true, Data: nil}
	w.Header().Set("content-type", "application/json")

	if packet, err = json.Marshal(response); err != nil {
		log.Error("Error marshalling response: ", err)
	}
	if _, err = w.Write(packet); err != nil {
		log.Error("Error sending response: ", err)
	}
}

func SendDataResponse(w http.ResponseWriter, data interface{}) {
	var packet []byte
	var err error

	response := &model.ResponseJson{Status: true, Data: data}
	w.Header().Set("content-type", "application/json")

	if packet, err = json.Marshal(response); err != nil {
		log.Error("Error marshalling response: ", err)
	}
	if _, err = w.Write(packet); err != nil {
		log.Error("Error sending response: ", err)
	}
}

func DoesArrayContain(haystack []string, needle string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Sigmoid(src float64) float64 {
	return 1.0 / (1.0 + math.Exp(-src))
}

func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func UnsafeAtoi(src string, alt int) int {
	if res, err := strconv.Atoi(src); err != nil {
		return alt
	} else {
		return res
	}
}

func UnsafeAtof(src string, alt float64) float64 {
	if res, err := strconv.ParseFloat(src, 64); err != nil {
		return alt
	} else {
		return res
	}
}
