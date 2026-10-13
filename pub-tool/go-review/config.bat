@echo off
set "RESTART_TITLE=XR Game Server Review Restart"
REM Remote server configuration
set REMOTE_HOST=18.144.165.177
set REMOTE_USER=ec2-user
set REMOTE_DIR=/home/ec2-user/xgameserver
REM Review server /tmp is small tmpfs; use REMOTE_STAGE on disk
set REMOTE_STAGE=/home/ec2-user/staging

REM Local path configuration
set LOCAL_PROJECT_PATH=D:\company-code\sara-live-server
set LOCAL_GO_SRC=%LOCAL_PROJECT_PATH%\go-src
set LOCAL_CONFIG_PATH=%LOCAL_PROJECT_PATH%\config\review\config.yaml
REM Hot restart phase seconds in CMS preload_cfgs; deploy/restart wait default 60+60
set LOCAL_BUILD_PATH=%LOCAL_PROJECT_PATH%\go-build

REM SSH key configuration
set SSH_KEY_PATH=D:\tools\ppk\live-review.ppk

REM Pinned SSH host key (prevents first-run prompts and host spoofing)
set "SSH_HOST_KEY=ssh-ed25519 255 SHA256:eoE8/boxqJSOy9fLU+J+pL+AxSJk4bUe06kHCTgAXPU"

REM Application name
set APP_NAME=xr-game-server

REM Remote sudo command (not used)
set SUDO_CMD=

REM Hot restart auth (must match CMS/server runtime cfg hotRestartAuth)
set HOT_RESTART_AUTH=nGH66S4TjBjQqCKyWJAM

REM SSH port
set REMOTE_PORT=22