# Cross-platform installer script for JDK Manager (Windows PowerShell)

$ErrorActionPreference = "Stop" # Stop on errors

$InstallDir = "$env:USERPROFILE\bin" # User-specific bin directory
$ProjectDir = (Get-Item -Path $PSScriptRoot).FullName # Directory where this script is located

Write-Host "--- JDK Manager Installer (Windows PowerShell) ---"

# 1. Build the Go binary
Write-Host "Building JDK Manager binary..."
$jdkExePath = Join-Path $ProjectDir "dist\jdk.exe"

# Ensure dist directory exists
New-Item -ItemType Directory -Force -Path (Split-Path $jdkExePath) | Out-Null

# Explicitly set CGO_ENABLED for the build command and use go build directly
$env:CGO_ENABLED = "1"
try {
    # Use go build directly instead of 'make build' to avoid 'make' dependency issues
    # Capture output to see any errors
    $buildOutput = & go build -ldflags "-X github.com/jdk-manager/cmd.version=1.0.0" -o $jdkExePath main.go 2>&1
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Go build failed. Output:`n$buildOutput"
        exit 1
    }
} catch {
    Write-Error "Failed to execute Go build command: $_"
    exit 1
} finally {
    # Unset CGO_ENABLED after build to avoid affecting subsequent commands if not needed
    Remove-Item Env:CGO_ENABLED -ErrorAction SilentlyContinue
}

if (-not (Test-Path $jdkExePath)) {
    Write-Error "Error: Build failed. 'dist\jdk.exe' not found after build attempt."
    exit 1
}

# 2. Create installation directory and copy binary
Write-Host "Creating installation directory: $InstallDir"
New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null

Write-Host "Copying 'jdk.exe' to $InstallDir..."
Copy-Item -Path $jdkExePath -Destination (Join-Path $InstallDir "jdk.exe") -Force
Write-Host "JDK Manager executable installed to $InstallDir\jdk.exe"

# 3. Add installation directory to user's PATH (persistent)
Write-Host "Adding $InstallDir to user's PATH environment variable..."
$currentPath = [System.Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -notlike "*$InstallDir*") {
    [System.Environment]::SetEnvironmentVariable("Path", "$currentPath;$InstallDir", "User")
    Write-Host "Added $InstallDir to user PATH. You may need to restart your system for this to take full effect in all applications."
} else {
    Write-Host "$InstallDir is already in user PATH. Skipping."
}

# 4. Set up PowerShell profile function for seamless 'jdk use'
Write-Host "Setting up PowerShell profile function for seamless 'jdk use'..."

$profilePath = $PROFILE
$profileDir = Split-Path $profilePath
if (-not (Test-Path $profileDir)) {
    New-Item -ItemType Directory -Force -Path $profileDir | Out-Null
}
if (-not (Test-Path $profilePath)) {
    New-Item -ItemType File -Force -Path $profilePath | Out-Null
}

$initScriptPath = Join-Path $ProjectDir "scripts\init-jdk-manager.ps1"
$profileContent = Get-Content $profilePath -Raw

# Check if the sourcing line already exists
$sourcingLine = ". `"$initScriptPath`""
if ($profileContent -notlike "*$sourcingLine*") {
    Add-Content -Path $profilePath -Value "`n# JDK Manager setup`n$sourcingLine`n"
    Write-Host "Added sourcing of '$initScriptPath' to your PowerShell profile ($profilePath)."
} else {
    Write-Host "Sourcing line for JDK Manager already exists in your PowerShell profile. Skipping."
}

Write-Host ""
Write-Host "--- Installation Complete! ---"
Write-Host "Please restart your PowerShell terminal to activate the 'jdk' command."
Write-Host "You can now use 'jdk install <version>' and 'jdk use <version>' directly."
