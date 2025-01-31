package main

import (
	"path/filepath"
)

func convertRelativePathToAbsolute(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}
	return filepath.Abs(path)
}

func getYamlFiles(path string) ([]string, error) {
	yamlFiles, err := filepath.Glob(filepath.Join(path, "*.yaml"))
	if err != nil {
		return nil, err
	}
	for _, yamlFile := range yamlFiles {
		absolutePath, err := convertRelativePathToAbsolute(yamlFile)
		if err != nil {
			return nil, err
		}
		yamlFiles = append(yamlFiles, absolutePath)
	}
	return yamlFiles, nil
}
