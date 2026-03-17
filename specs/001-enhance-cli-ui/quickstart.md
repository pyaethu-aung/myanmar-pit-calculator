# Quickstart: Enhance CLI UI

## Development Environment
1. Ensure You have Go 1.25.5+ installed.
2. Install dependencies:
   ```bash
   go get github.com/charmbracelet/huh
   go get github.com/charmbracelet/bubbles
   go get github.com/charmbracelet/bubbletea
   go get github.com/charmbracelet/lipgloss
   ```

## Local Execution
Run the enhanced CLI:
```bash
go run cmd/pitcalc_bubbletea/main.go
```

## Key Components
- `form := huh.NewForm(...)`: Defines the input steps.
- `t := table.New(...)`: Renders the results breakdown.
- `style := lipgloss.NewStyle()...`: Defines the visual theme.
