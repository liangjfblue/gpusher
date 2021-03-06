/**
 *
 * @author liangjf
 * @create on 2020/6/3
 * @version 1.0
 */
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/liangjfblue/gpusher/logic/models"

	"github.com/liangjfblue/gpusher/common/logger/log"
	"github.com/liangjfblue/gpusher/logic/common"

	pb "github.com/liangjfblue/gpusher/proto/gateway/rpc/v1"

	"github.com/liangjfblue/gpusher/common/push"
)

var (
	_receiver push.IQueueReceiver
)

func InitKafkaConsumer(ctx context.Context, brokerAddr []string) error {
	//消息队列
	_receiver = push.NewKafkaReceiver(brokerAddr)
	go func() {
		if err := _receiver.Init(); err != nil {
			return
		}
	}()

	go func() {
		if err := _receiver.Recv(dealMsg); err != nil {
			return
		}
	}()

	return nil
}

func StopKafkaConsumer() {
	_receiver.Stop()
}

func dealMsg(msg []byte) {
	var m push.PushMsg
	if err := json.Unmarshal(msg, &m); err != nil {
		log.GetLogger(common.LogicLog).Error("json msg error")
		return
	}

	//msg router to user, where gateway
	host, err := router(context.TODO(), &m)
	if err != nil {
		log.GetLogger(common.LogicLog).Error("uuid no gateway node, err:%s", err.Error())
		return
	}

	rpcClient, err := GetGatewayRpcClient(host)
	if err != nil {
		log.GetLogger(common.LogicLog).Error("get gateway rpc client error:%s", err.Error())
		return
	}

	if _, ok := push.AppM[m.Tag]; !ok {
		log.GetLogger(common.LogicLog).Error("no this app tag:%s", m.Tag)
		return
	}

	//每个消息打序号
	m.Body.MsgSeq = models.GetRedisPool().GenerateMsgSeq(m.Body.UUID)

	//TODO 持久化 add to message table
	if m.Body.OfflinePush {

	}

	//TODO 消息多久过期, 默认不过期
	if m.Body.ExpireTime > 0 {

	}

	//TODO 推入已发送待确认队列infightqueue

	switch m.Body.Type {
	case push.Push2One:
		err = pushOne(rpcClient, &m)
	case push.Push2App:
		err = pushApp(rpcClient, &m)
	case push.Push2All:
		err = pushAll(rpcClient, &m)
	default:
	}

	if err != nil {
		log.GetLogger(common.LogicLog).Error("gpusher: push err: %s", err.Error())
	}
}

func pushOne(rpcClient pb.GatewayClient, m *push.PushMsg) error {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()

	log.GetLogger(common.LogicLog).Error("gpusher: logic pushOne: %s", *m)

	if _, err := rpcClient.PushOne(ctx, &pb.PushOneRequest{
		AppId:     push.AppM[m.Tag],
		UUid:      m.Body.UUID,
		MsgId:     fmt.Sprint(time.Now().UnixNano()),
		Timestamp: fmt.Sprint(time.Now().UnixNano()),
		Content:   m.Body.String(),
	}); err != nil {
		return err
	}

	return nil
}

func pushApp(rpcClient pb.GatewayClient, m *push.PushMsg) error {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()

	log.GetLogger(common.LogicLog).Error("gpusher: logic pushApp: %s", *m)

	if _, err := rpcClient.PushApp(ctx, &pb.PushAppRequest{
		AppId:     push.AppM[m.Tag],
		MsgId:     fmt.Sprint(time.Now().UnixNano()),
		Timestamp: fmt.Sprint(time.Now().UnixNano()),
		Content:   m.Body.String(),
	}); err != nil {
		return err
	}

	return nil
}

func pushAll(rpcClient pb.GatewayClient, m *push.PushMsg) error {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()

	log.GetLogger(common.LogicLog).Error("gpusher: logic pushAll: %s", *m)

	if _, err := rpcClient.PushAll(ctx, &pb.PushAllRequest{
		MsgId:     fmt.Sprint(time.Now().UnixNano()),
		Timestamp: fmt.Sprint(time.Now().UnixNano()),
		Content:   m.Body.String(),
	}); err != nil {
		return err
	}

	return nil
}
