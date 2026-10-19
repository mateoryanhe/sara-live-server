@echo off

chcp 65001 >nul 2>&1

REM Double-click: install app-debug.apk to USB phone

cd /d "%~dp0"

call "%~dp0install.bat" debug

echo.

pause

