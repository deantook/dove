package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/deantook/dove/internal/model"
	"github.com/deantook/dove/internal/repository"
	"gorm.io/gorm"
)

// ProfileFieldTemplateService 系统资料字段模板服务接口
type ProfileFieldTemplateService interface {
	CreateTemplate(ctx context.Context, req *model.CreateProfileFieldTemplateRequest) (*model.ProfileFieldTemplateResponse, error)
	GetTemplateByID(ctx context.Context, id int) (*model.ProfileFieldTemplateResponse, error)
	GetTemplateByFieldKey(ctx context.Context, fieldKey string) (*model.ProfileFieldTemplateResponse, error)
	UpdateTemplate(ctx context.Context, id int, req *model.UpdateProfileFieldTemplateRequest) (*model.ProfileFieldTemplateResponse, error)
	DeleteTemplate(ctx context.Context, id int) error
	ListTemplates(ctx context.Context, category string, fieldType string, isActive *bool, page, pageSize int) ([]*model.ProfileFieldTemplateResponse, int64, error)
	GetTemplatesByCategory(ctx context.Context, category string) ([]*model.ProfileFieldTemplateResponse, error)
	ApplyTemplateToUser(ctx context.Context, templateID, userID int) (*ApplyTemplateResult, error)
	ApplyTemplatesToUser(ctx context.Context, templateIDs []int, userID int) (*ApplyTemplatesResult, error)
}

// ApplyTemplateResult 应用模板结果
type ApplyTemplateResult struct {
	FieldID   int    `json:"field_id"`
	FieldKey  string `json:"field_key"`
	FieldName string `json:"field_name"`
	Message   string `json:"message"`
}

// ApplyTemplatesResult 批量应用模板结果
type ApplyTemplatesResult struct {
	AppliedFields []ApplyTemplateResult `json:"applied_fields"`
	TotalCount    int                   `json:"total_count"`
	SuccessCount  int                   `json:"success_count"`
	FailedCount   int                   `json:"failed_count"`
	Message       string                `json:"message"`
}

// profileFieldTemplateService 系统资料字段模板服务实现
type profileFieldTemplateService struct {
	templateRepo repository.ProfileFieldTemplateRepository
	fieldRepo    repository.ProfileFieldRepository // 需要创建 ProfileFieldRepository
}

// NewProfileFieldTemplateService 创建系统资料字段模板服务实例
func NewProfileFieldTemplateService(
	templateRepo repository.ProfileFieldTemplateRepository,
	fieldRepo repository.ProfileFieldRepository,
) ProfileFieldTemplateService {
	return &profileFieldTemplateService{
		templateRepo: templateRepo,
		fieldRepo:    fieldRepo,
	}
}

// CreateTemplate 创建字段模板
func (s *profileFieldTemplateService) CreateTemplate(ctx context.Context, req *model.CreateProfileFieldTemplateRequest) (*model.ProfileFieldTemplateResponse, error) {
	// 检查字段标识是否已存在
	if _, err := s.templateRepo.GetByFieldKey(req.FieldKey); err == nil {
		return nil, errors.New("字段标识已存在")
	} else if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("查询字段模板失败: %w", err)
	}

	// 验证选项配置 JSON（如果提供）
	if req.Options != "" {
		var options map[string]interface{}
		if err := json.Unmarshal([]byte(req.Options), &options); err != nil {
			return nil, fmt.Errorf("选项配置格式错误: %w", err)
		}
	}

	// 验证验证规则 JSON（如果提供）
	if req.Validation != "" {
		var validation map[string]interface{}
		if err := json.Unmarshal([]byte(req.Validation), &validation); err != nil {
			return nil, fmt.Errorf("验证规则格式错误: %w", err)
		}
	}

	// 验证解锁规则 JSON（如果提供）
	if req.DefaultUnlockRules != "" {
		var unlockRules map[string]interface{}
		if err := json.Unmarshal([]byte(req.DefaultUnlockRules), &unlockRules); err != nil {
			return nil, fmt.Errorf("解锁规则格式错误: %w", err)
		}
	}

	template := &model.ProfileFieldTemplate{
		FieldKey:           req.FieldKey,
		FieldName:          req.FieldName,
		FieldType:          req.FieldType,
		IsRequired:         req.IsRequired,
		IsSearchable:       req.IsSearchable,
		IsPublic:           req.IsPublic,
		DefaultValue:       req.DefaultValue,
		Options:            req.Options,
		Validation:         req.Validation,
		DisplayOrder:       req.DisplayOrder,
		Icon:               req.Icon,
		Description:        req.Description,
		DefaultUnlockRules: req.DefaultUnlockRules,
		Category:           req.Category,
		IsActive:           true,
	}

	if err := s.templateRepo.Create(template); err != nil {
		return nil, fmt.Errorf("创建字段模板失败: %w", err)
	}

	return template.ToResponse(), nil
}

