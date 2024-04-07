package serve

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"r.tomng.dev/goserve/internal/logger"
	"r.tomng.dev/goserve/internal/server"
)

const (
	DefaultListenHost = "0.0.0.0"
	DefaultListenPort = "8080"
)

// Default run function for root command.
//
// Handles flags and continue to start a server
func HandleCommand(cmd *cobra.Command, args []string) {
	host := DefaultListenHost
	port := DefaultListenPort

	if len(args) > 0 {
		host, port = parseHostPort(args[0])
	}

	rootDir := getRootDir(cmd)

	corsEnabled, err := cmd.Flags().GetBool("cors")
	if err != nil {
		logger.Fatalf("Error getting 'cors' flag: %v\n", err)
	}

	// Set up and start server
	config := server.ServerConfig{
		Host:        host,
		Port:        port,
		RootDir:     rootDir,
		CorsEnabled: corsEnabled,
	}

	config.StartServer()
}

func getRootDir(cmd *cobra.Command) string {
	userRootDir, err := cmd.Flags().GetString("dir")
	if err != nil {
		logger.Fatalf("Error getting flag: %v\n", err)
	}

	userRootDir, err = filepath.EvalSymlinks(filepath.Clean(userRootDir))
	if err != nil {
		logger.Fatalf("Error resolving directory: %v\n", err)
	}

	userRootDir, err = filepath.Abs(userRootDir)
	if err != nil {
		logger.Fatalf("Error getting absolute path: %v\n", err)
	}

	return userRootDir
}
