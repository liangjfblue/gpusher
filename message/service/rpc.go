/**
 *
 * @author liangjf
 * @create on 2020/6/4
 * @version 1.0
 */
package service

import (
	"context"
	"errors"

	"github.com/liangjfblue/gpusher/message/models"

	pb "github.com/liangjfblue/gpusher/message/proto/rpc/v1"
)

var (
	ErrCtxDone = errors.New("ctx done")
)

//MessageRpc message rpc服务
type MessageRpc struct {
	models models.IModels
}

func NewMessageRpc(models models.IModels) *MessageRpc {
	return &MessageRpc{
		models: models,
	}
}

func (m *MessageRpc) SaveGatewayUUID(ctx context.Context, in *pb.SaveGatewayUUIDRequest) (out *pb.Respond, err error) {
	select {
	case <-ctx.Done():
		return nil, ErrCtxDone
	default:
	}

	defer func() {
		if err != nil {
			out.Code = 0
			out.Msg = "err: save gateway uuid"
		}
	}()

	out = &pb.Respond{
		Code: 1,
		Msg:  "ok",
	}

	if err = m.models.SaveGatewayUUID(in.UUID, in.GatewayAddr); err != nil {
		return
	}

	return
}

func (m *MessageRpc) SaveAppUUID(ctx context.Context, in *pb.SaveAppUUIDRequest) (out *pb.Respond, err error) {
	select {
	case <-ctx.Done():
		return nil, ErrCtxDone
	default:
	}

	defer func() {
		if err != nil {
			out.Code = 0
			out.Msg = "err: save app uuid"
		}
	}()

	out = &pb.Respond{
		Code: 1,
		Msg:  "ok",
	}

	if err = m.models.SaveAppUUID(in.UUID, in.AppTag); err != nil {
		return
	}

	return
}

func (m *MessageRpc) SaveExpireMsg(ctx context.Context, in *pb.SaveExpireMsgRequest) (out *pb.Respond, err error) {
	select {
	case <-ctx.Done():
		return nil, ErrCtxDone
	default:
	}

	//TODO
	return
}

func (m *MessageRpc) DeleteGatewayUUID(ctx context.Context, in *pb.DeleteGatewayUUIDRequest) (out *pb.Respond, err error) {
	select {
	case <-ctx.Done():
		return nil, ErrCtxDone
	default:
	}

	defer func() {
		if err != nil {
			out.Code = 0
			out.Msg = "err: delete gateway uuid"
		}
	}()

	out = &pb.Respond{
		Code: 1,
		Msg:  "ok",
	}

	if err = m.models.DeleteGatewayUUID(in.UUID); err != nil {
		return
	}

	return
}

func (m *MessageRpc) DeleteAppUUID(ctx context.Context, in *pb.DeleteAppUUIDRequest) (out *pb.Respond, err error) {
	select {
	case <-ctx.Done():
		return nil, ErrCtxDone
	default:
	}

	defer func() {
		if err != nil {
			out.Code = 0
			out.Msg = "err: delete app uuid"
		}
	}()

	out = &pb.Respond{
		Code: 1,
		Msg:  "ok",
	}

	if err = m.models.DeleteAppUUID(in.UUID, in.AppTag); err != nil {
		return
	}

	return
}

func (m *MessageRpc) DeleteExpireMsg(ctx context.Context, in *pb.DeleteExpireMsgRequest) (out *pb.Respond, err error) {
	select {
	case <-ctx.Done():
		return nil, ErrCtxDone
	default:
	}

	//TODO
	return nil, nil
}

func (m *MessageRpc) GetGatewayUUID(ctx context.Context, in *pb.GetGatewayUUIDRequest) (out *pb.GetGatewayUUIDRespond, err error) {
	select {
	case <-ctx.Done():
		return nil, ErrCtxDone
	default:
	}

	defer func() {
		if err != nil {
			out.Code = 0
			out.Msg = "err: delete app uuid"
			out.GatewayAddr = ""
		}
	}()

	var gatewayAddr string
	gatewayAddr, err = m.models.GetGatewayUUID(in.UUID)
	if err != nil {
		return
	}

	out = &pb.GetGatewayUUIDRespond{
		Code:        1,
		Msg:         "ok",
		GatewayAddr: gatewayAddr,
	}

	return
}

func (m *MessageRpc) GetAppUUID(ctx context.Context, in *pb.GetAppUUIDRequest) (out *pb.GetAppUUIDRespond, err error) {
	select {
	case <-ctx.Done():
		return nil, ErrCtxDone
	default:
	}

	defer func() {
		if err != nil {
			out.Code = 0
			out.Msg = "err: delete app uuid"
			out.UUIDs = nil
		}
	}()

	var gatewayAddrs []string
	gatewayAddrs, err = m.models.GetAppUUID(in.AppTag)
	if err != nil {
		return
	}

	out = &pb.GetAppUUIDRespond{
		Code:  1,
		Msg:   "ok",
		UUIDs: gatewayAddrs,
	}

	return
}
