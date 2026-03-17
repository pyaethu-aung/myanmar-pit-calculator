# Implementation Plan: Enhance CLI UI with Modern Terminal Aesthetics

**Branch**: `001-enhance-cli-ui` | **Date**: 2026-03-17 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/001-enhance-cli-ui/spec.md`

## Summary
The goal is to modernize the Myanmar PIT Calculator's CLI by replacing the manual state management for data collection with `charmbracelet/huh`. This will provide a more polished, form-based input experience organized into logical groups. Additionally, the CLI will support bilingual (English/Burmese) localization, "Copy to Clipboard" functionality using `atotto/clipboard`, and multi-format data export (TXT, JSON, CSV). The results summary will be restructured using `bubbles/table` and styled with `lipgloss` to provide a premium, readable breakdown of tax calculations.

## Technical Context

**Language/Version**: Go 1.25.5
**Primary Dependencies**: `github.com/charmbracelet/huh`, `github.com/charmbracelet/bubbletea`, `github.com/charmbracelet/lipgloss`, `github.com/charmbracelet/bubbles`, `github.com/atotto/clipboard`
**Storage**: N/A (State is ephemeral; exports are file-based)
**Testing**: `go test` for logic; manual verification for UI and localization.
**Target Platform**: Terminal (CLI)
**Project Type**: CLI tool
**Performance Goals**: Sub-10ms response to keypresses.
**Constraints**: Standard terminal emulator support (Unicode for Burmese).
**Scale/Scope**: Single CLI entry point (`cmd/pitcalc_bubbletea/main.go`).

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- [x] Principle: Library-First (Logic is in `pkg/pitcalc`)
- [x] Principle: CLI Interface (This feature enhances the primary interface)
- [x] Principle: Test-First (Logic already tested; UI tests are manual/interactive)

## Project Structure

### Documentation (this feature)

```text
specs/001-enhance-cli-ui/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
└── quickstart.md        # Phase 1 output
```

### Source Code (repository root)

```text
cmd/
└── pitcalc_bubbletea/
    └── main.go          # Main entry point for TUI
pkg/
└── pitcalc/
    └── calculator.go    # Tax calculation logic (Core)
```

**Structure Decision**: The feature is localized to the `cmd/pitcalc_bubbletea/main.go` file which handles the TUI. The core calculation logic remains in `pkg/pitcalc`.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| New Dependency (`huh`) | Native form handling & aesthetics | Manual `bubbletea` state machines are verbose and harder to style consistently. |
