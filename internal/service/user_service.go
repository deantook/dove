package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/deantook/dove/internal/model"
	"github.com/deantook/dove/internal/repository"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// UserService 用户服务接口
type UserService interface {
	CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.UserResponse, error)
	GetUserByID(ctx context.Context, id int) (*model.UserResponse, error)
	UpdateUser(ctx context.Context, id int, req *model.UpdateUserRequest) (*model.UserResponse, error)
	DeleteUser(ctx context.Context, id int) error
	ListUsers(ctx context.Context, page, pageSize int) ([]*model.UserResponse, int64, error)
}

// userService 用户服务实现
type userService struct {
	userRepo repository.UserRepository
	redis    *redis.Client
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo repository.UserRepository, redis *redis.Client) UserService {
	return &userService{
		userRepo: userRepo,
		redis:    redis,
	}
}

// CreateUser 创建用户
func (s *userService) CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.UserResponse, error) {
	// 检查用户名是否已存在
	if _, err := s.userRepo.GetByUsername(req.Username); err == nil {
		return nil, errors.New("用户名已存在")
	} else if err != gorm.ErrRecordNotFound {
		return nil, errors.New("查询用户失败")
	}

	// 检查手机号是否已存在
	if _, err := s.userRepo.GetByPhone(req.Phone); err == nil {
		return nil, errors.New("手机号已存在")
	} else if err != gorm.ErrRecordNotFound {
		return nil, errors.New("查询用户失败")
	}

	user := &model.User{
		Username: req.Username,
		Phone:    req.Phone,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("创建用户失败")
	}

	return user.ToResponse(), nil
}

// GetUserByID 根据 ID 获取用户
func (s *userService) GetUserByID(ctx context.Context, id int) (*model.UserResponse, error) {
	// 尝试从缓存获取
	cacheKey := fmt.Sprintf("user:%d", id)
	if s.redis != nil {
		// 这里可以添加缓存逻辑
		_ = cacheKey
	}

	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("用户不存在")
		}
		return nil, errors.New("查询用户失败")
	}

	return user.ToResponse(), nil
}

// UpdateUser 更新用户
func (s *userService) UpdateUser(ctx context.Context, id int, req *model.UpdateUserRequest) (*model.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, errors.New("用户不存在")
		}
		return nil, errors.New("更新用户失败")
	}

	// 如果更新用户名，检查是否重复
	if req.Username != "" && req.Username != user.Username {
		if _, err := s.userRepo.GetByUsername(req.Username); err == nil {
			return nil, errors.New("用户名已存在")
		} else if !errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, errors.New("查询用户失败")
		}
		user.Username = req.Username
	}

	// 如果更新手机号，检查是否重复
	if req.Phone != "" && req.Phone != user.Phone {
		if _, err := s.userRepo.GetByPhone(req.Phone); err == nil {
			return nil, errors.New("手机号已存在")
		} else if !errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, errors.New("查询用户失败")
		}
		user.Phone = req.Phone
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("更新用户失败")
	}

	// 清除缓存
	if s.redis != nil {
		cacheKey := fmt.Sprintf("user:%d", id)
		s.redis.Del(ctx, cacheKey)
	}

	return user.ToResponse(), nil
}

// DeleteUser 删除用户
func (s *userService) DeleteUser(ctx context.Context, id int) error {
	// 检查用户是否存在
	if _, err := s.userRepo.GetByID(id); err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return errors.New("用户不存在")
		}
		return errors.New("查询用户失败")
	}

	if err := s.userRepo.Delete(id); err != nil {
		return errors.New("删除用户失败")
	}

	// 清除缓存
	if s.redis != nil {
		cacheKey := fmt.Sprintf("user:%d", id)
		s.redis.Del(ctx, cacheKey)
	}

	return nil
}

// ListUsers 获取用户列表（分页）
func (s *userService) ListUsers(ctx context.Context, page, pageSize int) ([]*model.UserResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize
	users, total, err := s.userRepo.List(offset, pageSize)
	if err != nil {
		return nil, 0, errors.New("查询用户列表失败")
	}

	responses := make([]*model.UserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, user.ToResponse())
	}

	return responses, total, nil
}
