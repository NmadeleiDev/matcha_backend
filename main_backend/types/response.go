package types

type ResponseJson struct {
	Status bool        `json:"status"`
	Data   interface{} `json:"data"`
}
