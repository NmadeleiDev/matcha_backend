package dto

import (
	"fmt"
	"math/rand"

	"backend/db/realtimeDataDb"
	"backend/db/userFullDataStorage"
	"backend/db/userMetaDataStorage"
	"backend/model"
	"backend/utils"

	"github.com/sirupsen/logrus"
)

func GetUserDTO(user *model.FullUserData) *UserDTO {
	if len(user.Id) == 0 {
		logrus.Errorf("User passed for dto is empty: %v", user)
		return &UserDTO{err: fmt.Errorf("user passed for dto is empty")}
	}
	return &UserDTO{user: user}
}

type UserDTO struct {
	user	*model.FullUserData
	err		error
}

func (d *UserDTO) GetUser() *model.FullUserData {
	if d.err != nil {
		logrus.Errorf("dto is invalid: %v", d.err)
	}
	return d.user
}

func (d *UserDTO) GetError() error {
	return d.err
}

func (d *UserDTO) LoadUserData(isPublic bool) *UserDTO {
	if d.err != nil {
		logrus.Errorf("dto is invalid: %v", d.err)
		return d
	}
	data, err := userFullDataStorage.Manager.GetUserData(model.LoginData{Id: d.user.Id}, isPublic)
	if err != nil {
		d.err = err
		return d
	} else {
		data.IsOnline = realtimeDataDb.GetManager().IsUserOnline(d.user.Id)
		d.user = &data
	}
	return d
}

func (d *UserDTO) PrepareUserDataForClient() *UserDTO {
	if d.err != nil {
		logrus.Errorf("dto is invalid: %v", d.err)
		return d
	}

	d.user.Tags = userMetaDataStorage.Manager.GetTagsById(d.user.TagIds)
	d.user.ConvertFromDbCoords()

	d.user.Tags = userMetaDataStorage.Manager.GetTagsById(d.user.TagIds)
	if len(d.user.Avatar) == 0 && len(d.user.Images) > 0 {
		d.user.Avatar = d.user.Images[rand.Intn(len(d.user.Images))]
	}
	d.user.Rating = utils.Sigmoid(d.user.Rating)

	d.user.IsOnline = realtimeDataDb.GetManager().IsUserOnline(d.user.Id)
	return d
}
