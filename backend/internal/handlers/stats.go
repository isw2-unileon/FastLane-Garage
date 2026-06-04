// Package handlers contains HTTP request handlers for the API endpoints.
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/service"
)

// GetTopPartsStats
func GetTopPartsStats(ordersSvc service.OrdersService) gin.HandlerFunc {
	return func(c *gin.Context) {
		limitStr := c.DefaultQuery("limit", "5")
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			limit = 5
		}

		stats, err := ordersSvc.GetTopPartsStats(limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, stats)
	}
}
