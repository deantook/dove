package service

import (
	"context"
	"dove/internal/domain"
	"dove/internal/model"
	"dove/pkg/logger"
	"dove/pkg/pagination"
	"errors"
)

type userService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) domain.UserService {
	return &userService{repo: repo}
}

func (s *userService) Create(ctx context.Context, user *model.User) error {
	// 检查用户名是否已存在
	if _, err := s.repo.GetByUsername(ctx, user.Username); err == nil {
		logger.WarnWithTrace(ctx, "Username already exists", "username", user.Username)
		return errors.New("username already exists")
	}

	// 检查邮箱是否已存在
	if _, err := s.repo.GetByEmail(ctx, user.Email); err == nil {
		logger.WarnWithTrace(ctx, "Email already exists", "email", user.Email)
		return errors.New("email already exists")
	}

	if err := s.repo.Create(ctx, user); err != nil {
		logger.ErrorWithTrace(ctx, "Failed to create user", "error", err.Error(), "username", user.Username)
		return err
	}

	logger.InfoWithTrace(ctx, "User created successfully", "user_id", user.ID, "username", user.Username)
	return nil
}

func (s *userService) GetByID(ctx context.Context, id uint) (*model.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.ErrorWithTrace(ctx, "Failed to get user by ID", "error", err.Error(), "user_id", id)
		return nil, err
	}
	logger.InfoWithTrace(ctx, "User retrieved by ID", "user_id", id)
	return user, nil
}

func (s *userService) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		logger.ErrorWithTrace(ctx, "Failed to get user by username", "error", err.Error(), "username", username)
		return nil, err
	}
	logger.InfoWithTrace(ctx, "User retrieved by username", "username", username)
	return user, nil
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		logger.ErrorWithTrace(ctx, "Failed to get user by email", "error", err.Error(), "email", email)
		return nil, err
	}
	logger.InfoWithTrace(ctx, "User retrieved by email", "email", email)
	return user, nil
}

func (s *userService) GetAll(ctx context.Context) ([]model.User, error) {
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		logger.ErrorWithTrace(ctx, "Failed to get all users", "error", err.Error())
		return nil, err
	}
	logger.InfoWithTrace(ctx, "All users retrieved", "count", len(users))
	return users, nil
}

func (s *userService) GetAllWithPagination(ctx context.Context, page *pagination.PageRequest) (*pagination.PageResponse, error) {
	users, total, err := s.repo.GetAllWithPagination(ctx, page)
	if err != nil {
		logger.ErrorWithTrace(ctx, "Failed to get users with pagination", "error", err.Error(), "page", page.Page, "pageSize", page.PageSize)
		return nil, err
	}

	// 不返回密码
	for i := range users {
		users[i].Password = ""
	}

	pageResponse := pagination.NewPageResponse(users, total, page.Page, page.PageSize)
	logger.InfoWithTrace(ctx, "Users retrieved with pagination", "count", len(users), "total", total, "page", page.Page, "pageSize", page.PageSize)
	return pageResponse, nil
}

func (s *userService) Update(ctx context.Context, user *model.User) error {
	// 检查用户名是否已被其他用户使用
	if existingUser, err := s.repo.GetByUsername(ctx, user.Username); err == nil && existingUser.ID != user.ID {
		logger.WarnWithTrace(ctx, "Username already exists for update", "username", user.Username, "user_id", user.ID)
		return errors.New("username already exists")
	}

	// 检查邮箱是否已被其他用户使用
	if existingUser, err := s.repo.GetByEmail(ctx, user.Email); err == nil && existingUser.ID != user.ID {
		logger.WarnWithTrace(ctx, "Email already exists for update", "email", user.Email, "user_id", user.ID)
		return errors.New("email already exists")
	}

	if err := s.repo.Update(ctx, user); err != nil {
		logger.ErrorWithTrace(ctx, "Failed to update user", "error", err.Error(), "user_id", user.ID)
		return err
	}

	logger.InfoWithTrace(ctx, "User updated successfully", "user_id", user.ID, "username", user.Username)
	return nil
}

func (s *userService) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		logger.ErrorWithTrace(ctx, "Failed to delete user", "error", err.Error(), "user_id", id)
		return err
	}
	logger.InfoWithTrace(ctx, "User deleted successfully", "user_id", id)
	return nil
}
