@echo off



chcp 65001 >nul 2>&1



REM Double-click: install latest release APK to USB phone



cd /d "%~dp0"



call "%~dp0install.bat" release



echo.



pause



