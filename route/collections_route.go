package route

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/gookit/slog"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"main/models"
	"net/http"
	"strconv"
)

var collections, collectionsGlobal *mongo.Collection

type CollectionInfo struct {
	Id    int `bson:"id"`
	MaxId int `bson:"max_id"`
}

func getCollection(w http.ResponseWriter, r *http.Request) {
	var result models.Collection
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	filter := bson.D{{Key: "id", Value: id}}
	err := collections.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			w.WriteHeader(404)
			w.Write([]byte("collection not found"))
		} else {
			slog.FatalErr(err)
		}
		return
	}
	s, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(s)
}

func initCollectionRouter(r *chi.Mux) {
	collections = Client.Database(Env.DBName).Collection("collections")
	collectionsGlobal = Client.Database(Env.DBName).Collection("collections_global")
	r.Get("/collection/{id}", getCollection)
}
