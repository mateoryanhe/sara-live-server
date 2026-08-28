$ErrorActionPreference = 'Stop'
$root = Split-Path -Parent $PSScriptRoot
$cmsDir = Join-Path $root 'cms-测试服'
$bat = Join-Path $cmsDir 'upload.bat'
if (-not (Test-Path -LiteralPath $bat)) {
  Write-Host "Not found: $bat"
  Read-Host 'Press Enter to exit'
  exit 1
}
Set-Location -LiteralPath $cmsDir
& cmd /c "`"$bat`""
exit $LASTEXITCODE