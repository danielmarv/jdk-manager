package cmd

import (
	"fmt"
	"os"
	// "runtime" // No longer directly used here for OS-specific commands, manager handles it

	"github.com/jdk-manager/internal/jdk"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use <version>",
	Short: "Switch to a specific JDK version",
	Long: `Switch to a specific JDK version by setting JAVA_HOME and updating PATH.
This command will output the necessary shell commands that you need to run.
These commands are typically executed by the 'jdk' shell function set up by the installer.

Examples:
  jdk use 21     # Switch to JDK 21
  jdk use 17.0.8 # Switch to specific version`,
	Args: cobra.ExactArgs(1),
	Run:  runUse,
}

func init() {
	rootCmd.AddCommand(useCmd)
}

func runUse(cmd *cobra.Command, args []string) {
	version := args[0]

	manager, err := jdk.NewManager()
	checkError(err)

	// Check if version is installed
	installed, err := manager.IsInstalled(version)
	checkError(err)

	if !installed {
		fmt.Fprintf(os.Stderr, "Error: JDK %s is not installed.\n", version) // Print error to stderr
		fmt.Fprintf(os.Stderr, "Install it with: jdk install %s\n", version)
		os.Exit(1) // Exit with error code
	}

	// Get the JDK path
	jdkPath, err := manager.GetJDKPath(version)
	checkError(err)

	// Check if this version is already active via the symlink
	currentActiveVersion := manager.GetCurrentActiveJDKVersion()
	if currentActiveVersion == version {
		fmt.Fprintf(os.Stderr, "JDK %s is already active.\n", version)
		os.Exit(0)
	}

	// Generate and print environment commands based on OS
	// These commands will be executed by the shell function (e.g., in .bashrc or $PROFILE)
	manager.GenerateSymlinkCommands(jdkPath)

	// Removed: showJavaVersion() - this is now handled by the shell function after execution
}
