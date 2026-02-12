package handler

import (
	"net/http"
	"strconv"

	"github.com/deantook/dove/internal/model"
	"github.com/deantook/dove/internal/service"

	"github.com/gin-gonic/gin"
)

var userServiceInstance service.UserService

// getUserService 获取用户服务实例（延迟初始化）
func getUserService() service.UserService {
	if userServiceInstance == nil {
		userServiceInstance = service.NewUserService()
	}
	return userServiceInstance
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Username string `json:"username" binding:"required"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Phone      string `json:"phone"`
	Nickname   string `json:"nickname"`
	Avatar     string `json:"avatar"`
	Status     int    `json:"status"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}

// ListUsersResponse 用户列表响应
type ListUsersResponse struct {
	Users      []UserResponse `json:"users"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// CreateUser 创建用户
// @Summary      创建用户
// @Description  根据用户名创建新用户
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        body  body      handler.CreateUserRequest  true  "创建用户请求"
// @Success      200   {object}  object  "success, message, data(user)"
// @Failure      400   {object}  object  "success, message"
// @Router       /api/v1/users [post]
func CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	user, err := getUserService().CreateUser(req.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "用户创建成功",
		"data":    toUserResponse(user),
	})
}

// GetUser 获取用户详情
// @Summary      获取用户详情
// @Description  根据用户 ID 获取用户信息
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "用户 ID"
// @Success      200  {object}  object  "success, data(user)"
// @Failure      400  {object}  object  "success, message"
// @Failure      404  {object}  object  "success, message"
// @Router       /api/v1/users/{id} [get]
func GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的用户ID",
		})
		return
	}

	user, err := getUserService().GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    toUserResponse(user),
	})
}

// UpdateUser 更新用户
// @Summary      更新用户
// @Description  根据用户 ID 更新用户信息
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      int  true  "用户 ID"
// @Param        body  body      handler.UpdateUserRequest  true  "更新用户请求"
// @Success      200   {object}  object  "success, message, data(user)"
// @Failure      400   {object}  object  "success, message"
// @Router       /api/v1/users/{id} [put]
func UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的用户ID",
		})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	user, err := getUserService().UpdateUser(id, req.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "用户更新成功",
		"data":    toUserResponse(user),
	})
}

// DeleteUser 删除用户
// @Summary      删除用户
// @Description  根据用户 ID 删除用户
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "用户 ID"
// @Success      200  {object}  object  "success, message"
// @Failure      400  {object}  object  "success, message"
// @Router       /api/v1/users/{id} [delete]
func DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的用户ID",
		})
		return
	}

	if err := getUserService().DeleteUser(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "用户删除成功",
	})
}

// ListUsers 获取用户列表（分页）
// @Summary      用户列表
// @Description  分页获取用户列表
// @Tags         users
// @Produce      json
// @Param        page       query     int  false  "页码"       default(1)
// @Param        page_size  query     int  false  "每页条数"   default(10)
// @Success      200        {object}  handler.ListUsersResponse
// @Failure      500        {object}  object  "success, message"
// @Router       /api/v1/users [get]
func ListUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	users, total, err := getUserService().ListUsers(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = toUserResponse(&user)
	}

	c.JSON(http.StatusOK, ListUsersResponse{
		Users:      userResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

// toUserResponse 转换为用户响应结构
func toUserResponse(user *model.User) UserResponse {
	response := UserResponse{
		ID:         user.ID,
		Username:   user.Username,
		Phone:      user.Phone,
		Nickname:   user.Nickname,
		Avatar:     user.Avatar,
		Status:     user.Status,
		CreateTime: user.CreateTime.Format("2006-01-02 15:04:05"),
	}
	// 处理更新时间，如果为零值则不显示
	if !user.UpdateTime.IsZero() {
		response.UpdateTime = user.UpdateTime.Format("2006-01-02 15:04:05")
	}
	return response
}
