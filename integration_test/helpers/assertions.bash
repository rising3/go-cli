#!/usr/bin/env bash
# Custom assertion functions for mycli integration tests

#
# assert_success - Assert that the last command succeeded (exit status 0)
#
# Checks that $status is 0.
#
# Example:
#   run_mycli --help
#   assert_success
#
assert_success() {
    if [[ "$status" -ne 0 ]]; then
        echo "Expected success (exit status 0) but got: $status"
        echo "Output: $output"
        return 1
    fi
}

#
# assert_failure - Assert that the last command failed (exit status non-zero)
#
# Checks that $status is not 0.
#
# Example:
#   run_mycli --invalid-flag
#   assert_failure
#
assert_failure() {
    if [[ "$status" -eq 0 ]]; then
        echo "Expected failure (non-zero exit status) but command succeeded"
        echo "Output: $output"
        return 1
    fi
}

#
# assert_output - Assert that output exactly matches expected string
#
# Arguments:
#   $1 - Expected output string
#
# Example:
#   run_mycli echo "Hello"
#   assert_output "Hello"
#
assert_output() {
    local expected="$1"
    if [[ "$output" != "$expected" ]]; then
        echo "Output mismatch:"
        echo "Expected: '$expected'"
        echo "Got:      '$output'"
        return 1
    fi
}

#
# assert_output_contains - Assert that output contains a substring
#
# Arguments:
#   $1 - Substring that should be present in output
#
# Example:
#   run_mycli --help
#   assert_output_contains "Usage:"
#
assert_output_contains() {
    local substring="$1"
    if [[ "$output" != *"$substring"* ]]; then
        echo "Output does not contain expected substring:"
        echo "Expected substring: '$substring'"
        echo "Actual output: '$output'"
        return 1
    fi
}

#
# assert_output_regex - Assert that output matches a regular expression
#
# Arguments:
#   $1 - Regular expression pattern
#
# Example:
#   run_mycli --version
#   assert_output_regex "mycli version [0-9]+\.[0-9]+\.[0-9]+"
#
assert_output_regex() {
    local pattern="$1"
    if [[ ! "$output" =~ $pattern ]]; then
        echo "Output does not match regex pattern:"
        echo "Pattern: '$pattern'"
        echo "Output:  '$output'"
        return 1
    fi
}

#
# assert_line - Assert that a specific line matches expected content
#
# Arguments:
#   $1 - Line index (0-based)
#   $2 - Expected line content
#
# Example:
#   run_mycli echo -e "Line1\nLine2"
#   assert_line 0 "Line1"
#   assert_line 1 "Line2"
#
assert_line() {
    local line_index="$1"
    local expected="$2"
    
    # Check if line index is valid
    if [[ "$line_index" -ge "${#lines[@]}" ]]; then
        echo "Line index $line_index out of range (only ${#lines[@]} lines)"
        return 1
    fi
    
    local actual="${lines[$line_index]}"
    if [[ "$actual" != "$expected" ]]; then
        echo "Line $line_index mismatch:"
        echo "Expected: '$expected'"
        echo "Got:      '$actual'"
        return 1
    fi
}

#
# assert_file_exists - Assert that a file exists at the given path
#
# Arguments:
#   $1 - File path
#
# Example:
#   assert_file_exists "$TEST_CONFIG_HOME/mycli/default.yaml"
#
assert_file_exists() {
    local file_path="$1"
    if [[ ! -f "$file_path" ]]; then
        echo "File does not exist: $file_path"
        return 1
    fi
}
