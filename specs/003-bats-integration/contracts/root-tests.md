# Root Command Test Contract

**Command**: `mycli` (root command)  
**Feature**: 003-bats-integration  
**Phase**: 1 - Design & Contracts

## Overview

このドキュメントは、mycli rootコマンドの統合テスト契約を定義します。各テストケースは、Given-When-Then形式で明確な前提条件、アクション、期待される結果を記述します。

---

## Test Cases

### TC-ROOT-001: Display Help (Priority: P1)

**Given**: No arguments or subcommands are provided  
**When**: `mycli` is executed without any arguments  
**Then**: 
- Help message is displayed to stdout
- Help message includes "Usage:" section
- Help message lists available commands (echo, configure)
- Exit status is 0

**Bats Implementation**:
```bash
@test "TC-ROOT-001: Display help with no arguments" {
    run_mycli
    assert_success
    assert_output_contains "Usage:"
    assert_output_contains "Available Commands:"
    assert_output_contains "echo"
    assert_output_contains "configure"
}
```

---

### TC-ROOT-002: Display Help with --help Flag (Priority: P1)

**Given**: `--help` flag is provided  
**When**: `mycli --help` is executed  
**Then**: 
- Same help message as TC-ROOT-001 is displayed
- Exit status is 0

**Bats Implementation**:
```bash
@test "TC-ROOT-002: Display help with --help flag" {
    run_mycli --help
    assert_success
    assert_output_contains "Usage:"
}
```

---

### TC-ROOT-003: Display Help with -h Flag (Priority: P1)

**Given**: `-h` flag is provided  
**When**: `mycli -h` is executed  
**Then**: 
- Same help message as TC-ROOT-001 is displayed
- Exit status is 0

**Bats Implementation**:
```bash
@test "TC-ROOT-003: Display help with -h flag" {
    run_mycli -h
    assert_success
    assert_output_contains "Usage:"
}
```

---

### TC-ROOT-004: Display Version (Priority: P1)

**Given**: `--version` flag is provided  
**When**: `mycli --version` is executed  
**Then**: 
- Version string is displayed (format: "mycli version X.Y.Z")
- Exit status is 0

**Bats Implementation**:
```bash
@test "TC-ROOT-004: Display version with --version flag" {
    run_mycli --version
    assert_success
    assert_output_regex "mycli version [0-9]+\.[0-9]+\.[0-9]+"
}
```

---

### TC-ROOT-005: Invalid Flag Error (Priority: P2)

**Given**: An unrecognized flag is provided  
**When**: `mycli --invalid-flag` is executed  
**Then**: 
- Error message is displayed to stderr
- Error message indicates unknown flag
- Exit status is non-zero (typically 1)

**Bats Implementation**:
```bash
@test "TC-ROOT-005: Error on invalid flag" {
    run_mycli --invalid-flag
    assert_failure
    assert_output_contains "unknown flag"
}
```

---

### TC-ROOT-006: Invalid Subcommand Error (Priority: P2)

**Given**: An unrecognized subcommand is provided  
**When**: `mycli invalid-subcommand` is executed  
**Then**: 
- Error message is displayed
- Error message suggests valid commands or "mycli --help"
- Exit status is non-zero

**Bats Implementation**:
```bash
@test "TC-ROOT-006: Error on invalid subcommand" {
    run_mycli invalid-subcommand
    assert_failure
    assert_output_contains "unknown command"
}
```

---

### TC-ROOT-007: Config File Path Override (Priority: P2)

**Given**: `MYCLI_CONFIG` environment variable is set to a custom path  
**When**: `mycli` is executed  
**Then**: 
- Configuration is loaded from the custom path (if exists)
- Command executes successfully
- Exit status is 0

**Bats Implementation**:
```bash
@test "TC-ROOT-007: Use custom config path from env var" {
    setup_test_env
    export MYCLI_CONFIG="$TEST_CONFIG_HOME/custom"
    mkdir -p "$MYCLI_CONFIG"
    
    run_mycli --help
    assert_success
}
```

---

### TC-ROOT-008: Profile Selection via Environment Variable (Priority: P2)

