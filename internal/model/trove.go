package model

import (
	"time"

	"gorm.io/gorm"
)

// Trove 收藏模型
// @Description 收藏信息
type Trove struct {
	ID          uint           `json:"id" gorm:"primaryKey" example:"1"`
	Title       string         `json:"title" gorm:"size:255"`
	Description string         `json:"description" gorm:"text"`
	CreatedAt   time.Time      `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   time.Time      `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index" swaggerignore:"true"`
}
