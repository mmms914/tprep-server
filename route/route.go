package route

import (
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"main/bootstrap"
)

var Client *mongo.Client

func Setup(r *chi.Mux, app bootstrap.Application) {
	Client = app.Mongo
	r.Get("/getCollection/{id}", getCollection)
}
