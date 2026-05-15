@echo off
REM Remote server configuration
set REMOTE_HOST=44.244.204.122
set REMOTE_USER=ec2-user
set REMOTE_DIR=/home/ec2-user/xgameserver

REM Local path configuration
set LOCAL_PROJECT_PATH=D:\go-project\xrgameserver
set LOCAL_GO_SRC=%LOCAL_PROJECT_PATH%\go-src
set LOCAL_CONFIG_PATH=%LOCAL_PROJECT_PATH%\config\dev\config.yaml
set LOCAL_BUILD_PATH=%LOCAL_PROJECT_PATH%\go-build

REM SSH key configuration
set SSH_KEY_PATH=C:\Users\EDY\Downloads\tt.ppk

REM Application name
set APP_NAME=xr-game-server

REM Remote sudo command (not used)
set SUDO_CMD=

REM Graceful shutdown wait time in seconds
set SHUTDOWN_WAIT_TIME=5