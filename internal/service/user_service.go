package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/deantook/dove/internal/model"
	"github.com/deantook/dove/internal/repository"
	"github.com/deantook/dove/pkg/jwt"
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
	SendCode(ctx context.Context, req *model.SendCodeRequest) (*model.SendCodeResponse, error)
	LoginOrRegister(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error)
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

	now := time.Now()
	user := &model.User{
		Username:   req.Username,
		Phone:      req.Phone,
		CreateTime: now,
		UpdateTime: now,
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

// SendCode 发送验证码
func (s *userService) SendCode(ctx context.Context, req *model.SendCodeRequest) (*model.SendCodeResponse, error) {
	// 生成6位随机验证码
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := fmt.Sprintf("%06d", r.Intn(1000000))

	// 将验证码存储到 Redis，有效期5分钟
	if s.redis != nil {
		codeKey := fmt.Sprintf("sms:code:%s", req.Phone)
		if err := s.redis.Set(ctx, codeKey, code, 5*time.Minute).Err(); err != nil {
			return nil, errors.New("存储验证码失败")
		}
	}

	// 开发阶段返回验证码，生产环境不返回
	return &model.SendCodeResponse{
		Code: code,
	}, nil
}

// LoginOrRegister 登录或注册
func (s *userService) LoginOrRegister(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
	// 验证验证码
	if s.redis != nil {
		codeKey := fmt.Sprintf("sms:code:%s", req.Phone)
		storedCode, err := s.redis.Get(ctx, codeKey).Result()
		if err == redis.Nil {
			return nil, errors.New("验证码已过期或不存在")
		} else if err != nil {
			return nil, errors.New("验证验证码失败")
		}

		if storedCode != req.Code {
			return nil, errors.New("验证码错误")
		}

		// 验证成功后删除验证码
		s.redis.Del(ctx, codeKey)
	}

	// 查找用户是否存在
	user, err := s.userRepo.GetByPhone(req.Phone)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.New("查询用户失败")
	}

	// 如果用户不存在，创建新用户
	if err == gorm.ErrRecordNotFound {
		now := time.Now()
		user = &model.User{
			Phone:      req.Phone,
			Username:   req.Phone, // 默认用户名为手机号
			CreateTime: now,
			UpdateTime: now,
		}
		if err := s.userRepo.Create(user); err != nil {
			return nil, errors.New("创建用户失败")
		}
	}

	// 生成 JWT token
	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		return nil, errors.New("生成token失败")
	}

	return &model.LoginResponse{
		User:  user.ToResponse(),
		Token: token,
	}, nil
}
