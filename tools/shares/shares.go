package shares

import "github.com/llmcontext/gomcp/types"

type ShareToolConfig struct {
	ShareName string `json:"shareName"`
	Path      string `json:"path"`
}

type ShareTools struct {
	baseDirectory string
}

func NewShareTools(baseDirectory string) *ShareTools {
	return &ShareTools{
		baseDirectory: baseDirectory,
	}
}

func (s *ShareTools) Register(mcpServerDefinition types.McpSdkServerDefinition) error {

	return nil
}
