// Package service contains tests for the orders service business logic.
package service

import (
	"fmt"
	"testing"

	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/dto"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/models"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/repository"
)

// mockOrdersRepository is a mock implementation of OrdersRepository for testing.
type mockOrdersRepository struct {
	// orders is an in-memory store of orders for testing.
	orders map[uint]*models.Order
	// nextID tracks the next ID to assign to new orders.
	nextID uint

	// MockGetTopParts allows inyect custom behavior for GetTopParts in tests.
	MockGetTopParts func(limit int) ([]dto.TopPartResponse, error)
}

// NewMockOrdersRepository creates a new mock repository with initial test data.
func NewMockOrdersRepository() *mockOrdersRepository {
	return &mockOrdersRepository{
		orders: make(map[uint]*models.Order),
		nextID: 1,
	}
}

// FindAll returns all orders from the mock repository.
func (m *mockOrdersRepository) FindAll(status string, email string) ([]models.Order, error) {
	orders := make([]models.Order, 0, len(m.orders))
	for _, o := range m.orders {
		orders = append(orders, *o)
	}
	return orders, nil
}

// FindByID returns a single order by ID from the mock repository.
func (m *mockOrdersRepository) FindByID(id uint) (*models.Order, error) {
	order, exists := m.orders[id]
	if !exists {
		return nil, repository.ErrNotFound
	}
	return order, nil
}

// Create inserts a new order into the mock repository.
func (m *mockOrdersRepository) Create(order *models.Order) (*models.Order, error) {
	order.ID = m.nextID
	m.orders[m.nextID] = order
	m.nextID++
	return order, nil
}

// UpdateStatus changes the status of an order in the mock repository.
func (m *mockOrdersRepository) UpdateStatus(id uint, status models.OrderStatus) (*models.Order, error) {
	order, exists := m.orders[id]
	if !exists {
		return nil, repository.ErrNotFound
	}
	order.Status = status
	return order, nil
}

// Delete removes an order from the mock repository.
func (m *mockOrdersRepository) Delete(id uint) error {
	if _, exists := m.orders[id]; !exists {
		return repository.ErrNotFound
	}
	delete(m.orders, id)
	return nil
}

// mockPartsRepository is a mock implementation of PartsRepository for testing.
func (m *mockOrdersRepository) GetTopParts(limit int) ([]dto.TopPartResponse, error) {
	if m.MockGetTopParts != nil {
		return m.MockGetTopParts(limit)
	}
	return nil, nil
}

// TestGetAllOrders tests the GetAllOrders service method.
func TestGetAllOrders(t *testing.T) {
	// Arrange: Set up the mock repository with test data.
	mockRepo := NewMockOrdersRepository()
	_, err := mockRepo.Create(&models.Order{
		CustomerName:  "Juan García",
		CustomerEmail: "juan@example.com",
		Status:        models.OrderStatusPending,
		TotalPrice:    2500.00,
	})
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}
	_, err = mockRepo.Create(&models.Order{
		CustomerName:  "María López",
		CustomerEmail: "maria@example.com",
		Status:        models.OrderStatusCompleted,
		TotalPrice:    1200.00,
	})
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	// Create the service with the mock repository.
	svc := NewOrdersService(mockRepo)

	// Act: Call the service method.
	orders, err := svc.GetAllOrders("", "")
	// Assert: Verify the results.
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(orders) != 2 {
		t.Fatalf("expected 2 orders, got %d", len(orders))
	}
	if orders[0].CustomerName != "Juan García" {
		t.Fatalf("expected first order customer to be 'Juan García', got %s", orders[0].CustomerName)
	}
}

// TestGetAllOrdersWithStatusFilter tests filtering orders by status.
func TestGetAllOrdersWithStatusFilter(t *testing.T) {
	// Arrange: Set up the mock repository with multiple orders.
	mockRepo := NewMockOrdersRepository()
	_, err := mockRepo.Create(&models.Order{
		CustomerName:  "Juan García",
		CustomerEmail: "juan@example.com",
		Status:        models.OrderStatusPending,
		TotalPrice:    2500.00,
	})
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}
	_, err = mockRepo.Create(&models.Order{
		CustomerName:  "Maria Lopez",
		CustomerEmail: "maria@example.com",
		Status:        models.OrderStatusCompleted,
		TotalPrice:    1200.00,
	})
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	svc := NewOrdersService(mockRepo)

	// Act: Call the service with status filter.
	orders, err := svc.GetAllOrders("pending", "")
	// Assert: Verify only pending orders are returned.
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(orders) != 1 {
		t.Fatalf("expected 1 pending order, got %d", len(orders))
	}
	if orders[0].Status != "pending" {
		t.Fatalf("expected status to be 'pending', got %s", orders[0].Status)
	}
}

