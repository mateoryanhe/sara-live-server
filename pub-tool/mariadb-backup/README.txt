live_db 导入/导出
================

通过 SSH 在远程服务器本机执行 mariadb-dump / mariadb，再下载或上传备份文件。
无需开放 MariaDB 端口，无需 SSH 隧道。

配置文件（分开维护）
--------------------
  config-export.bat   导出源（默认：直播正式服 52.9.70.64:14501）
  config-import.bat   导入目标（默认：直播测试服 54.241.124.37:13307）

  修改 SSH、数据库密码、端口时，只改对应配置文件即可。

用法
----
  导出（从 config-export.bat 指定的服务器）:

    export.bat

  备份默认保存到 config-export.bat 中的 LOCAL_BACKUP_DIR
  文件名: live_db_yyyyMMdd_HHmmss.sql.gz

  导入（到 config-import.bat 指定的服务器）:

    import.bat D:\var\live_db_20260810_153000.sql.gz

  或在 config-import.bat 设置 IMPORT_SQL_FILE 后直接运行:

    import.bat

  命令行参数优先于 IMPORT_SQL_FILE。

说明
----
- 典型流程：正式服 export.bat -> 测试服 import.bat
- backup / D:\var 含生产数据，请勿提交 Git
- plink.exe、pscp.exe 需在本目录（可从 go-prod 复制）
