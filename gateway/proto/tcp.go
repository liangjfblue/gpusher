package proto

type ConnPayload struct {
	AppId int    `json:"appId"`
	Key   string `json:"key"`
	Token string `json:"token"`
}
