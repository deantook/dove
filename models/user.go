package models

import (
	"time"
)

// User 用户模型
type User struct {
	ID            int64      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Username      string     `gorm:"type:varchar(50);uniqueIndex;not null;column:username" json:"username"`
	Password      string     `gorm:"type:varchar(255);not null;column:password" json:"-"`
	Email         *string    `gorm:"type:varchar(100);uniqueIndex;column:email" json:"email"`
	Phone         *string    `gorm:"type:varchar(20);index;column:phone" json:"phone"`
	Nickname      *string    `gorm:"type:varchar(50);column:nickname" json:"nickname"`
	Avatar        *string    `gorm:"type:varchar(500);column:avatar" json:"avatar"`
	Status        int8       `gorm:"type:tinyint;default:1;column:status" json:"status"`
	Gender        *int8      `gorm:"type:tinyint;column:gender" json:"gender"`
	Birthday      *time.Time `gorm:"type:date;column:birthday" json:"birthday"`
	LastLoginTime *time.Time `gorm:"type:datetime;column:last_login_time" json:"last_login_time"`
	LastLoginIP   *string    `gorm:"type:varchar(50);column:last_login_ip" json:"last_login_ip"`
	CreateTime    time.Time  `gorm:"type:datetime;default:CURRENT_TIMESTAMP;column:create_time" json:"create_time"`
	UpdateTime    time.Time  `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:update_time" json:"update_time"`
	Deleted       int8       `gorm:"type:tinyint;default:0;column:deleted" json:"-"`
}

// TableName 指定表名
func (User) TableName() string {
	return "u_user"
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone" binding:"omitempty"`
	Nickname string `json:"nickname" binding:"omitempty,max=50"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}

// UserInfo 用户信息（不包含敏感信息）
type UserInfo struct {
	ID       int64      `json:"id"`
	Username string     `json:"username"`
	Email    *string    `json:"email"`
	Phone    *string    `json:"phone"`
	Nickname *string    `json:"nickname"`
	Avatar   *string    `json:"avatar"`
	Status   int8       `json:"status"`
	Gender   *int8      `json:"gender"`
	Birthday *time.Time `json:"birthday"`
}

// ToUserInfo 转换为用户信息
func (u *User) ToUserInfo() UserInfo {
	return UserInfo{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Phone:    u.Phone,
		Nickname: u.Nickname,
		Avatar:   u.Avatar,
		Status:   u.Status,
		Gender:   u.Gender,
		Birthday: u.Birthday,
	}
}
