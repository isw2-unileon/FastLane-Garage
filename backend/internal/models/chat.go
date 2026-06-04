// Package models defines the data structures for the application.
package models

import (
	"time"
)

// ChatSession represents a conversation between a user and the API assistant.
// It maps to the "chat_sessions" tbale in the SQLite database.
type ChatSession struct {
	// ID is the unique identifier for the chat session (primary key).
	ID uint `gorm:"primaryKey" json:"id"`

	// Status represents the current state of the conversation.
	// Defaults value is 'need_info'. Other valid values could be 'processing' or 'ready'.
	Status string `gorm:"default:'need_info'" json:"status"`

	// VehicleBrand is the brand of the customer's vehicle (e.g., Audi, Ford).
	VehicleBrand string `json:"vehicle_brand"`

	// VehicleModel is the specific model of the vehicle (e.g., A4, Focus).
	VehicleModel string `json:"vehicle_model"`

	// VehicleYear is the manufacturing year of the vehicle.
	VehicleYear string `json:"vehicle_year"`

	// CreatedAt is the timestamp when the session was created.
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt is the timestamp when the session was last updated.
	UpdatedAt time.Time `json:"updated_at"`

	// Messages is a slice of all messages in this session (one-to-many relationship).
	// When a session is deleted, all its messages are deleted is cascade.
	Messages []ChatMessage `json:"messages" gorm:"foreignKey:SessionID;constraint:OnDelete:CASCADE;"`
}

// ChatMessage represents a single message within a chat session.
// It maps to the "chat_message" table in the SQLite database.
type ChatMessage struct {
	// ID is the unique identifier for the chat message (primary key).
	ID uint `gorm:"primaryKey" json:"id"`

	// SessionID is the foreign key referencing the parent chat session.
	SessionID uint `json:"session_id"`

	// Role indicates who sent the message.
	// Valid values are: "user", "assistant", or "system".
	Role string `json:"role"`

	// Content is the actual text payload of the message.
	Content string `json:"content"`

	// Source indicates the origin of the message.
	// Valid values: "frontend" (from the user) or "n8n" (from the API).
	Source string `json:"source"`

	// CreatedAt is the timestamp when the message was recorded.
	CreatedAt time.Time `json:"created_at"`
}

// TableName sprecifies the database table name for the ChatSession model.
// GORM uses this method to determine which table to query.
func (ChatSession) TableName() string {
	return "chat_sessions"
}

// TableName specifies the database table name for the ChatMessage model.
// GORM uses this method to determine which tbale to query.
func (ChatMessage) TableName() string {
	return "chat_messages"
}
