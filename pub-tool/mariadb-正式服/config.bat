@echo off
REM 直播正式服 MariaDB 导入/导出配置

REM SSH（与 go-正式服 一致）
set REMOTE_HOST=52.9.70.64
set REMOTE_USER=ec2-user
set REMOTE_PORT=22
set SSH_KEY_PATH=D:\tools\ppk\live-test.ppk

REM 数据库（服务器本机连接；远程直连需 SSH 隧道，见 README）
set DB_HOST=127.0.0.1
set DB_PORT=14501
set DB_NAME=live_db
set DB_USER=live
set DB_PASSWORD=zp8JXc25eKK5wL4nnbhB

REM 本地备份目录
set LOCAL_BACKUP_DIR=%~dp0backup

REM 远程临时目录
set REMOTE_TMP_DIR=/tmp
