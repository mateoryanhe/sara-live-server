官网 SFTP 专用账号（正式服 52.9.70.64）
==========================================

账号信息
--------
  用户名: official-site
  主机:   52.9.70.64
  端口:   22
  协议:   SFTP（仅文件传输，不可 SSH 执行命令）
  私钥:   official-site-sftp.pem（本目录，勿提交 Git）

可操作目录
----------
  服务器路径: /home/ec2-user/cdn/official-site
  SFTP 登录后当前目录即为该目录（internal-sftp -d /files）

  可上传、删除该目录下文件；不可运行程序、不可访问 cms/images 等其它目录。

生成本地密钥（首次）
------------------
  在本目录执行:
    ssh-keygen -t ed25519 -f official-site-sftp -N "" -C "official-site-sftp@52.9.70.64"
    copy official-site-sftp official-site-sftp.pem

服务器创建 SFTP 账号（需 ec2-user sudo）
----------------------------------------
  scp keys/official-site-sftp.pub setup_sftp_user.sh ec2-user@52.9.70.64:/tmp/
  ssh ec2-user@52.9.70.64 "sed -i 's/\r$//' /tmp/setup_sftp_user.sh && sudo bash /tmp/setup_sftp_user.sh"

  若脚本报 pipefail 错误，多为 Windows CRLF，按上面 sed 处理后再执行。

连接示例
--------
  FileZilla / WinSCP:
    协议 SFTP, 主机 52.9.70.64, 用户 official-site
    认证选「密钥文件」，选 official-site-sftp.pem

  命令行:
    sftp -i official-site-sftp.pem -P 22 official-site@52.9.70.64

  注意: 正式服与测试服使用不同密钥；私钥请妥善保管。
