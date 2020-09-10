/**
 *
 * @author liangjf
 * @create on 2020/9/9
 * @version 1.0
 */
package models

import (
	"sync"

	"github.com/liangjfblue/gpusher/common/db"
)

var (
	onceDo sync.Once
	pool   *db.RedisPool
)

func InitRedisModel(nodes []string) error {
	var err error
	onceDo.Do(func() {
		pool, err = db.NewRedisCluster(nodes)
	})
	return err
}

func GetRedisPool() *db.RedisPool {
	return pool
}
