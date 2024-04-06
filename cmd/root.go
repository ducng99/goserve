/*
Copyright © 2024 Thomas Nguyen <tom@tomng.dev>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"r.tomng.dev/goserve/cmd/serve"
	"r.tomng.dev/goserve/internal/logger"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goserve [flags] [host:port]",
	Example: "goserve -cd /path/to/dir --https --cert cert.pem --key key.pem localhost:1337",
	Short: "Starts a web server to serve static files",
	Long:  "Starts a web server to serve static files, with options for HTTPS, directory, CORS, and more.",
	Run:   serve.HandleCommand,
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
