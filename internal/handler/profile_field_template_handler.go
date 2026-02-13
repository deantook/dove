package handler

import (
	"net/http"
	"strconv"

	"github.com/deantook/dove/internal/model"
	"github.com/deantook/dove/internal/service"
	"github.com/deantook/dove/pkg/response"
	"github.com/gin-gonic/gin"
)

// ProfileFieldTemplateHandler 系统资料字段模板处理器
type ProfileFieldTemplateHandler struct {
	templateService service.ProfileFieldTemplateService
}

// NewProfileFieldTemplateHandler 创建系统资料字段模板处理器实例
func NewProfileFieldTemplateHandler(templateService service.ProfileFieldTemplateService) *ProfileFieldTemplateHandler {
	return &ProfileFieldTemplateHandler{
		templateService: templateService,
	}
}

// CreateTemplate 创建字段模板
// @Summary 创建字段模板
// @Description 创建新的系统资料字段模板（管理员）
// @Tags profile-field-templates
// @Accept json
// @Produce json
// @Param template body model.CreateProfileFieldTemplateRequest true "字段模板信息"
// @Success 201 {object} response.Response{data=model.ProfileFieldTemplateResponse}
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/profile/field-templates [post]
func (h *ProfileFieldTemplateHandler) CreateTemplate(c *gin.Context) {
	var req model.CreateProfileFieldTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithCode(c, http.StatusBadRequest, 400, "参数错误", err.Error())
		return
	}

	template, err := h.templateService.CreateTemplate(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithCode(c, http.StatusCreated, "创建成功", template)
}

// GetTemplate 获取字段模板详情
// @Summary 获取字段模板详情
// @Description 根据 ID 获取字段模板详情
// @Tags profile-field-templates
// @Produce json
// @Param id path int true "字段模板 ID"
// @Success 200 {object} response.Response{data=model.ProfileFieldTemplateResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/profile/field-templates/{id} [get]
func (h *ProfileFieldTemplateHandler) GetTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorWithCode(c, http.StatusBadRequest, 400, "无效的字段模板 ID", err.Error())
		return
	}

	template, err := h.templateService.GetTemplateByID(c.Request.Context(), int(id))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "获取成功", template)
}

// GetTemplateByFieldKey 根据字段标识获取字段模板
// @Summary 根据字段标识获取字段模板
// @Description 根据字段标识获取字段模板详情
// @Tags profile-field-templates
// @Produce json
// @Param key path string true "字段标识"
// @Success 200 {object} response.Response{data=model.ProfileFieldTemplateResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/profile/field-templates/key/{key} [get]
func (h *ProfileFieldTemplateHandler) GetTemplateByFieldKey(c *gin.Context) {
	fieldKey := c.Param("key")
	if fieldKey == "" {
		response.ErrorWithCode(c, http.StatusBadRequest, 400, "字段标识不能为空", "")
		return
	}

	template, err := h.templateService.GetTemplateByFieldKey(c.Request.Context(), fieldKey)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "获取成功", template)
}

// UpdateTemplate 更新字段模板
// @Summary 更新字段模板
// @Description 更新字段模板信息（管理员）
// @Tags profile-field-templates
// @Accept json
// @Produce json
// @Param id path int true "字段模板 ID"
// @Param template body model.UpdateProfileFieldTemplateRequest true "字段模板信息"
// @Success 200 {object} response.Response{data=model.ProfileFieldTemplateResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/profile/field-templates/{id} [put]
func (h *ProfileFieldTemplateHandler) UpdateTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorWithCode(c, http.StatusBadRequest, 400, "无效的字段模板 ID", err.Error())
		return
	}

	var req model.UpdateProfileFieldTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithCode(c, http.StatusBadRequest, 400, "参数错误", err.Error())
		return
	}

	template, err := h.templateService.UpdateTemplate(c.Request.Context(), int(id), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "更新成功", template)
}

// DeleteTemplate 删除字段模板
// @Summary 删除字段模板
// @Description 删除字段模板（管理员，软删除）
// @Tags profile-field-templates
// @Produce json
// @Param id path int true "字段模板 ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/profile/field-templates/{id} [delete]
func (h *ProfileFieldTemplateHandler) DeleteTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorWithCode(c, http.StatusBadRequest, 400, "无效的字段模板 ID", err.Error())
		return
	}

	if err := h.templateService.DeleteTemplate(c.Request.Context(), int(id)); err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

