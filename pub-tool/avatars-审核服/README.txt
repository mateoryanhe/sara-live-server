头像资源上传 - 审核服

用法
  upload.bat

远程目录
  /home/ec2-user/cdn/images

重要
  审核服 /tmp 为 tmpfs，仅约 479MB，放不下大压缩包。
  本脚本上传到 /home/ec2-user/staging，解压到 cdn/images。
  勿用 /tmp 传 800MB+ 文件。

前置
  plink.exe / pscp.exe 放在本目录
  config.bat 中 SSH_KEY_PATH 指向 live-review.ppk
  pub-tool/avatars/ 下放置头像文件
