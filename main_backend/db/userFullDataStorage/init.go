package userFullDataStorage

import (
	"backend/dao"
)

var Manager dao.UserFullDataStorage

func Init() {
	Manager = &ManagerStruct{}
	Manager.MakeConnection()
	Manager.CreateUniqueIndex("username")
	Manager.CreateLocationIndex()
}
