@echo off
cd /d "%~dp0"
REM ASCII-only entry for taskbar launcher (avoid Chinese filename encoding issues)
call "%~dp0deploy.bat"
