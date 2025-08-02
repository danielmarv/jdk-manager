package main

import (
	"context"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	// Import your existing JDK manager logic if needed for future GUI features
	// "github.com/jdk-manager/internal/jdk"
)

//go:embed all:frontend/dist
var assets embed.FS

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// InstallCLI builds the CLI and sets up shell integration
func (a *App) InstallCLI(installPath string) (string, error) {
	runtime.LogInfo(a.ctx, fmt.Sprintf("Starting CLI installation to: %s", installPath))

	// Determine the root directory of the CLI project
	// This assumes the GUI installer is run from the root of the jdk-manager project
	// or that the CLI project is located relative to the GUI installer.
	// For simplicity, let's assume the CLI project is in the parent directory.
	cliProjectRoot := filepath.Join(filepath.Dir(os.Args[0]), "..")
	if _, err := os.Stat(filepath.Join(cliProjectRoot, "main.go")); os.IsNotExist(err) {
		// Fallback if not found in parent, try current working directory (e.g., during dev)
		cliProjectRoot, _ = os.Getwd()
	}
	runtime.LogInfo(a.ctx, fmt.Sprintf("CLI project root assumed to be: %s", cliProjectRoot))


	// 1. Build the CLI executable
	runtime.LogInfo(a.ctx, "Building CLI executable...")
	buildCmd := exec.Command("make", "build") // Assuming 'make' is available
	buildCmd.Dir = cliProjectRoot
	buildOutput, err := buildCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to build CLI: %s\n%s", err.Error(), string(buildOutput))
	}
	runtime.LogInfo(a.ctx, fmt.Sprintf("Build output:\n%s", string(buildOutput)))

	// Determine source path of the built executable
	var cliExeName string
	if runtime.GOOS == "windows" {
		cliExeName = "jdk.exe"
	} else {
		cliExeName = "jdk"
	}
	sourcePath := filepath.Join(cliProjectRoot, "dist", cliExeName)

	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return "", fmt.Errorf("built executable not found at: %s", sourcePath)
	}

	// 2. Copy the executable to the target installPath
	targetPath := filepath.Join(installPath, cliExeName)
	runtime.LogInfo(a.ctx, fmt.Sprintf("Copying %s to %s...", sourcePath, targetPath))
	
	// Ensure the target directory exists
	if err := os.MkdirAll(installPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create installation directory: %w", err)
	}

	// Copy the file
	input, err := os.ReadFile(sourcePath)
	if err != nil {
		return "", fmt.Errorf("failed to read source executable: %w", err)
	}
	err = os.WriteFile(targetPath, input, 0755) // 0755 for executable permissions
	if err != nil {
		return "", fmt.Errorf("failed to copy executable to %s: %w", targetPath, err)
	}
	
	// 3. Configure shell integration (similar to install.sh/install.ps1)
	runtime.LogInfo(a.ctx, "Configuring shell integration...")
	var shellConfigOutput string
	if runtime.GOOS == "windows" {
		psScriptPath := filepath.Join(cliProjectRoot, "scripts", "install.ps1")
		// Execute the PowerShell script. This might require admin privileges.
		// For a real installer, you'd use a tool like go-elevate or manifest to request admin.
		cmd := exec.Command("powershell.exe", "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", psScriptPath)
		cmd.Dir = cliProjectRoot // Run script from project root
		output, err := cmd.CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("failed to run Windows installer script: %s\n%s", err.Error(), string(output))
		}
		shellConfigOutput = string(output)
	} else {
		bashScriptPath := filepath.Join(cliProjectRoot, "scripts", "install.sh")
		// Execute the Bash script with sudo. This will prompt the user for password.
		cmd := exec.Command("sudo", "bash", bashScriptPath)
		cmd.Dir = cliProjectRoot // Run script from project root
		output, err := cmd.CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("failed to run Linux/macOS installer script: %s\n%s", err.Error(), string(output))
		}
		shellConfigOutput = string(output)
	}

	return fmt.Sprintf("JDK Manager CLI installed successfully to %s!\n\nShell Configuration Output:\n%s", targetPath, shellConfigOutput), nil
}

func main() {
	// Create application with options
	app := NewApp()

	err := wails.Run(&options.App{
		Title:  "JDK Manager Installer",
		Width:  800,
		Height: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
