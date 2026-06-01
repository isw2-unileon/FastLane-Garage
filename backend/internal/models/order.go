// Package models defines the data structures for the application.
package models

import "time"

// OrderStatus represents the current state of an order.
// Valid values are: pending, processing, completed, cancelled.
type OrderStatus string

const (
	// OrderStatusPending indicates the order has been created but not yet processed.
	OrderStatusPending OrderStatus = "pending"

	// OrderStatusProcessing indicates the order is being processed.
	OrderStatusProcessing OrderStatus = "processing"

	// OrderStatusCompleted indicates the order has been successfully completed.
	OrderStatusCompleted OrderStatus = "completed"

	// OrderStatusCancelled indicates the order has been cancelled.
	OrderStatusCancelled OrderStatus = "cancelled"
)

// Order represents a customer order in the system.
// It maps to the "orders" table in the SQLite database.
type Order struct {
	// ID is the unique identifier for the order (primary key).
	ID uint `gorm:"primaryKey" json:"id"`

	// CustomerName is the name of the customer who placed the order.
	CustomerName string `json:"customer_name"`

	// CustomerEmail is the email address of the cutomer.
	CustomerEmail string `json:"customer_email"`

	// Status represents the current state of the order.
	// Valid values: pending, processing, completed or cancelled.
	Status OrderStatus `json:"status"`

	// TotalPrice is the sum of all items in the order.
	TotalPrice float64 `json:"total_price"`

	// CreatedAt is the timestamp when the order was created.
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt is the timestamp when the order was last updated.
	UpdatedAt time.Time `json:"updated_at"`

	// OrderItems is a slice of items in this order (one-to-many relationship).
	OrderItems []OrderItem `json:"order_items" gorm:"foreignKey:OrderID"`
}

// OrderItem represents a single item (part) in an order.
// It maps to the "order_items" table in the SQLite database.
type OrderItem struct {
	// ID is the unique identifier for the order item (primary key).
	ID uint `gorm:"primaryKey" json:"id"`

	// OrderID is the foreign key referencing the parent order.
	OrderID uint `json:"order_id"`

	// PartID is the foreign key referencing the part that was ordered.
	PartID uint `json:"part_id"`

	// Quantity is the number of parts ordered.
	Quantity int `json:"quantity"`

	// UnitPrice is the price of a single part at the time of ordering.
	UnitPrice float64 `json:"unit_price"`

	// Part is a reference to the part details (populated via GORM association).
	Part *Part `json:"part,omitempty" gorm:"foreignKey:PartID"`
}

// TableName specifies the database table name for the Order model.
// GORM uses this method to determine which table to query.
func (Order) TableName() string {
	return "orders"
}

// TableName specifies the database table name for the OrderItem model.
// GORM uses this method to determine which table to query.
func (OrderItem) TableName() string {
	return "order_items"
}
