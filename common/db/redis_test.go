/**
 *
 * @author liangjf
 * @create on 2020/6/1
 * @version 1.0
 */
package db

import (
	"testing"
)

func TestNewRedisPool(t *testing.T) {
	//p := NewRedisPool("172.16.7.16:8001", 10, time.Duration(300))
	p, err := NewRedisCluster([]string{"172.16.7.16:8001", "172.16.7.16:8002", "172.16.7.16:8003", "172.16.7.16:8004", "172.16.7.16:8005", "172.16.7.16:8006"})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(p.Set("name", "liangjf"))

	t.Log(p.Get("name"))
}
