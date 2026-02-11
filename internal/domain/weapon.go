package domain

import (
	"context"
	"dove/internal/model"
	"dove/pkg/pagination"
)

// WeaponRepository Weapon 仓库接口
type WeaponRepository interface {
	Create(ctx context.Context, weapon *model.Weapon) error
	GetByID(ctx context.Context, id uint) (*model.Weapon, error)
	GetAll(ctx context.Context) ([]model.Weapon, error)
	GetAllWithPagination(ctx context.Context, page *pagination.PageRequest) ([]model.Weapon, int64, error)
	Update(ctx context.Context, weapon *model.Weapon) error
	Delete(ctx context.Context, id uint) error
}

// WeaponService Weapon 服务接口
type WeaponService interface {
	Create(ctx context.Context, weapon *model.Weapon) error
	GetByID(ctx context.Context, id uint) (*model.Weapon, error)
	GetAll(ctx context.Context) ([]model.Weapon, error)
	GetAllWithPagination(ctx context.Context, page *pagination.PageRequest) (*pagination.PageResponse, error)
	Update(ctx context.Context, weapon *model.Weapon) error
	Delete(ctx context.Context, id uint) error
}

// CreateWeaponRequest 创建 Weapon 请求
type CreateWeaponRequest struct {
	Name string `json:"name"`

	Level int `json:"level"`

	Content string `json:"content"`

	Type int `json:"type"`

	Story string `json:"story"`
}

// UpdateWeaponRequest 更新 Weapon 请求
type UpdateWeaponRequest struct {
	Name *string `json:"name"`

	Level *int `json:"level"`

	Content *string `json:"content"`

	Type *int `json:"type"`

	Story *string `json:"story"`
}
