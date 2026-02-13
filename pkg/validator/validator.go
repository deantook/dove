package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// RegisterPhoneValidator 注册手机号验证器
func RegisterPhoneValidator(v *validator.Validate) error {
	return v.RegisterValidation("phone", validatePhone)
}

// validatePhone 验证手机号（中国手机号格式：11位数字，以1开头）
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	// 中国手机号正则：11位数字，以1开头，第二位为3-9
	matched, _ := regexp.MatchString(`^1[3-9]\d{9}$`, phone)
	return matched
}
