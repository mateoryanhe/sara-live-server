param(
    [string]$PngPath = (Join-Path $PSScriptRoot 'deploy-cms.png'),
    [string]$IcoPath = (Join-Path $PSScriptRoot 'deploy-cms.ico')
)

Add-Type -AssemblyName System.Drawing
$bmp = New-Object System.Drawing.Bitmap($PngPath)
$h = $bmp.GetHicon()
$icon = [System.Drawing.Icon]::FromHandle($h)
$fs = New-Object System.IO.FileStream($IcoPath, [IO.FileMode]::Create)
$icon.Save($fs)
$fs.Close()
$bmp.Dispose()
Write-Host "Created: $IcoPath"
