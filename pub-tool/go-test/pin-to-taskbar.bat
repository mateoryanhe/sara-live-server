@echo off
setlocal EnableDelayedExpansion
cd /d "%~dp0"

set "TARGET=%~dp0deploy.bat"
if not exist "!TARGET!" (
  echo Error: deploy script not found: !TARGET!
  exit /b 1
)

powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0..\pin-to-taskbar.ps1" ^
  -DeployScript "!TARGET!" ^
  -ShortcutName "XR-Go-Test-Deploy" ^
  -LauncherExeName "XR-Go-Deploy.exe" ^
  -Description "XR Live Go test server deploy" ^
  -IconPath "%~dp0deploy-go.ico"
