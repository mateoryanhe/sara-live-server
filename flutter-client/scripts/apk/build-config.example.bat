@echo off
REM Copy to build-config.bat and edit (optional).

REM API (optional):
REM set "API_BASE=http://192.168.2.11:9443"

REM APK size — default in build.bat is arm64 (~1/3 of fat universal):
REM   arm64     one APK for most phones (recommended)
REM   split     app-armeabi-v7a / arm64-v8a / x86_64 — three files
REM   fat       universal single APK (largest, all ABIs)
REM set "APK_PACKAGING=arm64"

REM Build speed (optional):
REM   BUILD_FORCE_PUB=1     always flutter pub get (default: skip if pubspec unchanged)
REM   BUILD_GRADLE_STOP=1   gradlew --stop before each build (slower; only if Gradle stuck)
REM   BUILD_OPEN_FOLDER=0   do not open explorer after build
