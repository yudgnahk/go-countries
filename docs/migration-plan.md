# Migration and Scope Expansion Plan

This document is the working execution checklist for expanding project scope from code-only lookups to country identifier resolution, and for rebranding the repository/module safely.

## Goals

- Expand scope to support country identifiers:
  - ISO 3166-1 alpha-2 codes
  - ISO 3166-1 alpha-3 codes
  - CIOC codes
  - Full country names and aliases
- Keep behavior deterministic and testable.
- Rebrand repository and description.
- Assess and execute module path migration with controlled breaking change.

## Rollout Strategy

Execute in 3 PRs to reduce risk:

1. API/behavior hardening (no rename)
2. Repository rebrand and docs refresh
3. Optional module path migration (`go.mod`) with migration notes

---

## PR1 - Scope Expansion and Deterministic Behavior

### Objectives

- Introduce a unified resolver for mixed user input.
- Preserve existing APIs.
- Remove nondeterministic behavior.

### Changes

- `emoji_flags.go`
  - Add unified resolver API (example):
    - `ResolveFlag(input string) (flag string, matched string, kind string)`
  - Resolver precedence:
    1. Exact code (`alpha-2`, `alpha-3`, `CIOC`, special code)
    2. Exact alias/name
    3. Fuzzy fallback
  - Update `GetFlagByName` to perform exact alias lookup before fuzzy matching.
  - Make special reverse lookup deterministic for `GetCode("🏴")`.
    - Avoid unordered map iteration for return value selection.

- `emoji_flags_test.go`
  - Add tests for resolver precedence and expected outputs.
  - Add deterministic tests for special flag reverse lookup.
  - Add exact alias tests (`USA`, `UK`, `UAE`) independent of fuzzy behavior.

- `README.md`
  - Add resolver usage section.
  - Clarify `GetFlag` as strict code lookup.

- `playground/main.go`
  - Fix `VietNamCode` to `VietnamCode`.

### Validation

- Run: `go test ./...`
- Ensure no flaky tests caused by map iteration randomness.

### Suggested PR title

- `feat: add unified country identifier resolver and deterministic lookups`

---

## PR2 - Repository Rebrand and Metadata Update

### Objectives

- Reposition project publicly without immediate module-breaking change.

### Changes

- GitHub repository settings (manual)
  - Rename repository.
  - Update repository description and topics.

- `README.md`
  - Update title and positioning text.
  - Update badges and repository links.
  - Keep module import path unchanged if PR3 is not merged yet.

### Validation

- Confirm GitHub URL redirect works.
- Confirm badges and links render correctly.
- Run: `go test ./...`

### Suggested PR title

- `docs: rebrand repository and update project positioning`

---

## PR3 - Module Path Migration (`go.mod`) [Breaking]

### Objectives

- Align module path with renamed repository.

### Changes

- `go.mod`
  - Change module path:
    - from `github.com/yudgnahk/go-emoji-flags`
    - to `github.com/yudgnahk/<new-repo-name>`

- `crawler/crawler.go`
  - Update module import path references.

- `README.md`
  - Update `go get` command and import snippets.
  - Add migration section with copy/paste replacement guidance.

- Repository-wide cleanup
  - Replace all remaining references to old module path.

### Validation

- Run: `go test ./...`
- Run: `go list ./...`
- Smoke test in a temp module:
  - `go mod init tmp`
  - `go get github.com/yudgnahk/<new-repo-name>@<tag>`

### Suggested PR title

- `chore!: migrate module path to new repository name`

---

## Impact Assessment: Repo Rename vs `go.mod` Change

### Rename repository only

- Lower risk for consumers in short term.
- GitHub URL redirects generally preserve access.
- Existing Go module imports can continue if module path remains unchanged.

### Rename repository and change module path

- Breaking for all downstream consumers.
- Consumers must update imports to the new module path.
- Requires clear migration notes and release communication.

---

## Release Plan

1. Release after PR1 (functional improvement, non-breaking intent).
2. Release after PR2 (branding/docs refresh).
3. Release after PR3 with explicit breaking-change migration notes.

Include in changelog:

- New resolver API
- Deterministic special-flag behavior
- Import path migration instructions (if PR3 merged)

---

## Migration Notes Template (for README)

```text
Old import:
  github.com/yudgnahk/go-emoji-flags

New import:
  github.com/yudgnahk/<new-repo-name>
```

Recommended command-line replacement:

```bash
rg -l 'github.com/yudgnahk/go-emoji-flags' | xargs sed -i '' 's#github.com/yudgnahk/go-emoji-flags#github.com/yudgnahk/<new-repo-name>#g'
```

Note: verify `sed -i` behavior on your OS/shell before bulk replacement.

---

## Current Known Issues to Address During PR1

- `GetCode("🏴")` currently has ambiguous mapping behavior due to special flag reuse.
- `playground/main.go` references `VietNamCode` but generated constant is `VietnamCode`.
