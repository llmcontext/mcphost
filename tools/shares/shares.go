package shares

import (
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
