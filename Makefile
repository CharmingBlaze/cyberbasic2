# CyberBasic 2 — common dev targets (Unix-friendly; on Windows use Git Bash or WSL, or run the go commands directly).

GO       ?= go
EXAMPLES := $(wildcard examples/*.bas)

.PHONY: build test lint examples clean

build:
	$(GO) build -o cyberbasic .

test:
	$(GO) test ./...

lint:
	golangci-lint run ./...

# Compile-check every example program (no run — safe in headless CI).
examples: build
	@set -e; for f in $(EXAMPLES); do echo "==> $$f"; ./cyberbasic --lint "$$f"; done

clean:
	rm -f cyberbasic cyberbasic.exe
