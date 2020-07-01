/**
 *
 * @author liangjf
 * @create on 2020/6/1
 * @version 1.0
 */
package defind

const (
	//服务列表
	RedisKeyGatewayList = "gpusher_gateway_list" //网关列表		key:gpusher_gateway_list	value:gatewayAddr
	RedisKeyMessageList = "gpusher_message_list" //持久化服列表	key:gpusher_message_list	value:messageAddr
	RedisKeyLogicList   = "gpusher_logic_list"   //逻辑服列表	key:gpusher_logic_list		value:logicAddr

	//网关
	RedisKeyGatewayAllUUID = "gpusher_gateway_all_uuid" //set	key:gpusher_gateway_all_uuid	field:uuid	value:gatewayAddr
	RedisKeyGatewayAppUUID = "gpusher_gateway_app_uuid" //hash	key:gpusher_gateway_app_uuid	field:uuid	value:gatewayAddr

	RedisKeyGatewayAllUUIDNum = "gpusher_gateway_all_uuid_num" //string	key:gpusher_gateway_all_uuid_num	value:
	RedisKeyGatewayUUIDNum    = "gpusher_gateway_%s_uuid_num"  //string	key:gpusher_gateway_%s_uuid_num		value:

	//离线消息
	RedisKeyExpireMsg = "gpusher_expire_msg_uuid_"

	//用户token
	RedisKeyUUIDToken = "gpusher_uuid_token_" //hash 分片
)
