@echo off
REM 头像资源上传配置（使用 PuTTY .ppk 密钥，勿在此文件填写密码）

REM 远程服务器
set REMOTE_HOST=54.241.124.37
set REMOTE_USER=ec2-user
set REMOTE_PORT=22

REM SSH 密钥（PuTTY 格式 .ppk）
set SSH_KEY_PATH=D:\tools\ppk\live-test.ppk

REM 本地路径
set LOCAL_PROJECT_ROOT=D:\company-code\sara-live-server
set LOCAL_AVATARS_DIR=%LOCAL_PROJECT_ROOT%\pub-tool\avatars

REM 远程头像目录（对应 staticPaths /images）
set REMOTE_DIR_TEST=/home/ec2-user/cdn/images
set REMOTE_DIR_PROD=/home/ec2-user/cdn/images
