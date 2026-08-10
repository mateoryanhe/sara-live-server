@echo off
REM 官网部署配置（正式服 SFTP 上传）

REM 远程服务器
set REMOTE_HOST=52.9.70.64
set REMOTE_PORT=22

REM SFTP 专用账号（仅 official-site 目录，不可执行命令）
set SFTP_USER=official-site
set SFTP_KEY_PATH=%~dp0keys\official-site-sftp.pem

REM 本地路径
set LOCAL_PROJECT_ROOT=D:\company-code\sara-live-server
set LOCAL_SITE_DIR=%LOCAL_PROJECT_ROOT%\official-site
set STAGING_DIR=%~dp0_staging

REM 远程 SFTP 登录后即为官网根目录（chroot /files -> official-site）
set REMOTE_DIR=/

REM 注销账号等接口的 API 根地址（写入 js/site.js 的 apiBaseUrl）
REM 留空表示与官网同域（相对路径请求）
set API_BASE_URL=https://www.saralive.net
