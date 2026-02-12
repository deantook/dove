package handler

import (
	"net/http"

	"github.com/deantook/dove/internal/service"

	"github.com/gin-gonic/gin"
)

var verificationServiceInstance service.VerificationService

// getVerificationService 获取验证码服务实例（延迟初始化）
func getVerificationService() service.VerificationService {
	if verificationServiceInstance == nil {
		verificationServiceInstance = service.NewVerificationService()
	}
	return verificationServiceInstance
}

// SendVerificationCodeRequest 发送验证码请求
type SendVerificationCodeRequest struct {
	Phone string `json:"phone" binding:"required"`
}

// SendVerificationCodeResponse 发送验证码响应
type SendVerificationCodeResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"` // 开发阶段直接返回验证码
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Code     string `json:"code" binding:"required"`
	Nickname string `json:"nickname"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

// AuthResponse 认证响应
type AuthResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SendVerificationCode 发送验证码
// @Summary      发送验证码
// @Description  向指定手机号发送验证码（开发阶段直接返回验证码）
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      handler.SendVerificationCodeRequest  true  "发送验证码请求"
// @Success      200   {object}  handler.SendVerificationCodeResponse
// @Failure      400   {object}  handler.SendVerificationCodeResponse
// @Failure      500   {object}  handler.SendVerificationCodeResponse
// @Router       /api/v1/auth/send-code [post]
func SendVerificationCode(c *gin.Context) {
	var req SendVerificationCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, SendVerificationCodeResponse{
			Success: false,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	// 检查频率限制
	verificationService := getVerificationService()
	if err := verificationService.CheckRateLimit(req.Phone); err != nil {
		c.JSON(http.StatusBadRequest, SendVerificationCodeResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// 生成验证码
	code, err := verificationService.GenerateCode(req.Phone)
	if err != nil {
		c.JSON(http.StatusBadRequest, SendVerificationCodeResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// 存储验证码
	if err := verificationService.StoreCode(req.Phone, code); err != nil {
		c.JSON(http.StatusInternalServerError, SendVerificationCodeResponse{
			Success: false,
			Message: "发送验证码失败: " + err.Error(),
		})
		return
	}

	// 开发阶段直接返回验证码
	c.JSON(http.StatusOK, SendVerificationCodeResponse{
		Success: true,
		Message: "验证码已发送",
		Code:    code, // 开发阶段直接返回，生产环境应移除
	})
}

// Register 用户注册
// @Summary      用户注册
// @Description  使用手机号和验证码注册新用户
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      handler.RegisterRequest  true  "注册请求"
// @Success      200   {object}  handler.AuthResponse
// @Failure      400   {object}  handler.AuthResponse
// @Router       /api/v1/auth/register [post]
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, AuthResponse{
			Success: false,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	user, err := getUserService().RegisterByPhone(req.Phone, req.Code, req.Nickname)
	if err != nil {
		c.JSON(http.StatusBadRequest, AuthResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Success: true,
		Message: "注册成功",
		Data:    toUserResponse(user),
	})
}

// Login 用户登录
// @Summary      用户登录
// @Description  使用手机号和验证码登录
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      handler.LoginRequest  true  "登录请求"
// @Success      200   {object}  handler.AuthResponse
// @Failure      400   {object}  handler.AuthResponse
// @Router       /api/v1/auth/login [post]
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, AuthResponse{
			Success: false,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	user, err := getUserService().LoginByPhone(req.Phone, req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, AuthResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Success: true,
		Message: "登录成功",
		Data:    toUserResponse(user),
	})
}
