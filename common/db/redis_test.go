/**
 *
 * @author liangjf
 * @create on 2020/6/1
 * @version 1.0
 */
package db

import (
	"testing"

	"github.com/chasex/redis-go-cluster"
)

func TestNewRedisPool(t *testing.T) {
	p, err := NewRedisCluster([]string{"172.16.7.16:8001", "172.16.7.16:8002", "172.16.7.16:8003"})
	if err != nil {
		t.Fatal(err)
	}

	if _, err = p.Get("name"); err != nil {
		t.Fatal(err)
	} else {
		t.Log(redis.String(p.Get("name")))
	}
}
