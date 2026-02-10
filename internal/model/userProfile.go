package model

import (
	"time"

	"gorm.io/datatypes"
)

type UserProfile struct {
	ID           uint64         `gorm:"primaryKey;autoIncrement;comment:资料项ID"`
	UserID       uint64         `gorm:"index;not null;comment:资料所属用户ID"`
	ProfileType  string         `gorm:"size:64;not null;comment:资料类型 basic/interest/personality/private"`
	Content      datatypes.JSON `gorm:"type:json;not null;comment:资料内容JSON，结构由ProfileType决定"`
	PrivacyLevel int8           `gorm:"not null;default:1;comment:隐私级别 1=公开 2=可解锁 3=私密 4=需双方同意"`
	CreatedAt    time.Time      `gorm:"autoCreateTime;comment:资料创建时间"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime;comment:资料更新时间"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func (UserProfile) TableName() string {
	return "user_profile"
}
