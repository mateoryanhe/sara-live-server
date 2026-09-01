# Create Go + CMS deploy shortcuts, then prompt drag-to-taskbar (Win11 reliable)
$ErrorActionPreference = 'Stop'
$root = Split-Path -Parent $MyInvocation.MyCommand.Path
$repoRoot = Split-Path -Parent $root
$assetCandidates = @(
    (Join-Path $repoRoot 'assets')
    'C:\Users\hw\.cursor\projects\d-company-code-sara-live-server\assets'
)
$assets = $assetCandidates | Where-Object { Test-Path -LiteralPath $_ } | Select-Object -First 1

function Ensure-Icon {
    param(
        [string]$Dir,
        [string]$BaseName,
        [string]$AssetPngName
    )
    $png = Join-Path $Dir "$BaseName.png"
    $ico = Join-Path $Dir "$BaseName.ico"
    if ($assets) {
        $assetPng = Join-Path $assets $AssetPngName
        if (Test-Path -LiteralPath $assetPng) {
            Copy-Item -LiteralPath $assetPng -Destination $png -Force
        }
    }
    if ((Test-Path -LiteralPath $png) -and -not (Test-Path -LiteralPath $ico)) {
        & powershell -NoProfile -ExecutionPolicy Bypass -File (Join-Path $Dir 'convert-icon.ps1') `
            -PngPath $png -IcoPath $ico | Out-Null
    }
    if (-not (Test-Path -LiteralPath $ico)) {
        throw "Icon not found: $ico"
    }
    return $ico
}

function New-DeployShortcut {
    param(
        [string]$DeployScript,
        [string]$ShortcutName,
        [string]$LauncherExeName,
        [string]$Description,
        [string]$IconPath = ''
    )
    $args = @(
        '-NoProfile', '-ExecutionPolicy', 'Bypass',
        '-File', (Join-Path $root 'pin-to-taskbar.ps1'),
        '-DeployScript', $DeployScript,
        '-ShortcutName', $ShortcutName,
        '-LauncherExeName', $LauncherExeName,
        '-Description', $Description,
        '-NoPrompt'
    )
    if ($IconPath -ne '') {
        $args += @('-IconPath', $IconPath)
    }
    return (& powershell @args | Select-Object -Last 1)
}

Write-Host '========================================'
Write-Host ' XR Live - Deploy shortcuts (drag pin)'
Write-Host '========================================'
Write-Host ''

$goDir = Join-Path $root 'go-test'
$cmsDir = Join-Path $root 'cms-test'

$goIco = Ensure-Icon -Dir $goDir -BaseName 'deploy-go' -AssetPngName 'go-test-deploy-icon.png'
$cmsIco = Ensure-Icon -Dir $cmsDir -BaseName 'deploy-cms' -AssetPngName 'cms-test-deploy-icon.png'

$goScript = Join-Path $goDir 'deploy.bat'
if (-not (Test-Path -LiteralPath $goScript)) {
    throw "Go deploy script not found: $goScript"
}

Write-Host '[1/2] Go test deploy'
New-DeployShortcut -DeployScript $goScript -ShortcutName 'XR-Go-Test-Deploy' -LauncherExeName 'XR-Go-Deploy.exe' `
    -Description 'XR Live Go test server deploy' -IconPath $goIco | Out-Null
Write-Host ''

Write-Host '[2/2] CMS test deploy'
New-DeployShortcut -DeployScript (Join-Path $cmsDir 'upload.bat') -ShortcutName 'XR-CMS-Test-Deploy' -LauncherExeName 'XR-CMS-Deploy.exe' `
    -Description 'XR Live CMS test server build and upload' -IconPath $cmsIco | Out-Null
Write-Host ''

$shortcuts = @(
    (Join-Path $env:USERPROFILE 'Desktop\XR-Go-Test-Deploy.lnk')
    (Join-Path $env:USERPROFILE 'Desktop\XR-CMS-Test-Deploy.lnk')
) | Where-Object { Test-Path -LiteralPath $_ }
$names = ($shortcuts | ForEach-Object { [System.IO.Path]::GetFileNameWithoutExtension($_) }) -join "`n  - "
$text = @"
Win11 请手动拖到任务栏（最可靠）：

  1. 资源管理器将打开桌面快捷方式
  2. 依次把两个快捷方式拖到任务栏
  3. 松开即可固定

快捷方式：
  - $names
"@
Add-Type -AssemblyName System.Windows.Forms
[void][System.Windows.Forms.MessageBox]::Show($text, '固定到任务栏', 'OK', 'Information')

foreach ($path in $shortcuts) {
    Start-Process explorer.exe -ArgumentList "/select,`"$path`""
}

Write-Host 'Done. Drag both desktop shortcuts onto the taskbar.'
