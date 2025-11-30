#!/usr/bin/env bash
# Common test utilities for mycli integration tests

# Path to the mycli binary (convert to absolute path)
_script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
_project_root="$(cd "$_script_dir/../.." && pwd)"
export MYCLI_BINARY="${MYCLI_BINARY:-$_project_root/bin/mycli}"

#
# setup_test_env - Initialize isolated test environment
#
# Creates unique temporary directories and sets environment variables
# to ensure tests don't interfere with user's actual configuration.
#
# Sets the following global variables:
#   TEST_TEMP_DIR    - Unique temporary directory for this test run
#   TEST_CONFIG_HOME - Test-specific config directory
#   TEST_HOME        - Test-specific home directory
#
# Exports environment variables:
#   MYCLI_CONFIG - Points to test config directory
#   HOME         - Points to test home directory
#
setup_test_env() {
    # Create unique temporary directory
    TEST_TEMP_DIR=$(mktemp -d -t mycli-test.XXXXXX)
    
    # Create subdirectories for config and home
    TEST_CONFIG_HOME="${TEST_TEMP_DIR}/config"
    TEST_HOME="${TEST_TEMP_DIR}/home"
    
    mkdir -p "${TEST_CONFIG_HOME}/mycli"
    mkdir -p "${TEST_HOME}"
    
    # Override HOME to isolate from user's config
    # mycli will look for config at $HOME/.config/mycli/
    export HOME="${TEST_HOME}"
    
    # Note: Bats handles cleanup via teardown() function, trap not needed
    # Note: MYCLI_CONFIG not exported - let mycli use default $HOME/.config/mycli/
}

#
# cleanup_test_env - Remove temporary test directories
#
# Called automatically via trap or manually in teardown.
# Safe to call multiple times.
#
cleanup_test_env() {
    if [[ -n "${TEST_TEMP_DIR}" && -d "${TEST_TEMP_DIR}" ]]; then
        rm -rf "${TEST_TEMP_DIR}"
    fi
}

#
# teardown_test_env - Teardown test environment
#
# Should be called in the teardown() function of each test file.
#
teardown_test_env() {
    cleanup_test_env
}

#
# run_mycli - Execute mycli binary and capture output
#
# Executes the mycli binary with given arguments and stores results
# in Bats global variables: $status, $output, $lines
#
# Arguments:
#   $@ - Command line arguments to pass to mycli
#
# Example:
#   run_mycli somecommand "Hello World"
#   run_mycli --help
#   run_mycli configure --profile test
#
# Note: Uses array to prevent bash builtin conflicts when used with Bats' run command.
#
run_mycli() {
    local -a cmd_array=("$MYCLI_BINARY" "$@")
    run "${cmd_array[@]}"
}

#
# create_test_config - Create a test configuration file
#
# Arguments:
#   $1 - Profile name (default: "default")
#   $2 - YAML content for the config file
#
# Returns:
#   Prints the path to the created config file
#
# Example:
#   create_test_config "default" "key: value"
#   create_test_config "prod" "environment: production"
#
create_test_config() {
    local profile="${1:-default}"
    local content="${2:-}"
    local config_file="${TEST_CONFIG_HOME}/mycli/${profile}.yaml"
    
    # Ensure directory exists
    mkdir -p "$(dirname "$config_file")"
    
    # Write content to file
    printf '%s\n' "$content" > "$config_file"
    
    # Return the file path
    printf '%s\n' "$config_file"
}
