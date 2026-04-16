package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ========== ADDED: Prometheus HTTP Metrics ==========
var (
	// Counter: Total requests by method, path, status
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// Histogram: Request duration by method, path
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path"},
	)

	// Gauge: Current requests in flight
	httpRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of HTTP requests being processed",
		},
	)
)

// ========== ADDED: Gin Middleware for Instrumentation ==========
func prometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		httpRequestsInFlight.Inc()

		// Process request
		c.Next()

		// Record metrics after request completes
		httpRequestsInFlight.Dec()
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		// Use path template (e.g., /users/:id) not actual path (e.g., /users/123)
		// Prevents metric explosion from high cardinality
		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}

		httpRequestsTotal.WithLabelValues(c.Request.Method, path, status).Inc()
		httpRequestDuration.WithLabelValues(c.Request.Method, path).Observe(duration)
	}
}

func main() {
	// ========== CHANGED: Use Gin with middleware ==========
	r := gin.New()                    // REMOVED: gin.Default() (includes Logger/Recovery)
	r.Use(gin.Recovery())             // ADDED: Recovery only (we'll add custom logging via middleware)
	r.Use(prometheusMiddleware())     // ADDED: Prometheus instrumentation

	// ========== CHANGED: Metrics endpoint (no middleware) ==========
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// ========== SAME: Health endpoints (now auto-instrumented) ==========
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "api",
		})
	})

	r.GET("/ready", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"ready": true,
		})
	})

	// ========== ADDED: Example API endpoint (auto-instrumented) ==========
	r.GET("/api/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{
			"user_id": id,
			"name":    "User " + id,
		})
	})

	// ========== SAME: Port configuration ==========
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}