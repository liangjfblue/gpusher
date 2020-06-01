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
	red "github.com/garyburd/redigo/redis"
)

type RedisPool struct {
	pool    *red.Pool
	cluster redis.Cluster
}

func NewRedisPool(redisURL string, redisMaxIdle int, redisIdleTimeoutSec time.Duration) *RedisPool {
	p := new(RedisPool)

	p.pool = &red.Pool{
		MaxIdle:     redisMaxIdle,
		IdleTimeout: redisIdleTimeoutSec * time.Second,
		Dial: func() (red.Conn, error) {
			c, err := red.Dial("tcp", redisURL)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c red.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	return p
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
			KeepAlive:    200,
			AliveTime:    60 * time.Second,
		})

	if err != nil {
		return nil, err
	}

	return p, nil
}