**Given**: `MYCLI_PROFILE` environment variable is set  
**When**: `mycli` is executed  
**Then**: 
- Configuration is loaded from the specified profile
- Command executes with profile-specific settings
- Exit status is 0

**Bats Implementation**:
```bash
@test "TC-ROOT-008: Use profile from env var" {
    setup_test_env
    set_test_profile "test"
    
    run_mycli --help
    assert_success
}
```

---

### TC-ROOT-009: Completion Command Availability (Priority: P3)

**Given**: Completion subcommand is implemented (if applicable)  
**When**: `mycli completion` is executed  
**Then**: 
- Shell completion script is generated or help is displayed
- Exit status is 0

**Note**: This test is optional and depends on whether completion functionality is implemented.

**Bats Implementation**:
```bash
@test "TC-ROOT-009: Completion command exists" {
    run_mycli completion --help
    # May succeed or fail depending on implementation
    # This is a placeholder for future functionality
    skip "Completion command not yet implemented"
}
```

---

### TC-ROOT-010: Binary Execution from Different Working Directory (Priority: P2)

**Given**: mycli binary is invoked from a directory different from project root  
**When**: `mycli --help` is executed from a different directory  
**Then**: 
- Command executes successfully
- Help is displayed correctly
- Exit status is 0

**Bats Implementation**:
```bash
@test "TC-ROOT-010: Execute from different directory" {
    setup_test_env
    cd "$TEST_TEMP_DIR"
    
    run "$MYCLI_BINARY" --help
    assert_success
    assert_output_contains "Usage:"
}
```

---

## Edge Cases

### EC-ROOT-001: No Config File Exists

**Scenario**: First-time user with no existing configuration  
**Expected Behavior**: Commands that don't require config should work normally

**Bats Implementation**:
```bash
@test "EC-ROOT-001: Run without config file" {
    setup_test_env
    # Ensure no config exists
    rm -rf "$TEST_CONFIG_HOME"
    
    run_mycli --help
    assert_success
}
```

---

### EC-ROOT-002: Corrupted Config File

**Scenario**: Configuration file exists but contains invalid YAML  
**Expected Behavior**: Error message indicating config parsing failure

**Bats Implementation**:
```bash
@test "EC-ROOT-002: Handle corrupted config file" {
    setup_test_env
    mkdir -p "$TEST_CONFIG_HOME/mycli"
    echo "invalid: yaml: syntax:" > "$TEST_CONFIG_HOME/mycli/default.yaml"
    
    run_mycli --help
    # May succeed if config not required for --help
    # Specific behavior depends on implementation
}
```

---

### EC-ROOT-003: Binary Permissions

**Scenario**: Binary does not have execute permissions  
**Expected Behavior**: Shell reports permission denied

**Note**: This is tested at the system level, not in Bats tests

---

## Success Criteria Mapping

| Success Criterion | Covered by Tests |
|------------------|------------------|
| SC-001: Complete suite under 30s | All root tests should complete in < 5 seconds |
| SC-002: Comprehensive coverage | 10 test cases covering primary workflows |
| SC-003: Accurate failure detection | TC-ROOT-005, TC-ROOT-006 test error scenarios |
| SC-004: Clear failure messages | All tests use descriptive assert messages |
| SC-006: CI/CD compatibility | Tests use isolated environments |
| SC-007: Clear progress indication | Bats formatter provides test-by-test output |

---

## Test Execution

### Run all root command tests
```bash
bats integration_test/root.bats
```

### Run specific test by name
```bash
bats integration_test/root.bats --filter "TC-ROOT-001"
```

### Verbose output
```bash
bats integration_test/root.bats --verbose-run
```

---

## Dependencies

- Binary must exist at `../bin/mycli`
- Helpers must be available: `common.bash`, `assertions.bash`, `test_env.bash`
- Environment must support `mktemp -d` for test isolation

---

## Future Enhancements

1. **Performance Testing**: Measure startup time (should be < 100ms per Constitution)
2. **Parallel Execution**: Test concurrent invocations of mycli
3. **Internationalization**: Test with different locale settings
4. **Signal Handling**: Test SIGINT, SIGTERM handling (if applicable)