// GetTemplateByID 根据 ID 获取字段模板
func (s *profileFieldTemplateService) GetTemplateByID(ctx context.Context, id int) (*model.ProfileFieldTemplateResponse, error) {
	template, err := s.templateRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("字段模板不存在")
		}
		return nil, fmt.Errorf("查询字段模板失败: %w", err)
	}

	return template.ToResponse(), nil
}

// GetTemplateByFieldKey 根据字段标识获取字段模板
func (s *profileFieldTemplateService) GetTemplateByFieldKey(ctx context.Context, fieldKey string) (*model.ProfileFieldTemplateResponse, error) {
	template, err := s.templateRepo.GetByFieldKey(fieldKey)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("字段模板不存在")
		}
		return nil, fmt.Errorf("查询字段模板失败: %w", err)
	}

	return template.ToResponse(), nil
}

// UpdateTemplate 更新字段模板
func (s *profileFieldTemplateService) UpdateTemplate(ctx context.Context, id int, req *model.UpdateProfileFieldTemplateRequest) (*model.ProfileFieldTemplateResponse, error) {
	template, err := s.templateRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("字段模板不存在")
		}
		return nil, fmt.Errorf("查询字段模板失败: %w", err)
	}

	// 更新字段
	if req.FieldName != "" {
		template.FieldName = req.FieldName
	}
	if req.FieldType != "" {
		template.FieldType = req.FieldType
	}
	if req.IsRequired != nil {
		template.IsRequired = *req.IsRequired
	}
	if req.IsSearchable != nil {
		template.IsSearchable = *req.IsSearchable
	}
	if req.IsPublic != nil {
		template.IsPublic = *req.IsPublic
	}
	if req.DefaultValue != "" {
		template.DefaultValue = req.DefaultValue
	}
	if req.Options != "" {
		// 验证选项配置 JSON
		var options map[string]interface{}
		if err := json.Unmarshal([]byte(req.Options), &options); err != nil {
			return nil, fmt.Errorf("选项配置格式错误: %w", err)
		}
		template.Options = req.Options
	}
	if req.Validation != "" {
		// 验证验证规则 JSON
		var validation map[string]interface{}
		if err := json.Unmarshal([]byte(req.Validation), &validation); err != nil {
			return nil, fmt.Errorf("验证规则格式错误: %w", err)
		}
		template.Validation = req.Validation
	}
	if req.DisplayOrder != nil {
		template.DisplayOrder = *req.DisplayOrder
	}
	if req.Icon != "" {
		template.Icon = req.Icon
	}
	if req.Description != "" {
		template.Description = req.Description
	}
	if req.DefaultUnlockRules != "" {
		// 验证解锁规则 JSON
		var unlockRules map[string]interface{}
		if err := json.Unmarshal([]byte(req.DefaultUnlockRules), &unlockRules); err != nil {
			return nil, fmt.Errorf("解锁规则格式错误: %w", err)
		}
		template.DefaultUnlockRules = req.DefaultUnlockRules
	}
	if req.Category != "" {
		template.Category = req.Category
	}
	if req.IsActive != nil {
		template.IsActive = *req.IsActive
	}

	if err := s.templateRepo.Update(template); err != nil {
		return nil, fmt.Errorf("更新字段模板失败: %w", err)
	}

	return template.ToResponse(), nil
}

