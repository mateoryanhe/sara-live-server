@echo off
REM Change to the script's directory
cd /d "%~dp0"

REM Configuration parameters
set SERVER_IP=hzaicoin.fun
set SERVER_USER=root
set SERVER_PASSWORD=Appledev882116
set SERVER_PORT=22
set VUE_PROJECT_DIR=D:\go-project\xrgameserver\cms
set BUILD_OUTPUT_DIR=D:\root\cms
set BUILD_ENV=dev
set REMOTE_DIR=/root/p-cms
set ZIP_FILE=upload_package.zip

echo Starting build and upload process...
echo Vue Project Directory: %VUE_PROJECT_DIR%
echo Build Output Directory: %BUILD_OUTPUT_DIR%
echo Build Environment: %BUILD_ENV%
echo Remote Directory: %REMOTE_DIR%

REM Check if Vue project directory exists
if not exist "%VUE_PROJECT_DIR%" (
    echo Error: Vue project directory %VUE_PROJECT_DIR% does not exist
    pause
    exit /b 1
)

REM Build Vue project
echo Building Vue project with environment %BUILD_ENV%...
cd /d "%VUE_PROJECT_DIR%"

REM Check if package.json exists
if not exist "package.json" (
    echo Error: package.json not found in %VUE_PROJECT_DIR%
    pause
    exit /b 1
)

REM Determine the build command based on environment
if "%BUILD_ENV%"=="prod" (
    set BUILD_CMD=build:prod
) else (
    set BUILD_CMD=build:dev
)

echo Running build command: npm run %BUILD_CMD%
call npm run %BUILD_CMD%

if %ERRORLEVEL% NEQ 0 (
    echo Error: Vue build failed
    cd /d "%~dp0"
    pause
    exit /b 1
)

REM Return to script directory
cd /d "%~dp0"

REM Check if PuTTY tools exist in current directory
if not exist "plink.exe" (
    echo Error: plink.exe not found, please ensure plink.exe is in current directory
    pause
    exit /b 1
)

if not exist "pscp.exe" (
    echo Error: pscp.exe not found, please ensure pscp.exe is in current directory
    pause
    exit /b 1
)

REM Check if build output directory exists
if not exist "%BUILD_OUTPUT_DIR%" (
    echo Error: Build output directory %BUILD_OUTPUT_DIR% does not exist
    pause
    exit /b 1
)

REM Compress the built files
echo Compressing built files from %BUILD_OUTPUT_DIR%...
where 7z >nul 2>&1
if %ERRORLEVEL% EQU 0 (
    echo Compressing using 7-Zip...
    7z a -tzip "%ZIP_FILE%" "%BUILD_OUTPUT_DIR%\*"
) else (
    echo 7-Zip not found, using PowerShell compression...
    powershell -Command "if (Test-Path '%ZIP_FILE%') { Remove-Item '%ZIP_FILE%' }; Add-Type -AssemblyName System.IO.Compression.FileSystem; [System.IO.Compression.ZipFile]::CreateFromDirectory('%BUILD_OUTPUT_DIR%', '%ZIP_FILE%')"
)

if not exist "%ZIP_FILE%" (
    echo Error: Failed to create zip file
    pause
    exit /b 1
)

echo Deleting remote directory %REMOTE_DIR%...
plink.exe -ssh -l %SERVER_USER% -pw "%SERVER_PASSWORD%" -P %SERVER_PORT% -batch %SERVER_IP% "if [ -d '%REMOTE_DIR%' ]; then rm -rf '%REMOTE_DIR%'; echo 'Directory deleted'; else echo 'Directory does not exist'; fi"

if %ERRORLEVEL% EQU 0 (
    echo Remote directory operation successful, uploading zip file...
    
    REM Upload zip file to server
    pscp.exe -pw "%SERVER_PASSWORD%" "%ZIP_FILE%" %SERVER_USER%@%SERVER_IP%:/tmp/%ZIP_FILE%
    
    if %ERRORLEVEL% EQU 0 (
        echo Zip file uploaded successfully, extracting on server...
        
        REM Create remote directory and extract
        plink.exe -ssh -l %SERVER_USER% -pw "%SERVER_PASSWORD%" -P %SERVER_PORT% -batch %SERVER_IP% "mkdir -p '%REMOTE_DIR%' && cd '%REMOTE_DIR%' && unzip -o /tmp/%ZIP_FILE% && rm /tmp/%ZIP_FILE%"
        
        if %ERRORLEVEL% EQU 0 (
            echo Files extracted successfully!
        ) else (
            echo File extraction failed!
            REM Clean up uploaded zip file in case of extraction failure
            plink.exe -ssh -l %SERVER_USER% -pw "%SERVER_PASSWORD%" -P %SERVER_PORT% -batch %SERVER_IP% "rm /tmp/%ZIP_FILE%"
            pause
            del "%ZIP_FILE%"
            exit /b 1
        )
    ) else (
        echo Zip file upload failed!
        del "%ZIP_FILE%"
        pause
        exit /b 1
    )
) else (
    echo Delete remote directory failed!
    del "%ZIP_FILE%"
    pause
    exit /b 1
)

REM Clean up local zip file
del "%ZIP_FILE%"

echo Build and upload completed!
pause