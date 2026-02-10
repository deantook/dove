package domain

import "time"

type User struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement;comment:用户唯一ID"`
	Email        *string   `gorm:"size:128;uniqueIndex;comment:用户邮箱，用于登录，可为空"`
	Phone        *string   `gorm:"size:32;uniqueIndex;comment:用户手机号，用于登录，可为空"`
	Nickname     string    `gorm:"size:32;not null;comment:用户昵称，用于展示"`
	PasswordHash string    `gorm:"size:255;not null;comment:密码哈希值"`
	Status       int8      `gorm:"not null;default:1;comment:用户状态 1=正常 0=禁用"`
	CreatedAt    time.Time `gorm:"autoCreateTime;comment:用户创建时间"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime;comment:用户信息更新时间"`
}

func (User) TableName() string {
	return "user"
}
