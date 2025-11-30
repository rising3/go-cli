#!/usr/bin/env bats
# Integration tests for mycli echo command

load helpers/common
load helpers/assertions
load helpers/test_env

setup() {
    setup_test_env
}

teardown() {
    teardown_test_env
}

@test "TC-ECHO-001: Basic single argument output" {
    run_mycli echo "Hello"
    assert_success
    assert_output "Hello"
}

@test "TC-ECHO-002: Multiple arguments" {
    run_mycli echo Hello World
    assert_success
    assert_output "Hello World"
}

@test "TC-ECHO-003: No trailing newline with -n" {
    run_mycli echo -n "Hello"
    assert_success
    assert_output "Hello"
}

@test "TC-ECHO-004: Interpret escape sequences with -e" {
    run_mycli echo -e "Line1\nLine2"
    assert_success
    assert_line 0 "Line1"
    assert_line 1 "Line2"
}

@test "TC-ECHO-007: Empty string output" {
    run_mycli echo
    assert_success
    assert_output ""
}

@test "TC-ECHO-012: Display echo command help" {
    run_mycli echo --help
    assert_success
    assert_output_contains "Usage:"
    assert_output_contains "echo"
}

@test "TC-ECHO-013: Error on invalid flag" {
    run_mycli echo --invalid
    assert_failure
    assert_output_contains "unknown flag"
}
