@echo off
setlocal EnableDelayedExpansion
cd /d "%~dp0"

call "%~dp0config-export.bat"

if not exist "plink.exe" (
    echo Error: plink.exe not found. Copy from pub-tool\go-prod\
    pause
    exit /b 1
)
if not exist "pscp.exe" (
    echo Error: pscp.exe not found. Copy from pub-tool\go-prod\
    pause
    exit /b 1
)
if not exist "%SSH_KEY_PATH%" (
    echo Error: SSH key not found: %SSH_KEY_PATH%
    pause
    exit /b 1
)

if not exist "%LOCAL_BACKUP_DIR%" mkdir "%LOCAL_BACKUP_DIR%"

for /f "usebackq delims=" %%t in (`powershell -NoProfile -Command "Get-Date -Format 'yyyyMMdd_HHmmss'"`) do set TS=%%t
set REMOTE_FILE=%REMOTE_TMP_DIR%/live_db_%TS%.sql.gz
set LOCAL_FILE=%LOCAL_BACKUP_DIR%\live_db_%TS%.sql.gz

echo ========================================
echo Export live_db (%ENV_NAME%)
echo Server: %REMOTE_USER%@%REMOTE_HOST%
echo Database: %DB_NAME% @ %DB_HOST%:%DB_PORT%
echo Output: %LOCAL_FILE%
echo ========================================
echo.

echo [1/3] mariadb-dump on remote server...
plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T %REMOTE_USER%@%REMOTE_HOST% "mariadb-dump -h%DB_HOST% -P%DB_PORT% -u%DB_USER% -p%DB_PASSWORD% --single-transaction --routines --triggers --events --default-character-set=utf8mb4 %DB_NAME% | gzip -c > %REMOTE_FILE%"
if errorlevel 1 (
    echo Error: remote mariadb-dump failed
    pause
    exit /b 1
)

echo [2/3] Downloading backup...
pscp.exe -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% %REMOTE_USER%@%REMOTE_HOST%:%REMOTE_FILE% "%LOCAL_FILE%"
if errorlevel 1 (
    echo Error: download failed
    plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T %REMOTE_USER%@%REMOTE_HOST% "rm -f %REMOTE_FILE%"
    pause
    exit /b 1
)

echo [3/3] Cleaning remote temp file...
plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T %REMOTE_USER%@%REMOTE_HOST% "rm -f %REMOTE_FILE%"

for %%A in ("%LOCAL_FILE%") do set SIZE=%%~zA
echo.
echo Export OK
echo   File: %LOCAL_FILE%
echo   Size: !SIZE! bytes
echo.
echo Import: import.bat "%LOCAL_FILE%"
echo.
endlocal
pause
