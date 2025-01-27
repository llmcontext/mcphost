package shares

import (
	"fmt"

	"github.com/llmcontext/gomcp/types"
)

type ShareTools struct {
	baseDirectory string
}

func NewShareTools(baseDirectory string) *ShareTools {
	return &ShareTools{
		baseDirectory: baseDirectory,
	}
}

func (s *ShareTools) Register(mcpServerDefinition types.McpSdkServerDefinition) error {
	config := &ShareToolConfiguration{
		BaseDirectory: s.baseDirectory,
	}

	mcpToolsDefinition := mcpServerDefinition.WithTools(config, ShareToolInit)

	mcpToolsDefinition.AddTool("listShares", "List all share names", ListShares)

	return nil
}

func (s *ShareTools) AddShare(name string, path string) error {
	shares, err := loadShares(s.baseDirectory)
	if err != nil {
		return err
	}

	// check if the share already exists
	for _, share := range shares.Shares {
		if share.ShareName == name {
			return fmt.Errorf("share with same namealready exists: %s", name)
		}

		if share.Path == path {
			return fmt.Errorf("sharewith same path already exists: %s", share.ShareName)
		}
	}

	shares.Shares = append(shares.Shares, ShareDescription{
		ShareName: name,
		Path:      path,
	})

	err = saveShares(s.baseDirectory, shares)
	if err != nil {
		return err
	}

	return nil
}
