@echo off
setlocal EnableExtensions
cd /d "%~dp0"

REM ------------------------------------------------------------
REM Windows 环境一键脚本（仅 Win）
REM   setup.bat           加载到【当前 cmd】会话
REM   setup.bat install   写入【当前用户】永久环境变量（新开终端生效）
REM 先编辑同目录 config.bat 填写 Flutter / JDK / Android SDK 路径
REM ------------------------------------------------------------

call "%~dp0config.bat"
if errorlevel 1 exit /b 1

if not exist "%FLUTTER_HOME%\bin\flutter.bat" (
  echo [ERROR] Flutter 不存在: %FLUTTER_HOME%
  echo 请编辑 scripts\env\config.bat
  exit /b 1
)
if not exist "%JAVA_HOME%\bin\java.exe" (
  echo [ERROR] JDK 不存在: %JAVA_HOME%
  echo 请编辑 scripts\env\config.bat
  exit /b 1
)
if not exist "%ANDROID_SDK%" (
  echo [ERROR] Android SDK 不存在: %ANDROID_SDK%
  echo 请编辑 scripts\env\config.bat
  exit /b 1
)

if /I "%~1"=="install" goto :INSTALL
if /I "%~1"=="/install" goto :INSTALL
goto :SESSION

:SESSION
endlocal & (
  call "%~dp0config.bat"
  set "ANDROID_HOME=%ANDROID_SDK%"
  set "ANDROID_SDK_ROOT=%ANDROID_SDK%"
  set "PATH=%FLUTTER_HOME%\bin;%ANDROID_SDK%\platform-tools;%ANDROID_SDK%\cmdline-tools\latest\bin;%JAVA_HOME%\bin;%PATH%"
  if "%USE_CN_MIRROR%"=="1" (
    set "PUB_HOSTED_URL=https://pub.flutter-io.cn"
    set "FLUTTER_STORAGE_BASE_URL=https://storage.flutter-io.cn"
  )
)
if defined HTTP_PROXY (
  echo %HTTP_PROXY% | findstr /I "://" >nul
  if errorlevel 1 set "HTTP_PROXY=http://%HTTP_PROXY%"
)
if defined http_proxy (
  echo %http_proxy% | findstr /I "://" >nul
  if errorlevel 1 set "http_proxy=http://%http_proxy%"
)
if not defined NO_PROXY set "NO_PROXY=localhost,127.0.0.1,::1"

echo [OK] 当前会话环境已加载
echo FLUTTER_HOME=%FLUTTER_HOME%
echo JAVA_HOME=%JAVA_HOME%
echo ANDROID_HOME=%ANDROID_HOME%
where flutter 2>nul
where java 2>nul
echo.
echo 写入用户永久环境变量请执行: setup.bat install
exit /b 0

:INSTALL
powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0install.ps1"
exit /b %ERRORLEVEL%
