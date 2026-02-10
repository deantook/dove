package handlers

import (
	"dove/database"
	"dove/models"
	"dove/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// Register 用户注册
func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	result := database.DB.Where("username = ? AND deleted = 0", req.Username).First(&existingUser)
	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{
			"code":    409,
			"message": "用户名已存在",
		})
		return
	}

	// 如果提供了邮箱，检查邮箱是否已存在
	if req.Email != "" {
		result = database.DB.Where("email = ? AND deleted = 0", req.Email).First(&existingUser)
		if result.Error == nil {
			c.JSON(http.StatusConflict, gin.H{
				"code":    409,
				"message": "邮箱已被注册",
			})
			return
		}
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "密码加密失败",
			"error":   err.Error(),
		})
		return
	}

	// 创建用户
	user := models.User{
		Username: req.Username,
		Password: hashedPassword,
		Status:   1, // 默认正常状态
	}

	if req.Email != "" {
		user.Email = &req.Email
	}
	if req.Phone != "" {
		user.Phone = &req.Phone
	}
	if req.Nickname != "" {
		user.Nickname = &req.Nickname
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "注册失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "注册成功",
		"data":    user.ToUserInfo(),
	})
}

// Login 用户登录
func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 查找用户
	var user models.User
	result := database.DB.Where("username = ? AND deleted = 0", req.Username).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户名或密码错误",
		})
		return
	}

	// 检查用户状态
	if user.Status != 1 {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "账户已被禁用或锁定",
		})
		return
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户名或密码错误",
		})
		return
	}

	// 更新最后登录时间和IP
	now := time.Now()
	user.LastLoginTime = &now
	clientIP := c.ClientIP()
	user.LastLoginIP = &clientIP
	database.DB.Save(&user)

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成token失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"data": models.LoginResponse{
			Token: token,
			User:  user.ToUserInfo(),
		},
	})
}
