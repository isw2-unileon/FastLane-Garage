// Package handlers contains HTTP request handlers for the API endpoints.
package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/dto"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/repository"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/service"
)

// GetOrders returns a Gin handler function that fetches all orders.
// It accepts an OrdersService and returns a handler that can be registered to a Gin route.
func GetOrders(svc service.OrdersService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Call the service to get all orders.
		orders, err := svc.GetAllOrders()
		if err != nil {
			// If an error occurs, respond with HTTP 500 Interval Server Error.
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to fetch orders",
			})
			return
		}

		// Returns the orders as JSON with HTTP 200 OK status.
		c.JSON(http.StatusOK, gin.H{
			"data": orders,
		})
	}
}

// GetOrderByID returns a Gin handler function that fetches a single order by ID.
// The order ID os extracted form the URL path parameter (e.g., /api/orders/1).
func GetOrderByID(svc service.OrdersService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the order ID from the URL path.
		idParam := c.Param("id")

		// Convert the ID string to an unsigned integer.
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			// If the ID is not a valid number, return HTTP 400 Bad Request.
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid order ID format",
			})
			return
		}

		// Call the service to get the order by ID.
		order, err := svc.GetOrderByID(uint(id))
		if err != nil {
			// Check if the error is "record not found".
			if errors.Is(err, repository.ErrNotFound) {
				// If the order is not found, return HTTP 404 Not Found.
				c.JSON(http.StatusNotFound, gin.H{
					"error": "order not found",
				})
				return
			}

			// For other errors, return HTTP 500 Internal Server Error.
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to fetch order",
			})
			return
		}

		// Return the order as JSON with HTTP 200 OK status.
		c.JSON(http.StatusOK, gin.H{
			"data": order,
		})
	}
}

// CreateOrder returns a Gin handler function that creates a new order.
// The request body ir expected to be JSON matching CreateOrderRequest.
func CreateOrder(svc service.OrdersService, partsRepo repository.PartsRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse and validate the request body.
		var req dto.CreateOrderRequest

		// BindJSON parses the JSON and validates it using the binding tags defined in the DTO.
		if err := c.BindJSON(&req); err != nil {
			// If binding fails, return HTTP 400 Bad Request with the error message.
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Call the service to create the order, passing the parts repository for validation.
		order, err := svc.CreateOrder(&req, partsRepo)
		if err != nil {
			// If service returns an error, return HTTP 500 Interval Server Error.
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Return the created order with HTTP 201 Created status.
		// HTTP 201 indicates that a new resource has been created.
		c.JSON(http.StatusCreated, gin.H{
			"data": order,
		})
	}
}

// UpdateOrderStatus returns a Gin handler function that updates an orders's status.
// The order ID is extracted from the URL path, and the new status is in the request body.
func UpdateOrderStatus(svc service.OrdersService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the order ID from the URL path.
		idParam := c.Param("id")

		// Convert the ID string to an unsigned integer.
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			// If the ID is not a valid number, return HTTP 400 Bad Request.
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid order ID format",
			})
			return
		}

		// Parse and validate the request body.
		var req dto.UpdateOrderStatusRequest

		// BindJSON parses the JSON and validates it using the binding tags.
		if err := c.BindJSON(&req); err != nil {
			// If binding fails, return HTTP 400 Bad Request.
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Call the service to update the order status.
		order, err := svc.UpdateOrderStatus(uint(id), &req)
		if err != nil {
			// Check if the error is "record not found".
			if errors.Is(err, repository.ErrNotFound) {
				// If the order is not found, return HTTP 404 Not Found.
				c.JSON(http.StatusNotFound, gin.H{
					"error": "order not found",
				})
				return
			}

			// For other errors, return HTTP 500 Interval Server Error.
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Return the updated order with HTTP 200 OK status.
		c.JSON(http.StatusOK, gin.H{
			"data": order,
		})
	}
}

// DeleteOrder returns a Gin handler function that deletes an order by ID.
// The order ID is extracted from the URL path parameter.
func DeleteOrder(svc service.OrdersService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the order ID from the URL path.
		idParam := c.Param("id")

		// Convert the ID string to an unsigned integer.
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			// If the ID is not a valid number, return HTTP 400 Bad Request.
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid order ID format",
			})
			return
		}

		// Call the service to delete the order.
		if err := svc.DeleteOrder(uint(id)); err != nil {
			// Check if the error is "record not found".
			if errors.Is(err, repository.ErrNotFound) {
				// If the order is not found, return HTTP 404 Not Found.
				c.JSON(http.StatusNotFound, gin.H{
					"error": "order not found",
				})
				return
			}

			// For other errors, return HTTP 500 Interval Server Error.
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
