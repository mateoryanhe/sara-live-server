@echo off
REM 与 sara-live-server config/local/config.yaml 中 database 一致
set MARIADB_HOST_PORT=3391
set MARIADB_DATABASE=live_db
set MARIADB_USER=appuser
set MARIADB_PASSWORD=Appledev882116
set MARIADB_ROOT_PASSWORD=Appledev882116
set MARIADB_CONTAINER=sara-live-mariadb
set COMPOSE_PROJECT_NAME=sara-live-mariadb