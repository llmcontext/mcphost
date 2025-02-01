package prompts

import (
	"fmt"
	"os"

	"github.com/llmcontext/gomcp/pkg/prompts"
	"github.com/llmcontext/gomcp/types"
)

type PromptsTools struct {
	baseDirectory string
}

type PromptDescription struct {
	PromptYamlPath string `json:"yamlPath"`
}

type PromptsContext struct {
	Prompts []PromptDescription `json:"prompts"`
}

// Describes a prompt file and its contents
type PromptFile struct {
	YamlPath string
	Error    error
	Prompts  []*prompts.PromptDefinition
}

func NewPromptsTools(baseDirectory string) *PromptsTools {
	return &PromptsTools{
		baseDirectory: baseDirectory,
	}
}

func (p *PromptsTools) Register(mcpServerDefinition types.McpSdkServerDefinition) error {
	promptsList, err := p.ListPrompts()
	if err != nil {
		return err
	}
	for _, prompt := range promptsList {
		duplicatedPrompts, err := mcpServerDefinition.AddTemplateYamlFile(prompt.YamlPath)
		if err != nil {
			return err
		}
		for _, duplicatedPrompt := range duplicatedPrompts {
			// write on stderr
			fmt.Fprintf(os.Stderr, "duplicated prompt: %s, in file: %s\n",
				duplicatedPrompt.PromptName, duplicatedPrompt.FilePath)
		}
	}
	return nil
}

func (p *PromptsTools) ListPrompts() ([]PromptFile, error) {
	promptsFiles, err := loadPrompts(p.baseDirectory)
	if err != nil {
		return nil, err
	}

	result := []PromptFile{}

	for _, promptFile := range promptsFiles.Prompts {
		yamlPath := promptFile.PromptYamlPath
		prompts, err := prompts.LoadPromptYamlFile(yamlPath)
		if err != nil {
			result = append(result, PromptFile{
				YamlPath: yamlPath,
				Error:    err,
			})
			continue
		}

		promptFile := PromptFile{
			YamlPath: yamlPath,
			Prompts:  prompts.Prompts,
			Error:    nil,
		}
		result = append(result, promptFile)
	}

	return result, nil
}

func (f *PromptsContext) hasFile(yamlFile string) bool {
	for _, prompt := range f.Prompts {
		if prompt.PromptYamlPath == yamlFile {
			return true
		}
	}
	return false
}

func (p *PromptsTools) AddPrompts(yamlFiles []string) []string {
	out := []string{}
	promptsFiles, err := loadPrompts(p.baseDirectory)
	if err != nil {
		out = append(out, fmt.Sprintf("error loading prompts: %s", err))
		return out
	}
	for _, yamlFile := range yamlFiles {
		// check if the yaml file is already in the promptsFiles
		if promptsFiles.hasFile(yamlFile) {
			fmt.Println("prompt yaml file already exists:", yamlFile)
			continue
		}
		// load the prompt yaml file to see if it's valid
		_, err := prompts.LoadPromptYamlFile(yamlFile)
		if err != nil {
			out = append(out, fmt.Sprintf("error loading prompt yaml file: %s", err))
			continue
		}

		out = append(out, fmt.Sprintf("prompt yaml file added: %s", yamlFile))
		promptsFiles.Prompts = append(promptsFiles.Prompts, PromptDescription{
			PromptYamlPath: yamlFile,
		})
	}

	// we save the promptsFiles to the file system
	err = savePrompts(p.baseDirectory, promptsFiles)
	if err != nil {
		out = append(out, fmt.Sprintf("error saving prompts: %s", err))
	}
	return out
}
