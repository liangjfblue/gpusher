/**
 *
 * @author liangjf
 * @create on 2020/7/1
 * @version 1.0
 */
package message

import (
	"context"
	"errors"
	"fmt"

	"github.com/liangjfblue/gpusher/common/logger/log"
	"github.com/liangjfblue/gpusher/gateway/common"

	"github.com/liangjfblue/gpusher/gateway/api"
	pb "github.com/liangjfblue/gpusher/proto/message/rpc/v1"
)

func SaveGatewayUUID(uuid, gatewayAddr string) error {
	var err error
	//失败重试
	for i := 0; i < 3; i++ {

		if _, err = api.GetMessageRpcClient().SaveGatewayUUID(
			context.TODO(),
			&pb.SaveGatewayUUIDRequest{
				UUID:        uuid,
				GatewayAddr: gatewayAddr,
			}); err == nil {
			return nil
		}

		//故障转移failover, 重新和任意message连接
		if err = api.ReBalanceMessageRpcClient(); err == nil {
			return err
		}
	}

	if err != nil {
		return errors.New(fmt.Sprintf("SaveGatewayUUID: over 3 times"))
	}
	return nil
}

func DeleteGatewayUUID(uuid string) error {
	var err error
	//失败重试
	for i := 0; i < 3; i++ {
		if _, err = api.GetMessageRpcClient().DeleteGatewayUUID(
			context.TODO(),
			&pb.DeleteGatewayUUIDRequest{
				UUID: uuid,
			}); err == nil {
			return nil
		} else {
			log.GetLogger(common.GatewayLog).Error("DeleteGatewayUUID err:%s", err.Error())
		}

		//故障转移failover, 重新和任意message连接
		if err = api.ReBalanceMessageRpcClient(); err == nil {
			return err
		}
	}

	if err != nil {
		return errors.New(fmt.Sprintf("DeleteGatewayUUID: over 3 times"))
	}
	return nil
}

func ExpireGatewayUUID(uuid string) error {
	var err error
	//失败重试
	for i := 0; i < 3; i++ {
		//TODO ExpireGatewayUUID
		continue
		//故障转移failover, 重新和任意message连接
		if err = api.ReBalanceMessageRpcClient(); err == nil {
			return err
		}
	}

	if err != nil {
		return errors.New(fmt.Sprintf("ExpireGatewayUUID: over 3 times"))
	}
	return nil
}
