package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

// Config ...
type Config struct {
	App         string
	AppVersion  string
	Environment string // development, staging, production

	GRPCPort string

	DefaultOffset string
	DefaultLimit  string

	CarServiceGrpcHost string
	CarServiceGrpcPort string

	RentalServiceGrpcHost string
	RentalServiceGrpcPort string

	AuthorizationServiceGrpcHost string
	AuthorizationServiceGrpcPort string

	PostgresHost     string
	PostgresPort     int
	PostgresDatabase string
	PostgresUser     string
	PostgresPassword string
}

// Load ...
func Load() Config {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	config := Config{}
	config.App = cast.ToString(getOrReturnDefaultValue("APP", "rental_service_API"))
	config.AppVersion = cast.ToString(getOrReturnDefaultValue("APP_VERSION", "1.0.1"))
	config.Environment = cast.ToString(getOrReturnDefaultValue("ENVIRONMENT", "development"))

	config.GRPCPort = cast.ToString(getOrReturnDefaultValue("GRPC_PORT", ":7003"))

	config.DefaultOffset = cast.ToString(getOrReturnDefaultValue("DEFAULT_OFFSET", "0"))
	config.DefaultLimit = cast.ToString(getOrReturnDefaultValue("DEFAULT_LIMIT", "10"))

	config.CarServiceGrpcHost = cast.ToString(getOrReturnDefaultValue("CAR_SERVICE_GRPC_HOST", "localhost"))
	config.CarServiceGrpcPort = cast.ToString(getOrReturnDefaultValue("CAR_SERVICE_GRPC_PORT", ":7001"))

	config.RentalServiceGrpcHost = cast.ToString(getOrReturnDefaultValue("RENTAL_SERVICE_GRPC_HOST", "localhost"))
	config.RentalServiceGrpcPort = cast.ToString(getOrReturnDefaultValue("RENTAL_SERVICE_GRPC_PORT", ":7003"))

	config.AuthorizationServiceGrpcHost = cast.ToString(getOrReturnDefaultValue("AUTHORIZATION_SERVICE_GRPC_HOST", "localhost"))
	config.AuthorizationServiceGrpcPort = cast.ToString(getOrReturnDefaultValue("AUTHORIZATION_SERVICE_GRPC_PORT", ":7002"))

	config.PostgresHost = cast.ToString(getOrReturnDefaultValue("POSTGRES_HOST", "127.0.0.1"))
	config.PostgresPort = cast.ToInt(getOrReturnDefaultValue("POSTGRES_PORT", 5432))
	config.PostgresDatabase = cast.ToString(getOrReturnDefaultValue("POSTGRES_DATABASE", "rental_db"))
	config.PostgresUser = cast.ToString(getOrReturnDefaultValue("POSTGRES_USER", "rental_db_user"))
	config.PostgresPassword = cast.ToString(getOrReturnDefaultValue("POSTGRES_PASSWORD", "rental_db_password"))

	return config
}

func getOrReturnDefaultValue(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)

	if exists {
		return os.Getenv(key)
	}

	return defaultValue
}
