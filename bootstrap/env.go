package bootstrap

import (
	"fmt"
	"github.com/gookit/slog"
	"github.com/spf13/viper"
	"os"
)

type Env struct {
	Port int `mapstructure:"PORT"`

	DBName   string `mapstructure:"DB_NAME"`
	DBPort   int    `mapstructure:"DB_PORT"`
	MongoDir string `mapstructure:"MONGO_DIR"`

	MinioPort         int    `mapstructure:"MINIO_PORT"`
	MinioDir          string `mapstructure:"MINIO_DIR"`
	MinioRootUser     string `mapstructure:"MINIO_ROOT_USER"`
	MinioRootPassword string `mapstructure:"MINIO_ROOT_PASSWORD"`

	ContextTimeout         int    `mapstructure:"CONTEXT_TIMEOUT"`
	AccessTokenSecret      string `mapstructure:"ACCESS_TOKEN_SECRET"`
	AccessTokenExpiryHour  int    `mapstructure:"ACCESS_TOKEN_EXPIRY_HOUR"`
	RefreshTokenSecret     string `mapstructure:"REFRESH_TOKEN_SECRET"`
	RefreshTokenExpiryHour int    `mapstructure:"REFRESH_TOKEN_EXPIRY_HOUR"`
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
	slog.Info(fmt.Sprintf("The T-prep is running in %s env", os.Getenv("APP_ENV")))

	return &env
}
