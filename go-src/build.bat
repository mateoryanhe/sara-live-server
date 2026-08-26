@echo off
setlocal
cd /d "%~dp0"

set "BUILD_DIR=%~dp0..\go-build"
if not exist "%BUILD_DIR%" mkdir "%BUILD_DIR%"

go build -o "%BUILD_DIR%\xr-game-server" .
if errorlevel 1 exit /b 1

echo Built: %BUILD_DIR%\xr-game-server
