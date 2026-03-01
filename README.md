# gh-pr-url

GitHub CLI extension to return the url for the current branch.

## Installation

You need to have the GitHub CLI already installed: <https://cli.github.com> and you need to be logged in as well: `gh login`.

To install this extension run:

```console
gh extension install pierskarsenbarg/gh-pr-url
```

## Usage

To get the URL, run the following: `gh pr-url`:

```console
❯ gh pr-url
https://github.com/pierskarsenbarg/gh-pr-url
```

For all options, including shortened flags:

```console
❯ gh pr-url --help
GitHub CLI extension for getting PR for current branch

Usage:
  gh-pr-url [flags]

Flags:
  -h, --help   help for gh-pr-url
```
