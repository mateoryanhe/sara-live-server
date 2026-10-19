@echo off
setlocal EnableExtensions
chcp 65001 >nul 2>&1
cd /d "%~dp0"

REM ------------------------------------------------------------
REM Windows APK build
REM   build.bat [release|debug] [API_BASE_URL]
REM   build.bat http://192.168.2.11:9443
REM Optional: build-config.bat (copy from build-config.example.bat)
REM ------------------------------------------------------------

set "DEFAULT_API_BASE=https://www.bigtktool.shop"

if exist "%~dp0build-config.bat" (
  call "%~dp0build-config.bat"
  if defined API_BASE set "API_SOURCE=build-config.bat"
)

set "DO_CLEAN=0"
if /I "%~1"=="clean" (
  set "DO_CLEAN=1"
  set "MODE=%~2"
  set "API_BASE_ARG=%~3"
) else (
  set "MODE=%~1"
  set "API_BASE_ARG=%~2"
)
if not defined API_SOURCE set "API_SOURCE=app-default"

if not "%API_BASE_ARG%"=="" (
  set "API_BASE=%API_BASE_ARG%"
  set "API_SOURCE=command-line"
)

if /I "%MODE:~0,7%"=="http://" goto :arg1_is_api
if /I "%MODE:~0,8%"=="https://" goto :arg1_is_api
goto :after_arg1_api
:arg1_is_api
set "API_BASE=%MODE%"
set "MODE=release"
set "API_SOURCE=command-line"
:after_arg1_api

if "%MODE%"=="" set "MODE=release"
if /I "%MODE%"=="r" set "MODE=release"
if /I "%MODE%"=="d" set "MODE=debug"
if /I not "%MODE%"=="debug" if /I not "%MODE%"=="release" (
  echo [ERROR] Usage: build.bat [clean] [release^|debug] [API_BASE_URL]
  echo        Example: build.bat clean release
  exit /b 1
)

if defined API_BASE call :NormalizeApiBase
goto :after_normalize
:NormalizeApiBase
:trim_api_commas
if "%API_BASE:~-1%"=="," set "API_BASE=%API_BASE:~0,-1%" & goto trim_api_commas
if "%API_BASE:~-1%"=="/" set "API_BASE=%API_BASE:~0,-1%" & goto trim_api_commas
exit /b 0
:after_normalize

set "DART_DEFINES="
set "EFFECTIVE_API=%DEFAULT_API_BASE%"
set "EFFECTIVE_WS=wss://www.bigtktool.shop"
if defined API_BASE (
  set "DART_DEFINES=--dart-define=API_BASE=%API_BASE%"
  set "EFFECTIVE_API=%API_BASE%"
  set "EFFECTIVE_WS=%API_BASE:http://=ws://%"
  set "EFFECTIVE_WS=%EFFECTIVE_WS:https://=wss://%"
)

REM APK size: arm64 (default, phones) | split (3 ABIs) | fat (universal, largest)
if not defined APK_PACKAGING set "APK_PACKAGING=arm64"
set "APK_FLUTTER_FLAGS="
if /I "%APK_PACKAGING%"=="arm64" set "APK_FLUTTER_FLAGS=--target-platform android-arm64"
if /I "%APK_PACKAGING%"=="split" set "APK_FLUTTER_FLAGS=--split-per-abi"
if /I "%APK_PACKAGING%"=="fat" set "APK_FLUTTER_FLAGS="
if /I "%APK_PACKAGING%"=="universal" set "APK_FLUTTER_FLAGS="

set "SETUP_QUIET=1"
call "%~dp0..\env\setup.bat"
if errorlevel 1 exit /b 1
REM Prefer scripts\env\config.bat JDK (avoid JDK 25 on system PATH breaking Kotlin)
if defined JAVA_HOME set "PATH=%JAVA_HOME%\bin;%PATH%"

cd /d "%~dp0..\.."
if not exist "pubspec.yaml" (
  echo [ERROR] pubspec.yaml not found: %CD%
  exit /b 1
)

set "APK_DIR=%CD%\build\app\outputs\flutter-apk"
if /I "%MODE%"=="debug" (
  set "APK_FILE=%APK_DIR%\app-debug.apk"
) else if /I "%APK_PACKAGING%"=="split" (
  set "APK_FILE=%APK_DIR%\app-arm64-v8a-release.apk"
) else (
  set "APK_FILE=%APK_DIR%\app-release.apk"
)

