package config

import (
	"io"
	"os"
	"time"

	"github.com/abaxoth0/Vega/libs/go/packages/logger"
	"github.com/go-playground/validator"
	"gopkg.in/yaml.v3"
)

var log = logger.NewSource("CONFIG", logger.Default)

// Wrapper for time.ParseDuration. Panics on error.
func parseDuration(raw string) time.Duration {
	v, e := time.ParseDuration(raw)
	if e != nil {
		panic(e)
	}
	return v
}

type dbConfig struct {
	RawDefaultQueryTimeout string `yaml:"db-default-queuery-timeout" validate:"required"`
	SkipPostConnection     bool   `yaml:"db-skip-post-connection" validate:"exists"`
}

func (c *dbConfig) DefaultQueryTimeout() time.Duration {
	return parseDuration(c.RawDefaultQueryTimeout)
}

type cacheConfig struct {
	RawPoolTimeout      string `yaml:"cache-pool-timeout" validate:"required"`
	RawOperationTimeout string `yaml:"cache-operation-timeout" validate:"required"`
	RawTTL              string `yaml:"cache-ttl" validate:"required"`
}

func (c *cacheConfig) PoolTimeout() time.Duration {
	return parseDuration(c.RawPoolTimeout)
}

func (c *cacheConfig) OperationTimeout() time.Duration {
	return parseDuration(c.RawOperationTimeout)
}

func (c *cacheConfig) TTL() time.Duration {
	return parseDuration(c.RawTTL)
}

type debugConfig struct {
	Enabled           bool `yaml:"debug-mode" validate:"exists"`
	SafeDatabaseScans bool `yaml:"debug-safe-db-scans" validate:"exists"`
	LogDbQueries      bool `yaml:"debug-log-db-queries" validate:"exists"`
}

type appConfig struct {
	ShowLogs         bool   `yaml:"show-logs" validate:"exists"`
	TraceLogsEnabled bool   `yaml:"trace-logs" validate:"exists"`
	ServiceID        string `yaml:"service-id" validate:"required"`
}

type sentryConfig struct {
	TraceSampleRate float64 `yaml:"sentry-trace-sample-rate" validate:"required,min=0.0,max=1.0"`
}

type configs struct {
	dbConfig         `yaml:",inline"`
	cacheConfig      `yaml:",inline"`
	debugConfig      `yaml:",inline"`
	appConfig        `yaml:",inline"`
	sentryConfig     `yaml:",inline"`
}

var (
	DB     *dbConfig
	Cache  *cacheConfig
	Debug  *debugConfig
	App    *appConfig
	Sentry *sentryConfig
)

var isInit bool = false

func loadConfig(path string, dest *configs) {
	log.Info("Reading config file...", nil)

	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Failed to open config file", err.Error(), nil)
	}

	rawConfig, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("Failed to read config file", err.Error(), nil)
	}

	log.Info("Reading config file: OK", nil)

	log.Info("Parsing config file...", nil)

	if err := yaml.Unmarshal(rawConfig, dest); err != nil {
		log.Fatal("Failed to parse config file", err.Error(), nil)
	}

	log.Info("Parsing config file: OK", nil)

	log.Info("Validating config...", nil)

	validate := validator.New()
	validate.RegisterValidation("exists", func(fl validator.FieldLevel) bool {
		return true // Always pass (just ensure that the field exists)
	})

	if err := validate.Struct(dest); err != nil {
		log.Fatal("Failed to validate config", err.Error(), nil)
		os.Exit(1)
	}

	log.Info("Validating config: OK", nil)
}

func Init() {
	if isInit {
		log.Fatal("Failed to initialize config", "Config already initialized", nil)
	}

	log.Info("Initializing...", nil)

	configs := new(configs)

	loadConfig("config.yaml", configs)
	loadSecrets()

	DB 	   = &configs.dbConfig
	Cache  = &configs.cacheConfig
	Debug  = &configs.debugConfig
	App    = &configs.appConfig
	Sentry = &configs.sentryConfig

	log.Info("Initializing: OK", nil)

	isInit = true
}
