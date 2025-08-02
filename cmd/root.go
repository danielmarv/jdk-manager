package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "1.0.0"
	rootCmd = &cobra.Command{
		Use:   "jdk",
		Short: "A cross-platform JDK version manager",
		Long: `JDK Manager is a cross-platform command-line tool for managing multiple Java Development Kit (JDK) versions.
Similar to nvm for Node.js, it allows you to easily install, switch between, and manage different JDK versions.

Features:
- Install JDKs from Eclipse Adoptium
- List installed and available JDK versions  
- Switch between JDK versions seamlessly
- Cross-platform support (Linux, macOS, Windows)
- Environment variable management`,
		Version: version,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add version flag
	rootCmd.Flags().BoolP("version", "v", false, "Show version information")
	
	// Customize help
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "help [command]",
		Short:  "Help about any command",
		Hidden: true,
	})
}

// checkError is a helper function to handle errors consistently
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
