package utils

import (
	"backend/types"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func UnmarshalHttpBodyToUserData(w http.ResponseWriter, r *http.Request) (*types.FullUserData, bool) {
	container := types.FullUserData{}
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

func UnmarshalHttpBodyToLoginData(w http.ResponseWriter, r *http.Request) (*types.LoginData, bool) {
	container := types.LoginData{}
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

func UnmarshalHttpBodyToTags(w http.ResponseWriter, r *http.Request) (*types.Tags, bool) {
	container := types.Tags{}
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
