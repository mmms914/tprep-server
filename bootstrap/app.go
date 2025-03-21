package bootstrap

import (
	"main/database"
	"main/storage"
)

type Application struct {
	Env     *Env
	Mongo   database.Client
	Storage storage.Client
}

func App() Application {
	app := &Application{}
	app.Env = NewEnv()
	app.Mongo = NewMongoDatabase(app.Env)
	app.Storage = NewStorage(app.Env)
	return *app
}

func (app *Application) CloseDBConnection() {
	CloseMongoDBConnection(app.Mongo)
}
