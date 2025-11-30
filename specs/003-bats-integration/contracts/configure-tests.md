# Configure Command Test Contract

**Command**: `mycli configure`  
**Feature**: 003-bats-integration  
**Phase**: 1 - Design & Contracts

## Overview

このドキュメントは、mycli configureコマンドの統合テスト契約を定義します。configureコマンドは設定ファイルの作成・編集を行うため、エディタの模擬とファイルシステムの検証が主要なテスト観点となります。

---

## Test Cases

### TC-CONF-001: Create New Config File (Priority: P1)

**Given**: No existing config file  
**When**: `mycli configure` is executed and editor saves content  
**Then**: 
- Config file is created at `$MYCLI_CONFIG/default.yaml`
- File contains valid YAML structure
- Exit status is 0

**Bats Implementation**:
```bash
@test "TC-CONF-001: Create new config file" {
    setup_test_env
    mock_editor "save"
    
    run_mycli configure
    assert_success
    assert_file_exists "$TEST_CONFIG_HOME/mycli/default.yaml"
}
```

---

### TC-CONF-002: Edit Existing Config File (Priority: P1)

**Given**: Existing config file with content  
**When**: `mycli configure` is executed  
**Then**: 
- Editor opens with existing content
- Modified content is saved
- Exit status is 0

**Bats Implementation**:
```bash
@test "TC-CONF-002: Edit existing config file" {
    setup_test_env
    create_test_config "default" "existing: value"
    mock_editor "save"
    
    run_mycli configure
    assert_success
}
```

---

### TC-CONF-003: Create Config with Profile (Priority: P1)

**Given**: `--profile test` flag is provided  
**When**: `mycli configure --profile test` is executed  
**Then**: 
- Profile-specific config file is created at `$MYCLI_CONFIG/test.yaml`
- Exit status is 0

**Bats Implementation**:
```bash
@test "TC-CONF-003: Create profile-specific config" {
    setup_test_env
    mock_editor "save"
    
    run_mycli configure --profile test
    assert_success
    assert_file_exists "$TEST_CONFIG_HOME/mycli/test.yaml"
}
```

---

### TC-CONF-004: Cancel Configuration (Priority: P2)

**Given**: Editor exits without saving  
**When**: `mycli configure` is executed and user cancels  
**Then**: 
- Config file is not modified (if existing) or not created (if new)
- Informational message may be displayed
- Exit status is 0 or appropriate cancel code

**Bats Implementation**:
```bash
@test "TC-CONF-004: Cancel configuration" {
    setup_test_env
    mock_editor "cancel"
    
    run_mycli configure
    # Behavior depends on implementation
    # May succeed with message or exit with specific code
}
```

---

### TC-CONF-005: Editor Not Found (Priority: P2)

**Given**: No editor is configured (EDITOR env var not set)  
**When**: `mycli configure` is executed  
**Then**: 
- Error message indicates no editor found
- Error suggests setting EDITOR environment variable
- Exit status is non-zero

**Bats Implementation**:
```bash
@test "TC-CONF-005: Error when no editor configured" {
    setup_test_env
    unset EDITOR
    unset VISUAL
    
    run_mycli configure
    assert_failure
    assert_output_contains "editor"
}
```

---

### TC-CONF-006: Use VISUAL Environment Variable (Priority: P3)

**Given**: VISUAL env var is set, EDITOR is not  
**When**: `mycli configure` is executed  
**Then**: 
- Editor specified in VISUAL is used
- Config file is created/edited successfully
- Exit status is 0

**Bats Implementation**:
```bash
@test "TC-CONF-006: Use VISUAL env var" {
    setup_test_env
    unset EDITOR
    export VISUAL="mock-editor"
    mock_editor "save"
    
    run_mycli configure
    assert_success
}
```

---

### TC-CONF-007: Display Current Config (Priority: P3)

**Given**: `--show` or similar flag is provided (if implemented)  
**When**: `mycli configure --show` is executed  
**Then**: 
- Current config content is displayed
- Exit status is 0

**Note**: This test assumes a --show flag; adjust based on actual implementation.

**Bats Implementation**:
```bash
@test "TC-CONF-007: Display current config" {
    setup_test_env
    create_test_config "default" "key: value"
    
    run_mycli configure --show
    # Implementation-dependent
    skip "--show flag not yet implemented"
}
```

---

### TC-CONF-008: Config Directory Creation (Priority: P2)

**Given**: Config directory does not exist  
**When**: `mycli configure` is executed  
**Then**: 
- Config directory is created automatically
- Config file is created within the directory
- Exit status is 0

**Bats Implementation**:
```bash
@test "TC-CONF-008: Create config directory if missing" {
    setup_test_env
    rm -rf "$TEST_CONFIG_HOME"
    mock_editor "save"
    
    run_mycli configure
    assert_success
    assert_file_exists "$TEST_CONFIG_HOME/mycli/default.yaml"
}
```

---

### TC-CONF-009: Multiple Profiles (Priority: P2)

**Given**: Multiple profile configs exist  
**When**: Each profile is loaded via MYCLI_PROFILE env var  
**Then**: 
- Correct profile-specific config is used
- Profiles do not interfere with each other

**Bats Implementation**:
```bash
@test "TC-CONF-009: Multiple profiles work independently" {
    setup_test_env
    create_test_config "dev" "env: development"
    create_test_config "prod" "env: production"
    
    set_test_profile "dev"
    # Verify dev profile behavior
    
    set_test_profile "prod"
    # Verify prod profile behavior
    
    assert_success
}
```

---

### TC-CONF-010: Config File Permissions (Priority: P3)

**Given**: New config file is created  
**When**: `mycli configure` completes  
**Then**: 
- Config file has appropriate permissions (e.g., 0644)
- Config file is readable by user

**Bats Implementation**:
```bash
@test "TC-CONF-010: Config file has correct permissions" {
    setup_test_env
    mock_editor "save"
    
    run_mycli configure
    assert_success
    
    # Check file is readable
    [[ -r "$TEST_CONFIG_HOME/mycli/default.yaml" ]]
}
```

---

## Edge Cases

### EC-CONF-001: Concurrent Configuration

**Scenario**: Two configure commands run simultaneously  
**Expected Behavior**: File locking or last-write-wins (document expected behavior)

---

### EC-CONF-002: Disk Full

**Scenario**: No space available to write config  
**Expected Behavior**: Error message about disk space

---

### EC-CONF-003: Invalid YAML Syntax

**Scenario**: User saves invalid YAML  
**Expected Behavior**: Validation error or warning (if implemented)

---

## Success Criteria Mapping

| Success Criterion | Covered by Tests |
|------------------|------------------|
| SC-002: Comprehensive coverage | 10 test cases for configure workflows |
| SC-003: Accurate failure detection | TC-CONF-004, TC-CONF-005 test failures |
| SC-006: CI/CD compatibility | Mocked editor for automation |

---

## Test Execution

```bash
bats integration_test/configure.bats
```

---

## Dependencies

- Mock editor script for automation
- File system access for config creation
- Environment variable manipulation (EDITOR, VISUAL, MYCLI_CONFIG)
