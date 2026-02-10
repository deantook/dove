package model

import "time"

type ProfileUnlockProgress struct {
	ID           uint64     `gorm:"primaryKey;autoIncrement;comment:解锁进度记录ID"`
	RelationID   uint64     `gorm:"not null;index;comment:所属用户关系ID"`
	ProfileID    uint64     `gorm:"not null;comment:被解锁的资料ID"`
	CurrentValue int        `gorm:"not null;default:0;comment:当前解锁进度值"`
	TargetValue  int        `gorm:"not null;comment:解锁所需目标值"`
	UnlockType   string     `gorm:"size:32;not null;comment:解锁类型 points/level/task/time/mutual"`
	Unlocked     bool       `gorm:"not null;default:false;comment:是否已解锁"`
	UnlockedAt   *time.Time `gorm:"comment:解锁完成时间"`
	CreatedAt    time.Time  `gorm:"autoCreateTime;comment:进度创建时间"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime;comment:进度更新时间"`

	Relation UserRelation `gorm:"foreignKey:RelationID;constraint:OnDelete:CASCADE"`
	Profile  UserProfile  `gorm:"foreignKey:ProfileID;constraint:OnDelete:CASCADE"`
}

func (ProfileUnlockProgress) TableName() string {
	return "profile_unlock_progress"
}
