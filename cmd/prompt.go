package main

import (
	"fmt"
	"os"

	"github.com/llmcontext/mcphost/config"
	"github.com/llmcontext/mcphost/tools/prompts"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(promptCmd)
}

var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "list of prompts",
	Long:  "list all prompts",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		NewPromptsTools := prompts.NewPromptsTools(config.DefaultConfigurationDirectory)

		prompts, err := NewPromptsTools.ListPrompts()
		if err != nil {
			fmt.Println("error listing prompts:", err)
			os.Exit(1)
		}

		fmt.Println("Prompts:")
		for _, prompt := range prompts {
			fmt.Printf("- File: %s\n", prompt.YamlPath)
			if prompt.Error != nil {
				fmt.Printf("  Error: %v\n", prompt.Error)
			} else {
				for _, prompt := range prompt.Prompts {
					fmt.Printf("  Prompt: %s - %s\n", prompt.Name, prompt.Description)
				}
			}
		}
	},
}
