@echo off
setlocal EnableDelayedExpansion
cd /d "%~dp0"

call "%~dp0config.bat"

echo ========================================
echo   Go local build + start
echo ========================================
echo GO_SRC : %LOCAL_GO_SRC%
echo BUILD  : %LOCAL_BUILD_PATH%\%APP_EXE%
echo CONFIG : %LOCAL_CONFIG_FILE%
echo PORT   : %LOCAL_PORT%
echo.

if not exist "%LOCAL_GO_SRC%" (
  echo ERROR: go-src not found: %LOCAL_GO_SRC%
  exit /b 1
)

if not exist "%LOCAL_CONFIG_FILE%" (
  echo ERROR: config not found: %LOCAL_CONFIG_FILE%
  exit /b 1
)

where go >nul 2>&1
if errorlevel 1 (
  echo ERROR: go not found in PATH.
  exit /b 1
)

if not exist "%LOCAL_BUILD_PATH%" mkdir "%LOCAL_BUILD_PATH%"

echo [1/3] Stop process listening on port %LOCAL_PORT% ...
set KILLED=0
for /f "tokens=5" %%p in ('netstat -ano ^| findstr ":%LOCAL_PORT% " ^| findstr "LISTENING"') do (
  echo   killing PID %%p
  taskkill /F /PID %%p >nul 2>&1
  set KILLED=1
)
if "!KILLED!"=="0" (
  echo   no listener on %LOCAL_PORT%
) else (
  ping -n 2 127.0.0.1 >nul
)
echo.

echo [2/3] go build -o %LOCAL_BUILD_PATH%\%APP_EXE%
cd /d "%LOCAL_GO_SRC%"
if errorlevel 1 (
  echo ERROR: cannot cd to %LOCAL_GO_SRC%
  exit /b 1
)

set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=0

go build -o "%LOCAL_BUILD_PATH%\%APP_EXE%" .
if errorlevel 1 (
  echo ERROR: go build failed
  cd /d "%~dp0"
  exit /b 1
)
echo   build ok
cd /d "%~dp0"
echo.

echo [3/3] Start server (cwd=%LOCAL_CONFIG_DIR%)
if not exist "%LOCAL_BUILD_PATH%\%APP_EXE%" (
  echo ERROR: binary missing: %LOCAL_BUILD_PATH%\%APP_EXE%
  exit /b 1
)

start "xr-game-server-local" /D "%LOCAL_CONFIG_DIR%" "%LOCAL_BUILD_PATH%\%APP_EXE%"
ping -n 3 127.0.0.1 >nul

netstat -ano | findstr ":%LOCAL_PORT% " | findstr "LISTENING" >nul
if errorlevel 1 (
  echo WARN: port %LOCAL_PORT% not listening yet; check the new console window / D:\log
) else (
  echo OK: listening on %LOCAL_PORT%
)

echo.
echo Done. Config: %LOCAL_CONFIG_FILE%
exit /b 0
