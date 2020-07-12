package structuredDataStorage

import "backend/types"

var Manager types.StructuredDataStorage

func Init() {
	postgres := ManagerStruct{}
	postgres.MakeConnection()
	postgres.InitTables()
	Manager = &postgres
}
