package handlers

import (
	"database/sql"
	"net/http"

	"catalog-service/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	db *sql.DB
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

// HealthCheck handles GET /health
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	// Check database connection
	if h.db != nil {
		err := h.db.Ping()
		if err != nil {
			logger.WithError(err).WithFields(logrus.Fields{
				"component": "health",
				"action":    "check",
				"database":  "disconnected",
			}).Error("Database health check failed")

			c.JSON(http.StatusServiceUnavailable, gin.H{
				"data": gin.H{
					"status":   "unhealthy",
					"database": "disconnected",
					"service":  "catalog-service",
				},
				"error": err.Error(),
			})
			return
		}
	}

	logger.WithFields(logrus.Fields{
		"component": "health",
		"action":    "check",
		"status":    "healthy",
		"database":  "connected",
	}).Info("Health check successful")

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"status":   "healthy",
			"database": "connected",
			"service":  "catalog-service",
		},
	})
}
