直播正式服 live_db 导入/导出工具
================================

前置
----
1. 将 pub-tool\go-正式服\ 下的 plink.exe、pscp.exe 复制到本目录
2. 确认 config.bat 中 SSH_KEY_PATH 指向正确的 .ppk 密钥

数据库用户（已在正式服创建）
----------------------------
  用户: live
  密码: 见 config.bat
  权限: 仅 live_db（SELECT/INSERT/UPDATE/DELETE/CREATE/DROP/ALTER 等，可 dump/import）
  主机: live@'%' 与 live@'localhost'

导出
----
  export.bat

  流程: SSH 登录正式服 -> mariadb-dump -> gzip -> 下载到 backup\ 目录

导入（会覆盖正式服 live_db 数据，慎用）
--------------------------------------
  import.bat backup\live_db_20260810_153000.sql.gz

远程直连说明
------------
MariaDB 当前 bind_address=127.0.0.1，端口 14501，外网无法直接连库。
export.bat / import.bat 通过 SSH 在服务器本机执行，无需开放 14501 端口。

若需本机 MariaDB 客户端直连，可先建 SSH 隧道:
  ssh -F %USERPROFILE%\.ssh\config -L 14501:127.0.0.1:14501 直播正式服

然后:
  mariadb-dump -h127.0.0.1 -P14501 -ulive -p live_db > live_db.sql

重新创建 live 用户
------------------
  scp grant_live_user.sql 到服务器后:
  mariadb -h127.0.0.1 -P14501 -uroot -p < grant_live_user.sql
