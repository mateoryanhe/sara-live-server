@echo off
cd /d "%~dp0"
call "%~dp0config.bat"
plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T -hostkey "%SSH_HOST_KEY%" %REMOTE_USER%@%REMOTE_HOST% "ls -la --time-style=long-iso %REMOTE_DIR%/xr-game-server %REMOTE_DIR%/live-hot-restart-gw %REMOTE_DIR%/config.json 2>/dev/null; echo ---; pgrep -xo xr-game-server; ps -eo pid,lstart,etime,cmd | grep xr-game-server | grep -v grep"
