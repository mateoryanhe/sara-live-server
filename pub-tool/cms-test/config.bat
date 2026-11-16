@echo off
REM CMS 部署配置（语聊测试服；使用 PuTTY .ppk，勿在此文件填写密码）

REM 远程服务器（语聊测试服 ec2-user @ 63.185.66.53）
set REMOTE_HOST=63.185.66.53
set REMOTE_USER=ec2-user
set REMOTE_PORT=22
set REMOTE_STAGE=/home/ec2-user/staging

REM SSH 密钥（PuTTY .ppk，由 D:\root\yuliao.pem 用 puttygen 转换）
set SSH_KEY_PATH=D:\tools\ppk\yuliao.ppk

REM 本地路径
set LOCAL_PROJECT_ROOT=D:\company-code\sara-live-server
set VUE_PROJECT_DIR=%LOCAL_PROJECT_ROOT%\cms
set BUILD_OUTPUT_DIR=D:\root\cms

REM 各环境远程解压目录（可按实际服务器调整）
set REMOTE_DIR_TEST=/home/ec2-user/cdn/cms
set REMOTE_DIR_PROD=/home/ec2-user/cdn/cms
