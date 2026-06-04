// Package service contains business logic for orders operations.
package service

import (
	"fmt"
	"time"

	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/dto"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/models"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/repository"
)

// OrdersService defines the interface for order business operations.
// This abstraction allows us to mock the service in tests.
type OrdersService interface {
	// GetAllOrders retrieves all orders from the repository.
	GetAllOrders(status string, email string) ([]dto.OrderResponse, error)

	// GetOrderByID retrieves a single order by its ID.
	GetOrderByID(id uint) (*dto.OrderResponse, error)

	// CreateOrder validates the request and creates a new order.
	CreateOrder(req *dto.CreateOrderRequest, partsRepo repository.PartsRepository) (*dto.OrderResponse, error)

	// UpdateOrderStatus changes the status of an existing order.
	UpdateOrderStatus(id uint, req *dto.UpdateOrderStatusRequest) (*dto.OrderResponse, error)

	// DeleteOrder removes an order by ID.
	DeleteOrder(id uint) error

	// GetTopPartsStats retrieves stadistics about the most ordered parts.
	GetTopPartsStats(limit int) ([]dto.TopPartResponse, error)
}

// ordersServiceImpl is the concrete implementation of OrdersService.
type ordersServiceImpl struct {
	repo repository.OrdersRepository
}

// NewOrdersService creates a new instance of the orders service.
// It accepts an OrdersRepository and return an OrdersService interface.
func NewOrdersService(repo repository.OrdersRepository) OrdersService {
	return &ordersServiceImpl{
		repo: repo,
	}
}

// GetAllOrders retrieves all orders and converts them to DTOs.
func (s *ordersServiceImpl) GetAllOrders(status string, email string) ([]dto.OrderResponse, error) {
	// Fetch all orders from the repository.
	orders, err := s.repo.FindAll(status, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get all orders: %w", err)
	}

	// Apply filters in the service layer if needed.
	var filtered []models.Order
	for _, order := range orders {
		// Filter by status if provided.
		if status != "" && order.Status != models.OrderStatus(status) {
			continue
		}
		// Filter by email if provided.
		if email != "" && order.CustomerEmail != email {
			continue
		}
		filtered = append(filtered, order)
	}

	// Convert Order models to OrderResponse DTOs.
	responses := make([]dto.OrderResponse, len(filtered))
	for i, order := range filtered {
		responses[i] = s.orderToResponse(&order)
	}

	return responses, nil
}

// GetOrderByID retrieves a single order by ID and converts it to a DTO.
func (s *ordersServiceImpl) GetOrderByID(id uint) (*dto.OrderResponse, error) {
	// Fetch the order from the repository.
	order, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order with ID %d: %w", id, err)
	}

	// Convert to DTO and return.
	response := s.orderToResponse(order)
	return &response, nil
}

// CreateOrder validates the request and creates a new order with its items.
func (s *ordersServiceImpl) CreateOrder(req *dto.CreateOrderRequest, partsRepo repository.PartsRepository) (*dto.OrderResponse, error) {
	// Validate request.
	if req.CustomerName == "" || req.CustomerEmail == "" {
		return nil, fmt.Errorf("customer name and email are required")
	}

	if len(req.Items) == 0 {
		return nil, fmt.Errorf("order must contain at least one item")
	}

	// Validate parts exist and calculate total price.
	var totalPrice float64
	orderItems := make([]models.OrderItem, len(req.Items))

	for i, item := range req.Items {
		// Verify the part exists.
		part, err := partsRepo.FindByID(item.PartID)
		if err != nil {
			return nil, fmt.Errorf("part with ID %d not found: %w", item.PartID, err)
		}

		// Calculate line total.
		lineTotal := part.Price * float64(item.Quantity)
		totalPrice += lineTotal

		// Create order item.
		orderItems[i] = models.OrderItem{
			PartID:    item.PartID,
			Quantity:  item.Quantity,
			UnitPrice: part.Price,
		}
	}

	// Create the order model.
	order := &models.Order{
		CustomerName:  req.CustomerName,
		CustomerEmail: req.CustomerEmail,
		Status:        models.OrderStatusPending,
		TotalPrice:    totalPrice,
		OrderItems:    orderItems,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Create the order in the repository.
	createdOrder, err := s.repo.Create(order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Convert to DTO return.
	response := s.orderToResponse(createdOrder)
	return &response, nil
}

// UpdateOrderStatus changes the status of an existing order.
func (s *ordersServiceImpl) UpdateOrderStatus(id uint, req *dto.UpdateOrderStatusRequest) (*dto.OrderResponse, error) {
	// First, verify the order exists.
	if _, err := s.repo.FindByID(id); err != nil {
		return nil, err
	}

	// Convert string status to OrderStatus enum.
	status := models.OrderStatus(req.Status)

	// Update the order status.
	updatedOrder, err := s.repo.UpdateStatus(id, status)
	if err != nil {
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	// Convert to DTO and return.
	response := s.orderToResponse(updatedOrder)
	return &response, nil
}

// DeleteOrder removes an order by ID.
func (s *ordersServiceImpl) DeleteOrder(id uint) error {
	// Verify the order exists before deleting.
	if _, err := s.repo.FindByID(id); err != nil {
		return err
	}

	// Delete the order.
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	return nil
}

// orderToResponse converts an Order model to an OrderResponse DTO.
func (s *ordersServiceImpl) orderToResponse(order *models.Order) dto.OrderResponse {
	// Converts order items.
	itemResponses := make([]dto.OrderItemResponse, len(order.OrderItems))
	for i, item := range order.OrderItems {
		itemResponses[i] = s.orderItemToResponse(&item)
	}

	return dto.OrderResponse{
		ID:            order.ID,
		CustomerName:  order.CustomerName,
		CustomerEmail: order.CustomerEmail,
		Status:        string(order.Status),
		TotalPrice:    order.TotalPrice,
		CreatedAt:     order.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     order.UpdatedAt.Format(time.RFC3339),
		OrderItems:    itemResponses,
	}
}

// orderItemToResponse converts an OrderItem model to an OrderItemResponse DTO.
func (s *ordersServiceImpl) orderItemToResponse(item *models.OrderItem) dto.OrderItemResponse {
	response := dto.OrderItemResponse{
		ID:        item.ID,
		OrderID:   item.OrderID,
		PartID:    item.PartID,
		Quantity:  item.Quantity,
		UnitPrice: item.UnitPrice,
	}

	// Include part details if available.
	if item.Part != nil {
		partResponse := &dto.PartResponse{
			ID:       item.Part.ID,
			Name:     item.Part.Name,
			CarZone:  item.Part.CarZone,
			ImageURL: item.Part.ImageURL,
			Price:    item.Part.Price,
		}
		response.Part = partResponse
	}

	return response
}

// GetTopPartsStats
func (s *ordersServiceImpl) GetTopPartsStats(limit int) ([]dto.TopPartResponse, error) {
	if limit <= 0 {
		limit = 5
	}
	return s.repo.GetTopParts(limit)
}
