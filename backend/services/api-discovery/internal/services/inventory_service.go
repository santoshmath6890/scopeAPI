package services

import (
	"context"
	"fmt"
	"time"

	"scopeapi.local/backend/services/api-discovery/internal/models"
	"scopeapi.local/backend/services/api-discovery/internal/repository"
	"shared/logging"
)

type InventoryServiceInterface interface {
	GetAPIInventory(ctx context.Context, page, limit int, filters InventoryFilter) (*models.APIInventory, error)
	GetAPIDetails(ctx context.Context, apiID string) (*models.APIDetails, error)
	UpdateAPITags(ctx context.Context, apiID string, tags []string) error
	DeleteAPI(ctx context.Context, apiID string) error
	GetAPIStatistics(ctx context.Context) (*models.APIStatistics, error)
}

type InventoryService struct {
	repo   repository.InventoryRepositoryInterface
	logger logging.Logger
}

type InventoryFilter struct {
	Status   string `json:"status" form:"status"`
	Protocol string `json:"protocol" form:"protocol"`
	Domain   string `json:"domain" form:"domain"`
	Tags     string `json:"tags" form:"tags"`
	DateFrom string `json:"date_from" form:"date_from"`
	DateTo   string `json:"date_to" form:"date_to"`
}

func NewInventoryService(repo repository.InventoryRepositoryInterface, logger logging.Logger) InventoryServiceInterface {
	return &InventoryService{
		repo:   repo,
		logger: logger,
	}
}

func (s *InventoryService) GetAPIInventory(ctx context.Context, page, limit int, filters InventoryFilter) (*models.APIInventory, error) {
	// For now, return basic inventory without filtering
	inventory, err := s.repo.GetAPIs(ctx, page, limit)
	if err != nil {
		s.logger.Error("Failed to get API inventory", "error", err)
		return nil, fmt.Errorf("failed to get API inventory: %w", err)
	}

	return inventory, nil
}

func (s *InventoryService) GetAPIDetails(ctx context.Context, apiID string) (*models.APIDetails, error) {
	details, err := s.repo.GetAPIDetails(ctx, apiID)
	if err != nil {
		s.logger.Error("Failed to get API details", "error", err, "api_id", apiID)
		return nil, fmt.Errorf("failed to get API details: %w", err)
	}

	return details, nil
}

func (s *InventoryService) UpdateAPITags(ctx context.Context, apiID string, tags []string) error {
	// Get current API
	api, err := s.repo.GetAPI(ctx, apiID)
	if err != nil {
		s.logger.Error("Failed to get API for tag update", "error", err, "api_id", apiID)
		return fmt.Errorf("failed to get API: %w", err)
	}

	// Update tags
	api.Tags = tags
	api.UpdatedAt = time.Now()

	err = s.repo.UpdateAPI(ctx, api)
	if err != nil {
		s.logger.Error("Failed to update API tags", "error", err, "api_id", apiID)
		return fmt.Errorf("failed to update API tags: %w", err)
	}

	s.logger.Info("API tags updated", "api_id", apiID, "tags", tags)
	return nil
}

func (s *InventoryService) DeleteAPI(ctx context.Context, apiID string) error {
	err := s.repo.DeleteAPI(ctx, apiID)
	if err != nil {
		s.logger.Error("Failed to delete API", "error", err, "api_id", apiID)
		return fmt.Errorf("failed to delete API: %w", err)
	}

	s.logger.Info("API deleted", "api_id", apiID)
	return nil
}

func (s *InventoryService) GetAPIStatistics(ctx context.Context) (*models.APIStatistics, error) {
	stats, err := s.repo.GetAPIStatistics(ctx)
	if err != nil {
		s.logger.Error("Failed to get API statistics", "error", err)
		return nil, fmt.Errorf("failed to get API statistics: %w", err)
	}

	return stats, nil
}

// Removed all MetadataService-related functions from this file. Only inventory-specific logic remains.
