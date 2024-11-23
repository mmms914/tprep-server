package route

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gookit/slog"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"main/bootstrap"
	"main/models"
	"net/http"
	"time"
)

var Client *mongo.Client
var Env *bootstrap.Env
var globals *mongo.Collection
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

func getGlobalValues() {
	err := globals.FindOne(context.TODO(), bson.D{}).Decode(&globalValues)
	if err != nil {
		slog.FatalErr(err)
	}
}
func updateGlobalValues() {
	_, err := globals.ReplaceOne(context.TODO(), bson.D{}, globalValues)
	if err != nil {
		slog.FatalErr(err)
	}
}

func Setup(r *chi.Mux, app bootstrap.Application) {
	Client, Env = app.Mongo, app.Env
	globals = Client.Database(Env.DBName).Collection("globals") // ???
	r.Use(loggingMiddleware)
	initCardRouter(r)
	initCollectionRouter(r)
	getGlobalValues()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ping")) })
}
