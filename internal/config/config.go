package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	DBPath       string `mapstructure:"db_path"`
	GeminiAPIKey string `mapstructure:"gemini_api_key"`
}

func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(configDir, "kairos")
	if err := os.MkdirAll(configPath, 0755); err != nil {
		return nil, err
	}

	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Default DB path in ~/.local/share/kairos
	localSharePath := filepath.Join(home, ".local", "share", "kairos")
	viper.SetDefault("db_path", filepath.Join(localSharePath, "kairos.db"))
	viper.BindEnv("gemini_api_key", "GEMINI_API_KEY")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired or write default
			// For now, we just rely on defaults and env vars
		} else {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	if cfg.GeminiAPIKey == "" {
		fmt.Print("Gemini API Key not found. Please enter it: ")
		reader := bufio.NewReader(os.Stdin)
		key, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("failed to read input: %w", err)
		}
		key = strings.TrimSpace(key)
		cfg.GeminiAPIKey = key
		viper.Set("gemini_api_key", key)

		if err := viper.WriteConfigAs(filepath.Join(configPath, "config.yaml")); err != nil {
			return nil, fmt.Errorf("failed to save config: %w", err)
		}
		fmt.Print("\033[H\033[2J")
	}

	return &cfg, nil
}
