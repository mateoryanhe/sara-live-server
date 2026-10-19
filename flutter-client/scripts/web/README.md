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
