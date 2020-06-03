/**
 *
 * @author liangjf
 * @create on 2020/6/3
 * @version 1.0
 */
package push

import "fmt"

type PushMsg struct {
	Tag  string  `json:"tag"`
	Body MsgBody `json:"body"`
}

type MsgBody struct {
	Type        int    `json:"type"` //推送类型(个体, 同一个app, 全体)
	UUID        string `json:"uuid"`
	Content     string `json:"content"`
	ExpireTime  uint32 `json:"expireTime"`
	OfflinePush bool   `json:"offlinePush"`
}

func (p PushMsg) String() string {
	return fmt.Sprintf(
		"tag:%s, body.type:%d, body.uuid:%s, body.content:%s, body.expireTime:%d, body.offlinePush:%v",
		p.Tag, p.Body.Type, p.Body.UUID, p.Body.Content, p.Body.ExpireTime, p.Body.OfflinePush)
}

const (
	Push2One = iota + 1
	Push2App
	Push2All
)

const (
	MaxExpireTime = 3600 * 24 * 7 //消息最大过期时间7天
)
