package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// PrometheusCollector represents a Prometheus metrics collector
type PrometheusCollector struct {
	serviceName string
	registry    *prometheus.Registry
}

// NewPrometheusCollector creates a new Prometheus metrics collector
func NewPrometheusCollector(serviceName string) *PrometheusCollector {
	registry := prometheus.NewRegistry()
	
	// Register default metrics
	registry.MustRegister(
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		prometheus.NewGoCollector(),
	)

	return &PrometheusCollector{
		serviceName: serviceName,
		registry:    registry,
	}
}

// Handler returns the HTTP handler for metrics endpoint
func (c *PrometheusCollector) Handler() http.Handler {
	return promhttp.HandlerFor(c.registry, promhttp.HandlerOpts{})
}

// Registry returns the Prometheus registry
func (c *PrometheusCollector) Registry() *prometheus.Registry {
	return c.registry
}
