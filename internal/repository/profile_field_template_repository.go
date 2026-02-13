package repository

import (
	"github.com/deantook/dove/internal/model"
	"gorm.io/gorm"
)

// ProfileFieldTemplateRepository 系统资料字段模板仓储接口
type ProfileFieldTemplateRepository interface {
	Create(template *model.ProfileFieldTemplate) error
	GetByID(id int) (*model.ProfileFieldTemplate, error)
	GetByFieldKey(fieldKey string) (*model.ProfileFieldTemplate, error)
	Update(template *model.ProfileFieldTemplate) error
	Delete(id int) error
	List(category string, fieldType string, isActive *bool, offset, limit int) ([]*model.ProfileFieldTemplate, int64, error)
	GetByCategory(category string) ([]*model.ProfileFieldTemplate, error)
}

// profileFieldTemplateRepository 系统资料字段模板仓储实现
type profileFieldTemplateRepository struct {
	db *gorm.DB
}

// NewProfileFieldTemplateRepository 创建系统资料字段模板仓储实例
func NewProfileFieldTemplateRepository(db *gorm.DB) ProfileFieldTemplateRepository {
	return &profileFieldTemplateRepository{db: db}
}

// Create 创建字段模板
func (r *profileFieldTemplateRepository) Create(template *model.ProfileFieldTemplate) error {
	return r.db.Create(template).Error
}

// GetByID 根据 ID 获取字段模板
func (r *profileFieldTemplateRepository) GetByID(id int) (*model.ProfileFieldTemplate, error) {
	var template model.ProfileFieldTemplate
	err := r.db.Where("id = ?", id).First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// GetByFieldKey 根据字段标识获取字段模板
func (r *profileFieldTemplateRepository) GetByFieldKey(fieldKey string) (*model.ProfileFieldTemplate, error) {
	var template model.ProfileFieldTemplate
	err := r.db.Where("field_key = ?", fieldKey).First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// Update 更新字段模板
func (r *profileFieldTemplateRepository) Update(template *model.ProfileFieldTemplate) error {
	return r.db.Save(template).Error
}

// Delete 删除字段模板（软删除）
func (r *profileFieldTemplateRepository) Delete(id int) error {
	return r.db.Delete(&model.ProfileFieldTemplate{}, id).Error
}

// List 获取字段模板列表（分页）
func (r *profileFieldTemplateRepository) List(category string, fieldType string, isActive *bool, offset, limit int) ([]*model.ProfileFieldTemplate, int64, error) {
	var templates []*model.ProfileFieldTemplate
	var total int64

	query := r.db.Model(&model.ProfileFieldTemplate{})

	// 按分类筛选
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// 按字段类型筛选
	if fieldType != "" {
		query = query.Where("field_type = ?", fieldType)
	}

	// 按启用状态筛选
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取列表，按分类和显示顺序排序
	if err := query.Order("category ASC, display_order ASC, create_time DESC").
		Offset(offset).Limit(limit).Find(&templates).Error; err != nil {
		return nil, 0, err
	}

	return templates, total, nil
}

// GetByCategory 根据分类获取字段模板列表
func (r *profileFieldTemplateRepository) GetByCategory(category string) ([]*model.ProfileFieldTemplate, error) {
	var templates []*model.ProfileFieldTemplate
	err := r.db.Where("category = ? AND is_active = ?", category, true).
		Order("display_order ASC").
		Find(&templates).Error
	if err != nil {
		return nil, err
	}
	return templates, nil
}
