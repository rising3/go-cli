#!/usr/bin/env bats
# BATS integration tests for cat subcommand

load helpers/common
load helpers/test_env
load helpers/assertions

setup() {
    setup_test_env
    TEST_FILE="${TEST_TEMP_DIR}/test.txt"
    TEST_FILE2="${TEST_TEMP_DIR}/test2.txt"
}

teardown() {
    teardown_test_env
}

# T086 [P] Basic file display
@test "cat displays file content" {
    printf "line1\nline2\nline3\n" > "${TEST_FILE}"
    
    run_mycli cat "${TEST_FILE}"
    assert_success
    assert_line 0 "line1"
    assert_line 1 "line2"
    assert_line 2 "line3"
}

# T087 [P] Number all lines with -n flag
@test "cat with -n flag numbers all lines" {
    printf "first\n\nthird\n" > "${TEST_FILE}"
    
    run_mycli cat -n "${TEST_FILE}"
    assert_success
    assert_line 0 "     1  first"
    assert_line 1 "     2  "
    assert_line 2 "     3  third"
}

# T088 [P] Number non-empty lines with -b flag
@test "cat with -b flag numbers nonempty lines" {
    printf "first\n\nthird\n\nfifth\n" > "${TEST_FILE}"
    
    run_mycli cat -b "${TEST_FILE}"
    assert_success
    assert_output_contains "     1  first"
    assert_output_contains "     2  third"
    assert_output_contains "     3  fifth"
}

# T089 [P] Show line ends with -E flag
@test "cat with -E flag shows line ends" {
    printf "line1\nline2\n" > "${TEST_FILE}"
    
    run_mycli cat -E "${TEST_FILE}"
    assert_success
    assert_line 0 "line1$"
    assert_line 1 "line2$"
}

# T090 [P] Show tabs with -T flag
@test "cat with -T flag shows tabs" {
    printf "col1\tcol2\tcol3\n" > "${TEST_FILE}"
    
    run_mycli cat -T "${TEST_FILE}"
    assert_success
    assert_line 0 "col1^Icol2^Icol3"
}

# T091 [P] Show control characters with -v flag
@test "cat with -v flag shows control chars" {
    printf "Hello\x07World\n" > "${TEST_FILE}"
    
    run_mycli cat -v "${TEST_FILE}"
    assert_success
    assert_line 0 "Hello^GWorld"
}

# T092 [P] Show all with -A flag
@test "cat with -A flag shows all" {
    printf "test\ttab\x07\n" > "${TEST_FILE}"
    
    run_mycli cat -A "${TEST_FILE}"
    assert_success
    assert_line 0 "test^Itab^G$"
}

# T093 [P] Read from stdin
@test "cat from stdin" {
    run bash -c "echo 'stdin content' | ${MYCLI_BINARY} cat"
    assert_success
    assert_output "stdin content"
}

# T094 [P] Multiple files concatenation
@test "cat multiple files" {
    echo "File 1" > "${TEST_FILE}"
    echo "File 2" > "${TEST_FILE2}"
    
    run_mycli cat "${TEST_FILE}" "${TEST_FILE2}"
    assert_success
    assert_line 0 "File 1"
    assert_line 1 "File 2"
}

# T095 [P] Nonexistent file error
@test "cat nonexistent file error" {
    run_mycli cat /nonexistent/file.txt
    assert_failure
    assert_output_contains "cat:"
    assert_output_contains "/nonexistent/file.txt"
}

# T096 [P] Directory error
@test "cat directory error" {
    mkdir -p "${TEST_TEMP_DIR}/subdir"
    
    run_mycli cat "${TEST_TEMP_DIR}/subdir"
    assert_failure
    assert_output_contains "cat:"
}

# T097 [P] Partial error continues processing
@test "cat partial error continues" {
    echo "Valid content" > "${TEST_FILE}"
    
    run_mycli cat "${TEST_FILE}" /nonexistent.txt "${TEST_FILE}"
    assert_failure
    assert_output_contains "Valid content"
    assert_output_contains "cat:"
}

# Additional: Empty file
@test "cat empty file produces no output" {
    touch "${TEST_FILE}"
    
    run_mycli cat "${TEST_FILE}"
    assert_success
    assert_output ""
}

# Additional: Dash argument for stdin
@test "cat with - argument reads stdin" {
    echo "File 1" > "${TEST_FILE}"
    
    run bash -c "echo 'stdin' | ${MYCLI_BINARY} cat ${TEST_FILE} - ${TEST_FILE}"
    assert_success
    assert_output_contains "File 1"
    assert_output_contains "stdin"
}

# Additional: Combined flags
@test "cat with combined flags -nE" {
    printf "line1\nline2\n" > "${TEST_FILE}"
    
    run_mycli cat -nE "${TEST_FILE}"
    assert_success
    assert_line 0 "     1  line1$"
    assert_line 1 "     2  line2$"
}

# Additional: Help flag
@test "cat --help shows usage" {
    run_mycli cat --help
    assert_success
    assert_output_contains "Concatenate FILE"
}
