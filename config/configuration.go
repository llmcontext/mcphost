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

// TODO: create configuration file on startup

func getDefaultMcpHostConfigurationPath() string {
	return filepath.Join(DefaultConfigurationDirectory, "mcphost.json")
}

// LoadConfig loads the configuration from a file
func LoadMcpHostConfiguration() (*McpHostConfiguration, error) {
	configFilePath := getDefaultMcpHostConfigurationPath()

	// Check if the file exists
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("mcphost configuration file does not exist: %s", configFilePath)
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
