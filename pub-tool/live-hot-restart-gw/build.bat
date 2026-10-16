@echo off
setlocal
cd /d "%~dp0"
go build -o live-hot-restart-gw.exe .
if errorlevel 1 exit /b 1
echo OK: Windows %~dp0live-hot-restart-gw.exe
echo Linux prod: set GOOS=linux&& set GOARCH=amd64&& go build -o live-hot-restart-gw .
exit /b 0
