# Myanmar Personal Income Tax Calculator

## Project Layout

- `cmd/pitcalc/main.go`: Standard CLI mode (non-interactive)
- `cmd/pitcalc_bubbletea/main.go`: Interactive TUI mode with Bubble Tea
- `pkg/pitcalc`: Shared tax calculation library
- `main.go`: Ignored wrapper (contains `//go:build ignore`)

## Running the Application

You can run the calculator in two modes:

### Mode 1: Standard CLI (Non-interactive)

Run with simple input/output prompts:

```bash
make cli
```

Or directly:

```bash
go run ./cmd/pitcalc
```

### Mode 2: Interactive TUI (Bubble Tea)

Run with an interactive terminal user interface:

```bash
make bubbletea
```

Or directly:

```bash
go run ./cmd/pitcalc_bubbletea
```

## Building Binaries

Build both modes:

```bash
make build
```

This creates:
- `bin/pitcalc` - Standard CLI binary
- `bin/pitcalc-bubbletea` - Interactive TUI binary

Build individual binaries:

```bash
make build-cli        # CLI only
make build-bubbletea  # TUI only
```

## Testing

Run unit tests:

```bash
make test
```

Run tests with coverage report:

```bash
make test-coverage
```

This generates a `coverage.out` file that shows code coverage metrics for all packages:
- `pkg/pitcalc` - Tax calculation library
- `cmd/pitcalc` - Standard CLI mode
- `cmd/pitcalc_bubbletea` - Interactive TUI mode

## Continuous Integration

Unit tests are automatically run on GitHub Actions for:
- **Push events** to `main` branch
- **Pull requests** targeting the `main` branch

The workflow:
1. Sets up Go 1.25.2
2. Runs all unit tests with `make test`
3. Generates coverage report with `make test-coverage`

You can view the workflow in `.github/workflows/test.yml`

## Branch Naming Convention

Branch names must follow one of these prefixes for pull requests:

- **`feature/`** - New features (e.g., `feature/add-tax-brackets`)
- **`bugfix/`** - Bug fixes (e.g., `bugfix/fix-tax-calculation`)
- **`hotfix/`** - Hotfixes for production (e.g., `hotfix/critical-bug`)
- **`refactor/`** - Code refactoring (e.g., `refactor/improve-performance`)
- **`docs/`** - Documentation updates (e.g., `docs/update-readme`)
- **`test/`** - Test additions (e.g., `test/add-unit-tests`)
- **`ci/`** - CI/CD updates (e.g., `ci/add-github-actions`)

Use lowercase letters, numbers, hyphens, and underscores in branch names.

The branch naming convention is enforced by GitHub Actions (`.github/workflows/branch-name-check.yml`).

## Available Make Commands

- `make cli` - Run CLI mode
- `make bubbletea` - Run interactive TUI mode
- `make build` - Build both binaries
- `make build-cli` - Build CLI binary only
- `make build-bubbletea` - Build TUI binary only
- `make test` - Run all unit tests
- `make test-coverage` - Run tests with code coverage report
- `make clean` - Clean up built binaries and coverage files
- `make help` - Show all available commands
