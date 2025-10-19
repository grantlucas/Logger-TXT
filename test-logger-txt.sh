#!/bin/bash

# Test suite for logger-txt
# This script safely tests all logger-txt functionality without affecting production logs

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

# Test counters
TESTS_RUN=0
TESTS_PASSED=0
TESTS_FAILED=0

# Test environment setup
TEST_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEST_LOG_FILE="$TEST_DIR/test_log.txt"
LOGGER_SCRIPT="$TEST_DIR/logger-txt"

# Store original environment
ORIGINAL_LOGGERTXT_PATH="$LOGGERTXT_PATH"

print_header() {
    echo "================================================"
    echo "Logger-TXT Test Suite"
    echo "================================================"
    echo "Test directory: $TEST_DIR"
    echo "Test log file: $TEST_LOG_FILE"
    echo "Logger script: $LOGGER_SCRIPT"
    echo ""
}

setup_test_environment() {
    # Unset production environment variable to ensure we use local test file
    unset LOGGERTXT_PATH

    # Clean up any existing test log
    rm -f "$TEST_LOG_FILE"

    # Verify logger script exists and is executable
    if [[ ! -f "$LOGGER_SCRIPT" ]]; then
        echo -e "${RED}ERROR: logger-txt script not found at $LOGGER_SCRIPT${NC}"
        exit 1
    fi

    if [[ ! -x "$LOGGER_SCRIPT" ]]; then
        chmod +x "$LOGGER_SCRIPT"
    fi
}

cleanup_test_environment() {
    # Remove test log file
    rm -f "$TEST_LOG_FILE"

    # Restore original environment
    if [[ -n "$ORIGINAL_LOGGERTXT_PATH" ]]; then
        export LOGGERTXT_PATH="$ORIGINAL_LOGGERTXT_PATH"
    fi
}

run_test() {
    local test_name="$1"
    local test_function="$2"

    echo -n "Running $test_name... "
    TESTS_RUN=$((TESTS_RUN + 1))

    # Clean log file before each test to prevent pollution
    rm -f "$TEST_LOG_FILE"

    if $test_function; then
        echo -e "${GREEN}PASS${NC}"
        TESTS_PASSED=$((TESTS_PASSED + 1))
    else
        echo -e "${RED}FAIL${NC}"
        TESTS_FAILED=$((TESTS_FAILED + 1))
    fi
}

# Test basic logging functionality
test_basic_logging() {
    local test_message="Test log entry"

    # Use -f flag to specify our test log file
    "$LOGGER_SCRIPT" -f "$TEST_LOG_FILE" "$test_message" > /dev/null 2>&1

    # Check if log file was created and contains our message
    if [[ -f "$TEST_LOG_FILE" ]] && grep -q "$test_message" "$TEST_LOG_FILE"; then
        return 0
    else
        return 1
    fi
}

# Test append functionality
test_append_functionality() {
    local first_message="First log entry"
    local second_message="Second log entry"

    # Add first entry
    "$LOGGER_SCRIPT" -f "$TEST_LOG_FILE" "$first_message" > /dev/null 2>&1

    # Verify first entry exists
    if ! grep -q "$first_message" "$TEST_LOG_FILE"; then
        return 1
    fi

    # Add second entry (should append, not overwrite)
    "$LOGGER_SCRIPT" -f "$TEST_LOG_FILE" "$second_message" > /dev/null 2>&1

    # Verify both entries exist
    if grep -q "$first_message" "$TEST_LOG_FILE" && grep -q "$second_message" "$TEST_LOG_FILE"; then
        # Verify we have exactly 2 lines
        local line_count
        line_count=$(wc -l < "$TEST_LOG_FILE")
        if [[ $line_count -eq 2 ]]; then
            return 0
        fi
    fi

    return 1
}

# Test search across multiple entries
test_search_multiple_entries() {
    # Add various entries to test search across multiple lines
    "$LOGGER_SCRIPT" -f "$TEST_LOG_FILE" "Work task completed" > /dev/null 2>&1
    "$LOGGER_SCRIPT" -f "$TEST_LOG_FILE" "Personal project update" > /dev/null 2>&1
    "$LOGGER_SCRIPT" -f "$TEST_LOG_FILE" "Meeting with team" > /dev/null 2>&1
    "$LOGGER_SCRIPT" -f "$TEST_LOG_FILE" "Work review session" > /dev/null 2>&1
    "$LOGGER_SCRIPT" -f "$TEST_LOG_FILE" "Shopping trip" > /dev/null 2>&1

    # Search for "work" (case insensitive) - should find 2 entries
    local work_results
    work_results=$("$LOGGER_SCRIPT" -f "$TEST_LOG_FILE" -s "work")
    local work_count
    work_count=$(echo "$work_results" | grep -c "work\|Work")

    # Search for "project" - should find 1 entry
    local project_results
    project_results=$("$LOGGER_SCRIPT" -f "$TEST_LOG_FILE" -s "project")
    local project_count
    project_count=$(echo "$project_results" | grep -c "project")

    # Verify we found the expected number of matches
    if [[ $work_count -eq 2 ]] && [[ $project_count -eq 1 ]]; then
        return 0
    else
        return 1
    fi
}

