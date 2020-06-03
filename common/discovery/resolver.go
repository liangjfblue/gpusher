/**
 *
 * @author liangjf
 * @create on 2020/5/28
 * @version 1.0
 */
package discovery

import (
	"context"
	"fmt"
	"log"
	"sync"

	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc/resolver"
)

const (
	scheme = "etcd"
)

type etcdBuilder struct {
	service   string
	endpoints []string
}

func NewEtcdBuilder(endpoints []string, service string) resolver.Builder {
	return &etcdBuilder{
		endpoints: endpoints,
		service:   service,
	}
}

func (e *etcdBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints: e.endpoints,
	})
	if err != nil {
		log.Fatal("grpc-etcd: ", err)
		return nil, err
	}

	ed := &etcdResolver{
		client: client,
		cc:     cc,
		rn:     make(chan struct{}, 1),
	}

	go ed.watcher(fmt.Sprintf("/%s/%s/", target.Scheme, e.service))
	ed.ResolveNow(resolver.ResolveNowOptions{})
	return ed, nil
}

//Scheme 命名空间
func (e *etcdBuilder) Scheme() string {
	return scheme
}

//etcdResolver etcd resolver
type etcdResolver struct {
	client *clientv3.Client
	cc     resolver.ClientConn
	rn     chan struct{}
	wg     sync.WaitGroup
}

//ResolveNow 主动通知更新服务可用列表, 这里用不到, 直接用etcd watch机制
func (r *etcdResolver) ResolveNow(resolver.ResolveNowOptions) {}

//Close 关闭Resolver
func (r *etcdResolver) Close() {}

//watcher 监听服务变更
func (r *etcdResolver) watcher(path string) {
	defer r.wg.Done()

	state := resolver.State{}
	kvs := make(map[string]string)

	update := func() {
		for k, v := range kvs {
			state.Addresses = append(state.Addresses,
				resolver.Address{
					ServerName: k,
					Addr:       v,
				},
			)
		}
		r.cc.UpdateState(state)
	}
	resp, err := r.client.Get(context.TODO(), path, clientv3.WithPrefix())
	if err != nil {
		log.Fatal("grpc-etcd: ", err)
		return
	}
	for _, kv := range resp.Kvs {
		kvs[string(kv.Key)] = string(kv.Value)
	}

	//first update
	update()

	wp := r.client.Watch(context.TODO(), path, clientv3.WithPrefix(), clientv3.WithPrevKV())
	for {
		for v := range wp {
			if v.Err() != nil {
				log.Fatal("grpc-etcd: ", err)
				return
			}

			for _, event := range v.Events {
				switch event.Type {
				case clientv3.EventTypePut:
					kvs[string(event.Kv.Key)] = string(event.Kv.Value)
				case clientv3.EventTypeDelete:
					if _, ok := kvs[string(event.Kv.Key)]; ok {
						delete(kvs, string(event.Kv.Key))
					}
				}
			}
			update()
		}
	}
}
