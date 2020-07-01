/**
 *
 * @author liangjf
 * @create on 2020/7/1
 * @version 1.0
 */
package controllers

import (
	"fmt"
	"hash/crc32"
	"strings"

	"github.com/liangjfblue/gpusher/web/models"

	"github.com/gin-gonic/gin"
	"github.com/liangjfblue/gpusher/common/defind"
	"github.com/liangjfblue/gpusher/common/token"
)

type TokenReq struct {
	UUID     string `json:"uuid"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenResp struct {
	Token       string `json:"token"`
	GatewayAddr string `json:"gatewayAddr"`
}

//Token 用户登录前获取token, 用于后续的推送凭证
func Token(c *gin.Context) {
	var (
		err    error
		result Result
		req    TokenReq
		resp   TokenResp
	)

	if err = c.BindJSON(&req); err != nil {
		result.Failure(c, ErrPushMsg)
		return
	}

	//TODO 查找数据库,判断用户
	if req.Username != "test" && req.Password != "123456" {
		result.Failure(c, ErrUserNotFound)
		return
	}

	tk, err := token.SignToken(token.Context{UUID: req.UUID})
	if err != nil {
		result.Failure(c, ErrSignToken)
		return
	}

	//分片, 减小redis hash大小, 提高查询效率(单个元素最大值为 512 MB，推荐元素个数小于 8192， value 最大长度不超过 1 MB)
	index := crc32.ChecksumIEEE([]byte(req.UUID)) % 64
	key := defind.RedisKeyUUIDToken + fmt.Sprint(index)
	if err := models.GetRedisPool().HSet(key, req.UUID, tk); err != nil {
		result.Failure(c, ErrSaveToken2Redis)
		return
	}

	//负载均衡获取一个gateway给客户端
	gatewayAddr, err := models.GetAllGateway()
	if err != nil {
		result.Failure(c, ErrGatewayEmpty)
		return
	}

	resp.Token = tk
	resp.GatewayAddr = strings.Join(gatewayAddr, ",")
	result.Success(c, resp)
}
