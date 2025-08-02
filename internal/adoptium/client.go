package adoptium

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	adoptiumAPIBase = "https://api.adoptium.net/v3"
)

// Client handles communication with the Adoptium API
type Client struct {
	httpClient *http.Client
}

// Release represents a JDK release from Adoptium
type Release struct {
	VersionData VersionData `json:"version_data"`
	PreRelease  bool        `json:"prerelease"`
	Binaries    []Binary    `json:"binaries"`
}

// VersionData contains version information
type VersionData struct {
	Major    int `json:"major"`
	Minor    int `json:"minor"`
	Security int `json:"security"`
	Build    int `json:"build"`
}

// Binary represents a downloadable binary
type Binary struct {
	OS           string  `json:"os"`
	Architecture string  `json:"architecture"`
	ImageType    string  `json:"image_type"`
	Package      Package `json:"package"`
}

// Package contains download information
type Package struct {
	Name string `json:"name"`
	Link string `json:"link"`
	Size int64  `json:"size"`
}

// DownloadInfo contains information needed to download a JDK
type DownloadInfo struct {
	URL      string
	Filename string
	Size     int64
}

// NewClient creates a new Adoptium API client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetAvailableReleases fetches available JDK releases from Adoptium
func (c *Client) GetAvailableReleases() ([]Release, error) {
	url := fmt.Sprintf("%s/info/available_releases", adoptiumAPIBase)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch available releases: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var apiResponse struct {
		AvailableReleases []int `json:"available_releases"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	// Convert to Release structs
	var releases []Release
	for _, version := range apiResponse.AvailableReleases {
		releases = append(releases, Release{
			VersionData: VersionData{
				Major: version,
			},
		})
	}

	return releases, nil
}

// GetDownloadInfo gets download information for a specific JDK version
func (c *Client) GetDownloadInfo(version string) (*DownloadInfo, error) {
	// Parse version to get major version
	majorVersion, err := c.parseMajorVersion(version)
	if err != nil {
		return nil, fmt.Errorf("invalid version format: %w", err)
	}

	// Get current platform info
	osName := c.getOSName()
	arch := c.getArchitecture()

	// Fetch release information
	url := fmt.Sprintf("%s/assets/feature_releases/%d/ga", adoptiumAPIBase, majorVersion)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var releases []Release
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	// Find the best matching release
	for _, release := range releases {
		// Skip if specific version requested and doesn't match
		if c.isSpecificVersion(version) && !c.matchesVersion(release, version) {
			continue
		}

		// Find matching binary for current platform
		for _, binary := range release.Binaries {
			if binary.OS == osName && 
			   binary.Architecture == arch && 
			   binary.ImageType == "jdk" {
				return &DownloadInfo{
					URL:      binary.Package.Link,
					Filename: binary.Package.Name,
					Size:     binary.Package.Size,
				}, nil
			}
		}
	}

	return nil, fmt.Errorf("no suitable JDK found for version %s on %s/%s", version, osName, arch)
}

// parseMajorVersion extracts the major version number from a version string
func (c *Client) parseMajorVersion(version string) (int, error) {
	parts := strings.Split(version, ".")
	if len(parts) == 0 {
		return 0, fmt.Errorf("empty version string")
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid major version: %s", parts[0])
	}

	return major, nil
}

// isSpecificVersion checks if the version string specifies more than just major version
func (c *Client) isSpecificVersion(version string) bool {
	return strings.Contains(version, ".")
}

// matchesVersion checks if a release matches the requested version
func (c *Client) matchesVersion(release Release, requestedVersion string) bool {
	parts := strings.Split(requestedVersion, ".")
	
	// Check major version
	if len(parts) >= 1 {
		major, _ := strconv.Atoi(parts[0])
		if release.VersionData.Major != major {
			return false
		}
	}

	// Check minor version
	if len(parts) >= 2 {
		minor, _ := strconv.Atoi(parts[1])
		if release.VersionData.Minor != minor {
			return false
		}
	}

	// Check security version
	if len(parts) >= 3 {
		security, _ := strconv.Atoi(parts[2])
		if release.VersionData.Security != security {
			return false
		}
	}

	return true
}

// getOSName returns the OS name in Adoptium API format
func (c *Client) getOSName() string {
	switch runtime.GOOS {
	case "darwin":
		return "mac"
	case "windows":
		return "windows"
	case "linux":
		return "linux"
	default:
		return runtime.GOOS
	}
}

// getArchitecture returns the architecture in Adoptium API format
func (c *Client) getArchitecture() string {
	switch runtime.GOARCH {
	case "amd64":
		return "x64"
	case "arm64":
		return "aarch64"
	case "386":
		return "x32"
	default:
		return runtime.GOARCH
	}
}
