package model

import "time"

type UserInteraction struct {
	ID              uint64    `gorm:"primaryKey;autoIncrement;comment:互动记录ID"`
	RelationID      uint64    `gorm:"not null;index;comment:所属用户关系ID"`
	FromUserID      uint64    `gorm:"not null;index;comment:发起互动的用户ID"`
	InteractionType string    `gorm:"size:32;not null;comment:互动类型 text/voice/video/like/task"`
	Value           int       `gorm:"not null;default:0;comment:互动原始值，如条数或秒数"`
	Points          int       `gorm:"not null;default:0;comment:本次互动获得的亲密度积分"`
	IsAIAssisted    bool      `gorm:"not null;default:false;comment:是否AI辅助生成"`
	CreatedAt       time.Time `gorm:"autoCreateTime;comment:互动发生时间"`

	Relation UserRelation `gorm:"foreignKey:RelationID;constraint:OnDelete:CASCADE"`
}

func (UserInteraction) TableName() string {
	return "user_interaction"
}
