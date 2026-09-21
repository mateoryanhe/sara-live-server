@echo off
cd /d "%~dp0"
call "%~dp0restart.bat"
set "RESTART_EXIT_CODE=%errorlevel%"
echo.
if "%RESTART_EXIT_CODE%"=="0" (
    echo ===== production restart succeeded =====
) else (
    echo ===== production restart failed, exit code: %RESTART_EXIT_CODE% =====
)
pause
exit /b %RESTART_EXIT_CODE%
