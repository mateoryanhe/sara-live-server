@echo off
cd /d "%~dp0"
call "%~dp0restart.bat"
set "RESTART_EXIT_CODE=%errorlevel%"
echo.
if "%RESTART_EXIT_CODE%"=="0" (
    echo ===== review restart succeeded =====
) else (
    echo ===== review restart failed, exit code: %RESTART_EXIT_CODE% =====
)
if /i not "%~1"=="--nopause" pause
exit /b %RESTART_EXIT_CODE%