# Test type categorization
test_type_categorization() {
    local test_message="Work task completed"
    local test_type="WORK"

    "$LOGGER_SCRIPT" -f "$TEST_LOG_FILE" -t work "$test_message" > /dev/null 2>&1

    # Check if the log contains the type in uppercase
    if grep -q "$test_type" "$TEST_LOG_FILE" && grep -q "$test_message" "$TEST_LOG_FILE"; then
        return 0
    else
        return 1
    fi
}

# Test project categorization
test_project_categorization() {
    local test_message="Project milestone reached"
    local test_project="(TESTPROJ)"

    "$LOGGER_SCRIPT" -f "$TEST_LOG_FILE" -p testproj "$test_message" > /dev/null 2>&1

    # Check if the log contains the project in uppercase with parentheses
    if grep -q "$test_project" "$TEST_LOG_FILE" && grep -q "$test_message" "$TEST_LOG_FILE"; then
        return 0
    else
        return 1
    fi
}

# Test both type and project together
test_type_and_project() {
    local test_message="Combined categorization test"

    "$LOGGER_SCRIPT" -f "$TEST_LOG_FILE" -t work -p combo "$test_message" > /dev/null 2>&1

    # Check if the log contains both WORK and (COMBO)
    if grep -q "WORK (COMBO)" "$TEST_LOG_FILE" && grep -q "$test_message" "$TEST_LOG_FILE"; then
        return 0
    else
        return 1
    fi
}

# Test display functionality
test_display_functionality() {
    # Add multiple entries
    for i in {1..5}; do
        "$LOGGER_SCRIPT" -f "$TEST_LOG_FILE" "Entry number $i" > /dev/null 2>&1
    done

    # Test default display (should show last 10, but we only have 5+previous entries)
    local output
    output=$("$LOGGER_SCRIPT" -f "$TEST_LOG_FILE")
    local line_count
    line_count=$(echo "$output" | wc -l)

    # Should show multiple lines
    if [[ $line_count -gt 1 ]]; then
        return 0
    else
        return 1
    fi
}

# Test display count functionality
test_display_count() {
    # Add multiple entries first
    for i in {1..5}; do
        "$LOGGER_SCRIPT" -f "$TEST_LOG_FILE" "Entry number $i" > /dev/null 2>&1
    done

    # Test showing specific number of lines
    local output
    output=$("$LOGGER_SCRIPT" -f "$TEST_LOG_FILE" -c 2)
    local line_count
    line_count=$(echo "$output" | wc -l)

    # Should show exactly 2 lines
    if [[ $line_count -eq 2 ]]; then
        return 0
    else
        return 1
    fi
}

# Test default 10-line limit with more than 10 entries
test_default_ten_line_limit() {
    # Add 15 entries to test the 10-line default limit
    for i in {1..15}; do
        "$LOGGER_SCRIPT" -f "$TEST_LOG_FILE" "Entry number $i" > /dev/null 2>&1
    done

    # Get default output (should be last 10 lines)
    local output
    output=$("$LOGGER_SCRIPT" -f "$TEST_LOG_FILE")
    local line_count
    line_count=$(echo "$output" | wc -l)

    # Should show exactly 10 lines
    if [[ $line_count -eq 10 ]]; then
        # Verify it shows the last 10 entries (6-15)
        if echo "$output" | grep -q "Entry number 15" && echo "$output" | grep -q "Entry number 6" && ! echo "$output" | grep -q "Entry number 5"; then
            return 0
        fi
    fi

    return 1
}

# Test case-insensitive search
test_search_insensitive() {
    # Add multiple entries to properly test search
    "$LOGGER_SCRIPT" -f "$TEST_LOG_FILE" "First entry" > /dev/null 2>&1
    "$LOGGER_SCRIPT" -f "$TEST_LOG_FILE" "SearchTest entry" > /dev/null 2>&1
    "$LOGGER_SCRIPT" -f "$TEST_LOG_FILE" "Another entry" > /dev/null 2>&1
    "$LOGGER_SCRIPT" -f "$TEST_LOG_FILE" "Final entry" > /dev/null 2>&1

    # Search for lowercase version of entry
    local output
    output=$("$LOGGER_SCRIPT" -f "$TEST_LOG_FILE" -s "searchtest")

    if echo "$output" | grep -q "SearchTest entry"; then
        return 0
    else
        return 1
    fi
}

