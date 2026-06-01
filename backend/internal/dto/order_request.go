// Package dto defines request and response data transfer objects for API endpoints.
package dto

// OrderItemRequest represents a single item in an order creation request.
type OrderItemRequest struct {
	// PartID is the ID of the part being ordered.
	PartID uint `json:"part_id" binding:"required,gt=0"`

	// Quantity is the number of parts to order.
	Quantity int `json:"quantity" binding:"required,gt=0,max=1000"`
}

// CreatedOrderRequest represents the request payload for creating a new order.
type CreateOrderRequest struct {
	// CustomerName is the name of the customer placing the order.
	CustomerName string `json:"customer_name" binding:"required,min=1,max=100"`

	// CustomerEmail is the email address of the customer.
	CustomerEmail string `json:"customer_email" binding:"required,email"`

	// Items is a slice of parts to include in the order.
	Items []OrderItemRequest `json:"items" binding:"required,min=1,max=100"`
}

// UpdateOrderStatusRequest represents the request payload for updating an order's status.
type UpdateOrderStatusRequest struct {
	// Status is the new status fot the order.
	// Valid values: pending, processing, completed or cancelled.
	Status string `json:"status" binding:"required,oneof=pending processing completed cancelled"`
}

// OrderItemResponse represents a single item in an order response.
type OrderItemResponse struct {
	// ID is the unique identifier for the order item.
	ID uint `json:"id"`

	// OrderID is the ID of the parent order.
	OrderID uint `json:"order_id"`

	// PartID is the ID of the part.
	PartID uint `json:"part_id"`

	// Quantity is the number of parts ordered.
	Quantity int `json:"quantity"`

	// UnitPrice is the price of a single part at the time or ordering.
	UnitPrice float64 `json:"unit_price"`

	// Part contains the part details.
	Part *PartResponse `json:"part,omitempty"`
}

// OrderResponse represents the response payload for an order.
type OrderResponse struct {
	// ID is the unique identifier of the order.
	ID uint `json:"id"`

	// CustomerName is the name of the customer.
	CustomerName string `json:"customer_name"`

	// CustomerEmail is the email address of the customer.
	CustomerEmail string `json:"customer_email"`

	// Status is the current status of the order.
	Status string `json:"status"`

	// TotalPrice is the total cost of the order.
	TotalPrice float64 `json:"total_price"`

	// CreatedAt is when the order was created.
	CreatedAt string `json:"created_at"`

	// UpdatedAt is when the order was last updated.
	UpdatedAt string `json:"updated_at"`

	// OrderItems contains all items in the order.
	OrderItems []OrderItemResponse `json:"order_item,omitempty"`
}
