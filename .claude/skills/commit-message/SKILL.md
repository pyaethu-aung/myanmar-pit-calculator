---
name: commit-message
description: Use when creating or amending git commits. Enforces atomic commits, the 50/72 subject/body rule, and Conventional Commits format.
allowed-tools: Bash(git log:*) Bash(git diff:*) Bash(git status:*) Bash(git add:*)
---

# Commit Message Rules

Follow these rules for every commit.

## Working tree status
```!
git status --short
```

## Staged changes
```!
git diff --staged
```

If the staged diff above is empty, inspect the working tree status above:

- If there are unstaged or untracked changes, infer which files belong to
  the same logical change based on their names and paths, then stage them
  with `git add <files>` and proceed.
- If there are no changes at all, stop and tell the user:
  "Nothing to commit — working tree is clean."
- If the changes span multiple unrelated concerns, stage only the files
  that form one logical change, tell the user what you staged and why,
  and note that the remaining files should be committed separately.

## Recent commit history
```!
git log --oneline -10
```

Use the history above to match this project's existing commit style
(types, scopes, level of detail in subjects). If the project deviates
from Conventional Commits, follow the project's established pattern.

## 1. Atomic Commits

Each commit must represent one logical, self-contained change.

- **One reason to exist**: do not mix a bug fix with a refactor or a dependency bump with a feature
- **Builds at every commit**: the codebase must compile and tests must pass at every commit
- **Reviewable in isolation**: a reviewer should understand the change without needing context from adjacent commits

If `git diff --staged` spans multiple concerns, stop and tell the user
which concerns you see, then ask them to stage and commit each one
separately using `git add -p` or by staging specific files.

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
- **Never** add a `Co-Authored-By` trailer — omit it even if suggested

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
fix(auth): prevent session token from expiring prematurely
```

Feature with body:
```
feat(api): add pagination support to list endpoints

Without pagination, list endpoints return all records in a single
response. This causes memory spikes and slow response times as
data grows. Adds cursor-based pagination with a default page size
of 20.
```

Dependency bump:
```
chore(deps): bump golang.org/x/sys from v0.38.0 to v0.43.0
```

CI change:
```
ci: add golangci-lint to pull request workflow
```
