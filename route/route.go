package route

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gookit/slog"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"main/bootstrap"
	"net/http"
	"time"
)

var Client *mongo.Client
var Env *bootstrap.Env

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		slog.Debugf("%s %s %d %dB in %v", r.Method, r.URL.Path,
			ww.Status(), ww.BytesWritten(), time.Since(start))
	})
}

func Setup(r *chi.Mux, app bootstrap.Application) {
	Client, Env = app.Mongo, app.Env

	r.Use(loggingMiddleware)

	r.Get("/getCollection/{id}", getCollection)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ping")) })
}
