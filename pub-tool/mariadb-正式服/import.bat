@echo off
setlocal EnableDelayedExpansion
cd /d "%~dp0"

call "%~dp0config.bat"

set SQL_FILE=%~1
if "%SQL_FILE%"=="" (
    echo Usage: %~nx0 ^<sql_or_sql.gz_file^>
    echo Example: %~nx0 backup\live_db_20260810_153000.sql.gz
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
set REMOTE_FILE=%REMOTE_TMP_DIR%/live_db_import_%TS%.sql
set REMOTE_GZ=%REMOTE_FILE%.gz

echo ========================================
echo Import live_db to prod
echo Server: %REMOTE_USER%@%REMOTE_HOST%
echo Database: %DB_NAME% (port %DB_PORT%)
echo Source: %SQL_FILE%
echo ========================================
echo.
echo WARNING: This will overwrite data in live_db on production!
set /p CONFIRM=Type YES to continue: 
if /i not "%CONFIRM%"=="YES" (
    echo Cancelled.
    pause
    exit /b 1
)

echo [1/4] Uploading SQL file...
pscp.exe -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% "%SQL_FILE%" %REMOTE_USER%@%REMOTE_HOST%:%REMOTE_GZ%
if errorlevel 1 (
    echo Error: upload failed
    pause
    exit /b 1
)

echo [2/4] Decompressing on remote (if needed)...
plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T %REMOTE_USER%@%REMOTE_HOST% "gunzip -c %REMOTE_GZ% > %REMOTE_FILE% 2>/dev/null || cp %REMOTE_GZ% %REMOTE_FILE%"

echo [3/4] Importing into live_db...
plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T %REMOTE_USER%@%REMOTE_HOST% "mariadb -h%DB_HOST% -P%DB_PORT% -u%DB_USER% -p%DB_PASSWORD% --default-character-set=utf8mb4 %DB_NAME% < %REMOTE_FILE%"
if errorlevel 1 (
    echo Error: import failed
    plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T %REMOTE_USER%@%REMOTE_HOST% "rm -f %REMOTE_FILE% %REMOTE_GZ%"
    pause
    exit /b 1
)

echo [4/4] Cleaning remote temp files...
plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T %REMOTE_USER%@%REMOTE_HOST% "rm -f %REMOTE_FILE% %REMOTE_GZ%"

echo.
echo Import OK
echo.
endlocal
pause
