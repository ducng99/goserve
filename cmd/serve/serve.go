package serve

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/ducng99/goserve/internal/logger"
	"github.com/ducng99/goserve/internal/server"
	"github.com/ducng99/goserve/internal/tmpl/dirview/themes"
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

	dirViewTheme, err := cmd.Flags().GetString("index-theme")
	if err != nil {
		logger.Fatalf("Error getting 'index-theme' flag: %v\n", err)
	}
	if !themes.Exists(dirViewTheme) {
		cmd.Help()
		fmt.Printf("Invalid value for 'index-theme' flag: %s\n", dirViewTheme)
		os.Exit(1)
	}

	httpsEnabled, err := cmd.Flags().GetBool("https")
	if err != nil {
		logger.Fatalf("Error getting 'https' flag: %v\n", err)
	}

	sslCert, err := cmd.Flags().GetString("sslcert")
	if err != nil {
		logger.Fatalf("Error getting 'sslcert' flag: %v\n", err)
	}

	sslKey, err := cmd.Flags().GetString("sslkey")
	if err != nil {
		logger.Fatalf("Error getting 'sslkey' flag: %v\n", err)
	}

	proxyToAddr, err := cmd.Flags().GetString("proxy")
	if err != nil {
		logger.Fatalf("Error getting 'proxy' flag: %v\n", err)
	}

	proxyHeadersEnabled, err := cmd.Flags().GetBool("proxy-headers")
	if err != nil {
		logger.Fatalf("Error getting 'proxy-headers' flag: %v\n", err)
	}

	proxyIgnoreRedirect, err := cmd.Flags().GetBool("proxy-ignore-redirect")
	if err != nil {
		logger.Fatalf("Error getting 'proxy-ignore-redirect' flag: %v\n", err)
	}

	// Set up and start server
	config := server.ServerConfig{
		Host:                host,
		Port:                port,
		RootDir:             rootDir,
		CorsEnabled:         corsEnabled,
		DirViewTheme:        dirViewTheme,
		HttpsEnabled:        httpsEnabled,
		CertPath:            sslCert,
		KeyPath:             sslKey,
		ProxyToAddr:         proxyToAddr,
		ProxyHeadersEnabled: proxyHeadersEnabled,
		ProxyIgnoreRedirect: proxyIgnoreRedirect,
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
