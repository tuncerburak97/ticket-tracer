package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

var (
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"path", "method"},
	)

	httpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)
)

func init() {
	prometheus.MustRegister(httpRequests, httpDuration)
}

func promMiddleware(ctx *fiber.Ctx) error {
	start := time.Now()

	err := ctx.Next()

	duration := time.Since(start).Seconds()

	httpRequests.WithLabelValues(ctx.Path(), ctx.Method()).Inc()
	httpDuration.WithLabelValues(ctx.Path()).Observe(duration)

	return err
}
