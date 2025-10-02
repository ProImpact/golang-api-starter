package config

import (
	"apistarter/internal/env"
	"fmt"
	"log/slog"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Configuration struct {
	Port           string
	DatabaseConfig DatabaseConfig
	Mode           string
	AppName        string
	AllowedOrigins []string
	Key            string
}

type DatabaseConfig struct {
	UserName     string
	Password     string
	Host         string
	Port         int
	DatabaseName string
	Version      int
}

func init() {
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

func NewConfiguration() *Configuration {
	var cfg Configuration
	cfg.Port = viper.GetString("port")
	cfg.DatabaseConfig.DatabaseName = viper.GetString("db.name")
	cfg.DatabaseConfig.Host = viper.GetString("db.host")
	cfg.DatabaseConfig.Password = viper.GetString("db.password")
	cfg.DatabaseConfig.UserName = viper.GetString("db.username")
	cfg.DatabaseConfig.Port = viper.GetInt("db.port")
	cfg.DatabaseConfig.Version = viper.GetInt("db.version")
	cfg.Mode = viper.GetString("mode")
	cfg.AppName = viper.GetString("app-name")
	cfg.AllowedOrigins = viper.GetStringSlice("allowed_origins")
	cfg.Key = viper.GetString("key")
	env.KEY = cfg.Key
	env.APP_NAME = cfg.AppName
	slog.Info("allowed origins", "origins", strings.Join(cfg.AllowedOrigins, ","))
	gin.SetMode(cfg.Mode)
	return &cfg
}
