package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/gookit/slog"
	"main/api/route"
	"main/bootstrap"
	"net/http"
	"os"
)

func main() {
	slog.SetExitFunc(os.Exit)

	app := bootstrap.App()
	env := app.Env

	r := chi.NewRouter()
	route.Setup(r, app)

	slog.Infof("Listening on port %d", env.Port)
	slog.FatalErr(http.ListenAndServe(fmt.Sprintf(":%d", env.Port), r))
}
