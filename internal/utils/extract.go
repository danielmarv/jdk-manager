package utils

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ExtractArchive extracts a tar.gz or zip archive to the specified directory
func ExtractArchive(archivePath, destDir string) (string, error) {
	if strings.HasSuffix(archivePath, ".tar.gz") || strings.HasSuffix(archivePath, ".tgz") {
		return extractTarGz(archivePath, destDir)
	} else if strings.HasSuffix(archivePath, ".zip") {
		return extractZip(archivePath, destDir)
	}
	
	return "", fmt.Errorf("unsupported archive format: %s", archivePath)
}

// extractTarGz extracts a tar.gz archive
func extractTarGz(archivePath, destDir string) (string, error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return "", fmt.Errorf("failed to open archive: %w", err)
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return "", fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	var rootDir string
	
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("failed to read tar entry: %w", err)
		}

		// Get the root directory name from the first entry
		if rootDir == "" {
			parts := strings.Split(header.Name, "/")
			if len(parts) > 0 {
				rootDir = parts[0]
			}
		}

		target := filepath.Join(destDir, header.Name)

		// Ensure the target is within destDir (security check)
		if !strings.HasPrefix(target, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return "", fmt.Errorf("invalid file path: %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return "", fmt.Errorf("failed to create directory: %w", err)
			}
		case tar.TypeReg:
			if err := extractTarFile(tr, target, os.FileMode(header.Mode)); err != nil {
				return "", fmt.Errorf("failed to extract file %s: %w", header.Name, err)
			}
		}
	}

	if rootDir == "" {
		return "", fmt.Errorf("could not determine root directory")
	}

	return filepath.Join(destDir, rootDir), nil
}

// extractZip extracts a zip archive
func extractZip(archivePath, destDir string) (string, error) {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return "", fmt.Errorf("failed to open zip archive: %w", err)
	}
	defer r.Close()

	var rootDir string

	for _, f := range r.File {
		// Get the root directory name from the first entry
		if rootDir == "" {
			parts := strings.Split(f.Name, "/")
			if len(parts) > 0 {
				rootDir = parts[0]
			}
		}

		target := filepath.Join(destDir, f.Name)

		// Ensure the target is within destDir (security check)
		if !strings.HasPrefix(target, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return "", fmt.Errorf("invalid file path: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, f.FileInfo().Mode()); err != nil {
				return "", fmt.Errorf("failed to create directory: %w", err)
			}
			continue
		}

		if err := extractZipFile(f, target); err != nil {
			return "", fmt.Errorf("failed to extract file %s: %w", f.Name, err)
		}
	}

	if rootDir == "" {
		return "", fmt.Errorf("could not determine root directory")
	}

	return filepath.Join(destDir, rootDir), nil
}

// extractTarFile extracts a single file from a tar archive
func extractTarFile(tr *tar.Reader, target string, mode os.FileMode) error {
	// Create parent directories
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return err
	}

	f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, mode)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, tr)
	return err
}

// extractZipFile extracts a single file from a zip archive
func extractZipFile(f *zip.File, target string) error {
	// Create parent directories
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return err
	}

	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	outFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, f.FileInfo().Mode())
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, rc)
	return err
}
