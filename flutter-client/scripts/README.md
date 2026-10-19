# Flutter client 脚本（Windows）

所有脚本集中在 `flutter-client/scripts/` 下。

| 目录 | 用途 |
|------|------|
| [`env/`](env/) | Flutter / JDK / Android SDK 环境变量 |
| [`web/`](web/) | Chrome Web 调试（观众 8080 / 主播 8081） |
| [`apk/`](apk/) | 打 release / debug APK、USB 安装 |

## 常用入口（在 `scripts\` 下双击）

| 文件 | 说明 |
|------|------|
| `build-apk.bat` | Release APK → `apk\build.bat` |
| `build-apk-debug.bat` | Debug APK |
| `install-apk.bat` | 安装 release 到 USB 真机 |
| `install-apk-debug.bat` | 安装 debug |
| `start-web.bat` | Web 观众（8080） |
| `start-web-audience.bat` | **Web 观众端（8080，常用）** |
| `start-web-anchor.bat` | Web 主播（8081） |
| `start-web-config.example.bat` | 示例 → 复制为 `web\start-web-config.bat`（含 `API_BASE`、日志目录） |
| `dev-usb.bat` | USB 真机 **flutter run**（热重载，日常改 UI 比打 release APK 快） |

根目录同名 `build-apk.bat`、`start-web-*.bat` 等仅为转发到上表，方便在 `flutter-client\` 根下双击。

## 首次使用

1. 编辑 [`env/config.bat`](env/config.bat) 中的 Flutter / JDK / Android SDK 路径。
2. `call scripts\env\setup.bat` 或双击 `env\一键设置环境.bat`（install 写入用户环境变量）。
3. Web：`scripts\start-web-audience.bat`；APK：`scripts\build-apk.bat`（可选 `scripts\apk\build-config.bat` 设 `API_BASE`）。
4. 真机安装：USB 调试连接后 `scripts\install-apk.bat`，或 `apk\一键打包并安装.bat`。

APK 产物：`build\app\outputs\flutter-apk\app-release.apk` / `app-debug.apk`。
