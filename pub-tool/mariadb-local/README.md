# 本地 MariaDB — **sara-live-server**（Docker）

| 项 | 值 |
|---|---|
| 容器 | `sara-live-mariadb` |
| 端口 | **3391** |
| 库 / 用户 | `live_db` / `appuser` |

与语聊仓 **yuliao** 的 `yuliao-mariadb:3390` 分开，可同时 `go-local` + `cms-local`。

双仓一键起 DB：`D:\company-code\yl\pub-tool\mariadb-local\start-both-db.bat`