@echo off
cd /d "%~dp0"
call "%~dp0create-docker.bat"
exit /b %ERRORLEVEL%