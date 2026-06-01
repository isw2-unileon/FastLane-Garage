// Package handlers contains HTTP request handlers for the API endpoints.
package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/repository"
)

// parseID extracts and validates an ID from URL parameters.
// Returns the parsed ID or an error response if invalid.
func parseID(c *gin.Context) (uint, bool) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id format",
		})
		return 0, false
	}
	return uint(id), true
}

// handleNotFoundError checks if an error is "not found" and responds appropriately.
// Returns true if it was a not-found error, false otherwise.
func handleNotFoundError(c *gin.Context, err error, entityType string) bool {
	if errors.Is(err, repository.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": entityType + " not found",
		})
		return true
	}
	return false
}

// handleInternalError responds with a 500 error message.
func handleInternalError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": message,
	})
}
