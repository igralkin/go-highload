package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Две метрики из задания: общий счётчик запросов и гистограмма времени ответа.
var (
	TotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint"},
	)

	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Request duration in seconds",
		},
		[]string{"method", "endpoint"},
	)
)

// Регистрация метрик в Prometheus default registry.
func init() {
	prometheus.MustRegister(TotalRequests)
	prometheus.MustRegister(RequestDuration)
}

// MetricsMiddleware считает RPS и latency по методу и пути.
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Увеличиваем счётчик запросов
		TotalRequests.WithLabelValues(r.Method, r.URL.Path).Inc()

		// Отдаём запрос дальше по цепочке
		next.ServeHTTP(w, r)

		// Фиксируем задержку
		elapsed := time.Since(start).Seconds()
		RequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(elapsed)
	})
}

// Handler возвращает http.Handler для /metrics.
func Handler() http.Handler {
	return promhttp.Handler()
}
