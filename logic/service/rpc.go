/**
 *
 * @author liangjf
 * @create on 2020/6/3
 * @version 1.0
 */
package service

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/liangjfblue/gpusher/logic/common"

	"github.com/coreos/etcd/clientv3"

	pb "github.com/liangjfblue/gpusher/proto/gateway/rpc/v1"
	"google.golang.org/grpc"
)

var (
	ErrNoThisGatewayNode = errors.New("no this gateway node")
)

var (
	mu                 sync.RWMutex
	_gatewayRpcClients map[string]pb.GatewayClient
	_conns             map[string]*grpc.ClientConn
)

func init() {
	_gatewayRpcClients = make(map[string]pb.GatewayClient)
	_conns = make(map[string]*grpc.ClientConn)
}

func InitGatewayRpcClient(etcdAddr []string) error {
	gatewayHots, err := getAllGateway(etcdAddr)
	if err != nil {
		return err
	}

	if len(gatewayHots) <= 0 {
		//不需要退出
		return nil
	}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*20)
	defer cancel()

	for _, host := range gatewayHots {
		if err := addGatewayNode(ctx, host); err != nil {
			continue
		}
	}

	return nil
}

func CLoseRpcClient() {
	mu.Lock()
	defer mu.Unlock()
	for h, c := range _conns {
		_ = c.Close()
		delete(_conns, h)
	}
}

func GetGatewayRpcClient(host string) (pb.GatewayClient, error) {
	if _, ok := _gatewayRpcClients[host]; !ok {
		return nil, ErrNoThisGatewayNode
	}

	return _gatewayRpcClients[host], nil
}

func addGatewayNode(ctx context.Context, host string) error {
	cc, err := grpc.DialContext(ctx, host, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}

	_conns[host] = cc
	_gatewayRpcClients[host] = pb.NewGatewayClient(cc)
	return nil
}

func updateGatewayInfo(host string, op string) {
	mu.Lock()
	defer mu.Unlock()

	switch op {
	case "put":
		if _, ok := _conns[host]; !ok {
			ctx, cancel := context.WithTimeout(context.TODO(), time.Second*3)
			defer cancel()
			if err := addGatewayNode(ctx, host); err != nil {
				return
			}
		}
	case "delete":
		if v, ok := _conns[host]; ok {
			_ = v.Close()
		}

		delete(_conns, host)
		delete(_gatewayRpcClients, host)
	}
}

func getAllGateway(etcdAddr []string) (map[string]string, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints: etcdAddr,
	})
	if err != nil {
		return nil, err
	}

	kvs := make(map[string]string)
	resp, err := client.Get(context.TODO(), common.GatewayPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	for _, kv := range resp.Kvs {
		kvs[string(kv.Key)] = string(kv.Value)
	}

	go watcher(client, common.GatewayPrefix, kvs)

	return kvs, nil
}

func watcher(client *clientv3.Client, prefix string, kvs map[string]string) {
	wp := client.Watch(context.TODO(), prefix, clientv3.WithPrefix(), clientv3.WithPrevKV())
	for {
		for v := range wp {
			if v.Err() != nil {
				return
			}

			for _, event := range v.Events {
				switch event.Type {
				case clientv3.EventTypePut:
					kvs[string(event.Kv.Key)] = string(event.Kv.Value)
					updateGatewayInfo(string(event.Kv.Key), "put")
				case clientv3.EventTypeDelete:
					if _, ok := kvs[string(event.Kv.Key)]; ok {
						delete(kvs, string(event.Kv.Key))
						updateGatewayInfo(string(event.Kv.Key), "delete")
					}
				}
			}
		}
	}
}
