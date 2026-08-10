@echo off
setlocal EnableDelayedExpansion
cd /d "%~dp0"

call "%~dp0config-import.bat"

set SQL_FILE=%~1
if "%SQL_FILE%"=="" set SQL_FILE=%IMPORT_SQL_FILE%
if "%SQL_FILE%"=="" (
    echo Usage: %~nx0 [file.sql or file.sql.gz]
    echo   或在 config-import.bat 中设置 IMPORT_SQL_FILE
    echo Example: %~nx0 D:\var\live_db_20260810_153000.sql.gz
    pause
    exit /b 1
)
if not exist "%SQL_FILE%" (
    echo Error: file not found: %SQL_FILE%
    pause
    exit /b 1
)

if not exist "plink.exe" (
    echo Error: plink.exe not found. Copy from pub-tool\go-正式服\
    pause
    exit /b 1
)
if not exist "pscp.exe" (
    echo Error: pscp.exe not found. Copy from pub-tool\go-正式服\
    pause
    exit /b 1
)
if not exist "%SSH_KEY_PATH%" (
    echo Error: SSH key not found: %SSH_KEY_PATH%
    pause
    exit /b 1
)

for /f "usebackq delims=" %%t in (`powershell -NoProfile -Command "Get-Date -Format 'yyyyMMdd_HHmmss'"`) do set TS=%%t
set REMOTE_UPLOAD=%REMOTE_TMP_DIR%/live_db_import_%TS%.upload
set REMOTE_SQL=%REMOTE_TMP_DIR%/live_db_import_%TS%.sql

echo ========================================
echo Import live_db (%ENV_NAME%)
echo Server: %REMOTE_USER%@%REMOTE_HOST%
echo Database: %DB_NAME% @ %DB_HOST%:%DB_PORT%
echo Source: %SQL_FILE%
echo ========================================
echo.
echo WARNING: This OVERWRITES data in live_db on %ENV_NAME%!
set /p CONFIRM=Type YES to continue: 
if /i not "%CONFIRM%"=="YES" (
    echo Cancelled.
    pause
    exit /b 1
)

echo [1/4] Uploading...
pscp.exe -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% "%SQL_FILE%" %REMOTE_USER%@%REMOTE_HOST%:%REMOTE_UPLOAD%
if errorlevel 1 (
    echo Error: upload failed
    pause
    exit /b 1
)

echo [2/4] Preparing SQL on remote...
set SQL_EXT=
for %%I in ("%SQL_FILE%") do set SQL_EXT=%%~xI
if /i "!SQL_EXT!"==".gz" (
    plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T %REMOTE_USER%@%REMOTE_HOST% "gunzip -c %REMOTE_UPLOAD% > %REMOTE_SQL%"
) else (
    plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T %REMOTE_USER%@%REMOTE_HOST% "cp %REMOTE_UPLOAD% %REMOTE_SQL%"
)
if errorlevel 1 (
    echo Error: remote prepare failed
    plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T %REMOTE_USER%@%REMOTE_HOST% "rm -f %REMOTE_UPLOAD% %REMOTE_SQL%"
    pause
    exit /b 1
)

echo [3/4] Importing...
plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T %REMOTE_USER%@%REMOTE_HOST% "mariadb -h%DB_HOST% -P%DB_PORT% -u%DB_USER% -p%DB_PASSWORD% --default-character-set=utf8mb4 %DB_NAME% < %REMOTE_SQL%"
if errorlevel 1 (
    echo Error: import failed
    plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T %REMOTE_USER%@%REMOTE_HOST% "rm -f %REMOTE_UPLOAD% %REMOTE_SQL%"
    pause
    exit /b 1
)

echo [4/4] Cleaning remote temp files...
plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T %REMOTE_USER%@%REMOTE_HOST% "rm -f %REMOTE_UPLOAD% %REMOTE_SQL%"

echo.
echo Import OK
echo.
endlocal
pause
