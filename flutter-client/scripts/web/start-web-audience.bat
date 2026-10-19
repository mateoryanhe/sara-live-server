@echo off
setlocal
REM Audience (caller), Chrome 8080.

if not exist "%~dp0start-web-config.bat" (
  if exist "%~dp0start-web-config.example.bat" (
    copy /Y "%~dp0start-web-config.example.bat" "%~dp0start-web-config.bat" >nul
    echo [setup] Created start-web-config.bat - edit HTTP_LOG_DIR_* then re-run if needed.
  )
)

if exist "%~dp0start-web-config.bat" call "%~dp0start-web-config.bat"

if defined HTTP_LOG_DIR_AUDIENCE set "HTTP_LOG_DIR=%HTTP_LOG_DIR_AUDIENCE%"
if not defined HTTP_LOG_DIR set "HTTP_LOG_DIR=%~dp0..\..\logs\web-audience"
if not defined HTTP_LOG_ENABLED set "HTTP_LOG_ENABLED=true"

call "%~dp0_start-web-core.bat" Audience 8080 profile-audience
