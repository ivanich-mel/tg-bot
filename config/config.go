package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	BotToken string `mapstructure:"tg_token"`
	DB       PostgresConfig
}
type PostgresConfig struct {
	Driver string `mapstructure:"driver"`
	Host   string `mapstructure:"host"`
	User   string `mapstructure:"user"`
	Pass   string `mapstructure:"pass"`
	Name   string `mapstructure:"name"`
	Port   string `mapstructure:"port"`
}

func Init() (*Config, error) {
	var config Config

	if err := setViper(); err != nil {
		return nil, err
	}
	if err := unmarshal(&config); err != nil {
		return nil, err
	}
	config.BotToken = os.Getenv("TGBOT_TOKEN")
	return &config, nil
}

func setViper() error {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

func unmarshal(config *Config) error {
	if err := viper.Unmarshal(&config); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("postgres", &config.DB); err != nil {
		return err
	}
	return nil
}
func (p *PostgresConfig) GetDBSource() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		p.User,
		p.Pass,
		p.Host,
		p.Port,
		p.Name,
	)
}
