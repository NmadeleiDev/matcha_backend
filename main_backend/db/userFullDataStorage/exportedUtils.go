package userFullDataStorage

import (
	"backend/db/userMetaDataStorage"
	"backend/model"

	"github.com/sirupsen/logrus"
)

func (m *ManagerStruct) GetUserData(loginData model.LoginData, isPublic bool) (model.FullUserData, error) {
	var variant string

	if isPublic {
		variant = "public"
	} else {
		variant = "private"
	}
	userData, err := Manager.GetFullUserData(loginData, variant)
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
		userData.Tags = userMetaDataStorage.Manager.GetTagsById(userData.TagIds)
		userData.ConvertFromDbCoords()
		return userData, nil
	}
}
