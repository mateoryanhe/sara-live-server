CMS 审核服一键构建上传说明（PuTTY + .ppk 密钥）

脚本说明
  upload.bat  - 构建 Vue（review 环境）并上传到审核服
  config.bat  - 服务器、密钥、目录等配置

环境对应
  review  - npm run build:review  -> 加载 cms/.env.review
            API: https://v1.saralive.net
            CMS: https://review.saralive.net/

前置条件
  1. 本目录放置 plink.exe、pscp.exe（PuTTY 工具，可从 cms-test 复制）
  2. config.bat 中 SSH_KEY_PATH 指向有效的 .ppk 私钥
  3. 已安装 Node.js / npm，且 cms 目录可正常构建
  4. 审核服已安装 unzip

配置项（config.bat）
  REMOTE_HOST=18.144.165.177   直播审核服
  REMOTE_DIR=/home/ec2-user/cdn/cms
  SSH_KEY_PATH=D:\tools\ppk\1v1-review.ppk  审核服专用（勿用 live-test.ppk）

密钥说明
  审核服 SSH 使用 1v1.pem，已转换为 D:\tools\ppk\1v1-review.ppk。
  plink.exe / pscp.exe 已放入本目录。
