package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

/*
This whole handler is for an experimental metrics handling for the frontend.
TL;DR is that the frontend sends its own browser metrics to the catalog service, which it will then
expose at /metrics so prometheus can scrape it.
*/

// Frontend metrics data structures
type FrontendMetric struct {
	Name      string            `json:"name"`
	Value     float64           `json:"value"`
	Timestamp int64             `json:"timestamp"`
	Labels    map[string]string `json:"labels"`
}

type BusinessEvent struct {
	Event      string                 `json:"event"`
	Properties map[string]interface{} `json:"properties"`
	Timestamp  int64                  `json:"timestamp"`
}

type ErrorEvent struct {
	Error     string `json:"error"`
	Stack     string `json:"stack,omitempty"`
	Component string `json:"component,omitempty"`
	Timestamp int64  `json:"timestamp"`
	URL       string `json:"url"`
}

type FrontendMetricsPayload struct {
	PerformanceMetrics []FrontendMetric `json:"performance_metrics"`
	BusinessEvents     []BusinessEvent  `json:"business_events"`
	ErrorEvents        []ErrorEvent     `json:"error_events"`
	Timestamp          int64            `json:"timestamp"`
	SessionID          string           `json:"session_id"`
}

// Prometheus metrics for frontend data
var (
	// Performance metrics
	frontendLCPHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "frontend_lcp_seconds",
			Help: "Frontend Largest Contentful Paint timing",
		},
		[]string{"url", "user_agent_type", "session_id"},
	)

	frontendFIDHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "frontend_fid_milliseconds",
			Help: "Frontend First Input Delay timing",
		},
		[]string{"url", "user_agent_type", "session_id"},
	)

	frontendCLSHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "frontend_cls_score",
			Help: "Frontend Cumulative Layout Shift score",
		},
		[]string{"url", "user_agent_type", "session_id"},
	)

	frontendPageLoadHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "frontend_page_load_seconds",
			Help: "Frontend page load timing",
		},
		[]string{"url", "user_agent_type", "session_id"},
	)

	frontendAPICallHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "frontend_api_call_duration_seconds",
			Help: "Frontend API call duration from user perspective",
		},
		[]string{"endpoint", "success", "url", "user_agent_type", "session_id"},
	)

	frontendReactQueryHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "frontend_react_query_duration_seconds",
			Help: "Frontend React Query operation duration",
		},
		[]string{"query_key", "from_cache", "url", "user_agent_type", "session_id"},
	)

	// Business event counters
	frontendPageViewsCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "frontend_page_views_total",
			Help: "Total frontend page views",
		},
		[]string{"page", "url", "session_id"},
	)

	frontendProductViewsCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "frontend_product_views_total",
			Help: "Total frontend product views",
		},
		[]string{"product_id", "product_name", "url", "session_id"},
	)

	// Error tracking
	frontendErrorsCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "frontend_errors_total",
			Help: "Total frontend errors",
		},
		[]string{"error_type", "component", "url", "session_id"},
	)
)

// FrontendMetricsHandler handles frontend metrics
type FrontendMetricsHandler struct{}

// NewFrontendMetricsHandler creates a new frontend metrics handler
func NewFrontendMetricsHandler() *FrontendMetricsHandler {
	return &FrontendMetricsHandler{}
}

// Helper function to extract user agent type
func getUserAgentType(userAgent string) string {
	if len(userAgent) == 0 {
		return "unknown"
	}

	// Simple user agent classification
	switch {
	case len(userAgent) > 20 && userAgent[:20] == "Mozilla/5.0 (iPhone":
		return "mobile_ios"
	case len(userAgent) > 25 && userAgent[:25] == "Mozilla/5.0 (Macintosh;":
		return "desktop_mac"
	case len(userAgent) > 15 && userAgent[:15] == "Mozilla/5.0 (X1":
		return "desktop_linux"
	case len(userAgent) > 20 && userAgent[:20] == "Mozilla/5.0 (Windows":
		return "desktop_windows"
	default:
		return "other"
	}
}

// HandleFrontendMetrics processes frontend metrics and converts them to Prometheus format
func (h *FrontendMetricsHandler) HandleFrontendMetrics(c *gin.Context) {
	var payload FrontendMetricsPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		logrus.WithError(err).Error("Failed to decode frontend metrics payload")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	// Process performance metrics
	for _, metric := range payload.PerformanceMetrics {
		userAgentType := getUserAgentType(metric.Labels["user_agent"])
		url := metric.Labels["url"]

		switch metric.Name {
		case "frontend_lcp_seconds":
			frontendLCPHistogram.WithLabelValues(url, userAgentType, payload.SessionID).Observe(metric.Value)
		case "frontend_fid_milliseconds":
			frontendFIDHistogram.WithLabelValues(url, userAgentType, payload.SessionID).Observe(metric.Value)
		case "frontend_cls_score":
			frontendCLSHistogram.WithLabelValues(url, userAgentType, payload.SessionID).Observe(metric.Value)
		case "frontend_page_load_seconds":
			frontendPageLoadHistogram.WithLabelValues(url, userAgentType, payload.SessionID).Observe(metric.Value)
		case "frontend_api_call_duration_seconds":
			endpoint := metric.Labels["endpoint"]
			success := metric.Labels["success"]
			frontendAPICallHistogram.WithLabelValues(endpoint, success, url, userAgentType, payload.SessionID).Observe(metric.Value)
		case "frontend_react_query_duration_seconds":
			queryKey := metric.Labels["query_key"]
			fromCache := metric.Labels["from_cache"]
			frontendReactQueryHistogram.WithLabelValues(queryKey, fromCache, url, userAgentType, payload.SessionID).Observe(metric.Value)
		}
	}

	// Process business events
	for _, event := range payload.BusinessEvents {
		url := ""
		if urlProp, ok := event.Properties["url"].(string); ok {
			url = urlProp
		}

		switch event.Event {
		case "page_view":
			if page, ok := event.Properties["page"].(string); ok {
				frontendPageViewsCounter.WithLabelValues(page, url, payload.SessionID).Inc()
			}
		case "product_view":
			if productID, ok := event.Properties["product_id"].(string); ok {
				productName := ""
				if name, ok := event.Properties["product_name"].(string); ok {
					productName = name
				}
				frontendProductViewsCounter.WithLabelValues(productID, productName, url, payload.SessionID).Inc()
			}
		}
	}

	// Process error events
	for _, errorEvent := range payload.ErrorEvents {
		errorType := errorEvent.Error
		component := errorEvent.Component
		if component == "" {
			component = "unknown"
		}
		frontendErrorsCounter.WithLabelValues(errorType, component, errorEvent.URL, payload.SessionID).Inc()
	}

	// Log the metrics reception
	logrus.WithFields(logrus.Fields{
		"session_id":          payload.SessionID,
		"performance_metrics": len(payload.PerformanceMetrics),
		"business_events":     len(payload.BusinessEvents),
		"error_events":        len(payload.ErrorEvents),
		"timestamp":           time.Unix(payload.Timestamp/1000, 0),
	}).Info("Processed frontend metrics")

	// Return success
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Frontend metrics processed successfully",
		"metrics": gin.H{
			"performance": len(payload.PerformanceMetrics),
			"business":    len(payload.BusinessEvents),
			"errors":      len(payload.ErrorEvents),
		},
	})
}
