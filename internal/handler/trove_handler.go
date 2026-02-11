package handler

import (
	"strconv"

	"dove/internal/domain"
	"dove/internal/model"
	"dove/pkg/pagination"
	"dove/pkg/response"

	"github.com/gin-gonic/gin"
)

type TroveHandler struct {
	troveService domain.TroveService
}

func NewTroveHandler(troveService domain.TroveService) *TroveHandler {
	return &TroveHandler{troveService: troveService}
}

// Create CreateTrove godoc
// @Summary      创建Trove
// @Description  创建新的Trove
// @Tags         troves
// @Accept       json
// @Produce      json
// @Param        trove body domain.CreateTroveRequest true "Trove信息"
// @Success      201  {object}  response.Response{data=model.Trove}
// @Failure      400  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /troves/ [post]
func (h *TroveHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var req domain.CreateTroveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	trove := &model.Trove{

		Title: req.Title,

		Description: req.Description,
	}

	if err := h.troveService.Create(ctx, trove); err != nil {
		response.DatabaseError(c, err.Error())
		return
	}

	response.Created(c, trove)
}

// GetByID GetTroveByID godoc
// @Summary      获取Trove详情
// @Description  根据ID获取Trove详细信息
// @Tags         troves
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "TroveID"
// @Success      200  {object}  response.Response{data=model.Trove}
// @Failure      400  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /troves/{id} [get]
func (h *TroveHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid ID")
		return
	}

	trove, err := h.troveService.GetByID(ctx, uint(id))
	if err != nil {
		response.NotFound(c, "Trove not found")
		return
	}

	response.Success(c, trove)
}

// GetAll GetAllTroves godoc
// @Summary      获取所有Trove
// @Description  获取所有Trove列表（支持分页、排序和搜索）
// @Tags         troves
// @Accept       json
// @Produce      json
// @Param        page       query     int    false  "页码，默认为1"  minimum(1)
// @Param        page_size  query     int    false  "每页大小，默认为10，最大100"  minimum(1) maximum(100)
// @Param        sort_by    query     string false  "排序字段：ID, "
// @Param        sort_order query     string false  "排序方向：asc, desc，默认为desc"
// @Param        keyword    query     string false  "搜索关键词"
// @Param        search_by  query     string false  "搜索字段：Title, Description, ，不指定则在所有字段中搜索"
// @Success      200  {object}  response.Response{data=pagination.PageResponse{data=[]model.Trove}}
// @Failure      400  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /troves/ [get]
func (h *TroveHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()

	// 解析分页参数
	pageReq := pagination.ParsePageRequest(c)

	// 使用分页获取Trove列表
	pageResponse, err := h.troveService.GetAllWithPagination(ctx, pageReq)
	if err != nil {
		response.DatabaseError(c, err.Error())
		return
	}

	response.Success(c, pageResponse)
}

// Update UpdateTrove godoc
// @Summary      更新Trove
// @Description  根据ID更新Trove信息
// @Tags         troves
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "TroveID"
// @Param        trove body domain.UpdateTroveRequest true "Trove更新信息"
// @Success      200  {object}  response.Response{data=model.Trove}
// @Failure      400  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /troves/{id} [put]
func (h *TroveHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid ID")
		return
	}

	var req domain.UpdateTroveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	trove, err := h.troveService.GetByID(ctx, uint(id))
	if err != nil {
		response.NotFound(c, "Trove not found")
		return
	}

	if req.Title != nil {
		trove.Title = *req.Title
	}

	if req.Description != nil {
		trove.Description = *req.Description
	}

	if err := h.troveService.Update(ctx, trove); err != nil {
		response.DatabaseError(c, err.Error())
		return
	}

	response.Success(c, trove)
}

// Delete DeleteTrove godoc
// @Summary      删除Trove
// @Description  根据ID删除Trove
// @Tags         troves
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "TroveID"
// @Success      200  {object}  response.Response{data=map[string]interface{}}
// @Failure      400  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /troves/{id} [delete]
func (h *TroveHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid ID")
		return
	}

	if err := h.troveService.Delete(ctx, uint(id)); err != nil {
		response.DatabaseError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "Trove deleted successfully"})
}
