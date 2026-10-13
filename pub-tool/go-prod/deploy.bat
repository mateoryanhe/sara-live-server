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

REM Create temporary directory (unique name avoids stale config when temp_deploy is file-locked)
set TEMP_DIR=%LOCAL_BUILD_PATH%\temp_deploy_%RANDOM%_%TIME:~6,2%%TIME:~3,2%%TIME:~0,2%
mkdir "%TEMP_DIR%"
if not exist "%TEMP_DIR%" (
    echo Error: Failed to create temp dir %TEMP_DIR%
    exit /b 1
)

REM Copy executable and configuration files to temporary directory
copy /y "%LOCAL_BUILD_PATH%\%APP_NAME%" "%TEMP_DIR%\%APP_NAME%" >nul
if not exist "%LOCAL_CONFIG_PATH%" (
    echo Error: Missing config: %LOCAL_CONFIG_PATH%
    rmdir /s /q "%TEMP_DIR%"
    exit /b 1
)
copy /y "%LOCAL_CONFIG_PATH%" "%TEMP_DIR%\config.yaml" >nul
if not exist "%TEMP_DIR%\config.yaml" (
    echo Error: config.yaml not copied into deployment package
    rmdir /s /q "%TEMP_DIR%"
    exit /b 1
)

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

echo Uploading deployment package to %REMOTE_STAGE%...

REM Upload to disk staging first, avoid the small /tmp tmpfs and app-dir permission issues
plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T %REMOTE_USER%@%REMOTE_HOST% "mkdir -p %REMOTE_STAGE%"
if %errorlevel% neq 0 (
    echo Error: Failed to prepare remote staging directory %REMOTE_STAGE%
    exit /b 1
)
pscp.exe -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% "%DEPLOY_PACKAGE%" %REMOTE_USER%@%REMOTE_HOST%:%REMOTE_STAGE%/deploy_package.zip
if %errorlevel% neq 0 (
    echo Error: File upload failed
    plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T %REMOTE_USER%@%REMOTE_HOST% "rm -f %REMOTE_STAGE%/deploy_package.zip"
    exit /b 1
)

REM Ensure remote directory exists, fix ownership (start.sh may leave root-owned files)
echo Preparing remote directory...
plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T %REMOTE_USER%@%REMOTE_HOST% "sudo mkdir -p %REMOTE_DIR% && sudo chown -R %REMOTE_USER%:%REMOTE_USER% %REMOTE_DIR% && rm -f %REMOTE_DIR%/%APP_NAME% %REMOTE_DIR%/config.yaml %REMOTE_DIR%/%APP_NAME%.log %REMOTE_DIR%/deploy_package.zip"
if %errorlevel% neq 0 (
    echo Error: Failed to prepare remote directory %REMOTE_DIR%
    plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T %REMOTE_USER%@%REMOTE_HOST% "rm -f %REMOTE_STAGE%/deploy_package.zip"
    exit /b 1
)

REM Extract non-interactively; the trap removes the remote package on success, error or interruption
echo Extracting deployment package...
plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T %REMOTE_USER%@%REMOTE_HOST% "ZIP=%REMOTE_STAGE%/deploy_package.zip; cleanup() { rm -f $ZIP; }; trap cleanup EXIT HUP INT TERM; unzip -oq $ZIP -d %REMOTE_DIR% </dev/null"
if %errorlevel% neq 0 (
    echo Error: Remote extraction failed
    exit /b 1
)


echo.
echo ================================
echo Step 3: Hot restart or cold start
echo ================================

echo Setting execution permissions...
plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T -hostkey "%SSH_HOST_KEY%" %REMOTE_USER%@%REMOTE_HOST% "chmod +x %REMOTE_DIR%/%APP_NAME%"

echo Creating remote start script...
plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T -hostkey "%SSH_HOST_KEY%" %REMOTE_USER%@%REMOTE_HOST% "printf '%%s\n' '#!/bin/sh' 'cd %REMOTE_DIR%' 'nohup ./%APP_NAME% >/dev/null 2>&1 &' > %REMOTE_DIR%/start.sh && chmod +x %REMOTE_DIR%/start.sh"
if %errorlevel% neq 0 (
    echo Error: Failed to create remote start script
    exit /b 1
)

call "%~dp0restart.bat"
if errorlevel 1 (
    echo Error: Restart after deploy failed
    exit /b 1
)
echo Program start/restart completed.

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
ping -n 3 127.0.0.1 >nul
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
