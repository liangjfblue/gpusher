/**
 *
 * @author liangjf
 * @create on 2020/6/1
 * @version 1.0
 */
package defind

const (
	AppAll = iota
	AppTest
	AppPushSystem
)

var (
	_appM map[int]string
)

func init() {
	RegisterApp(AppAll, "all-test")
	RegisterApp(AppTest, "app-test")
}

func RegisterApp(appId int, appName string) {
	if _appM == nil {
		_appM = make(map[int]string)
	}
	_appM[appId] = appName
}

func GetApp(appId int) string {
	if v, ok := _appM[appId]; ok {
		return v
	}
	return ""
}
