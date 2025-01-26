package main

import (
	"fmt"
	"os"

	"github.com/llmcontext/gomcp"
	"github.com/llmcontext/mcphost/config"
	"github.com/spf13/cobra"
)

var (
	debug   bool
	rootCmd = &cobra.Command{
		Use:   "mcphost",
		Short: "Model Context Protocol host",
		Run: func(cmd *cobra.Command, args []string) {
			// we read the configuration file
			conf, err := config.LoadMcpHostConfiguration()
			if err != nil {
				fmt.Println("Error reading configuration file:", err)
				os.Exit(1)
			}

			// we create the MCP server
			mcpServer, err := gomcp.NewMcpServer(config.ServerName, config.ServerVersion)
			if err != nil {
				fmt.Println("Error creating MCP server:", err)
				os.Exit(1)
			}
			mcpServer.SetDebugLevel("debug", "mcphost.log")

			// start the server
			// start the server
			transport := mcpServer.StdioTransport()
			err = mcpServer.Start(transport)
			if err != nil {
				fmt.Println("Error starting MCP server:", err)
				os.Exit(1)
			}
		},
	}
)

func init() {
	rootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug mode")
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
