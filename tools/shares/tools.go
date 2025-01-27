package shares

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/llmcontext/gomcp/pkg/jsonschema"
	"github.com/llmcontext/gomcp/types"
)

type ShareToolConfiguration struct {
	BaseDirectory string
}

type ShareDescription struct {
	ShareName string `json:"shareName"`
	Path      string `json:"path"`
}

type ShareToolContext struct {
	Shares []ShareDescription `json:"shares"`
}

func loadShares(baseDirectory string) (*ShareToolContext, error) {
	filePath := filepath.Join(baseDirectory, "shares.json")

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return &ShareToolContext{
			Shares: []ShareDescription{},
		}, nil
	}

	// read the file and parse it into a list of ShareToolConfig
	// let's generate the schema from the config struct
	configSchema, err := jsonschema.GetSchemaFromAny(&ShareToolContext{})
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

	var config ShareToolContext
	err = json.Unmarshal(jsonBytes, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal shares configuration file %s: %v", filePath, err)
	}

	return &config, nil
}

func ShareToolInit(ctx context.Context, config *ShareToolConfiguration) (*ShareToolContext, error) {
	return loadShares(config.BaseDirectory)
}

type ListSharesInput struct {
}

func ListShares(ctx context.Context, config *ShareToolContext, input *ListSharesInput, output types.ToolCallResult) error {
	logger := types.GetLogger(ctx)
	logger.Debug("ListShares", types.LogArg{
		"config": config,
	})
	shareNames := []string{}
	for _, share := range config.Shares {
		shareNames = append(shareNames, share.ShareName)
	}
	output.AddJSONTextContent(shareNames)
	return nil
}
