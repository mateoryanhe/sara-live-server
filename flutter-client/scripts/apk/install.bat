@echo off

setlocal EnableExtensions EnableDelayedExpansion

chcp 65001 >nul 2>&1

cd /d "%~dp0"



REM ------------------------------------------------------------

REM 安装 APK 到 USB 已连接的真机（adb install -r）

REM   install.bat              release（默认 app-release.apk）

REM   install.bat debug        app-debug.apk

REM   install.bat "D:\path\app.apk"

REM 需：手机开 USB 调试、已授权本机；platform-tools 在 config.bat 的 ANDROID_SDK 下

REM ------------------------------------------------------------



set "MODE=release"

set "APK_ARG=%~1"



if /I "%APK_ARG%"=="debug" set "MODE=debug" & set "APK_ARG="

if /I "%APK_ARG%"=="d" set "MODE=debug" & set "APK_ARG="

if /I "%APK_ARG%"=="release" set "MODE=release" & set "APK_ARG="

if /I "%APK_ARG%"=="r" set "MODE=release" & set "APK_ARG="



if exist "%~dp0build-config.bat" call "%~dp0build-config.bat"

if not defined APK_PACKAGING set "APK_PACKAGING=arm64"



set "SETUP_QUIET=1"

call "%~dp0..\env\setup.bat"

if errorlevel 1 exit /b 1



cd /d "%~dp0..\.."

if not exist "pubspec.yaml" (

  echo [ERROR] pubspec.yaml not found: %CD%

  exit /b 1

)



set "APK_DIR=%CD%\build\app\outputs\flutter-apk"

set "APK_FILE="



if not "%APK_ARG%"=="" (

  if exist "%APK_ARG%" (

    set "APK_FILE=%APK_ARG%"

  ) else (

    echo [ERROR] APK not found: %APK_ARG%

    exit /b 1

  )

) else if /I "%MODE%"=="debug" (

  set "APK_FILE=%APK_DIR%\app-debug.apk"

) else if /I "%APK_PACKAGING%"=="split" (

  if exist "%APK_DIR%\app-arm64-v8a-release.apk" (

    set "APK_FILE=%APK_DIR%\app-arm64-v8a-release.apk"

  ) else (

    set "APK_FILE=%APK_DIR%\app-release.apk"

  )

) else (

  set "APK_FILE=%APK_DIR%\app-release.apk"

)



where adb >nul 2>&1

if errorlevel 1 (

  echo [ERROR] adb not in PATH. Check ANDROID_SDK in scripts\env\config.bat

  exit /b 1

)



echo.

echo ========================================

echo   Sara Live - Install APK ^(USB^)

echo ========================================

echo   APK : %APK_FILE%

echo ========================================

echo.



if not exist "%APK_FILE%" (

  echo [ERROR] APK 不存在，请先打包：

  echo         scripts\apk\build.bat release

  echo         或双击 一键打包.bat

  exit /b 1

)



echo [1/3] adb devices ...

adb devices

echo.



set "DEVICE_COUNT=0"

set "DEVICE_SERIAL="

for /f "skip=1 tokens=1,2" %%A in ('adb devices 2^>nul') do (

  if "%%B"=="device" (

    set /a DEVICE_COUNT+=1

    if not defined DEVICE_SERIAL set "DEVICE_SERIAL=%%A"

  )

)



if !DEVICE_COUNT!==0 (

  echo [ERROR] 未检测到已授权的设备。

  echo         - 手机用 USB 连接电脑

  echo         - 设置里开启「USB 调试」

  echo         - 弹窗点「允许 USB 调试」

  echo         - 若显示 unauthorized，拔线重插并在手机上确认

  exit /b 1

)



if !DEVICE_COUNT! GTR 1 (

  if defined ANDROID_SERIAL (

    set "DEVICE_SERIAL=%ANDROID_SERIAL%"

    echo       多台设备，使用 ANDROID_SERIAL=%ANDROID_SERIAL%

  ) else (

    echo [WARN] 检测到 !DEVICE_COUNT! 台设备，将安装到第一台: !DEVICE_SERIAL!

    echo        指定设备: set ANDROID_SERIAL=序列号 后再运行本脚本

  )

)



set "ADB_TARGET="

if defined DEVICE_SERIAL set "ADB_TARGET=-s !DEVICE_SERIAL!"



echo [2/3] Installing ...

set "INSTALL_FLAGS=-r"

if /I "%MODE%"=="debug" set "INSTALL_FLAGS=-r -t"



adb %ADB_TARGET% install %INSTALL_FLAGS% "%APK_FILE%"

set "INSTALL_EC=%ERRORLEVEL%"



if not "%INSTALL_EC%"=="0" (

  echo.

  echo [FAIL] adb install 失败 ^(exit %INSTALL_EC%^)

  echo        常见原因：签名不一致 ^(先卸载旧包^)、空间不足、包损坏

  echo        卸载: adb %ADB_TARGET% uninstall net.saralive.sara_live_app

  exit /b 1

)



echo.

echo [3/3] Done.

for %%A in ("%APK_FILE%") do echo       Installed : %%~nxA ^(%%~zA bytes^)

echo       Package   : net.saralive.sara_live_app

echo.

echo ========================================

echo   INSTALL OK

echo ========================================

echo.

exit /b 0

