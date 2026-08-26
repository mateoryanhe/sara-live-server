CMS 一键构建上传说明（PuTTY + .ppk 密钥）

脚本说明
  upload.bat [test|prod]  - 构建 Vue 并上传到远程服务器
  config.bat              - 服务器、密钥、目录等配置（修改配置请编辑此文件）

环境对应
  test  - npm run build:test  -> 加载 cms/.env.test
  prod  - npm run build:prod  -> 加载 cms/.env.production

默认不传参数时为 test 环境。

前置条件
  1. 本目录放置 plink.exe、pscp.exe（PuTTY 工具）
  2. config.bat 中 SSH_KEY_PATH 指向有效的 .ppk 私钥
  3. 已安装 Node.js / npm，且 cms 目录可正常构建
  4. 远程服务器已安装 unzip

配置项（config.bat）
  REMOTE_HOST / REMOTE_USER / REMOTE_PORT  - SSH 连接
  SSH_KEY_PATH                             - PuTTY .ppk 密钥路径
  VUE_PROJECT_DIR / BUILD_OUTPUT_DIR       - 项目与构建输出目录
  REMOTE_DIR_TEST / REMOTE_DIR_PROD        - 各环境远程部署目录

部署流程
  1. npm run build:test / build:prod
  2. 压缩 D:\root\cms 构建产物（优先 tar，路径兼容 Linux）
  3. 上传到 /tmp 后解压到临时目录
  4. 仅覆盖包内顶层文件/目录到 REMOTE_DIR，保留 log-query-export 等服务器目录
  5. 校验远程 index.html 存在

注意
  旧脚本会先 rm -rf 整个 REMOTE_DIR，若解压失败会导致 /cms 目录被清空。
  服务器 domainSites 映射: cms.bigtktool.shop -> /home/ec2-user/cdn/cms
  访问示例: https://cms.bigtktool.shop/

密钥格式
  若密钥为 OpenSSH (.pem)，请用 PuTTYgen 转为 .ppk 后写入 SSH_KEY_PATH。
