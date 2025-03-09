package bootstrap

import (
	"fmt"
	"github.com/gookit/slog"
	"github.com/spf13/viper"
)

type Env struct {
	AppEnv string `mapstructure:"APP_ENV"`
	Port   int    `mapstructure:"PORT"`

	LocalMongoURI  string `mapstructure:"LOCAL_MONGO_URI"`
	DockerMongoURI string `mapstructure:"DOCKER_MONGO_URI"`
	DBName         string `mapstructure:"DB_NAME"`
	DBPort         int    `mapstructure:"DB_PORT"`

	MinioPort      int    `mapstructure:"MINIO_PORT"`
	MongoDir       string `mapstructure:"MONGO_DIR"`
	LocalMinioURI  string `mapstructure:"LOCAL_MINIO_URI"`
	DockerMinioURI string `mapstructure:"DOCKER_MINIO_URI"`
	MinioAccessKey string `mapstructure:"MINIO_ACCESS_KEY"`
	MinioSecretKey string `mapstructure:"MINIO_SECRET_KEY"`

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
	slog.Info(fmt.Sprintf("The T-prep is running in %s env", env.AppEnv))

	return &env
}
