@echo off
setlocal EnableDelayedExpansion
cd /d "%~dp0"

REM 用法: upload.bat [test|prod]  默认 test
set DEPLOY_ENV=%~1
if "%DEPLOY_ENV%"=="" set DEPLOY_ENV=test

if /i not "%DEPLOY_ENV%"=="test" if /i not "%DEPLOY_ENV%"=="prod" (
    echo Usage: %~nx0 [test^|prod]
    echo   test - 使用 cms/.env.test 构建并部署到测试目录
    echo   prod - 使用 cms/.env.production 构建并部署到生产目录
    exit /b 1
)

call "%~dp0config.bat"

if /i "%DEPLOY_ENV%"=="prod" (
    set BUILD_CMD=build:prod
    set REMOTE_DIR=%REMOTE_DIR_PROD%
    set ZIP_FILE=upload_package_prod.zip
) else (
    set BUILD_CMD=build:test
    set REMOTE_DIR=%REMOTE_DIR_TEST%
    set ZIP_FILE=upload_package_test.zip
)

echo ========================================
echo CMS Build and Upload
echo Environment: %DEPLOY_ENV%
echo Build command: npm run %BUILD_CMD%
echo Target server: %REMOTE_USER%@%REMOTE_HOST%:%REMOTE_PORT%
echo Remote directory: %REMOTE_DIR%
echo SSH key: %SSH_KEY_PATH%
echo ========================================
echo.

if not exist "%VUE_PROJECT_DIR%" (
    echo Error: Vue project directory does not exist: %VUE_PROJECT_DIR%
    exit /b 1
)

if not exist "%SSH_KEY_PATH%" (
    echo Error: SSH key file does not exist: %SSH_KEY_PATH%
    echo Please set SSH_KEY_PATH in config.bat to your .ppk file.
    exit /b 1
)

if not exist "plink.exe" (
    echo Error: plink.exe not found in %~dp0
    exit /b 1
)

if not exist "pscp.exe" (
    echo Error: pscp.exe not found in %~dp0
    exit /b 1
)

REM ---------- Build Vue ----------
echo [1/5] Building Vue project...
cd /d "%VUE_PROJECT_DIR%"

if not exist "package.json" (
    echo Error: package.json not found in %VUE_PROJECT_DIR%
    exit /b 1
)

echo Running: npm run %BUILD_CMD%
call npm run %BUILD_CMD%
if errorlevel 1 (
    echo Error: Vue build failed
    cd /d "%~dp0"
    exit /b 1
)

cd /d "%~dp0"

if not exist "%BUILD_OUTPUT_DIR%" (
    echo Error: Build output directory does not exist: %BUILD_OUTPUT_DIR%
    exit /b 1
)

if not exist "%BUILD_OUTPUT_DIR%\index.html" (
    echo Error: index.html not found in build output: %BUILD_OUTPUT_DIR%
    exit /b 1
)

REM ---------- SSH connect (ppk) ----------
echo [2/5] Testing SSH connection...
plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -hostkey "*" %REMOTE_USER%@%REMOTE_HOST% "echo ok" >nul 2>&1
if errorlevel 1 (
    echo First connection, accepting host key...
    echo y | plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% %REMOTE_USER%@%REMOTE_HOST% "exit" >nul 2>&1
)

plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch %REMOTE_USER%@%REMOTE_HOST% "echo Connection successful"
if errorlevel 1 (
    echo Error: Cannot connect to %REMOTE_HOST% with key %SSH_KEY_PATH%
    exit /b 1
)

REM ---------- Zip (prefer tar for Linux-friendly paths) ----------
echo [3/5] Compressing build output...
if exist "%ZIP_FILE%" del /f /q "%ZIP_FILE%"

set ZIP_OK=0
where tar >nul 2>&1
if %ERRORLEVEL% EQU 0 (
    tar -a -cf "%ZIP_FILE%" -C "%BUILD_OUTPUT_DIR%" .
    if not errorlevel 1 set ZIP_OK=1
)