// ListTemplates 获取字段模板列表
// @Summary 获取字段模板列表
// @Description 分页获取字段模板列表
// @Tags profile-field-templates
// @Produce json
// @Param category query string false "字段分类"
// @Param field_type query string false "字段类型"
// @Param is_active query bool false "是否启用"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=response.ListResponse{list=[]model.ProfileFieldTemplateResponse}}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/profile/field-templates [get]
func (h *ProfileFieldTemplateHandler) ListTemplates(c *gin.Context) {
	category := c.Query("category")
	fieldType := c.Query("field_type")
	isActiveStr := c.Query("is_active")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var isActive *bool
	if isActiveStr != "" {
		active, err := strconv.ParseBool(isActiveStr)
		if err == nil {
			isActive = &active
		}
	}

	templates, total, err := h.templateService.ListTemplates(c.Request.Context(), category, fieldType, isActive, page, pageSize)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessList(c, templates, total, page, pageSize)
}

// GetTemplatesByCategory 根据分类获取字段模板列表
// @Summary 根据分类获取字段模板列表
// @Description 根据分类获取启用的字段模板列表
// @Tags profile-field-templates
// @Produce json
// @Param category path string true "字段分类"
// @Success 200 {object} response.Response{data=[]model.ProfileFieldTemplateResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/profile/field-templates/category/{category} [get]
func (h *ProfileFieldTemplateHandler) GetTemplatesByCategory(c *gin.Context) {
	category := c.Param("category")
	if category == "" {
		response.ErrorWithCode(c, http.StatusBadRequest, 400, "分类不能为空", "")
		return
	}

	templates, err := h.templateService.GetTemplatesByCategory(c.Request.Context(), category)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "获取成功", templates)
}

// ApplyTemplateToUser 将字段模板应用到用户
// @Summary 应用字段模板到用户
// @Description 将字段模板应用到当前用户，在 profile_fields 表中创建一条记录
// @Tags profile-field-templates
// @Accept json
// @Produce json
// @Param id path int true "字段模板 ID"
// @Param request body map[string]interface{} false "请求参数" example({"user_id": 123})
// @Success 200 {object} response.Response{data=service.ApplyTemplateResult}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/profile/field-templates/{id}/apply [post]
func (h *ProfileFieldTemplateHandler) ApplyTemplateToUser(c *gin.Context) {
	idStr := c.Param("id")
	templateID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorWithCode(c, http.StatusBadRequest, 400, "无效的字段模板 ID", err.Error())
		return
	}

	// 从请求中获取 user_id，如果没有则从 token 中获取当前用户ID
	var req struct {
		UserID int `json:"user_id"`
	}
	userID := 0
	if err := c.ShouldBindJSON(&req); err == nil {
		userID = req.UserID
	}

	// TODO: 从 JWT token 中获取当前用户ID（如果 userID 为 0）
	// 这里暂时使用 userID，实际应该从中间件中获取
	if userID == 0 {
		response.ErrorWithCode(c, http.StatusBadRequest, 400, "用户ID不能为空", "")
		return
	}

	result, err := h.templateService.ApplyTemplateToUser(c.Request.Context(), int(templateID), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "应用成功", result)
}

// ApplyTemplatesToUser 批量将字段模板应用到用户
// @Summary 批量应用字段模板到用户
// @Description 批量将字段模板应用到当前用户
// @Tags profile-field-templates
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "请求参数" example({"template_ids": [1, 2, 3], "user_id": 123})
// @Success 200 {object} response.Response{data=service.ApplyTemplatesResult}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/profile/field-templates/apply [post]
func (h *ProfileFieldTemplateHandler) ApplyTemplatesToUser(c *gin.Context) {
	var req struct {
		TemplateIDs []int `json:"template_ids" binding:"required"`
		UserID      int   `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithCode(c, http.StatusBadRequest, 400, "参数错误", err.Error())
		return
	}

	userID := req.UserID
	// TODO: 从 JWT token 中获取当前用户ID（如果 userID 为 0）
	if userID == 0 {
		response.ErrorWithCode(c, http.StatusBadRequest, 400, "用户ID不能为空", "")
		return
	}

	result, err := h.templateService.ApplyTemplatesToUser(c.Request.Context(), req.TemplateIDs, userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "批量应用完成", result)
}
