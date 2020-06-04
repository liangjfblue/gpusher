package proto

type ConnPayload struct {
	AppId int    `json:"appId"` //应用ID
	UUID  string `json:"uuid"`  //用户uuid
	Key   string `json:"key"`   //订阅key
	Token string `json:"token"` //推送token
}
