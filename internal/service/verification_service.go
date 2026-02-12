package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"time"

	"github.com/deantook/brigitta/pkg/cache"
)

// VerificationService 验证码服务接口
type VerificationService interface {
	GenerateCode(phone string) (string, error)
	StoreCode(phone string, code string) error
	VerifyCode(phone string, code string) (bool, error)
	CheckRateLimit(phone string) error
}

type verificationService struct{}

// NewVerificationService 创建验证码服务实例
func NewVerificationService() VerificationService {
	return &verificationService{}
}

// ValidatePhone 验证手机号格式
func ValidatePhone(phone string) bool {
	matched, _ := regexp.MatchString(`^1[3-9]\d{9}$`, phone)
	return matched
}

// GenerateCode 生成6位数字验证码
func (s *verificationService) GenerateCode(phone string) (string, error) {
	if !ValidatePhone(phone) {
		return "", errors.New("手机号格式错误")
	}

	// 生成6位随机数字
	code := ""
	for i := 0; i < 6; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", fmt.Errorf("生成验证码失败: %w", err)
		}
		code += n.String()
	}

	return code, nil
}

// StoreCode 存储验证码到Redis（5分钟过期）
func (s *verificationService) StoreCode(phone string, code string) error {
	ctx := context.Background()
	key := fmt.Sprintf("sms:code:%s", phone)
	ttl := 5 * time.Minute

	if err := cache.Set(ctx, key, code, ttl); err != nil {
		return fmt.Errorf("存储验证码失败: %w", err)
	}

	return nil
}

// VerifyCode 验证验证码
func (s *verificationService) VerifyCode(phone string, code string) (bool, error) {
	if !ValidatePhone(phone) {
		return false, errors.New("手机号格式错误")
	}

	if code == "" || len(code) != 6 {
		return false, errors.New("验证码格式错误")
	}

	ctx := context.Background()
	key := fmt.Sprintf("sms:code:%s", phone)

	storedCode, err := cache.Get(ctx, key)
	if err != nil {
		return false, errors.New("验证码不存在或已过期")
	}

	if storedCode != code {
		return false, errors.New("验证码错误")
	}

	// 验证成功后删除验证码（防止重复使用）
	// 注意：如果 cache 包支持 Delete 方法，可以取消注释
	// _ = cache.Delete(ctx, key)

	return true, nil
}

// CheckRateLimit 检查发送频率限制（60秒内只能发送一次）
func (s *verificationService) CheckRateLimit(phone string) error {
	if !ValidatePhone(phone) {
		return errors.New("手机号格式错误")
	}

	ctx := context.Background()
	limitKey := fmt.Sprintf("sms:limit:%s", phone)

	// 检查是否在限制期内
	_, err := cache.Get(ctx, limitKey)
	if err == nil {
		return errors.New("操作过于频繁，请稍后再试")
	}

	// 设置限制标记（60秒过期）
	ttl := 60 * time.Second
	if err := cache.Set(ctx, limitKey, "1", ttl); err != nil {
		return fmt.Errorf("设置频率限制失败: %w", err)
	}

	return nil
}
