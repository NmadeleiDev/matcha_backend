package structuredDataStorage

import "backend/model"

var Manager model.StructuredDataStorage

func Init() {
	postgres := ManagerStruct{}
	postgres.MakeConnection()
	postgres.InitTables()
	Manager = &postgres
}
