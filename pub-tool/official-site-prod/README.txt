官网一键上传说明（正式服 SFTP）

脚本说明
  upload.bat  - 打包 official-site 静态页并通过 SFTP 上传到正式服
  config.bat  - 服务器、密钥、API 地址等配置

前置条件
  1. Windows 已安装 OpenSSH Client（自带 sftp 命令）
  2. keys\official-site-sftp.pem 已生成并在服务器配置 SFTP 账号（见 keys\README.txt）
  3. config.bat 中 API_BASE_URL 可按需填写注销接口根地址，留空表示与官网同域

配置项（config.bat）
  REMOTE_HOST / REMOTE_PORT  - SFTP 连接
  SFTP_USER / SFTP_KEY_PATH  - SFTP 专用账号与私钥
  LOCAL_SITE_DIR             - 本地 official-site 目录
  API_BASE_URL               - 写入 js/site.js 的 apiBaseUrl

用法
  upload.bat

部署流程
  1. 复制 official-site 到 %TEMP% 临时目录
  2. 写入 js/site.js 中的 apiBaseUrl
  3. 通过 SFTP 逐文件上传（自动 mkdir 子目录）

说明
  远程目录 /home/ec2-user/cdn/official-site，与 config/prod/config.yaml 中 staticPaths 一致
  访问示例: https://www.saralive.net/official-site/index.html
