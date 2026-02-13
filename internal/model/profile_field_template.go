package model

import (
	"time"

	"gorm.io/gorm"
)

// ProfileFieldTemplate 系统资料字段模板模型
// 存储系统预设的单个字段类型定义，用户引用后会在 profile_fields 表中复制一条记录
type ProfileFieldTemplate struct {
	ID                 int            `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	FieldKey           string         `gorm:"column:field_key;type:varchar(100);uniqueIndex" json:"field_key"`
	FieldName          string         `gorm:"column:field_name;type:varchar(100)" json:"field_name"`
	FieldType          string         `gorm:"column:field_type;type:varchar(50)" json:"field_type"`
	IsRequired         bool           `gorm:"column:is_required;type:tinyint(1);default:0" json:"is_required"`
	IsSearchable       bool           `gorm:"column:is_searchable;type:tinyint(1);default:0" json:"is_searchable"`
	IsPublic           bool           `gorm:"column:is_public;type:tinyint(1);default:0" json:"is_public"`
	DefaultValue       string         `gorm:"column:default_value;type:text" json:"default_value"`
	Options            string         `gorm:"column:options;type:text" json:"options"`
	Validation         string         `gorm:"column:validation;type:text" json:"validation"`
	DisplayOrder       int            `gorm:"column:display_order;type:int;default:0" json:"display_order"`
	Icon               string         `gorm:"column:icon;type:varchar(500)" json:"icon"`
	Description        string         `gorm:"column:description;type:varchar(500)" json:"description"`
	DefaultUnlockRules string         `gorm:"column:default_unlock_rules;type:text" json:"default_unlock_rules"`
	Category           string         `gorm:"column:category;type:varchar(50)" json:"category"`
	IsActive           bool           `gorm:"column:is_active;type:tinyint(1);default:1" json:"is_active"`
	CreateTime         time.Time      `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime         time.Time      `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (ProfileFieldTemplate) TableName() string {
	return "profile_field_templates"
}

// CreateProfileFieldTemplateRequest 创建字段模板请求
type CreateProfileFieldTemplateRequest struct {
	FieldKey           string `json:"field_key" binding:"required,min=1,max=100" example:"education"`
	FieldName          string `json:"field_name" binding:"required,min=1,max=100" example:"学历"`
	FieldType          string `json:"field_type" binding:"required" example:"SELECT_SINGLE"`
	IsRequired         bool   `json:"is_required" example:"false"`
	IsSearchable       bool   `json:"is_searchable" example:"true"`
	IsPublic           bool   `json:"is_public" example:"false"`
	DefaultValue       string `json:"default_value" binding:"omitempty" example:""`
	Options            string `json:"options" binding:"omitempty" example:"{\"options\":[{\"key\":\"bachelor\",\"label\":\"本科\"}]}"`
	Validation         string `json:"validation" binding:"omitempty" example:"{}"`
	DisplayOrder       int    `json:"display_order" example:"10"`
	Icon               string `json:"icon" binding:"omitempty,max=500" example:""`
	Description        string `json:"description" binding:"omitempty,max=500" example:"最高学历"`
	DefaultUnlockRules string `json:"default_unlock_rules" binding:"omitempty" example:"{\"unlock_type\":\"CHAT\",\"conditions\":{\"message_count\":50}}"`
	Category           string `json:"category" binding:"omitempty,max=50" example:"教育背景"`
}

// UpdateProfileFieldTemplateRequest 更新字段模板请求
type UpdateProfileFieldTemplateRequest struct {
	FieldName          string `json:"field_name" binding:"omitempty,min=1,max=100" example:"学历"`
	FieldType          string `json:"field_type" binding:"omitempty" example:"SELECT_SINGLE"`
	IsRequired         *bool  `json:"is_required" example:"false"`
	IsSearchable       *bool  `json:"is_searchable" example:"true"`
	IsPublic           *bool  `json:"is_public" example:"false"`
	DefaultValue       string `json:"default_value" binding:"omitempty" example:""`
	Options            string `json:"options" binding:"omitempty" example:"{}"`
	Validation         string `json:"validation" binding:"omitempty" example:"{}"`
	DisplayOrder       *int   `json:"display_order" example:"10"`
	Icon               string `json:"icon" binding:"omitempty,max=500" example:""`
	Description        string `json:"description" binding:"omitempty,max=500" example:"最高学历"`
	DefaultUnlockRules string `json:"default_unlock_rules" binding:"omitempty" example:"{}"`
	Category           string `json:"category" binding:"omitempty,max=50" example:"教育背景"`
	IsActive           *bool  `json:"is_active" example:"true"`
}

// ProfileFieldTemplateResponse 字段模板响应
type ProfileFieldTemplateResponse struct {
	ID                 int       `json:"id"`
	FieldKey           string    `json:"field_key"`
	FieldName          string    `json:"field_name"`
	FieldType          string    `json:"field_type"`
	IsRequired         bool      `json:"is_required"`
	IsSearchable       bool      `json:"is_searchable"`
	IsPublic           bool      `json:"is_public"`
	DefaultValue       string    `json:"default_value"`
	Options            string    `json:"options"`
	Validation         string    `json:"validation"`
	DisplayOrder       int       `json:"display_order"`
	Icon               string    `json:"icon"`
	Description        string    `json:"description"`
	DefaultUnlockRules string    `json:"default_unlock_rules"`
	Category           string    `json:"category"`
	IsActive           bool      `json:"is_active"`
	CreateTime         time.Time `json:"create_time"`
	UpdateTime         time.Time `json:"update_time"`
}

// ToResponse 转换为响应格式
func (t *ProfileFieldTemplate) ToResponse() *ProfileFieldTemplateResponse {
	return &ProfileFieldTemplateResponse{
		ID:                 t.ID,
		FieldKey:           t.FieldKey,
		FieldName:          t.FieldName,
		FieldType:          t.FieldType,
		IsRequired:         t.IsRequired,
		IsSearchable:       t.IsSearchable,
		IsPublic:           t.IsPublic,
		DefaultValue:       t.DefaultValue,
		Options:            t.Options,
		Validation:         t.Validation,
		DisplayOrder:       t.DisplayOrder,
		Icon:               t.Icon,
		Description:        t.Description,
		DefaultUnlockRules: t.DefaultUnlockRules,
		Category:           t.Category,
		IsActive:           t.IsActive,
		CreateTime:         t.CreateTime,
		UpdateTime:         t.UpdateTime,
	}
}

// ApplyToUser 将字段模板应用到用户，返回 ProfileField 结构
// 用户引用后会在 profile_fields 表中创建一条记录
func (t *ProfileFieldTemplate) ApplyToUser(userID int) *ProfileField {
	return &ProfileField{
		UserID:       userID,
		FieldKey:     t.FieldKey,
		FieldName:    t.FieldName,
		FieldType:    t.FieldType,
		IsSystem:     true, // 来自系统模板
		IsRequired:   t.IsRequired,
		IsSearchable: t.IsSearchable,
		IsPublic:     t.IsPublic,
		DefaultValue: t.DefaultValue,
		Options:      t.Options,
		Validation:   t.Validation,
		DisplayOrder: t.DisplayOrder,
		Icon:         t.Icon,
		Description:  t.Description,
	}
}

// ProfileField 资料字段模型（用于应用模板时创建）
type ProfileField struct {
	ID           int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID       int       `gorm:"column:user_id;type:int;default:0" json:"user_id"`
	FieldKey     string    `gorm:"column:field_key;type:varchar(100)" json:"field_key"`
	FieldName    string    `gorm:"column:field_name;type:varchar(100)" json:"field_name"`
	FieldType    string    `gorm:"column:field_type;type:varchar(50)" json:"field_type"`
	IsSystem     bool      `gorm:"column:is_system;type:tinyint(1);default:0" json:"is_system"`
	IsRequired   bool      `gorm:"column:is_required;type:tinyint(1);default:0" json:"is_required"`
	IsSearchable bool      `gorm:"column:is_searchable;type:tinyint(1);default:0" json:"is_searchable"`
	IsPublic     bool      `gorm:"column:is_public;type:tinyint(1);default:0" json:"is_public"`
	DefaultValue string    `gorm:"column:default_value;type:text" json:"default_value"`
	Options      string    `gorm:"column:options;type:text" json:"options"`
	Validation   string    `gorm:"column:validation;type:text" json:"validation"`
	DisplayOrder int       `gorm:"column:display_order;type:int;default:0" json:"display_order"`
	Icon         string    `gorm:"column:icon;type:varchar(500)" json:"icon"`
	Description  string    `gorm:"column:description;type:varchar(500)" json:"description"`
	CreateTime   time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
}

// TableName 指定表名
func (ProfileField) TableName() string {
	return "profile_fields"
}
