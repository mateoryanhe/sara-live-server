@echo off
REM 双击：打 debug APK
cd /d "%~dp0"
call "%~dp0build.bat" debug
echo.
pause
