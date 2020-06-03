/**
 *
 * @author liangjf
 * @create on 2020/6/3
 * @version 1.0
 */
package service

import (
	"github.com/liangjfblue/gpusher/common/push"
)

func router(msg *push.PushMsg) (string, error) {
	//TODO 根据type appId, uuid查看gateway地址
	if msg.Body.UUID == "liangjf" {
		return "127.0.0.1:7771", nil
	}

	return "", nil
}
