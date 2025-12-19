package config

import (
	"os"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
)

type secrets struct {
	PrimaryDatabaseHost     string `validate:"required"`
	PrimaryDatabasePort     string `validate:"required"`
	PrimaryDatabaseName     string `validate:"required"`
	PrimaryDatabaseUser     string `validate:"required"`
	PrimaryDatabasePassword string `validate:"required"`

	ReplicaDatabaseHost     string `validate:"required"`
	ReplicaDatabasePort     string `validate:"required"`
	ReplicaDatabaseName     string `validate:"required"`
	ReplicaDatabaseUser     string `validate:"required"`
	ReplicaDatabasePassword string `validate:"required"`

	CacheURI      string `validate:"required"`
	CachePassword string `validate:"required"`
	CacheDB       int    `validate:"exists"`

	SentryDSN string `validate:"required"`
}

var Secret secrets

func getEnv(key string) string {
	env, _ := os.LookupEnv(key)
	log.Info("Loaded: "+key, nil)
	return env
}

func loadSecrets() {
	log.Info("Loading environment vairables...", nil)

	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to load environment vairables", err.Error(), nil)
	}

	requiredEnvVars := []string{
		"PRIMARY_DB_HOST",
		"PRIMARY_DB_PORT",
		"PRIMARY_DB_NAME",
		"PRIMARY_DB_USER",
		"PRIMARY_DB_PASSWORD",

		"REPLICA_DB_HOST",
		"REPLICA_DB_PORT",
		"REPLICA_DB_NAME",
		"REPLICA_DB_USER",
		"REPLICA_DB_PASSWORD",

		"CACHE_URI",
		"CACHE_PASSWORD",
		"CACHE_DB",

		"SENTRY_DSN",
	}

	// Check is all required env variables exists
	for _, variable := range requiredEnvVars {
		if _, exists := os.LookupEnv(variable); !exists {
			log.Fatal(
				"Failed to load environment variables",
				"Missing required env variable: "+variable,
				nil,
			)
		}
	}

	cacheDB, err := strconv.ParseInt(getEnv("CACHE_DB"), 10, 64)

	if err != nil {
		log.Fatal("Failed to parse CACHE_DB env variable", err.Error(), nil)
	}

	Secret.PrimaryDatabaseHost 	   = getEnv("PRIMARY_DB_HOST")
	Secret.PrimaryDatabasePort     = getEnv("PRIMARY_DB_PORT")
	Secret.PrimaryDatabaseName     = getEnv("PRIMARY_DB_NAME")
	Secret.PrimaryDatabaseUser 	   = getEnv("PRIMARY_DB_USER")
	Secret.PrimaryDatabasePassword = getEnv("PRIMARY_DB_PASSWORD")

	Secret.ReplicaDatabaseHost 	   = getEnv("REPLICA_DB_HOST")
	Secret.ReplicaDatabasePort 	   = getEnv("REPLICA_DB_PORT")
	Secret.ReplicaDatabaseName 	   = getEnv("REPLICA_DB_NAME")
	Secret.ReplicaDatabaseUser 	   = getEnv("REPLICA_DB_USER")
	Secret.ReplicaDatabasePassword = getEnv("REPLICA_DB_PASSWORD")

	Secret.CacheURI 	 = getEnv("CACHE_URI")
	Secret.CachePassword = getEnv("CACHE_PASSWORD")
	Secret.CacheDB 		 = int(cacheDB)

	Secret.SentryDSN = getEnv("SENTRY_DSN")

	log.Info("Loading environment vairables: OK", nil)

	log.Info("Validating secrets...", nil)

	validate := validator.New()

	validate.RegisterValidation("exists", func(fl validator.FieldLevel) bool {
		return true // Always pass (just ensure that the field exists)
	})

	if err := validate.Struct(Secret); err != nil {
		log.Fatal("Secrets validation failed", err.Error(), nil)
	}

	log.Info("Validating secrets: OK", nil)
}
