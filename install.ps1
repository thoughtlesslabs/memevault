$ErrorActionPreference = "Stop"

$repo = "thoughtlesslabs/memevault"
$os = "windows"
$arch = "amd64"
$binaryName = "memevault_${os}_${arch}.exe"
$url = "https://github.com/$repo/releases/latest/download/$binaryName"

$installDir = "$HOME\.memevault\bin"
if (!(Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir | Out-Null
}

$destPath = "$installDir\memevault.exe"

Write-Host "Downloading Memevault..."
Invoke-WebRequest -Uri $url -OutFile $destPath

# Add to PATH if needed
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -notlike "*$installDir*") {
    Write-Host "Adding $installDir to User PATH..."
    [Environment]::SetEnvironmentVariable("Path", "$userPath;$installDir", "User")
    $env:PATH += ";$installDir"
    Write-Host "Path updated. You may need to restart your terminal."
}

Write-Host "Successfully installed Memevault to $destPath"
Write-Host "Run 'memevault --help' to get started!"
