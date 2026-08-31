@echo off
cd /d "%~dp0"
powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0..\pin-to-taskbar.ps1" ^
  -DeployScript "%~dp0upload.bat" ^
  -ShortcutName "XR-CMS-Test-Deploy" ^
  -LauncherExeName "XR-CMS-Deploy.exe" ^
  -Description "XR Live CMS test server build and upload" ^
  -IconPath "%~dp0deploy-cms.ico"
