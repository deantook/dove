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

type troveRepository struct{}

func NewTroveRepository() domain.TroveRepository {
	return &troveRepository{}
}

func (r *troveRepository) Create(ctx context.Context, trove *model.Trove) error {
	return database.DB.WithContext(ctx).Create(trove).Error
}

func (r *troveRepository) GetByID(ctx context.Context, id uint) (*model.Trove, error) {
	var trove model.Trove
	err := database.DB.WithContext(ctx).First(&trove, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("trove not found")
		}
		return nil, err
	}
	return &trove, nil
}

func (r *troveRepository) GetAll(ctx context.Context) ([]model.Trove, error) {
	var troves []model.Trove
	err := database.DB.WithContext(ctx).Find(&troves).Error
	return troves, err
}

func (r *troveRepository) GetAllWithPagination(ctx context.Context, page *pagination.PageRequest) ([]model.Trove, int64, error) {
	var troves []model.Trove
	var total int64

	// 构建查询
	query := database.DB.WithContext(ctx)

	// 添加搜索条件
	if page.HasSearch() {
		// 验证搜索字段
		allowedSearchFields := []string{"Title", "Description"}
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

			searchConditions = append(searchConditions, "Title LIKE ?")
			searchArgs = append(searchArgs, "%"+keyword+"%")

			searchConditions = append(searchConditions, "Description LIKE ?")
			searchArgs = append(searchArgs, "%"+keyword+"%")

			query = query.Where(strings.Join(searchConditions, " OR "), searchArgs...)
		}
	}

	// 获取总记录数
	if err := query.Model(&model.Trove{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 添加排序
	if page.HasSort() {
		// 验证排序字段
		allowedFields := []string{"ID"}
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
	err := query.Offset(offset).Limit(limit).Find(&troves).Error

	return troves, total, err
}

func (r *troveRepository) Update(ctx context.Context, trove *model.Trove) error {
	// 检查记录是否存在
	var existingTrove model.Trove
	if err := database.DB.WithContext(ctx).First(&existingTrove, trove.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("trove not found")
		}
		return err
	}

	return database.DB.WithContext(ctx).Save(trove).Error
}

func (r *troveRepository) Delete(ctx context.Context, id uint) error {
	// 检查记录是否存在
	var trove model.Trove
	if err := database.DB.WithContext(ctx).First(&trove, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("trove not found")
		}
		return err
	}

	return database.DB.WithContext(ctx).Delete(&trove).Error
}
