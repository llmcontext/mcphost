package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/llmcontext/gomcp/pkg/jsonschema"
)

type LoggingInfo struct {
	File       string `json:"file,omitempty"`
	Level      string `json:"level,omitempty"`
	WithStderr bool   `json:"withStderr,omitempty"`
}

type McpHostConfiguration struct {
	ConfigVersion int          `json:"v"`
	Logging       *LoggingInfo `json:"logging,omitempty"`
}

func getDefaultMcpHostConfigurationPath() string {
	return filepath.Join(DefaultConfigurationDirectory, "mcphost.json")
}

// LoadConfig loads the configuration from a file
func LoadMcpHostConfiguration() (*McpHostConfiguration, error) {
	configFilePath := getDefaultMcpHostConfigurationPath()

	// Check if the file exists
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// create the file and parent directories
		if err := CreateFileWithDirs(configFilePath); err != nil {
			return nil, fmt.Errorf("failed to create configuration file: %v", err)
		}
		mcpHostConfig := &McpHostConfiguration{
			ConfigVersion: 1,
			Logging: &LoggingInfo{
				File:       "mcphost.log",
				Level:      "debug",
				WithStderr: true,
			},
		}
		jsonBytes, err := json.MarshalIndent(mcpHostConfig, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal configuration file: %v", err)
		}
		if err := os.WriteFile(configFilePath, jsonBytes, 0644); err != nil {
			return nil, fmt.Errorf("failed to write configuration file: %v", err)
		}
	}

	// let's generate the schema from the config struct
	configSchema, err := jsonschema.GetSchemaFromAny(&McpHostConfiguration{})
	if err != nil {
		return nil, fmt.Errorf("failed to generate schema from config struct: %v", err)
	}
	// let's check that the file is a valid json file
	jsonBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	err = jsonschema.ValidateJsonSchemaWithBytes(configSchema, jsonBytes)
	if err != nil {
		return nil, err
	}

	var config McpHostConfiguration
	err = json.Unmarshal(jsonBytes, &config)
	if err != nil {
		return nil, err
	}

	// update the file path to be absolute
	if config.Logging != nil && config.Logging.File != "" {
		config.Logging.File = updateFilePath(config.Logging.File)
	}
	return &config, nil
}

func updateFilePath(path string) string {
	if path == "" {
		return path
	}
	if !filepath.IsAbs(path) {
		path = filepath.Join(DefaultConfigurationDirectory, path)
	}
	return path
}

// CreateFileWithDirs creates a file and all necessary parent directories.
// Returns an error if the file cannot be created or if directories cannot be created.
func CreateFileWithDirs(path string) error {
	// Create parent directories if they don't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directories for %s: %v", path, err)
	}

	// Create or truncate the file
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", path, err)
	}
	defer file.Close()

	return nil
}
