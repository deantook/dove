package domain

import (
	"context"
	"dove/internal/model"
	"dove/pkg/pagination"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uint) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetAll(ctx context.Context) ([]model.User, error)
	GetAllWithPagination(ctx context.Context, page *pagination.PageRequest) ([]model.User, int64, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint) error
}

type UserService interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uint) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetAll(ctx context.Context) ([]model.User, error)
	GetAllWithPagination(ctx context.Context, page *pagination.PageRequest) (*pagination.PageResponse, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint) error
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required" example:"john_doe"`
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"123456"`
	Nickname string `json:"nickname" example:"John Doe"`
	Avatar   string `json:"avatar" example:"https://example.com/avatar.jpg"`
	Status   int    `json:"status" example:"1"`
}

type UpdateUserRequest struct {
	Username string `json:"username" example:"john_doe"`
	Email    string `json:"email" binding:"omitempty,email" example:"john@example.com"`
	Password string `json:"password" binding:"omitempty,min=6" example:"123456"`
	Nickname string `json:"nickname" example:"John Doe"`
	Avatar   string `json:"avatar" example:"https://example.com/avatar.jpg"`
	Status   *int   `json:"status" example:"1"`
}
