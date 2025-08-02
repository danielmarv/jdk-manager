package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jdk-manager/internal/jdk"
)

func TestManagerIntegration(t *testing.T) {
	// Skip integration tests in short mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "jdk-integration-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set HOME to temp directory for testing
	originalHome := os.Getenv("HOME")
	if originalHome == "" {
		originalHome = os.Getenv("USERPROFILE") // Windows
	}
	
	os.Setenv("HOME", tempDir)
	if originalHome != "" {
		defer os.Setenv("HOME", originalHome)
	} else {
		defer os.Unsetenv("HOME")
	}

	// Test manager creation
	manager, err := jdk.NewManager()
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// Test that .jdks directory was created
	jdksDir := manager.GetJDKsDir()
	expectedDir := filepath.Join(tempDir, ".jdks")
	if jdksDir != expectedDir {
		t.Fatalf("Expected JDKs dir %s, got %s", expectedDir, jdksDir)
	}

	if _, err := os.Stat(jdksDir); os.IsNotExist(err) {
		t.Fatalf("JDKs directory should be created: %s", jdksDir)
	}

	// Test listing empty directory
	versions, err := manager.ListInstalled()
	if err != nil {
		t.Fatalf("Failed to list installed versions: %v", err)
	}

	if len(versions) != 0 {
		t.Fatalf("Expected 0 versions in empty directory, got %d", len(versions))
	}

	// Test checking non-existent version
	installed, err := manager.IsInstalled("21")
	if err != nil {
		t.Fatalf("Failed to check if version is installed: %v", err)
	}

	if installed {
		t.Fatal("Version 21 should not be installed in empty directory")
	}
}
