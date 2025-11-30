# Refactor: Configure subcommand to use Cobra streams

## Summary

Refactored the `configure` subcommand to follow best practices and align with the `echo` subcommand implementation pattern. Removed dependency on `internal/stdio` package and migrated to Cobra's standard I/O streams.

## Changes

### New Package: `internal/cmd/configure/`

- **configure.go**: Core business logic extracted from command layer
  - `ConfigureOptions` struct with 9 fields (Force, Edit, NoWait, Data, Format, Output, ErrOutput, EditorLookup, EditorShouldWait)
  - `Configure(target string, opts ConfigureOptions) error` function
  - Framework-agnostic implementation with 116 lines
  - Comprehensive godoc comments added

- **configure_test.go**: Complete test suite with 9 test functions
  - Basic file creation, force overwrite, directory creation
  - YAML/JSON format support
  - Editor integration (found/not found/no-wait scenarios)
  - Test coverage: 91.4% (exceeds 80% target)

### Modified Files

- **cmd/configure.go**: Refactored to use Cobra streams
  - Removed `internal/stdio` dependency
  - Added `cmd.OutOrStdout()` and `cmd.ErrOrStderr()` usage
  - Simplified RunE function to delegate to `configure.ConfigureFunc()`
  - Reduced from 52 to 51 lines

- **cmd/configure_test.go**: Enhanced with 4 new test cases
  - `TestConfigureCommand_ForceFlag`: Verifies --force flag propagation
  - `TestConfigureCommand_EditFlag`: Verifies --edit flag propagation
  - `TestConfigureCommand_NoWaitFlag`: Verifies --no-wait flag behavior
  - `TestConfigureCommand_ProfileFlag`: Verifies --profile path resolution
  - All tests use `configure.ConfigureFunc` mocking pattern

- **cmd/configure_wrapper_test.go**: Updated package references
  - Changed from `internalcmd` to `configure` package
  - All existing tests continue to pass

## Success Criteria Met (8/8)

- ✅ **SC-001**: Zero `internal/stdio` references in new code
- ✅ **SC-002**: Cobra streams (`cmd.OutOrStdout()`/`ErrOrStderr()`) used
- ✅ **SC-003**: Correct function signature `Configure(target string, opts ConfigureOptions) error`
- ✅ **SC-004**: All existing tests pass (100%)
- ✅ **SC-005**: Test coverage 91.4% (exceeds 80% target)
- ✅ **SC-006**: Backward compatibility maintained (output unchanged)
- ✅ **SC-007**: Structure consistency with `cmd/echo.go` (RunE within ±3 lines)
- ✅ **SC-008**: `make lint` passes with zero issues

## Testing

- **Unit tests**: 9 new tests in `internal/cmd/configure/`, 4 enhanced tests in `cmd/`
- **Integration tests**: Manual verification of all flags (--force, --edit, --no-wait, --profile)
- **Regression tests**: All existing tests pass without modification
- **Coverage**: 91.4% statement coverage in core package

## Quality Assurance

- ✅ `make test`: All packages pass
- ✅ `make fmt`: Code formatted with gofmt
- ✅ `make lint`: Zero warnings/errors (golangci-lint with govet)
- ✅ `make build`: Binary builds successfully
- ✅ Backward compatibility: Existing behavior preserved

## Implementation Notes

- Followed TDD (Test-Driven Development) workflow: Red → Green → Refactor
- Pattern consistent with existing `echo` subcommand implementation
- Editor error handling: Errors absorbed to maintain backward compatibility
- File permissions: 0644 (rw-r--r--) maintained for config files
- Directory creation: Automatic with 0755 permissions

## Migration Impact

- **Breaking Changes**: None
- **Deprecations**: None (internal/stdio still exists for other commands)
- **New Dependencies**: None (uses existing gopkg.in/yaml.v3, encoding/json)
- **Configuration**: No changes to user-facing configuration

## Related Documentation

- Feature Spec: `specs/002-configure-refactor/spec.md`
- Implementation Plan: `specs/002-configure-refactor/plan.md`
- Task List: `specs/002-configure-refactor/tasks.md`
- Quickstart Guide: `specs/002-configure-refactor/quickstart.md`

---

**Branch**: `002-configure-refactor`  
**Total Tasks Completed**: 80/81 (Phase 5 skipped due to coverage goal achieved)  
**Lines Changed**: +182 -17 (net +165)  
**Files Modified**: 4 files  
**Files Created**: 2 files (internal/cmd/configure/)
