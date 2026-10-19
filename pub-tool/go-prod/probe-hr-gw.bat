@echo off
cd /d "%~dp0"
call "%~dp0config.bat"
echo === SSH localhost :18900/health ===
REM curl -w 的 %% 留给远端；单引号 URL 避免 bash 解析 ?auth=
plink.exe -ssh -i "%SSH_KEY_PATH%" -P %REMOTE_PORT% -batch -T -hostkey "%SSH_HOST_KEY%" %REMOTE_USER%@%REMOTE_HOST% "curl -sS -m 5 http://127.0.0.1:18900/health; echo; http=$(curl -sS -m 5 -o /dev/null -w '%%{http_code}' 'http://127.0.0.1:18900/hotRestart?auth=bad'); echo hotRestart_bad_auth_http=$http; ss -tlnp 2>/dev/null | grep 18900 || netstat -tlnp 2>/dev/null | grep 18900; pgrep -af %HR_GW_NAME% || true"
