@echo off
REM 官网部署配置（测试服 SFTP 上传）

REM 远程服务器
set REMOTE_HOST=54.241.124.37
set REMOTE_PORT=22

REM SFTP 专用账号（仅 official-site 目录，不可执行命令）
set SFTP_USER=official-site
set SFTP_KEY_PATH=%~dp0keys\official-site-sftp.pem

REM 本地路径
set LOCAL_PROJECT_ROOT=D:\company-code\sara-live-server
set LOCAL_SITE_DIR=%LOCAL_PROJECT_ROOT%\official-site
set STAGING_DIR=%~dp0_staging

REM 远程 SFTP 登录后即为官网根目录（chroot /files -> official-site）
set REMOTE_DIR_TEST=/
set REMOTE_DIR_PROD=/

REM 注销账号等接口的 API 根地址（写入 js/site.js 的 apiBaseUrl）
set API_BASE_URL_TEST=
set API_BASE_URL_PROD=
