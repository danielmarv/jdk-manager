package jdk

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewManager(t *testing.T) {
	manager, err := NewManager()
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	if manager == nil {
		t.Fatal("Manager should not be nil")
	}

	if manager.jdksDir == "" {
		t.Fatal("JDKs directory should not be empty")
	}

	// Check if directory was created
	if _, err := os.Stat(manager.jdksDir); os.IsNotExist(err) {
		t.Fatalf("JDKs directory should be created: %s", manager.jdksDir)
	}
}

func TestGetJDKsDir(t *testing.T) {
	manager, err := NewManager()
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	jdksDir := manager.GetJDKsDir()
	if jdksDir == "" {
		t.Fatal("JDKs directory should not be empty")
	}

	if !filepath.IsAbs(jdksDir) {
		t.Fatal("JDKs directory should be absolute path")
	}
}

func TestListInstalled_EmptyDirectory(t *testing.T) {
	manager, err := NewManager()
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	versions, err := manager.ListInstalled()
	if err != nil {
		t.Fatalf("Failed to list installed versions: %v", err)
	}

	if len(versions) != 0 {
		t.Fatalf("Expected 0 versions, got %d", len(versions))
	}
}

func TestIsInstalled_NonExistent(t *testing.T) {
	manager, err := NewManager()
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	installed, err := manager.IsInstalled("21")
	if err != nil {
		t.Fatalf("Failed to check if version is installed: %v", err)
	}

	if installed {
		t.Fatal("Version 21 should not be installed")
	}
}

func TestIsValidJDK(t *testing.T) {
	manager, err := NewManager()
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// Test with non-existent directory
	if manager.isValidJDK("/non/existent/path") {
		t.Fatal("Non-existent path should not be valid JDK")
	}

	// Test with empty directory
	tempDir, err := os.MkdirTemp("", "jdk-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	if manager.isValidJDK(tempDir) {
		t.Fatal("Empty directory should not be valid JDK")
	}
}
