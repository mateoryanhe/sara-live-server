@echo off
setlocal EnableDelayedExpansion
cd /d "%~dp0"

call "%~dp0config.bat"

echo ========================================
echo   CMS local start (Vite)
echo ========================================
echo CMS dir : %CMS_DIR%
echo Dev URL : %CMS_DEV_URL%
echo API     : %CMS_API_HINT%  (.env.development)
echo.

if not exist "%CMS_DIR%\package.json" (
  echo ERROR: cms package.json not found: %CMS_DIR%
  exit /b 1
)

where node >nul 2>&1
if errorlevel 1 (
  echo ERROR: node not found in PATH. Install Node.js 20+ then retry.
  exit /b 1
)

where npm >nul 2>&1
if errorlevel 1 (
  echo ERROR: npm not found in PATH.
  exit /b 1
)

cd /d "%CMS_DIR%"
if errorlevel 1 (
  echo ERROR: cannot cd to %CMS_DIR%
  exit /b 1
)

if not exist "node_modules\" (
  echo node_modules missing, running npm install ...
  call npm install
  if errorlevel 1 (
    echo ERROR: npm install failed
    exit /b 1
  )
  echo.
)

echo Starting: npm run dev
echo Tip: keep local Go API on %CMS_API_HINT% (config/local address 9443)
echo.

call npm run dev
exit /b %ERRORLEVEL%
