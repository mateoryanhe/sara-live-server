# Create deploy launcher + desktop shortcut, then open Explorer for manual drag-to-taskbar (Win11 reliable)
param(
    [Parameter(Mandatory = $true)]
    [string]$DeployScript,
    [Parameter(Mandatory = $true)]
    [string]$ShortcutName,
    [Parameter(Mandatory = $true)]
    [string]$Description,
    [Parameter(Mandatory = $true)]
    [string]$LauncherExeName,
    [string]$IconPath = '',
    [switch]$NoPrompt
)

$ErrorActionPreference = 'Stop'

function Build-DeployLauncher {
    param(
        [string]$DeployDir,
        [string]$ScriptName,
        [string]$ExeName
    )
    $launcherSrc = Join-Path $PSScriptRoot 'taskbar-launcher\main.go'
    if (-not (Test-Path -LiteralPath $launcherSrc)) {
        throw "Launcher source not found: $launcherSrc"
    }
    $exePath = Join-Path $DeployDir $ExeName
    $buildArgs = @(
        'build', '-o', $exePath,
        '-ldflags', "-X main.scriptName=$ScriptName",
        $launcherSrc
    )
    & go @buildArgs
    if ($LASTEXITCODE -ne 0) {
        throw "go build failed for $ExeName"
    }
    if (-not (Test-Path -LiteralPath $exePath)) {
        throw "Launcher exe not created: $exePath"
    }
    return $exePath
}

function New-ExeShortcut {
    param(
        [string]$LinkPath,
        [string]$ExePath,
        [string]$WorkDir,
        [string]$Description,
        [string]$IconLocation
    )
    $WshShell = New-Object -ComObject WScript.Shell
    $sc = $WshShell.CreateShortcut($LinkPath)
    $sc.TargetPath = $ExePath
    $sc.Arguments = ''
    $sc.WorkingDirectory = $WorkDir
    $sc.WindowStyle = 1
    $sc.Description = $Description
    $sc.IconLocation = $IconLocation
    $sc.Save()
    return $LinkPath
}

function Show-DragPinPrompt {
    param(
        [string[]]$ShortcutPaths,
        [string]$Title = '固定到任务栏'
    )
    $names = ($ShortcutPaths | ForEach-Object { [System.IO.Path]::GetFileNameWithoutExtension($_) }) -join "`n  - "
    $text = @"
Win11 请手动拖到任务栏（最可靠）：

  1. 资源管理器已打开桌面快捷方式
  2. 用鼠标把快捷方式拖到任务栏
  3. 松开即可固定

快捷方式：
  - $names
"@
    Add-Type -AssemblyName System.Windows.Forms
    [void][System.Windows.Forms.MessageBox]::Show($text, $Title, 'OK', 'Information')

    foreach ($path in $ShortcutPaths) {
        if (Test-Path -LiteralPath $path) {
            Start-Process explorer.exe -ArgumentList "/select,`"$path`""
        }
    }
}

$DeployScript = $ExecutionContext.SessionState.Path.GetUnresolvedProviderPathFromPSPath($DeployScript)
$deployDir = Split-Path -Parent $DeployScript
$scriptLeaf = Split-Path -Leaf $DeployScript
if (-not (Test-Path -LiteralPath $DeployScript)) {
    throw "Not found: $DeployScript"
}

$baseName = [System.IO.Path]::GetFileNameWithoutExtension($ShortcutName)
if ($ShortcutName.EndsWith('.lnk')) {
    $baseName = [System.IO.Path]::GetFileNameWithoutExtension($ShortcutName)
}

$exePath = Build-DeployLauncher -DeployDir $deployDir -ScriptName $scriptLeaf -ExeName $LauncherExeName

$iconLocation = "$exePath,0"
if ($IconPath -ne '') {
    $IconPath = $ExecutionContext.SessionState.Path.GetUnresolvedProviderPathFromPSPath($IconPath)
    if (Test-Path -LiteralPath $IconPath) {
        $iconLocation = $IconPath
    }
}

$desktopLnk = Join-Path $env:USERPROFILE "Desktop\$baseName.lnk"
$startMenuDir = Join-Path $env:APPDATA 'Microsoft\Windows\Start Menu\Programs\XR Live'
$startMenuLnk = Join-Path $startMenuDir "$baseName.lnk"
New-Item -ItemType Directory -Force -Path $startMenuDir | Out-Null

New-ExeShortcut -LinkPath $desktopLnk -ExePath $exePath -WorkDir $deployDir -Description $Description -IconLocation $iconLocation
New-ExeShortcut -LinkPath $startMenuLnk -ExePath $exePath -WorkDir $deployDir -Description $Description -IconLocation $iconLocation

Write-Host "Launcher exe: $exePath"
Write-Host "Deploy script: $DeployScript"
Write-Host "Icon: $iconLocation"
Write-Host "Desktop: $desktopLnk"
Write-Host "Start Menu: $startMenuLnk"

if (-not $NoPrompt) {
    Show-DragPinPrompt -ShortcutPaths @($desktopLnk)
    Write-Host 'Drag the desktop shortcut onto the taskbar to pin.'
} else {
    Write-Host 'Shortcut ready (prompt skipped).'
}

return $desktopLnk
