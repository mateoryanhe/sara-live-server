param(
    [Parameter(Mandatory = $false)]
    [string]$CountriesGo = "",

    [Parameter(Mandatory = $false)]
    [string]$OutDir = "",

    [Parameter(Mandatory = $false)]
    [string]$CdnBase = "https://flagcdn.com/w80"
)

$ErrorActionPreference = "Stop"
$here = Split-Path -Parent $MyInvocation.MyCommand.Path

if ([string]::IsNullOrWhiteSpace($CountriesGo)) {
    $CountriesGo = Join-Path $here "..\..\go-src\constants\country\countries.go"
}
if ([string]::IsNullOrWhiteSpace($OutDir)) {
    $OutDir = Join-Path $here "flags"
}

$CountriesGo = [System.IO.Path]::GetFullPath($CountriesGo)
$OutDir = [System.IO.Path]::GetFullPath($OutDir)

if (-not (Test-Path -LiteralPath $CountriesGo)) {
    Write-Error "countries.go not found: $CountriesGo"
}

if (-not (Test-Path -LiteralPath $OutDir)) {
    New-Item -ItemType Directory -Path $OutDir | Out-Null
}

$text = Get-Content -LiteralPath $CountriesGo -Raw -Encoding UTF8
$matches = [regex]::Matches($text, 'newCountry\(\s*"([A-Za-z]{2})"')
$codes = New-Object System.Collections.Generic.List[string]
$seen = @{}
foreach ($m in $matches) {
    $code = $m.Groups[1].Value.ToLowerInvariant()
    if (-not $seen.ContainsKey($code)) {
        $seen[$code] = $true
        [void]$codes.Add($code)
    }
}

if ($codes.Count -eq 0) {
    Write-Error "no country codes parsed from $CountriesGo"
}

Write-Host "codes=$($codes.Count) out=$OutDir cdn=$CdnBase"

$ok = 0
$fail = 0
foreach ($code in $codes) {
    $url = "$CdnBase/$code.png"
    $dest = Join-Path $OutDir "$code.png"
    try {
        Invoke-WebRequest -Uri $url -OutFile $dest -UseBasicParsing
        $ok++
        Write-Host "OK  $code"
    }
    catch {
        $fail++
        Write-Warning "FAIL $code  $url  $($_.Exception.Message)"
        if (Test-Path -LiteralPath $dest) {
            Remove-Item -LiteralPath $dest -Force
        }
    }
}

Write-Host ""
Write-Host "done ok=$ok fail=$fail -> $OutDir"
if ($fail -gt 0) {
    exit 1
}
