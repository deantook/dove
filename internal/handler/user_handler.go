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
	CreateTime string `json:"create_time"`
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
// POST /api/users
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
// GET /api/users/:id
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
// PUT /api/users/:id
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
// DELETE /api/users/:id
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
// GET /api/users?page=1&page_size=10
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
	return UserResponse{
		ID:         user.ID,
		Username:   user.Username,
		CreateTime: user.CreateTime.Format("2006-01-02 15:04:05"),
	}
}
