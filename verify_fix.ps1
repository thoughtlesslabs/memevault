$ErrorActionPreference = "Stop"

$tmp = New-Item -ItemType Directory -Path "$env:TEMP\envault_test_$(Get-Random)"
$env:USERPROFILE = $tmp
Write-Host "Testing in $tmp"

try {
    # Build
    Write-Host "Building..."
    Set-Location "C:\Users\jdiet\Coding\envault"
    go build -o "$tmp\envault.exe" .

    # 1. First run
    Write-Host "Running init (1st time)..."
    Set-Location "$tmp"
    & ".\envault.exe" init
    if (-not (Test-Path "secrets.jpg")) { 
        throw "Failed to create first secrets.jpg" 
    }
    Write-Host "First secrets.jpg created."

    # 2. Second run
    Write-Host "Running init (2nd time, new subfolder)..."
    $sub = New-Item -ItemType Directory -Path "$tmp\project2"
    Set-Location $sub
    & "..\envault.exe" init
    if (-not (Test-Path "secrets.jpg")) { 
        throw "Failed to create second secrets.jpg" 
    }
    Write-Host "Second secrets.jpg created."
    
    Write-Host "VERIFICATION SUCCESSFUL"
}
catch {
    Write-Error $_
    exit 1
}
finally {
    Remove-Item -Recurse -Force $tmp
}
