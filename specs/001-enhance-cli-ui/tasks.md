# Tasks: Enhance CLI UI with Modern Terminal Aesthetics

**Input**: Design documents from `/specs/001-enhance-cli-ui/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 Initialize Go project with `huh` and `clipboard` dependencies in `go.mod`
- [x] T002 [P] Configure `lipgloss` styles and color tokens in `cmd/pitcalc_bubbletea/main.go`
- [x] T003 [P] Setup bilingual translation map for EN/MY localization in `cmd/pitcalc_bubbletea/main.go`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T004 Define `huh.Form` structure and groups in `cmd/pitcalc_bubbletea/main.go`
- [x] T005 [P] Implement state transition logic from form to result in `cmd/pitcalc_bubbletea/main.go`
- [x] T006 [P] Create rounding and currency formatting utilities in `cmd/pitcalc_bubbletea/main.go`

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Visually Appealing Interface (Priority: P1) 🎯 MVP

**Goal**: Transform the manual state machine into a polished, grouped `huh` form.

**Independent Test**: Run `go run cmd/pitcalc_bubbletea/main.go` and verify the grouped form renders with `lipgloss` styling.

### Implementation for User Story 1

- [x] T007 [US1] Implement Language selection group (EN/MY) at form start in `cmd/pitcalc_bubbletea/main.go`
- [x] T008 [US1] Implement "Income" group (Salary, Bonus) in `cmd/pitcalc_bubbletea/main.go`
- [x] T009 [US1] Implement "Reliefs" group (Spouse, Children, Parents) in `cmd/pitcalc_bubbletea/main.go`
- [x] T010 [US1] Implement "Other" group (SSB, Life Insurance) in `cmd/pitcalc_bubbletea/main.go`
- [x] T011 [US1] Wrap form in a `lipgloss` styled header banner in `cmd/pitcalc_bubbletea/main.go`

**Checkpoint**: At this point, the data collection flow is modernized and functional.

---

## Phase 4: User Story 2 - Clear Validation Feedback (Priority: P2)

**Goal**: Ensure all inputs are validated in real-time with helpful error messages.

**Independent Test**: Enter invalid values (e.g., negative salary) and verify red error text appears immediately.

### Implementation for User Story 2

- [x] T012 [US2] Add numeric validation to salary and bonus fields in `cmd/pitcalc_bubbletea/main.go`
- [x] T013 [US2] Add count validation for dependents (children >= 0, parents 0-2) in `cmd/pitcalc_bubbletea/main.go`
- [x] T014 [US2] Add SSB cap validation (max 30,000 MMK/month) in `cmd/pitcalc_bubbletea/main.go`
- [x] T015 [US2] Implement localized error messages for all validation rules in `cmd/pitcalc_bubbletea/main.go`

**Checkpoint**: Inputs are now robust and provide instant feedback to the user.

---

## Phase 5: User Story 3 - Structured Summary & Export (Priority: P3)

**Goal**: Display results in a structured table and provide export options.

**Independent Test**: Complete the form, view the table, copy to clipboard, and export a JSON file.

### Implementation for User Story 3

- [x] T016 [US3] Implement `bubbles/table` for the final tax breakdown in `cmd/pitcalc_bubbletea/main.go`
- [x] T017 [US3] Implement side-by-side layout for Income vs Reliefs using `lipgloss` in `cmd/pitcalc_bubbletea/main.go`
- [x] T018 [US3] Add "Copy to Clipboard" keybinding using `atotto/clipboard` in `cmd/pitcalc_bubbletea/main.go`
- [x] T019 [US3] Implement export menu for TXT, JSON, and CSV formats in `cmd/pitcalc_bubbletea/main.go`
- [x] T020 [US3] Implement file system writers for each export format in `cmd/pitcalc_bubbletea/main.go`

**Checkpoint**: Final result is premium, readable, and portable.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] T021 [P] Ensure responsive layouts handle terminal window resizing in `cmd/pitcalc_bubbletea/main.go`
- [ ] T022 [P] Final audit of Burmese translations across all groups in `cmd/pitcalc_bubbletea/main.go`
- [ ] T023 Code cleanup, removing old `model.step` logic in `cmd/pitcalc_bubbletea/main.go`
- [ ] T024 Perform final manual verification of all features together

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
- **Polish (Final Phase)**: Depends on all user stories being complete

### Parallel Opportunities

- All Setup tasks (T001-T003) can run in parallel
- Foundational tasks T005 and T006 can run in parallel after T004
- User Stories can be integrated incrementally, but US1 is the core driver

---

## Implementation Strategy

### MVP First (User Story 1 & 2)

1. Complete Setup and Foundation.
2. Implement the grouped `huh` form (US1).
3. Add real-time validation (US2).
4. **STOP and VALIDATE**: Verify the input flow is superior to the old wizard.

### Incremental Delivery

1. Foundation → Inputs ready.
2. US1 + US2 → Polished Data Collection.
3. US3 → Polished Results & Portability.
4. Polish → Final Branding and Responsive UX.
