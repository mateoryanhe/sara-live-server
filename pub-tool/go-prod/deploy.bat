@echo off
setlocal EnableDelayedExpansion
REM Set script to execute in current directory
cd /d "%~dp0"

REM XR Game Server Auto Deployment Script
REM Purpose: Compile Go program and deploy to remote server using PuTTY tools

REM Load configuration
call "%~dp0config.bat"

echo Starting deployment of XR Game Server...
echo Target Server: %REMOTE_HOST%
echo Target User: %REMOTE_USER%
echo Target Directory: %REMOTE_DIR%
echo.

REM Check necessary files and directories
if not exist "%LOCAL_GO_SRC%" (
    echo Error: Source code directory does not exist: %LOCAL_GO_SRC%
    exit /b 1
)

if not exist "%LOCAL_CONFIG_PATH%" (
    echo Error: Configuration directory does not exist: %LOCAL_CONFIG_PATH%
    exit /b 1
)

if not exist "%SSH_KEY_PATH%" (
    echo Error: SSH key does not exist: %SSH_KEY_PATH%
    exit /b 1
)

REM Check if PuTTY tools exist
if not exist "plink.exe" (
    echo Error: plink.exe does not exist in current directory
    exit /b 1
)

if not exist "pscp.exe" (
    echo Error: pscp.exe does not exist in current directory
    exit /b 1
)

REM Create build directory
if not exist "%LOCAL_BUILD_PATH%" mkdir "%LOCAL_BUILD_PATH%"

echo.
echo ================================
echo Step 1: Compile Go program for Linux
echo ================================
cd /d "%LOCAL_GO_SRC%"

REM Set GOOS environment variable to linux
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0

echo Compiling Go program...
go build -o "%LOCAL_BUILD_PATH%\%APP_NAME%" .
if %errorlevel% neq 0 (
    echo Error: Go program compilation failed
    cd /d "%~dp0"
    exit /b 1
)
echo Compilation completed!
cd /d "%~dp0"

echo.
echo ================================
echo Step 2: Prepare deployment package
echo ================================
set DEPLOY_PACKAGE=%LOCAL_BUILD_PATH%\deploy_package.zip

REM Create temporary directory
set TEMP_DIR=%LOCAL_BUILD_PATH%\temp_deploy
if exist "%TEMP_DIR%" rmdir /s /q "%TEMP_DIR%"
mkdir "%TEMP_DIR%"

REM Copy executable and configuration files to temporary directory
copy "%LOCAL_BUILD_PATH%\%APP_NAME%" "%TEMP_DIR%\%APP_NAME%" >nul
if exist "%LOCAL_CONFIG_PATH%" copy "%LOCAL_CONFIG_PATH%" "%TEMP_DIR%\config.yaml" >nul

REM Package deployment files
echo Creating deployment package...
pushd "%TEMP_DIR%"
"C:\Program Files\7-Zip\7z.exe" a -tzip "%DEPLOY_PACKAGE%" .\* >nul 2>nul
if %errorlevel% neq 0 (
    REM If 7-Zip is not available, try using PowerShell to package
    powershell -Command "Compress-Archive -Path .\* -DestinationPath '%DEPLOY_PACKAGE%' -Force" >nul 2>nul
    if %errorlevel% neq 0 (
        echo Warning: Unable to create ZIP package, trying alternative method
        echo Will transfer files directly...
    )
)
popd

echo Deployment package prepared!



echo.
echo ================================
echo Step 4: Upload new files
echo ================================



REM First, accept the host key by running a command with Auto_AcceptHostKeys option
echo Accepting host key for server...
plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T -hostkey "*" %REMOTE_USER%@%REMOTE_HOST% "echo 'Host key accepted'" >nul 2>&1
if %errorlevel% neq 0 (
    REM If the command failed but it's due to host key, try to connect without batch mode to accept it
    echo Connecting for the first time to accept host key...
    echo y | plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% %REMOTE_USER%@%REMOTE_HOST% "exit" >nul 2>&1
)

