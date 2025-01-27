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

	mcpToolsDefinition.AddTool("listShares", "List all share names. A share is the name of a directory that is shared with the client", ListShares)
	mcpToolsDefinition.AddTool("getShareFiles",
		`Get a detailed listing of all files and directories in a specified share name
Results clearly distinguish between files and directories with [FILE] and [DIR] prefixes
This tool is essential for understanding directory structure and  finding specific files within a directory.`,
		GetShareFiles)
	mcpToolsDefinition.AddTool("readFile",
		`Read the complete contents of a file from the share.
Handles various text encodings and provides detailed error messages if the file cannot be read.
Use this tool when you need to examine the contents of a single file.`,
		ReadFile)

	return nil
}

func (s *ShareTools) ListShares() ([]ShareDescription, error) {
	shares, err := loadShares(s.baseDirectory)
	if err != nil {
		return nil, err
	}
	return shares.Shares, nil
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
			return fmt.Errorf("share with same path already exists: %s", share.ShareName)
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

func (s *ShareTools) DeleteShareForDirectory(path string) (*ShareDescription, error) {
	shares, err := loadShares(s.baseDirectory)
	if err != nil {
		return nil, err
	}

	var theShare *ShareDescription
	for i, share := range shares.Shares {
		if share.Path == path {
			shares.Shares = append(shares.Shares[:i], shares.Shares[i+1:]...)
			theShare = &share
			break
		}
	}

	err = saveShares(s.baseDirectory, shares)
	if err != nil {
		return nil, err
	}

	if theShare == nil {
		return nil, fmt.Errorf("share not found for path: %s", path)
	}

	return theShare, nil
}
