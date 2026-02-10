package handlers

import (
	"dove/internal/data/domain"
	"dove/pkg/common"
	"dove/pkg/database"
	"dove/pkg/utils"
	"errors"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PhoneRegisterRequest 手机号注册请求
type PhoneRegisterRequest struct {
	Phone    string `json:"phone" binding:"required"`    // 手机号，必填
	Password string `json:"password" binding:"required"` // 密码，必填
	Nickname string `json:"nickname"`                    // 昵称，可选
}

// PhoneLoginRequest 手机号登录请求
type PhoneLoginRequest struct {
	Phone    string `json:"phone" binding:"required"`    // 手机号，必填
	Password string `json:"password" binding:"required"` // 密码，必填
}

// validatePhone 验证手机号格式（中国大陆手机号）
func validatePhone(phone string) bool {
	// 去除空格
	phone = strings.TrimSpace(phone)
	// 中国大陆手机号正则：1开头，11位数字
	matched, _ := regexp.MatchString(`^1[3-9]\d{9}$`, phone)
	return matched
}

// PhoneRegister 手机号注册接口
func PhoneRegister(c *gin.Context) {
	var req PhoneRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 验证手机号格式
	if !validatePhone(req.Phone) {
		common.Error(c, common.ErrInvalidPhone.Code, common.ErrInvalidPhone.Message)
		return
	}

	// 验证密码长度
	if len(req.Password) < 6 || len(req.Password) > 50 {
		common.Error(c, common.ErrInvalidPassword.Code, common.ErrInvalidPassword.Message)
		return
	}

	// 检查手机号是否已存在
	var existingUser domain.User
	result := database.DB.Where("phone = ?", req.Phone).First(&existingUser)
	if result.Error == nil {
		common.Error(c, common.ErrPhoneExists.Code, common.ErrPhoneExists.Message)
		return
	}
	if result.Error != nil && !errors.Is(gorm.ErrRecordNotFound, result.Error) {
		common.InternalServerError(c, "查询用户失败")
		return
	}

	// 生成密码哈希
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		common.Error(c, common.ErrPasswordHash.Code, common.ErrPasswordHash.Message)
		return
	}

	// 设置默认昵称
	nickname := req.Nickname
	if nickname == "" {
		nickname = "用户" + req.Phone[len(req.Phone)-4:] // 使用手机号后4位作为默认昵称
	}

	// 创建用户
	user := domain.User{
		Phone:        &req.Phone,
		PasswordHash: passwordHash,
		Nickname:     nickname,
		Status:       1, // 正常状态
	}

	if err := database.DB.Create(&user).Error; err != nil {
		common.InternalServerError(c, "注册失败: "+err.Error())
		return
	}

	// 返回用户信息（不包含密码）
	common.SuccessWithMessage(c, "注册成功", gin.H{
		"id":       user.ID,
		"phone":    user.Phone,
		"nickname": user.Nickname,
		"status":   user.Status,
	})
}

// PhoneLogin 手机号登录接口
func PhoneLogin(c *gin.Context) {
	var req PhoneLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 验证手机号格式
	if !validatePhone(req.Phone) {
		common.Error(c, common.ErrInvalidPhone.Code, common.ErrInvalidPhone.Message)
		return
	}

	// 查找用户
	var user domain.User
	result := database.DB.Where("phone = ?", req.Phone).First(&user)
	if result.Error != nil {
		if errors.Is(gorm.ErrRecordNotFound, result.Error) {
			common.Error(c, common.ErrPhoneNotFound.Code, common.ErrPhoneNotFound.Message)
		} else {
			common.InternalServerError(c, "查询用户失败")
		}
		return
	}

	// 检查用户状态
	if user.Status != 1 {
		common.Error(c, common.ErrAccountDisabled.Code, common.ErrAccountDisabled.Message)
		return
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		common.Error(c, common.ErrPasswordError.Code, common.ErrPasswordError.Message)
		return
	}

	// 生成JWT token
	token, err := utils.GenerateToken(int64(user.ID), user.Nickname)
	if err != nil {
		common.Error(c, common.ErrTokenGenerate.Code, common.ErrTokenGenerate.Message)
		return
	}

	// 返回token和用户信息
	common.SuccessWithMessage(c, "登录成功", gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"phone":    user.Phone,
			"nickname": user.Nickname,
			"status":   user.Status,
		},
	})
}
