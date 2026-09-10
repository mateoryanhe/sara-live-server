@echo off
REM 双击：打 release APK
cd /d "%~dp0"
call "%~dp0build.bat" release
echo.
pause
