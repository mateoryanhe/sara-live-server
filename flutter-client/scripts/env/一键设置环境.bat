@echo off
REM 双击：按 config.bat 写入当前用户永久环境变量
cd /d "%~dp0"
call "%~dp0setup.bat" install
echo.
pause
