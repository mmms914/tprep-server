package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
	"net/http"
	"strconv"
)

type Card struct {
	LocalID  int    `json:"local_id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type Collection struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	IsPublic bool   `json:"is_public"`
	Cards    []Card `json:"cards"`
}

var client *mongo.Client

func getCollection(w http.ResponseWriter, r *http.Request) {
	var result Collection
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	filter := bson.D{{"id", id}}
	collections := client.Database("tprep").Collection("collections")
	err := collections.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			w.WriteHeader(404)
			w.Write([]byte("collection not found"))
		} else {
			log.Fatal(err)
		}
		return
	}
	s, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte(s))
}

func main() {
	var mongoURI = "mongodb://localhost:27017"
	var err error
	client, err = mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		fmt.Println(err)
	}
	r := chi.NewRouter()
	r.Get("/getCollection/{id}", getCollection)
	fmt.Println("Listening on port 3000")
	http.ListenAndServe(":3000", r)

}
