@echo off
cd /d "%~dp0"
echo Starting local MariaDB on port 3390...
docker compose up -d
if %errorlevel% neq 0 (
    echo Failed to start. Is Docker running?
    pause
    exit /b 1
)
echo.
echo Waiting for MariaDB to be ready...
timeout /t 8 /nobreak >nul
docker compose ps
echo.
echo Local MariaDB:
echo   Host: 127.0.0.1
echo   Port: 3390
echo   Database: live_db
echo   appuser (from host): Appledev882116
echo.
echo Connect:
echo   docker exec -it sara-live-mariadb mariadb -uroot -pAppledev882116 live_db
echo   mariadb -uappuser -pAppledev882116 -h127.0.0.1 -P3390 live_db
echo.
pause
