官网一键上传说明（PuTTY + .ppk 密钥）

脚本说明
  upload.bat [test|prod]  - 打包 official-site 静态页并上传到远程服务器
  config.bat              - 服务器、密钥、目录、API 地址等配置

环境对应
  test  - 使用 API_BASE_URL_TEST 写入 js/site.js，部署到 REMOTE_DIR_TEST
  prod  - 使用 API_BASE_URL_PROD 写入 js/site.js，部署到 REMOTE_DIR_PROD

默认不传参数时为 test 环境。

前置条件
  1. 本目录放置 plink.exe、pscp.exe（PuTTY 工具，可与 cms-测试服 共用）
  2. config.bat 中 SSH_KEY_PATH 指向有效的 .ppk 私钥
  3. 远程服务器已安装 unzip
  4. 服务器 serverRoot 或 domainSites 已配置可访问官网目录

配置项（config.bat）
  REMOTE_HOST / REMOTE_USER / REMOTE_PORT  - SSH 连接
  SSH_KEY_PATH                             - PuTTY .ppk 密钥路径
  LOCAL_SITE_DIR                           - 本地 official-site 目录
  REMOTE_DIR_TEST / REMOTE_DIR_PROD        - 各环境远程部署目录
  API_BASE_URL_TEST / API_BASE_URL_PROD    - 注销账号等接口 API 根地址
                                             留空表示与官网同域（相对路径请求）

用法示例
  upload.bat           部署测试环境（默认）
  upload.bat test      部署测试环境
  upload.bat prod      部署生产环境

部署流程
  1. 复制 official-site 到 _staging 临时目录
  2. 按环境写入 js/site.js 中的 apiBaseUrl
  3. 压缩后 SCP 上传
  4. 远程解压到临时目录，仅覆盖包内文件（不 rm -rf 整目录，避免无关子目录权限导致失败）

说明
  远程目录 /home/ec2-user/cdn/official-site,与 config/dev staticSites 中 web.bigtktool.shop 一致
  API 根域 www.bigtktool.shop 不在 staticSites 中,默认走 Go API 路由
  访问示例: https://web.bigtktool.shop/

密钥格式
  若密钥为 OpenSSH (.pem)，请用 PuTTYgen 转为 .ppk 后写入 SSH_KEY_PATH。
