/**
 *
 * @author liangjf
 * @create on 2020/9/10
 * @version 1.0
 */
package models

import (
	"github.com/jinzhu/gorm"
)

// TBUser 用户表
type TBUser struct {
	gorm.Model
	Username    string `gorm:"column:username;not null;index:index_username" json:"username" description:"username"`
	Password    string `gorm:"column:password;not null;" json:"password" description:"password"`
	UUID        string `gorm:"column:uuid;not null;index:index_uuid" json:"uuid" description:"uuid"`
	UserID      string `gorm:"column:user_id;not null;" json:"user_id" description:"用户id"`
	Phone       string `gorm:"column:phone;not null;" json:"phone" description:"手机号"`
	IsAvailable int8   `gorm:"column:is_available; default:1; not null;" json:"isavailable"  description:"是否可用 1可用 0不可用"`
}

// TableName 表名
func (t *TBUser) TableName() string {
	return "tb_user"
}

// AddTBUser 插入记录
func (t *TBUser) AddTBUser() error {
	return GetDB().Create(t).Error
}

// GetTBUser 查找一个
func GetTBUser(query map[string]interface{}) (*TBUser, error) {
	var user TBUser
	err := GetDB().Where(query).First(&user).Error
	return &user, err
}

// DeleteTBUser 删除记录
func DeleteTBUser(uuid string) error {
	user := TBUser{
		UUID: uuid,
	}
	return GetDB().Delete(&user).Error
}

// GetAllTBUsers 获取所有记录
func GetAllTBUsers(query map[string]interface{}, offset int32, limit int32) ([]TBUser, error) {
	var users []TBUser
	err := GetDB().Where(query).Offset(offset).Limit(limit).Find(&users).Error
	return users, err
}

// UpdateTBUser 更新记录
func (t *TBUser) UpdateTBUser() error {
	return GetDB().Save(t).Error
}
