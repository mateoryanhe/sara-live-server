@echo off
REM CMS 部署配置（使用 PuTTY .ppk 密钥，勿在此文件填写密码）

REM 远程服务器
set REMOTE_HOST=50.18.253.123
set REMOTE_USER=ec2-user
set REMOTE_PORT=22

REM SSH 密钥（PuTTY 格式 .ppk）
set SSH_KEY_PATH=D:\tools\ppk\live-test.ppk

REM 本地路径
set LOCAL_PROJECT_ROOT=D:\company-code\app-code\live-server
set VUE_PROJECT_DIR=%LOCAL_PROJECT_ROOT%\cms
set BUILD_OUTPUT_DIR=D:\root\cms

REM 各环境远程解压目录（可按实际服务器调整）
set REMOTE_DIR_TEST=/home/ec2-user/cdn/cms
set REMOTE_DIR_PROD=/home/ec2-user/cdn/cms
