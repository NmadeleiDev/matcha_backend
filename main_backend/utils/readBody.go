package utils

import (
	"backend/model"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func UnmarshalHttpBodyToUserData(w http.ResponseWriter, r *http.Request) (*model.FullUserData, bool) {
	container := model.FullUserData{}
	requestData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Error("Can't read request body: ", err)
		SendFailResponse(w, "error reading body")
		return nil, false
	}
	err = json.Unmarshal(requestData, &container)
	if err != nil {
		logrus.Error("Can't read request body: ", err)
		SendFailResponse(w, "error reading body")
		return nil, false
	}
	return &container, true
}

func UnmarshalHttpBodyToLoginData(w http.ResponseWriter, r *http.Request) (*model.LoginData, bool) {
	container := model.LoginData{}
	requestData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Error("Can't read request body: ", err)
		SendFailResponse(w, "error reading body")
		return nil, false
	}
	err = json.Unmarshal(requestData, &container)
	if err != nil {
		logrus.Error("Can't read request body: ", err)
		SendFailResponseWithCode(w, "error reading body", http.StatusInternalServerError)
		return nil, false
	}
	return &container, true
}

func UnmarshalHttpBodyToTags(w http.ResponseWriter, r *http.Request) (*model.Tags, bool) {
	container := model.Tags{}
	requestData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Error("Can't read request body: ", err)
		SendFailResponse(w, "error reading body")
		return nil, false
	}
	err = json.Unmarshal(requestData, &container)
	if err != nil {
		logrus.Error("Can't read request body: ", err)
		SendFailResponse(w, "error reading body")
		return nil, false
	}
	return &container, true
}

func UnmarshalHttpBodyToChat(w http.ResponseWriter, r *http.Request) (*model.Chat, bool) {
	container := model.Chat{}
	requestData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Error("Can't read request body: ", err)
		SendFailResponse(w, "error reading body")
		return nil, false
	}
	err = json.Unmarshal(requestData, &container)
	if err != nil {
		logrus.Error("Can't read request body: ", err)
		SendFailResponse(w, "error reading body")
		return nil, false
	}
	return &container, true
}

func UnmarshalHttpBodyToReport(w http.ResponseWriter, r *http.Request) (*model.Report, bool) {
	container := model.Report{}
	requestData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Error("Can't read request body: ", err)
		SendFailResponse(w, "error reading body")
		return nil, false
	}
	err = json.Unmarshal(requestData, &container)
	if err != nil {
		logrus.Error("Can't read request body: ", err)
		SendFailResponse(w, "error reading body")
		return nil, false
	}
	return &container, true
}

