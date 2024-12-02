package route

import (
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"main/api/middleware"
	"main/bootstrap"
	"main/database"
	"main/models"
	"net/http"
)

// var Env *bootstrap.Env
var collections, collectionsGlobal *mongo.Collection
var globalValues models.GlobalValues

func Setup(r *chi.Mux, app bootstrap.Application) {
	client, Env := app.Mongo, app.Env
	collections, collectionsGlobal = database.InitDatabase(client.Database(Env.DBName))

	r.Use(middleware.LoggingMiddleware)
	initCardRouter(r)
	initCollectionRouter(r)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ping")) })
}
