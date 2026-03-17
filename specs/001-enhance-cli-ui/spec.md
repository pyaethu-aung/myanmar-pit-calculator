# Feature Specification: Enhance CLI UI with Modern Terminal Aesthetics

## Summary
The current CLI provides a functional step-by-step wizard but lacks modern terminal aesthetics. This feature aims to significantly improve the user experience by introducing polished styling, structural layouts, and better interactive components using `charmbracelet/huh` and other Charm libraries.

## Clarifications
### Session 2026-03-17
- Q: Should the UI support Burmese localization? → A: Yes, bilingual support (English/Burmese) but default to English.
- Q: Should results be exportable? → A: Yes, support both "Copy to Clipboard" and "Export to File".
- Q: How should form fields be organized? → A: Use `huh.Group` to organize fields into logical sections.
- Q: Which library should be used for clipboard support? → A: Use `atotto/clipboard`.
- Q: What file formats should be supported for export? → A: Support TXT, JSON, and CSV with a user choice.

## User Stories
- As a user, I want a visually appealing interface so that using the calculator is more pleasant.
- As a user, I want clear validation feedback so that I know when my input is incorrect.
- As a user, I want a structured summary of my tax calculation so that I can easily understand the results.

## Requirements
### 1. Data Collection (Forms)
- Refactor data collection to use `charmbracelet/huh`.
- Implement built-in input validation.
- Use beautiful default themes and keybinding management.

### 2. Summary Screen
- Replace linear `fmt.Printf` blocks with structured `bubbles/table`.
- Show Income and Reliefs side-by-side or clearly delineated.
- Display detailed breakdowns of all user inputs (Gross Income, Basic Relief, Parents, Spouse, Children, SSB) instead of just total sums.
- Include a formatted tax breakdown table (Limit vs. Amount).
- Cohesive color scheme using soft colors.

### 3. Branding & Polish
- Add an eye-catching title banner using `lipgloss`.
- Ensure responsive layout for different terminal sizes.

## Component Updates
- `main.go`: Integrate `huh` for forms, `lipgloss` for styling, and `bubbles/table` for results.

## Acceptance Criteria
- CLI UI should be interactive and visually enhanced.
- Inputs must be validated correctly.
- Summary screen must display information in a structured, readable table.
- Title banner must be visible at the top.
