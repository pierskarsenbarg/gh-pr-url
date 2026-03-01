# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`gh-pr-url` is a GitHub CLI extension that returns the pull request URL for the current git branch. It is installed via `gh extension install pierskarsenbarg/gh-pr-url` and invoked as `gh pr-url`.

## Commands

```bash
# Build
go build ./...

# Run locally (from within a git repo)
go run .

# Tidy dependencies
go mod tidy

# Install as a gh extension from local source
gh extension install .
```

## Architecture

This is a minimal single-command CLI built with [Cobra](https://github.com/spf13/cobra):

- `main.go` — entry point, delegates to `cmd.Execute()`
- `cmd/root.go` — all logic lives here in the root Cobra command

The command flow in `cmd/root.go`:

1. Opens the local git repo using `go-git` to determine the current HEAD branch
2. Uses `github.com/cli/go-gh/v2/pkg/repository` to identify the GitHub remote repo (owner/name)
3. Creates an authenticated REST client via `github.com/cli/go-gh/v2/pkg/api` (uses the user's existing `gh` auth)
4. Fetches repo metadata to get the default branch; exits silently if HEAD is on the default branch
5. Queries the GitHub REST API for open PRs with `head=<owner>:<branch>` and prints the HTML URL

Key dependencies:

- `github.com/cli/go-gh/v2` — gh CLI integration (auth, repo detection)
- `github.com/go-git/go-git/v6` — local git operations
- `github.com/google/go-github/v42` — GitHub REST API types
- `github.com/spf13/cobra` — CLI framework
