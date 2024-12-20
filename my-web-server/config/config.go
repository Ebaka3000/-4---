package config

import (
    "github.com/spf13/viper"
)

type Config struct {
    DatabaseURL string
    NATSURL     string
}

func LoadConfig() (*Config, error) {
    viper.SetConfigName("config")
    viper.AddConfigPath(".")
    viper.AutomaticEnv()

    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }

    config := &Config{
        DatabaseURL: viper.GetString("DATABASE_URL"),
        NATSURL:     viper.GetString("NATS_URL"),
    }

    return config, nil
}