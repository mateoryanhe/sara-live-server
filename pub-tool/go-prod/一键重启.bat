@echo off
cd /d "%~dp0"
call restart.bat
set RESTART_EXIT_CODE=%errorlevel%
echo.
if %RESTART_EXIT_CODE% equ 0 (
    echo ===== production restart succeeded =====
) else (
    echo ===== production restart failed, exit code: %RESTART_EXIT_CODE% =====
)
pause
exit /b %RESTART_EXIT_CODE%
