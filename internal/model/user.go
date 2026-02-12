package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID         int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Username   string    `gorm:"column:username;type:varchar(255)" json:"username"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
}

// TableName 指定表名
func (User) TableName() string {
	return "u_user"
}

// BeforeCreate 创建前钩子
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.CreateTime.IsZero() {
		u.CreateTime = time.Now()
	}
	return nil
}
