$ErrorActionPreference = 'Stop'
$root = Split-Path -Parent $PSScriptRoot
$goDir = Join-Path $root 'go-审核服'
$bat = Join-Path $goDir '一键部署.bat'
if (-not (Test-Path -LiteralPath $bat)) {
  Write-Host "Not found: $bat"
  Read-Host 'Press Enter to exit'
  exit 1
}
Set-Location -LiteralPath $goDir
& cmd /c "`"$bat`""
exit $LASTEXITCODE
