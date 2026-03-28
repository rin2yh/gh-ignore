# gh-ignore

A [gh](https://cli.github.com/) extension to fetch GitHub gitignore templates.

## Installation

Requires [GitHub CLI](https://cli.github.com/) (`gh`).

```bash
gh extension install rin2yh/gh-ignore
```

To upgrade to the latest version:

```bash
gh extension upgrade gh-ignore
```

## Usage

### List available templates

```bash
gh ignore list
```

### Fetch a template

```bash
gh ignore get <name>
```

The template is appended to `.gitignore` by default. Template names are case-insensitive.

```bash
gh ignore get go
gh ignore get python
```

### Options

| Option | Description |
| --- | --- |
| `-o, --output <file>` | Output file (default: `.gitignore`) |

```bash
gh ignore get Go -o my.gitignore
```
