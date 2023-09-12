$ErrorActionPreference = 'Stop'

function runCommand {
    param (
        [string] $command
    )

    Write-Host $command
    Invoke-Expression $command   
}

function cleanProject {
    if (Test-Path -Path .\bin) {
        Write-Host "removing bin"
        Remove-Item -Force -Recurse -Path .\bin
    }
}

function testProject {
    runCommand "go test -v ./test"
}

function buildProject {
    $env:GOOS="windows"
    $env:GOARCH="amd64"
     
    runCommand "go build -v -o bin/winquit.exe ./cmd/winquit"
}

if (($args.Count -gt 0) -and ($args[0] -eq "clean")) {
    cleanProject
    Exit 0
}

if (($args.Count -gt 0) -and ($args[0] -eq "test")) { 
    testProject
    exit 0
}

if (($args.Count -gt 0) -and ($args[0] -eq "build")) { 
    buildProject
    exit 0
}

cleanProject
buildProject