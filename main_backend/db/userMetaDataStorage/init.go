package userMetaDataStorage

import (
	"backend/dao"
)

var Manager dao.UserMetaDataStorage

func Init() {
	postgres := ManagerStruct{}
	postgres.MakeConnection()
	postgres.InitTables()
	Manager = &postgres
}
