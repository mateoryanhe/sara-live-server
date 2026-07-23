@echo off
REM 官网部署配置（使用 PuTTY .ppk 密钥，勿在此文件填写密码）

REM 远程服务器
set REMOTE_HOST=50.18.253.123
set REMOTE_USER=ec2-user
set REMOTE_PORT=22

REM SSH 密钥（PuTTY 格式 .ppk）
set SSH_KEY_PATH=D:\tools\ppk\live-test.ppk

REM 本地路径
set LOCAL_PROJECT_ROOT=D:\company-code\sara-live-server
set LOCAL_SITE_DIR=%LOCAL_PROJECT_ROOT%\official-site
set STAGING_DIR=%~dp0_staging

REM 各环境远程解压目录（可按实际服务器 / 域名静态根目录调整）
set REMOTE_DIR_TEST=/home/ec2-user/cdn/official-site
set REMOTE_DIR_PROD=/home/ec2-user/cdn/official-site

REM 注销账号等接口的 API 根地址（写入 js/site.js 的 apiBaseUrl）
REM 留空表示与官网同域；测试服可填如 https://cai.hzaicoin.fun
set API_BASE_URL_TEST=
set API_BASE_URL_PROD=
