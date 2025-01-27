package main

import (
	"fmt"
	"os"

	"github.com/llmcontext/mcphost/config"
	"github.com/llmcontext/mcphost/tools/shares"
	"github.com/spf13/cobra"
)

func init() {
	shareCmd.AddCommand(shareDeleteCmd)
}

var shareDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete share",
	Long:  "delete a share for current directory",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		path, err := os.Getwd()
		if err != nil {
			fmt.Println("error getting current directory:", err)
			os.Exit(1)
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
		fmt.Println("path:", path)

		shareTools := shares.NewShareTools(config.DefaultConfigurationDirectory)

		share, err := shareTools.DeleteShareForDirectory(path)
		if err != nil {
			fmt.Println("error deleting share:", err)
			os.Exit(1)
		}
		fmt.Printf("share deleted successfully: %s\n", share.ShareName)
	},
}
