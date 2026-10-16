# live-hot-restart-gw（独立进程）

与 **xr-game-server 无关**，单独部署在直播正式服本机。监听 **18900**，收到带密钥的 HTTP 后，请求本机：

`https://127.0.0.1/internal/hotRestart?auth=…`

从而触发主程序热重启（与 `pub-tool/go-prod/restart.bat` 一致）。

## 部署（正式服 ec2-user）

1. 本机 `build.bat` 得到 `live-hot-restart-gw`（Linux 在 `go-src` 外执行  
   `GOOS=linux GOARCH=amd64 go build -o live-hot-restart-gw .`）
2. 上传到例如 `/home/ec2-user/xgameserver/live-hot-restart-gw` + 同目录 `config.json`（从 `config.example.json` 复制）
3. `config.json` 中 **`auth`** = CMS「热重启密钥」；**`listen`** 默认 `0.0.0.0:18900`
4. AWS 安全组放行 **18900/TCP**（建议仅管理员 IP）
5. 用 systemd 或 `nohup ./live-hot-restart-gw &` 常驻（主 Go 进程重启时本程序仍可单独运行）

## 手机触发

```http
GET http://54.177.166.38:18900/hotRestart?auth=你的热重启密钥
```

成功：`{"ok":true,"reason":"{\"ok\":true,\"reason\":\"accepted\"}",...}`

## 安全

- 明文 HTTP，密钥会进 URL/日志；仅紧急运维，务必强密钥 + 限制 IP。
- `config.json` 勿提交 Git。
