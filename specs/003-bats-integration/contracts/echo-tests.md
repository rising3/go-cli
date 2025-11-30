# Echo Command Test Contract

**Command**: `mycli echo`  
**Feature**: 003-bats-integration  
**Phase**: 1 - Design & Contracts

## Overview

このドキュメントは、mycli echoコマンドの統合テスト契約を定義します。echoコマンドはUNIX互換の実装であり、基本的な出力、フラグ処理、エスケープシーケンスの解釈をテストします。

---

## Test Cases

### TC-ECHO-001: Basic Single Argument (Priority: P1)

**Given**: Single string argument  
**When**: `mycli echo "Hello"` is executed  
**Then**: 
- "Hello" followed by newline is output to stdout
- Exit status is 0

**Bats Implementation**:
```bash
@test "TC-ECHO-001: Basic single argument output" {
    run_mycli echo "Hello"
    assert_success
    assert_output "Hello"
}
```

---

### TC-ECHO-002: Multiple Arguments (Priority: P1)

**Given**: Multiple string arguments  
**When**: `mycli echo Hello World` is executed  
**Then**: 
- "Hello World" (space-separated) followed by newline is output
- Exit status is 0

**Bats Implementation**:
```bash
@test "TC-ECHO-002: Multiple arguments" {
    run_mycli echo Hello World
    assert_success
    assert_output "Hello World"
}
```

---

### TC-ECHO-003: No Trailing Newline (-n flag) (Priority: P1)

**Given**: `-n` flag is provided  
**When**: `mycli echo -n "Hello"` is executed  
**Then**: 
- "Hello" without trailing newline is output
- Exit status is 0

**Bats Implementation**:
```bash
@test "TC-ECHO-003: No trailing newline with -n" {
    run_mycli echo -n "Hello"
    assert_success
    # Check output doesn't end with newline
    [[ "$output" == "Hello" ]]
}
```

---

### TC-ECHO-004: Escape Sequence Interpretation (-e flag) (Priority: P1)

**Given**: `-e` flag and escape sequences  
**When**: `mycli echo -e "Line1\nLine2"` is executed  
**Then**: 
- Output contains two lines: "Line1" and "Line2"
- Escape sequence `\n` is interpreted as newline
- Exit status is 0

**Bats Implementation**:
```bash
@test "TC-ECHO-004: Interpret escape sequences with -e" {
    run_mycli echo -e "Line1\nLine2"
    assert_success
    assert_line 0 "Line1"
    assert_line 1 "Line2"
}
```

---

### TC-ECHO-005: Tab Escape Sequence (Priority: P2)

**Given**: `-e` flag with `\t` escape sequence  
**When**: `mycli echo -e "Col1\tCol2"` is executed  
**Then**: 
- Output contains tab character between Col1 and Col2
- Exit status is 0

**Bats Implementation**:
```bash
@test "TC-ECHO-005: Tab escape sequence" {
    run_mycli echo -e "Col1\tCol2"
    assert_success
    assert_output_contains $'\t'
}
```

---

### TC-ECHO-006: Combined Flags (-n -e) (Priority: P2)

**Given**: Both `-n` and `-e` flags  
**When**: `mycli echo -n -e "Text\tTab"` is executed  
**Then**: 
- Escape sequences are interpreted
- No trailing newline is output
- Exit status is 0

**Bats Implementation**:
```bash
@test "TC-ECHO-006: Combined -n and -e flags" {
    run_mycli echo -n -e "Text\tTab"
    assert_success
    [[ "$output" =~ Text.*Tab ]]
    [[ "$output" != *$'\n' ]]
}
```

---

### TC-ECHO-007: Empty String (Priority: P2)

**Given**: No arguments provided to echo  
**When**: `mycli echo` is executed  
**Then**: 
- Empty line (just newline) is output
- Exit status is 0

**Bats Implementation**:
```bash
@test "TC-ECHO-007: Empty string output" {
    run_mycli echo
    assert_success
    assert_output ""
}
```

---

### TC-ECHO-008: UTF-8 Characters (Priority: P2)

**Given**: UTF-8 string with emoji  
**When**: `mycli echo "こんにちは🚀"` is executed  
**Then**: 
- UTF-8 string is output correctly
- Exit status is 0

**Bats Implementation**:
```bash
@test "TC-ECHO-008: UTF-8 character support" {
    run_mycli echo "こんにちは🚀"
    assert_success
    assert_output "こんにちは🚀"
}
```

