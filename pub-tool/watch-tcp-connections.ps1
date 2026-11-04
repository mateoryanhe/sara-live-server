#Requires -Version 5.1
<#
.SYNOPSIS
  排查「上不了网 / TCP 4231 端口耗尽」：按进程统计 TCP 连接，可选循环写日志。

.EXAMPLE
  .\watch-tcp-connections.ps1
  .\watch-tcp-connections.ps1 -Top 20
  .\watch-tcp-connections.ps1 -Loop -IntervalSec 60 -DurationMinutes 60 -LogFile "$env:USERPROFILE\Logs\tcp-watch\session.log"
#>
param(
    [int] $Top = 15,
    [switch] $Loop,
    [int] $IntervalSec = 60,
    [int] $DurationMinutes = 0,
    [string] $LogFile = '',
    [int] $WarnTotalTcp = 12000
)

function Get-TcpConnectionSnapshot {
    $conns = Get-NetTCPConnection -ErrorAction SilentlyContinue |
        Where-Object { $_.State -in @('Established', 'TimeWait', 'CloseWait', 'SynSent', 'FinWait1', 'FinWait2', 'LastAck') }

    $byState = $conns | Group-Object State | Sort-Object Count -Descending |
        ForEach-Object { [PSCustomObject]@{ State = $_.Name; Count = $_.Count } }

    $byProc = $conns |
        Group-Object -Property OwningProcess |
        ForEach-Object {
            $proc = Get-Process -Id $_.Name -ErrorAction SilentlyContinue
            [PSCustomObject]@{
                PID     = [int]$_.Name
                Count   = $_.Count
                Name    = if ($proc) { $proc.ProcessName } else { '?' }
                Path    = if ($proc -and $proc.Path) { $proc.Path } else { '' }
            }
        } |
        Sort-Object Count -Descending

    $dynamic = netsh int ipv4 show dynamicport tcp 2>$null | Out-String
    $startPort = 49152
    $numPorts = 16384
    if ($dynamic -match 'Start Port\s*:\s*(\d+)') { $startPort = [int]$Matches[1] }
    if ($dynamic -match 'Number of Ports\s*:\s*(\d+)') { $numPorts = [int]$Matches[1] }

    [PSCustomObject]@{
        Time       = Get-Date -Format 'yyyy-MM-dd HH:mm:ss'
        TotalTcp   = $conns.Count
        ByState    = $byState
        ByProcess  = $byProc
        PortStart  = $startPort
        PortCount  = $numPorts
        WarnTotal  = $WarnTotalTcp
    }
}

function Format-SnapshotText {
    param($Snap)
    $lines = New-Object System.Collections.Generic.List[string]
    $lines.Add("=== $($Snap.Time) ===")
    $lines.Add("TCP sockets (tracked states): $($Snap.TotalTcp)  |  dynamic ports: $($Snap.PortStart)-$($Snap.PortStart + $Snap.PortCount - 1) ($($Snap.PortCount) ports)")
    if ($Snap.TotalTcp -ge $Snap.WarnTotal) {
        $lines.Add("*** WARNING: high socket count (approaching port exhaustion). Event: Tcpip 4231 ***")
    }
    if ($Snap.ByState) {
        $stateLine = ($Snap.ByState | ForEach-Object { "$($_.State)=$($_.Count)" }) -join ', '
        $lines.Add("By state: $stateLine")
    }
    $lines.Add("Top processes:")
    $rank = 0
    foreach ($p in $Snap.ByProcess | Select-Object -First $Top) {
        $rank++
        $path = if ($p.Path) { "  $($p.Path)" } else { '' }
        $lines.Add(("  {0,2}. PID {1,6}  {2,6} conn  {3}{4}" -f $rank, $p.PID, $p.Count, $p.Name, $path))
    }
    $lines.Add('')
    return ($lines -join "`r`n")
}

function Show-Snapshot {
    param($Snap)
    $text = Format-SnapshotText $Snap
    Write-Host $text
    if ($LogFile) {
        Add-Content -LiteralPath $LogFile -Value $text -Encoding UTF8
    }
}

$script:LoopStarted = $null
do {
    $snap = Get-TcpConnectionSnapshot
    Show-Snapshot $snap
    if (-not $Loop) { break }
    if ($DurationMinutes -gt 0 -and $script:LoopStarted -eq $null) {
        $script:LoopStarted = Get-Date
    }
    if ($DurationMinutes -gt 0) {
        $elapsed = (Get-Date) - $script:LoopStarted
        if ($elapsed.TotalMinutes -ge $DurationMinutes) {
            $tail = "=== stopped after $DurationMinutes min ==="
            Write-Host $tail
            if ($LogFile) { Add-Content -LiteralPath $LogFile -Value $tail -Encoding UTF8 }
            break
        }
    }
    Start-Sleep -Seconds $IntervalSec
} while ($true)
