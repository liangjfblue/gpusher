/**
 *
 * @author liangjf
 * @create on 2020/6/4
 * @version 1.0
 */
package models

import (
	"github.com/chasex/redis-go-cluster"
	"github.com/liangjfblue/gpusher/common/db"
	"github.com/liangjfblue/gpusher/common/defind"
)

type RedisModel struct {
	pool *db.RedisPool

	nodes []string
}

func NewRedisModel(nodes []string) IModels {
	return &RedisModel{
		nodes: nodes,
	}
}

//初始化
func (m *RedisModel) Init() error {
	var err error
	m.pool, err = db.NewRedisCluster(m.nodes)
	return err
}

//SaveGatewayUUID 保存网关uuid映射
func (m *RedisModel) SaveGatewayUUID(uuid string, gatewayAddr string) error {
	return m.pool.HSet(defind.RedisKeyGatewayAllUUID, uuid, gatewayAddr)
}

//SaveAppUUID 保存App和uuid映射
func (m *RedisModel) SaveAppUUID(uuid string, appTag string) error {
	return m.pool.HSet(defind.RedisKeyGatewayAppUUID+":"+appTag, uuid, "")
}

//SaveExpireMsg 保存离线消息
func (m *RedisModel) SaveExpireMsg(uuid string, msgId string, msg string, expireTime int64) error {
	return nil
}

//DeleteGatewayUUID 删除网关uuid映射
func (m *RedisModel) DeleteGatewayUUID(uuid string) error {
	_, err := m.pool.HDel(defind.RedisKeyGatewayAllUUID, uuid)
	return err
}

//DeleteAppUUID 删除AppTag和uuid映射
func (m *RedisModel) DeleteAppUUID(uuid string, appTag string) error {
	_, err := m.pool.HDel(defind.RedisKeyGatewayAppUUID, appTag, uuid)
	return err
}

//DeleteExpireMsg 删除离线消息
func (m *RedisModel) DeleteExpireMsg(uuid string, msgId string) error {
	return nil
}

//获取网关uuid映射
func (m *RedisModel) GetGatewayUUID(uuid string) (string, error) {
	return redis.String(m.pool.HGet(defind.RedisKeyGatewayAllUUID, uuid))
}

//获取App和uuid映射
func (m *RedisModel) GetAppUUID(appTag string) ([]string, error) {
	return redis.Strings(m.pool.HGetAll(defind.RedisKeyGatewayAppUUID + ":" + appTag))
}
