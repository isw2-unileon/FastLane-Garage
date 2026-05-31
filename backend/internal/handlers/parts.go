// Package handlers contains HTTP request handlers for the API endpoints.
package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/dto"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/service"
	"gorm.io/gorm"
)

// GetParts returns a Gin handler function that fetches all parts from the database.
// It accepts a PartsService and returns a handler that can be registered
// to a Gin route. This pattern (handler factory) allows dependency injection of the service.
func GetParts(svc service.PartsService) gin.HandlerFunc {
	// Return the actual handler function that will be called for each HTTP request.
	return func(c *gin.Context) {
		// Call the service to get all parts.
		parts, err := svc.GetAllParts()
		if err != nil {
			// If an error occurs, respond with HTTP 500 Internal Server Error.
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to fetch parts",
			})
			return
		}

		// Return the parts as JSON with HTTP 200 OK status.
		// The response format wraps parts in a "data" field for API consistency.
		c.JSON(http.StatusOK, gin.H{
			"data": parts,
		})
	}
}

// GetPartByID returns a Gin handler function that fetches a single part by ID.
// The part ID is extracted from the URL path parameter (e.g., /api/parts/1).
func GetPartByID(svc service.PartsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the part ID from the URL path.
		idParam := c.Param("id")

		// Convert the ID string to an unsigned integer.
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			// If the ID is not a valid number, return HTTP 400 Bad Request.
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid part ID format",
			})
			return
		}

		// Call the service to get the part by ID.
		part, err := svc.GetPartByID(uint(id))
		if err != nil {
			// Check if the error is "record not found".
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// If the part is not found, return HTTP 404 Not Found.
				c.JSON(http.StatusNotFound, gin.H{
					"error": "part not found",
				})
				return
			}

			// For other errors, return HTTP 500 Internal Server Error.
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to fetch part",
			})
			return
		}

		// Return the part as JSON with HTTP 200 OK status.
		c.JSON(http.StatusOK, gin.H{
			"data": part,
		})
	}
}

// CreatePart returns a Gin handler function that creates a new part.
// The request body is expected to be JSON matching CreatePartRequest.
func CreatePart(svc service.PartsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse and validate the request body.
		var req dto.CreatePartRequest

		// BindJSON parses the JSON and validates it using the binding tags defined in the DTO.
		if err := c.BindJSON(&req); err != nil {
			// If binding fails (invalid JSON or validation error),
			// return HTTP 400 Bad Request with the error message.
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Call the service to create the part.
		part, err := svc.CreatePart(&req)
		if err != nil {
			// If service returns an error, return HTTP 500 Internal Server Error.
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Return the created part with HTTP 201 Created status.
		// HTTP 201 indicates that a new resource has been created.
		c.JSON(http.StatusCreated, gin.H{
			"data": part,
		})
	}
}

// UpdatePart returns a Gin handler function that updates an existing part.
// The part ID is extracted from the URL path, and the new data is in the request body.
func UpdatePart(svc service.PartsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the part ID from the URL path.
		idParam := c.Param("id")

		// Convert the ID string to an unsigned integer.
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			// If the ID is not a valid number, return HTTP 400 Bad Request.
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid part ID format",
			})
			return
		}

		// Parse and validate the request body.
		var req dto.UpdatePartRequest

		// BindJSON parses the JSON and validates it using the binding tags.
		if err := c.BindJSON(&req); err != nil {
			// If binding fails, return HTTP 400 Bad Request.
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Call the service to update the part.
		part, err := svc.UpdatePart(uint(id), &req)
		if err != nil {
			// Check for specific error types.
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// If the part is not found, return HTTP 404 Not Found.
				c.JSON(http.StatusNotFound, gin.H{
					"error": "part not found",
				})
				return
			}

			// For other errors, return HTTP 500 Internal Server Error.
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Return the updated part with HTTP 200 OK status.
		c.JSON(http.StatusOK, gin.H{
			"data": part,
		})
	}
}

// DeletePart returns a Gin handler function that deletes a part by ID.
// The part ID is extracted from the URL path parameter.
func DeletePart(svc service.PartsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the part ID from the URL path.
		idParam := c.Param("id")

		// Convert the ID string to an unsigned integer.
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			// If the ID is not a valid number, return HTTP 400 Bad Request.
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid part ID format",
			})
			return
		}

		// Call the service to delete the part.
		if err := svc.DeletePart(uint(id)); err != nil {
			// Check if the error is "record not found".
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// If the part is not found, return HTTP 404 Not Found.
				c.JSON(http.StatusNotFound, gin.H{
					"error": "part not found",
				})
				return
			}

			// For other errors, return HTTP 500 Internal Server Error.
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Return HTTP 204 No Content for successful deletion.
		// HTTP 204 indicates success but no content to return.
		c.Status(http.StatusNoContent)
	}
}
