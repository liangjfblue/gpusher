/**
 *
 * @author liangjf
 * @create on 2020/6/4
 * @version 1.0
 */
package api

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/liangjfblue/gpusher/common/logger/log"

	"github.com/liangjfblue/gpusher/gateway/common"
	"github.com/liangjfblue/gpusher/logic/api"

	pb "github.com/liangjfblue/gpusher/proto/message/rpc/v1"

	"github.com/liangjfblue/gpusher/common/discovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

var (
	_MessageRpcClient pb.MessageClient
	_conns            *grpc.ClientConn
	_etcdAddr         []string
)

func InitMessageClientRpc(ctx context.Context, etcdAddr []string, serviceName string) error {
	_etcdAddr = etcdAddr

	rl := discovery.NewEtcdBuilder(etcdAddr, serviceName)
	resolver.Register(rl)

	conn, err := grpc.DialContext(
		ctx,
		fmt.Sprintf("%s:///%s", rl.Scheme(), serviceName),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return err
	}

	_MessageRpcClient = pb.NewMessageClient(conn)
	_conns = conn

	return nil
}

//CloseRpcClient release message rpc client
func CloseRpcClient() {
	_ = _conns.Close()
}

//GetMessageRpcClient get message rpc client
func GetMessageRpcClient() pb.MessageClient {
	return _MessageRpcClient
}

//ReBalanceMessageRpcClient balance message rpc client
func ReBalanceMessageRpcClient() error {
	log.GetLogger(common.GatewayLog).Debug("reconnect message rpc")

	CloseRpcClient()

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*4)
	defer cancel()
	for i := 0; i < 3; i++ {
		if err := api.InitMessageClientRpc(ctx, _etcdAddr, common.MessageServiceName); err == nil {
			log.GetLogger(common.GatewayLog).Debug("reconnect message rpc ok")
			return nil
		}
	}
	return errors.New("gateway rpc to message err")
}
