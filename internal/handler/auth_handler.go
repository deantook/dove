package handler

import (
	"net/http"
	"strings"
	"time"

	"dove/internal/domain"
	"dove/internal/model"
	"dove/pkg/jwt"
	"dove/pkg/logger"
	"dove/pkg/password"
	"dove/pkg/redis"
	"dove/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userService domain.UserService
}

func NewAuthHandler(userService domain.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

// Register godoc
// @Summary      用户注册
// @Description  新用户注册
// @Tags         认证与校验
// @Accept       json
// @Produce      json
// @Param        register body RegisterRequest true "注册信息"
// @Success      201  {object}  response.Response{data=model.User}
// @Failure      400  {object}  response.Response
// @Failure      409  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorWithTrace(ctx, "Register validation failed", "error", err.Error())
		response.ValidationError(c, err.Error())
		return
	}

	// 检查用户名是否已存在
	if _, err := h.userService.GetByUsername(ctx, req.Username); err == nil {
		logger.WarnWithTrace(ctx, "Username already exists", "username", req.Username)
		response.Error(c, http.StatusConflict, "Username already exists")
		return
	}

	// 检查邮箱是否已存在
	if _, err := h.userService.GetByEmail(ctx, req.Email); err == nil {
		logger.WarnWithTrace(ctx, "Email already exists", "email", req.Email)
		response.Error(c, http.StatusConflict, "Email already exists")
		return
	}

	// 验证密码强度
	if err := password.ValidatePassword(req.Password); err != nil {
		logger.WarnWithTrace(ctx, "Password validation failed", "error", err.Error(), "username", req.Username)
		response.BadRequest(c, err.Error())
		return
	}

	// 加密密码
	hashedPassword, err := password.HashPassword(req.Password)
	if err != nil {
		logger.ErrorWithTrace(ctx, "Failed to hash password", "error", err.Error(), "username", req.Username)
		response.InternalServerError(c, "Failed to process password")
		return
	}

	// 创建新用户
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword, // 使用加密后的密码
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		Status:   1, // 默认启用状态
	}

	if err := h.userService.Create(ctx, user); err != nil {
		logger.ErrorWithTrace(ctx, "Failed to create user", "error", err.Error(), "username", req.Username)
		response.DatabaseError(c, "Failed to create user")
		return
	}

	// 不返回密码
	user.Password = ""

	logger.InfoWithTrace(ctx, "User registered successfully", "user_id", user.ID, "username", user.Username)
	response.Created(c, user)
}

type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"john_doe"`
	Password string `json:"password" binding:"required" example:"123456"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required" example:"john_doe"`
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"123456"`
	Nickname string `json:"nickname" example:"John Doe"`
	Avatar   string `json:"avatar" example:"https://example.com/avatar.jpg"`
}

type LoginResponse struct {
	Token     string      `json:"token"`
	User      *model.User `json:"user"`
	ExpiresIn int64       `json:"expires_in"`
}

// Login godoc
// @Summary      用户登录
// @Description  用户登录并返回 JWT token
// @Tags         认证与校验
// @Accept       json
// @Produce      json
// @Param        login body LoginRequest true "登录信息"
// @Success      200  {object}  response.Response{data=LoginResponse}
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// 验证用户名和密码
	user, err := h.userService.GetByUsername(ctx, req.Username)
	if err != nil {
		response.Unauthorized(c, "Invalid credentials")
		return
	}

	// 使用 bcrypt 验证密码
	if !password.CheckPassword(req.Password, user.Password) {
		logger.WarnWithTrace(ctx, "Invalid password", "username", req.Username)
		response.Unauthorized(c, "Invalid credentials")
		return
	}

	// 生成 JWT token
	token, err := jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		response.InternalServerError(c, "Failed to generate token")
		return
	}

	// 将 token 存储到 Redis（可选，用于 token 黑名单功能）
	tokenKey := "token:" + token
	expiration := time.Duration(24) * time.Hour // 24小时
	if err := redis.Set(tokenKey, "valid", expiration); err != nil {
		// 这里只是记录错误，不影响登录流程
		// log.Printf("Failed to store token in Redis: %v", err)
	}

	// 不返回密码
	user.Password = ""

	loginResponse := LoginResponse{
		Token:     token,
		User:      user,
		ExpiresIn: int64(expiration.Seconds()),
	}

	response.Success(c, loginResponse)
}

// Logout godoc
// @Summary      用户登出
// @Description  用户登出，将 token 加入黑名单
// @Tags         认证与校验
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response{data=map[string]interface{}}
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// 从请求头获取 token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		response.BadRequest(c, "Authorization header is required")
		return
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		response.BadRequest(c, "Invalid authorization header format")
		return
	}

	token := tokenParts[1]

	// 将 token 加入黑名单（设置较短的过期时间）
	tokenKey := "token:" + token
	if err := redis.Set(tokenKey, "blacklisted", 24*time.Hour); err != nil {
		response.InternalServerError(c, "Failed to logout")
		return
	}

	response.Success(c, gin.H{"message": "Logout successful"})
}

// Profile godoc
// @Summary      获取用户信息
// @Description  获取当前登录用户的信息
// @Tags         认证与校验
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response{data=model.User}
// @Failure      401  {object}  response.Response
// @Router       /auth/profile [get]
func (h *AuthHandler) Profile(c *gin.Context) {
	ctx := c.Request.Context()
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	user, err := h.userService.GetByID(ctx, userID.(uint))
	if err != nil {
		response.NotFound(c, "User not found")
		return
	}

	// 不返回密码
	user.Password = ""

	response.Success(c, user)
}
