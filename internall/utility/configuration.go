package utility

import (
	"github.com/IosifSuzuki/todo/internall/logger"
	"github.com/IosifSuzuki/todo/internall/model"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
)

const (
	keyPostgresUser     = "POSTGRES_USER"
	keyPostgresPassword = "POSTGRES_USER"
	keyPostgresDb       = "POSTGRES_DB"
	keyDatabaseHost     = "DATABASE_HOST"
	keySecretKey        = "SECRET_KEY"
)

type Configuration struct {
	DB        model.DBConfig
	SecretKey string
}

var Config Configuration

func Setup() {
	if err := godotenv.Load(".env"); err != nil {
		logger.Fatal("Couldn't load env file", zap.Error(err))
	}
	var (
		postgresUser     = os.Getenv(keyPostgresUser)
		postgresPassword = os.Getenv(keyPostgresPassword)
		postgresDB       = os.Getenv(keyPostgresDb)
		databaseHost     = os.Getenv(keyDatabaseHost)
		secretKey        = os.Getenv(keySecretKey)
	)
	var dbConfig = model.DBConfig{
		UserName: postgresUser,
		Password: postgresPassword,
		DBName:   postgresDB,
		DBHost:   databaseHost,
	}
	Config = Configuration{
		DB:        dbConfig,
		SecretKey: secretKey,
	}
}
