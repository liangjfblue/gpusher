/**
 *
 * @author liangjf
 * @create on 2020/6/3
 * @version 1.0
 */
package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

//-----------------------------------------//
type Result struct {
	Code int         `json:"Code"`
	Data interface{} `json:"Data,omitempty"`
	Msg  string      `json:"Msg"`
}

func (r *Result) Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Result{
		Code: 1,
		Data: data,
		Msg:  SuccessMsg,
	})
}

func (r *Result) Failure(c *gin.Context, err error) {
	if e, ok := err.(*Errno); ok {
		c.JSON(http.StatusOK, Result{
			Code: 0,
			Data: map[string]interface{}{
				"code": e.Code,
				"msg":  e.Msg,
			},
			Msg: "error",
		})
	} else {
		c.JSON(http.StatusOK, Result{
			Code: 0,
			Data: map[string]interface{}{
				"code": -1,
				"msg":  "system error",
			},
			Msg: "error",
		})
	}
}

//----------------------------------------//
var (
	_codes = map[int]struct{}{}
)

func New(e int) int {
	if e <= 0 {
		panic("code must greater than zero")
	}
	return add(e)
}

func add(e int) int {
	if _, ok := _codes[e]; ok {
		panic(fmt.Sprintf("code: %d already exist", e))
	}
	_codes[e] = struct{}{}
	return e
}

type Errno struct {
	Code int
	Msg  string
}

func (e Errno) Error() string {
	return e.Msg
}

var (
	SuccessMsg = "ok"

	ErrPushMsg                   = Errno{Code: New(100), Msg: "推送失败"}
	ErrPushMsgTagEmpty           = Errno{Code: New(101), Msg: "缺失推送应用tag"}
	ErrPushMsgBodyUUIDEmpty      = Errno{Code: New(102), Msg: "缺失推送uuid"}
	ErrPushMsgBodyTypeEmpty      = Errno{Code: New(103), Msg: "缺失推送类型"}
	ErrPushMsgBodyContentEmpty   = Errno{Code: New(104), Msg: "缺失推送内容"}
	ErrPushMsgBodyExpireTimeOver = Errno{Code: New(105), Msg: "超过最大过期时间"}

	ErrSignToken       = Errno{Code: New(150), Msg: "生成token失败"}
	ErrUserNotFound    = Errno{Code: New(151), Msg: "用户不存在/密码错误"}
	ErrSaveToken2Redis = Errno{Code: New(152), Msg: "生成token失败"}
	ErrGatewayEmpty    = Errno{Code: New(153), Msg: "网关列表空"}
)
