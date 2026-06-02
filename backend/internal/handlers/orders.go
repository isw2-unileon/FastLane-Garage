// Package handlers contains HTTP request handlers for the API endpoints.
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/dto"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/repository"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/service"
)

// GetOrders returns a Gin handler function that fetches all orders.
// It accepts an OrdersService and returns a handler that can be registered to a Gin route.
func GetOrders(svc service.OrdersService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract optional query parameters for filtering.
		status := c.Query("status")
		email := c.Query("email")

		// Call the service to get all orders.
		orders, err := svc.GetAllOrders(status, email)
		if err != nil {
			handleInternalError(c, "failed to fetch orders")
			return
		}

		// Return the orders as JSON with HTTP 200 OK status.
		c.JSON(http.StatusOK, gin.H{
			"data": orders,
		})
	}
}

// GetOrderByID returns a Gin handler function that fetches a single order by ID.
// The order ID is extracted from the URL path parameter (e.g., /api/orders/1).
func GetOrderByID(svc service.OrdersService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract and validate the order ID from the URL path.
		id, ok := parseID(c)
		if !ok {
			return
		}

		// Call the service to get the order by ID.
		order, err := svc.GetOrderByID(id)
		if err != nil {
			if handleNotFoundError(c, err, "order") {
				return
			}

			handleInternalError(c, "failed to fetch order")
			return
		}

		// Return the order as JSON with HTTP 200 OK status.
		c.JSON(http.StatusOK, gin.H{
			"data": order,
		})
	}
}

// CreateOrder returns a Gin handler function that creates a new order.
// The request body is expected to be JSON matching CreateOrderRequest.
func CreateOrder(svc service.OrdersService, partsRepo repository.PartsRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse and validate the request body.
		var req dto.CreateOrderRequest

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Call the service to create the order.
		order, err := svc.CreateOrder(&req, partsRepo)
		if err != nil {
			handleInternalError(c, err.Error())
			return
		}

		// Return the created order with HTTP 201 Created status.
		c.JSON(http.StatusCreated, gin.H{
			"data": order,
		})
	}
}

// UpdateOrderStatus returns a Gin handler function that updates an order's status.
// The order ID is extracted from the URL path, and the new status is in the request body.
func UpdateOrderStatus(svc service.OrdersService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract and validate the order ID from the URL path.
		id, ok := parseID(c)
		if !ok {
			return
		}

		// Parse and validate the request body.
		var req dto.UpdateOrderStatusRequest

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Call the service to update the order status.
		order, err := svc.UpdateOrderStatus(id, &req)
		if err != nil {
			if handleNotFoundError(c, err, "order") {
				return
			}

			handleInternalError(c, "failed to update order status")
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
		// Extract and validate the order ID from the URL path.
		id, ok := parseID(c)
		if !ok {
			return
		}

		// Call the service to delete the order.
		if err := svc.DeleteOrder(id); err != nil {
			if handleNotFoundError(c, err, "order") {
				return
			}

			handleInternalError(c, "failed to delete order")
			return
		}

		// Return HTTP 204 No Content for successful deletion.
		c.Status(http.StatusNoContent)
	}
}