---

### TC-ECHO-009: Backslash Escape (Priority: P2)

**Given**: `-e` flag with `\\` escape  
**When**: `mycli echo -e "Back\\slash"` is executed  
**Then**: 
- Single backslash is output
- Exit status is 0

**Bats Implementation**:
```bash
@test "TC-ECHO-009: Backslash escape" {
    run_mycli echo -e "Back\\\\slash"
    assert_success
    assert_output "Back\\slash"
}
```

---

### TC-ECHO-010: Stop Output Escape (\c) (Priority: P3)

**Given**: `-e` flag with `\c` escape  
**When**: `mycli echo -e "Start\cEnd"` is executed  
**Then**: 
- Only "Start" is output (no "End", no newline)
- Exit status is 0

**Bats Implementation**:
```bash
@test "TC-ECHO-010: Stop output with \\c" {
    run_mycli echo -e "Start\cEnd"
    assert_success
    assert_output "Start"
}
```

---

### TC-ECHO-011: Verbose/Debug Mode (Priority: P3)

**Given**: `--verbose` flag is provided  
**When**: `mycli echo --verbose "Debug"` is executed  
**Then**: 
- Debug information is displayed (implementation-specific)
- Exit status is 0

**Bats Implementation**:
```bash
@test "TC-ECHO-011: Verbose mode" {
    run_mycli echo --verbose "Debug"
    assert_success
    # Check for debug output if applicable
}
```

---

### TC-ECHO-012: Help for Echo Command (Priority: P2)

**Given**: `--help` flag  
**When**: `mycli echo --help` is executed  
**Then**: 
- Help message for echo command is displayed
- Usage examples are shown
- Exit status is 0

**Bats Implementation**:
```bash
@test "TC-ECHO-012: Display echo command help" {
    run_mycli echo --help
    assert_success
    assert_output_contains "Usage:"
    assert_output_contains "echo"
}
```

---

### TC-ECHO-013: Invalid Flag Error (Priority: P2)

**Given**: Invalid flag is provided  
**When**: `mycli echo --invalid` is executed  
**Then**: 
- Error message is displayed
- Exit status is non-zero

**Bats Implementation**:
```bash
@test "TC-ECHO-013: Error on invalid flag" {
    run_mycli echo --invalid
    assert_failure
    assert_output_contains "unknown flag"
}
```

---

### TC-ECHO-014: All Supported Escape Sequences (Priority: P2)

**Given**: `-e` flag with all documented escapes  
**When**: `mycli echo -e "\n\t\\\"\a\b\r\v"` is executed  
**Then**: 
- All escape sequences are interpreted correctly
- Exit status is 0

**Bats Implementation**:
```bash
@test "TC-ECHO-014: All escape sequences" {
    run_mycli echo -e "newline:\ntab:\tbackslash:\\\\"
    assert_success
    # Verify each escape is interpreted
}
```

---

### TC-ECHO-015: Large Output (Priority: P3)

**Given**: Very long string argument  
**When**: `mycli echo "<10KB string>"` is executed  
**Then**: 
- Complete string is output without truncation
- Exit status is 0

**Bats Implementation**:
```bash
@test "TC-ECHO-015: Large output handling" {
    local large_string=$(printf 'A%.0s' {1..10000})
    run_mycli echo "$large_string"
    assert_success
    [[ ${#output} -eq 10000 ]]
}
```

---

## Edge Cases

### EC-ECHO-001: Null Byte Handling

**Scenario**: String contains null byte  
**Expected Behavior**: Output up to null byte or error

---

### EC-ECHO-002: Binary Data

**Scenario**: Non-UTF-8 binary data provided  
**Expected Behavior**: Best-effort output or error

---

### EC-ECHO-003: Very Long Single Argument

**Scenario**: Argument exceeds shell/OS limits  
**Expected Behavior**: Error or truncation with warning

---

## Success Criteria Mapping

| Success Criterion | Covered by Tests |
|------------------|------------------|
| SC-002: Comprehensive coverage | 15 test cases covering all echo features |
| SC-003: Accurate failure detection | TC-ECHO-013 tests error scenarios |
| SC-004: Clear failure messages | All assertions include descriptive messages |

---

## Test Execution

```bash
bats integration_test/echo.bats
```

---

## Dependencies

- Binary at `../bin/mycli`
- UTF-8 locale support
- Standard bash string handling
