/**
 *
 * @author liangjf
 * @create on 2020/9/10
 * @version 1.0
 */
package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/liangjfblue/gpusher/web/config"
)

var (
	_db *gorm.DB
)

// InitMysql 初始化mysql
func InitMysqlPool(mysqlConf *config.Mysql) {
	var err error
	str := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", mysqlConf.User, mysqlConf.Password, mysqlConf.Addr, mysqlConf.Db)
	_db, err = gorm.Open("mysql", str)
	if err != nil {
		panic(err)
	}

	_db.LogMode(false)
	_db.DB().SetMaxIdleConns(mysqlConf.MaxIdleConns)
	_db.DB().SetMaxOpenConns(mysqlConf.MaxOpenConns)
	return
}

func GetDB() *gorm.DB {
	return _db
}
