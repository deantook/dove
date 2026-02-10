package model

import "time"

type InteractionLevel struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement;comment:等级配置ID"`
	Level     int       `gorm:"not null;uniqueIndex;comment:关系等级数值"`
	LevelName string    `gorm:"size:32;not null;comment:等级名称"`
	MinScore  int       `gorm:"not null;comment:进入该等级的最小亲密度积分"`
	MaxScore  int       `gorm:"not null;comment:该等级的最大亲密度积分"`
	CreatedAt time.Time `gorm:"autoCreateTime;comment:等级配置创建时间"`
}

func (InteractionLevel) TableName() string {
	return "interaction_level"
}
