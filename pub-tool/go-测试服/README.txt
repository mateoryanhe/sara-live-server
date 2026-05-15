XR Game Server 一键部署脚本 (PuTTY版)
================================

概述
--------------------------------
本脚本用于将XR Game Server从Windows环境编译并部署到远程Linux服务器。
使用PuTTY工具链进行SSH连接和文件传输。

部署流程
--------------------------------
1. 编译Go程序为Linux可执行文件
2. 打包可执行文件和配置文件
3. 连接远程服务器并停止旧程序 (先使用kill -15，等待5秒，再使用kill -9)
4. 上传新的部署包到远程服务器
5. 解压并在远程服务器上启动新程序

文件说明
--------------------------------
- config.bat: 配置文件，包含服务器信息、路径等
- deploy.bat: 主要的批处理部署脚本
- 一键部署.bat: 用户友好的启动脚本
- README.txt: 本说明文件
- plink.exe: PuTTY SSH客户端工具
- pscp.exe: PuTTY 安全文件传输工具

使用前准备
--------------------------------
1. 确保已安装Go编译环境
2. 确保PuTTY工具链已放置在本目录 (plink.exe 和 pscp.exe)
3. 确保SSH密钥文件存在 (D:\prod-key\ytt.pem)
4. 确保SSH密钥为PuTTY兼容格式(.ppk)，如果不是，请按以下步骤转换
5. 确保源代码路径正确 (D:\go-project\xrgameserver\go-src)
6. 确保配置文件路径正确 (D:\go-project\xrgameserver\config\dev)

SSH密钥转换说明
--------------------------------
如果您的密钥是OpenSSH格式(.pem)，需要转换为PuTTY格式(.ppk)：
1. 下载PuTTYgen工具
2. 运行puttygen.exe
3. 点击"Load"加载您的ytt.pem文件
4. 如果提示密钥是OpenSSH格式，点击"Yes"确认转换
5. 点击"Save private key"保存为ytt.ppk文件
6. 将ytt.ppk重命名为ytt.pem或更新config.bat中的SSH_KEY_PATH指向新文件

配置说明
--------------------------------
在config.bat中可配置以下参数：
- REMOTE_HOST: 远程服务器IP地址
- REMOTE_USER: 远程服务器用户名
- REMOTE_DIR: 远程服务器部署目录
- LOCAL_PROJECT_PATH: 本地项目路径
- LOCAL_GO_SRC: 本地Go源代码路径
- LOCAL_CONFIG_PATH: 本地配置文件路径
- SSH_KEY_PATH: SSH密钥路径
- APP_NAME: 应用程序名称
- SUDO_CMD: 远程sudo命令

使用方法
--------------------------------
1. 双击"一键部署.bat"运行脚本
2. 脚本将自动完成整个部署过程

注意事项
--------------------------------
- 确保防火墙允许SSH连接 (端口22)
- 确保远程服务器允许sudo权限执行程序
- 部署过程中请勿关闭命令行窗口
- 部署完成后可在远程服务器查看 /var/log/xr-game-server.log 日志
- 程序和配置文件将部署到 /home/ec2-user 目录下