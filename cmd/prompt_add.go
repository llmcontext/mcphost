package main

import (
	"fmt"
	"os"

	"github.com/llmcontext/mcphost/config"
	"github.com/llmcontext/mcphost/tools/prompts"
	"github.com/spf13/cobra"
)

func init() {
	promptAddCmd.Flags().StringP("path", "p", "", "the path to the yaml file containing the prompt")
	promptCmd.AddCommand(promptAddCmd)
}

var promptAddCmd = &cobra.Command{
	Use:   "add",
	Short: "add prompts",
	Long:  "add prompts to the mcphost server, if no path is provided, we try to add all the .yaml files in the current directory",
	Run: func(cmd *cobra.Command, args []string) {
		yamlFiles := []string{}
		var err error
		// retrieve the name and path from the flags
		path, _ := cmd.Flags().GetString("path")
		if path != "" {
			absPath, err := convertRelativePathToAbsolute(path)
			if err != nil {
				fmt.Println("error converting relative path to absolute:", err)
				os.Exit(1)
			}
			yamlFiles = append(yamlFiles, absPath)
		} else {
			yamlFiles, err = getYamlFiles(path)
			if err != nil {
				fmt.Println("error getting yaml files:", err)
				os.Exit(1)
			}
		}
		if len(yamlFiles) == 0 {
			fmt.Println("no yaml files found")
			os.Exit(1)
		}
		promptsTools := prompts.NewPromptsTools(config.DefaultConfigurationDirectory)
		out := promptsTools.AddPrompts(yamlFiles)
		for _, msg := range out {
			fmt.Println(msg)
		}
	},
}
