@echo off
chcp 65001 >nul 2>&1
cd /d "%~dp0"
call "%~dp0dev-usb.bat"
set "EC=%ERRORLEVEL%"
echo.
if not "%EC%"=="0" (
  echo ========================================
  echo   DEV USB FAILED - see errors above
  echo ========================================
)
pause
exit /b %EC%
