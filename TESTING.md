# Testing Documentation for gomodsync

## Refactoring for Testability

The code has been refactored to separate business logic from CLI/IO operations, making it highly testable.

### Extracted Functions

1. **ParseGoMod** - Parses go.mod file data
2. **BuildVersionMap** - Creates a version map from a modfile
3. **CompareVersions** - Compares versions and identifies changes
4. **ApplyVersionChanges** - Applies version updates to a modfile
5. **SyncVersions** - Main business logic orchestrator

## Test Suite Overview

### TestBuildVersionMap
Tests the creation of version maps from modfiles:
- Simple dependencies
- Indirect dependencies
- Empty go.mod files

### TestCompareVersions
Tests version comparison logic:
- Versions that need updating
- No changes needed
- Partial overlap between files
- Reference with additional dependencies

### TestApplyVersionChanges
Tests applying version changes:
- Single change application
- Multiple changes application
- No changes (edge case)
- Verifies correct version updates

### TestSyncVersions
Integration tests for the complete sync workflow:
- Successful sync with multiple changes
- No changes needed (files already in sync)
- Partial overlap scenarios
- Target with more dependencies than reference

### TestParseGoMod
Tests go.mod file parsing:
- Valid go.mod files
- Invalid go.mod files
- Minimal go.mod files

## Running Tests

```bash
# Run all tests
go test

# Run with verbose output
go test -v

# Run with coverage
go test -cover

# Generate detailed coverage report
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Test Coverage

Current coverage: **27.9%** of statements

The coverage focuses on business logic functions. Lower overall percentage is due to:
- `main()` function handling CLI arguments and file I/O (not unit tested)
- Error handling paths in main (requires integration tests)

Business logic functions have near 100% coverage:
- BuildVersionMap
- CompareVersions
- ApplyVersionChanges
- SyncVersions

## Testing Philosophy

1. **Unit tests** for pure business logic functions
2. **Integration tests** for the SyncVersions workflow
3. **Edge cases** covered: empty files, no changes, partial overlaps
4. **Error cases** tested where applicable

## Future Test Improvements

Potential additions:
1. Integration tests with actual file I/O
2. Tests for concurrent modifications
3. Performance benchmarks for large go.mod files
4. Tests for malformed version strings
5. Tests for complex dependency chains

## Verification

The refactored code maintains full compatibility with the original implementation:

```bash
$ go test
PASS
ok      gomodsync  0.254s

$ ./gomodsync sync -target gomod.test -reference gomod.reference -dry-run
# Successfully identifies and reports 14 version differences
```