# Test case-sensitive search
test_search_sensitive() {
    # Add multiple entries including case variations
    "$LOGGER_SCRIPT" -f "$TEST_LOG_FILE" "casesensitive lowercase" > /dev/null 2>&1
    "$LOGGER_SCRIPT" -f "$TEST_LOG_FILE" "CaseSensitive Entry" > /dev/null 2>&1
    "$LOGGER_SCRIPT" -f "$TEST_LOG_FILE" "CASESENSITIVE UPPER" > /dev/null 2>&1

    # Search for exact case
    local output
    output=$("$LOGGER_SCRIPT" -f "$TEST_LOG_FILE" -S "CaseSensitive")

    if echo "$output" | grep -q "CaseSensitive Entry" && ! echo "$output" | grep -q "casesensitive lowercase"; then
        return 0
    else
        return 1
    fi
}

# Test case-sensitive search failure
test_search_sensitive_fail() {
    # Add entries with different cases
    "$LOGGER_SCRIPT" -f "$TEST_LOG_FILE" "CaseSensitive Entry" > /dev/null 2>&1
    "$LOGGER_SCRIPT" -f "$TEST_LOG_FILE" "Another entry" > /dev/null 2>&1

    # Search for wrong case (should not find CaseSensitive Entry)
    local output
    output=$("$LOGGER_SCRIPT" -f "$TEST_LOG_FILE" -S "casesensitive")

    # Should not find the CaseSensitive entry
    if ! echo "$output" | grep -q "CaseSensitive Entry"; then
        return 0
    else
        return 1
    fi
}

# Test delete functionality - successful deletion with "Y"
test_delete_functionality() {
    local test_message="Entry to be deleted"

    # Add an entry
    "$LOGGER_SCRIPT" -f "$TEST_LOG_FILE" "$test_message" > /dev/null 2>&1

    # Verify it exists
    if ! grep -q "$test_message" "$TEST_LOG_FILE"; then
        return 1
    fi

    # Delete last entry (simulate 'Y' response)
    echo "Y" | "$LOGGER_SCRIPT" -f "$TEST_LOG_FILE" -x > /dev/null 2>&1

    # Verify it's gone
    if ! grep -q "$test_message" "$TEST_LOG_FILE"; then
        return 0
    else
        return 1
    fi
}

# Test delete functionality - cancelled deletion with "n"
test_delete_cancelled_n() {
    local test_message="Entry that should not be deleted"

    # Add an entry
    "$LOGGER_SCRIPT" -f "$TEST_LOG_FILE" "$test_message" > /dev/null 2>&1

    # Verify it exists
    if ! grep -q "$test_message" "$TEST_LOG_FILE"; then
        return 1
    fi

    # Cancel deletion (simulate 'n' response)
    echo "n" | "$LOGGER_SCRIPT" -f "$TEST_LOG_FILE" -x > /dev/null 2>&1

    # Verify it still exists
    if grep -q "$test_message" "$TEST_LOG_FILE"; then
        return 0
    else
        return 1
    fi
}

# Test delete functionality - cancelled deletion with random input
test_delete_cancelled_random() {
    local test_message="Entry that should not be deleted with random input"

    # Add an entry
    "$LOGGER_SCRIPT" -f "$TEST_LOG_FILE" "$test_message" > /dev/null 2>&1

    # Verify it exists
    if ! grep -q "$test_message" "$TEST_LOG_FILE"; then
        return 1
    fi

    # Cancel deletion (simulate random response that's not Y)
    echo "xyz" | "$LOGGER_SCRIPT" -f "$TEST_LOG_FILE" -x > /dev/null 2>&1

    # Verify it still exists
    if grep -q "$test_message" "$TEST_LOG_FILE"; then
        return 0
    else
        return 1
    fi
}

# Test file creation
test_file_creation() {
    local temp_log="/tmp/test_logger_creation.txt"
    rm -f "$temp_log"

    # Create log in non-existent file
    "$LOGGER_SCRIPT" -f "$temp_log" "File creation test" > /dev/null 2>&1

    # Check if file was created and is readable/writable
    if [[ -f "$temp_log" ]] && [[ -r "$temp_log" ]] && [[ -w "$temp_log" ]]; then
        rm -f "$temp_log"
        return 0
    else
        rm -f "$temp_log"
        return 1
    fi
}

# Test help and version functions
test_help_version() {
    # Test help output
    local help_output
    help_output=$("$LOGGER_SCRIPT" -h 2>&1)
    local help_exit=$?

    # Test version output
    local version_output
    version_output=$("$LOGGER_SCRIPT" -V 2>&1)
    local version_exit=$?

    # Test lowercase version flag
    local version_v_output
    version_v_output=$("$LOGGER_SCRIPT" -v 2>&1)
    local version_v_exit=$?

    # Verify output content and exit codes
    if echo "$help_output" | grep -q "Usage:" && \
       echo "$version_output" | grep -q "Logger-TXT" && \
       echo "$version_v_output" | grep -q "Logger-TXT" && \
       [[ $help_exit -eq 0 ]] && \
       [[ $version_exit -eq 0 ]] && \
       [[ $version_v_exit -eq 0 ]]; then
        return 0
    else
        return 1
    fi
}

