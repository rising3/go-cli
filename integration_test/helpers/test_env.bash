#!/usr/bin/env bash
# Test environment configuration utilities for mycli integration tests

#
# mock_editor - Set up a mock editor for configure command testing
#
# Creates a temporary shell script that simulates editor behavior.
# Sets the EDITOR environment variable to point to this mock script.
#
# Arguments:
#   $1 - Behavior mode:
#        "save"   - Simulate successful edit and save
#        "cancel" - Simulate user canceling without changes
#        "error"  - Simulate editor error
#
# Example:
#   mock_editor "save"
#   run_mycli configure
#
#   mock_editor "cancel"
#   run_mycli configure
#
mock_editor() {
    local behavior="${1:-save}"
    local mock_script="${TEST_TEMP_DIR}/mock-editor.sh"
    
    case "$behavior" in
        save)
            # Create a script that exits successfully (simulates save)
            cat > "$mock_script" <<'EOF'
#!/bin/sh
# Mock editor - save behavior
# If file exists, append a comment; if new, create with default content
if [ -f "$1" ]; then
    echo "# Edited by test" >> "$1"
else
    echo "# Test configuration" > "$1"
    echo "test: true" >> "$1"
fi
exit 0
EOF
            ;;
        cancel)
            # Create a script that exits with code 1 (simulates cancel)
            cat > "$mock_script" <<'EOF'
#!/bin/sh
# Mock editor - cancel behavior
# Don't modify the file, exit with failure
exit 1
EOF
            ;;
        error)
            # Create a script that exits with error code
            cat > "$mock_script" <<'EOF'
#!/bin/sh
# Mock editor - error behavior
echo "Error: Editor failed" >&2
exit 2
EOF
            ;;
        *)
            echo "Unknown mock_editor behavior: $behavior" >&2
            return 1
            ;;
    esac
    
    # Make the mock script executable
    chmod +x "$mock_script"
    
    # Set EDITOR to use our mock
    export EDITOR="$mock_script"
}

#
# set_test_profile - Set the MYCLI_PROFILE environment variable
#
# Arguments:
#   $1 - Profile name
#
# Example:
#   set_test_profile "production"
#   run_mycli configure
#
set_test_profile() {
    local profile="$1"
    export MYCLI_PROFILE="$profile"
}
