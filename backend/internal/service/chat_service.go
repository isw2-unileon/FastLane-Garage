// Package service contains the business logic for the application.
package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/dto"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/models"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/repository"
)

// ChatService defines the interface for chat business operations.
type ChatService interface {
	// CreateSession initializes a new chat with optional vehicle context.
	CreateSession(req *dto.CreateChatSessionRequest) (*dto.ChatSessionResponse, error)

	// GetSessionHistory retrieves a session and its full message history.
	GetSessionHistory(sessionID uint) (*dto.ChatSessionResponse, error)

	// SendMessage processes a user message, calls the n8n webhook, and returns the sync response.
	SendMessage(sessionID uint, req *dto.SendChatMessageRequest) (*dto.ChatSyncResponse, error)
}

// chatServiceImpl is the concrete implementation of ChatService.
type chatServiceImpl struct {
	repo          repository.ChatRepository
	n8nwebhookURL string
	httpClient    *http.Client
}

// n8nRequest represents the payload sent to the n8n webhook.
type n8nRequest struct {
	SessionID string `json:"sessionId"`
	Message   string `json:"message"`
	Vehicle   struct {
		Brand string `json:"brand"`
		Model string `json:"model"`
		Year  string `json:"year"`
	} `json:"vehicle"`
	Parts []dto.PartItemDTO `json:"parts,omitempty"`
}

// n8nResponse represents the expected JSON response from the n8n webhook.
type n8nResponse struct {
	Reply   string `json:"reply"`
	Status  string `json:"status"`
	Vehicle struct {
		Brand string `json:"brand"`
		Model string `json:"model"`
		Year  string `json:"year"`
	} `json:"vehicle"`
}

// NewChatService creates a new instance of the chat service.
// It accepts a ChatRepository and the n8n webhook URL.
func NewChatService(repo repository.ChatRepository, n8nURL string) ChatService {
	// If no URL is provided, fallback to a default local n8n webhook URL.
	if n8nURL == "" {
		n8nURL = "http://localhost:5678/webhook/chat"
	}

	return &chatServiceImpl{
		repo:          repo,
		n8nwebhookURL: n8nURL,
		// Using a custom HTTP client with a timeout is a best practice
		// to prevent the Go server from hanging if n8n is slow or down.
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CreateSession initializes a new chat session in the database.
func (s *chatServiceImpl) CreateSession(req *dto.CreateChatSessionRequest) (*dto.ChatSessionResponse, error) {
	partsJSON := "[]"
	if len(req.Parts) > 0 {
		bytes, err := json.Marshal(req.Parts)
		if err == nil {
			partsJSON = string(bytes)
		}
	}

	session := &models.ChatSession{
		Status:       "need_info",
		VehicleBrand: req.VehicleBrand,
		VehicleModel: req.VehicleModel,
		VehicleYear:  req.VehicleYear,
		Parts:        partsJSON,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	createdSession, err := s.repo.CreateSession(session)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	response := s.sessionToResponse(createdSession)
	return &response, nil
}

// GetSessionHistory retrieves a chat session and its messages.
func (s *chatServiceImpl) GetSessionHistory(sessionID uint) (*dto.ChatSessionResponse, error) {
	session, err := s.repo.GetSessionByID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	response := s.sessionToResponse(session)
	return &response, nil
}

// SendMessage handles the core chat workflow: save user message, call n8n, save AI response.
func (s *chatServiceImpl) SendMessage(sessionID uint, req *dto.SendChatMessageRequest) (*dto.ChatSyncResponse, error) {
	// 1. Verify the session exists.
	session, err := s.repo.GetSessionByID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	// 2. Save the user's message to the database.
	userMsg := &models.ChatMessage{
		SessionID: sessionID,
		Role:      "user",
		Content:   req.Text,
		Source:    "frontend",
		CreatedAt: time.Now(),
	}
	createdUserMsg, err := s.repo.CreateMessage(userMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to save user message: %w", err)
	}

	// 3. Prepare the payload for n8n.
	payload := n8nRequest{
		SessionID: strconv.FormatUint(uint64(sessionID), 10),
		Message:   req.Text,
	}
	payload.Vehicle.Brand = session.VehicleBrand
	payload.Vehicle.Model = session.VehicleModel
	payload.Vehicle.Year = session.VehicleYear

	if session.Parts != "" && session.Parts != "[]" {
		var parts []dto.PartItemDTO
		_ = json.Unmarshal([]byte(session.Parts), &parts)
		payload.Parts = parts
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal n8n request: %w", err)
	}

	// 4. Call the n8n webhook synchronously.
	httpResp, err := s.httpClient.Post(s.n8nwebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to reach n8n webhook: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("n8n returned non-200 status code: %d", httpResp.StatusCode)
	}

	// 5. Decode the n8n response.
	var n8nRes n8nResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&n8nRes); err != nil {
		return nil, fmt.Errorf("failed to decode n8n response: %w", err)
	}

	// 6. Update the session with the new context provided by n8n.
	session.Status = n8nRes.Status
	session.VehicleBrand = n8nRes.Vehicle.Brand
	session.VehicleModel = n8nRes.Vehicle.Model
	session.VehicleYear = n8nRes.Vehicle.Year
	session.UpdatedAt = time.Now()

	if err := s.repo.UpdateSession(session); err != nil {
		return nil, fmt.Errorf("failed to update session state: %w", err)
	}

	// 7. Save the assistant's reply to the database.
	assistantMsg := &models.ChatMessage{
		SessionID: sessionID,
		Role:      "assistant",
		Content:   n8nRes.Reply,
		Source:    "n8n",
		CreatedAt: time.Now(),
	}
	createdAssistantMsg, err := s.repo.CreateMessage(assistantMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to save assistant message: %w", err)
	}

	// 8. Build and return the final DTO for the frontend.
	syncResponse := &dto.ChatSyncResponse{
		UserMessage:      s.messageToResponse(createdUserMsg),
		AssistantMessage: s.messageToResponse(createdAssistantMsg),
		Session:          s.sessionToResponse(session),
	}
	return syncResponse, nil
}

// sessionToResponse converts a ChatSession model to a ChatSessionResponse DTO.
func (s *chatServiceImpl) sessionToResponse(session *models.ChatSession) dto.ChatSessionResponse {
	response := dto.ChatSessionResponse{
		ID:           session.ID,
		Status:       session.Status,
		VehicleBrand: session.VehicleBrand,
		VehicleModel: session.VehicleModel,
		VehicleYear:  session.VehicleYear,
		CreatedAt:    session.CreatedAt.Format(time.RFC3339),
	}

	if session.Parts != "" && session.Parts != "[]" {
		var parts []dto.PartItemDTO
		_ = json.Unmarshal([]byte(session.Parts), &parts)
		response.Parts = parts
	}

	// Convert associated messages if they exist.
	if len(session.Messages) > 0 {
		msgResponses := make([]dto.ChatMessageResponse, len(session.Messages))
		for i, msg := range session.Messages {
			msgResponses[i] = s.messageToResponse(&msg)
		}
		response.Messages = msgResponses
	}

	return response
}

// messageToResponse converts a ChatMessage model to a ChatMessageResponse DTO.
func (s *chatServiceImpl) messageToResponse(msg *models.ChatMessage) dto.ChatMessageResponse {
	return dto.ChatMessageResponse{
		ID:        msg.ID,
		Role:      msg.Role,
		Content:   msg.Content,
		CreatedAt: msg.CreatedAt.Format(time.RFC3339),
	}
}
