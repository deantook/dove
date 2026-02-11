package service

import (
	"context"
	"dove/internal/domain"
	"dove/internal/model"
	"dove/pkg/logger"
	"dove/pkg/pagination"
)

type weaponService struct {
	repo domain.WeaponRepository
}

func NewWeaponService(repo domain.WeaponRepository) domain.WeaponService {
	return &weaponService{repo: repo}
}

func (s *weaponService) Create(ctx context.Context, weapon *model.Weapon) error {
	if err := s.repo.Create(ctx, weapon); err != nil {
		logger.ErrorWithTrace(ctx, "Failed to create weapon", "error", err.Error())
		return err
	}

	logger.InfoWithTrace(ctx, "Weapon created successfully", "weapon_id", weapon.ID)
	return nil
}

func (s *weaponService) GetByID(ctx context.Context, id uint) (*model.Weapon, error) {
	weapon, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.ErrorWithTrace(ctx, "Failed to get weapon by ID", "error", err.Error(), "weapon_id", id)
		return nil, err
	}
	logger.InfoWithTrace(ctx, "Weapon retrieved by ID", "weapon_id", id)
	return weapon, nil
}

func (s *weaponService) GetAll(ctx context.Context) ([]model.Weapon, error) {
	weapons, err := s.repo.GetAll(ctx)
	if err != nil {
		logger.ErrorWithTrace(ctx, "Failed to get all weapons", "error", err.Error())
		return nil, err
	}
	logger.InfoWithTrace(ctx, "All weapons retrieved", "count", len(weapons))
	return weapons, nil
}

func (s *weaponService) GetAllWithPagination(ctx context.Context, page *pagination.PageRequest) (*pagination.PageResponse, error) {
	weapons, total, err := s.repo.GetAllWithPagination(ctx, page)
	if err != nil {
		logger.ErrorWithTrace(ctx, "Failed to get weapons with pagination", "error", err.Error(), "page", page.Page, "pageSize", page.PageSize)
		return nil, err
	}

	pageResponse := pagination.NewPageResponse(weapons, total, page.Page, page.PageSize)
	logger.InfoWithTrace(ctx, "Weapons retrieved with pagination", "count", len(weapons), "total", total, "page", page.Page, "pageSize", page.PageSize)
	return pageResponse, nil
}

func (s *weaponService) Update(ctx context.Context, weapon *model.Weapon) error {
	// 检查是否存在
	_, err := s.repo.GetByID(ctx, weapon.ID)
	if err != nil {
		logger.ErrorWithTrace(ctx, "Failed to get existing weapon for update", "error", err.Error(), "weapon_id", weapon.ID)
		return err
	}

	if err := s.repo.Update(ctx, weapon); err != nil {
		logger.ErrorWithTrace(ctx, "Failed to update weapon", "error", err.Error(), "weapon_id", weapon.ID)
		return err
	}

	logger.InfoWithTrace(ctx, "Weapon updated successfully", "weapon_id", weapon.ID)
	return nil
}

func (s *weaponService) Delete(ctx context.Context, id uint) error {
	// 检查是否存在
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		logger.ErrorWithTrace(ctx, "Failed to get weapon for deletion", "error", err.Error(), "weapon_id", id)
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		logger.ErrorWithTrace(ctx, "Failed to delete weapon", "error", err.Error(), "weapon_id", id)
		return err
	}
	logger.InfoWithTrace(ctx, "Weapon deleted successfully", "weapon_id", id)
	return nil
}
