package model

import "time"

type UserRelation struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement;comment:用户关系ID"`
	User1ID   uint64    `gorm:"not null;index;comment:用户1的ID"`
	User2ID   uint64    `gorm:"not null;index;comment:用户2的ID"`
	Intimacy  int       `gorm:"not null;default:0;comment:亲密度积分"`
	Level     int       `gorm:"not null;default:1;comment:关系等级"`
	Status    int8      `gorm:"not null;default:1;comment:关系状态 1=正常 0=已解除"`
	CreatedAt time.Time `gorm:"autoCreateTime;comment:关系建立时间"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;comment:关系更新时间"`

	User1 User `gorm:"foreignKey:User1ID;constraint:OnDelete:CASCADE"`
	User2 User `gorm:"foreignKey:User2ID;constraint:OnDelete:CASCADE"`
}

func (UserRelation) TableName() string {
	return "user_relation"
}
