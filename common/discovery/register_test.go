/**
 *
 * @author liangjf
 * @create on 2020/5/28
 * @version 1.0
 */
package discovery

import (
	"context"
	"testing"
	"time"
)

func TestRegister(t *testing.T) {
	e := NewRegister(
		[]string{"172.16.7.16:9002", "172.16.7.16:9004", "172.16.7.16:9006"},
		0,
	)

	if err := e.Register(context.TODO(), ServiceDesc{
		ServiceName: "grpc-etcd",
		Host:        "127.0.0.1",
		Port:        8899,
		TTL:         time.Second * 3,
	}); err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second * 3)

	if err := e.UnRegister(context.TODO(), ServiceDesc{
		ServiceName: "grpc-etcd",
	}); err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second * 5)
}
