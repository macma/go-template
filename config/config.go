package config

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Config *EnvConfig
var DbConnection *sql.DB

const lvsIteration = 3

func findFileLevelsUp(filename string, lvs int) (*string, error) {
	for i := 0; i < lvs; i++ {
		if _, err := os.Stat(filename); err == nil {
			log.Info(fmt.Sprintf("found the env file, path: %v; it is %d th layer above the path.", filename, i))
			return &filename, nil
		}
		filename = "../" + filename
	}
	return nil, errors.New("No files found")
}

func InitViperConfig(cfg *EnvConfig, envfile string) {
	viper.SetDefault("rabbit_protocol", "amqp")
	viper.SetDefault("POSTGRES_DATABASE", "postgresdatabase")
	viper.SetDefault("POSTGRES_ADDRESS", "postgresql:5432")
	viper.SetDefault("POSTGRES_SSLMODE", "require")
	viper.SetDefault("POSTGRES_MAX_OPEN_CONNS", "20")
	viper.SetDefault("POSTGRES_MAX_IDLE_CONNS", "2")
	viper.SetDefault("MIGRATION_PATH", "./db/migrations")

	filePath, err := findFileLevelsUp(envfile, lvsIteration)
	if err != nil {
		log.Error(fmt.Sprintf("could not find the config file: %s", envfile))
	}
	if filePath != nil {
		log.Info(fmt.Sprintf("Trying to read the config file: %s", *filePath))
		viper.SetConfigFile(*filePath)
		err = viper.ReadInConfig()
		if err != nil {
			log.Error(fmt.Sprintf("could not read config file: %v", err))
			return
		}
	}

	viper.AutomaticEnv()

	cfg.PostgresUsername = viper.GetString("POSTGRES_USERNAME")
	cfg.PostgresPassword = viper.GetString("POSTGRES_PASSWORD")
	cfg.PostgresDatabase = viper.GetString("POSTGRES_DATABASE")
	cfg.PostgresAddress = viper.GetString("POSTGRES_ADDRESS")
	cfg.PostgresSslMode = viper.GetString("POSTGRES_SSLMODE")
	cfg.PostgresMaxOpenConns = viper.GetInt("POSTGRES_MAX_OPEN_CONNS")
	cfg.PostgresMaxIdleConns = viper.GetInt("POSTGRES_MAX_IDLE_CONNS")
	cfg.MigrationPath = viper.GetString("MIGRATION_PATH")
}

//EnvConfig environment variables to parse to get config
type EnvConfig struct {
	PostgresUsername     string `env:"POSTGRES_USERNAME,required"`
	PostgresPassword     string `env:"POSTGRES_PASSWORD,required"`
	PostgresDatabase     string `env:"POSTGRES_DATABASE" envDefault:"fiat-integration"`
	PostgresAddress      string `env:"POSTGRES_ADDRESS" envDefault:"postgresql:5432"`
	PostgresSslMode      string `env:"POSTGRES_SSLMODE" envDefault:"require"`
	PostgresMaxOpenConns int    `env:"POSTGRES_MAX_OPEN_CONNS" envDefault:"20"`
	PostgresMaxIdleConns int    `env:"POSTGRES_MAX_IDLE_CONNS" envDefault:"2"`
	MigrationPath        string `env:"MIGRATION_PATH" envDefault:"./db/migrations"`
}
