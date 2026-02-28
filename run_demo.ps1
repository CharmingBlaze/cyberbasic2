# Build and run the first-game demo. From repo root: .\run_demo.ps1
Set-Location $PSScriptRoot
go build -o cyberbasic .
.\cyberbasic.exe examples\first_game.bas
