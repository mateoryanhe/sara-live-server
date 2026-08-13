头像资源一键上传说明（PuTTY + .ppk 密钥）

脚本说明
  upload.bat [test|prod]  - 上传本地头像到远程服务器
  config.bat              - 服务器、密钥、目录等配置（修改配置请编辑此文件）

本地头像目录
  pub-tool/avatars/
  将 demo 头像或业务头像文件放在该目录后执行 upload.bat

远程目录
  测试/生产默认: /home/ec2-user/cdn/images
  对应静态访问路径: /images/文件名.png（扩展名与上传文件一致）
  说明: 该目录可能由 sudo 启动的服务创建(属主 root),脚本会使用 sudo unzip 写入

前置条件
  1. 本目录放置 plink.exe、pscp.exe（PuTTY 工具，与 cms-测试服 相同）
  2. config.bat 中 SSH_KEY_PATH 指向有效的 .ppk 私钥
  3. ec2-user 具备免密 sudo（与 go-测试服 部署脚本相同）
  4. 远程 serverRoot 为 /home/ec2-user/cdn（见 config/test/config.yaml）
  5. 远程服务器已安装 unzip

配置项（config.bat）
  REMOTE_HOST / REMOTE_USER / REMOTE_PORT  - SSH 连接
  SSH_KEY_PATH                             - PuTTY .ppk 密钥路径
  LOCAL_AVATARS_DIR                        - 本地头像目录
  REMOTE_DIR_TEST / REMOTE_DIR_PROD        - 各环境远程目录

用法示例
  upload.bat           上传到测试服（默认）
  upload.bat test      上传到测试服
  upload.bat prod      上传到生产目录

说明
  - 脚本会压缩 pub-tool/avatars 下全部文件并覆盖上传到远程目录
  - 不会删除远程目录中未包含在本次压缩包里的其它文件
  - 数据库 avatar 字段可存文件名: demo_avatar_1.png（走 /images）
