package utils

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExtractArchive_UnsupportedFormat(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "extract-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	_, err = ExtractArchive("test.rar", tempDir)
	if err == nil {
		t.Fatal("Expected error for unsupported format")
	}

	if !strings.Contains(err.Error(), "unsupported archive format") {
		t.Fatalf("Expected unsupported format error, got: %v", err)
	}
}

func TestExtractTarGz(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "extract-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test tar.gz file
	tarGzPath := filepath.Join(tempDir, "test.tar.gz")
	if err := createTestTarGz(tarGzPath); err != nil {
		t.Fatalf("Failed to create test tar.gz: %v", err)
	}

	// Extract the archive
	extractedPath, err := ExtractArchive(tarGzPath, tempDir)
	if err != nil {
		t.Fatalf("Failed to extract tar.gz: %v", err)
	}

	// Verify extraction
	if !strings.HasPrefix(extractedPath, tempDir) {
		t.Fatalf("Extracted path should be within temp directory")
	}

	// Check if test file exists
	testFilePath := filepath.Join(extractedPath, "test.txt")
	if _, err := os.Stat(testFilePath); os.IsNotExist(err) {
		t.Fatalf("Test file should exist after extraction")
	}

	// Verify file content
	content, err := os.ReadFile(testFilePath)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	if string(content) != "Hello, World!" {
		t.Fatalf("Expected 'Hello, World!', got '%s'", string(content))
	}
}

func TestExtractZip(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "extract-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test zip file
	zipPath := filepath.Join(tempDir, "test.zip")
	if err := createTestZip(zipPath); err != nil {
		t.Fatalf("Failed to create test zip: %v", err)
	}

	// Extract the archive
	extractedPath, err := ExtractArchive(zipPath, tempDir)
	if err != nil {
		t.Fatalf("Failed to extract zip: %v", err)
	}

	// Verify extraction
	if !strings.HasPrefix(extractedPath, tempDir) {
		t.Fatalf("Extracted path should be within temp directory")
	}

	// Check if test file exists
	testFilePath := filepath.Join(extractedPath, "test.txt")
	if _, err := os.Stat(testFilePath); os.IsNotExist(err) {
		t.Fatalf("Test file should exist after extraction")
	}

	// Verify file content
	content, err := os.ReadFile(testFilePath)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	if string(content) != "Hello, World!" {
		t.Fatalf("Expected 'Hello, World!', got '%s'", string(content))
	}
}

// Helper function to create a test tar.gz file
func createTestTarGz(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	gzWriter := gzip.NewWriter(file)
	defer gzWriter.Close()

	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	// Add a directory
	header := &tar.Header{
		Name:     "test-jdk/",
		Mode:     0755,
		Typeflag: tar.TypeDir,
	}
	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}

	// Add a test file
	content := "Hello, World!"
	header = &tar.Header{
		Name: "test-jdk/test.txt",
		Mode: 0644,
		Size: int64(len(content)),
	}
	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}

	_, err = tarWriter.Write([]byte(content))
	return err
}

// Helper function to create a test zip file
func createTestZip(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	// Add a directory
	_, err = zipWriter.Create("test-jdk/")
	if err != nil {
		return err
	}

	// Add a test file
	writer, err := zipWriter.Create("test-jdk/test.txt")
	if err != nil {
		return err
	}

	_, err = io.WriteString(writer, "Hello, World!")
	return err
}
