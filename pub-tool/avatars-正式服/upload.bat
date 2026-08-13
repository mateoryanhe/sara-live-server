@echo off
setlocal EnableDelayedExpansion
cd /d "%~dp0"

REM 用法: upload.bat [test|prod]  默认 test
set DEPLOY_ENV=%~1
if "%DEPLOY_ENV%"=="" set DEPLOY_ENV=test

if /i not "%DEPLOY_ENV%"=="test" if /i not "%DEPLOY_ENV%"=="prod" (
    echo Usage: %~nx0 [test^|prod]
    echo   test - 上传头像到测试服目录
    echo   prod - 上传头像到生产目录
    pause
    exit /b 1
)

call "%~dp0config.bat"

if /i "%DEPLOY_ENV%"=="prod" (
    set REMOTE_DIR=%REMOTE_DIR_PROD%
    set ZIP_FILE=avatars_upload_prod.zip
) else (
    set REMOTE_DIR=%REMOTE_DIR_TEST%
    set ZIP_FILE=avatars_upload_test.zip
)

echo ========================================
echo Avatar Upload
echo Environment: %DEPLOY_ENV%
echo Local directory: %LOCAL_AVATARS_DIR%
echo Target server: %REMOTE_USER%@%REMOTE_HOST%:%REMOTE_PORT%
echo Remote directory: %REMOTE_DIR%
echo SSH key: %SSH_KEY_PATH%
echo ========================================
echo.

if not exist "%LOCAL_AVATARS_DIR%" (
    echo Error: Local avatars directory does not exist: %LOCAL_AVATARS_DIR%
    pause
    exit /b 1
)

dir /b "%LOCAL_AVATARS_DIR%\*.*" >nul 2>&1
if errorlevel 1 (
    echo Error: No files found in %LOCAL_AVATARS_DIR%
    pause
    exit /b 1
)

if not exist "%SSH_KEY_PATH%" (
    echo Error: SSH key file does not exist: %SSH_KEY_PATH%
    echo Please set SSH_KEY_PATH in config.bat to your .ppk file.
    pause
    exit /b 1
)

if not exist "plink.exe" (
    echo Error: plink.exe not found in %~dp0
    pause
    exit /b 1
)

if not exist "pscp.exe" (
    echo Error: pscp.exe not found in %~dp0
    pause
    exit /b 1
)

echo [1/3] Testing SSH connection...
plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -hostkey "*" %REMOTE_USER%@%REMOTE_HOST% "echo ok" >nul 2>&1
if errorlevel 1 (
    echo First connection, accepting host key...
    echo y | plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% %REMOTE_USER%@%REMOTE_HOST% "exit" >nul 2>&1
)

plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch %REMOTE_USER%@%REMOTE_HOST% "echo Connection successful"
if errorlevel 1 (
    echo Error: Cannot connect to %REMOTE_HOST% with key %SSH_KEY_PATH%
    pause
    exit /b 1
)

echo [2/3] Compressing avatars...
if exist "%ZIP_FILE%" del /f /q "%ZIP_FILE%"

set ZIP_TOOL=
where 7z >nul 2>&1
if %ERRORLEVEL% EQU 0 set ZIP_TOOL=7z
if not defined ZIP_TOOL if exist "C:\Program Files\7-Zip\7z.exe" set "ZIP_TOOL=C:\Program Files\7-Zip\7z.exe"

set "ZIP_PATH=%~dp0%ZIP_FILE%"
if defined ZIP_TOOL (
    "%ZIP_TOOL%" a -tzip "!ZIP_PATH!" "!LOCAL_AVATARS_DIR!\*" >nul
    if errorlevel 1 (
        echo Error: Failed to create zip with 7-Zip
        echo Zip path: !ZIP_PATH!
        pause
        exit /b 1
    )
) else (
    where tar >nul 2>&1
    if %ERRORLEVEL% EQU 0 (
        tar -a -cf "%ZIP_FILE%" -C "%LOCAL_AVATARS_DIR%" .
    ) else (
        echo Warning: 7-Zip/tar not found, using PowerShell zip
        powershell -NoProfile -Command "Add-Type -AssemblyName System.IO.Compression.FileSystem; [System.IO.Compression.ZipFile]::CreateFromDirectory('%LOCAL_AVATARS_DIR%', '%CD%\%ZIP_FILE%')"
    )
)

if not exist "%ZIP_FILE%" (
    echo Error: Failed to create zip file
    pause
    exit /b 1
)

echo [3/3] Uploading to server...
plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch %REMOTE_USER%@%REMOTE_HOST% "sudo mkdir -p '%REMOTE_DIR%'"
if errorlevel 1 (
    echo Error: Failed to prepare remote directory
    del "%ZIP_FILE%"
    pause
    exit /b 1
)

pscp.exe -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% "%ZIP_FILE%" %REMOTE_USER%@%REMOTE_HOST%:/tmp/%ZIP_FILE%
if errorlevel 1 (
    echo Error: Zip upload failed
    del "%ZIP_FILE%"
    pause
    exit /b 1
)

REM images 目录可能属主为 root(服务 sudo 启动),需 sudo 解压
plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch %REMOTE_USER%@%REMOTE_HOST% "sudo unzip -o /tmp/%ZIP_FILE% -d '%REMOTE_DIR%'; ec=$?; rm -f /tmp/%ZIP_FILE%; if [ $ec -eq 0 ] || [ $ec -eq 1 ]; then exit 0; else exit $ec; fi"
if errorlevel 1 (
    echo Error: Remote extraction failed. Details:
    plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch %REMOTE_USER%@%REMOTE_HOST% "ls -la '%REMOTE_DIR%' 2>&1; ls -la /tmp/%ZIP_FILE% 2>&1"
    plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch %REMOTE_USER%@%REMOTE_HOST% "rm -f /tmp/%ZIP_FILE%"
    del "%ZIP_FILE%"
    pause
    exit /b 1
)

del "%ZIP_FILE%"
echo.
echo Avatar upload completed! [%DEPLOY_ENV%] -^> %REMOTE_DIR%
echo Access example: /images/demo_avatar_1.png
pause
endlocal
