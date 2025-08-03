package jdk

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jdk-manager/internal/adoptium"
	"github.com/jdk-manager/internal/utils"
	"github.com/mitchellh/go-homedir"
)

// Manager handles JDK installation and management
type Manager struct {
	jdksDir string
	symlinkPath string // New field for the 'current' symlink path
}

// NewManager creates a new JDK manager instance
func NewManager() (*Manager, error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	jdksDir := filepath.Join(homeDir, ".jdks")
	symlinkPath := filepath.Join(jdksDir, "current") // Symlink will be inside .jdks

	// Create .jdks directory if it doesn't exist
	if err := os.MkdirAll(jdksDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create JDKs directory: %w", err)
	}

	return &Manager{
		jdksDir: jdksDir,
		symlinkPath: symlinkPath,
	}, nil
}

// GetJDKsDir returns the JDKs installation directory
func (m *Manager) GetJDKsDir() string {
	return m.jdksDir
}

// GetSymlinkPath returns the path to the 'current' symlink
func (m *Manager) GetSymlinkPath() string {
	return m.symlinkPath
}

// ListInstalled returns a list of installed JDK versions
func (m *Manager) ListInstalled() ([]string, error) {
	entries, err := os.ReadDir(m.jdksDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read JDKs directory: %w", err)
	}

	var versions []string
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() != "current" { // Exclude the 'current' symlink directory
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

// Uninstall removes a specific JDK version
func (m *Manager) Uninstall(version string) error {
	jdkPath := filepath.Join(m.jdksDir, version)

	// Check if the directory exists
	if _, err := os.Stat(jdkPath); os.IsNotExist(err) {
		return fmt.Errorf("JDK %s is not installed at %s", version, jdkPath)
	}

	fmt.Printf("Uninstalling JDK %s from %s...\n", version, jdkPath)
	if err := os.RemoveAll(jdkPath); err != nil {
		return fmt.Errorf("failed to remove JDK %s: %w", version, err)
	}

	return nil
}

// GetCurrentActiveJDKVersion attempts to determine the currently active JDK version
// by resolving the 'current' symlink.
func (m *Manager) GetCurrentActiveJDKVersion() string {
	// Check if the symlink exists
	linkInfo, err := os.Lstat(m.symlinkPath)
	if err != nil {
		return "" // Symlink doesn't exist or error reading it
	}

	// Check if it's actually a symlink
	if linkInfo.Mode()&os.ModeSymlink == 0 {
		return "" // Not a symlink
	}

	// Resolve the symlink target
	targetPath, err := os.Readlink(m.symlinkPath)
	if err != nil {
		return "" // Error resolving symlink
	}

	// Ensure the target path is within the .jdks directory
	if !strings.HasPrefix(targetPath, m.jdksDir) {
		return "" // Symlink points outside our managed directory
	}

	// Extract version from path
	rel, err := filepath.Rel(m.jdksDir, targetPath)
	if err != nil {
		return ""
	}

	parts := strings.Split(rel, string(filepath.Separator))
	if len(parts) > 0 {
		return parts[0]
	}

	return ""
}

// GenerateSymlinkCommands generates shell commands to create/update the 'current' symlink
// and set JAVA_HOME/PATH. These commands are intended to be executed by the shell.
func (m *Manager) GenerateSymlinkCommands(targetJDKPath string) {
	symlinkPath := m.GetSymlinkPath()
	symlinkBinPath := filepath.Join(symlinkPath, "bin")

	fmt.Printf("# Commands to activate JDK %s:\n", filepath.Base(targetJDKPath))

	// Remove existing symlink if it exists
	switch runtime.GOOS {
	case "windows":
		// Use cmd /C rmdir /S /Q for robustness, even if it's not a symlink
		fmt.Printf("cmd /C rmdir /S /Q \"%s\" 2>$null\n", symlinkPath)
		// Create new symlink
		fmt.Printf("cmd /C mklink /D \"%s\" \"%s\"\n", symlinkPath, targetJDKPath)
		// Set JAVA_HOME and PATH
		fmt.Printf("$env:JAVA_HOME = \"%s\"\n", symlinkPath)
		fmt.Printf("$env:PATH = \"%s;$env:PATH\"\n", symlinkBinPath)
	default: // Linux, macOS
		fmt.Printf("rm -f \"%s\"\n", symlinkPath) // Remove existing symlink
		fmt.Printf("ln -s \"%s\" \"%s\"\n", targetJDKPath, symlinkPath) // Create new symlink
		fmt.Printf("export JAVA_HOME=\"%s\"\n", symlinkPath)
		fmt.Printf("export PATH=\"$JAVA_HOME/bin:$PATH\"\n")
	}
}

// GenerateClearEnvCommands generates shell commands to clear JAVA_HOME and remove the symlink.
func (m *Manager) GenerateClearEnvCommands() {
	symlinkPath := m.GetSymlinkPath()

	fmt.Println("# Commands to clear active JDK environment:")
	switch runtime.GOOS {
	case "windows":
		fmt.Printf("$env:JAVA_HOME = \"\"\n")
		fmt.Printf("$env:PATH = ($env:PATH -split ';') -notmatch '%s'\n", strings.ReplaceAll(symlinkPath, `\`, `\\`)) // Remove symlink path from PATH
		fmt.Printf("cmd /C rmdir /S /Q \"%s\" 2>$null\n", symlinkPath)
	default: // Linux, macOS
		fmt.Printf("unset JAVA_HOME\n")
		fmt.Printf("export PATH=$(echo $PATH | sed -e 's|%s/bin:||g')\n", symlinkPath) // Remove symlink path from PATH
		fmt.Printf("rm -f \"%s\"\n", symlinkPath)
	}
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
