package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	API         APIConfig         `mapstructure:"api"`
	Editor      EditorConfig      `mapstructure:"editor"`
	Preferences PreferencesConfig `mapstructure:"preferences"`
}

// APIConfig holds API-related configuration
type APIConfig struct {
	BaseURL string `mapstructure:"base_url"`
	Timeout int    `mapstructure:"timeout"` // in seconds
}

// EditorConfig holds editor-related configuration
type EditorConfig struct {
	ExternalEditor string `mapstructure:"external_editor"`
}

// PreferencesConfig holds user preferences
type PreferencesConfig struct {
	DefaultNoteType  string `mapstructure:"default_note_type"`
	AutoSaveInterval int    `mapstructure:"auto_save_interval"` // in seconds
	Theme            string `mapstructure:"theme"`
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig() (*Config, error) {
	// Set default values
	viper.SetDefault("api.base_url", "http://localhost:8080") // for localhost testing
	// viper.SetDefault("api.base_url", "API_SERVER") // change here for 'prod' server
	viper.SetDefault("api.timeout", 30)
	viper.SetDefault("editor.external_editor", os.Getenv("EDITOR"))
	if viper.GetString("editor.external_editor") == "" {
		viper.SetDefault("editor.external_editor", "vi")
	}
	viper.SetDefault("preferences.default_note_type", "note")
	viper.SetDefault("preferences.auto_save_interval", 30)
	viper.SetDefault("preferences.theme", "dark")

	// Set config file path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("get home dir: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config", "kg-cli")
	configFile := filepath.Join(configDir, "config.yaml")

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return nil, fmt.Errorf("create config dir: %w", err)
	}

	// Set config file name and path
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	// Read config file if it exists
	if _, err := os.Stat(configFile); err == nil {
		if err := viper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("read config file: %w", err)
		}
	}

	// Allow environment variables to override config
	viper.SetEnvPrefix("KG_CLI")
	viper.AutomaticEnv()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &config, nil
}

// SaveConfig saves the current configuration to file
func SaveConfig(config *Config) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("get home dir: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config", "kg-cli")
	configFile := filepath.Join(configDir, "config.yaml")

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	// Set config values
	viper.Set("api.base_url", config.API.BaseURL)
	viper.Set("api.timeout", config.API.Timeout)
	viper.Set("editor.external_editor", config.Editor.ExternalEditor)
	viper.Set("preferences.default_note_type", config.Preferences.DefaultNoteType)
	viper.Set("preferences.auto_save_interval", config.Preferences.AutoSaveInterval)
	viper.Set("preferences.theme", config.Preferences.Theme)

	// Write config file
	if err := viper.SafeWriteConfigAs(configFile); err != nil {
		// If file exists, use WriteConfig instead
		if err := viper.WriteConfigAs(configFile); err != nil {
			return fmt.Errorf("write config file: %w", err)
		}
	}

	return nil
}
