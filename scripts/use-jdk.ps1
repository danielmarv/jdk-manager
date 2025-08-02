# This script sources the output of 'jdk use' to update environment variables
# in the current PowerShell session.
#
# Usage: . .\scripts\use-jdk.ps1 <version>
# Example: . .\scripts\use-jdk.ps1 21

param (
    [Parameter(Mandatory=$true)]
    [string]$Version
)

$jdkExePath = ".\dist\jdk.exe" # Adjust path if your executable is elsewhere, e.g., ".\dist\windows-amd64\jdk.exe"

if (-not (Test-Path $jdkExePath)) {
    Write-Error "JDK Manager executable not found at $jdkExePath. Please build the project first."
    exit 1
}

# Execute the 'jdk use' command and capture its output
# Use Out-String to join the array of lines into a single string
$commands = & $jdkExePath use $Version | Out-String

if ($LASTEXITCODE -ne 0) {
    Write-Error "Failed to get environment commands from 'jdk use $Version'."
    exit 1
}

# Execute the captured commands in the current session
Invoke-Expression $commands

Write-Host "JDK $Version is now active in this PowerShell session."
Write-Host "Current JAVA_HOME: $env:JAVA_HOME"
java -version # Verify the change
