// Package handlers contains HTTP request handlers for the API endpoints.
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/models"
	"gorm.io/gorm"
)

// GetParts returns a Gin handler function that fetches all parts from the database.
// It accepts a GORM database instance and returns a handler that can be registered
// to a Gin route. This pattern (handler factory) allows dependency injection of the DB.
func GetParts(db *gorm.DB) gin.HandlerFunc {
	// Return the actual handler function that will be called for each HTTP request.
	return func(c *gin.Context) {
		// Initialize an empty slice to store parts retrieved from the database.
		var parts []models.Part

		// Query the database to fetch all parts.
		// The Find method retrieves all records without any filtering conditions.
		if err := db.Find(&parts).Error; err != nil {
			// If an error occurs (database connection failure, etc.),
			// respond with HTTP 500 Internal Server Error.
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to fetch parts",
			})
			return
		}

		// Safety check: ensure parts is never nil in the JSON response.
		// If the database returns nil for an empty result, convert it to an empty slice.
		if parts == nil {
			parts = []models.Part{}
		}

		// Return the parts as JSON with HTTP 200 OK status.
		// The response format wraps parts in a "data" field for API consistency.
		c.JSON(http.StatusOK, gin.H{
			"data": parts,
		})
	}
}
