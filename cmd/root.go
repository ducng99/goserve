package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"r.tomng.dev/goserve/cmd/serve"
	"r.tomng.dev/goserve/internal/logger"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     fmt.Sprintf("goserve [flags] [host:port]\n\nDefault host:port is \"%s:%s\"", serve.DefaultListenHost, serve.DefaultListenPort),
	Example: "goserve -cd /path/to/dir --https --sslcert full-cert.crt --sslkey private-key.key localhost:8443",
	Short:   "Starts a web server to serve static files",
	Long:    "Starts a web server to serve static files, with options for HTTPS, directory, CORS, and more.",
	Run:     serve.HandleCommand,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	rootCmd.Version = version

	err := rootCmd.Execute()
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

func init() {
	flags := rootCmd.Flags()

	// Generic server configs
	flags.StringP("dir", "d", ".", "Directory to serve")
	flags.BoolP("cors", "c", false, "Set CORS headers")
	flags.String("index-theme", "pretty", "Directory index page theme. Available themes: basic, pretty")

	// HTTPS
	sslFlag := flags.BoolP("ssl", "s", false, "Use HTTPS server")
	flags.BoolVar(sslFlag, "https", *sslFlag, "Alias for --ssl")
	flags.String("sslcert", "", "Path to a full certificate file")
	flags.String("sslkey", "", "Path to a private key file")
	rootCmd.MarkFlagsRequiredTogether("sslcert", "sslkey")

	// Other
	flags.BoolVar(&logger.LogNoColor, "no-color", false, "Disable colored log output")
}
