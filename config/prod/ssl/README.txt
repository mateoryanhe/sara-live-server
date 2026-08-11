Cloudflare Origin CA 证书 (正式服)
================================

文件说明
- cf-origin.pem  Cloudflare 签发的源站证书 (对应 yuan.txt)
- cf-origin.key  源站私钥 (对应 siyao.txt)

服务器路径 (与 config/prod/config.yaml 一致)
- /home/ec2-user/ssl/cf-origin/cf-origin.pem
- /home/ec2-user/ssl/cf-origin/cf-origin.key

Cloudflare 控制台
- SSL/TLS -> Overview: Full (strict)
- 源站证书仅用于 Cloudflare 到源站 HTTPS,浏览器不信任此证书

部署
- 运行 pub-tool/go-正式服/deploy.bat 会自动上传证书到服务器
- 或手动:
  mkdir -p /home/ec2-user/ssl/cf-origin
  chmod 644 cf-origin.pem
  chmod 600 cf-origin.key

证书域名: *.saralive.net, saralive.net
