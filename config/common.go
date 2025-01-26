package config

import (
	"path/filepath"
)

type LoggingInfo struct {
	File       string `json:"file,omitempty"`
	Level      string `json:"level,omitempty"`
	WithStderr bool   `json:"withStderr,omitempty"`
}

func updateFilePath(path string) string {
	if path == "" {
		return path
	}
	if !filepath.IsAbs(path) {
		path = filepath.Join(DefaultHubConfigurationDirectory, path)
	}
	return path
}

func (c *LoggingInfo) UpdateFilePaths() {
	c.File = updateFilePath(c.File)
}
