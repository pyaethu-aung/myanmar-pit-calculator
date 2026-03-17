<!--
Sync Impact Report:
- Version: 1.0.0 → 1.1.0 (Principle removal & testing expansion)
- Removed sections: Semantic Branch Naming
- Modified sections: Atomic Commits (Added 50/72 rule), Test Coverage & Quality Gates (Increased min coverage to 90%)
- Templates status: All aligned.
- Follow-up: None
-->

# Myanmar PIT Calculator Constitution

## Core Principles

### I. Core Logic Isolation (NON-NEGOTIABLE)
Code MUST respect the separation between the terminal UI layer (`cmd/`) and the core calculation logic (`pkg/`).

**Rules:**
- **Core Logic** (`pkg/pitcalc/`): Contains pure business calculations. MUST NOT import UI dependencies or perform pure I/O unless defined by interfaces.
- **UI Layer** (`cmd/pitcalc_bubbletea/`): Handles all `charmbracelet/huh` forms, state management, and `bubbles/table` rendering.
- Cross-layer violations (e.g., `pkg/` formatting terminal output with `lipgloss`) are FORBIDDEN.

**Rationale:** Separation ensures the tax calculation logic remains pure, independently testable, and reusable if future interfaces (e.g., web, API) are developed.

### II. State Management & Aesthetics (REQUIRED)
The CLI UI MUST deliver a premium user experience using standard Charm ecosystem tools.

**Rules:**
- All data collection MUST use `charmbracelet/huh` for form-based input, organized into logical groups.
- Manual `bubbletea` state machines for data collection are FORBIDDEN in new code.
- Output summaries MUST be structured using `bubbles/table` and styled with `lipgloss`.
- The interface MUST support bilingual (English/Burmese) localization, defaulting to English.

**Rationale:** A consistent, modern TUI reduces cognitive load and provides a robust, error-resistant user experience.

### III. Atomic Commits (REQUIRED)
Commits MUST be atomic—one logical task per commit with a clear, descriptive message adhering to the 50/72 rule.

**Rules:**
- One commit = one complete, self-contained change.
- Commit messages MUST follow the 50/72 rule: The first line (subject) must be 50 characters or less, followed by a blank line, and the body must be wrapped at 72 characters.
- Commit messages follow convention: `<type>: <description>`.
- Avoid mixed-concern commits.

**Rationale:** Atomic commits simplify code review, enable precise debugging, and create meaningful project history. The 50/72 rule ensures commit messages are highly readable in Git logs.

### IV. Test Coverage & Quality Gates (NON-NEGOTIABLE)
All business logic MUST have comprehensive unit tests.

**Rules:**
- **Unit tests**: Required for all calculation logic (e.g., `pkg/pitcalc/*_test.go`).
- Minimum coverage: 90% for the `pkg/pitcalc` package.
- TUI layer (`cmd/`) requires manual/interactive verification defined in `tasks.md` or automated UI testing where feasible.

**Rationale:** The calculator handles financial rules. Comprehensive testing prevents regressions and ensures trust in the results.

### V. Error Handling & Validation (REQUIRED)
Input validation MUST be handled proactively, and errors MUST be clearly communicated to the user.

**Rules:**
- `huh.Form` fields MUST implement robust validation functions preventing invalid inputs (e.g., negative salaries, impossible dependent counts).
- Validation errors MUST be localized and displayed inline.
- Panics are restricted to unrecoverable startup failures only.

**Rationale:** Immediate, clear feedback prevents confusing failures downstream and enhances the premium feel of the CLI.

### VI. Export & Clipboard Hygiene (REQUIRED)
Export functionality MUST be cross-platform and reliable.

**Rules:**
- Clipboard interactions MUST use `github.com/atotto/clipboard`.
- File exports MUST handle variable permissions and invalid paths gracefully without crashing the application.
- Structured exports (JSON/CSV) MUST use standard Go library encoders.

**Rationale:** Users rely on exporting their tax breakdowns. Silent failures or panics during export break user trust.

### VII. Dependency Management (NON-NEGOTIABLE)
Dependencies MUST be kept minimal and secure.

**Rules:**
- Critical-severity vulnerabilities MUST be resolved before opening a PR.
- `go mod tidy` MUST be run before committing.

**Rationale:** Catching issues locally before push prevents vulnerable code from entering the repository.

### VIII. Documentation Maintenance (REQUIRED)
Project documentation MUST be kept current with code changes.

**Rules:**
- `README.md` MUST be updated when setup steps, UI flows, or features change.
- The `specs/` directory MUST accurately reflect the implemented behavior.
- Stale or misleading documentation is treated as a bug.

**Rationale:** Outdated documentation wastes onboarding time and leads to tribal knowledge dependency.

## Development Workflow

### Code Review Requirements
- All changes require peer review before merge.
- Reviewers MUST verify: constitution compliance, test coverage for logic, and visual consistency for TUI.

### Deployment Standards
- The default branch is `main`.
- Feature branches diverge from and merge into `main` via PRs.

### Technology Stack Constraints
- **Language**: Go 1.25.5+
- **Framework**: `charmbracelet/huh`, `bubbletea`, `lipgloss`
- **Testing**: Go standard library `testing` package
- **Exports**: `atotto/clipboard`

## Governance

This constitution supersedes all conflicting practices. All pull requests, code reviews, and design discussions MUST reference and comply with these principles.

**Amendment Process:**
1. Propose changes via issue/discussion.
2. Require consensus from maintainers.
3. Update constitution version.
4. Propagate changes to dependent templates (plan, spec, tasks).

**Tech Debt Grandfathering Policy:**
- Existing code that predates an amendment is exempt until modified ("boy scout rule").
- New files and new functions MUST always comply with the current constitution.

**Version**: 1.1.0 | **Ratified**: 2026-03-17 | **Last Amended**: 2026-03-17
