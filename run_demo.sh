#!/usr/bin/env bash
# Build and run the first-game demo. From repo root: ./run_demo.sh
set -e
cd "$(dirname "$0")"
go build -o cyberbasic .
./cyberbasic examples/first_game.bas
