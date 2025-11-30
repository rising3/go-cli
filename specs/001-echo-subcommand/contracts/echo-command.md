# Echo Command Interface Contract

**Feature**: Echo サブコマンド実装  
**Date**: 2025-11-30  
**Version**: 1.0.0

## Command Specification

### Command Name
`mycli echo`

### Synopsis
```
mycli echo [flags] [args...]
```

### Description
Output arguments separated by spaces to stdout, with a trailing newline by default. Supports UNIX-standard options for newline suppression and escape sequence interpretation.

---

## Flags

### `-n, --no-newline`
- **Type**: Boolean flag
- **Default**: `false`
- **Description**: Suppress the trailing newline character.
- **Example**: 
  ```bash
  mycli echo -n "Hello"
  # Output: Hello (no newline)
  ```

### `-e, --escape`
- **Type**: Boolean flag
- **Default**: `false`
- **Description**: Enable interpretation of backslash escape sequences.
- **Supported Escape Sequences**:
  - `\n` - Newline
  - `\t` - Horizontal tab
  - `\\` - Backslash
  - `\"` - Double quote
  - `\a` - Alert (bell)
  - `\b` - Backspace
  - `\c` - Suppress further output (including trailing newline)
  - `\r` - Carriage return
  - `\v` - Vertical tab
- **Example**:
  ```bash
  mycli echo -e "Line1\nLine2"
  # Output:
  # Line1
  # Line2
  ```

### `-v, --verbose`
- **Type**: Boolean flag
- **Default**: `false`
- **Description**: Enable verbose debug output to stderr.
- **Example**:
  ```bash
  mycli echo --verbose "Hello"
  # stderr: [VERBOSE] Processing echo command with args: [Hello]
  # stdout: Hello
  ```

### `-h, --help`
- **Type**: Boolean flag
- **Default**: `false`
- **Description**: Display help information for the echo command.
- **Behavior**: Cobra automatically provides this flag.

---

## Arguments

### `args...`
- **Type**: String arguments (variadic)
- **Description**: Text arguments to be echoed to stdout, separated by spaces.
- **Constraints**:
  - Zero or more arguments (empty args results in blank line output)
  - No encoding restrictions (UTF-8 assumed)
  - Maximum argument count: Limited by OS (typically ~2MB command line)
- **Examples**:
  ```bash
  mycli echo "Hello" "World"
  # Output: Hello World
  
  mycli echo
  # Output: (blank line with newline)
  
  mycli echo "Special: !@#$%^&*()"
  # Output: Special: !@#$%^&*()
  ```

---

## Standard Streams

### Standard Output (stdout)
- **Purpose**: Normal command output (echoed text)
- **Encoding**: UTF-8
- **Termination**: Newline character (`\n`) by default, unless `-n` flag or `\c` escape is used

### Standard Error (stderr)
- **Purpose**: 
  - Error messages (invalid flags, command failures)
  - Help messages (when triggered by errors)
  - Verbose debug logs (`--verbose` flag)
- **Encoding**: UTF-8

---

## Exit Codes

| Exit Code | Meaning | Trigger Condition |
|-----------|---------|-------------------|
| `0` | Success | Command executed successfully |
| `1` | Error | Invalid flag specified, or other execution error |

---

## Behavior Specifications

### Flag Combinations

#### `-n` and `-e` together
```bash
mycli echo -n -e "Tab\there"
# Output: Tab	here (no trailing newline)
```

#### Order independence
```bash
mycli echo -e -n "Same\nOutput"
# Output (same as above):
# Same
# Output (no trailing newline)
```

### Edge Cases

#### Empty arguments
```bash
mycli echo "" "test"
# Output:  test (leading space due to empty first arg)
```

#### Large argument count
```bash
mycli echo $(seq 1 10000)
# Output: 1 2 3 ... 10000
# Requirement: Must complete in <100ms, use <100MB memory
```

#### Invalid escape sequences (with `-e`)
```bash
mycli echo -e "Invalid\zSequence"
# Output: Invalid\zSequence (literal backslash-z)
```

#### Escape sequence `\c` (suppress all further output)
```bash
mycli echo -e "Before\cAfter"
# Output: Before (everything after \c is suppressed, including newline)
```

#### Argument separator `--`
```bash
mycli echo -n -- -e
# Output: -e (treats -e as argument, not flag)
```

### Error Handling

#### Invalid flag
```bash
mycli echo -x "test"
# stderr: Error: unknown shorthand flag: 'x' in -x
# stderr: Usage: mycli echo [flags] [args...]
# ... (help message displayed automatically)
# Exit code: 1
```

#### No errors for valid inputs
- All valid combinations of flags and arguments succeed
- Empty arguments are valid (outputs blank line)

---

## Help Message Format

### Short Description
```
Output text to stdout (UNIX echo clone)
```

### Long Description
```
Output arguments separated by spaces, with a trailing newline by default.
Supports -n (suppress newline) and -e (interpret escape sequences).
```

### Usage Examples
```
Examples:
  mycli echo "Hello, World!"
  mycli echo -n "No newline"
  mycli echo -e "Line1\nLine2"
```

---

## Performance Requirements

- **Startup Time**: <100ms for command execution (SC-001)
- **Help Display**: <50ms for `mycli echo --help` (SC-003)
- **Memory Usage**: <100MB for 10,000 arguments (SC-004)

---

## UNIX Compatibility

This command aims for compatibility with GNU coreutils `echo` for the following features:
- Basic text output with space-separated arguments
- `-n` flag for newline suppression
- `-e` flag for escape sequence interpretation

**Note**: The following GNU `echo` features are NOT supported:
- `-E` flag (disable escape interpretation, redundant with default behavior)
- `--version` flag (use `mycli version` instead)

---

## Testing Contract

### Unit Test Coverage
- ✅ All escape sequences (`\n`, `\t`, `\\`, `\"`, `\a`, `\b`, `\c`, `\r`, `\v`)
- ✅ Flag combinations (`-n`, `-e`, `-n -e`)
- ✅ Empty arguments and large argument counts
- ✅ Invalid escape sequences (literal output)
- ✅ Error handling (invalid flags)

### Integration Test Coverage
- ✅ Cobra command integration (`cmd.Execute()`)
- ✅ Stdout/stderr separation
- ✅ Exit code verification
- ✅ Help message display

---

## Constitution Compliance

### TDD必須
- ✅ Tests written before implementation (Red-Green-Refactor)

### パッケージ責務分離
- ✅ `cmd/echo.go`: Cobra integration only
- ✅ `internal/echo/`: Pure logic (framework-independent)

### ユーザーエクスペリエンス
- ✅ Clear help messages with examples
- ✅ Automatic help display on errors

### パフォーマンス要件
- ✅ <100ms execution, <50ms help, <100MB memory for 10K args

---

## Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | 2025-11-30 | Initial contract specification |
