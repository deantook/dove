package handler

import (
	"strconv"

	"dove/internal/domain"
	"dove/internal/model"
	"dove/pkg/pagination"
	"dove/pkg/password"
	"dove/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService domain.UserService
}

func NewUserHandler(userService domain.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Create CreateUser godoc
// @Summary      创建用户
// @Description  创建新用户
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        user body domain.CreateUserRequest true "用户信息"
// @Success      201  {object}  response.Response{data=model.User}
// @Failure      400  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /users [post]
func (h *UserHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var req domain.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// 验证密码强度
	if err := password.ValidatePassword(req.Password); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 加密密码
	hashedPassword, err := password.HashPassword(req.Password)
	if err != nil {
		response.InternalServerError(c, "Failed to process password")
		return
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword, // 使用加密后的密码
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		Status:   req.Status,
	}

	if err := h.userService.Create(ctx, user); err != nil {
		response.DatabaseError(c, err.Error())
		return
	}

	// 不返回密码
	user.Password = ""
	response.Created(c, user)
}

// GetByID GetUserByID godoc
// @Summary      获取用户详情
// @Description  根据ID获取用户详细信息
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "用户ID"
// @Success      200  {object}  response.Response{data=model.User}
// @Failure      400  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid ID")
		return
	}

	user, err := h.userService.GetByID(ctx, uint(id))
	if err != nil {
		response.NotFound(c, "User not found")
		return
	}

	// 不返回密码
	user.Password = ""
	response.Success(c, user)
}

// GetAll GetAllUsers godoc
// @Summary      获取所有用户
// @Description  获取所有用户列表（支持分页、排序和搜索）
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        page       query     int    false  "页码，默认为1"  minimum(1)
// @Param        page_size  query     int    false  "每页大小，默认为10，最大100"  minimum(1) maximum(100)
// @Param        sort_by    query     string false  "排序字段：id, username, email, nickname, status, created_at, updated_at"
// @Param        sort_order query     string false  "排序方向：asc, desc，默认为desc"
// @Param        keyword    query     string false  "搜索关键词"
// @Param        search_by  query     string false  "搜索字段：username, email, nickname，不指定则在所有字段中搜索"
// @Success      200  {object}  response.Response{data=pagination.PageResponse{data=[]model.User}}
// @Failure      400  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /users/ [get]
func (h *UserHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()

	// 解析分页参数
	pageReq := pagination.ParsePageRequest(c)

	// 使用分页获取用户列表
	pageResponse, err := h.userService.GetAllWithPagination(ctx, pageReq)
	if err != nil {
		response.DatabaseError(c, err.Error())
		return
	}

	response.Success(c, pageResponse)
}

// Update UpdateUser godoc
// @Summary      更新用户
// @Description  根据ID更新用户信息
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "用户ID"
// @Param        user body domain.UpdateUserRequest true "用户更新信息"
// @Success      200  {object}  response.Response{data=model.User}
// @Failure      400  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid ID")
		return
	}

	var req domain.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	user, err := h.userService.GetByID(ctx, uint(id))
	if err != nil {
		response.NotFound(c, "User not found")
		return
	}

	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Password != "" {
		// 验证密码强度
		if err := password.ValidatePassword(req.Password); err != nil {
			response.BadRequest(c, err.Error())
			return
		}

		// 加密密码
		hashedPassword, err := password.HashPassword(req.Password)
		if err != nil {
			response.InternalServerError(c, "Failed to process password")
			return
		}
		user.Password = hashedPassword
	}
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Status != nil {
		user.Status = *req.Status
	}

	if err := h.userService.Update(ctx, user); err != nil {
		response.DatabaseError(c, err.Error())
		return
	}

	// 不返回密码
	user.Password = ""
	response.Success(c, user)
}

// Delete DeleteUser godoc
// @Summary      删除用户
// @Description  根据ID删除用户
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "用户ID"
// @Success      200  {object}  response.Response{data=map[string]interface{}}
// @Failure      400  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.userService.Delete(ctx, uint(id)); err != nil {
		response.DatabaseError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "User deleted successfully"})
}
