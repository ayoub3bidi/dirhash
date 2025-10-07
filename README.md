# Dirhash
Calculating the checksum of a directory made easy.

## Dirhash CLI
Compute a deterministic checksum of a directory by hashing file contents and their relative paths, with support for glob-based exclusions.

### Usage
```sh
dirhash sha256 [--ignore pattern]... <directory>

# Examples
dirhash sha256 .
dirhash sha256 -x node_modules/** -x "**/*.log" .
LOG_LEVEL=debug dirhash sha256 .

# JSON output
dirhash sha256 -o json -x "**/*.log" .
```

### Ignore patterns
- Supports doublestar `**` for recursive matching via `github.com/bmatcuk/doublestar`.
- Patterns are resolved relative to the provided directory argument.
- Examples:
  - `-x .git/**` to exclude the entire Git directory
  - `-x "**/*.tmp"` to exclude all temporary files


### How to build
 1. First download dependencies
 ```sh
 go mod download
 ```
 2. Build the CLI
 ```sh
 make build-cli
 ```

### Notes for windows users
- Use quotes around patterns with `*` to avoid shell expansion: `-x "**/*.log"`.
- Paths are handled using Go's `filepath`, so separators are normalized automatically.
