@echo off
REM 导出源：直播正式服 live_db

set ENV_NAME=直播正式服

REM SSH
set REMOTE_HOST=52.9.70.64
set REMOTE_USER=ec2-user
set REMOTE_PORT=22
set SSH_KEY_PATH=D:\tools\ppk\live-test.ppk

REM MariaDB（服务器本机，与 config/prod/config.yaml 一致）
set DB_HOST=127.0.0.1
set DB_PORT=14501
set DB_NAME=live_db
set DB_USER=root
set DB_PASSWORD=c63eac559a03dece518e3eb7a601b30e

REM 本地备份输出目录
set LOCAL_BACKUP_DIR=D:\var

REM 远程临时目录
set REMOTE_TMP_DIR=/tmp
