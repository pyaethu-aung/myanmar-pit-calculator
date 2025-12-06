.PHONY: help cli bubbletea run-cli run-bubbletea build build-cli build-bubbletea test test-coverage clean

help:
	@echo "Myanmar PIT Calculator - Available commands:"
	@echo ""
	@echo "  make cli              Run CLI mode (standard input/output)"
	@echo "  make bubbletea        Run interactive mode (bubble tea TUI)"
	@echo "  make build            Build both binaries"
	@echo "  make build-cli        Build CLI binary"
	@echo "  make build-bubbletea  Build interactive mode binary"
	@echo "  make test             Run all unit tests"
	@echo "  make test-coverage    Run tests with coverage report"
	@echo "  make clean            Clean up binaries and coverage files"
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

test:
	go test ./...

test-coverage:
	go test ./... -v -coverprofile=coverage.out -covermode=atomic
	@echo ""
	@echo "Coverage report generated: coverage.out"
	@go tool cover -func=coverage.out | grep total | awk '{print "Total coverage: " $$3}'

clean:
	rm -f coverage.out
	rm -rf bin/
