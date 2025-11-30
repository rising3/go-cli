# API Contract: Configure Function

**Module**: `internal/cmd/configure`  
**Feature**: [spec.md](../spec.md) | **Data Model**: [data-model.md](../data-model.md)  
**Created**: 2025-11-30

## Function Signature

```go
package configure

import "io"

// Configure creates or overwrites a configuration file at the specified target path.
func Configure(target string, opts ConfigureOptions) error
```

## Input Contract

### Parameters

#### `target` (string)

- **Type**: `string`
- **Required**: Yes
- **Description**: Absolute path to the configuration file to be created or overwritten
- **Constraints**:
  - Must be an absolute path (not validated by function, caller's responsibility)
  - Parent directory will be created if it doesn't exist
  - Example: `/Users/username/.config/mycli/default.yaml`

#### `opts` (ConfigureOptions)

- **Type**: `configure.ConfigureOptions`
- **Required**: Yes
- **Description**: Configuration options struct containing all necessary parameters
- **Fields**:
  - `Force` (bool): Overwrite existing file if true
  - `Edit` (bool): Launch editor after file creation if true
  - `NoWait` (bool): Don't wait for editor to exit if true (requires Edit=true)
  - `Data` (map[string]interface{}): Configuration data to marshal
  - `Format` (string): Output format - "yaml", "yml", or "json"
  - `Output` (io.Writer): Writer for informational messages (not used currently, reserved for future)
  - `ErrOutput` (io.Writer): Writer for error messages and status updates
  - `EditorLookup` (func() (string, []string, error)): Function to detect editor
  - `EditorShouldWait` (func(string, []string) bool): Function to determine if should wait for editor

### Preconditions

- `opts.Data` must not be nil (can be empty map)
- `opts.Format` must be "yaml", "yml", or "json" (not validated, caller's responsibility)
- `opts.ErrOutput` must not be nil
- `opts.EditorLookup` must not be nil if `opts.Edit` is true
- `opts.EditorShouldWait` must not be nil if `opts.Edit` is true

## Output Contract

### Return Values

#### `error`

- **Type**: `error`
- **Description**: Error if file operations fail, nil on success
- **Possible Values**:
  - `nil`: Success (file created/overwritten, editor launched if requested)
  - `error`: File system error (directory creation, file write, permission denied, etc.)
  - **Note**: Editor detection errors are NOT returned; they are logged to `opts.ErrOutput` and function returns nil

### Side Effects

#### File System

- Creates parent directory at `filepath.Dir(target)` with permission 0755 if it doesn't exist
- Creates file at `target` with permission 0644 if it doesn't exist
- Overwrites file at `target` if `opts.Force` is true
- Removes existing file before writing if `opts.Force` is true

#### I/O Streams

- Writes messages to `opts.ErrOutput`:
  - `"Config already exists, skipping initialization: <target>"` if file exists and Force=false
  - `"Wrote config: <target>"` on successful file write
  - `"No editor found: <error>"` if editor detection fails (Edit=true)

#### Process Execution

- Launches editor process if `opts.Edit` is true and editor is found:
  - Editor's stdin/stdout/stderr bound to `os.Stdin`, `os.Stdout`, `os.Stderr`
  - Waits for editor to exit if `opts.EditorShouldWait()` returns true
  - Runs in background if `opts.EditorShouldWait()` returns false

## Behavior Specification

### Normal Flow

1. **Directory Creation**
   - Call `os.MkdirAll(filepath.Dir(target), 0o755)`
   - If error: return error immediately

2. **File Existence Check**
   - Call `os.Stat(target)`
   - If file exists and `opts.Force` is false:
     - Write `"Config already exists, skipping initialization: <target>"` to `opts.ErrOutput`
     - Return nil (not an error)
   - If file exists and `opts.Force` is true:
     - Call `os.Remove(target)` (ignore errors)

3. **Data Marshaling**
   - If `opts.Format` is "yaml" or "yml": call `yaml.Marshal(opts.Data)`
   - Otherwise: call `json.MarshalIndent(opts.Data, "", "  ")`
   - If error: return error

4. **File Writing**
   - Call `os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)`
   - Write marshaled data
   - Close file
   - If error: return error
   - Write `"Wrote config: <target>"` to `opts.ErrOutput`

5. **Editor Launch** (if `opts.Edit` is true)
   - Call `opts.EditorLookup()`
   - If error:
     - Write `"No editor found: <error>"` to `opts.ErrOutput`
     - Return nil (error absorbed)
   - Create `exec.Cmd` with editor and target as argument
   - Set `cmd.Stdin = os.Stdin`, `cmd.Stdout = os.Stdout`, `cmd.Stderr = os.Stderr`
   - Determine wait based on `opts.EditorShouldWait(editor, args)`
   - Call `proc.Run(cmd, shouldWait, opts.ErrOutput)`
   - If error: return error

6. **Return**
   - Return nil on success

### Edge Cases

#### File Already Exists (Force=false)

- **Input**: target="/path/to/existing.yaml", opts.Force=false
- **Output**: nil (no error)
- **Side Effect**: Message written to ErrOutput, file not modified

#### File Already Exists (Force=true)

- **Input**: target="/path/to/existing.yaml", opts.Force=true
- **Output**: nil (assuming write succeeds)
- **Side Effect**: Existing file removed and recreated with new content

#### Parent Directory Doesn't Exist

- **Input**: target="/path/to/new/dir/config.yaml"
- **Output**: nil (assuming directory creation succeeds)
- **Side Effect**: Directory "/path/to/new/dir" created with permission 0755

#### Editor Not Found (Edit=true)

- **Input**: opts.Edit=true, EditorLookup returns error
- **Output**: nil (error absorbed)
- **Side Effect**: Error message written to ErrOutput, file still created successfully

#### Editor Launch Fails (Edit=true)

- **Input**: opts.Edit=true, editor found but launch fails
- **Output**: error from proc.Run
- **Side Effect**: File created, but editor launch error returned

#### Invalid Format

- **Input**: opts.Format="xml"
- **Output**: json.MarshalIndent is called (default behavior)
- **Side Effect**: File created with JSON format

## Error Handling Strategy

### Errors that STOP execution (return error)

- Directory creation failure (`os.MkdirAll`)
- Data marshaling failure (`yaml.Marshal`, `json.MarshalIndent`)
- File creation/write failure (`os.OpenFile`, `Write`)
- Editor launch failure (`proc.Run`) - **BUT** editor detection failure is absorbed

### Errors that are ABSORBED (logged and return nil)

- Editor detection failure (`opts.EditorLookup`) - file creation is considered successful

### Rationale

Editor detection failure doesn't prevent the primary goal (config file creation) from succeeding. Users can still manually edit the file. This maintains backward compatibility with existing behavior.

## Testing Contract

### Required Test Cases

1. **Basic file creation** (file doesn't exist)
   - Assert file exists after call
   - Assert file content matches marshaled Data
   - Assert file permission is 0644

2. **File exists, Force=false**
   - Assert function returns nil
   - Assert file content unchanged
   - Assert message written to ErrOutput

3. **File exists, Force=true**
   - Assert function returns nil
   - Assert file content updated
   - Assert message written to ErrOutput

4. **Directory creation**
   - Assert parent directories created with permission 0755

5. **YAML format**
   - Assert yaml.Marshal is used
   - Assert correct YAML syntax in file

6. **JSON format**
   - Assert json.MarshalIndent is used
   - Assert correct JSON syntax with 2-space indentation

7. **Edit=true, editor found**
   - Mock EditorLookup to return test editor
   - Mock proc.ExecCommand to capture command
   - Assert editor command constructed correctly
   - Assert stdin/stdout/stderr bound to os.*

8. **Edit=true, editor not found**
   - Mock EditorLookup to return error
   - Assert function returns nil (error absorbed)
   - Assert error message written to ErrOutput

9. **NoWait=true**
   - Mock EditorShouldWait to return false
   - Assert proc.Run called with shouldWait=false

10. **Directory creation failure**
    - Mock os.MkdirAll to fail
    - Assert function returns error

11. **File write failure**
    - Mock os.OpenFile to fail
    - Assert function returns error

### Mock Requirements

- `opts.Output`: `bytes.Buffer` (currently unused)
- `opts.ErrOutput`: `bytes.Buffer` to capture messages
- `opts.EditorLookup`: Custom function returning test values
- `opts.EditorShouldWait`: Custom function returning test values
- `proc.ExecCommand`: Mock via test helper (existing pattern)

## Dependencies

### External Packages

- `os`: File operations
- `os/exec`: Command execution
- `io`: Writer interface
- `path/filepath`: Path manipulation
- `gopkg.in/yaml.v3`: YAML marshaling
- `encoding/json`: JSON marshaling

### Internal Packages

- `github.com/rising3/go-cli/internal/proc`: Process execution utilities

### No Direct Dependencies

- `github.com/spf13/cobra`: Not used in internal package
- `github.com/rising3/go-cli/internal/stdio`: Removed (replaced with io.Writer)
- `github.com/rising3/go-cli/internal/editor`: Used only via EditorLookup function injection

## Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2025-11-30 | Initial contract definition for refactoring |

## References

- [Feature Specification](../spec.md)
- [Data Model](../data-model.md)
- [Echo Command Contract](../../001-echo-subcommand/contracts/echo-command.md) (reference pattern)
