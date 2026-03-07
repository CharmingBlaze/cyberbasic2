# Build and run the first-game demo. From repo root: .\run_demo.ps1 [--debug]
# --debug: enable render trace (BeginDrawing, SyncFrame, DrawObject, etc.)
param([switch]$Debug)
Set-Location $PSScriptRoot
go build -o cyberbasic.exe .
if ($Debug) {
  .\cyberbasic.exe --debug examples\first_game.bas
} else {
  .\cyberbasic.exe examples\first_game.bas
}
