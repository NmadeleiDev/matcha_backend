package userDataStorage

import "backend/types"

var Manager types.UserDataStorage

func Init() {
	Manager = &ManagerStruct{}
	Manager.MakeConnection()
}
