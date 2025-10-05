package monitoring

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"strconv"
	"time"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"method", "path", "status"})

	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request duration in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path", "status"})

	httpRequestInFlight = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "http_request_in_flight",
		Help: "Current number of HTTP requests in flight",
	})

	websocketConnections = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "websocket_connections_total",
		Help: "Current number of WebSocket connections",
	})

	kafkaMessageSent = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "kafka_messages_sent_total",
		Help: "Total number of Kafka messages sent",
	}, []string{"topic"})

	redisOperations = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "redis_operations_total",
		Help: "Total number of Redis operations",
	}, []string{"operation", "status"})
)

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()

		if path == "/metrics" || path == "/debug/pprof/*pprof" {
			c.Next()
			return
		}

		httpRequestInFlight.Inc()
		defer httpRequestInFlight.Dec()

		c.Next()

		status := strconv.Itoa(c.Writer.Status())
		httpRequestsTotal.WithLabelValues(c.Request.Method, path, status).Inc()

		httpRequestDuration.WithLabelValues(c.Request.Method, path, status).Observe(time.Since(start).Seconds())
	}
}

func IncrementWebSocketConnections() {
	websocketConnections.Inc()
}

func DecrementWebSocketConnections() {
	websocketConnections.Dec()
}

func IncrementKafkaMessagesSent(topic string) {
	kafkaMessageSent.WithLabelValues().Inc()
}

func IncrementRedisOperations(operation, status string) {
	redisOperations.WithLabelValues(operation, status).Inc()
}
