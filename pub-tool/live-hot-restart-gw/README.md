# live-hot-restart-gw（独立进程）

与 **xr-game-server 无关**，单独部署在直播正式服本机。监听 **18900**，收到带密钥的 HTTP 后：

| 路径 | 行为（对齐 `pub-tool/go-prod/restart.bat`） |
|------|---------------------------------------------|
| `/` 、`/hotRestart` | 无 `xr-game-server` 进程 → **`sudo start.sh` 冷启**；有进程 → 请求 `https://127.0.0.1/internal/hotRestart`；热重启失败且连不上/非 ok → **停旧进程 + 冷启**（需 `hotRestartFallbackCold`） |
| `/coldStart` | 始终冷启：有进程则先 SIGTERM/SIGKILL，再 `sudo start.sh` |

主 Go **进程已挂** 时，热重启接口不可用；走上述 **冷启分支**（与 SSH 执行 `start.sh` 相同）。

## 部署（正式服 ec2-user）

1. 本机 `build.bat` 得到 `live-hot-restart-gw`（Linux：`GOOS=linux GOARCH=amd64 go build -o live-hot-restart-gw .`）
2. 上传到 `/home/ec2-user/xgameserver/live-hot-restart-gw` + 同目录 `config.json`（从 `config.example.json` 复制）
3. **`auth`** = CMS「热重启密钥」；**`coldStartScript`** = 与 deploy 生成的 `start.sh` 同路径（默认 `/home/ec2-user/xgameserver/start.sh`）
4. ec2-user 需能 **`sudo` 无密码执行 `start.sh`**（与现有 restart.bat 一致）
5. AWS 安全组放行 **18900/TCP**（建议仅管理员 IP）
6. `nohup ./live-hot-restart-gw >> live-hot-restart-gw.log 2>&1 &` 或 systemd 常驻

## 手机 / 外网触发

热重启或智能重启（无进程则冷启）：

```http
GET http://54.177.166.38:18900/hotRestart?auth=你的热重启密钥
```

强制冷启：

```http
GET http://54.177.166.38:18900/coldStart?auth=你的热重启密钥
```

成功示例：`{"ok":true,"mode":"cold","pid":12345,"reason":"started"}` 或 `{"ok":true,"mode":"hot",...}`

## 安全

- 明文 HTTP，密钥会进 URL/日志；仅紧急运维，务必强密钥 + 限制 IP。
- `config.json` 勿提交 Git。
