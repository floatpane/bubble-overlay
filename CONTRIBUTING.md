# Contributing to bubble-overlay

Thanks for your interest in contributing! This guide covers the basics.

## Getting Started

### Prerequisites

- [Go 1.26+](https://go.dev/dl/)
- A terminal that supports ANSI / SGR for manual testing

### Setup

```bash
git clone https://github.com/floatpane/bubble-overlay.git
cd bubble-overlay
go mod tidy
```

### Build & Test

```bash
go build ./...
go test ./...
```

### Linting

```bash
gofmt -l .
go vet ./...
golangci-lint run
```

## Making Changes

### Branch Naming

Create a branch from `master` using one of these prefixes:

- `feature/` — new functionality
- `fix/` — bug fixes
- `docs/` — documentation changes
- `refactor/` — code restructuring without behavior changes

### Commit Messages

We use [Conventional Commits](https://www.conventionalcommits.org/):

```
type(scope): short description
```

Common types: `feat`, `fix`, `docs`, `test`, `ci`, `chore`.

Examples:

```
feat(line): handle zero-width joiners in overlay width math
fix(block): preserve trailing newline when overlay extends past base
docs: add lipgloss popup example
```

### Before Submitting a PR

1. `gofmt -l .` is clean.
2. `go vet ./...` is clean.
3. `go test ./...` passes.
4. Keep PRs focused — one logical change per PR.
5. Write a clear PR description: **what** changed and **why**.

## Reporting Bugs

Open an issue using the bug report template. Include:

- A minimal reproducer (`base`, `overlay`, `row`, `col`)
- Expected vs. actual output (paste with visible escape sequences if possible)
- Go version, OS, terminal emulator

## Requesting Features

Open an issue using the feature request template. Describe the problem and your proposed solution.

## AI Policy

We welcome contributions that use AI-assisted tools (Copilot, Claude, ChatGPT, etc.). Contributors are fully responsible for any code they submit, regardless of how it was written.

**What we expect:**

- **Understand what you submit.** You should be able to explain every line of your PR.
- **Review AI output carefully.** AI tools produce plausible-looking code that is sometimes subtly wrong, insecure, or off-pattern. Verify before committing.
- **No AI-generated issues, reviews, or comments.** Discussions should be genuine human communication.
- **No AI-generated tests that don't actually test anything.** Tests must validate behavior, not just exist for coverage.
- **Attribute when appropriate.** A brief mention in the PR description is appreciated but not required.

**What we won't accept:**

- Bulk PRs of AI-generated refactors or "improvements" that weren't requested.
- Code that introduces hallucinated dependencies, APIs, or patterns.
- Contributions where the author clearly doesn't understand the changes.

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating, you agree to uphold a welcoming and respectful environment.
