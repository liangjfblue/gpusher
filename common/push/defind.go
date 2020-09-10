/**
 *
 * @author liangjf
 * @create on 2020/6/3
 * @version 1.0
 */
package push

import (
	"encoding/json"
	"fmt"
)

type PushMsg struct {
	Tag  string  `json:"tag"`
	Body MsgBody `json:"body"`
}

type MsgBody struct {
	Type        int    `json:"type"`   //推送类型(个体, 同一个app, 全体)
	MsgSeq      uint64 `json:"msgSeq"` //消息序号, 用于ack和持久化
	UUID        string `json:"uuid"`
	Content     string `json:"content"`
	ExpireTime  uint32 `json:"expireTime"`
	OfflinePush bool   `json:"offlinePush"`
}

func (p PushMsg) String() string {
	return fmt.Sprintf(
		"tag:%s, body.type:%d, body.type:%d, body.uuid:%s, body.content:%s, body.expireTime:%d, body.offlinePush:%v",
		p.Tag, p.Body.Type, p.Body.MsgSeq, p.Body.UUID, p.Body.Content, p.Body.ExpireTime, p.Body.OfflinePush)
}

func (m MsgBody) String() string {
	d, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(d)
}

//推送类型
const (
	Push2One = iota + 1
	Push2App
	Push2All
)

const (
	MaxExpireTime = 3600 * 24 * 7 //消息最大过期时间7天
)

//app应用列表
//TODO save in etcd
const (
	AppGpusher = iota + 1000
)

var (
	AppM = map[string]int32{
		"app_gpusher": AppGpusher,
	}
)
