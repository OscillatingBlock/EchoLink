// config/config.go
package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

// Config holds all app configuration
type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`

	Postgres struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"postgres"`

	MLService struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"ml_service"`

	JWT struct {
		Secret string `mapstructure:"secret"`
	} `mapstructure:"jwt"`
}

// Load config from config.yaml or env
func Load() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	// Allow env override (e.g. POSTGRES_URL)
	viper.AutomaticEnv()

	// Defaults (MVP)
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("ml_service.url", "http://localhost:8000")
	viper.SetDefault("jwt.secret", "hackathon-secret-change-me")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Printf("Config file error: %v", err)
		}
	}

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		log.Fatalf("Unable to decode config: %v", err)
	}

	// Required: Postgres URL
	if c.Postgres.URL == "" {
		c.Postgres.URL = os.Getenv("POSTGRES_URL")
	}
	if c.Postgres.URL == "" {
		log.Fatal("POSTGRES_URL is required")
	}

	return &c
}

