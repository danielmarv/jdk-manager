package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/jdk-manager/internal/jdk"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use <version>",
	Short: "Switch to a specific JDK version",
	Long: `Switch to a specific JDK version by setting JAVA_HOME and updating PATH.
This command will output the necessary export commands that you need to run.

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

	// Generate environment commands based on OS
	switch runtime.GOOS {
	case "windows":
		generateWindowsCommands(version, jdkPath)
	default: // Linux, macOS, and other Unix-like systems
		generateUnixCommands(version, jdkPath)
	}
}

func generateUnixCommands(version, jdkPath string) {
	fmt.Printf("export JAVA_HOME=\"%s\"\n", jdkPath)
	fmt.Printf("export PATH=\"$JAVA_HOME/bin:$PATH\"\n")
	fmt.Printf("echo \"✓ JDK %s is now active!\"\n", version) // Informative message after successful switch
	fmt.Printf("java -version\n") // Show current Java version
}

func generateWindowsCommands(version, jdkPath string) {
	fmt.Printf("$env:JAVA_HOME = \"%s\"\n", jdkPath)
	fmt.Printf("$env:PATH = \"$env:JAVA_HOME\\bin;$env:PATH\"\n")
	fmt.Printf("Write-Host \"✓ JDK %s is now active!\"\n", version) // Informative message after successful switch
	fmt.Printf("java -version\n") // Show current Java version
}
