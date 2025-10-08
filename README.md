# Dirhash

[![CI](https://img.shields.io/github/actions/workflow/status/Think-iT-Labs/dirhash/ci.yml?branch=main&label=CI)](https://github.com/Think-iT-Labs/dirhash/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/tag/Think-iT-Labs/dirhash?label=release)](https://github.com/Think-iT-Labs/dirhash/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/Think-iT-Labs/dirhash)](https://go.dev/)
[![License](https://img.shields.io/github/license/Think-iT-Labs/dirhash)](LICENSE)

Calculating the checksum of a directory made easy.

## Dirhash CLI
Compute a deterministic checksum of a directory by hashing file contents and their relative paths, with support for glob-based exclusions.

### Usage
```sh
dirhash sha256 [--ignore pattern]... <directory>
```

### Examples
```sh
# Basic: hash current directory
dirhash sha256 .

# Ignore common folders and logs
dirhash sha256 -x node_modules/** -x ".git/**" -x "**/*.log" .

# Show which files are processed (debug logging)
LOG_LEVEL=debug dirhash sha256 .

# JSON output with per-file hashes (consumable by CI tools)
dirhash sha256 -o json -x "**/*.log" .
```

### Ignore patterns
- Supports doublestar `**` for recursive matching via `github.com/bmatcuk/doublestar`.
- Patterns are resolved relative to the provided directory argument.
- Examples:
  - `-x .git/**` to exclude the entire Git directory
  - `-x vendor/**` to exclude vendored dependencies
  - `-x "**/*.tmp"` to exclude all temporary files
  - `-x "**/*.log"` to exclude all log files

Tips:
- Quote patterns containing `*` or `**` to avoid shell expansion.
- Combine multiple `-x` flags to refine your selection.


### Build locally
 1. First download dependencies
 ```sh
 go mod download
 ```
 2. Build the CLI
 ```sh
 make build-cli
 ```

### Exit codes
- 0: success
- non-zero: fatal error occurred (see stderr for details)

### Notes for Windows users
- Use quotes around patterns with `*` to avoid shell expansion: `-x "**/*.log"`.
- Paths are handled via Go's `filepath`; separators are normalized automatically.
