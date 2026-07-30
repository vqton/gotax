package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	JWTSecret            string   `mapstructure:"jwt_secret"`
	DatabaseURL          string   `mapstructure:"database_url"`
	ServerPort           string   `mapstructure:"server_port"`
	TrustedProxies       []string `mapstructure:"trusted_proxies"`
	CORSOrigins          []string `mapstructure:"cors_origins"`
	LogLevel             string   `mapstructure:"log_level"`
	LogFormat            string   `mapstructure:"log_format"`
	GORMMaxOpenConns     int      `mapstructure:"gorm_max_open_conns"`
	GORMMaxIdleConns     int      `mapstructure:"gorm_max_idle_conns"`
	GORMConnMaxLifetimeMS int    `mapstructure:"gorm_conn_max_lifetime_ms"`
}

func defaults(v *viper.Viper) {
	v.SetDefault("jwt_secret", "")
	v.SetDefault("database_url", "")
	v.SetDefault("server_port", ":8080")
	v.SetDefault("trusted_proxies", []string{"127.0.0.1", "::1"})
	v.SetDefault("cors_origins", []string{"http://localhost:8080", "http://127.0.0.1:8080"})
	v.SetDefault("log_level", "info")
	v.SetDefault("log_format", "console")
	v.SetDefault("gorm_max_open_conns", 20)
	v.SetDefault("gorm_max_idle_conns", 5)
	v.SetDefault("gorm_conn_max_lifetime_ms", 1800000)
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	defaults(v)

	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("read config: %w", err)
		}
		fmt.Println("config: no config.yaml found, using env + defaults")
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}
