// Package handlers contains the HTTP request handlers for the application.
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/dto"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/service"
)

// CreateChatSession handles the POST request to initialize a new chat session.
func CreateChatSession(chatSvc service.ChatService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CreateChatSessionRequest

		// Bind JSON payload to the request DTO.
		if err := c.ShouldBindJSON(&req); err != nil && err.Error() != "EOF" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		response, err := chatSvc.CreateSession(&req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, response)
	}
}

// GetChatSessionHistory handles the GET request to retrieve a chat session's history.
func GetChatSessionHistory(chatSvc service.ChatService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract session ID from the URL parameter
		idParam := c.Param("id")
		sessionID, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID format"})
			return
		}

		response, err := chatSvc.GetSessionHistory(uint(sessionID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

// SendChatMessage handles the POST request when a user sends a message.
func SendChatMessage(chatSvc service.ChatService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		sessionID, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID format"})
			return
		}

		var req dto.SendChatMessageRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload. 'text' field is required."})
			return
		}

		response, err := chatSvc.SendMessage(uint(sessionID), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, response)
	}
}