// TestGetOrderByID tests the GetOrderByID service method.
func TestGetOrderByID(t *testing.T) {
	// Arrange: Set up the mock repository.
	mockRepo := NewMockOrdersRepository()
	created, err := mockRepo.Create(&models.Order{
		CustomerName:  "Juan García",
		CustomerEmail: "juan@example.com",
		Status:        models.OrderStatusPending,
		TotalPrice:    2500.00,
	})
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	svc := NewOrdersService(mockRepo)

	// Act: Call the service method with the created order's ID.
	order, err := svc.GetOrderByID(created.ID)
	// Assert: Verify the results.
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if order.CustomerName != "Juan García" {
		t.Fatalf("expected customer name to be 'Juan García', got %s", order.CustomerName)
	}
}

// TestCreateOrder tests the CreateOrder service method.
func TestCreateOrder(t *testing.T) {
	// Arrange: Set up the mock repositories.
	mockOrdersRepo := NewMockOrdersRepository()
	mockPartsRepo := NewMockPartsRepository()

	// Create a part to include in the order.
	part, err := mockPartsRepo.Create(&models.Part{
		Name:    "Motor V6",
		CarZone: "motor",
		Price:   2500.00,
	})
	if err != nil {
		t.Fatalf("failed to create part: %v", err)
	}

	svc := NewOrdersService(mockOrdersRepo)

	req := &dto.CreateOrderRequest{
		CustomerName:  "Juan García",
		CustomerEmail: "juan@example.com",
		Items: []dto.OrderItemRequest{
			{
				PartID:   part.ID,
				Quantity: 1,
			},
		},
	}

	// Act: Call the service method.
	order, err := svc.CreateOrder(req, mockPartsRepo)
	// Assert: Verify the results.
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if order.CustomerName != "Juan García" {
		t.Fatalf("expected customer name to be 'Juan García', got %s", order.CustomerName)
	}
	if order.TotalPrice != 2500.00 {
		t.Fatalf("expected total price to be 2500.00, got %f", order.TotalPrice)
	}
}

// TestCreateOrderWithInvalidEmail tests that CreateOrder rejects invalid emails.
func TestCreateOrderWithEmptyEmail(t *testing.T) {
	// Arrange: Set up the mock repositories.
	mockOrdersRepo := NewMockOrdersRepository()
	mockPartsRepo := NewMockPartsRepository()

	svc := NewOrdersService(mockOrdersRepo)

	req := &dto.CreateOrderRequest{
		CustomerName:  "Juan García",
		CustomerEmail: "", // Empty email
		Items: []dto.OrderItemRequest{
			{
				PartID:   1,
				Quantity: 1,
			},
		},
	}

	// Act: Call the service method.
	order, err := svc.CreateOrder(req, mockPartsRepo)

	// Assert: Verify that an error was returned.
	if err == nil {
		t.Fatal("expected an error for empty email, but got none")
	}
	if order != nil {
		t.Fatal("expected no order to be created for empty email")
	}
}

// TestCreateOrderWithNonexistentPart tests that CreateOrder rejects non-existent parts.
func TestCreateOrderWithNonexistentPart(t *testing.T) {
	// Arrange: Set up the mock repositories.
	mockOrdersRepo := NewMockOrdersRepository()
	mockPartsRepo := NewMockPartsRepository()

	svc := NewOrdersService(mockOrdersRepo)

	req := &dto.CreateOrderRequest{
		CustomerName:  "Juan García",
		CustomerEmail: "juan@example.com",
		Items: []dto.OrderItemRequest{
			{
				PartID:   999, // Non-existent part ID
				Quantity: 1,
			},
		},
	}

	// Act: Call the service method.
	order, err := svc.CreateOrder(req, mockPartsRepo)

	// Assert: Verify that an error was returned.
	if err == nil {
		t.Fatal("expected an error for non-existent part, but got none")
	}
	if order != nil {
		t.Fatal("expected no order to be created for non-existent part")
	}
}

// TestUpdateOrderStatus tests the UpdateOrderStatus service method.
func TestUpdateOrderStatus(t *testing.T) {
	// Arrange: Set up the mock repository with an order.
	mockRepo := NewMockOrdersRepository()
	created, _ := mockRepo.Create(&models.Order{
		CustomerName:  "Juan García",
		CustomerEmail: "juan@example.com",
		Status:        models.OrderStatusPending,
		TotalPrice:    2500.00,
	})

	svc := NewOrdersService(mockRepo)

	req := &dto.UpdateOrderStatusRequest{
		Status: "processing",
	}

	// Act: Call the service method.
	order, err := svc.UpdateOrderStatus(created.ID, req)
	// Assert: Verify the results.
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if order.Status != "processing" {
		t.Fatalf("expected updated status to be 'processing', got %s", order.Status)
	}
}

