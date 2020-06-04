/**
 *
 * @author liangjf
 * @create on 2020/6/1
 * @version 1.0
 */
package db

var (
	_dbMap map[string]IDB
)

type IDB interface {
	Init()
}

func init() {
	RegisterDB("mysql", NewMysql())
}

func RegisterDB(name string, db IDB) {
	if _dbMap == nil {
		_dbMap = make(map[string]IDB)
	}
	_dbMap[name] = db
}

func GetDB(name string) IDB {
	return _dbMap[name]
}
