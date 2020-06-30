/**
 *
 * @author liangjf
 * @create on 2020/6/1
 * @version 1.0
 */
package rpc

import (
	"context"
	"errors"

	"github.com/liangjfblue/gpusher/common/defind"
	"github.com/liangjfblue/gpusher/common/logger/log"
	"github.com/liangjfblue/gpusher/gateway/common"
	"github.com/liangjfblue/gpusher/gateway/service/connect"
	pb "github.com/liangjfblue/gpusher/proto/gateway/rpc/v1"
)

type GatewayRPC struct {
}

//New 创建新的客户端channel, 用于推送
func (g *GatewayRPC) New(ctx context.Context, in *pb.NewRequest) (out *pb.Respond, err error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("ctx done")
	default:
	}

	if _, _, err = connect.GetClientChannel().New(int(in.AppId), in.UUid); err != nil {
		log.GetLogger(common.GatewayLog).Error("get userConn err:%s", err.Error())
		return nil, err
	}

	out = &pb.Respond{
		Code: 1,
		Msg:  "ok",
	}

	return out, err
}

//Close 关闭客户端channel
func (g *GatewayRPC) Close(ctx context.Context, in *pb.CloseRequest) (out *pb.Respond, err error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("ctx done")
	default:
	}

	userConn, err := connect.GetClientChannel().Get(int(in.AppId), in.UUid, false)
	if err != nil {
		log.GetLogger(common.GatewayLog).Error("get userConn err:%s", err.Error())
		return nil, err
	}

	if err := userConn.Close(); err != nil {
		log.GetLogger(common.GatewayLog).Error("close err:%s", err.Error())
		return nil, err
	}

	out = &pb.Respond{
		Code: 1,
		Msg:  "ok",
	}

	return out, err
}

//PushApp 推送给某用户
func (g *GatewayRPC) PushOne(ctx context.Context, in *pb.PushOneRequest) (out *pb.Respond, err error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("ctx done")
	default:
	}

	c, err := connect.GetClientChannel().Get(int(in.AppId), in.UUid, false)
	if err != nil {
		log.GetLogger(common.GatewayLog).Error("get userConn err:%s", err.Error())
		return nil, err
	}

	if err := c.PushMsg(int(in.AppId), in.UUid, []byte(in.Content)); err != nil {
		log.GetLogger(common.GatewayLog).Error("push one err:%s", err.Error())
		return nil, err
	}

	out = &pb.Respond{
		Code: 1,
		Msg:  "ok",
	}

	return out, err
}

//PushApp 推送给某App所有用户
func (g *GatewayRPC) PushApp(ctx context.Context, in *pb.PushAppRequest) (out *pb.Respond, err error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("ctx done")
	default:
	}

	cs, err := connect.GetClientChannel().GetApp(int(in.AppId))
	if err != nil {
		log.GetLogger(common.GatewayLog).Error("get userConn err:%s", err.Error())
		return nil, err
	}

	for _, conn := range cs {
		if err := conn.PushMsg(int(in.AppId), defind.GetApp(int(in.AppId)), []byte(in.Content)); err != nil {
			log.GetLogger(common.GatewayLog).Error("push app err:%s", err.Error())
			return nil, err
		}
	}

	out = &pb.Respond{
		Code: 1,
		Msg:  "ok",
	}

	return out, err
}

//PushAll 推送给所有App所有用户
func (g *GatewayRPC) PushAll(ctx context.Context, in *pb.PushAllRequest) (out *pb.Respond, err error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("ctx done")
	default:
	}

	cs, err := connect.GetClientChannel().GetAll()
	if err != nil {
		log.GetLogger(common.GatewayLog).Error("get userConn err:%s", err.Error())
		return nil, err
	}

	for _, conn := range cs {
		if err := conn.PushMsg(defind.AppAll, defind.GetApp(defind.AppAll), []byte(in.Content)); err != nil {
			log.GetLogger(common.GatewayLog).Error("push one err:%s", err.Error())
			return nil, err
		}
	}

	out = &pb.Respond{
		Code: 1,
		Msg:  "ok",
	}

	return out, err
}
