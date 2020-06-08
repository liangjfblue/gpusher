/**
 *
 * @author liangjf
 * @create on 2020/6/3
 * @version 1.0
 */
package service

import (
	"context"

	pb "github.com/liangjfblue/gpusher/message/proto/rpc/v1"

	"github.com/liangjfblue/gpusher/common/push"
	"github.com/liangjfblue/gpusher/logic/api"
)

//router 查找uuid所在的gateway
func router(ctx context.Context, msg *push.PushMsg) (string, error) {
	resp, err := api.GetMessageRpcClient().GetGatewayUUID(ctx, &pb.GetGatewayUUIDRequest{
		UUID: msg.Body.UUID,
	})
	if err != nil {
		return "", err
	}

	return resp.GatewayAddr, nil
}
