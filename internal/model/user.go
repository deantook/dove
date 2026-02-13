package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID         int            `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Username   string         `gorm:"column:username;type:varchar(255)" json:"username"`
	Phone      string         `gorm:"column:phone;type:varchar(20);uniqueIndex" json:"phone"`
	Avatar     string         `gorm:"column:avatar;type:varchar(500)" json:"avatar"`
	Status     int            `gorm:"column:status;type:tinyint;default:1" json:"status"`
	CreateTime time.Time      `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime time.Time      `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (User) TableName() string {
	return "u_user"
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50" example:"john_doe"`
	Phone    string `json:"phone" binding:"required,phone" example:"13800138000"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Username string `json:"username" binding:"omitempty,min=3,max=50" example:"john_doe"`
	Phone    string `json:"phone" binding:"omitempty,phone" example:"13800138000"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID         int       `json:"id"`
	Username   string    `json:"username"`
	Phone      string    `json:"phone"`
	Avatar     string    `json:"avatar"`
	Status     int       `json:"status"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

// ToResponse 转换为响应格式
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:         u.ID,
		Username:   u.Username,
		Phone:      u.Phone,
		Avatar:     u.Avatar,
		Status:     u.Status,
		CreateTime: u.CreateTime,
		UpdateTime: u.UpdateTime,
	}
}

// SendCodeRequest 发送验证码请求
type SendCodeRequest struct {
	Phone string `json:"phone" binding:"required,phone" example:"13800138000"`
}

// SendCodeResponse 发送验证码响应
type SendCodeResponse struct {
	Code string `json:"code" example:"123456"` // 验证码（开发阶段返回，生产环境不返回）
}

// LoginRequest 登录/注册请求
type LoginRequest struct {
	Phone string `json:"phone" binding:"required,phone" example:"13800138000"`
	Code  string `json:"code" binding:"required,len=6" example:"123456"`
}

// LoginResponse 登录/注册响应
type LoginResponse struct {
	User  *UserResponse `json:"user"`
	Token string        `json:"token"`
}
