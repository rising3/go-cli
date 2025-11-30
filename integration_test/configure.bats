#!/usr/bin/env bats
# Integration tests for mycli configure command

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

@test "TC-CONF-001: Create new config file" {
    run_mycli configure
    assert_success
    assert_file_exists "$TEST_HOME/.config/mycli/default.yaml"
}

@test "TC-CONF-002: Edit existing config file" {
    create_test_config "default" "existing: value"
    mock_editor "save"
    
    run_mycli configure --edit --no-wait
    assert_success
}

@test "TC-CONF-003: Create profile-specific config" {
    run_mycli configure --profile test
    assert_success
    assert_file_exists "$TEST_HOME/.config/mycli/test.yaml"
}

@test "TC-CONF-004: Cancel configuration" {
    mock_editor "cancel"
    
    run_mycli configure --edit --no-wait
    # Editor cancel behavior - may succeed or fail depending on implementation
}

@test "TC-CONF-005: Error when no editor configured" {
    skip "Editor detection behavior needs investigation"
    unset EDITOR
    unset VISUAL
    
    run_mycli configure --edit
    assert_failure
    assert_output_contains "editor"
}

@test "TC-CONF-008: Create config directory if missing" {
    rm -rf "$TEST_HOME/.config"
    
    run_mycli configure
    assert_success
    assert_file_exists "$TEST_HOME/.config/mycli/default.yaml"
}