REM Test connection
echo Testing server connection...
plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T %REMOTE_USER%@%REMOTE_HOST% "echo 'Connection successful'"
if %errorlevel% neq 0 (
    echo Error: Unable to connect to server %REMOTE_HOST%
    exit /b 1
)

echo Uploading deployment package to /tmp...

REM Upload to /tmp first (ec2-user always writable), avoid permission denied on app dir
pscp.exe -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% "%DEPLOY_PACKAGE%" %REMOTE_USER%@%REMOTE_HOST%:/tmp/deploy_package.zip
if %errorlevel% neq 0 (
    echo Error: File upload failed
    exit /b 1
)

REM Ensure remote directory exists, fix ownership (sudo 启动可能导致目录属主为 root)
echo Preparing remote directory...
plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T %REMOTE_USER%@%REMOTE_HOST% "sudo mkdir -p %REMOTE_DIR% && sudo chown -R %REMOTE_USER%:%REMOTE_USER% %REMOTE_DIR% && rm -f %REMOTE_DIR%/%APP_NAME% %REMOTE_DIR%/config.yaml %REMOTE_DIR%/%APP_NAME%.log %REMOTE_DIR%/deploy_package.zip"
if %errorlevel% neq 0 (
    echo Error: Failed to prepare remote directory %REMOTE_DIR%
    exit /b 1
)

REM Extract deployment package from /tmp into target directory
echo Extracting deployment package...
plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T %REMOTE_USER%@%REMOTE_HOST% "unzip -o /tmp/deploy_package.zip -d %REMOTE_DIR% && rm -f /tmp/deploy_package.zip"
if %errorlevel% neq 0 (
    echo Error: Remote extraction failed
    plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T %REMOTE_USER%@%REMOTE_HOST% "rm -f /tmp/deploy_package.zip"
    exit /b 1
)


echo.
echo ================================
echo Step 3: Hot restart or cold start
echo ================================

echo Setting execution permissions...
plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T %REMOTE_USER%@%REMOTE_HOST% "chmod +x %REMOTE_DIR%/%APP_NAME%"

echo Creating remote start script...
plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T %REMOTE_USER%@%REMOTE_HOST% "printf '%%s\n' '#!/bin/sh' 'cd %REMOTE_DIR%' 'nohup ./%APP_NAME% >/dev/null 2>&1 &' > %REMOTE_DIR%/start.sh && chmod +x %REMOTE_DIR%/start.sh"
if %errorlevel% neq 0 (
    echo Error: Failed to create remote start script
    exit /b 1
)

echo Checking if %APP_NAME% is running...
set OLD_PID=
for /f "delims=" %%i in ('plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T %REMOTE_USER%@%REMOTE_HOST% "pgrep -xo %APP_NAME% 2>/dev/null || true"') do set OLD_PID=%%i

if defined OLD_PID goto do_hot_restart

echo No running process, cold starting...
plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T %REMOTE_USER%@%REMOTE_HOST% "sudo %REMOTE_DIR%/start.sh"
if !errorlevel! neq 0 (
    echo Error: Remote start script failed
    exit /b 1
)
goto hot_restart_done

:do_hot_restart
echo Process found ^(PID: !OLD_PID!^), triggering GoFrame hot restart...

REM Read hot restart timeouts from LOCAL_CONFIG_PATH (config.bat -> config/prod/config.yaml)
set HOT_RESTART_FLUSH_TIMEOUT=60
set HOT_RESTART_EXIT_TIMEOUT=60
for /f "usebackq tokens=2 delims=:" %%a in (`findstr /i /c:"hotRestartFlushTimeout" "%LOCAL_CONFIG_PATH%"`) do (
  for /f "tokens=1" %%b in ("%%a") do set HOT_RESTART_FLUSH_TIMEOUT=%%b
)
for /f "usebackq tokens=2 delims=:" %%a in (`findstr /i /c:"hotRestartExitTimeout" "%LOCAL_CONFIG_PATH%"`) do (
  for /f "tokens=1" %%b in ("%%a") do set HOT_RESTART_EXIT_TIMEOUT=%%b
)
set /a HOT_RESTART_WAIT_MAX=HOT_RESTART_FLUSH_TIMEOUT+HOT_RESTART_EXIT_TIMEOUT+3
echo Hot restart config from %LOCAL_CONFIG_PATH%: flush=!HOT_RESTART_FLUSH_TIMEOUT!s exit=!HOT_RESTART_EXIT_TIMEOUT!s waitMax=!HOT_RESTART_WAIT_MAX!s

plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T %REMOTE_USER%@%REMOTE_HOST% "curl -sf -k 'https://127.0.0.1/internal/hotRestart?auth=%HOT_RESTART_AUTH%' || exit 1"
if !errorlevel! neq 0 (
    echo Error: Hot restart API call failed
    exit /b 1
)
echo Hot restart triggered, waiting for new process ^(max !HOT_RESTART_WAIT_MAX!s, flush=!HOT_RESTART_FLUSH_TIMEOUT!s exit=!HOT_RESTART_EXIT_TIMEOUT!s^)...
set /a WAIT_LEFT=HOT_RESTART_WAIT_MAX
:wait_hot_restart
timeout /t 1 /nobreak >nul
set /a WAIT_LEFT-=1
set NEW_PID=
for /f "delims=" %%i in ('plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T %REMOTE_USER%@%REMOTE_HOST% "pgrep -xo %APP_NAME% 2>/dev/null || true"') do set NEW_PID=%%i
if not defined NEW_PID (
    if !WAIT_LEFT! leq 0 goto hot_restart_timeout
    goto wait_hot_restart
)
if not "!NEW_PID!"=="!OLD_PID!" goto hot_restart_done
if !WAIT_LEFT! leq 0 goto hot_restart_timeout
goto wait_hot_restart

:hot_restart_timeout
echo Warning: Hot restart wait timeout, checking process state...

:hot_restart_done
echo Program start/restart command completed.

echo.
echo ================================
echo Deployment Complete!
echo ================================
echo Server: %REMOTE_HOST%
echo Program: %APP_NAME%
echo.

REM Clean up local temporary files
if exist "%TEMP_DIR%" rmdir /s /q "%TEMP_DIR%"
if exist "%DEPLOY_PACKAGE%" del "%DEPLOY_PACKAGE%"

echo Verifying if program is running...
timeout /t 2 /nobreak >nul
echo.

set REMOTE_PID=
for /f "delims=" %%i in ('plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T %REMOTE_USER%@%REMOTE_HOST% "pgrep -xo %APP_NAME%" 2^>nul') do set REMOTE_PID=%%i

if not defined REMOTE_PID (
    echo ERROR: %APP_NAME% process not found!
    echo Recent error log:
    plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T %REMOTE_USER%@%REMOTE_HOST% "tail -30 /home/ec2-user/log/error*.log 2>/dev/null || echo '(no error log)'"
    exit /b 1
)

echo [Process] Started OK
echo   PID: !REMOTE_PID!
plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T %REMOTE_USER%@%REMOTE_HOST% "ps -p !REMOTE_PID! -o pid,etime,cmd --no-headers 2>/dev/null || pgrep -af %APP_NAME%"

echo.
echo [Port 443]
plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T %REMOTE_USER%@%REMOTE_HOST% "ss -tlnp 2>/dev/null | grep ':443' || (echo not listening && exit 1)"
if errorlevel 1 (
    echo ERROR: port 443 not listening
    echo Recent error log:
    plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T %REMOTE_USER%@%REMOTE_HOST% "tail -30 /home/ec2-user/log/error*.log 2>/dev/null || echo '(no error log)'"
    exit /b 1
)

echo.
echo ================================
echo Verify OK: %APP_NAME% ^(PID: !REMOTE_PID!^) is running
echo ================================
endlocal
