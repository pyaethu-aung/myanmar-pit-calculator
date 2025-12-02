.PHONY: help cli bubbletea run-cli run-bubbletea build build-cli build-bubbletea

help:
	@echo "Myanmar PIT Calculator - Available commands:"
	@echo ""
	@echo "  make cli              Run CLI mode (standard input/output)"
	@echo "  make bubbletea        Run interactive mode (bubble tea TUI)"
	@echo "  make build-cli        Build CLI binary"
	@echo "  make build-bubbletea  Build interactive mode binary"
	@echo "  make help             Show this help message"

cli:
	go run ./cmd/pitcalc

bubbletea:
	go run ./cmd/pitcalc_bubbletea

run-cli: cli

run-bubbletea: bubbletea

build-cli:
	go build -o bin/pitcalc ./cmd/pitcalc

build-bubbletea:
	go build -o bin/pitcalc-bubbletea ./cmd/pitcalc_bubbletea

build: build-cli build-bubbletea
	@echo "✅ Built both binaries in bin/"
