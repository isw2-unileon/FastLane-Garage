// Package repository handles all database operations for business logic.
package repository

import (
	"errors"
	"fmt"

	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/models"
	"gorm.io/gorm"
)

// ChatRepository defines the interface for chat data operations.
// Using an interface allows us to mock the repository in service tests.
type ChatRepository interface {
	// CreateSession inserts a new chat session into the database.
	// Returns the created session with its generated ID.
	CreateSession(session *models.ChatSession) (*models.ChatSession, error)

	// GetSessionByID retrieves a single chat session by its ID,
	// including all its associated messages.
	// Returns an error if the session is not found.
	GetSessionByID(id uint) (*models.ChatSession, error)

	// UpdateSession modifies an existing chat session (e.g., updating status or vehicle data).
	UpdateSession(session *models.ChatSession) error

	// CreateMessage inserts a new chat message into the database.
	// Returns the created message with its generated ID.
	CreateMessage(message *models.ChatMessage) (*models.ChatMessage, error)
}

// chatRepositoryImpl is the concrete implementation of ChatRepository using GORM.
type chatRepositoryImpl struct {
	db *gorm.DB
}

// NewChatRepository creates a new instance of the chat repository.
// It accepts a GORM database connection and returns a ChatRepository interface.
func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepositoryImpl{
		db: db,
	}
}

// CreateSession inserts a new chat session into the database.
func (r *chatRepositoryImpl) CreateSession(session *models.ChatSession) (*models.ChatSession, error) {
	if err := r.db.Create(session).Error; err != nil {
		return nil, fmt.Errorf("failed to create chat session: %w", err)
	}
	return session, nil
}

// GetSessionByID retrieves a single chat session by its ID along with its message.
func (r *chatRepositoryImpl) GetSessionByID(id uint) (*models.ChatSession, error) {
	var session models.ChatSession

	// Preload("Messages") automatically fetches the associated chat messages
	// and populates the Messages slice in the ChatSession struct.
	// We order messages by CreatedAt to ensure they appear in chronological order.
	err := r.db.Preload("Messages", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at ASC")
	}).First(&session, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound // ErrNotFound is defined in parts_repository.go
		}
		return nil, fmt.Errorf("failed to fetch chat session with ID %d: %w", id, err)
	}

	return &session, nil
}

// UpdateSession updates the fields of an existing chat session.
func (r *chatRepositoryImpl) UpdateSession(session *models.ChatSession) error {
	// Omit("Messages") prevents GORM from trying to update or replace
	// the associated messages when we only want to update the session metadata.
	if err := r.db.Omit("Messages").Save(session).Error; err != nil {
		return fmt.Errorf("failed to update chat session with ID %d: %w", session.ID, err)
	}
	return nil
}

// CreateMessage inserts a new chat message associated with a session.
func (r *chatRepositoryImpl) CreateMessage(message *models.ChatMessage) (*models.ChatMessage, error) {
	if err := r.db.Create(message).Error; err != nil {
		return nil, fmt.Errorf("failed to create chat message: %w", err)
	}
	return message, nil
}
