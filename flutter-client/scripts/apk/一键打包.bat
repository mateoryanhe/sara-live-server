@echo off
chcp 65001 >nul 2>&1
REM Double-click: release APK (see build.bat for API_BASE / progress)
cd /d "%~dp0"
call "%~dp0build.bat" release
echo.
pause
