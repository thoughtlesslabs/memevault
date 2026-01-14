$ErrorActionPreference = "Stop"

# Use absolute path to the binary we just built
$memevault = "C:\Users\jdiet\Coding\envault\memevault.exe"
$tmpDir = "C:\Users\jdiet\Coding\envault\tmp_local_test"

# Cleanup from previous runs
if (Test-Path $tmpDir) { Remove-Item -Recurse -Force $tmpDir }
New-Item -ItemType Directory -Path $tmpDir | Out-Null

try {
    Write-Host "--- 1. Testing Safe Init (New Key) ---"
    $env:HOME = "$tmpDir\user1"
    New-Item -ItemType Directory -Path $env:HOME | Out-Null
    
    & $memevault init
    if (-not (Test-Path "$env:HOME\.memevault\keys\memevault.key")) { throw "Key not created" }
    $key1 = Get-Content "$env:HOME\.memevault\keys\memevault.key"

    Write-Host "--- 2. Testing Safe Init (Protect Existing Key) ---"
    # Run init again
    & $memevault init
    $key2 = Get-Content "$env:HOME\.memevault\keys\memevault.key"
    
    # Compare raw content
    if ($key1 -ne $key2) { throw "CRITICAL: Init overwrote existing key!" }
    Write-Host "Success: Key preserved."

    Write-Host "--- 3. Testing Set/Get ---"
    & $memevault set TEST_KEY "Hello World"
    $val = (& $memevault get TEST_KEY).Trim()
    if ($val -ne "Hello World") { throw "Get returned wrong value: $val" }
    
    $list = (& $memevault get)
    if ($list -notmatch "TEST_KEY=Hello World") { throw "List missing key" }
    Write-Host "Success: Set/Get works."

    Write-Host "--- 4. Testing Access Control ---"
    # Create User 2
    $env:HOME = "$tmpDir\user2"
    New-Item -ItemType Directory -Path $env:HOME | Out-Null
    & $memevault init
    $pub2 = (& $memevault keys show).Trim()

    # Switch back to User 1
    $env:HOME = "$tmpDir\user1"
    
    # User 1 adds User 2
    & $memevault grant User2 $pub2
    
    # Switch to User 2 and try to read vault from User 1's directory
    $env:HOME = "$tmpDir\user2"
    $vaultPath = "$tmpDir\user1\secrets.jpg"
    
    $val2 = (& $memevault get TEST_KEY --vault $vaultPath).Trim()
    if ($val2 -ne "Hello World") { throw "User 2 failed to read shared vault" }

    Write-Host "Success: Access granted and working."
    Write-Host "ALL CHECKS PASSED"

} catch {
    Write-Host "FAILED: $_"
    exit 1
} finally {
     # Cleanup
     if (Test-Path $tmpDir) { Remove-Item -Recurse -Force $tmpDir }
}
