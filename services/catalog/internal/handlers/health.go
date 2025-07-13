package handlers

import (
	"net/http"
	"time"

	"catalog-service/internal/db"

	"github.com/gin-gonic/gin"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Service   string            `json:"service"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version"`
	Checks    map[string]string `json:"checks"`
}

// HealthCheck handles GET /health requests (simple version)
func HealthCheck(c *gin.Context) {
	response := HealthResponse{
		Status:    "healthy",
		Service:   "catalog",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Checks:    map[string]string{},
	}

	c.JSON(http.StatusOK, response)
}

// HealthCheckWithDB handles GET /health requests with database connectivity check
func HealthCheckWithDB(database *db.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		checks := make(map[string]string)
		status := "healthy"
		httpStatus := http.StatusOK

		// Check database connectivity
		if err := database.HealthCheck(); err != nil {
			checks["database"] = "unhealthy: " + err.Error()
			status = "unhealthy"
			httpStatus = http.StatusServiceUnavailable
		} else {
			checks["database"] = "healthy"
		}

		response := HealthResponse{
			Status:    status,
			Service:   "catalog",
			Timestamp: time.Now(),
			Version:   "1.0.0",
			Checks:    checks,
		}

		c.JSON(httpStatus, response)
	}
}
