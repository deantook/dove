package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
// @Description 用户信息
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey" example:"1"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null;size:50" example:"john_doe"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null;size:100" example:"john@example.com"`
	Password  string         `json:"password,omitempty" gorm:"not null;size:255" swaggerignore:"true"`
	Nickname  string         `json:"nickname" gorm:"size:50" example:"John Doe"`
	Avatar    string         `json:"avatar" gorm:"size:255" example:"https://example.com/avatar.jpg"`
	Status    int            `json:"status" gorm:"default:1" example:"1"` // 1: 正常, 0: 禁用
	CreatedAt time.Time      `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time      `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index" swaggerignore:"true"`
}
