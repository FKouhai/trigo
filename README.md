# Trigo

> Read as tree-go (spanish word that means wheat) — a simple pet project that aims to work similarly to `tree` command.

## Usage

```
trigo [flags]

Flags:
  -a, --all                   Show hidden files and directories
  -d, --dir string            Directory to print tree for (default ".")
  -e, --exclude stringArray   Exclude files or directories by name
  -h, --help                  help for trigo
```

### Examples

```sh
# Print tree of the current directory
trigo

# Print tree of a specific directory
trigo -d /home/user/projects

# Show hidden files and directories
trigo --all

# Exclude specific directories
trigo --exclude vendor --exclude node_modules

# Combine flags
trigo -d ~/projects --all --exclude .git
```

## Features

- Prints a tree structure of the filesystem similar to the `tree` command
- Hidden files and directories (dotfiles) are excluded by default
- Respects `.gitignore` when run inside a git repository
- Exclude specific files or directories by name with `--exclude`

## Build

### With Nix

```sh
nix build
./result/bin/trigo
```

### With Go

```sh
go build -o trigo .
./trigo
```
