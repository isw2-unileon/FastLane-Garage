// Package service contains tests for the parts service business logic.
package service

import (
	"testing"

	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/dto"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/models"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/repository"
)

// mockPartsRepository is a mock implementation of PartsRepository for testing.
type mockPartsRepository struct {
	// parts is an in-memory store of parts for testing.
	parts map[uint]*models.Part
	// nextID tracks the next ID to assign to new parts.
	nextID uint
}

// NewMockPartsRepository creates a new mock repository with initial test data.
func NewMockPartsRepository() *mockPartsRepository {
	return &mockPartsRepository{
		parts:  make(map[uint]*models.Part),
		nextID: 1,
	}
}

// FindAll returns all parts from the mock repository.
func (m *mockPartsRepository) FindAll(zone string, name string) ([]models.Part, error) {
	parts := make([]models.Part, 0, len(m.parts))
	for _, p := range m.parts {
		parts = append(parts, *p)
	}
	return parts, nil
}

// FindByID returns a single part by ID from the mock repository.
func (m *mockPartsRepository) FindByID(id uint) (*models.Part, error) {
	part, exists := m.parts[id]
	if !exists {
		return nil, repository.ErrNotFound
	}
	return part, nil
}

// Create inserts a new part into the mock repository.
func (m *mockPartsRepository) Create(part *models.Part) (*models.Part, error) {
	part.ID = m.nextID
	m.parts[m.nextID] = part
	m.nextID++
	return part, nil
}

// Update modifies an existing part in the mock repository.
func (m *mockPartsRepository) Update(id uint, updates map[string]interface{}) (*models.Part, error) {
	part, exists := m.parts[id]
	if !exists {
		return nil, repository.ErrNotFound
	}

	// Apply updates to the part.
	if name, ok := updates["name"]; ok {
		part.Name = name.(string)
	}
	if carZone, ok := updates["car_zone"]; ok {
		part.CarZone = carZone.(string)
	}
	if imageURL, ok := updates["image_url"]; ok {
		part.ImageURL = imageURL.(string)
	}
	if price, ok := updates["price"]; ok {
		part.Price = price.(float64)
	}

	return part, nil
}

// Delete removes a part from the mock repository.
func (m *mockPartsRepository) Delete(id uint) error {
	if _, exists := m.parts[id]; !exists {
		return repository.ErrNotFound
	}
	delete(m.parts, id)
	return nil
}

// TestGetAllParts tests the GetAllParts service method.
func TestGetAllParts(t *testing.T) {
	// Arrange: Set up the mock repository with test data.
	mockRepo := NewMockPartsRepository()
	_, err := mockRepo.Create(&models.Part{
		Name:    "Motor V6",
		CarZone: "motor",
		Price:   2500.00,
	})
	if err != nil {
		t.Fatalf("failed to create part: %v", err)
	}
	_, err = mockRepo.Create(&models.Part{
		Name:    "Turbo",
		CarZone: "motor",
		Price:   1200.00,
	})
	if err != nil {
		t.Fatalf("failed to create part: %v", err)
	}

	// Create the servie with the mock repository.
	svc := NewPartsService(mockRepo)

	// Act: Call the service method.
	parts, err := svc.GetAllParts("", "")
	// Assert: Verify the results.
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(parts) != 2 {
		t.Fatalf("expected 2 parts, got %d", len(parts))
	}
	if parts[0].Name != "Motor V6" {
		t.Fatalf("expected first part name to be 'Motor V6', got %s", parts[0].Name)
	}
}

// TestGetPartByID tests the GetPartByID service method.
func TestGetPartByID(t *testing.T) {
	// Arrange: Set up the mock repository.
	mockRepo := NewMockPartsRepository()
	created, err := mockRepo.Create(&models.Part{
		Name:    "Motor V6",
		CarZone: "motor",
		Price:   2500.00,
	})
	if err != nil {
		t.Fatalf("failed to create part: %v", err)
	}

	svc := NewPartsService(mockRepo)

	// Act: Call the service method with the created parts's ID.
	part, err := svc.GetPartByID(created.ID)
	// Assert: Verify the results.
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if part.Name != "Motor V6" {
		t.Fatalf("expected part name to be 'Motor V6', got %s", part.Name)
	}
}

// TestCreatePart tests the CreatePart service method.
func TestCreatePart(t *testing.T) {
	// Arrange: Set up the mock repository and service.
	mockRepo := NewMockPartsRepository()
	svc := NewPartsService(mockRepo)

	req := &dto.CreatePartRequest{
		Name:    "Disco de freno",
		CarZone: "frenos",
		Price:   120.00,
	}

	// Act: Call the service method.
	part, err := svc.CreatePart(req)
	// Assert: Verify the results.
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if part.Name != "Disco de freno" {
		t.Fatalf("expected part name to be 'Disco de freno', got %s", part.Name)
	}
	if part.Price != 120.00 {
		t.Fatalf("expected price to be 120.00, got %f", part.Price)
	}
}

// TestCreatePartWithInvalidPrice tests that CreatePart rejects invalid prices.
func TestCreatePartWithInvalidPrice(t *testing.T) {
	// Arrange: Set up the mock repository and service.
	mockRepo := NewMockPartsRepository()
	svc := NewPartsService(mockRepo)

	req := &dto.CreatePartRequest{
		Name:    "Invalid Part",
		CarZone: "motor",
		Price:   -100, // Invalid: negative price
	}

	// Act: Call the service method.
	part, err := svc.CreatePart(req)

	// Assert: Verify that an error was returned.
	if err == nil {
		t.Fatal("expected an error for negative price, but got none")
	}
	if part != nil {
		t.Fatal("expected no part to be created for invalid price")
	}
}

// TestUpdatePart tests the UpdatePart service method.
func TestUpdatePart(t *testing.T) {
	// Arrange: Set up the mock repository with a part.
	mockRepo := NewMockPartsRepository()
	created, _ := mockRepo.Create(&models.Part{
		Name:    "Motor V6",
		CarZone: "motor",
		Price:   2500.00,
	})

	svc := NewPartsService(mockRepo)

	newPrice := 2800.00
	req := &dto.UpdatePartRequest{
		Name:  "Motor V8",
		Price: &newPrice,
	}

	// Act: Call the service method.
	part, err := svc.UpdatePart(created.ID, req)
	// Assert: Verify the results.
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if part.Name != "Motor V8" {
		t.Fatalf("expected updated name to be 'Motor V8', got %s", part.Name)
	}
	if part.Price != 2800.00 {
		t.Fatalf("expected updated price to be 2800.00, got %f", part.Price)
	}
}

// TestDeletePart tests the DeletePart service method.
func TestDeletePart(t *testing.T) {
	// Arrange: Set up the mock repository with a part.
	mockRepo := NewMockPartsRepository()
	created, _ := mockRepo.Create(&models.Part{
		Name:    "Motor V6",
		CarZone: "motor",
		Price:   2500.00,
	})

	svc := NewPartsService(mockRepo)

	// Act: Call the service method to delete the part.
	err := svc.DeletePart(created.ID)
	// Assert: Verify no error occurred.
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify the part was actually deleted.
	_, err = svc.GetPartByID(created.ID)
	if err == nil {
		t.Fatal("expected an error when retrieving deleted part")
	}
}
