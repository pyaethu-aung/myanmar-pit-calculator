//go:build ignore
// +build ignore

package main

// This file is intentionally ignored by the Go toolchain so there is no
// duplicate `main` package when building the project. The real CLI
// entrypoint is at ./cmd/pitcalc.

// To run the CLI during development:
//
//   go run ./cmd/pitcalc --income 5000000

func main() {}
