package main

import (
	"fmt"
	"os"

	"github.com/llmcontext/gomcp"
	"github.com/llmcontext/mcphost/config"
	"github.com/llmcontext/mcphost/tools/shares"
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
			mcpServerDefinition := gomcp.NewMcpServerDefinition(config.ServerName, config.ServerVersion)
			mcpServerDefinition.SetDebugLevel(conf.Logging.Level, conf.Logging.File)

			// we register the tools
			shareTools := shares.NewShareTools(config.DefaultConfigurationDirectory)
			if err := shareTools.Register(mcpServerDefinition); err != nil {
				fmt.Println("Error registering shares tools:", err)
				os.Exit(1)
			}

			mcp, err := gomcp.NewModelContextProtocolServer(mcpServerDefinition)
			if err != nil {
				fmt.Println("Error creating MCP server:", err)
				os.Exit(1)
			}

			// start the server
			transport := mcp.StdioTransport()
			err = mcp.Start(transport)
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
