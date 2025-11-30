#!/usr/bin/env bats
# Integration tests for mycli root command

# Load helper functions
load helpers/common
load helpers/assertions
load helpers/test_env

# Run before each test
setup() {
    setup_test_env
}

# Run after each test
teardown() {
    teardown_test_env
}

@test "TC-ROOT-001: Display help with no arguments" {
    run_mycli
    assert_success
    assert_output_contains "Usage:"
    assert_output_contains "Available Commands:"
    assert_output_contains "echo"
    assert_output_contains "configure"
}

@test "TC-ROOT-002: Display help with --help flag" {
    run_mycli --help
    assert_success
    assert_output_contains "Usage:"
}

@test "TC-ROOT-003: Display help with -h flag" {
    run_mycli -h
    assert_success
    assert_output_contains "Usage:"
}

@test "TC-ROOT-004: Display version with --version flag" {
    run_mycli --version
    assert_success
    assert_output_regex "mycli version"
}

@test "TC-ROOT-005: Error on invalid flag" {
    run_mycli --invalid-flag
    assert_failure
    assert_output_contains "unknown flag"
}

@test "TC-ROOT-006: Error on invalid subcommand" {
    run_mycli invalid-subcommand
    assert_failure
    assert_output_contains "unknown command"
}

@test "TC-ROOT-007: Use custom config path from env var" {
    export MYCLI_CONFIG="$TEST_CONFIG_HOME/custom"
    mkdir -p "$MYCLI_CONFIG"
    
    run_mycli --help
    assert_success
}

@test "TC-ROOT-008: Use profile from env var" {
    set_test_profile "test"
    
    run_mycli --help
    assert_success
}

@test "TC-ROOT-010: Execute from different directory" {
    cd "$TEST_TEMP_DIR"
    
    run "$MYCLI_BINARY" --help
    assert_success
    assert_output_contains "Usage:"
}
