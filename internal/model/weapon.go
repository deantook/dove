package model

import (
	"time"

	"gorm.io/gorm"
)

// Weapon 武器模型
// @Description 武器信息
type Weapon struct {
	ID        uint           `json:"id" gorm:"primaryKey" example:"1"`
	Name      string         `json:"name" gorm:"uniqueIndex;not null;size:50" example:"AK47"`
	Level     int            `json:"level" gorm:"default:1" example:"1"`
	Content   string         `json:"content" gorm:"text" example:"AK47 is a popular weapon in the world"`
	Type      int            `json:"type" gorm:"default:1" example:"1"`
	Story     string         `json:"story" gorm:"text" example:"AK47 is a popular weapon in the world"`
	CreatedAt time.Time      `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time      `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index" swaggerignore:"true"`
}
