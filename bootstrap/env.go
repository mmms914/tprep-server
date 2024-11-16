package bootstrap

import (
	"github.com/gookit/slog"
	"github.com/spf13/viper"
)

type Env struct {
	AppEnv   string `mapstructure:"APP_ENV"`
	Port     int    `mapstructure:"PORT"`
	MongoURI string `mapstructure:"MONGO_URI"`
	DBName   string `mapstructure:"DB_NAME"`
	DBPort   int    `mapstructure:"DB_PORT"`
}

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		slog.Fatal("Can't find the config file", err)
	}
	err = viper.Unmarshal(&env)
	if err != nil {
		slog.Fatal("Environment can't be loaded", err)
	}
	if env.AppEnv == "admin" {
		slog.Info("The T-prep is running in admin env")
	}
	return &env
}
