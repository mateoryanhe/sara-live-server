@echo off
setlocal EnableExtensions EnableDelayedExpansion
cd /d "%~dp0"

REM XR Game Server production restart script.
REM This script does not compile, package, upload or replace any files.

call "%~dp0config.bat"

if not defined REMOTE_HOST (
    echo ERROR: REMOTE_HOST is not configured.
    exit /b 1
)
if not defined REMOTE_USER (
    echo ERROR: REMOTE_USER is not configured.
    exit /b 1
)
if not defined REMOTE_PORT set REMOTE_PORT=22
if not defined REMOTE_DIR (
    echo ERROR: REMOTE_DIR is not configured.
    exit /b 1
)
if not defined APP_NAME (
    echo ERROR: APP_NAME is not configured.
    exit /b 1
)
if not defined HOT_RESTART_AUTH (
    echo ERROR: HOT_RESTART_AUTH is not configured.
    exit /b 1
)
if not defined SSH_HOST_KEY (
    echo ERROR: SSH_HOST_KEY is not configured.
    exit /b 1
)
if not exist "%SSH_KEY_PATH%" (
    echo ERROR: SSH key does not exist: %SSH_KEY_PATH%
    exit /b 1
)
if not exist "%~dp0plink.exe" (
    echo ERROR: plink.exe does not exist: %~dp0plink.exe
    exit /b 1
)

echo ========================================
echo XR Game Server Production Restart
echo ========================================
echo Server: %REMOTE_USER%@%REMOTE_HOST%:%REMOTE_PORT%
echo App:    %REMOTE_DIR%/%APP_NAME%
echo.

echo [1/4] Testing SSH connection...
"%~dp0plink.exe" -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T -hostkey "%SSH_HOST_KEY%" %REMOTE_USER%@%REMOTE_HOST% "echo connected" >nul
if errorlevel 1 (
    echo ERROR: Unable to connect to %REMOTE_HOST%.
    echo Check the network, SSH key, port, and SSH_HOST_KEY in config.bat.
    exit /b 1
)

echo [2/4] Checking the existing server program...
"%~dp0plink.exe" -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T -hostkey "%SSH_HOST_KEY%" %REMOTE_USER%@%REMOTE_HOST% "test -x %REMOTE_DIR%/%APP_NAME%"
if errorlevel 1 (
    echo ERROR: Remote executable does not exist or is not executable: %REMOTE_DIR%/%APP_NAME%
    exit /b 1
)

set HOT_RESTART_FLUSH_TIMEOUT=60
set HOT_RESTART_EXIT_TIMEOUT=60
if exist "%LOCAL_CONFIG_PATH%" (
    for /f "usebackq tokens=2 delims=:" %%a in (`findstr /i /c:"hotRestartFlushTimeout" "%LOCAL_CONFIG_PATH%"`) do (
        for /f "tokens=1" %%b in ("%%a") do set HOT_RESTART_FLUSH_TIMEOUT=%%b
    )
    for /f "usebackq tokens=2 delims=:" %%a in (`findstr /i /c:"hotRestartExitTimeout" "%LOCAL_CONFIG_PATH%"`) do (
        for /f "tokens=1" %%b in ("%%a") do set HOT_RESTART_EXIT_TIMEOUT=%%b
    )
)
set /a HOT_RESTART_WAIT_MAX=HOT_RESTART_FLUSH_TIMEOUT+HOT_RESTART_EXIT_TIMEOUT+3

set OLD_PID=
for /f "delims=" %%i in ('plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T -hostkey "%SSH_HOST_KEY%" %REMOTE_USER%@%REMOTE_HOST% "pgrep -xo %APP_NAME% 2>/dev/null || true"') do set OLD_PID=%%i

if not defined OLD_PID goto cold_start

echo [3/4] Process PID !OLD_PID! found. Triggering graceful hot restart...
"%~dp0plink.exe" -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T -hostkey "%SSH_HOST_KEY%" %REMOTE_USER%@%REMOTE_HOST% "curl -sf -k 'https://127.0.0.1/internal/hotRestart?auth=%HOT_RESTART_AUTH%' >/dev/null"
if errorlevel 1 (
    echo ERROR: Hot restart API call failed. The running process was left untouched.
    goto show_error_log
)

echo Waiting for the new process, maximum !HOT_RESTART_WAIT_MAX! seconds...
set /a WAIT_LEFT=HOT_RESTART_WAIT_MAX

:wait_hot_restart
timeout /t 1 /nobreak >nul
set /a WAIT_LEFT-=1
set NEW_PID=
for /f "delims=" %%i in ('plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T -hostkey "%SSH_HOST_KEY%" %REMOTE_USER%@%REMOTE_HOST% "pgrep -xo %APP_NAME% 2>/dev/null || true"') do set NEW_PID=%%i
if defined NEW_PID if not "!NEW_PID!"=="!OLD_PID!" goto process_started
if !WAIT_LEFT! leq 0 (
    echo ERROR: Timed out waiting for a new process. Old PID: !OLD_PID!, current PID: !NEW_PID!
    goto show_error_log
)
goto wait_hot_restart

:cold_start
echo [3/4] No running process was found. Starting the existing server binary...
"%~dp0plink.exe" -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T -hostkey "%SSH_HOST_KEY%" %REMOTE_USER%@%REMOTE_HOST% "test -x %REMOTE_DIR%/start.sh && sudo %REMOTE_DIR%/start.sh"
if errorlevel 1 (
    echo ERROR: Cold start failed. Ensure %REMOTE_DIR%/start.sh exists; run deployment once if it is missing.
    goto show_error_log
)

set /a WAIT_LEFT=30
:wait_cold_start
timeout /t 1 /nobreak >nul
set /a WAIT_LEFT-=1
set NEW_PID=
for /f "delims=" %%i in ('plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T -hostkey "%SSH_HOST_KEY%" %REMOTE_USER%@%REMOTE_HOST% "pgrep -xo %APP_NAME% 2>/dev/null || true"') do set NEW_PID=%%i
if defined NEW_PID goto process_started
if !WAIT_LEFT! leq 0 (
    echo ERROR: Timed out waiting for the process to start.
    goto show_error_log
)
goto wait_cold_start

:process_started
echo [4/4] Process PID !NEW_PID! started. Waiting for port 443...
set /a READY_WAIT_LEFT=30

:wait_ready
"%~dp0plink.exe" -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T -hostkey "%SSH_HOST_KEY%" %REMOTE_USER%@%REMOTE_HOST% "ss -tlnp 2>/dev/null | grep -q ':443'"
if not errorlevel 1 goto restart_success
timeout /t 1 /nobreak >nul
set /a READY_WAIT_LEFT-=1
if !READY_WAIT_LEFT! leq 0 (
    echo ERROR: Process is running but port 443 did not become ready.
    goto show_error_log
)
goto wait_ready

:restart_success
echo.
echo Restart completed successfully.
echo Old PID: !OLD_PID!
echo New PID: !NEW_PID!
"%~dp0plink.exe" -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T -hostkey "%SSH_HOST_KEY%" %REMOTE_USER%@%REMOTE_HOST% "ps -p !NEW_PID! -o pid,etime,cmd --no-headers"
exit /b 0

:show_error_log
echo.
echo Recent server error log:
"%~dp0plink.exe" -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T -hostkey "%SSH_HOST_KEY%" %REMOTE_USER%@%REMOTE_HOST% "tail -30 /home/ec2-user/log/error*.log 2>/dev/null || echo '(no error log)'"
exit /b 1
