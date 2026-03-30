# install.ps1 — Portman installer for Windows

$ErrorActionPreference = "Stop"

$repo = "NoaTamburrini/portman"

# Detect architecture
$arch = if ([Environment]::Is64BitOperatingSystem) {
    if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { "arm64" } else { "x86_64" }
} else {
    Write-Error "Unsupported architecture: 32-bit is not supported"
    exit 1
}

Write-Host "Installing portman for Windows/${arch}..."

# Get latest release version
$release = Invoke-RestMethod -Uri "https://api.github.com/repos/${repo}/releases/latest"
$version = $release.tag_name

if (-not $version) {
    Write-Error "Failed to get latest version"
    exit 1
}

Write-Host "Latest version: ${version}"

# Download
$archiveName = "portman_Windows_${arch}.zip"
$downloadUrl = "https://github.com/${repo}/releases/download/${version}/${archiveName}"

Write-Host "Downloading from: ${downloadUrl}"

$tmpDir = Join-Path $env:TEMP "portman-install"
if (Test-Path $tmpDir) { Remove-Item -Recurse -Force $tmpDir }
New-Item -ItemType Directory -Path $tmpDir | Out-Null

try {
    Invoke-WebRequest -Uri $downloadUrl -OutFile (Join-Path $tmpDir "portman.zip")

    # Extract
    Expand-Archive -Path (Join-Path $tmpDir "portman.zip") -DestinationPath $tmpDir -Force

    # Install to user-local bin directory
    $installDir = Join-Path $env:LOCALAPPDATA "portman"
    if (-not (Test-Path $installDir)) {
        New-Item -ItemType Directory -Path $installDir | Out-Null
    }

    $installPath = Join-Path $installDir "portman.exe"
    Copy-Item -Path (Join-Path $tmpDir "portman.exe") -Destination $installPath -Force

    # Add to PATH if not already there
    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($userPath -notlike "*$installDir*") {
        [Environment]::SetEnvironmentVariable("Path", "$userPath;$installDir", "User")
        Write-Host "Added ${installDir} to your PATH (restart your terminal to apply)"
    }

    Write-Host "portman installed successfully to ${installPath}"
    Write-Host "Run 'portman' to get started!"
} finally {
    Remove-Item -Recurse -Force $tmpDir -ErrorAction SilentlyContinue
}