call :PrintBanner

echo [1/3] flutter pub get ...
echo.
call flutter pub get
if errorlevel 1 (
  echo [FAIL] flutter pub get
  exit /b 1
)
echo.
echo [OK] pub get done.
echo.

echo [2/3] Build settings
if defined API_BASE (
  echo       API baked in APK : %EFFECTIVE_API%
  echo       Source           : %API_SOURCE% + --dart-define
) else (
  echo       API baked in APK : %EFFECTIVE_API% ^(app default, no dart-define^)
  echo       Source           : lib/app/app_config.dart
  echo       Tip              : copy build-config.example.bat -^> build-config.bat
)
echo       WebSocket        : %EFFECTIVE_WS%
echo       After install    : App can override in Settings - API server
echo       APK packaging    : %APK_PACKAGING% ^(fat=universal 3 ABIs, largest^)
echo.

call :PrepareGradle
if errorlevel 1 exit /b 1

echo [3/3] flutter build apk --%MODE% ^(Gradle assemble, often 5-15 min^)
echo       First build or after upgrade may take longer. Do not close this window.
echo       Output folder    : %APK_DIR%
echo       JAVA_HOME        : %JAVA_HOME%
echo.

if /I "%MODE%"=="debug" (
  call flutter build apk --debug %APK_FLUTTER_FLAGS% %DART_DEFINES%
) else (
  call flutter build apk --release %APK_FLUTTER_FLAGS% %DART_DEFINES%
)
set "BUILD_ERR=%ERRORLEVEL%"
echo.

if not "%BUILD_ERR%"=="0" (
  echo [FAIL] flutter build exited %BUILD_ERR%
  echo.
  echo If you see Kotlin "incremental caches" / "Storage is already registered":
  echo   scripts\apk\build.bat clean %MODE%
  echo   Close Android Studio / other Gradle builds, then retry.
  exit /b %BUILD_ERR%
)

call :PrintSuccess
exit /b 0

:PrepareGradle
echo [2b] Gradle / Kotlin prep ...
if exist "%CD%\android\gradlew.bat" (
  pushd "%CD%\android"
  call gradlew.bat --stop >nul 2>&1
  popd
  echo       gradlew --stop done
) else (
  echo       gradlew.bat not found, skip --stop
)
if "%DO_CLEAN%"=="1" (
  echo       flutter clean ^(requested^) ...
  call flutter clean
  if errorlevel 1 exit /b 1
)
exit /b 0

:PrintBanner
echo.
echo ========================================
echo   Sara Live - APK build
echo ========================================
echo   Mode        : %MODE%
echo   Project     : %CD%
echo   API_BASE    : %EFFECTIVE_API%
echo   APK folder  : %APK_DIR%
echo   Packaging   : %APK_PACKAGING%
echo ========================================
echo.
exit /b 0

:PrintSuccess
echo ========================================
echo   BUILD OK
echo ========================================
echo   API in APK  : %EFFECTIVE_API%
echo   WS (derive) : %EFFECTIVE_WS%
echo.
if exist "%APK_FILE%" (
  echo   Main APK    : %APK_FILE%
  for %%A in ("%APK_FILE%") do echo   Size        : %%~zA bytes
) else (
  echo   Expected    : %APK_FILE%
  echo   ^(file missing — check errors above^)
)
if /I "%APK_PACKAGING%"=="split" (
  echo.
  echo   Install on phone: app-arm64-v8a-release.apk ^(most devices^)
)
if /I "%APK_PACKAGING%"=="arm64" (
  echo   ^(arm64-only; set APK_PACKAGING=fat in build-config for universal^)
)
echo.
echo   All APKs in folder:
dir /-C "%APK_DIR%\*.apk" 2>nul
if errorlevel 1 echo   ^(no .apk found^)
echo.
echo   Open folder : explorer "%APK_DIR%"
echo ========================================
echo.
if not defined BUILD_OPEN_FOLDER set "BUILD_OPEN_FOLDER=1"
if "%BUILD_OPEN_FOLDER%"=="1" if exist "%APK_DIR%" start "" explorer "%APK_DIR%"
exit /b 0
