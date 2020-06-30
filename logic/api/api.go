/**
 *
 * @author liangjf
 * @create on 2020/6/4
 * @version 1.0
 */
package api

import (
	"context"
	"fmt"

	pb "github.com/liangjfblue/gpusher/proto/message/rpc/v1"

	"github.com/liangjfblue/gpusher/common/discovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

var (
	_MessageRpcClient pb.MessageClient
	_conns            []*grpc.ClientConn
)

func init() {
	_conns = make([]*grpc.ClientConn, 0)
}

func InitMessageClientRpc(ctx context.Context, etcdAddr []string, serviceName string) error {
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

	_conns = append(_conns, conn)

	return nil
}

func CloseRpcClient() {
	for _, conn := range _conns {
		_ = conn.Close()
	}
}

func GetMessageRpcClient() pb.MessageClient {
	return _MessageRpcClient
}
