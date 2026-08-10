@echo off
setlocal EnableDelayedExpansion
cd /d "%~dp0"

REM 用法: upload.bat [test|prod]  默认 test
set DEPLOY_ENV=%~1
if "%DEPLOY_ENV%"=="" set DEPLOY_ENV=test

if /i not "%DEPLOY_ENV%"=="test" if /i not "%DEPLOY_ENV%"=="prod" (
    echo Usage: %~nx0 [test^|prod]
    echo   test - 使用 API_BASE_URL_TEST
    echo   prod - 使用 API_BASE_URL_PROD
    pause
    exit /b 1
)

call "%~dp0config.bat"

REM OpenSSH sftp 无法 lcd 含中文路径，staging 放到 TEMP
set "STAGING_DIR=%TEMP%\official-site-upload-staging"

if /i "%DEPLOY_ENV%"=="prod" (
    set API_BASE_URL=%API_BASE_URL_PROD%
) else (
    set API_BASE_URL=%API_BASE_URL_TEST%
)

echo ========================================
echo Official Site Upload (SFTP)
echo Environment: %DEPLOY_ENV%
echo Local directory: %LOCAL_SITE_DIR%
echo Staging directory: %STAGING_DIR%
echo API base URL: [%API_BASE_URL%]
echo Target: %SFTP_USER%@%REMOTE_HOST%:%REMOTE_PORT%
echo SFTP key: %SFTP_KEY_PATH%
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

if not exist "%SFTP_KEY_PATH%" (
    echo Error: SFTP key not found: %SFTP_KEY_PATH%
    pause
    exit /b 1
)

where sftp >nul 2>&1
if errorlevel 1 (
    echo Error: OpenSSH sftp not found. Enable Windows OpenSSH Client.
    pause
    exit /b 1
)

REM ---------- Prepare staging ----------
echo [1/3] Preparing staging files...
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

REM ---------- SFTP upload ----------
echo [2/3] Uploading via SFTP...

REM OpenSSH 要求私钥仅当前用户可读
icacls "%SFTP_KEY_PATH%" /inheritance:r >nul 2>&1
icacls "%SFTP_KEY_PATH%" /grant:r "%USERNAME%:(R)" >nul 2>&1

set "SFTP_BATCH=%TEMP%\official-site-upload-%RANDOM%.sftp"
set "STAGING_SFTP=%STAGING_DIR:\=/%"

REM put -r * 在远程子目录不存在时会失败；先生成 mkdir，再逐文件 put
powershell -NoProfile -Command ^
  "$staging = '%STAGING_DIR%';" ^
  "$out = '%SFTP_BATCH%';" ^
  "$stagingSftp = ($staging -replace '\\', '/');" ^
  "$lines = @('lcd ' + $stagingSftp);" ^
  "Get-ChildItem -LiteralPath $staging -Directory -Recurse | Sort-Object { $_.FullName.Length } | ForEach-Object {" ^
  "  $rel = $_.FullName.Substring($staging.Length).TrimStart('\').Replace('\', '/');" ^
  "  if ($rel) { $lines += ('-mkdir ' + $rel) }" ^
  "};" ^
  "Get-ChildItem -LiteralPath $staging -File -Recurse | ForEach-Object {" ^
  "  $rel = $_.FullName.Substring($staging.Length).TrimStart('\').Replace('\', '/');" ^
  "  $local = $_.FullName.Replace('\', '/');" ^
  "  $lines += ('put \"' + $local + '\" \"' + $rel + '\"');" ^
  "};" ^
  "$lines += 'bye';" ^
  "[IO.File]::WriteAllLines($out, $lines)"
if errorlevel 1 (
    echo Error: Failed to generate SFTP batch file
    pause
    exit /b 1
)

sftp -i "%SFTP_KEY_PATH%" -P %REMOTE_PORT% -o BatchMode=yes -o StrictHostKeyChecking=accept-new -b "%SFTP_BATCH%" %SFTP_USER%@%REMOTE_HOST%
if errorlevel 1 (
    echo Error: SFTP connection or upload failed
    del /f /q "%SFTP_BATCH%" 2>nul
    pause
    exit /b 1
)

del /f /q "%SFTP_BATCH%" 2>nul

echo [3/3] Cleaning up...
if exist "%STAGING_DIR%" rmdir /s /q "%STAGING_DIR%"

echo.
echo Official site upload completed! [%DEPLOY_ENV%]
echo Uploaded via %SFTP_USER% to /home/ec2-user/cdn/official-site
echo Access example: /official-site/index.html
pause
endlocal
