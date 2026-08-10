官网 SFTP 专用账号（测试服 54.241.124.37）
==========================================

账号信息
--------
  用户名: official-site
  主机:   54.241.124.37
  端口:   22
  协议:   SFTP（仅文件传输，不可 SSH 执行命令）
  私钥:   official-site-sftp.pem（本目录）

可操作目录
----------
  服务器路径: /home/ec2-user/cdn/official-site
  SFTP 登录后当前目录即为该目录（internal-sftp -d /files）

  可上传、删除该目录下文件；不可运行程序、不可访问 cms/images 等其它目录。

连接示例
--------
  FileZilla / WinSCP:
    协议 SFTP, 主机 54.241.124.37, 用户 official-site
    认证选「密钥文件」，选 official-site-sftp.pem

  命令行:
    sftp -i official-site-sftp.pem -P 22 official-site@54.241.124.37

  注意: 首次连接需接受主机指纹；私钥文件请妥善保管，勿提交 Git。

服务器重装/重建账号
------------------
  在本机执行（需 ec2-user sudo 权限）:
    scp keys/official-site-sftp.pub setup_sftp_user.sh 直播测试服:/tmp/
    ssh 直播测试服 "sudo bash /tmp/setup_sftp_user.sh"

  SSH 别名见 ~/.ssh/config 中 Host 直播测试服
