@echo off
cd /d "%~dp0"
powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0pin-all-taskbar.ps1"
exit /b %ERRORLEVEL%
