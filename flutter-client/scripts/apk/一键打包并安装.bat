@echo off

chcp 65001 >nul 2>&1

REM Double-click: build release then adb install

cd /d "%~dp0"

call "%~dp0build.bat" release

if errorlevel 1 (

  echo.

  pause

  exit /b 1

)

call "%~dp0install.bat" release

echo.

pause

