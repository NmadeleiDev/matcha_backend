package userFullDataStorage

import (
	"backend/model"

	"github.com/sirupsen/logrus"
)

func (m *ManagerStruct) GetUserData(loginData model.LoginData, isPublic bool) (model.FullUserData, error) {
	userData, err := Manager.GetFullUserData(loginData, isPublic)
	if err != nil {
		logrus.Error("Failed to get user data")
		return model.FullUserData{}, err
	} else {
		if userData.BannedUserIds == nil {
			userData.BannedUserIds = []string{}
		}
		if userData.LikedBy == nil {
			userData.LikedBy = []string{}
		}
		if userData.LookedBy == nil {
			userData.LookedBy = []string{}
		}
		if userData.Matches == nil {
			userData.Matches = []string{}
		}
		if userData.Tags == nil {
			userData.Tags = []string{}
		}
		if userData.Images == nil {
			userData.Images = []string{}
		}
		return userData, nil
	}
}
