package jdk

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jdk-manager/internal/adoptium"
	"github.com/jdk-manager/internal/utils"
	"github.com/mitchellh/go-homedir"
)

// Manager handles JDK installation and management
type Manager struct {
	jdksDir string
}

// NewManager creates a new JDK manager instance
func NewManager() (*Manager, error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	jdksDir := filepath.Join(homeDir, ".jdks")
	
	// Create .jdks directory if it doesn't exist
	if err := os.MkdirAll(jdksDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create JDKs directory: %w", err)
	}

	return &Manager{
		jdksDir: jdksDir,
	}, nil
}

// GetJDKsDir returns the JDKs installation directory
func (m *Manager) GetJDKsDir() string {
	return m.jdksDir
}

// ListInstalled returns a list of installed JDK versions
func (m *Manager) ListInstalled() ([]string, error) {
	entries, err := os.ReadDir(m.jdksDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read JDKs directory: %w", err)
	}

	var versions []string
	for _, entry := range entries {
		if entry.IsDir() {
			// Verify it's a valid JDK installation
			jdkPath := filepath.Join(m.jdksDir, entry.Name())
			if m.isValidJDK(jdkPath) {
				versions = append(versions, entry.Name())
			}
		}
	}

	return versions, nil
}

// IsInstalled checks if a specific JDK version is installed
func (m *Manager) IsInstalled(version string) (bool, error) {
	jdkPath := filepath.Join(m.jdksDir, version)
	
	// Check if directory exists and contains a valid JDK
	if _, err := os.Stat(jdkPath); os.IsNotExist(err) {
		return false, nil
	}

	return m.isValidJDK(jdkPath), nil
}

// GetJDKPath returns the full path to a specific JDK version
func (m *Manager) GetJDKPath(version string) (string, error) {
	jdkPath := filepath.Join(m.jdksDir, version)
	
	if !m.isValidJDK(jdkPath) {
		return "", fmt.Errorf("JDK %s is not properly installed", version)
	}

	return jdkPath, nil
}

// Install downloads and installs a JDK version
func (m *Manager) Install(version string, downloadInfo *adoptium.DownloadInfo) error {
	installPath := filepath.Join(m.jdksDir, version)
	
	// Remove existing installation if it exists
	if _, err := os.Stat(installPath); err == nil {
		if err := os.RemoveAll(installPath); err != nil {
			return fmt.Errorf("failed to remove existing installation: %w", err)
		}
	}

	// Create temporary directory for download
	tempDir, err := os.MkdirTemp("", "jdk-install-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Download the JDK archive
	archivePath := filepath.Join(tempDir, downloadInfo.Filename)
	fmt.Printf("Downloading %s...\n", downloadInfo.Filename)
	
	if err := utils.DownloadFile(downloadInfo.URL, archivePath); err != nil {
		return fmt.Errorf("failed to download JDK: %w", err)
	}

	// Extract the archive
	fmt.Println("Extracting JDK...")
	extractedPath, err := utils.ExtractArchive(archivePath, tempDir)
	if err != nil {
		return fmt.Errorf("failed to extract JDK: %w", err)
	}

	// Move to final location
	if err := os.Rename(extractedPath, installPath); err != nil {
		return fmt.Errorf("failed to move JDK to installation directory: %w", err)
	}

	// Verify installation
	if !m.isValidJDK(installPath) {
		return fmt.Errorf("JDK installation verification failed")
	}

	return nil
}

// isValidJDK checks if a directory contains a valid JDK installation
func (m *Manager) isValidJDK(jdkPath string) bool {
	// Check for java executable
	var javaExe string
	if runtime.GOOS == "windows" {
		javaExe = "java.exe"
	} else {
		javaExe = "java"
	}

	javaPath := filepath.Join(jdkPath, "bin", javaExe)
	if _, err := os.Stat(javaPath); err != nil {
		return false
	}

	// Check for javac executable (to ensure it's a JDK, not just JRE)
	var javacExe string
	if runtime.GOOS == "windows" {
		javacExe = "javac.exe"
	} else {
		javacExe = "javac"
	}

	javacPath := filepath.Join(jdkPath, "bin", javacExe)
	if _, err := os.Stat(javacPath); err != nil {
		return false
	}

	return true
}
