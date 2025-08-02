package cmd

import (
	"fmt"
	"sort"
	// "strconv"
	"strings"

	"github.com/jdk-manager/internal/adoptium"
	"github.com/spf13/cobra"
)

var listRemoteCmd = &cobra.Command{
	Use:   "list-remote",
	Short: "List available JDK versions from Adoptium",
	Long:  `Fetch and display available JDK versions from the Eclipse Adoptium API.`,
	Run:   runListRemote,
}

var (
	showAll bool
	ltsOnly bool
)

func init() {
	listRemoteCmd.Flags().BoolVar(&showAll, "all", false, "Show all versions (including pre-release)")
	listRemoteCmd.Flags().BoolVar(&ltsOnly, "lts", false, "Show only LTS versions")
	rootCmd.AddCommand(listRemoteCmd)
}

func runListRemote(cmd *cobra.Command, args []string) {
	fmt.Println("Fetching available JDK versions from Adoptium...")

	client := adoptium.NewClient()
	releases, err := client.GetAvailableReleases()
	checkError(err)

	if len(releases) == 0 {
		fmt.Println("No JDK versions available.")
		return
	}

	// Filter releases based on flags
	var filteredReleases []adoptium.Release
	for _, release := range releases {
		// Skip pre-release versions unless --all is specified
		if !showAll && release.PreRelease {
			continue
		}

		// Show only LTS versions if --lts is specified
		if ltsOnly && !isLTSVersion(release.VersionData.Major) {
			continue
		}

		filteredReleases = append(filteredReleases, release)
	}

	// Sort by version (descending)
	sort.Slice(filteredReleases, func(i, j int) bool {
		return filteredReleases[i].VersionData.Major > filteredReleases[j].VersionData.Major
	})

	fmt.Printf("\nAvailable JDK versions:\n")
	for _, release := range filteredReleases {
		versionStr := fmt.Sprintf("%d", release.VersionData.Major)
		if release.VersionData.Minor > 0 {
			versionStr += fmt.Sprintf(".%d", release.VersionData.Minor)
		}
		if release.VersionData.Security > 0 {
			versionStr += fmt.Sprintf(".%d", release.VersionData.Security)
		}

		markers := []string{}
		if isLTSVersion(release.VersionData.Major) {
			markers = append(markers, "LTS")
		}
		if release.PreRelease {
			markers = append(markers, "pre-release")
		}

		markerStr := ""
		if len(markers) > 0 {
			markerStr = fmt.Sprintf(" (%s)", strings.Join(markers, ", "))
		}

		fmt.Printf("  %s%s\n", versionStr, markerStr)
	}

	fmt.Printf("\nUse 'jdk install <version>' to install a specific version.\n")
	if !showAll {
		fmt.Printf("Use 'jdk list-remote --all' to see all versions including pre-releases.\n")
	}
	if !ltsOnly {
		fmt.Printf("Use 'jdk list-remote --lts' to see only LTS versions.\n")
	}
}

// isLTSVersion returns true if the given major version is an LTS version
func isLTSVersion(major int) bool {
	ltsVersions := []int{8, 11, 17, 21} // Known LTS versions
	for _, lts := range ltsVersions {
		if major == lts {
			return true
		}
	}
	// Future LTS versions follow a pattern: every 3 years starting from 21
	// So: 21, 24, 27, 30, etc. (21 + 3*n where n >= 1)
	if major > 21 && (major-21)%3 == 0 {
		return true
	}
	return false
}
