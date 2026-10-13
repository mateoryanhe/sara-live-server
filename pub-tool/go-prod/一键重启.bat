@echo off
REM Double-click: outer cmd /k keeps window open to read logs
if /i "%~1"=="--nopause" goto :main
if /i not "%~1"=="__inner__" (
    cmd.exe /k call "%~f0" __inner__ %*
    exit /b 0
)
shift /1

:main
setlocal EnableExtensions EnableDelayedExpansion
cd /d "%~dp0"
if not exist "%~dp0restart.bat" (
    echo ERROR: restart.bat not found in %~dp0
    goto :failed
)

if /i not "%~1"=="--nopause" set "XR_RESTART_NOPAUSE=1"

call "%~dp0restart.bat"
set "RESTART_EXIT_CODE=!ERRORLEVEL!"
echo.
if !RESTART_EXIT_CODE! equ 0 (
    echo ===== production restart succeeded =====
    endlocal
    if /i "%~1"=="--nopause" exit /b 0
    echo.
    echo Type exit to close this window.
    exit /b 0
)

:failed
if not defined RESTART_EXIT_CODE set "RESTART_EXIT_CODE=1"
echo ===== production restart failed, exit code: !RESTART_EXIT_CODE! =====
echo.
echo Scroll up to read the log. Type exit to close this window.
endlocal
if /i "%~1"=="--nopause" exit /b 1
exit /b 1
