# 创建 Go/CMS 测试服部署快捷方式，并固定到 Win11 任务栏
# 注意: .bat 文件本身不能「固定到任务栏」，必须先创建 .lnk 快捷方式
# 用法: 双击 pub-tool/pin-deploy-taskbar.bat

$ErrorActionPreference = 'Stop'

$root = Split-Path -Parent $MyInvocation.MyCommand.Path
$cmdExe = Join-Path $env:WINDIR 'System32\cmd.exe'
$taskbarDir = Join-Path $env:APPDATA 'Microsoft\Internet Explorer\Quick Launch\User Pinned\TaskBar'

$shortcuts = @(
    @{
        Name = 'XR-Live-Test-Deploy'
        Bat  = Join-Path $root 'deploy-go-test.bat'
        Desc = 'Go test server deploy'
    },
    @{
        Name = 'XR-Live-CMS-Deploy'
        Bat  = Join-Path $root 'deploy-cms-test.bat'
        Desc = 'CMS test server build and upload'
    }
)

function New-DeployShortcut {
    param(
        [string]$Path,
        [string]$BatPath,
        [string]$WorkDir,
        [string]$Description
    )
    if (-not (Test-Path -LiteralPath $BatPath)) {
        throw "Launcher not found: $BatPath"
    }
    $WshShell = New-Object -ComObject WScript.Shell
    $sc = $WshShell.CreateShortcut($Path)
    # Win11 只能固定 .lnk；目标设为 cmd.exe 才出现「固定到任务栏」
    $sc.TargetPath = $cmdExe
    $sc.Arguments = '/k ""' + $BatPath + '""'
    $sc.WorkingDirectory = $WorkDir
    $sc.WindowStyle = 1
    $sc.Description = $Description
    $sc.IconLocation = "$cmdExe,0"
    $sc.Save()
}

function Pin-ToTaskbar {
    param([string]$ShortcutPath)
    if (-not (Test-Path -LiteralPath $ShortcutPath)) {
        throw "Shortcut not found: $ShortcutPath"
    }
    New-Item -ItemType Directory -Force -Path $taskbarDir | Out-Null
    $dest = Join-Path $taskbarDir (Split-Path $ShortcutPath -Leaf)
    Copy-Item -LiteralPath $ShortcutPath -Destination $dest -Force

    try {
        $shell = New-Object -ComObject Shell.Application
        $folder = $shell.Namespace((Split-Path -Parent $ShortcutPath))
        $item = $folder.ParseName((Split-Path $ShortcutPath -Leaf))
        $item.InvokeVerb('taskbarpin')
        return $true
    } catch {
        Write-Host "taskbarpin failed: $($_.Exception.Message)"
        return $false
    }
}

Write-Host '========================================'
Write-Host '  XR Live - Pin deploy shortcuts'
Write-Host '========================================'
Write-Host "Project: $root"
Write-Host ''
Write-Host 'Tip: Do NOT pin .bat files. Pin the .lnk shortcuts on Desktop.'
Write-Host ''

$created = @()
foreach ($s in $shortcuts) {
    $desktopLnk = Join-Path $env:USERPROFILE "Desktop\$($s.Name).lnk"
    Write-Host "Creating shortcut: $desktopLnk"
    New-DeployShortcut -Path $desktopLnk -BatPath $s.Bat -WorkDir $root -Description $s.Desc
    Write-Host "Pinning to taskbar: $($s.Name)"
    $pinned = Pin-ToTaskbar -ShortcutPath $desktopLnk
    if ($pinned) {
        Write-Host "OK (taskbar): $($s.Name)"
    } else {
        Write-Host "Shortcut created. Manually: drag .lnk to taskbar, or right-click .lnk -> Pin to taskbar"
    }
    $created += $desktopLnk
    Write-Host ''
}

Write-Host 'Desktop shortcuts:'
foreach ($p in $created) {
    Write-Host "  $p"
}
Write-Host ''
Write-Host 'If taskbar still empty: drag XR-Live-*.lnk from Desktop onto the taskbar.'
Read-Host 'Press Enter to close'
