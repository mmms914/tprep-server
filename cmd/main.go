package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"main/bootstrap"
	"main/route"
	"net/http"
)

func main() {
	app := bootstrap.App()
	env := app.Env
	r := chi.NewRouter()
	fmt.Println("Listening on port 3000")
	route.Setup(r, app)
	http.ListenAndServe(env.ServerAddress, r)
}
