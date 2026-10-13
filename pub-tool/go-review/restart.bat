@echo off
setlocal EnableExtensions EnableDelayedExpansion
cd /d "%~dp0"

REM Hot restart first; cold restart fallback. Detect fork via process count / newest PID.

call "%~dp0config.bat"

set "PLINK=%~dp0plink.exe"
if not exist "%PLINK%" if exist "%~dp0..\go-prod\plink.exe" set "PLINK=%~dp0..\go-prod\plink.exe"
if not defined RESTART_TITLE set "RESTART_TITLE=XR Game Server Restart"
if not defined SSH_HOST_KEY (
    echo ERROR: SSH_HOST_KEY not set in config.bat
    exit /b 1
)
set "PLINK_SSH=%PLINK% -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T -hostkey "%SSH_HOST_KEY%" %REMOTE_USER%@%REMOTE_HOST%"

if not exist "%SSH_KEY_PATH%" (
    echo ERROR: SSH key missing: %SSH_KEY_PATH%
    exit /b 1
)
if not exist "%PLINK%" (
    echo ERROR: plink.exe missing: %PLINK%
    exit /b 1
)

goto :xr_restart_main

:remote_pgrep_state
set "%~1="
set "HOT_PROC_COUNT=0"
set "CUR_PGREP_O="
set "PGREP_TMP=%TEMP%\xr_pgrep_%RANDOM%.txt"
del "!PGREP_TMP!" 2>nul
%PLINK_SSH% "c=$(pgrep -x %APP_NAME% 2>/dev/null | wc -l | tr -d ' '); n=$(pgrep -x %APP_NAME% 2>/dev/null | sort -n | tail -1); o=$(pgrep -xo %APP_NAME% 2>/dev/null); printf 'XR_PGREP|%%s|%%s|%%s\n' \"$c\" \"$n\" \"$o\"" > "!PGREP_TMP!" 2>nul
for /f "usebackq tokens=1,2,3,4 delims=|" %%a in (`findstr /b "XR_PGREP" "!PGREP_TMP!" 2^>nul`) do (
    set "HOT_PROC_COUNT=%%b"
    set "%~1=%%c"
    set "CUR_PGREP_O=%%d"
)
del "!PGREP_TMP!" 2>nul
exit /b 0

:remote_hot_restart_failed
set "HOT_RESTART_FAILED=0"
set "HR_FAIL_TMP=%TEMP%\xr_hr_fail_%RANDOM%.txt"
del "!HR_FAIL_TMP!" 2>nul
%PLINK_SSH% "tail -400 /home/ec2-user/log/detail.log 2>/dev/null | grep -F 'hot_restart RestartAllServer failed' | grep -q . && echo 'XR_HR_FAIL|yes' || echo 'XR_HR_FAIL|no'" > "!HR_FAIL_TMP!" 2>nul
for /f "usebackq tokens=1,2 delims=|" %%a in (`findstr /b "XR_HR_FAIL" "!HR_FAIL_TMP!" 2^>nul`) do (
    if /i "%%b"=="yes" set "HOT_RESTART_FAILED=1"
)
del "!HR_FAIL_TMP!" 2>nul
exit /b 0

:remote_hot_fork_signal
set "HOT_FORK_SIGNAL=0"
set "HR_SIG_TMP=%TEMP%\xr_hr_sig_%RANDOM%.txt"
del "!HR_SIG_TMP!" 2>nul
%PLINK_SSH% "tail -400 /home/ec2-user/log/detail.log 2>/dev/null | grep -F 'hot_restart RestartAllServer ok' | grep -q . && echo 'XR_HR_SIG|yes' || echo 'XR_HR_SIG|no'" > "!HR_SIG_TMP!" 2>nul
for /f "usebackq tokens=1,2 delims=|" %%a in (`findstr /b "XR_HR_SIG" "!HR_SIG_TMP!" 2^>nul`) do (
    if /i "%%b"=="yes" set "HOT_FORK_SIGNAL=1"
)
del "!HR_SIG_TMP!" 2>nul
exit /b 0

:probe_hot_restart
set "HOT_RESTART_HTTP="
set "HOT_RESTART_BODY="
set "HR_TMP=%TEMP%\xr_hotrestart_%RANDOM%.txt"
del "!HR_TMP!" 2>nul
%PLINK_SSH% "http=$(curl -sk -o /tmp/xr-hr.json -w '%%{http_code}' 'https://127.0.0.1/internal/hotRestart?auth=%HOT_RESTART_AUTH%'); body=$(cat /tmp/xr-hr.json 2>/dev/null | tr -d '\r\n'); printf 'XR_HOT_RESTART|%%s|%%s\n' \"$http\" \"$body\"" > "!HR_TMP!" 2>&1
for /f "usebackq tokens=1,2,3 delims=|" %%a in (`findstr /b "XR_HOT_RESTART" "!HR_TMP!" 2^>nul`) do (
    set "HOT_RESTART_HTTP=%%b"
    set "HOT_RESTART_BODY=%%c"
)
if not defined HOT_RESTART_HTTP (
    echo WARN: Hot restart parse failed:
    type "!HR_TMP!" 2>nul
    set "HOT_RESTART_HTTP=000"
    set "HOT_RESTART_BODY=parse_error"
)
del "!HR_TMP!" 2>nul
exit /b 0

