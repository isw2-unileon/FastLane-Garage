// Package handlers contains HTTP request handlers for the API endpoints.
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/dto"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/service"
)

// GetParts returns a Gin handler function that fetches all parts from the database.
// It accepts a PartsService and returns a handler that can be registered
// to a Gin route. This pattern (handler factory) allows dependency injection of the service.
func GetParts(svc service.PartsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract optional query parameters for filtering.
		zone := c.Query("zone")
		name := c.Query("name")

		// Call the service to get all parts.
		parts, err := svc.GetAllParts(zone, name)
		if err != nil {
			handleInternalError(c, "failed to fetch parts")
			return
		}

		// Return the parts as JSON with HTTP 200 OK status.
		c.JSON(http.StatusOK, gin.H{
			"data": parts,
		})
	}
}

// GetPartByID returns a Gin handler function that fetches a single part by ID.
// The part ID is extracted from the URL path parameter (e.g., /api/parts/1).
func GetPartByID(svc service.PartsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract and validate the part ID from the URL path.
		id, ok := parseID(c)
		if !ok {
			return
		}

		// Call the service to get the part by ID.
		part, err := svc.GetPartByID(id)
		if err != nil {
			// Check if the error is "record not found".
			if handleNotFoundError(c, err, "part") {
				return
			}

			handleInternalError(c, "failed to fetch part")
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

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Call the service to create the part.
		part, err := svc.CreatePart(&req)
		if err != nil {
			handleInternalError(c, err.Error())
			return
		}

		// Return the created part with HTTP 201 Created status.
		c.JSON(http.StatusCreated, gin.H{
			"data": part,
		})
	}
}

// UpdatePart returns a Gin handler function that updates an existing part.
// The part ID is extracted from the URL path, and the new data is in the request body.
func UpdatePart(svc service.PartsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract and validate the part ID from the URL path.
		id, ok := parseID(c)
		if !ok {
			return
		}

		// Parse and validate the request body.
		var req dto.UpdatePartRequest

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Call the service to update the part.
		part, err := svc.UpdatePart(id, &req)
		if err != nil {
			if handleNotFoundError(c, err, "part") {
				return
			}

			handleInternalError(c, err.Error())
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
		// Extract and validate the part ID from the URL path.
		id, ok := parseID(c)
		if !ok {
			return
		}

		// Call the service to delete the part.
		if err := svc.DeletePart(id); err != nil {
			if handleNotFoundError(c, err, "part") {
				return
			}

			handleInternalError(c, "failed to delete part")
			return
		}

		// Return HTTP 204 No Content for successful deletion.
		c.Status(http.StatusNoContent)
	}
}
