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
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("collection not found"))
		} else {
			slog.FatalErr(err)
		}
		return
	}
	s, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(s)
}

func createCollection(w http.ResponseWriter, r *http.Request) {
	var newCollection models.Collection
	err := json.NewDecoder(r.Body).Decode(&newCollection)
	if err != nil || newCollection.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request payload or collection name"))
		return
	}
	collections := Client.Database(Env.DBName).Collection("collections")
	newCollection.ID = globalValues.MaxCollectionID + 1
	globalValues.MaxCollectionID += 1
	updateGlobalValues()
	_, err = collections.InsertOne(context.TODO(), newCollection)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	response, _ := json.Marshal(newCollection)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func editCollection(w http.ResponseWriter, r *http.Request) {
	var editedCollection, existingCollection models.Collection
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	err := json.NewDecoder(r.Body).Decode(&editedCollection)
	if err != nil || editedCollection.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request payload or collection name"))
		return
	}
	collections := Client.Database(Env.DBName).Collection("collections")
	filter := bson.D{{"id", id}}
	update := bson.D{
		{"$set", bson.D{
			{"name", editedCollection.Name},
			{"is_public", editedCollection.IsPublic},
		}},
	}
	editedCollection = existingCollection
	err = collections.FindOneAndUpdate(context.TODO(), filter, update).Decode(&existingCollection)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("collection not found"))
		return
	}
	response, _ := json.Marshal(editedCollection)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func deleteCollection(w http.ResponseWriter, r *http.Request) {
	var result models.Collection
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	filter := bson.D{{Key: "id", Value: id}}
	err := collections.FindOneAndDelete(context.TODO(), filter).Decode(&result)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("collection not found"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("collection successfully deleted"))
}

func initCollectionRouter(r *chi.Mux) {
	collections = Client.Database(Env.DBName).Collection("collections")
	collectionsGlobal = Client.Database(Env.DBName).Collection("collections_global")
	r.Get("/collection/{id}", getCollection)
	r.Post("/collection", createCollection)
	r.Put("/collection/{id}", editCollection)
	r.Delete("/collection/{id}", deleteCollection)
}
