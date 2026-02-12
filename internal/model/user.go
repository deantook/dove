package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID         int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Username   string    `gorm:"column:username;type:varchar(255)" json:"username"`
	Phone      string    `gorm:"column:phone;type:varchar(20);uniqueIndex" json:"phone"`
	Nickname   string    `gorm:"column:nickname;type:varchar(100)" json:"nickname"`
	Avatar     string    `gorm:"column:avatar;type:varchar(500)" json:"avatar"`
	Status     int       `gorm:"column:status;type:tinyint;default:1" json:"status"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`
}

// TableName 指定表名
func (User) TableName() string {
	return "u_user"
}

// BeforeCreate 创建前钩子
func (u *User) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if u.CreateTime.IsZero() {
		u.CreateTime = now
	}
	if u.UpdateTime.IsZero() {
		u.UpdateTime = now
	}
	if u.Status == 0 {
		u.Status = 1 // 默认启用状态
	}
	return nil
}

// BeforeUpdate 更新前钩子
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.UpdateTime = time.Now()
	return nil
}
