@echo off
REM ============================================================
REM  Copy to:  start-web-config.bat  (same folder as this file)
REM  Then edit below. Re-run start-web-audience / anchor.
REM ============================================================

REM Backend API root (no trailing slash). Baked into this Web run via --dart-define.
REM Local Go: http://127.0.0.1:9443   Test: https://www.bigtktool.shop
set "API_BASE=https://www.bigtktool.shop"

REM One dir for both roles (optional):
REM set "HTTP_LOG_DIR=%~dp0..\..\logs\web-shared"

REM Separate dirs (recommended for 1v1 test):
set "HTTP_LOG_DIR_AUDIENCE=%~dp0..\..\logs\web-audience"
set "HTTP_LOG_DIR_ANCHOR=%~dp0..\..\logs\web-anchor"

set "HTTP_LOG_ENABLED=true"

REM Open log folder in Explorer when launcher starts (1=yes 0=no)
set "HTTP_LOG_OPEN_FOLDER=1"
