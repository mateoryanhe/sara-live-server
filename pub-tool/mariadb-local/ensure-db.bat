@echo off
setlocal EnableDelayedExpansion
cd /d "%~dp0"
call "%~dp0config.bat"

echo ========================================
echo   sara MariaDB local (port %MARIADB_HOST_PORT%)
echo ========================================

where docker >nul 2>&1 || (echo ERROR: docker not found. & exit /b 1)

set CONTAINER=
for /f "delims=" %%n in ('docker ps --format "{{.Names}}"') do (
  for /f "tokens=*" %%p in ('docker port "%%n" 3306/tcp 2^>nul') do (
    echo %%p | findstr ":%MARIADB_HOST_PORT%" >nul && set CONTAINER=%%n
  )
)

if not defined CONTAINER (
  echo No container on %MARIADB_HOST_PORT%. Run create-docker.bat first.
  exit /b 1
)

echo Using container: !CONTAINER!

docker exec !CONTAINER! mariadb -u%MARIADB_USER% -p%MARIADB_PASSWORD% -e "CREATE DATABASE IF NOT EXISTS %MARIADB_DATABASE% CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;" 2>nul
if errorlevel 1 (
  docker exec !CONTAINER! mariadb -uroot -p%MARIADB_ROOT_PASSWORD% -e "CREATE DATABASE IF NOT EXISTS %MARIADB_DATABASE% CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci; CREATE USER IF NOT EXISTS '%MARIADB_USER%'@'%%' IDENTIFIED BY '%MARIADB_PASSWORD%'; GRANT ALL PRIVILEGES ON %MARIADB_DATABASE%.* TO '%MARIADB_USER%'@'%%'; FLUSH PRIVILEGES;"
)

docker exec !CONTAINER! mariadb -u%MARIADB_USER% -p%MARIADB_PASSWORD% -e "USE %MARIADB_DATABASE%; SELECT COUNT(*) AS table_count FROM information_schema.tables WHERE table_schema='%MARIADB_DATABASE%';"
echo Done. DSN: 127.0.0.1:%MARIADB_HOST_PORT%/%MARIADB_DATABASE%
exit /b 0