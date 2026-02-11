package handler

import (
	"strconv"

	"dove/internal/domain"
	"dove/internal/model"
	"dove/pkg/pagination"
	"dove/pkg/response"

	"github.com/gin-gonic/gin"
)

type WeaponHandler struct {
	weaponService domain.WeaponService
}

func NewWeaponHandler(weaponService domain.WeaponService) *WeaponHandler {
	return &WeaponHandler{weaponService: weaponService}
}

// Create CreateWeapon godoc
// @Summary      创建Weapon
// @Description  创建新的Weapon
// @Tags         武器
// @Accept       json
// @Produce      json
// @Param        weapon body domain.CreateWeaponRequest true "Weapon信息"
// @Success      201  {object}  response.Response{data=model.Weapon}
// @Failure      400  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /weapons/ [post]
func (h *WeaponHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var req domain.CreateWeaponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	weapon := &model.Weapon{

		Name: req.Name,

		Level: req.Level,

		Content: req.Content,

		Type: req.Type,

		Story: req.Story,
	}

	if err := h.weaponService.Create(ctx, weapon); err != nil {
		response.DatabaseError(c, err.Error())
		return
	}

	response.Created(c, weapon)
}

// GetByID GetWeaponByID godoc
// @Summary      获取Weapon详情
// @Description  根据ID获取Weapon详细信息
// @Tags         武器
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "WeaponID"
// @Success      200  {object}  response.Response{data=model.Weapon}
// @Failure      400  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /weapons/{id} [get]
func (h *WeaponHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid ID")
		return
	}

	weapon, err := h.weaponService.GetByID(ctx, uint(id))
	if err != nil {
		response.NotFound(c, "Weapon not found")
		return
	}

	response.Success(c, weapon)
}

// GetAll GetAllWeapons godoc
// @Summary      获取所有Weapon
// @Description  获取所有Weapon列表（支持分页、排序和搜索）
// @Tags         武器
// @Accept       json
// @Produce      json
// @Param        page       query     int    false  "页码，默认为1"  minimum(1)
// @Param        page_size  query     int    false  "每页大小，默认为10，最大100"  minimum(1) maximum(100)
// @Param        sort_by    query     string false  "排序字段：ID, Name, "
// @Param        sort_order query     string false  "排序方向：asc, desc，默认为desc"
// @Param        keyword    query     string false  "搜索关键词"
// @Param        search_by  query     string false  "搜索字段：Name, ，不指定则在所有字段中搜索"
// @Success      200  {object}  response.Response{data=pagination.PageResponse{data=[]model.Weapon}}
// @Failure      400  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /weapons/ [get]
func (h *WeaponHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()

	// 解析分页参数
	pageReq := pagination.ParsePageRequest(c)

	// 使用分页获取Weapon列表
	pageResponse, err := h.weaponService.GetAllWithPagination(ctx, pageReq)
	if err != nil {
		response.DatabaseError(c, err.Error())
		return
	}

	response.Success(c, pageResponse)
}

// Update UpdateWeapon godoc
// @Summary      更新Weapon
// @Description  根据ID更新Weapon信息
// @Tags         武器
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "WeaponID"
// @Param        weapon body domain.UpdateWeaponRequest true "Weapon更新信息"
// @Success      200  {object}  response.Response{data=model.Weapon}
// @Failure      400  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /weapons/{id} [put]
func (h *WeaponHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid ID")
		return
	}

	var req domain.UpdateWeaponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	weapon, err := h.weaponService.GetByID(ctx, uint(id))
	if err != nil {
		response.NotFound(c, "Weapon not found")
		return
	}

	if req.Name != nil {
		weapon.Name = *req.Name
	}

	if req.Level != nil {
		weapon.Level = *req.Level
	}

	if req.Content != nil {
		weapon.Content = *req.Content
	}

	if req.Type != nil {
		weapon.Type = *req.Type
	}

	if req.Story != nil {
		weapon.Story = *req.Story
	}

	if err := h.weaponService.Update(ctx, weapon); err != nil {
		response.DatabaseError(c, err.Error())
		return
	}

	response.Success(c, weapon)
}

// Delete DeleteWeapon godoc
// @Summary      删除Weapon
// @Description  根据ID删除Weapon
// @Tags         武器
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "WeaponID"
// @Success      200  {object}  response.Response{data=map[string]interface{}}
// @Failure      400  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /weapons/{id} [delete]
func (h *WeaponHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid ID")
		return
	}

	if err := h.weaponService.Delete(ctx, uint(id)); err != nil {
		response.DatabaseError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "Weapon deleted successfully"})
}
