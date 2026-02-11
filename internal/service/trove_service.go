package service

import (
	"context"
	"dove/internal/domain"
	"dove/internal/model"
	"dove/pkg/logger"
	"dove/pkg/pagination"
)

type troveService struct {
	repo domain.TroveRepository
}

func NewTroveService(repo domain.TroveRepository) domain.TroveService {
	return &troveService{repo: repo}
}

func (s *troveService) Create(ctx context.Context, trove *model.Trove) error {
	if err := s.repo.Create(ctx, trove); err != nil {
		logger.ErrorWithTrace(ctx, "Failed to create trove", "error", err.Error())
		return err
	}

	logger.InfoWithTrace(ctx, "Trove created successfully", "trove_id", trove.ID)
	return nil
}

func (s *troveService) GetByID(ctx context.Context, id uint) (*model.Trove, error) {
	trove, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.ErrorWithTrace(ctx, "Failed to get trove by ID", "error", err.Error(), "trove_id", id)
		return nil, err
	}
	logger.InfoWithTrace(ctx, "Trove retrieved by ID", "trove_id", id)
	return trove, nil
}

func (s *troveService) GetAll(ctx context.Context) ([]model.Trove, error) {
	troves, err := s.repo.GetAll(ctx)
	if err != nil {
		logger.ErrorWithTrace(ctx, "Failed to get all troves", "error", err.Error())
		return nil, err
	}
	logger.InfoWithTrace(ctx, "All troves retrieved", "count", len(troves))
	return troves, nil
}

func (s *troveService) GetAllWithPagination(ctx context.Context, page *pagination.PageRequest) (*pagination.PageResponse, error) {
	troves, total, err := s.repo.GetAllWithPagination(ctx, page)
	if err != nil {
		logger.ErrorWithTrace(ctx, "Failed to get troves with pagination", "error", err.Error(), "page", page.Page, "pageSize", page.PageSize)
		return nil, err
	}

	pageResponse := pagination.NewPageResponse(troves, total, page.Page, page.PageSize)
	logger.InfoWithTrace(ctx, "Troves retrieved with pagination", "count", len(troves), "total", total, "page", page.Page, "pageSize", page.PageSize)
	return pageResponse, nil
}

func (s *troveService) Update(ctx context.Context, trove *model.Trove) error {
	// 检查是否存在
	_, err := s.repo.GetByID(ctx, trove.ID)
	if err != nil {
		logger.ErrorWithTrace(ctx, "Failed to get existing trove for update", "error", err.Error(), "trove_id", trove.ID)
		return err
	}

	if err := s.repo.Update(ctx, trove); err != nil {
		logger.ErrorWithTrace(ctx, "Failed to update trove", "error", err.Error(), "trove_id", trove.ID)
		return err
	}

	logger.InfoWithTrace(ctx, "Trove updated successfully", "trove_id", trove.ID)
	return nil
}

func (s *troveService) Delete(ctx context.Context, id uint) error {
	// 检查是否存在
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		logger.ErrorWithTrace(ctx, "Failed to get trove for deletion", "error", err.Error(), "trove_id", id)
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		logger.ErrorWithTrace(ctx, "Failed to delete trove", "error", err.Error(), "trove_id", id)
		return err
	}
	logger.InfoWithTrace(ctx, "Trove deleted successfully", "trove_id", id)
	return nil
}
