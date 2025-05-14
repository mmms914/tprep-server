package main

import (
	"fmt"
	"main/api/route"
	"main/bootstrap"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gookit/slog"
)

func main() {
	slog.SetExitFunc(os.Exit)

	app := bootstrap.App()
	env := app.Env

	db := app.Mongo.Database(env.DBName)
	defer app.CloseDBConnection()

	s3 := app.Storage

	timeout := time.Duration(env.ContextTimeout) * time.Second

	r := chi.NewRouter()

	route.Setup(env, timeout, db, s3, r)

	slog.Infof("Listening on port %d", env.Port)
	slog.FatalErr(http.ListenAndServe(fmt.Sprintf(":%d", env.Port), r))
}
