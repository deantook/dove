package handler

import (
	"net/http"
	"strconv"

	"github.com/deantook/dove/internal/model"
	"github.com/deantook/dove/internal/service"
	"github.com/deantook/dove/pkg/response"
	"github.com/gin-gonic/gin"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler 创建用户处理器实例
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUser 创建用户
// @Summary 创建用户
// @Description 创建新用户
// @Tags users
// @Accept json
// @Produce json
// @Param user body model.CreateUserRequest true "用户信息"
// @Success 201 {object} response.Response{data=model.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithCode(c, 500, 500, "无效的用户 ID", err.Error())
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithCode(c, http.StatusCreated, "创建成功", user)
}

// GetUser 获取用户详情
// @Summary 获取用户详情
// @Description 根据 ID 获取用户详情
// @Tags users
// @Produce json
// @Param id path int true "用户 ID"
// @Success 200 {object} response.Response{data=model.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorWithCode(c, 500, 500, "无效的用户 ID", err.Error())
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), int(id))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "获取成功", user)
}

// UpdateUser 更新用户
// @Summary 更新用户
// @Description 更新用户信息
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "用户 ID"
// @Param user body model.UpdateUserRequest true "用户信息"
// @Success 200 {object} response.Response{data=model.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorWithCode(c, 500, 500, "无效的用户 ID", err.Error())
		return
	}

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithCode(c, 500, 500, "参数错误", err.Error())
		return
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), int(id), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "更新成功", user)
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 删除用户（软删除）
// @Tags users
// @Produce json
// @Param id path int true "用户 ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorWithCode(c, 500, 500, "无效的用户 ID", err.Error())
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), int(id)); err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

// ListUsers 获取用户列表
// @Summary 获取用户列表
// @Description 分页获取用户列表
// @Tags users
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=response.ListResponse{list=[]model.UserResponse}}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	users, total, err := h.userService.ListUsers(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessList(c, users, total, page, pageSize)
}

// SendCode 发送验证码
// @Summary 发送验证码
// @Description 发送手机验证码（开发阶段返回验证码）
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.SendCodeRequest true "发送验证码请求"
// @Success 200 {object} response.Response{data=model.SendCodeResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/auth/send-code [post]
func (h *UserHandler) SendCode(c *gin.Context) {
	var req model.SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithCode(c, http.StatusBadRequest, 400, "参数错误", err.Error())
		return
	}

	result, err := h.userService.SendCode(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "验证码发送成功", result)
}

// LoginOrRegister 登录或注册
// @Summary 登录或注册
// @Description 使用手机号和验证码登录或注册（如果用户不存在则自动注册）
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "登录请求"
// @Success 200 {object} response.Response{data=model.LoginResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/auth/login [post]
func (h *UserHandler) LoginOrRegister(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithCode(c, http.StatusBadRequest, 400, "参数错误", err.Error())
		return
	}

	result, err := h.userService.LoginOrRegister(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "登录成功", result)
}
