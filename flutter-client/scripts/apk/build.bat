@echo off
setlocal EnableExtensions
cd /d "%~dp0"

REM ------------------------------------------------------------
REM Windows 打 APK（仅 Win）
REM   build.bat           release（默认）
REM   build.bat release
REM   build.bat debug
REM 依赖同级 ..\env\config.bat 路径配置
REM ------------------------------------------------------------

set "MODE=%~1"
if "%MODE%"=="" set "MODE=release"
if /I "%MODE%"=="r" set "MODE=release"
if /I "%MODE%"=="d" set "MODE=debug"
if /I not "%MODE%"=="release" if /I not "%MODE%"=="debug" (
  echo [ERROR] 用法: build.bat [release^|debug]
  exit /b 1
)

REM 加载 Flutter / JDK / Android SDK 到当前会话
call "%~dp0..\env\setup.bat"
if errorlevel 1 exit /b 1

REM scripts\apk -> flutter-client 根目录
cd /d "%~dp0..\.."
if not exist "pubspec.yaml" (
  echo [ERROR] 未找到 pubspec.yaml，当前目录: %CD%
  exit /b 1
)

echo === flutter pub get ===
call flutter pub get
if errorlevel 1 exit /b 1

if /I "%MODE%"=="debug" (
  echo === flutter build apk --debug ===
  call flutter build apk --debug
  if errorlevel 1 exit /b 1
  echo.
  echo APK: %CD%\build\app\outputs\flutter-apk\app-debug.apk
) else (
  echo === flutter build apk --release ===
  call flutter build apk --release
  if errorlevel 1 exit /b 1
  echo.
  echo APK:
  dir /b build\app\outputs\flutter-apk\*.apk 2>nul
  echo Full path: %CD%\build\app\outputs\flutter-apk\app-release.apk
)

exit /b 0
