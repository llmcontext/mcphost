package shares

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

func saveShares(baseDirectory string, shares *ShareToolContext) error {
	filePath := filepath.Join(baseDirectory, "shares.json")
	jsonBytes, err := json.MarshalIndent(shares, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal shares: %v", err)
	}
	return os.WriteFile(filePath, jsonBytes, 0644)
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

type GetShareFilesInput struct {
	ShareName string
}

func GetShareFiles(ctx context.Context, config *ShareToolContext, input *GetShareFilesInput, output types.ToolCallResult) error {
	files := []string{}
	found := false
	for _, share := range config.Shares {
		if share.ShareName == input.ShareName {
			found = true
			shareFiles, err := listFilesRecursively(share.Path)
			if err != nil {
				return err
			}
			files = append(files, shareFiles...)
		}
	}
	if !found {
		return fmt.Errorf("share %s not found", input.ShareName)
	}
	output.AddJSONTextContent(files)
	return nil
}

func listFilesRecursively(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if path == root {
			return nil
		}

		// Get path relative to root
		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %v", err)
		}

		if info.IsDir() {
			// ignore hidden directories
			if !strings.HasPrefix(relPath, ".") {
				files = append(files, fmt.Sprintf("[DIR] %s", relPath))
			}
		} else {
			files = append(files, fmt.Sprintf("[FILE] %s", relPath))
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %v", err)
	}

	return files, nil
}
