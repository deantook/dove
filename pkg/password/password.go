package password

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 使用 bcrypt 加密密码
func HashPassword(password string) (string, error) {
	// 使用默认的 cost 值 (10)，可以根据需要调整
	// cost 值越高，加密越安全，但性能消耗也越大
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword 验证密码是否匹配
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// HashPasswordWithCost 使用指定的 cost 值加密密码
func HashPasswordWithCost(password string, cost int) (string, error) {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ValidatePassword 验证密码强度
func ValidatePassword(password string) error {
	// 这里可以添加密码强度验证逻辑
	// 例如：长度、复杂度等要求

	if len(password) < 6 {
		return ErrPasswordTooShort
	}

	return nil
}

// 密码相关错误
var (
	ErrPasswordTooShort = &PasswordError{Message: "密码长度至少6位"}
)

// PasswordError 密码相关错误
type PasswordError struct {
	Message string
}

func (e *PasswordError) Error() string {
	return e.Message
}
