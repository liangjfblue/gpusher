/**
 *
 * @author liangjf
 * @create on 2020/6/1
 * @version 1.0
 */
package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type mysql struct {
	Addr        string
	User        string
	Password    string
	Db          string
	MaxIdleConn int
	MaxOpenConn int
}

func NewMysql() IDB {
	return &mysql{}
}

func (d *mysql) Init() {
	var (
		err  error
		addr string
	)

	addr = d.Addr
	if addr == "" {
		addr = "127.0.0.1"
	}
	str := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", d.User, d.Password, addr, d.Db)
	DB, err := gorm.Open("mysql", str)
	if err != nil {
		panic(err)
	}

	DB.LogMode(true)
	DB.SingularTable(true)
	DB.DB().SetMaxIdleConns(d.MaxIdleConn)
	DB.DB().SetMaxOpenConns(d.MaxOpenConn)
}
