@echo off
REM 导入目标：直播测试服 live_db（可按实际环境修改）

set ENV_NAME=直播测试服

REM SSH
set REMOTE_HOST=54.241.124.37
set REMOTE_USER=ec2-user
set REMOTE_PORT=22
set SSH_KEY_PATH=D:\tools\ppk\live-test.ppk

REM MariaDB（服务器本机，与 config/dev/config.yaml 一致）
set DB_HOST=127.0.0.1
set DB_PORT=13307
set DB_NAME=kyc
set DB_USER=root
set DB_PASSWORD=Appledev882116

REM 本地备份文件（数据源）；留空则 import.bat 必须传参
REM 示例: set IMPORT_SQL_FILE=D:\var\live_db_20260810_153000.sql.gz
set IMPORT_SQL_FILE=D:\var\live_db_20260810_162251.sql.gz

REM 远程临时目录
set REMOTE_TMP_DIR=/tmp