// DeleteTemplate 删除字段模板
func (s *profileFieldTemplateService) DeleteTemplate(ctx context.Context, id int) error {
	_, err := s.templateRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("字段模板不存在")
		}
		return fmt.Errorf("查询字段模板失败: %w", err)
	}

	return s.templateRepo.Delete(id)
}

// ListTemplates 获取字段模板列表
func (s *profileFieldTemplateService) ListTemplates(ctx context.Context, category string, fieldType string, isActive *bool, page, pageSize int) ([]*model.ProfileFieldTemplateResponse, int64, error) {
	offset := (page - 1) * pageSize
	templates, total, err := s.templateRepo.List(category, fieldType, isActive, offset, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("查询字段模板列表失败: %w", err)
	}

	responses := make([]*model.ProfileFieldTemplateResponse, len(templates))
	for i, template := range templates {
		responses[i] = template.ToResponse()
	}

	return responses, total, nil
}

// GetTemplatesByCategory 根据分类获取字段模板列表
func (s *profileFieldTemplateService) GetTemplatesByCategory(ctx context.Context, category string) ([]*model.ProfileFieldTemplateResponse, error) {
	templates, err := s.templateRepo.GetByCategory(category)
	if err != nil {
		return nil, fmt.Errorf("查询字段模板列表失败: %w", err)
	}

	responses := make([]*model.ProfileFieldTemplateResponse, len(templates))
	for i, template := range templates {
		responses[i] = template.ToResponse()
	}

	return responses, nil
}

// ApplyTemplateToUser 将字段模板应用到用户
// 在 profile_fields 表中创建一条记录
func (s *profileFieldTemplateService) ApplyTemplateToUser(ctx context.Context, templateID, userID int) (*ApplyTemplateResult, error) {
	template, err := s.templateRepo.GetByID(templateID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("字段模板不存在")
		}
		return nil, fmt.Errorf("查询字段模板失败: %w", err)
	}

	if !template.IsActive {
		return nil, errors.New("字段模板未启用")
	}

	// 检查用户是否已经应用过该字段模板
	existingField, err := s.fieldRepo.GetByUserIDAndFieldKey(userID, template.FieldKey)
	if err == nil && existingField != nil {
		return nil, fmt.Errorf("用户已应用该字段模板，字段ID: %d", existingField.ID)
	} else if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("查询用户字段失败: %w", err)
	}

	// 将模板应用到用户，创建 profile_field 记录
	field := template.ApplyToUser(userID)
	if err := s.fieldRepo.Create(field); err != nil {
		return nil, fmt.Errorf("创建用户字段失败: %w", err)
	}

	// TODO: 如果模板有默认解锁规则，需要创建 unlock_rule 记录

	return &ApplyTemplateResult{
		FieldID:   field.ID,
		FieldKey:  field.FieldKey,
		FieldName: field.FieldName,
		Message:   "字段模板应用成功",
	}, nil
}

// ApplyTemplatesToUser 批量将字段模板应用到用户
func (s *profileFieldTemplateService) ApplyTemplatesToUser(ctx context.Context, templateIDs []int, userID int) (*ApplyTemplatesResult, error) {
	result := &ApplyTemplatesResult{
		AppliedFields: make([]ApplyTemplateResult, 0),
		TotalCount:    len(templateIDs),
		SuccessCount:  0,
		FailedCount:   0,
	}

	for _, templateID := range templateIDs {
		applyResult, err := s.ApplyTemplateToUser(ctx, templateID, userID)
		if err != nil {
			result.FailedCount++
			continue
		}
		result.SuccessCount++
		result.AppliedFields = append(result.AppliedFields, *applyResult)
	}

	if result.SuccessCount > 0 {
		result.Message = fmt.Sprintf("成功应用 %d 个字段模板", result.SuccessCount)
	} else {
		result.Message = "没有成功应用任何字段模板"
	}

	return result, nil
}
