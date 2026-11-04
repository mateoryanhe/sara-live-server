#Requires -Version 5.1
# 登录后由计划任务调用：写 1 小时 TCP 连接采样到 %USERPROFILE%\Logs\tcp-watch\
$ErrorActionPreference = 'Stop'
$logDir = Join-Path $env:USERPROFILE 'Logs\tcp-watch'
New-Item -ItemType Directory -Force -Path $logDir | Out-Null
$logFile = Join-Path $logDir ("tcp-{0:yyyy-MM-dd-HHmmss}.log" -f (Get-Date))
$scriptPath = Join-Path $PSScriptRoot 'watch-tcp-connections.ps1'
& $scriptPath -Loop -IntervalSec 60 -DurationMinutes 60 -Top 15 -LogFile $logFile
