package prompts

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/llmcontext/gomcp/pkg/jsonschema"
)

func savePrompts(baseDirectory string, prompts *PromptsContext) error {
	filePath := filepath.Join(baseDirectory, "prompts.json")
	jsonBytes, err := json.MarshalIndent(prompts, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal prompts: %v", err)
	}
	return os.WriteFile(filePath, jsonBytes, 0644)
}

func loadPrompts(baseDirectory string) (*PromptsContext, error) {
	filePath := filepath.Join(baseDirectory, "prompts.json")

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return &PromptsContext{
			Prompts: []PromptDescription{},
		}, nil
	}

	// read the file and parse it into a list of ShareToolConfig
	// let's generate the schema from the config struct
	configSchema, err := jsonschema.GetSchemaFromAny(&PromptsContext{})
	if err != nil {
		return nil, fmt.Errorf("failed to generate schema from config struct: %v", err)
	}
	// let's check that the file is a valid json file
	jsonBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read shares configuration file %s: %v", filePath, err)
	}

	err = jsonschema.ValidateJsonSchemaWithBytes(configSchema, jsonBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to validate shares configuration file %s: %v", filePath, err)
	}

	var config PromptsContext
	err = json.Unmarshal(jsonBytes, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal shares configuration file %s: %v", filePath, err)
	}

	return &config, nil
}
