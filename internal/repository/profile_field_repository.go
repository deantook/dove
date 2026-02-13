package repository

import (
	"github.com/deantook/dove/internal/model"
	"gorm.io/gorm"
)

// ProfileFieldRepository 资料字段仓储接口
type ProfileFieldRepository interface {
	Create(field *model.ProfileField) error
	GetByID(id int) (*model.ProfileField, error)
	GetByUserIDAndFieldKey(userID int, fieldKey string) (*model.ProfileField, error)
	GetByUserID(userID int) ([]*model.ProfileField, error)
	Update(field *model.ProfileField) error
	Delete(id int) error
	List(userID int, offset, limit int) ([]*model.ProfileField, int64, error)
}

// profileFieldRepository 资料字段仓储实现
type profileFieldRepository struct {
	db *gorm.DB
}

// NewProfileFieldRepository 创建资料字段仓储实例
func NewProfileFieldRepository(db *gorm.DB) ProfileFieldRepository {
	return &profileFieldRepository{db: db}
}

// Create 创建字段
func (r *profileFieldRepository) Create(field *model.ProfileField) error {
	return r.db.Create(field).Error
}

// GetByID 根据 ID 获取字段
func (r *profileFieldRepository) GetByID(id int) (*model.ProfileField, error) {
	var field model.ProfileField
	err := r.db.Where("id = ?", id).First(&field).Error
	if err != nil {
		return nil, err
	}
	return &field, nil
}

// GetByUserIDAndFieldKey 根据用户ID和字段标识获取字段
func (r *profileFieldRepository) GetByUserIDAndFieldKey(userID int, fieldKey string) (*model.ProfileField, error) {
	var field model.ProfileField
	err := r.db.Where("user_id = ? AND field_key = ?", userID, fieldKey).First(&field).Error
	if err != nil {
		return nil, err
	}
	return &field, nil
}

// GetByUserID 根据用户ID获取所有字段
func (r *profileFieldRepository) GetByUserID(userID int) ([]*model.ProfileField, error) {
	var fields []*model.ProfileField
	err := r.db.Where("user_id = ?", userID).
		Order("display_order ASC").
		Find(&fields).Error
	if err != nil {
		return nil, err
	}
	return fields, nil
}

// Update 更新字段
func (r *profileFieldRepository) Update(field *model.ProfileField) error {
	return r.db.Save(field).Error
}

// Delete 删除字段
func (r *profileFieldRepository) Delete(id int) error {
	return r.db.Delete(&model.ProfileField{}, id).Error
}

// List 获取字段列表（分页）
func (r *profileFieldRepository) List(userID int, offset, limit int) ([]*model.ProfileField, int64, error) {
	var fields []*model.ProfileField
	var total int64

	query := r.db.Model(&model.ProfileField{}).Where("user_id = ?", userID)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取列表
	if err := query.Order("display_order ASC").
		Offset(offset).Limit(limit).Find(&fields).Error; err != nil {
		return nil, 0, err
	}

	return fields, total, nil
}