:sleep_seconds
ping -n %~1 127.0.0.1 >nul
exit /b 0

:show_error_log
echo.
echo Recent hotRestart access log:
%PLINK_SSH% "grep '/internal/hotRestart' /home/ec2-user/log/access.log 2>/dev/null | tail -8 || echo none"
echo.
echo Recent server error log:
%PLINK_SSH% "tail -15 /home/ec2-user/log/error.log 2>/dev/null || echo none"
exit /b 1

:xr_restart_main
echo ========================================
echo %RESTART_TITLE%
echo ========================================
echo Server: %REMOTE_USER%@%REMOTE_HOST%:%REMOTE_PORT%
echo App:    %REMOTE_DIR%/%APP_NAME%
echo.

echo [1/4] Testing SSH connection...
%PLINK_SSH% "echo connected" >nul
if errorlevel 1 (
    echo ERROR: SSH connection failed.
    exit /b 1
)

echo [2/4] Checking the existing server program...
%PLINK_SSH% "test -x %REMOTE_DIR%/%APP_NAME%"
if errorlevel 1 (
    echo ERROR: Missing remote binary: %REMOTE_DIR%/%APP_NAME%
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
if not defined HOT_FORK_GRACE set HOT_FORK_GRACE=90
set /a HOT_PID_WAIT_MAX=HOT_RESTART_FLUSH_TIMEOUT+HOT_FORK_GRACE
set /a HOT_TICK=0
set HOT_RESTART_FAILED=0
set HOT_FORK_SIGNAL=0
set HOT_DID_FORK=0

call :remote_pgrep_state OLD_PID
set USE_COLD=0
set DO_HOT_WAIT=0
if not defined OLD_PID set USE_COLD=1

if not "!USE_COLD!"=="0" goto after_hot_probe
echo [3/4] Process PID !OLD_PID! found. Trying graceful hot restart...
call :probe_hot_restart
if "!HOT_RESTART_HTTP!"=="200" (
    echo Hot restart accepted. Waiting for fork ^(max !HOT_PID_WAIT_MAX!s, flush=!HOT_RESTART_FLUSH_TIMEOUT!s^)...
    set /a WAIT_LEFT=HOT_PID_WAIT_MAX
    set DO_HOT_WAIT=1
    goto after_hot_probe
)
echo WARN: Hot restart rejected HTTP=!HOT_RESTART_HTTP! body=!HOT_RESTART_BODY!
echo !HOT_RESTART_BODY!| findstr /i "unauthorized" >nul
if not errorlevel 1 (
    echo ERROR: Fix HOT_RESTART_AUTH in config.bat ^(CMS preloadCfg hotRestartAuth^).
    goto show_error_log
)
echo !HOT_RESTART_BODY!| findstr /i "already in progress" >nul
if not errorlevel 1 (
    echo Hot restart already in progress, waiting for fork...
    set /a WAIT_LEFT=HOT_PID_WAIT_MAX
    set DO_HOT_WAIT=1
    goto after_hot_probe
)
set USE_COLD=1

:after_hot_probe
if "!DO_HOT_WAIT!"=="1" goto wait_hot_loop
if "!USE_COLD!"=="1" goto cold_restart
goto show_error_log

:wait_hot_loop
call :sleep_seconds 1
set /a WAIT_LEFT-=1
set /a HOT_TICK+=1
call :remote_pgrep_state NEW_PID
call :remote_hot_restart_failed
if "!HOT_RESTART_FAILED!"=="1" (
    echo ERROR: detail.log reports RestartAllServer failure; skip cold kill. Fix server/hot-restart first.
    goto show_error_log
)
call :remote_hot_fork_signal
if "!HOT_FORK_SIGNAL!"=="1" if defined NEW_PID if not "!NEW_PID!"=="!OLD_PID!" (
    echo Hot fork OK: detail.log RestartAllServer ok, PID !NEW_PID! ^(was !OLD_PID!^)
    set HOT_DID_FORK=1
    set DO_HOT_WAIT=0
    goto process_started
)
if "!HOT_FORK_SIGNAL!"=="1" if !HOT_PROC_COUNT! geq 2 (
    echo Hot fork OK: !HOT_PROC_COUNT! processes, newest PID !NEW_PID! ^(old !OLD_PID! exits in ~!HOT_RESTART_EXIT_TIMEOUT!s^)
    set HOT_DID_FORK=1
    set DO_HOT_WAIT=0
    goto process_started
)
if defined HOT_PROC_COUNT if !HOT_PROC_COUNT! geq 2 (
    echo Hot fork OK: !HOT_PROC_COUNT! processes, newest PID !NEW_PID!
    set HOT_DID_FORK=1
    set DO_HOT_WAIT=0
    goto process_started
)
set /a MOD=!HOT_TICK! %% 5
if !MOD!==0 (
    echo   ... waiting !WAIT_LEFT!s left, count=!HOT_PROC_COUNT!, newest=!NEW_PID!, pgrep-xo=!CUR_PGREP_O!
)
if !WAIT_LEFT! leq 0 (
    echo WARN: No fork detected in !HOT_PID_WAIT_MAX!s. Falling back to cold restart...
    goto cold_restart
)
goto wait_hot_loop

:cold_restart
set DO_HOT_WAIT=0
call :remote_pgrep_state COLD_PID
if defined HOT_PROC_COUNT if !HOT_PROC_COUNT! geq 2 if defined OLD_PID (
    echo [3/4] Hot fork likely active ^(!HOT_PROC_COUNT! processes^). SIGTERM old PID !OLD_PID! only...
    %PLINK_SSH% "sudo kill -TERM !OLD_PID! 2>/dev/null; sleep 3; kill -0 !OLD_PID! 2>/dev/null && sudo kill -KILL !OLD_PID!; sleep 1"
    call :remote_pgrep_state NEW_PID
    if defined NEW_PID goto process_started
)
if defined OLD_PID (
    echo [3/4] Cold restart: stopping PID !OLD_PID! ...
    %PLINK_SSH% "sudo kill -TERM !OLD_PID! 2>/dev/null; sleep 3; pgrep -x %APP_NAME% >/dev/null 2>&1 && sudo pkill -9 -x %APP_NAME%; sleep 1"
) else (
    echo [3/4] No running process. Cold starting...
)
%PLINK_SSH% "test -x %REMOTE_DIR%/start.sh && sudo %REMOTE_DIR%/start.sh"
if errorlevel 1 (
    echo ERROR: start.sh failed. Run deploy once to create %REMOTE_DIR%/start.sh
    goto show_error_log
)

set /a WAIT_LEFT=45
set /a COLD_TICK=0
:cold_wait_loop
call :sleep_seconds 1
set /a WAIT_LEFT-=1
set /a COLD_TICK+=1
call :remote_pgrep_state NEW_PID
if defined NEW_PID goto process_started
set /a MOD=!COLD_TICK! %% 5
if !MOD!==0 echo   ... cold start waiting !WAIT_LEFT!s
if !WAIT_LEFT! leq 0 (
    echo ERROR: Timed out waiting for process after cold start.
    goto show_error_log
)
goto cold_wait_loop

:process_started
if "!DO_HOT_WAIT!"=="1" (
    echo ERROR: Hot restart wait ended without fork ^(process count !HOT_PROC_COUNT!^).
    goto show_error_log
)
echo [4/4] Process PID !NEW_PID! running. Waiting for port 443...
set /a READY_WAIT_LEFT=45
:ready_loop
%PLINK_SSH% "ss -tlnp 2>/dev/null | grep -q ':443' && curl -sk -o /dev/null -w '%%{http_code}' --max-time 3 'https://127.0.0.1/internal/hotRestart?auth=invalid' | grep -qE '^403$|^200$'"
if not errorlevel 1 goto wait_old_process_exit
call :sleep_seconds 1
set /a READY_WAIT_LEFT-=1
if !READY_WAIT_LEFT! leq 0 (
    echo ERROR: Port 443 not listening.
    goto show_error_log
)
goto ready_loop

:wait_old_process_exit
if not "!HOT_DID_FORK!"=="1" goto restart_success
if not defined OLD_PID goto restart_success
if not defined NEW_PID goto restart_success
echo [5/5] Phase2: waiting for old PID !OLD_PID! to exit ^(max !HOT_RESTART_EXIT_TIMEOUT!s + 15s grace^)...
set /a OLD_EXIT_WAIT_LEFT=HOT_RESTART_EXIT_TIMEOUT+15
:old_exit_loop
%PLINK_SSH% "kill -0 !OLD_PID! 2>/dev/null"
if errorlevel 1 (
    echo Old PID !OLD_PID! exited.
    call :remote_pgrep_state NEW_PID
    goto restart_success
)
call :sleep_seconds 1
set /a OLD_EXIT_WAIT_LEFT-=1
set /a MOD=!OLD_EXIT_WAIT_LEFT! %% 10
if !MOD!==0 echo   ... old PID !OLD_PID! still running, !OLD_EXIT_WAIT_LEFT!s left
if !OLD_EXIT_WAIT_LEFT! leq 0 (
    echo WARN: Old PID !OLD_PID! still alive after phase2 wait; continuing ^(check detail.log^).
    goto restart_success
)
goto old_exit_loop

:restart_success
echo.
echo Restart completed successfully.
if defined OLD_PID echo Old PID: !OLD_PID!
echo New PID: !NEW_PID!
%PLINK_SSH% "ps -p !NEW_PID! -o pid,etime,cmd --no-headers"
exit /b 0
