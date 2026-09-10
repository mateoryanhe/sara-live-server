# 由 setup.bat install 调用：读取同目录 config.bat，写入当前用户永久环境变量
param(
  [switch]$Yes
)
$ErrorActionPreference = 'Stop'
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$configBat = Join-Path $scriptDir 'config.bat'
if (-not (Test-Path -LiteralPath $configBat)) {
  throw "缺少 config.bat: $configBat"
}

$cfg = @{}
Get-Content -LiteralPath $configBat | ForEach-Object {
  if ($_ -match '^\s*set\s+"([^=]+)=([^"]*)"\s*$') {
    $cfg[$Matches[1]] = $Matches[2]
  } elseif ($_ -match '^\s*set\s+([^=]+)=(.*)\s*$') {
    $cfg[$Matches[1]] = $Matches[2].Trim()
  }
}

$flutterHome = $cfg['FLUTTER_HOME']
$javaHome = $cfg['JAVA_HOME']
$androidSdk = $cfg['ANDROID_SDK']
$useCn = $cfg['USE_CN_MIRROR']

function Assert-Dir([string]$Path, [string]$Label) {
  if (-not (Test-Path -LiteralPath $Path)) {
    throw "$Label 不存在: $Path （请先编辑 config.bat）"
  }
}
Assert-Dir (Join-Path $flutterHome 'bin') 'Flutter'
Assert-Dir (Join-Path $javaHome 'bin') 'JDK'
Assert-Dir $androidSdk 'Android SDK'

$prepend = @(
  (Join-Path $flutterHome 'bin'),
  (Join-Path $androidSdk 'platform-tools'),
  (Join-Path $androidSdk 'cmdline-tools\latest\bin'),
  (Join-Path $javaHome 'bin')
)

Write-Host "即将写入【当前用户】永久环境变量:"
Write-Host "  FLUTTER_HOME = $flutterHome"
Write-Host "  JAVA_HOME    = $javaHome"
Write-Host "  ANDROID_HOME = $androidSdk"
Write-Host "  PATH 前置    = $($prepend -join ' ; ')"
if ($useCn -eq '1') { Write-Host "  中国镜像     = 开启" }

if (-not $Yes) {
  $confirm = Read-Host "确认写入? (Y/N)"
  if ($confirm -notmatch '^[Yy]$') {
    Write-Host "已取消"
    exit 0
  }
}

[Environment]::SetEnvironmentVariable('FLUTTER_HOME', $flutterHome, 'User')
[Environment]::SetEnvironmentVariable('JAVA_HOME', $javaHome, 'User')
[Environment]::SetEnvironmentVariable('ANDROID_HOME', $androidSdk, 'User')
[Environment]::SetEnvironmentVariable('ANDROID_SDK_ROOT', $androidSdk, 'User')
if ($useCn -eq '1') {
  [Environment]::SetEnvironmentVariable('PUB_HOSTED_URL', 'https://pub.flutter-io.cn', 'User')
  [Environment]::SetEnvironmentVariable('FLUTTER_STORAGE_BASE_URL', 'https://storage.flutter-io.cn', 'User')
}

$userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
if ([string]::IsNullOrWhiteSpace($userPath)) { $userPath = '' }
$parts = @($userPath -split ';' | Where-Object { $_ -and $_.Trim() })
$norm = $prepend | ForEach-Object { $_.TrimEnd('\') }
$parts = $parts | Where-Object { $norm -notcontains $_.TrimEnd('\') }
$newPath = ($norm + $parts) -join ';'
[Environment]::SetEnvironmentVariable('Path', $newPath, 'User')

Write-Host ""
Write-Host "[OK] 用户环境变量已写入。请关闭并重新打开终端 / IDE 后再跑 flutter doctor。"
