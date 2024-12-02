package route

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gookit/slog"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"main/bootstrap"
	"main/database"
	"main/models"
	"net/http"
	"time"
)

var Env *bootstrap.Env
var collections, collectionsGlobal *mongo.Collection
var globalValues models.GlobalValues

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
	client, Env := app.Mongo, app.Env
	collections, collectionsGlobal = database.InitDatabase(client.Database(Env.DBName))

	r.Use(loggingMiddleware)
	initCardRouter(r)
	initCollectionRouter(r)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ping")) })
}
