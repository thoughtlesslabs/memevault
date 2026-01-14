# Envault Installer for Windows

$ErrorActionPreference = "Stop"

$installDir = "$env:USERPROFILE\.envault\bin"
if (!(Test-Path $installDir)) {
    New-Item -ItemType Directory -Force -Path $installDir | Out-Null
}

Write-Host "Installing Envault to $installDir..."

# In production, this would download from GitHub Releases.
# For local dev context, we copy the build artifact if it exists.
$localBuild = "$PWD\envault.exe"
if (Test-Path $localBuild) {
    Copy-Item $localBuild "$installDir\envault.exe" -Force
    Write-Host "Copied local build."
} else {
    Write-Host "Downloading latest release..."
    # Placeholder for actual download logic
    # Invoke-WebRequest -Uri "https://github.com/.../releases/latest/download/envault-windows-amd64.exe" -OutFile "$installDir\envault.exe"
    Write-Warning "No local build found and download is mocked. Please build first."
    exit 1
}

# Add to PATH (User)
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -notlike "*$installDir*") {
    Write-Host "Adding $installDir to User PATH..."
    [Environment]::SetEnvironmentVariable("Path", "$userPath;$installDir", "User")
    Write-Host "Success! Restart your terminal to use 'envault'."
} else {
    Write-Host "$installDir is already in your PATH."
}

Write-Host "Installation Complete! Try: envault --help"
