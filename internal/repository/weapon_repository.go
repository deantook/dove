package repository

import (
	"context"
	"dove/internal/domain"
	"dove/internal/model"
	"dove/pkg/database"
	"dove/pkg/pagination"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type weaponRepository struct{}

func NewWeaponRepository() domain.WeaponRepository {
	return &weaponRepository{}
}

func (r *weaponRepository) Create(ctx context.Context, weapon *model.Weapon) error {
	return database.DB.WithContext(ctx).Create(weapon).Error
}

func (r *weaponRepository) GetByID(ctx context.Context, id uint) (*model.Weapon, error) {
	var weapon model.Weapon
	err := database.DB.WithContext(ctx).First(&weapon, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("weapon not found")
		}
		return nil, err
	}
	return &weapon, nil
}

func (r *weaponRepository) GetAll(ctx context.Context) ([]model.Weapon, error) {
	var weapons []model.Weapon
	err := database.DB.WithContext(ctx).Find(&weapons).Error
	return weapons, err
}

func (r *weaponRepository) GetAllWithPagination(ctx context.Context, page *pagination.PageRequest) ([]model.Weapon, int64, error) {
	var weapons []model.Weapon
	var total int64

	// 构建查询
	query := database.DB.WithContext(ctx)

	// 添加搜索条件
	if page.HasSearch() {
		// 验证搜索字段
		allowedSearchFields := []string{"Name"}
		if !page.ValidateSearchField(allowedSearchFields) {
			return nil, 0, fmt.Errorf("invalid search field: %s", page.GetSearchBy())
		}

		// 如果指定了搜索字段，使用该字段进行搜索
		if page.GetSearchBy() != "" {
			searchField := page.GetSearchBy()
			keyword := page.GetKeyword()
			query = query.Where(fmt.Sprintf("%s LIKE ?", searchField), "%"+keyword+"%")
		} else {
			// 如果没有指定搜索字段，在所有可搜索字段中搜索
			keyword := page.GetKeyword()
			searchConditions := []string{}
			searchArgs := []interface{}{}

			searchConditions = append(searchConditions, "Name LIKE ?")
			searchArgs = append(searchArgs, "%"+keyword+"%")

			query = query.Where(strings.Join(searchConditions, " OR "), searchArgs...)
		}
	}

	// 获取总记录数
	if err := query.Model(&model.Weapon{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 添加排序
	if page.HasSort() {
		// 验证排序字段
		allowedFields := []string{"ID", "Name"}
		if !page.ValidateSortField(allowedFields) {
			return nil, 0, fmt.Errorf("invalid sort field: %s", page.GetSortBy())
		}

		// 构建排序语句
		sortClause := page.GetSortBy()
		if page.GetSortOrder() == "desc" {
			sortClause += " DESC"
		} else {
			sortClause += " ASC"
		}
		query = query.Order(sortClause)
	} else {
		// 默认按创建时间倒序
		query = query.Order("created_at DESC")
	}

	// 获取分页数据
	offset := page.GetOffset()
	limit := page.GetLimit()
	err := query.Offset(offset).Limit(limit).Find(&weapons).Error

	return weapons, total, err
}

func (r *weaponRepository) Update(ctx context.Context, weapon *model.Weapon) error {
	// 检查记录是否存在
	var existingWeapon model.Weapon
	if err := database.DB.WithContext(ctx).First(&existingWeapon, weapon.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("weapon not found")
		}
		return err
	}

	return database.DB.WithContext(ctx).Save(weapon).Error
}

func (r *weaponRepository) Delete(ctx context.Context, id uint) error {
	// 检查记录是否存在
	var weapon model.Weapon
	if err := database.DB.WithContext(ctx).First(&weapon, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("weapon not found")
		}
		return err
	}

	return database.DB.WithContext(ctx).Delete(&weapon).Error
}