// TestUpdateOrderStatusToCompleted tests changing order status to completed.
func TestUpdateOrderStatusToCompleted(t *testing.T) {
	// Arrange: Set up the mock repository with an order.
	mockRepo := NewMockOrdersRepository()
	created, _ := mockRepo.Create(&models.Order{
		CustomerName:  "Juan García",
		CustomerEmail: "juan@example.com",
		Status:        models.OrderStatusProcessing,
		TotalPrice:    2500.00,
	})

	svc := NewOrdersService(mockRepo)

	req := &dto.UpdateOrderStatusRequest{
		Status: "completed",
	}

	// Act: Call the service method.
	order, err := svc.UpdateOrderStatus(created.ID, req)
	// Assert: Verify the results.
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if order.Status != "completed" {
		t.Fatalf("expected updated status to be 'completed', got %s", order.Status)
	}
}

// TestDeleteOrder tests the DeleteOrder service method.
func TestDeleteOrder(t *testing.T) {
	// Arrange: Set up the mock repository with an order.
	mockRepo := NewMockOrdersRepository()
	created, _ := mockRepo.Create(&models.Order{
		CustomerName:  "Juan García",
		CustomerEmail: "juan@example.com",
		Status:        models.OrderStatusPending,
		TotalPrice:    2500.00,
	})

	svc := NewOrdersService(mockRepo)

	// Act: Call the service method to delete the order.
	err := svc.DeleteOrder(created.ID)
	// Assert: Verify no error occurred.
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify the order was actually deleted.
	_, err = svc.GetOrderByID(created.ID)
	if err == nil {
		t.Fatal("expected an error when retrieving deleted order")
	}
}

// TestDeleteNonexistentOrder tests that DeleteOrder returns error for non-existent order.
func TestDeleteNonexistentOrder(t *testing.T) {
	// Arrange: Set up the mock repository (empty).
	mockRepo := NewMockOrdersRepository()
	svc := NewOrdersService(mockRepo)

	// Act: Try to delete a non-existent order.
	err := svc.DeleteOrder(999)

	// Assert: Verify that an error was returned.
	if err == nil {
		t.Fatal("expected an error when deleting non-existent order, but got none")
	}
}

func TestOrdersService_GetTopPartsStats(t *testing.T) {
	tests := []struct {
		name          string
		limitInput    int
		expectedLimit int
		mockReturn    []dto.TopPartResponse
		mockErr       error
		expectError   bool
	}{
		{
			name:          "Éxito con límite explícito",
			limitInput:    3,
			expectedLimit: 3,
			mockReturn: []dto.TopPartResponse{
				{PartID: 1, PartName: "Motor V6", TotalOrdered: 10},
				{PartID: 2, PartName: "Turbo", TotalOrdered: 5},
				{PartID: 3, PartName: "Discos de freno", TotalOrdered: 2},
			},
			mockErr:     nil,
			expectError: false,
		},
		{
			name:          "Éxito con fallback a límite por defecto (5)",
			limitInput:    0,
			expectedLimit: 5,
			mockReturn:    []dto.TopPartResponse{{PartID: 1, PartName: "Motor V8", TotalOrdered: 10}},
			mockErr:       nil,
			expectError:   false,
		},
		{
			name:          "Fallo en el repositorio",
			limitInput:    5,
			expectedLimit: 5,
			mockReturn:    nil,
			mockErr:       fmt.Errorf("database error"),
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockOrdersRepository{
				MockGetTopParts: func(limit int) ([]dto.TopPartResponse, error) {
					if limit != tt.expectedLimit {
						t.Errorf("se esperaba el límite %d, pero se recibió %d", tt.expectedLimit, limit)
					}
					return tt.mockReturn, tt.mockErr
				},
			}

			svc := NewOrdersService(mockRepo)

			result, err := svc.GetTopPartsStats(tt.limitInput)

			if tt.expectError {
				if err == nil {
					t.Error("se esperaba un error pero no ocurrió ninguno")
				}
			} else {
				if err != nil {
					t.Errorf("no se esperaba error, se obtuvo: %v", err)
				}
				if len(result) != len(tt.mockReturn) {
					t.Errorf("se esperaban %d resultados, se obtuvieron %d", len(tt.mockReturn), len(result))
				}
			}
		})
	}
}
