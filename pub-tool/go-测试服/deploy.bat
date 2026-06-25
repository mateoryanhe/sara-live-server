@echo off
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
    pause
    exit /b 1
)

if not exist "%LOCAL_CONFIG_PATH%" (
    echo Error: Configuration directory does not exist: %LOCAL_CONFIG_PATH%
    pause
    exit /b 1
)

if not exist "%SSH_KEY_PATH%" (
    echo Error: SSH key does not exist: %SSH_KEY_PATH%
    pause
    exit /b 1
)

REM Check if PuTTY tools exist
if not exist "plink.exe" (
    echo Error: plink.exe does not exist in current directory
    pause
    exit /b 1
)

if not exist "pscp.exe" (
    echo Error: pscp.exe does not exist in current directory
    pause
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
    pause
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
plink.exe -ssh -i "%SSH_KEY_PATH%" -batch -hostkey "*" %REMOTE_USER%@%REMOTE_HOST% "echo 'Host key accepted'" >nul 2>&1
if %errorlevel% neq 0 (
    REM If the command failed but it's due to host key, try to connect without batch mode to accept it
    echo Connecting for the first time to accept host key...
    echo y | plink.exe -ssh -i "%SSH_KEY_PATH%" %REMOTE_USER%@%REMOTE_HOST% "exit" >nul 2>&1
)

REM Test connection
echo Testing server connection...
plink.exe -ssh -i "%SSH_KEY_PATH%" -batch %REMOTE_USER%@%REMOTE_HOST% "echo 'Connection successful'"
if %errorlevel% neq 0 (
    echo Error: Unable to connect to server %REMOTE_HOST%
    pause
    exit /b 1
)

echo Uploading deployment package to /tmp...

REM Upload to /tmp first (ec2-user always writable), avoid permission denied on app dir
pscp.exe -i "%SSH_KEY_PATH%" "%DEPLOY_PACKAGE%" %REMOTE_USER%@%REMOTE_HOST%:/tmp/deploy_package.zip
if %errorlevel% neq 0 (
    echo Error: File upload failed
    pause
    exit /b 1
)

REM Ensure remote directory exists, fix ownership (sudo 启动可能导致目录属主为 root)
echo Preparing remote directory...
plink.exe -ssh -i "%SSH_KEY_PATH%" -batch %REMOTE_USER%@%REMOTE_HOST% "sudo mkdir -p %REMOTE_DIR% && sudo chown -R %REMOTE_USER%:%REMOTE_USER% %REMOTE_DIR% && rm -f %REMOTE_DIR%/%APP_NAME% %REMOTE_DIR%/config.yaml %REMOTE_DIR%/%APP_NAME%.log %REMOTE_DIR%/deploy_package.zip"
if %errorlevel% neq 0 (
    echo Error: Failed to prepare remote directory %REMOTE_DIR%
    pause
    exit /b 1
)

REM Extract deployment package from /tmp into target directory
echo Extracting deployment package...
plink.exe -ssh -i "%SSH_KEY_PATH%" -batch %REMOTE_USER%@%REMOTE_HOST% "unzip -o /tmp/deploy_package.zip -d %REMOTE_DIR% && rm -f /tmp/deploy_package.zip"
if %errorlevel% neq 0 (
    echo Error: Remote extraction failed
    plink.exe -ssh -i "%SSH_KEY_PATH%" -batch %REMOTE_USER%@%REMOTE_HOST% "rm -f /tmp/deploy_package.zip"
    pause
    exit /b 1
)


echo.
echo ================================
echo Step 3: Connect to remote server and stop old program
echo ================================



echo Stopping old program on remote server...
REM Send SIGTERM signal (kill -15, equivalent to kill -13)
plink.exe -ssh -i "%SSH_KEY_PATH%" -batch %REMOTE_USER%@%REMOTE_HOST% "sudo  pkill -15 %APP_NAME% 2>nul || pkill -TERM %APP_NAME% 2>nul"
echo Waiting %SHUTDOWN_WAIT_TIME% seconds for graceful shutdown...
timeout /t %SHUTDOWN_WAIT_TIME% /nobreak >nul

REM Send SIGKILL signal (kill -9) to force stop any remaining processes
echo Sending force stop command...
plink.exe -ssh -i "%SSH_KEY_PATH%" -batch %REMOTE_USER%@%REMOTE_HOST% "sudo pkill -9 %APP_NAME% 2>nul || pkill -KILL %APP_NAME% 2>nul"




echo.
echo ================================
echo Step 5: Extract and start program
echo ================================

REM Set execution permissions
echo Setting execution permissions...
plink.exe -ssh -i "%SSH_KEY_PATH%" -batch %REMOTE_USER%@%REMOTE_HOST% "chmod +x %REMOTE_DIR%/%APP_NAME%"

REM No need to create session directory
REM Skipping session directory creation as per requirement

REM Start program
echo Starting program...
plink.exe -ssh -i "%SSH_KEY_PATH%" -batch %REMOTE_USER%@%REMOTE_HOST% "sudo nohup %REMOTE_DIR%/%APP_NAME% > /dev/null 2>&1  &"

echo.
echo ================================
echo Deployment Complete!
echo ================================
echo Server: %REMOTE_HOST%
echo Program: %APP_NAME% is running
echo Log file: %REMOTE_DIR%/%APP_NAME%.log
echo.

REM Clean up local temporary files
if exist "%TEMP_DIR%" rmdir /s /q "%TEMP_DIR%"
if exist "%DEPLOY_PACKAGE%" del "%DEPLOY_PACKAGE%"

echo Verifying if program is running...
timeout /t 5 /nobreak >nul
plink.exe -ssh -i "%SSH_KEY_PATH%" -batch %REMOTE_USER%@%REMOTE_HOST% "sudo ps aux | grep %APP_NAME% | grep -v grep"

echo.
echo Press any key to exit...
pause >nul