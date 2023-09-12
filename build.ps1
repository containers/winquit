$ErrorActionPreference = 'Stop'

function runCommand {
    param (
        [string] $command
    )

    Write-Host $command
    Invoke-Expression $command   
}

Remove-Item -Force -Recurse -Path .\bin -ErrorAction Ignore 
if (($args.Count -gt 0) -and ($args[0] -eq "clean")) {
    Exit 0
}

$env:GOOS="windows"
$env:GOARCH="amd64"
 
runCommand "go build -v -o bin/winquit.exe ./cmd/winquit"

if (($args.Count -gt 0) -and ($args[0] -eq "test")) {  
    runCommand "go test -v ./test"
}