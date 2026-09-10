# Windows 打 APK（`scripts/apk`）

依赖 [`../env/config.bat`](../env/config.bat) 中的 Flutter / JDK / Android SDK 路径。

## 用法

| 操作 | 命令 |
|------|------|
| Release（默认） | `scripts\apk\build.bat` 或 `build.bat release` |
| Debug | `scripts\apk\build.bat debug` |
| 双击 Release | `一键打包.bat` |
| 双击 Debug | `一键打包-debug.bat` |

产物一般在：

- Release：`build\app\outputs\flutter-apk\app-release.apk`
- Debug：`build\app\outputs\flutter-apk\app-debug.apk`
