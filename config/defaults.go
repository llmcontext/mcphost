package config

import (
	"os"
	"path/filepath"
)

const (
	ServerName    = "mcphost"
	ServerVersion = "0.1.0"
)

var DefaultConfigurationDirectory = filepath.Join(os.Getenv("HOME"), ".mcphost")
