/**
 *
 * @author liangjf
 * @create on 2020/6/1
 * @version 1.0
 */
package db

import (
	"time"

	"github.com/chasex/redis-go-cluster"
)

type RedisPool struct {
	cluster *redis.Cluster
}

func NewRedisCluster(nodes []string) (*RedisPool, error) {
	p := new(RedisPool)

	var err error
	p.cluster, err = redis.NewCluster(
		&redis.Options{
			StartNodes:   nodes,
			ConnTimeout:  50 * time.Millisecond,
			ReadTimeout:  50 * time.Millisecond,
			WriteTimeout: 50 * time.Millisecond,
			KeepAlive:    16,
			AliveTime:    60 * time.Second,
		})
	if err != nil {
		return nil, err
	}

	return p, nil
}
