@echo off
cd /d "%~dp0"
call "%~dp0config.bat"
plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T -hostkey "%SSH_HOST_KEY%" %REMOTE_USER%@%REMOTE_HOST% "pkill -x %HR_GW_NAME% 2>/dev/null; sleep 1; cd %REMOTE_DIR% && nohup ./%HR_GW_NAME% >> %REMOTE_DIR%/live-hot-restart-gw.log 2>&1 & sleep 2; pgrep -af %HR_GW_NAME%; ss -tlnp 2>/dev/null | grep 18900 || true"
