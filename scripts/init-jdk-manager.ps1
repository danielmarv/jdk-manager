# This script defines a 'jdk' function to manage JDK versions in PowerShell.
# It is designed to be sourced in your PowerShell profile ($PROFILE) by the installer.

function jdk {
    param (
        [Parameter(Mandatory=$true, Position=0)]
        [string]$Command,

        [Parameter(ValueFromRemainingArguments=$true)]
        [string[]]$Args
    )

    # This path should point to the globally installed jdk.exe
    # The installer script will ensure this path is correct.
    $jdkExePath = Join-Path $env:USERPROFILE "bin\jdk.exe" 

    if (-not (Test-Path $jdkExePath)) {
        Write-Error "JDK Manager executable not found at $jdkExePath. Please run the installer script."
        return
    }

    if ($Command -eq "use") {
        # For 'jdk use', capture output and execute it in the current session
        $commandsToExecute = & $jdkExePath $Command $Args | Out-String
        
        if ($LASTEXITCODE -ne 0) {
            # Error message already printed by jdk.exe to stderr
            return
        }
        
        Invoke-Expression $commandsToExecute
        Write-Host "✓ JDK is now active!"
        java -version # Verify the change
    } else {
        # For all other commands, just pass them through to the executable
        & $jdkExePath $Command $Args
    }
}

# Optional: Add tab completion for the 'jdk' function (basic example)
# This requires the 'completion' command to be implemented in jdk.exe
Register-ArgumentCompleter -CommandName jdk -ScriptBlock {
    param($commandName, $wordToComplete, $cursorPosition)
    
    $jdkExePath = Join-Path $env:USERPROFILE "bin\jdk.exe" # Adjust path
    if (-not (Test-Path $jdkExePath)) {
        return @()
    }

    # Get completions from the actual jdk.exe
    # Note: Cobra's completion command might need specific arguments for PowerShell
    # This is a basic example and might need refinement based on Cobra's output.
    $completions = & $jdkExePath completion powershell -- "$wordToComplete"
    return $completions | ForEach-Object {
        [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
    }
}
