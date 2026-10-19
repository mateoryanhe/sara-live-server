@echo off
setlocal EnableDelayedExpansion
cd /d "%~dp0"
call "%~dp0config.bat"

echo ========================================
echo Deploy %HR_GW_NAME% ONLY (no xr-game-server)
echo Server: %REMOTE_USER%@%REMOTE_HOST%:%REMOTE_PORT%
echo Dir:    %REMOTE_DIR%
echo ========================================

if not exist "%LOCAL_HR_GW_SRC%\main.go" (
    echo ERROR: missing %LOCAL_HR_GW_SRC%
    exit /b 1
)
if not exist "%SSH_KEY_PATH%" (
    echo ERROR: SSH key missing: %SSH_KEY_PATH%
    exit /b 1
)
if not exist "plink.exe" (
    echo ERROR: plink.exe missing in %~dp0
    exit /b 1
)
if not exist "pscp.exe" (
    echo ERROR: pscp.exe missing in %~dp0
    exit /b 1
)
if not exist "%LOCAL_BUILD_PATH%" mkdir "%LOCAL_BUILD_PATH%"

set "GW_CFG=%LOCAL_BUILD_PATH%\live-hot-restart-gw.config.json"
set "GW_BIN=%LOCAL_BUILD_PATH%\%HR_GW_NAME%"

echo [1/4] Build Linux %HR_GW_NAME%...
cd /d "%LOCAL_HR_GW_SRC%"
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
go build -o "%GW_BIN%" .
if errorlevel 1 (
    echo ERROR: go build failed
    exit /b 1
)
cd /d "%~dp0"

echo [2/4] Write config.json (cold start + hot restart)...
powershell -NoProfile -Command ^
  "$j=@{listen='%HR_GW_LISTEN%';auth='%HOT_RESTART_AUTH%';hotRestartURL='https://127.0.0.1/internal/hotRestart';minIntervalSec=60;coldStartScript='%REMOTE_DIR%/start.sh';appName='%APP_NAME%';hotRestartFallbackCold=$true;coldStartWaitSec=45}|ConvertTo-Json -Compress; [IO.File]::WriteAllText('%GW_CFG%',$j)"
if not exist "%GW_CFG%" (
    echo ERROR: failed to write %GW_CFG%
    exit /b 1
)

echo [3/4] Upload binary + config.json...
plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T -hostkey "%SSH_HOST_KEY%" %REMOTE_USER%@%REMOTE_HOST% "echo ok" >nul
if errorlevel 1 (
    echo ERROR: SSH failed
    exit /b 1
)
pscp.exe -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% "%GW_BIN%" %REMOTE_USER%@%REMOTE_HOST%:%REMOTE_DIR%/%HR_GW_NAME%
if errorlevel 1 (
    echo ERROR: upload binary failed
    exit /b 1
)
pscp.exe -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% "%GW_CFG%" %REMOTE_USER%@%REMOTE_HOST%:%REMOTE_DIR%/config.json
if errorlevel 1 (
    echo ERROR: upload config.json failed
    exit /b 1
)

plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T -hostkey "%SSH_HOST_KEY%" %REMOTE_USER%@%REMOTE_HOST% "chmod +x %REMOTE_DIR%/%HR_GW_NAME%"
if errorlevel 1 (
    echo ERROR: chmod failed
    exit /b 1
)

echo [4/4] Restart gateway...
call "%~dp0restart-hr-gw.bat"
if errorlevel 1 exit /b 1

echo.
echo ===== %HR_GW_NAME% deploy OK (xr-game-server untouched) =====
call "%~dp0probe-hr-gw.bat"
exit /b 0
