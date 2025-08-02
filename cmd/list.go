package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jdk-manager/internal/jdk"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed JDK versions",
	Long:  `List all JDK versions currently installed in ~/.jdks directory.`,
	Run:   runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) {
	manager, err := jdk.NewManager()
	checkError(err)

	versions, err := manager.ListInstalled()
	checkError(err)

	if len(versions) == 0 {
		fmt.Println("No JDK versions installed.")
		fmt.Printf("Install a JDK version with: %s install <version>\n", os.Args[0])
		return
	}

	// Sort versions
	sort.Strings(versions)

	// Get current version
	currentVersion := getCurrentVersion(manager)

	fmt.Println("Installed JDK versions:")
	for _, version := range versions {
		marker := "  "
		if version == currentVersion {
			marker = "* " // Mark current version
		}
		fmt.Printf("%s%s\n", marker, version)
	}

	if currentVersion != "" {
		fmt.Printf("\nCurrent: %s\n", currentVersion)
	}
}

// getCurrentVersion attempts to determine the currently active JDK version
func getCurrentVersion(manager *jdk.Manager) string {
	javaHome := os.Getenv("JAVA_HOME")
	if javaHome == "" {
		return ""
	}

	jdksDir := manager.GetJDKsDir()
	if !strings.HasPrefix(javaHome, jdksDir) {
		return ""
	}

	// Extract version from path
	rel, err := filepath.Rel(jdksDir, javaHome)
	if err != nil {
		return ""
	}

	parts := strings.Split(rel, string(filepath.Separator))
	if len(parts) > 0 {
		return parts[0]
	}

	return ""
}
