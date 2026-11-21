@echo off
setlocal EnableDelayedExpansion
cd /d "%~dp0"
call "%~dp0config.bat"

echo Stopping Sara CMS Vite on port 5174 (and child node) ...

set FOUND=0
for /f "tokens=5" %%a in ('netstat -ano ^| findstr ":5174" ^| findstr LISTENING') do (
  set PID=%%a
  set FOUND=1
  echo Killing PID !PID! listening on 5174 ...
  taskkill /PID !PID! /T /F >nul 2>&1
)

if "%FOUND%"=="0" (
  echo No process listening on 5174.
) else (
  echo Done. Re-run start.bat for http://127.0.0.1:5174
)

endlocal
exit /b 0