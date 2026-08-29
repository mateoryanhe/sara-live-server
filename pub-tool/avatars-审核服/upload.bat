@echo off
setlocal EnableDelayedExpansion
cd /d "%~dp0"

call "%~dp0config.bat"
set ZIP_FILE=avatars_upload_review.zip

echo ========================================
echo Avatar Upload [review]
echo Local directory: %LOCAL_AVATARS_DIR%
echo Target server: %REMOTE_USER%@%REMOTE_HOST%:%REMOTE_PORT%
echo Remote directory: %REMOTE_DIR%
echo Staging: %REMOTE_STAGE% (not /tmp)
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
plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch %REMOTE_USER%@%REMOTE_HOST% "echo Connection successful"
if errorlevel 1 (
    echo Error: Cannot connect to %REMOTE_HOST%
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
) else (
    where tar >nul 2>&1
    if %ERRORLEVEL% EQU 0 (
        tar -a -cf "%ZIP_FILE%" -C "%LOCAL_AVATARS_DIR%" .
    ) else (
        powershell -NoProfile -Command "Add-Type -AssemblyName System.IO.Compression.FileSystem; [System.IO.Compression.ZipFile]::CreateFromDirectory('%LOCAL_AVATARS_DIR%', '%CD%\%ZIP_FILE%')"
    )
)

if not exist "%ZIP_FILE%" (
    echo Error: Failed to create zip file
    pause
    exit /b 1
)

echo [3/3] Uploading to server...
plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch %REMOTE_USER%@%REMOTE_HOST% "mkdir -p '%REMOTE_STAGE%' && sudo mkdir -p '%REMOTE_DIR%'"
if errorlevel 1 (
    echo Error: Failed to prepare remote directories
    del "%ZIP_FILE%"
    pause
    exit /b 1
)

pscp.exe -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% "%ZIP_FILE%" %REMOTE_USER%@%REMOTE_HOST%:%REMOTE_STAGE%/%ZIP_FILE%
if errorlevel 1 (
    echo Error: Zip upload failed
    del "%ZIP_FILE%"
    pause
    exit /b 1
)

plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch %REMOTE_USER%@%REMOTE_HOST% "sudo unzip -o %REMOTE_STAGE%/%ZIP_FILE% -d '%REMOTE_DIR%'; ec=$?; rm -f %REMOTE_STAGE%/%ZIP_FILE%; if [ $ec -eq 0 ] || [ $ec -eq 1 ]; then exit 0; else exit $ec; fi"
if errorlevel 1 (
    echo Error: Remote extraction failed
    plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch %REMOTE_USER%@%REMOTE_HOST% "ls -la '%REMOTE_DIR%' 2>&1; ls -la %REMOTE_STAGE%/%ZIP_FILE% 2>&1"
    plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch %REMOTE_USER%@%REMOTE_HOST% "rm -f %REMOTE_STAGE%/%ZIP_FILE%"
    del "%ZIP_FILE%"
    pause
    exit /b 1
)

del "%ZIP_FILE%"
echo.
echo Avatar upload completed! -^> %REMOTE_DIR%
echo Access example: https://reviewpng.saralive.net/demo_avatar_1.png
pause
endlocal
