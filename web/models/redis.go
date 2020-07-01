/**
 *
 * @author liangjf
 * @create on 2020/7/1
 * @version 1.0
 */
package models

import (
	"sync"

	"github.com/liangjfblue/gpusher/common/db"
)

var (
	_pool  *db.RedisPool
	onceDo sync.Once
)

func InitRedisPool(nodes []string) error {
	var err error
	onceDo.Do(func() {
		_pool, err = db.NewRedisCluster(nodes)
	})
	return err
}

func GetRedisPool() *db.RedisPool {
	return _pool
}
