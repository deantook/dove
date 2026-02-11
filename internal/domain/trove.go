package domain

import (
	"context"
	"dove/internal/model"
	"dove/pkg/pagination"
)

// TroveRepository Trove 仓库接口
type TroveRepository interface {
	Create(ctx context.Context, trove *model.Trove) error
	GetByID(ctx context.Context, id uint) (*model.Trove, error)
	GetAll(ctx context.Context) ([]model.Trove, error)
	GetAllWithPagination(ctx context.Context, page *pagination.PageRequest) ([]model.Trove, int64, error)
	Update(ctx context.Context, trove *model.Trove) error
	Delete(ctx context.Context, id uint) error
}

// TroveService Trove 服务接口
type TroveService interface {
	Create(ctx context.Context, trove *model.Trove) error
	GetByID(ctx context.Context, id uint) (*model.Trove, error)
	GetAll(ctx context.Context) ([]model.Trove, error)
	GetAllWithPagination(ctx context.Context, page *pagination.PageRequest) (*pagination.PageResponse, error)
	Update(ctx context.Context, trove *model.Trove) error
	Delete(ctx context.Context, id uint) error
}

// CreateTroveRequest 创建 Trove 请求
type CreateTroveRequest struct {
	Title string `json:"title"`

	Description string `json:"description"`
}

// UpdateTroveRequest 更新 Trove 请求
type UpdateTroveRequest struct {
	Title *string `json:"title"`

	Description *string `json:"description"`
}
