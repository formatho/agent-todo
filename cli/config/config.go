package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	ServerURL string `mapstructure:"server_url"`
	Token     string `mapstructure:"token"`
	APIKey    string `mapstructure:"api_key"`
	Insecure  bool   `mapstructure:"insecure"`
}

var cfg *Config

func Init() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".agent-todo")
	configPath := filepath.Join(configDir, "config.yaml")

	// Create config directory if it doesn't exist
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("error creating config directory: %w", err)
		}
	}

	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// Set defaults
	viper.SetDefault("server_url", "http://localhost:8080")
	viper.SetDefault("insecure", false)

	// Read config file if it exists
	if _, err := os.Stat(configPath); err == nil {
		if err := viper.ReadInConfig(); err != nil {
			return fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Unmarshal config
	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return fmt.Errorf("error unmarshaling config: %w", err)
	}

	return nil
}

func Get() *Config {
	return cfg
}

func Save() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".agent-todo", "config.yaml")

	viper.Set("server_url", cfg.ServerURL)
	viper.Set("token", cfg.Token)
	viper.Set("api_key", cfg.APIKey)
	viper.Set("insecure", cfg.Insecure)

	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}

func SetServerURL(url string) {
	cfg.ServerURL = url
}

func SetToken(token string) {
	cfg.Token = token
}

func SetAPIKey(apiKey string) {
	cfg.APIKey = apiKey
}

func GetServerURL() string {
	return cfg.ServerURL
}

func GetToken() string {
	return cfg.Token
}

func GetAPIKey() string {
	return cfg.APIKey
}

func IsInsecure() bool {
	return cfg.Insecure
}

func ClearAuth() {
	cfg.Token = ""
	cfg.APIKey = ""
}