# Test timestamp format
test_timestamp_format() {
    "$LOGGER_SCRIPT" -f "$TEST_LOG_FILE" "Timestamp test" > /dev/null 2>&1

    # Check if timestamp matches expected format (DD/MM/YY HH:MM TZ)
    if grep -qE '[0-9]{2}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2} [+-][0-9]{4}' "$TEST_LOG_FILE"; then
        return 0
    else
        return 1
    fi
}

# Test delete functionality with spaces in path (Issue #30)
test_delete_with_spaces_in_path() {
    local temp_dir="/tmp/logger test dir"
    local temp_log="$temp_dir/test log.txt"

    # Create directory with spaces
    mkdir -p "$temp_dir"

    # Add two entries
    "$LOGGER_SCRIPT" -f "$temp_log" "First entry" > /dev/null 2>&1
    "$LOGGER_SCRIPT" -f "$temp_log" "Entry to be deleted" > /dev/null 2>&1

    # Verify both entries exist
    if ! grep -q "First entry" "$temp_log" || ! grep -q "Entry to be deleted" "$temp_log"; then
        rm -rf "$temp_dir"
        return 1
    fi

    # Delete last entry (simulate 'Y' response)
    echo "Y" | "$LOGGER_SCRIPT" -f "$temp_log" -x > /dev/null 2>&1

    # Verify first entry still exists and second is gone
    if grep -q "First entry" "$temp_log" && ! grep -q "Entry to be deleted" "$temp_log"; then
        rm -rf "$temp_dir"
        return 0
    else
        rm -rf "$temp_dir"
        return 1
    fi
}

# Test file creation with spaces in path (Issue #30)
test_file_creation_with_spaces() {
    local temp_dir="/tmp/logger test dir 2"
    local temp_log="$temp_dir/new log.txt"

    # Create directory with spaces
    mkdir -p "$temp_dir"

    # Remove log file if it exists
    rm -f "$temp_log"

    # Create log in non-existent file with spaces in path
    "$LOGGER_SCRIPT" -f "$temp_log" "File creation test with spaces" > /dev/null 2>&1

    # Check if file was created and is readable/writable
    if [[ -f "$temp_log" ]] && [[ -r "$temp_log" ]] && [[ -w "$temp_log" ]] && grep -q "File creation test with spaces" "$temp_log"; then
        rm -rf "$temp_dir"
        return 0
    else
        rm -rf "$temp_dir"
        return 1
    fi
}

print_summary() {
    echo ""
    echo "================================================"
    echo "Test Summary"
    echo "================================================"
    echo "Tests run: $TESTS_RUN"
    echo -e "Tests passed: ${GREEN}$TESTS_PASSED${NC}"
    echo -e "Tests failed: ${RED}$TESTS_FAILED${NC}"

    if [[ $TESTS_FAILED -eq 0 ]]; then
        echo -e "\n${GREEN}All tests passed!${NC}"
        return 0
    else
        echo -e "\n${RED}Some tests failed!${NC}"
        return 1
    fi
}

# Main execution
main() {
    print_header
    setup_test_environment

    # Trap to ensure cleanup on exit
    trap cleanup_test_environment EXIT

    # Run all tests
    run_test "Basic Logging" test_basic_logging
    run_test "Append Functionality" test_append_functionality
    run_test "Type Categorization" test_type_categorization
    run_test "Project Categorization" test_project_categorization
    run_test "Type and Project Combined" test_type_and_project
    run_test "Display Functionality" test_display_functionality
    run_test "Display Count" test_display_count
    run_test "Default 10-Line Limit" test_default_ten_line_limit
    run_test "Search Multiple Entries" test_search_multiple_entries
    run_test "Case-Insensitive Search" test_search_insensitive
    run_test "Case-Sensitive Search" test_search_sensitive
    run_test "Case-Sensitive Search (No Match)" test_search_sensitive_fail
    run_test "Delete Functionality (Y)" test_delete_functionality
    run_test "Delete Cancelled (n)" test_delete_cancelled_n
    run_test "Delete Cancelled (Random)" test_delete_cancelled_random
    run_test "File Creation" test_file_creation
    run_test "Help and Version" test_help_version
    run_test "Timestamp Format" test_timestamp_format
    run_test "Delete with Spaces in Path (Issue #30)" test_delete_with_spaces_in_path
    run_test "File Creation with Spaces (Issue #30)" test_file_creation_with_spaces

    print_summary
}

# Run the tests
main "$@"
