package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jdk-manager/internal/adoptium"
	"github.com/jdk-manager/internal/jdk"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install <version>",
	Short: "Install a JDK version",
	Long: `Download and install a JDK version from Eclipse Adoptium.
	
Examples:
  jdk install 21        # Install JDK 21 (latest)
  jdk install 17.0.8    # Install specific version
  jdk install 11        # Install JDK 11 (latest)`,
	Args: cobra.ExactArgs(1),
	Run:  runInstall,
}

var (
	forceInstall bool
)

func init() {
	installCmd.Flags().BoolVarP(&forceInstall, "force", "f", false, "Force reinstall even if version exists")
	rootCmd.AddCommand(installCmd)
}

func runInstall(cmd *cobra.Command, args []string) {
	version := args[0]
	
	// Validate version format
	if !isValidVersion(version) {
		checkError(fmt.Errorf("invalid version format: %s", version))
	}

	manager, err := jdk.NewManager()
	checkError(err)

	// Check if already installed
	if !forceInstall {
		installed, err := manager.IsInstalled(version)
		checkError(err)
		
		if installed {
			fmt.Printf("JDK %s is already installed.\n", version)
			fmt.Printf("Use --force to reinstall or 'jdk use %s' to switch to it.\n", version)
			return
		}
	}

	fmt.Printf("Installing JDK %s...\n", version)

	// Get download info from Adoptium
	client := adoptium.NewClient()
	downloadInfo, err := client.GetDownloadInfo(version)
	checkError(err)

	if downloadInfo == nil {
		checkError(fmt.Errorf("JDK version %s not found", version))
	}

	// Install the JDK
	err = manager.Install(version, downloadInfo)
	checkError(err)

	fmt.Printf("✓ JDK %s installed successfully!\n", version)
	fmt.Printf("Use 'jdk use %s' to switch to this version.\n", version)
}

// isValidVersion checks if the version string is in a valid format
func isValidVersion(version string) bool {
	// Allow formats like: 21, 17.0.8, 11.0.20
	// Standard JDK versions follow major[.minor[.security]] format
	parts := strings.Split(version, ".")
	
	for _, part := range parts {
		if _, err := strconv.Atoi(part); err != nil {
			return false
		}
	}
	
	return len(parts) >= 1 && len(parts) <= 3  // Allow up to 3 parts: major.minor.security
}
