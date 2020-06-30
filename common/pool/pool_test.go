/**
 *
 * @author liangjf
 * @create on 2020/6/8
 * @version 1.0
 */
package pool

import (
	"context"
	"net"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	p, err := NewPool(
		WithInitCap(5),
		WithMaxCap(20),
		WithBuilder(func(ctx context.Context) (net.Conn, error) {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}

			timeout := 200 * time.Millisecond
			if t, ok := ctx.Deadline(); ok {
				timeout = time.Until(t) //t.Sub(time.Now())
			}

			return net.DialTimeout("tcp", "172.16.7.17:9090", timeout)
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()

	conn, err := p.Get()
	if err != nil {
		t.Fatal(err)
	}

	n, err := conn.Write([]byte("hello"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(n)

	recv := make([]byte, 64)
	n, err = conn.Read(recv)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(recv[:n]))
}
