# Windows 环境脚本（`scripts/env`）

## 配置

编辑 [`config.bat`](config.bat)：

```bat
set "FLUTTER_HOME=D:\tools\flutter"
set "JAVA_HOME=D:\tools\jdk-17"
set "ANDROID_SDK=D:\tools\android-sdk"
set "USE_CN_MIRROR=1"
```

## 用法

| 操作 | 命令 |
|------|------|
| 当前 cmd 会话临时生效 | `call scripts\env\setup.bat` |
| 写入用户永久环境变量 | `scripts\env\setup.bat install` 或双击 `一键设置环境.bat` |

永久写入后需**新开**终端 / IDE，再执行 `flutter doctor`。
