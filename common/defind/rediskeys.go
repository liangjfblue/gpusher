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
	RedisKeyGatewayAllUUID = "gpusher_gateway_%s_all_uuid"    //hash	key:gpusher_gateway_%s_all_uuid	field:uuid	value:gatewayAddr
	RedisKeyGatewayAppUUID = "gpusher_gateway_%s_app_%s_uuid" //hash	key:gpusher_gateway_%s_app_%s_uuid	field:uuid	value:gatewayAddr

	RedisKeyGatewayUUIDNum = "gpusher_gateway_%s_uuid_num" //string	key:gpusher_gateway_%s_uuid_num		value:
)
