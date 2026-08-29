@echo off
REM 头像资源上传配置 - 审核服（使用 PuTTY .ppk 密钥）

REM 远程服务器（直播审核服）
set REMOTE_HOST=18.144.165.177
set REMOTE_USER=ec2-user
set REMOTE_PORT=22

REM SSH 密钥（PuTTY 格式 .ppk）
set SSH_KEY_PATH=D:\tools\ppk\live-review.ppk

REM 本地路径
set LOCAL_PROJECT_ROOT=D:\company-code\sara-live-server
set LOCAL_AVATARS_DIR=%LOCAL_PROJECT_ROOT%\pub-tool\avatars

REM 远程头像目录；审核服勿用 /tmp（tmpfs 仅 ~479MB）
set REMOTE_DIR=/home/ec2-user/cdn/images
set REMOTE_STAGE=/home/ec2-user/staging
