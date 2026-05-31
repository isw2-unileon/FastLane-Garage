// Package service contains business logic for parts operations.
package service

import (
	"fmt"

	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/dto"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/models"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/repository"
)

// PartsService defines the interface for part business operations.
// This abstraction allows us to mock the service in tests.
type PartsService interface {
	// GetAllParts retrieves all parts from the repository.
	GetAllParts() ([]dto.PartResponse, error)

	// GetPartByID retrieves a single part by its ID.
	GetPartByID(id uint) (*dto.PartResponse, error)

	// CreatePart validates the request and creates a new part.
	CreatePart(req *dto.CreatePartRequest) (*dto.PartResponse, error)

	// UpdatePart validates the request and updates an existing part.
	UpdatePart(id uint, req *dto.UpdatePartRequest) (*dto.PartResponse, error)

	// DeletePart removes a part by ID.
	DeletePart(id uint) error
}

// partsServiceImpl is the concrete implementation of PartsService.
type partsServiceImpl struct {
	repo repository.PartsRepository
}

// NewPartsService creates a new instance of the parts service.
// It accepts a PartsRepository and returns a PartsService interface.
func NewPartsService(repo repository.PartsRepository) PartsService {
	return &partsServiceImpl{
		repo: repo,
	}
}

// GetAllParts retrieves all parts and converts them to DTOs.
// DTOs are used for responses to ensure a consistent API contract.
func (s *partsServiceImpl) GetAllParts() ([]dto.PartResponse, error) {
	// Fetch all parts from the repository.
	parts, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all parts: %w", err)
	}

	// Convert Part models to PartResponse DTOs.
	responses := make([]dto.PartResponse, len(parts))
	for i, part := range parts {
		responses[i] = s.partToResponse(&part)
	}

	return responses, nil
}

// GetPartByID retrieves a single part by ID and converts it to a DTO.
func (s *partsServiceImpl) GetPartByID(id uint) (*dto.PartResponse, error) {
	// Fetch the part from the repository.
	part, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get part with ID %d: %w", id, err)
	}

	// Convert to DTO and return.
	response := s.partToResponse(part)
	return &response, nil
}

// CreatePart validates the request data and creates a new part.
// Business logic can be added here (e.g., check for duplicates, apply discounts, etc.).
func (s *partsServiceImpl) CreatePart(req *dto.CreatePartRequest) (*dto.PartResponse, error) {
	// Validate request (this is already done by Gin binding,
	// but we can add additional business logic here).
	if req.Name == "" || req.CarZone == "" {
		return nil, fmt.Errorf("part name and car zone are required")
	}

	if req.Price <= 0 {
		return nil, fmt.Errorf("price must be greater than zero")
	}

	// Convert DTO to model.
	part := &models.Part{
		Name:     req.Name,
		CarZone:  req.CarZone,
		ImageURL: req.ImageURL,
		Price:    req.Price,
	}

	// Create the part in the repository.
	createdPart, err := s.repo.Create(part)
	if err != nil {
		return nil, fmt.Errorf("failed to create part: %w", err)
	}

	// Convert to DTO and return.
	response := s.partToResponse(createdPart)
	return &response, nil
}

// UpdatePart validates the request and updates an existing part.
// Only provided fields are updated (partial updates).
func (s *partsServiceImpl) UpdatePart(id uint, req *dto.UpdatePartRequest) (*dto.PartResponse, error) {
	// First, verify the part exists.
	if _, err := s.repo.FindByID(id); err != nil {
		return nil, err
	}

	// Build update map with only non-empty fields.
	updates := make(map[string]interface{})

	// Only add fields that were provided in the request.
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.CarZone != "" {
		updates["car_zone"] = req.CarZone
	}
	if req.ImageURL != "" {
		updates["image_url"] = req.ImageURL
	}
	if req.Price != nil && *req.Price > 0 {
		updates["price"] = *req.Price
	}

	// If no fields were provided, return an error.
	if len(updates) == 0 {
		return nil, fmt.Errorf("no fields to update provided")
	}

	// Update the part in the repository.
	updatedPart, err := s.repo.Update(id, updates)
	if err != nil {
		return nil, fmt.Errorf("failed to update part: %w", err)
	}

	// Convert to DTO and return.
	response := s.partToResponse(updatedPart)
	return &response, nil
}

// DeletePart removes a part by ID.
func (s *partsServiceImpl) DeletePart(id uint) error {
	// Verify the part exists before deleting.
	if _, err := s.repo.FindByID(id); err != nil {
		return err
	}

	// Delete the part from the repository.
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete part: %w", err)
	}

	return nil
}

// partToResponse converts a Part model to a PartResponse DTO.
// This method ensures consistent conversion logic throughout the service.
func (s *partsServiceImpl) partToResponse(part *models.Part) dto.PartResponse {
	return dto.PartResponse{
		ID:       part.ID,
		Name:     part.Name,
		CarZone:  part.CarZone,
		ImageURL: part.ImageURL,
		Price:    part.Price,
	}
}
