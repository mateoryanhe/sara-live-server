@echo off
REM 语聊测试服 Go 发版（yuliao 分支；勿改回直播测试服 54.241.124.37）
set REMOTE_HOST=63.185.66.53
set REMOTE_USER=ec2-user
set REMOTE_DIR=/home/ec2-user/xgameserver
REM /tmp 为容量较小的 tmpfs，发版包统一放到磁盘 staging
set REMOTE_STAGE=/home/ec2-user/staging

REM Local path configuration
set LOCAL_PROJECT_PATH=D:\company-code\sara-live-server
set LOCAL_GO_SRC=%LOCAL_PROJECT_PATH%\go-src
set LOCAL_CONFIG_PATH=%LOCAL_PROJECT_PATH%\config\yuliao\config.yaml
REM 热重启一/二阶段秒数在 CMS preload_cfgs; deploy/restart 脚本等待默认各 60
set LOCAL_BUILD_PATH=%LOCAL_PROJECT_PATH%\go-build

REM SSH key configuration（D:\root\yuliao.pem -> D:\tools\ppk\yuliao.ppk）
set SSH_KEY_PATH=D:\tools\ppk\yuliao.ppk

REM Pinned SSH host key for restart.bat (plink -hostkey)
set "SSH_HOST_KEY=ssh-ed25519 255 SHA256:9IMFNzd8vX1/NbhLoUBxCzMNJtnkC9XZInJEydJIEfg"

set "RESTART_TITLE=XR Game Server Yuliao Test"

REM Application name
set APP_NAME=xr-game-server

REM Remote sudo command (not used)
set SUDO_CMD=

REM Hot restart auth (must match CMS/server runtime cfg hotRestartAuth, default nGH66S4TjBjQqCKyWJAM)
set HOT_RESTART_AUTH=nGH66S4TjBjQqCKyWJAM

REM SSH port
set REMOTE_PORT=22
