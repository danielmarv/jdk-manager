package cmd

import (
	"fmt"
	"os"

	"github.com/jdk-manager/internal/jdk"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall <version>",
	Short: "Uninstall a JDK version",
	Long: `Remove a specific JDK version from your system.
	
Examples:
  jdk uninstall 21        # Uninstall JDK 21
  jdk uninstall 17.0.8    # Uninstall specific version`,
	Args: cobra.ExactArgs(1),
	Run:  runUninstall,
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}

func runUninstall(cmd *cobra.Command, args []string) {
	version := args[0]

	manager, err := jdk.NewManager()
	checkError(err)

	// Check if version is installed before attempting to uninstall
	installed, err := manager.IsInstalled(version)
	checkError(err)

	if !installed {
		fmt.Fprintf(os.Stderr, "Error: JDK %s is not installed. Nothing to uninstall.\n", version)
		os.Exit(1)
	}

	err = manager.Uninstall(version)
	checkError(err)

	fmt.Printf("✓ JDK %s uninstalled successfully!\n", version)
}
