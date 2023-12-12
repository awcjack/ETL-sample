package config

import (
	"errors"

	"github.com/spf13/viper"
)

// Config struct that contain all config for this application
type Config struct {
	Application ApplicationConfig
	Datasource  []DataSourceConfig
	Database    DatabaseConfig
}

// Application config
type ApplicationConfig struct {
	// logger level (all possible value: panic, fatal, error, warn, warning, info, debug and trace)
	LogLevel string
	// channel size for keeping data in memory before storing to persist storage
	ProcessPipelineSize int
	// flag to enable bulk insert
	BulkInsert bool
	// bulk insert number in each batch (ignored if BulkInsert flag is disabled)
	BulkInsertSize int
	// trigger bulk insert in every x second if didn't fill the bulk insert size (ignored if BulkInsert flag is disabled)
	BulkInsertInterval int
}

// Data source config
type DataSourceConfig struct {
	// data source name
	Name string
	// data source type ["http", "file", etc]
	Type string
	// transformer type ["random-data-api", "random-data-api-v2", etc]
	Transformer string
	// source ["http url", "file path", etc]
	Source string
}

// Database config
type DatabaseConfig struct {
	// database type (possible to switching from postgresql to mysql/mongodb/memory etc)
	Type string
	// connection string of database
	ConnectionString string
}

// Loading all confing from environment
// Using os environment due to this application is expected to be deployed to docker/k8s
// setting environment is easiest way to config application in docker/k8s comparing reading config file
func LoadConfig() (*Config, error) {
	// reading config.json from current path
	viper.AddConfigPath("./")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.ReadInConfig()

	c := &Config{}

	// Database Config
	c.Database.Type = getStringConfigWithDefault("Database.Type", "postgresql")

	c.Database.ConnectionString = viper.GetString("Database.ConnectionString")
	if c.Database.ConnectionString == "" {
		return nil, errors.New("missing connection string")
	}

	// Application Config
	c.Application.LogLevel = getStringConfigWithDefault("Application.LogLevel", "info")

	c.Application.ProcessPipelineSize = getIntConfigWithDefault("Application.ProcessPipelineSize", 10)

	c.Application.BulkInsert = viper.GetBool("Application.BulkInsert")

	c.Application.BulkInsertSize = viper.GetInt("Application.BulkInsertSize")

	c.Application.BulkInsertInterval = viper.GetInt("Application.BulkInsertInterval")

	// Data source Config
	viper.UnmarshalKey("Datasource", &c.Datasource)

	return c, nil
}

// get string config from environment
// if value is not found fomr environemnt, defaultValue will be used
func getStringConfigWithDefault(key, defaultValue string) string {
	value := viper.GetString(key)
	if value == "" {
		return defaultValue
	}

	return value
}

// get string config from environment
// if value is not found fomr environemnt, defaultValue will be used
func getIntConfigWithDefault(key string, defaultValue int) int {
	value := viper.GetInt(key)
	if value == 0 {
		return defaultValue
	}

	return value
}
