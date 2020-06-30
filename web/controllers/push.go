/**
 *
 * @author liangjf
 * @create on 2020/6/2
 * @version 1.0
 */
package controllers

import (
	"github.com/liangjfblue/gpusher/common/logger/log"
	"github.com/liangjfblue/gpusher/common/push"
	"github.com/liangjfblue/gpusher/web/common"

	"github.com/gin-gonic/gin"
	"github.com/liangjfblue/gpusher/web/service"
)

//PushMsg 接收推送消息, 写入消息队列, logic消费消息转发给网关客户端
func PushMsg(c *gin.Context) {
	var (
		err    error
		result Result
		req    push.PushMsg
	)

	if err = c.BindJSON(&req); err != nil {
		result.Failure(c, ErrPushMsg)
		return
	}

	if err := checkParam(&req); err != nil {
		result.Failure(c, err)
		return
	}

	msg := push.PushMsg{
		Tag:  req.Tag,
		Body: req.Body,
	}
	service.GetPush("kafka").Push(&msg)

	log.GetLogger(common.WebLog).Debug("gpusher: web send msg: %s", msg)

	result.Success(c, nil)
}

func checkParam(req *push.PushMsg) error {
	if req.Tag == "" {
		return ErrPushMsgTagEmpty
	}

	if req.Body.Type < push.Push2One || req.Body.Type > push.Push2All {
		return ErrPushMsgBodyTypeEmpty
	}

	if req.Body.UUID == "" {
		return ErrPushMsgBodyUUIDEmpty
	}

	if req.Body.Content == "" {
		return ErrPushMsgBodyContentEmpty
	}

	if req.Body.ExpireTime > push.MaxExpireTime {
		return ErrPushMsgBodyExpireTimeOver
	}

	return nil
}
