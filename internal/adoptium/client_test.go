package adoptium

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Fatal("Client should not be nil")
	}

	if client.httpClient == nil {
		t.Fatal("HTTP client should not be nil")
	}
}

func TestParseMajorVersion(t *testing.T) {
	client := NewClient()

	tests := []struct {
		version  string
		expected int
		hasError bool
	}{
		{"21", 21, false},
		{"17.0.8", 17, false},
		{"11.0.20", 11, false},
		{"8", 8, false},
		{"", 0, true},
		{"invalid", 0, true},
		{"21.0.1.2", 21, false},
	}

	for _, test := range tests {
		result, err := client.parseMajorVersion(test.version)
		
		if test.hasError {
			if err == nil {
				t.Errorf("Expected error for version %s, but got none", test.version)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for version %s: %v", test.version, err)
			}
			if result != test.expected {
				t.Errorf("Expected %d for version %s, got %d", test.expected, test.version, result)
			}
		}
	}
}

func TestIsSpecificVersion(t *testing.T) {
	client := NewClient()

	tests := []struct {
		version  string
		expected bool
	}{
		{"21", false},
		{"17.0.8", true},
		{"11.0.20", true},
		{"8", false},
		{"21.0", true},
	}

	for _, test := range tests {
		result := client.isSpecificVersion(test.version)
		if result != test.expected {
			t.Errorf("Expected %v for version %s, got %v", test.expected, test.version, result)
		}
	}
}

func TestGetOSName(t *testing.T) {
	client := NewClient()
	osName := client.getOSName()
	
	if osName == "" {
		t.Fatal("OS name should not be empty")
	}

	// Should be one of the supported OS names
	validOSNames := []string{"linux", "mac", "windows"}
	found := false
	for _, validOS := range validOSNames {
		if osName == validOS {
			found = true
			break
		}
	}
	
	if !found {
		t.Logf("OS name '%s' is not in the standard list, but this might be expected for other platforms", osName)
	}
}

func TestGetArchitecture(t *testing.T) {
	client := NewClient()
	arch := client.getArchitecture()
	
	if arch == "" {
		t.Fatal("Architecture should not be empty")
	}

	// Should be one of the supported architectures
	validArchs := []string{"x64", "aarch64", "x32"}
	found := false
	for _, validArch := range validArchs {
		if arch == validArch {
			found = true
			break
		}
	}
	
	if !found {
		t.Logf("Architecture '%s' is not in the standard list, but this might be expected for other platforms", arch)
	}
}

func TestMatchesVersion(t *testing.T) {
	client := NewClient()

	release := Release{
		VersionData: VersionData{
			Major:    17,
			Minor:    0,
			Security: 8,
		},
	}

	tests := []struct {
		requestedVersion string
		expected         bool
	}{
		{"17", true},
		{"17.0", true},
		{"17.0.8", true},
		{"17.0.9", false},
		{"18", false},
		{"16", false},
	}

	for _, test := range tests {
		result := client.matchesVersion(release, test.requestedVersion)
		if result != test.expected {
			t.Errorf("Expected %v for version %s against release 17.0.8, got %v", 
				test.expected, test.requestedVersion, result)
		}
	}
}
