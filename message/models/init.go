/**
 *
 * @author liangjf
 * @create on 2020/6/4
 * @version 1.0
 */
package models

type IModels interface {
	//初始化
	Init() error
	//保存网关uuid映射
	SaveGatewayUUID(string, string) error
	//保存App和uuid映射
	SaveAppUUID(string, string) error
	//保存离线消息
	SaveExpireMsg(string, string, string, int64) error
	//删除网关uuid映射
	DeleteGatewayUUID(string) error
	//删除AppTag和uuid映射
	DeleteAppUUID(string, string) error
	//删除离线消息
	DeleteExpireMsg(string, string) error
	//获取网关uuid映射
	GetGatewayUUID(string) (string, error)
	//获取App和uuid映射
	GetAppUUID(string) ([]string, error)
}
