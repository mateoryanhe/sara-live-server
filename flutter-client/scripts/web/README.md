# Web 调试（Chrome）

| 操作 | 推荐双击 |
|------|----------|
| **观众端（8080）** | **`scripts\start-web-audience.bat`**（同目录下 `scripts\web\start-web-audience.bat` 等价） |
| 主播端（8081） | `scripts\start-web-anchor.bat` |
| 默认启动观众 | `scripts\start-web.bat` |

**API 与日志**：编辑 `scripts\web\start-web-config.bat`（无则从 `start-web-config.example.bat` 复制）。  
- `API_BASE` — 后端根地址，如 `http://127.0.0.1:9443` 或 `https://www.bigtktool.shop`  
- `HTTP_LOG_DIR_AUDIENCE` / `HTTP_LOG_DIR_ANCHOR` — HTTP 调试日志目录  

该配置文件在 `.gitignore`，每人本地可不同。运行时仍可在 App **我的 → API 服务器** 覆盖（若已实现）。

环境：优先 `scripts\env\setup.bat`；若无则回退 `D:\tools\env-flutter.bat`。

根目录与 `scripts\start-web*.bat` 均为转发到 `scripts\web\`，方便双击。

### `Account Web Data` / profile 被占用（errno 1224）

同一 Chrome 用户目录只能被一个进程使用（`%LOCALAPPDATA%\SaraLiveWeb\profile-audience` 等）。

- 先 **关掉** 上一次 `flutter run` 弹出的 Sara Live Chrome 窗口，或在启动器 CMD 里 **Ctrl+C** 结束。
- 不要 **同时开两个** 观众端（两个 8080）或共用同一 profile 的主播/观众脚本。
- 重新双击 `start-web-audience.bat`；脚本会尝试结束仍占用该 profile 的 Chrome/Chromium 并清理 Singleton 锁文件。
- 仍失败：任务管理器结束所有 **Google Chrome**（若窗口标题含 SaraLive），或暂时关掉杀毒软件对 `%LOCALAPPDATA%\SaraLiveWeb` 的实时扫描后再试。
