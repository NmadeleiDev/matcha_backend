package userDataStorage

import "backend/model"

var Manager model.UserDataStorage

func Init() {
	Manager = &ManagerStruct{}
	Manager.MakeConnection()
}
