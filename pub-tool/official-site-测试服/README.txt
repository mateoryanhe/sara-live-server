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
  4. 远程删除 REMOTE_DIR 后重新解压（目录独立,不影响 cms / upload）

说明
  REMOTE_DIR 默认 /home/ec2-user/cdn/official-site,与 server.staticPaths 中官网目录一致
  cms、images 已拆分到独立目录,清空官网目录不会影响其他资源

访问路径
  默认远程目录 /home/ec2-user/cdn/official-site
  若 serverRoot 为 /home/ec2-user/cdn，则访问：
    https://你的域名/official-site/index.html

  若使用 domainSites 将独立域名指到该目录，可改为部署到 site 根目录并调整 REMOTE_DIR_*。

密钥格式
  若密钥为 OpenSSH (.pem)，请用 PuTTYgen 转为 .ppk 后写入 SSH_KEY_PATH。
