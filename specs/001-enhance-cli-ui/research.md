# Research: Modern CLI Enhancement with Charmbracelet

## Decision: Use `charmbracelet/huh` for Form Handling
Rationale: `huh` provides a high-level API for creating validated terminal forms. It handles cursor management, input validation, and theming out of the box, which reduces the boilerplate code compared to a manual `bubbletea` state machine.

Alternatives considered:
- Stick to existing `bubbletea` state machine: Rejected because it requires manual implementation of focus management, error rendering, and keybinding hints for every field.
- `survey`: Rejected as it doesn't integrate as seamlessly with the `bubbletea` event loop as `huh` does.

## Decision: Use `bubbles/table` for Summary Results
Rationale: `bubbles/table` provides a structured, scrollable table view that is perfect for displaying tax breakdowns. It supports column headers and easy styling via `lipgloss`.

Alternatives considered:
- `fmt.Printf`: Rejected because it lacks structure and doesn't scale well for complex layouts.
- Custom `lipgloss` layouts: A combination of `lipgloss` and `table` will be used for the best results.

## Decision: Branding with `lipgloss`
Rationale: `lipgloss` allows for sophisticated terminal styling (borders, colors, padding). A rounded border banner will be used for the title.

## Decision: Localization Strategy
Rationale: Since the requirement is bilingual (EN/MY) with English as default, we will implement a lightweight translation map or struct within `cmd/pitcalc_bubbletea`. This avoids heavy i18n frameworks while providing the necessary flexibility.

## Decision: Clipboard Integration
Rationale: We will use `github.com/atotto/clipboard` as it is the standard and most reliable library for cross-platform clipboard access in Go.

## Decision: Export Functionality
Rationale: After the calculation, we will provide a choice to export the result. 
- **TXT**: Plain text with terminal-like structure.
- **JSON**: Machine-readable format.
- **CSV**: Structured data for spreadsheets.
The user will select these via a `huh` menu after the summary screen.
