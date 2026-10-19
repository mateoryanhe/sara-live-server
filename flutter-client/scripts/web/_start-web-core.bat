@echo off
setlocal EnableDelayedExpansion
REM Internal: called by start-web-*.bat
REM Args: %1=role label  %2=web-port  %3=chrome profile  %4=optional HTTP_LOG_DIR override

set "ROLE_LABEL=%~1"
set "WEB_PORT=%~2"
set "PROFILE_NAME=%~3"
if not "%~4"=="" set "HTTP_LOG_DIR=%~4"

if "%ROLE_LABEL%"=="" set "ROLE_LABEL=Web"
if "%WEB_PORT%"=="" set "WEB_PORT=8080"
if "%PROFILE_NAME%"=="" set "PROFILE_NAME=default"

set "CLIENT_ROOT=%~dp0..\.."
call "%~dp0..\env\setup.bat"
if errorlevel 1 (
  if exist "D:\tools\env-flutter.bat" (
    call D:\tools\env-flutter.bat
  )
  if errorlevel 1 (
    echo ERROR: load scripts\env\setup.bat or D:\tools\env-flutter.bat
    pause
    exit /b 1
  )
)

cd /d "%CLIENT_ROOT%"
if not exist "pubspec.yaml" (
  echo ERROR: flutter-client root not found: %CD%
  pause
  exit /b 1
)

where flutter >nul 2>&1
if errorlevel 1 (
  echo ERROR: flutter not in PATH. Edit scripts\env\config.bat
  pause
  exit /b 1
)

set "CHROME_PROFILE=%LOCALAPPDATA%\SaraLiveWeb\%PROFILE_NAME%"
if not exist "%CHROME_PROFILE%" mkdir "%CHROME_PROFILE%" 2>nul

if not defined API_BASE set "API_BASE=http://127.0.0.1:9443"

if not defined HTTP_LOG_ENABLED set "HTTP_LOG_ENABLED=true"
if not defined HTTP_LOG_DIR (
  set "HTTP_LOG_DIR=%LOCALAPPDATA%\SaraLiveWeb\http-logs\%PROFILE_NAME%"
)
if not exist "%HTTP_LOG_DIR%" mkdir "%HTTP_LOG_DIR%" 2>nul
set "HTTP_LOG_DIR_FWD=%HTTP_LOG_DIR:\=/%"

call :WriteLogReadme
if not defined HTTP_LOG_OPEN_FOLDER set "HTTP_LOG_OPEN_FOLDER=1"
if "%HTTP_LOG_OPEN_FOLDER%"=="1" start "" explorer "%HTTP_LOG_DIR%"

echo ========================================
echo   Sara Live Flutter Web - %ROLE_LABEL%
echo ========================================
echo Project : %cd%
echo Device  : chrome
echo Port    : %WEB_PORT%
echo Profile : %CHROME_PROFILE%
echo API     : %API_BASE%
echo HttpLog : enabled=%HTTP_LOG_ENABLED%
echo Log dir : %HTTP_LOG_DIR%
echo.
echo WHERE TO SEE HTTP LOGS (Web):
echo   1^) THIS cmd window - blocks appear after each API call
echo   2^) Chrome F12 - Console tab
echo   3^) Folder opened in Explorer - read WHERE-TO-SEE-LOGS.txt
echo   Edit: scripts\web\start-web-config.bat ^(API_BASE, HTTP_LOG_DIR_*^)
echo.
echo 1v1 test: run Audience (8080) then Anchor (8081); separate login + WS.
echo.

echo === flutter pub get ===
call flutter pub get
if errorlevel 1 (
  echo ERROR: flutter pub get failed
  pause
  exit /b 1
)

echo.
echo === flutter run -d chrome ===
echo Keep this window open. Ctrl+C to stop.
echo.

call flutter run -d chrome ^
  --web-port=%WEB_PORT% ^
  --dart-define=API_BASE=%API_BASE% ^
  --dart-define=DEV_ROLE=%PROFILE_NAME% ^
  --dart-define=HTTP_LOG_ENABLED=%HTTP_LOG_ENABLED% ^
  --dart-define=HTTP_LOG_DIR=%HTTP_LOG_DIR_FWD% ^
  --web-browser-flag=--user-data-dir=%CHROME_PROFILE% ^
  --web-browser-flag=--window-name=SaraLive-%PROFILE_NAME%

set ERR=%ERRORLEVEL%
echo.
if not "%ERR%"=="0" (
  echo ERROR: flutter run exited %ERR%
  pause
)
exit /b %ERR%

:WriteLogReadme
> "%HTTP_LOG_DIR%\WHERE-TO-SEE-LOGS.txt" (
  echo Sara Live - HTTP debug logs
  echo Role: %ROLE_LABEL%
  echo Directory: %HTTP_LOG_DIR%
  echo.
  echo Web ^(Chrome on Win11^):
  echo   - Keep the launcher CMD window open; each POST prints a block here.
  echo   - Chrome DevTools ^(F12^) - Console - same HTTP blocks.
  echo.
  echo Android:
  echo   - Files: http_android_YYYY-MM-DD.log when HTTP_LOG_DIR is set.
  echo.
  echo Customize:
  echo   Edit flutter-client\scripts\web\start-web-config.bat
  echo   API_BASE, HTTP_LOG_DIR_AUDIENCE / HTTP_LOG_DIR_ANCHOR or HTTP_LOG_DIR
  echo   HTTP_LOG_ENABLED=true^|false
)
exit /b 0
