@echo off
setlocal EnableDelayedExpansion
cd /d "%~dp0"

REM 用法: upload.bat [test|prod]  默认 test
set DEPLOY_ENV=%~1
if "%DEPLOY_ENV%"=="" set DEPLOY_ENV=test

if /i not "%DEPLOY_ENV%"=="test" if /i not "%DEPLOY_ENV%"=="prod" (
    echo Usage: %~nx0 [test^|prod]
    echo   test - 部署到测试目录，使用 config.bat 中 API_BASE_URL_TEST
    echo   prod - 部署到生产目录，使用 config.bat 中 API_BASE_URL_PROD
    pause
    exit /b 1
)

call "%~dp0config.bat"

if /i "%DEPLOY_ENV%"=="prod" (
    set REMOTE_DIR=%REMOTE_DIR_PROD%
    set API_BASE_URL=%API_BASE_URL_PROD%
    set ZIP_FILE=official_site_upload_prod.zip
) else (
    set REMOTE_DIR=%REMOTE_DIR_TEST%
    set API_BASE_URL=%API_BASE_URL_TEST%
    set ZIP_FILE=official_site_upload_test.zip
)

echo ========================================
echo Official Site Upload
echo Environment: %DEPLOY_ENV%
echo Local directory: %LOCAL_SITE_DIR%
echo Staging directory: %STAGING_DIR%
echo API base URL: [%API_BASE_URL%]
echo Target server: %REMOTE_USER%@%REMOTE_HOST%:%REMOTE_PORT%
echo Remote directory: %REMOTE_DIR%
echo SSH key: %SSH_KEY_PATH%
echo ========================================
echo.

if not exist "%LOCAL_SITE_DIR%" (
    echo Error: Official site directory does not exist: %LOCAL_SITE_DIR%
    pause
    exit /b 1
)

if not exist "%LOCAL_SITE_DIR%\index.html" (
    echo Error: index.html not found in %LOCAL_SITE_DIR%
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

REM ---------- Prepare staging ----------
echo [1/4] Preparing staging files...
if exist "%STAGING_DIR%" rmdir /s /q "%STAGING_DIR%"
mkdir "%STAGING_DIR%" >nul 2>&1
if errorlevel 1 (
    echo Error: Failed to create staging directory: %STAGING_DIR%
    pause
    exit /b 1
)

xcopy /E /I /Y /Q "%LOCAL_SITE_DIR%\*" "%STAGING_DIR%\" >nul
if errorlevel 1 (
    echo Error: Failed to copy official-site to staging
    pause
    exit /b 1
)

if not exist "%STAGING_DIR%\js\site.js" (
    echo Error: site.js not found in staging directory
    pause
    exit /b 1
)

set "SITE_JS=%STAGING_DIR%\js\site.js"
powershell -NoProfile -Command "$path = '%SITE_JS%'; $content = Get-Content -LiteralPath $path -Raw -Encoding UTF8; $content = $content -replace 'apiBaseUrl:\s*''[^'']*''', 'apiBaseUrl: ''%API_BASE_URL%'''; [System.IO.File]::WriteAllText($path, $content, [System.Text.UTF8Encoding]::new($false))"
if errorlevel 1 (
    echo Error: Failed to patch apiBaseUrl in site.js
    pause
    exit /b 1
)

REM ---------- SSH connect (ppk) ----------
echo [2/4] Testing SSH connection...
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

REM ---------- Zip ----------
echo [3/4] Compressing staging files...
if exist "%ZIP_FILE%" del /f /q "%ZIP_FILE%"

set ZIP_TOOL=
where 7z >nul 2>&1
if %ERRORLEVEL% EQU 0 set ZIP_TOOL=7z
if not defined ZIP_TOOL if exist "C:\Program Files\7-Zip\7z.exe" set "ZIP_TOOL=C:\Program Files\7-Zip\7z.exe"

set "ZIP_PATH=%~dp0%ZIP_FILE%"
if defined ZIP_TOOL (
    "%ZIP_TOOL%" a -tzip "!ZIP_PATH!" "!STAGING_DIR!\*" >nul
    if errorlevel 1 (
        echo Error: Failed to create zip with 7-Zip
        echo Zip path: !ZIP_PATH!
        pause
        exit /b 1
    )
) else (
    where tar >nul 2>&1
    if %ERRORLEVEL% EQU 0 (
        tar -a -cf "%ZIP_FILE%" -C "%STAGING_DIR%" .
    ) else (
        echo Warning: 7-Zip/tar not found, using PowerShell zip ^(Linux unzip may warn^)
        powershell -NoProfile -Command "Add-Type -AssemblyName System.IO.Compression.FileSystem; [System.IO.Compression.ZipFile]::CreateFromDirectory('%STAGING_DIR%', '%CD%\%ZIP_FILE%')"
    )
)

if not exist "%ZIP_FILE%" (
    echo Error: Failed to create zip file
    pause
    exit /b 1
)

if "%REMOTE_DIR%"=="" (
    echo Error: REMOTE_DIR is not configured in config.bat
    pause
    exit /b 1
)
if /i "%REMOTE_DIR%"=="/" (
    echo Error: REMOTE_DIR cannot be /
    pause
    exit /b 1
)

REM ---------- Upload and extract ----------
echo [4/4] Uploading to server...
plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch %REMOTE_USER%@%REMOTE_HOST% "mkdir -p '%REMOTE_DIR%'"
if errorlevel 1 (
    echo Error: Failed to prepare remote directory: %REMOTE_DIR%
    plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch %REMOTE_USER%@%REMOTE_HOST% "ls -ld '%REMOTE_DIR%' 2>&1"
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

REM 解压到临时目录后仅覆盖包内文件，避免 rm -rf 整目录时因无关子目录权限失败
REM unzip 在仅有 warning 时返回 1，需与真正失败区分
plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch %REMOTE_USER%@%REMOTE_HOST% "STAGE=/tmp/official-site-stage-$$; REMOTE='%REMOTE_DIR%'; ZIP=/tmp/%ZIP_FILE%; rm -rf $STAGE; mkdir -p $STAGE $REMOTE; unzip -o $ZIP -d $STAGE; ec=$?; if [ $ec -ne 0 ] && [ $ec -ne 1 ]; then rm -rf $STAGE; rm -f $ZIP; exit $ec; fi; for f in $STAGE/*; do [ -e $f ] || continue; rm -rf $REMOTE/$(basename $f); done; cp -a $STAGE/. $REMOTE/; rm -rf $STAGE; rm -f $ZIP; exit 0"
if errorlevel 1 (
    echo Error: Remote extraction failed. Details:
    plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch %REMOTE_USER%@%REMOTE_HOST% "ls -la /tmp/%ZIP_FILE% 2>&1; unzip -t /tmp/%ZIP_FILE% 2>&1 || true"
    del "%ZIP_FILE%"
    pause
    exit /b 1
)

if exist "%STAGING_DIR%" rmdir /s /q "%STAGING_DIR%"
del "%ZIP_FILE%"
echo.
echo Official site upload completed! [%DEPLOY_ENV%] -^> %REMOTE_DIR%
echo Access example: /official-site/index.html
pause
endlocal
