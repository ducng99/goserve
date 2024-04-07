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
	Example: "goserve -cd /path/to/dir --https --cert cert.pem --key key.pem localhost:1337",
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
	flags.StringP("dir", "d", ".", "Directory to serve")

	flags.Bool("https", false, "Enable HTTPS")
	flags.String("cert", "", "Path to certificate file")
	flags.String("key", "", "Path to key file")
	rootCmd.MarkFlagsRequiredTogether("cert", "key")

	flags.BoolP("cors", "c", false, "Set CORS headers")

	flags.BoolVar(&logger.LogNoColor, "nocolor", false, "Disable colored log output")
}
