---
name: commit-message
description: Use when creating or amending git commits. Enforces atomic commits, the 50/72 subject/body rule, and Conventional Commits format.
allowed-tools: Bash(git log:*) Bash(git diff:*) Bash(git status:*)
---

# Commit Message Rules

Follow these rules for every commit in this project.

## 1. Atomic Commits

Each commit must represent one logical, self-contained change.

- **One reason to exist**: do not mix a bug fix with a refactor or a dependency bump with a feature
- **Builds at every commit**: the codebase must compile and tests must pass at every commit
- **Reviewable in isolation**: a reviewer should understand the change without needing context from adjacent commits

**Wrong**: `"fix login bug and update dependencies and refactor auth"`
**Right**: three separate commits, one per concern

## 2. The 50/72 Rule

| Part         | Rule                                      |
|--------------|-------------------------------------------|
| Subject line | 50 characters or fewer (hard limit: 72)   |
| Body lines   | Wrap at 72 characters                     |
| Separator    | Always one blank line between subject and body |

## 3. Conventional Commits

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Types

| Type       | When to use                                             | SemVer |
|------------|---------------------------------------------------------|--------|
| `feat`     | New feature for the user                                | MINOR  |
| `fix`      | Bug fix for the user                                    | PATCH  |
| `chore`    | Maintenance, dependency updates, tooling                | —      |
| `ci`       | CI/CD configuration changes                             | —      |
| `docs`     | Documentation only                                      | —      |
| `refactor` | Code change that neither fixes a bug nor adds a feature | —      |
| `perf`     | Performance improvement                                 | PATCH  |
| `test`     | Adding or correcting tests                              | —      |
| `style`    | Formatting, whitespace, missing semicolons              | —      |
| `build`    | Build system or external dependency changes             | —      |
| `revert`   | Reverts a previous commit                               | —      |

### Scope (optional)

A noun in parentheses describing the section of the codebase:

```
feat(auth): add OAuth2 login flow
fix(ui): correct button alignment on mobile
chore(deps): bump golang.org/x/sys from v0.38.0 to v0.43.0
```

### Subject line rules

- Imperative mood: "add feature" not "added" or "adds"
- No capital letter after the colon
- No trailing period
- 50 characters or fewer

### Body rules

- Explain **what** and **why**, not how
- Wrap at 72 characters
- Separated from subject by one blank line

### Footer rules

- Format: `Token: value` (hyphens in token names, not spaces)
- Common tokens: `Fixes`, `Closes`, `Refs`, `Reviewed-by`, `BREAKING CHANGE`
- `BREAKING CHANGE` must be uppercase

### Breaking changes

Use `!` before the colon, or a `BREAKING CHANGE:` footer, or both:

```
feat(api)!: change income input from string to integer

BREAKING CHANGE: income values must now be integers.
String-based input is no longer accepted.
```

## Examples

Simple fix (subject only):
```
fix(calculator): handle negative income edge case
```

Feature with body:
```
feat(ui): add responsive table for tax breakdown

Replaces the fixed-width layout with a lipgloss table that
adapts to the terminal viewport. Improves readability on
narrow terminals.
```

Dependency bump:
```
chore(deps): bump charmbracelet/bubbletea from v1.3.9 to v1.3.10
```

CI change:
```
ci: add golangci-lint to pull request workflow
```
