@echo off
setlocal EnableExtensions EnableDelayedExpansion
chcp 65001 >nul 2>&1
cd /d "%~dp0"

REM USB dev: flutter run with hot reload. Use build.bat for release APK.

set "EXIT_CODE=0"

if exist "%~dp0build-config.bat" call "%~dp0build-config.bat"

set "SETUP_QUIET=1"
call "%~dp0..\env\setup.bat"
if errorlevel 1 (
  echo [FAIL] env setup
  set "EXIT_CODE=1"
  goto :Done
)
if defined JAVA_HOME set "PATH=%JAVA_HOME%\bin;%PATH%"

cd /d "%~dp0..\.."
if not exist "pubspec.yaml" (
  echo [ERROR] pubspec.yaml not found: %CD%
  set "EXIT_CODE=1"
  goto :Done
)

set "DEVICE_ARG="
if not "%~1"=="" (
  set "DEVICE_ARG=-d %~1"
  goto :RunFlutter
)

where adb >nul 2>&1
if errorlevel 1 goto :RunFlutter

set "ADB_COUNT=0"
set "ADB_SERIAL="
for /f "skip=1 tokens=1,2" %%A in ('adb devices 2^>nul') do (
  if "%%B"=="device" (
    set /a ADB_COUNT+=1
    if not defined ADB_SERIAL set "ADB_SERIAL=%%A"
  )
)
if !ADB_COUNT!==1 (
  set "DEVICE_ARG=-d !ADB_SERIAL!"
  echo [info] Using USB device: !ADB_SERIAL!
) else if !ADB_COUNT! GTR 1 (
  if defined ANDROID_SERIAL (
    set "DEVICE_ARG=-d !ANDROID_SERIAL!"
    echo [info] Using ANDROID_SERIAL=!ANDROID_SERIAL!
  ) else (
    echo [WARN] Multiple USB devices. Run: dev-usb.bat ^<serial^>
    adb devices
    set "EXIT_CODE=1"
    goto :Done
  )
)

:RunFlutter
echo.
echo ========================================
echo   Sara Live - USB dev ^(flutter run^)
echo ========================================
echo   Hot reload: r   Hot restart: R   Quit: q
echo ========================================
echo.

call flutter devices
echo.

if not defined DEVICE_ARG (
  echo [info] Pick device in terminal, or connect one USB phone only.
  echo.
)

call flutter run !DEVICE_ARG!
if errorlevel 1 set "EXIT_CODE=1"

:Done
if not "!EXIT_CODE!"=="0" (
  echo.
  echo [FAIL] dev-usb failed. Fix errors above and retry.
)
exit /b !EXIT_CODE!
