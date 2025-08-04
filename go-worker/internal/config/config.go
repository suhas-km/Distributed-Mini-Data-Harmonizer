package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	// Server settings
	Port int `json:"port"`
	Host string `json:"host"`

	// Worker settings
	WorkerCount int `json:"worker_count"`
	QueueSize   int `json:"queue_size"`

	// File paths
	InputDir  string `json:"input_dir"`
	OutputDir string `json:"output_dir"`

	// API settings
	PythonAPIURL string `json:"python_api_url"`
}

// LoadConfig loads configuration from environment variables or defaults
func LoadConfig() (*Config, error) {
	config := &Config{
		Port:         getEnvAsInt("PORT", 8081),
		Host:         getEnvAsString("HOST", "0.0.0.0"),
		WorkerCount:  getEnvAsInt("WORKER_COUNT", 3),
		QueueSize:    getEnvAsInt("QUEUE_SIZE", 100),
		InputDir:     getEnvAsString("INPUT_DIR", "../uploads"),
		OutputDir:    getEnvAsString("OUTPUT_DIR", "../results"),
		PythonAPIURL: getEnvAsString("PYTHON_API_URL", "http://localhost:8080/api/v1"),
	}

	// Create directories if they don't exist
	if err := os.MkdirAll(config.InputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create input directory: %w", err)
	}

	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	return config, nil
}

// LoadConfigFromFile loads configuration from a JSON file
func LoadConfigFromFile(filePath string) (*Config, error) {
	// Expand file path if it contains ~
	if filePath[:1] == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		filePath = filepath.Join(home, filePath[1:])
	}

	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// Helper functions to get environment variables with defaults
func getEnvAsString(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
