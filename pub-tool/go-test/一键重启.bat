@echo off
cd /d "%~dp0"
echo XR Game Server Test — hot restart only ^(no build/upload^)
echo.
call restart.bat
echo.
echo ===== finished — close this window when done =====
cmd /k