if !ZIP_OK! EQU 0 (
    set ZIP_TOOL=
    where 7z >nul 2>&1
    if %ERRORLEVEL% EQU 0 set ZIP_TOOL=7z
    if not defined ZIP_TOOL if exist "C:\Program Files\7-Zip\7z.exe" set "ZIP_TOOL=C:\Program Files\7-Zip\7z.exe"
    if defined ZIP_TOOL (
        set "ZIP_PATH=%~dp0%ZIP_FILE%"
        pushd "%BUILD_OUTPUT_DIR%"
        "%ZIP_TOOL%" a -tzip "!ZIP_PATH!" * >nul
        popd
        if exist "%~dp0%ZIP_FILE%" set ZIP_OK=1
    )
)

if !ZIP_OK! EQU 0 (
    echo Error: Failed to create zip file. Please install tar or 7-Zip.
    exit /b 1
)

if not exist "%ZIP_FILE%" (
    echo Error: Failed to create zip file
    exit /b 1
)

if "%REMOTE_DIR%"=="" (
    echo Error: REMOTE_DIR is not configured in config.bat
    del "%ZIP_FILE%"
    exit /b 1
)
if /i "%REMOTE_DIR%"=="/" (
    echo Error: REMOTE_DIR cannot be /
    del "%ZIP_FILE%"
    exit /b 1
)

REM ---------- Upload and extract ----------
echo [4/5] Uploading to server...
plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch %REMOTE_USER%@%REMOTE_HOST% "mkdir -p '%REMOTE_DIR%'"
if errorlevel 1 (
    echo Error: Failed to prepare remote directory
    del "%ZIP_FILE%"
    exit /b 1
)

pscp.exe -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% "%ZIP_FILE%" %REMOTE_USER%@%REMOTE_HOST%:/tmp/%ZIP_FILE%
if errorlevel 1 (
    echo Error: Zip upload failed
    del "%ZIP_FILE%"
    exit /b 1
)

REM 解压到临时目录后仅覆盖包内顶层文件/目录，保留 log-query-export 等服务器侧目录
REM unzip 在仅有 warning 时返回 1，需与真正失败区分
plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch %REMOTE_USER%@%REMOTE_HOST% "STAGE=/tmp/cms-stage-$$; REMOTE='%REMOTE_DIR%'; ZIP=/tmp/%ZIP_FILE%; rm -rf $STAGE; mkdir -p $STAGE $REMOTE; unzip -o $ZIP -d $STAGE; ec=$?; if [ $ec -ne 0 ] && [ $ec -ne 1 ]; then rm -rf $STAGE; rm -f $ZIP; exit $ec; fi; for f in $STAGE/*; do [ -e $f ] || continue; rm -rf $REMOTE/$(basename $f); done; cp -a $STAGE/. $REMOTE/; rm -rf $STAGE; rm -f $ZIP; test -f $REMOTE/index.html"
if errorlevel 1 (
    echo Error: Remote extraction failed or index.html missing. Details:
    plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch %REMOTE_USER%@%REMOTE_HOST% "ls -la '%REMOTE_DIR%' 2>&1; ls -la /tmp/%ZIP_FILE% 2>&1; unzip -t /tmp/%ZIP_FILE% 2>&1 || true"
    del "%ZIP_FILE%"
    exit /b 1
)

echo [5/5] Verifying remote deployment...
plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch %REMOTE_USER%@%REMOTE_HOST% "ls -la '%REMOTE_DIR%' | head -20"
if errorlevel 1 (
    echo Warning: Could not list remote directory after deploy
)

del "%ZIP_FILE%"
echo.
echo Build and upload completed! [%DEPLOY_ENV%] -^> %REMOTE_DIR%
echo Access: https://cms.bigtktool.shop/
endlocal
