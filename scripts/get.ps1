#!/usr/bin/env pwsh

$ErrorActionPreference = 'Stop'

$BinPath = "${Home}\.holesail-proxy\bin"
$Zip = "$RuntimerPath\holesail-proxy.zip"
$Exe = "$RuntimerPath\holesail-proxy.exe"
$OldExe = "$env:Temp\holesail-proxy-old.exe"

$Target = "windows-amd64"

$DownloadUrl = "https://github.com/CyberL1/holesail-proxy/releases/latest/download/holesal-proxy-${Target}.zip"

if (!(Test-Path $BinPath)) {
  New-Item $BinPath -ItemType Directory | Out-Null
}

curl.exe -Lo $Zip $DownloadUrl

if (Test-Path $Exe) {
  Move-Item -Path $Exe -Destination $OldExe -Force
}

Expand-Archive -LiteralPath $Zip -DestinationPath $BinPath -Force
Remove-Item $Zip

$User = [System.EnvironmentVariableTarget]::User
$Path = [System.Environment]::GetEnvironmentVariable('Path', $User)

if (!(";${Path};".ToLower() -like "*;${BinPath};*".ToLower())) {
  [System.Environment]::SetEnvironmentVariable('Path', "${Path};${BinPath}", $User)
  $Env:Path += ";${BinPath}"
}

Write-Output "Holesail proxy was installed to $Exe"
Write-Output "Run 'holesail-proxy up' to get started"
