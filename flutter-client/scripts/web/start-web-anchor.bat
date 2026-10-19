@echo off
setlocal
REM Anchor (callee), Chrome 8081.

if not exist "%~dp0start-web-config.bat" (
  if exist "%~dp0start-web-config.example.bat" (
    copy /Y "%~dp0start-web-config.example.bat" "%~dp0start-web-config.bat" >nul
    echo [setup] Created start-web-config.bat - edit HTTP_LOG_DIR_* then re-run if needed.
  )
)

if exist "%~dp0start-web-config.bat" call "%~dp0start-web-config.bat"

if defined HTTP_LOG_DIR_ANCHOR set "HTTP_LOG_DIR=%HTTP_LOG_DIR_ANCHOR%"
if not defined HTTP_LOG_DIR set "HTTP_LOG_DIR=%~dp0..\..\logs\web-anchor"
if not defined HTTP_LOG_ENABLED set "HTTP_LOG_ENABLED=true"

call "%~dp0_start-web-core.bat" Anchor 8081 profile-anchor
