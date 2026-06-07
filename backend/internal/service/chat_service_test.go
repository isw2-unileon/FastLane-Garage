// Package service contains tests for the chat service business logic.
package service

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/dto"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/models"
)

// mockChatRepository is a mock implementation of ChatRepository for testing.
type mockChatRepository struct {
	MockCreateSession  func(session *models.ChatSession) (*models.ChatSession, error)
	MockGetSessionByID func(id uint) (*models.ChatSession, error)
	MockUpdateSession  func(session *models.ChatSession) error
	MockCreateMessage  func(message *models.ChatMessage) (*models.ChatMessage, error)
}

func (m *mockChatRepository) CreateSession(session *models.ChatSession) (*models.ChatSession, error) {
	if m.MockCreateSession != nil {
		return m.MockCreateSession(session)
	}
	session.ID = 1
	return session, nil
}

func (m *mockChatRepository) GetSessionByID(id uint) (*models.ChatSession, error) {
	if m.MockGetSessionByID != nil {
		return m.MockGetSessionByID(id)
	}
	return &models.ChatSession{ID: id}, nil
}

func (m *mockChatRepository) UpdateSession(session *models.ChatSession) error {
	if m.MockUpdateSession != nil {
		return m.MockUpdateSession(session)
	}
	return nil
}

func (m *mockChatRepository) CreateMessage(message *models.ChatMessage) (*models.ChatMessage, error) {
	if m.MockCreateMessage != nil {
		return m.MockCreateMessage(message)
	}
	message.ID = 1
	return message, nil
}

// TestCreateSession tests the initialization of a new chat session.
func TestCreateSession(t *testing.T) {
	// Arrange: setup the mock repository
	mockRepo := &mockChatRepository{
		MockCreateSession: func(session *models.ChatSession) (*models.ChatSession, error) {
			session.ID = 99
			session.CreatedAt = time.Now()
			return session, nil
		},
	}
	svc := NewChatService(mockRepo, "")

	req := &dto.CreateChatSessionRequest{
		VehicleBrand: "Audi",
		VehicleModel: "A4",
	}

	// Act: call the service method
	res, err := svc.CreateSession(req)
	// Assert: verify the results
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.ID != 99 {
		t.Errorf("expected session ID 99, got %d", res.ID)
	}
	if res.VehicleBrand != "Audi" {
		t.Errorf("expected brand Audi, got %s", res.VehicleBrand)
	}
}

// TestGetSessionHistory tests retrieving a session with its messages.
func TestGetSessionHistory(t *testing.T) {
	mockRepo := &mockChatRepository{
		MockGetSessionByID: func(id uint) (*models.ChatSession, error) {
			return &models.ChatSession{
				ID:           id,
				VehicleBrand: "Ford",
				Messages: []models.ChatMessage{
					{ID: 1, Role: "user", Content: "Hello"},
					{ID: 2, Role: "assistant", Content: "Hi there"},
				},
			}, nil
		},
	}
	svc := NewChatService(mockRepo, "")

	res, err := svc.GetSessionHistory(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.VehicleBrand != "Ford" {
		t.Errorf("expected brand Ford, got %s", res.VehicleBrand)
	}
	if len(res.Messages) != 2 {
		t.Errorf("expected 2 messages, got %d", len(res.Messages))
	}
}

// TestSendMessage tests the full synchronous flow including the n8n webhook call.
func TestSendMessage(t *testing.T) {
	// 1. Create a fake HTTP server to simulate n8n.
	// This server will return a valid JSON response whenever it's called.
	fakeN8nServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Mock n8n response payload
		_, _ = w.Write([]byte(`{
			"reply": "I am a mocked AI",
			"status": "ready",
			"vehicle": {"brand": "Tesla", "model": "Model 3", "year": "2023"}
		}`))
	}))
	defer fakeN8nServer.Close() // Ensure the server is shut down after the test

	// 2. Set up the mock repository
	mockRepo := &mockChatRepository{
		MockGetSessionByID: func(id uint) (*models.ChatSession, error) {
			return &models.ChatSession{ID: id, VehicleBrand: "Tesla"}, nil
		},
		MockUpdateSession: func(session *models.ChatSession) error {
			// Verify that the service updated the status based on n8n's response
			if session.Status != "ready" {
				t.Errorf("expected session status 'ready', got '%s'", session.Status)
			}
			return nil
		},
	}

	// 3. Instantiate the service, injecting the fake n8n server URL instead of the real one!
	svc := NewChatService(mockRepo, fakeN8nServer.URL)

	req := &dto.SendChatMessageRequest{
		Text: "I need a new battery",
	}

	// 4. Act: Call the SendMessage method
	res, err := svc.SendMessage(1, req)
	// 5. Assert: Verify the entire flow succeeded
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.AssistantMessage.Content != "I am a mocked AI" {
		t.Errorf("expected assistant message 'I am a mocked AI', got '%s'", res.AssistantMessage.Content)
	}

	if res.Session.Status != "ready" {
		t.Errorf("expected final session status 'ready', got '%s'", res.Session.Status)
	}
}
