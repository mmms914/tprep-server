package middleware

import (
	mapset "github.com/deckarep/golang-set"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

var uniqueUsers mapset.Set
var lastSync time.Time

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w, http.StatusOK}
}

func (rw *ResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of get requests.",
	},
	[]string{"path"},
)

var FavouriteButtonClicks = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "favourite_button_clicks_total",
		Help: "The number of clicks on the 'Favourites' buttons, by type",
	},
	[]string{"button_type"},
)

var UserTime = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "user_time_seconds",
		Help:    "The time that user spent, by type",
		Buckets: []float64{60, 180, 300, 600, 1800},
	},
	[]string{"time_type"},
)

var TrainingsCount = prometheus.NewHistogram(
	prometheus.HistogramOpts{
		Name:    "user_trainings_count",
		Help:    "The count of trainings in last using",
		Buckets: []float64{0, 1, 2, 3, 4, 5, 10},
	},
)

var UsersCount = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "unique_users_count",
		Help: "The count of unique users in that day",
	},
)

var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "response_status_total",
		Help: "Status of HTTP response",
	},
	[]string{"status"},
)

var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_response_time_seconds",
	Help: "Duration of HTTP requests.",
}, []string{"path"})

func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := chi.RouteContext(r.Context())
		path := route.RoutePattern()
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)

		statusCode := rw.statusCode

		responseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
		totalRequests.WithLabelValues(path).Inc()

		timer.ObserveDuration()

		checkUniqueUsers()
	})
}

func addPrometheusUser(userID string) {
	uniqueUsers.Add(userID)
	UsersCount.Set(float64(uniqueUsers.Cardinality()))
}

func checkUniqueUsers() {
	if uniqueUsers == nil {
		uniqueUsers = mapset.NewSet()
	}

	if lastSync.Day() != time.Now().Day() {
		uniqueUsers.Clear()
	}
	lastSync = time.Now()
}

//nolint:gochecknoinits // middleware
func init() {
	prometheus.Register(FavouriteButtonClicks)
	prometheus.Register(UserTime)
	prometheus.Register(TrainingsCount)
	prometheus.Register(UsersCount)
	prometheus.Register(totalRequests)
	prometheus.Register(responseStatus)
	prometheus.Register(httpDuration)
}
