package main

import (
	"fmt"
	"os"

	"github.com/llmcontext/mcphost/config"
	"github.com/llmcontext/mcphost/tools/shares"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(shareCmd)
}

var shareCmd = &cobra.Command{
	Use:   "share",
	Short: "list of shares",
	Long:  "list all shares",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		shareTools := shares.NewShareTools(config.DefaultConfigurationDirectory)

		shares, err := shareTools.ListShares()
		if err != nil {
			fmt.Println("error listing shares:", err)
			os.Exit(1)
		}

		for _, share := range shares {
			fmt.Printf("%s: %s\n", share.ShareName, share.Path)
		}
	},
}
