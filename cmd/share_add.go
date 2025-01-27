package main

import (
	"fmt"
	"os"

	"github.com/llmcontext/mcphost/config"
	"github.com/llmcontext/mcphost/tools/shares"
	"github.com/spf13/cobra"
)

func init() {
	shareAddCmd.Flags().StringP("name", "n", "", "the name of the share")
	shareAddCmd.Flags().StringP("path", "p", "", "the path to the directory to share (default: current directory)")
	shareCmd.AddCommand(shareAddCmd)
}

var shareAddCmd = &cobra.Command{
	Use:   "add",
	Short: "share directory",
	Long:  "add a directory to the mcphost server to share with the client",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		// retrieve the name and path from the flags
		name, _ := cmd.Flags().GetString("name")
		path, _ := cmd.Flags().GetString("path")
		if name == "" {
			fmt.Println("name is required")
			return
		}
		if path == "" {
			path, err = os.Getwd()
			if err != nil {
				fmt.Println("error getting current directory:", err)
				os.Exit(1)
			}
		}
		path, err = convertRelativePathToAbsolute(path)
		if err != nil {
			fmt.Println("error converting relative path to absolute:", err)
			os.Exit(1)
		}

		// check if the path exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Println("path does not exist:", path)
			os.Exit(1)
		}
		fmt.Println("name:", name)
		fmt.Println("path:", path)

		shareTools := shares.NewShareTools(config.DefaultConfigurationDirectory)

		err = shareTools.AddShare(name, path)
		if err != nil {
			fmt.Println("error adding share:", err)
			os.Exit(1)
		}
	},
}
