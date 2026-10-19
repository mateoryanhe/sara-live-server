@echo off
REM Copy to build-config.bat and edit (optional).

REM API (optional):
REM set "API_BASE=http://192.168.2.11:9443"

REM APK size — default in build.bat is arm64 (~1/3 of fat universal):
REM   arm64     one APK for most phones (recommended)
REM   split     app-armeabi-v7a / arm64-v8a / x86_64 — three files
REM   fat       universal single APK (largest, all ABIs)
REM set "APK_PACKAGING=arm64"
