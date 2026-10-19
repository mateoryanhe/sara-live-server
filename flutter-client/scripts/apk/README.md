# Windows 打 APK（`scripts/apk`）

依赖 [`../env/config.bat`](../env/config.bat) 中的 Flutter / JDK / Android SDK 路径。

## 用法

| 操作 | 命令 |
|------|------|
| Release（默认） | `scripts\build-apk.bat` 或 `scripts\apk\build.bat` |
| Debug | `scripts\build-apk-debug.bat` 或 `scripts\apk\build.bat debug` |
| 指定 API 根地址 | `build.bat release http://192.168.2.11:9443` |
| 仅 URL（release） | `build.bat http://192.168.2.11:9443` |
| 双击 Release | `一键打包.bat` |
| 双击 Debug | `一键打包-debug.bat` |
| **安装到 USB 手机** | `scripts\install-apk.bat` 或 `scripts\apk\install.bat` |
| 双击安装 Release | `scripts\apk\一键安装.bat` |
| 双击安装 Debug | `scripts\apk\一键安装-debug.bat` |
| 打包 + 安装 | `scripts\apk\一键打包并安装.bat` |

### 安装到真机（adb）

1. 编辑 [`../env/config.bat`](../env/config.bat)，保证 `ANDROID_SDK` 正确（内含 `platform-tools\adb.exe`）。
2. 手机：**开发者选项 → USB 调试** 打开，USB 连接电脑，弹窗点允许。
3. 先要有 APK（`build.bat` 或 `一键打包.bat`），再运行安装脚本。
4. 多台手机时：先 `adb devices` 看序列号，再 `set ANDROID_SERIAL=序列号` 后执行 `install.bat`。

命令示例：

```bat
scripts\apk\install.bat
scripts\apk\install.bat debug
scripts\apk\install.bat "D:\path\custom.apk"
```

默认安装 `build\app\outputs\flutter-apk\app-release.apk`（与 `APK_PACKAGING=arm64` 一致；`split` 时优先 `app-arm64-v8a-release.apk`）。

### API_BASE（编译进 APK 的默认服务器）

1. 复制 `build-config.example.bat` → `build-config.bat`，设置 `API_BASE=...`（推荐内网打包）。
2. 或命令行第二个参数 / 单独传 URL（见上表）。
3. 未设置时使用 `lib/app/app_config.dart` 里的编译默认值；安装后仍可在 App **我的 → API 服务器** 覆盖。

打包时会传 `--dart-define=API_BASE=...` 给 Flutter。

产物一般在 `build\app\outputs\flutter-apk\`：

| 模式 | 文件 | 说明 |
|------|------|------|
| **arm64**（默认） | `app-release.apk` | 仅 64 位手机，体积约为「通用包」的 1/3 |
| split | `app-arm64-v8a-release.apk` 等 3 个 | 按 CPU 分包 |
| fat / universal | `app-release.apk` | 含 arm32+arm64+x86，最大（Agora 会很大） |

在 `build-config.bat` 里设 `APK_PACKAGING=fat` 可恢复旧的全架构单包。

Release 已开 R8 压缩（`proguard-rules.pro`）。Debug：`app-debug.apk`。

**Agora 体积**：Dart 侧只用进房/推流/订阅/预览，未用美颜、空间音频、屏幕共享等；`android/app/build.gradle.kts` 的 `packaging.jniLibs.excludes` 会去掉 `*_extension.so`（仍保留 `libagora-rtc-sdk.so` 等核心库）。若以后要开某扩展，从 excludes 里去掉对应 so 即可。

## 常见构建错误

**`agora_rtc_engine:checkReleaseAarMetadata` / compileSdk 31**

工程已在 `android/build.gradle.kts` 为所有 Flutter 插件强制 `compileSdk = 36`。请在 Android SDK Manager 安装 **Android 14/15 API 34–36** 平台包。

**Kotlin incremental cache / `Storage is already registered`**

1. 关闭 Android Studio、其它正在编译的终端。
2. 执行（推荐）：

```bat
scripts\apk\build.bat clean release
```

或 `scripts\apk\clean.bat` 后再 `build.bat`。

工程已在 `android/gradle.properties` 关闭 Kotlin 增量编译并限制并行；仍失败时确认 `where java` 首项为 `scripts\env\config.bat` 里的 **JDK 17**，不要混用 JDK 25。
