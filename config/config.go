package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv     string
	Port       int
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
}

func Load() (Config, error) {
	v := viper.New()
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	_ = v.ReadInConfig()

	v.SetDefault("APP_ENV", "dev")
	v.SetDefault("PORT", 8080)
	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", 5432)
	v.SetDefault("DB_USER", "postgres")
	v.SetDefault("DB_PASSWORD", "postgres")
	v.SetDefault("DB_NAME", "contest")

	cfg := Config{
		AppEnv:     v.GetString("APP_ENV"),
		Port:       v.GetInt("PORT"),
		DBHost:     v.GetString("DB_HOST"),
		DBPort:     v.GetInt("DB_PORT"),
		DBUser:     v.GetString("DB_USER"),
		DBPassword: v.GetString("DB_PASSWORD"),
		DBName:     v.GetString("DB_NAME"),
	}

	if cfg.Port <= 0 {
		return Config{}, fmt.Errorf("invalid PORT")
	}
	if cfg.DBPort <= 0 {
		return Config{}, fmt.Errorf("invalid DB_PORT")
	}

	return cfg, nil
}
