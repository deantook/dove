package service

import (
	"errors"
	"fmt"

	"github.com/deantook/dove/internal/model"

	"github.com/deantook/brigitta/pkg/database"
	"gorm.io/gorm"
)

// UserService 用户服务接口
type UserService interface {
	CreateUser(username string) (*model.User, error)
	GetUserByID(id int) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	UpdateUser(id int, username string) (*model.User, error)
	DeleteUser(id int) error
	ListUsers(page, pageSize int) ([]model.User, int64, error)
}

type userService struct {
	db *gorm.DB
}

// NewUserService 创建用户服务实例
func NewUserService() UserService {
	return &userService{
		db: database.GetDB(),
	}
}

// CreateUser 创建用户
func (s *userService) CreateUser(username string) (*model.User, error) {
	if username == "" {
		return nil, errors.New("用户名不能为空")
	}

	// 检查用户名是否已存在
	var existingUser model.User
	if err := s.db.Where("username = ?", username).First(&existingUser).Error; err == nil {
		return nil, fmt.Errorf("用户名 %s 已存在", username)
	}

	user := &model.User{
		Username: username,
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	return user, nil
}

// GetUserByID 根据ID获取用户
func (s *userService) GetUserByID(id int) (*model.User, error) {
	var user model.User
	if err := s.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在，ID: %d", id)
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return &user, nil
}

// GetUserByUsername 根据用户名获取用户
func (s *userService) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在，用户名: %s", username)
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return &user, nil
}

// UpdateUser 更新用户
func (s *userService) UpdateUser(id int, username string) (*model.User, error) {
	if username == "" {
		return nil, errors.New("用户名不能为空")
	}

	// 检查用户是否存在
	user, err := s.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	// 检查新用户名是否已被其他用户使用
	var existingUser model.User
	if err := s.db.Where("username = ? AND id != ?", username, id).First(&existingUser).Error; err == nil {
		return nil, fmt.Errorf("用户名 %s 已被使用", username)
	}

	// 更新用户
	user.Username = username
	if err := s.db.Save(user).Error; err != nil {
		return nil, fmt.Errorf("更新用户失败: %w", err)
	}

	return user, nil
}

// DeleteUser 删除用户
func (s *userService) DeleteUser(id int) error {
	// 检查用户是否存在
	_, err := s.GetUserByID(id)
	if err != nil {
		return err
	}

	if err := s.db.Delete(&model.User{}, id).Error; err != nil {
		return fmt.Errorf("删除用户失败: %w", err)
	}

	return nil
}

// ListUsers 分页查询用户列表
func (s *userService) ListUsers(page, pageSize int) ([]model.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	var users []model.User
	var total int64

	// 获取总数
	if err := s.db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询用户总数失败: %w", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := s.db.Order("id DESC").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("查询用户列表失败: %w", err)
	}

	return users, total, nil
}
