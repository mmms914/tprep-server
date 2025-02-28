package route

import (
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"main/api/middleware"
	"main/bootstrap"
	"main/database"
	"net/http"
	"time"
)

func Setup(env *bootstrap.Env, timeout time.Duration, db database.Database, r *chi.Mux) {
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CORSHandler)
	r.Use(middleware.PrometheusMiddleware)
	// public methods
	r.Group(func(r chi.Router) {
		NewSignupRouter(env, timeout, db, r)
		NewLoginRouter(env, timeout, db, r)
		NewRefreshTokenRouter(env, timeout, db, r)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ping")) })
		r.Handle("/metrics", promhttp.Handler())
	})

	// private methods
	r.Group(func(r chi.Router) {
		r.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret))
		NewCollectionRouter(env, timeout, db, r)
		NewUserRouter(env, timeout, db, r)
		NewGlobalRouter(env, timeout, r)
	})
}
