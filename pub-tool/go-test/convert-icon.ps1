param(
    [string]$PngPath = (Join-Path $PSScriptRoot 'deploy-go.png'),
    [string]$IcoPath = (Join-Path $PSScriptRoot 'deploy-go.ico')
)

Add-Type -AssemblyName System.Drawing
$bmp = New-Object System.Drawing.Bitmap($PngPath)
$h = $bmp.GetHicon()
$icon = [System.Drawing.Icon]::FromHandle($h)
$fs = New-Object System.IO.FileStream($IcoPath, [System.IO.FileMode]::Create)
$icon.Save($fs)
$fs.Close()
$bmp.Dispose()
Write-Host "Created: $IcoPath"
