@echo off
setlocal
cd /d "%~dp0"
if not exist "live-hot-restart-gw.exe" call "%~dp0build.bat"
if not exist "config.json" (
  echo Copy config.example.json to config.json and set auth.
  exit /b 1
)
live-hot-restart-gw.exe
