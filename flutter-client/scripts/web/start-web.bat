@echo off
echo.
echo Dual Web debug (1v1 call):
echo   scripts\web\start-web-audience.bat  - Chrome 8080, caller
echo   scripts\web\start-web-anchor.bat    - Chrome 8081, callee
echo.
echo Default: audience.
echo.
call "%~dp0start-web-audience.bat"
