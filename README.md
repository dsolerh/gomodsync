# gomodsync - Go Module Version Comparator and Sync Tool

A command-line tool that compares and synchronizes dependency versions between two `go.mod` files.

## Features

- **Check** dependency versions between two `go.mod` files
- **Sync** versions from a reference file to a target file
- **Remote references** - Fetch reference go.mod from URLs (GitHub, GitLab, etc.)
- Dry-run mode to preview changes before applying them
- Strict mode to enforce exact dependency matching
- Handles both direct and indirect dependencies
- Syncs Go version between files
- Preserves file permissions and structure
- Clear output showing all version differences

## Installation

```bash
# Build to bin directory
go build -o bin/gomodsync
```

Or install directly:

```bash
go install
```

## Usage

### Commands

#### sync - Synchronize versions

Updates the target go.mod file with versions from the reference file.

```bash
./bin/gomodsync sync -target <target-go.mod> -reference <reference-go.mod> [-dry-run] [-verbose]
```

**Options:**
- `-target`: Path to the target go.mod file to be modified (required)
- `-reference`: Path or URL to the reference go.mod file with desired versions (required)
- `-dry-run`: Show changes without modifying the target file (optional)
- `-verbose`: Show detailed list of all changes (optional)

**Example:**
```bash
# Preview changes (minimal output)
./bin/gomodsync sync -target ./project/go.mod -reference ./reference/go.mod -dry-run

# Preview changes (detailed output)
./bin/gomodsync sync -target ./project/go.mod -reference ./reference/go.mod -dry-run -verbose

# Apply changes
./bin/gomodsync sync -target ./project/go.mod -reference ./reference/go.mod

# Apply changes with verbose output
./bin/gomodsync sync -target ./project/go.mod -reference ./reference/go.mod -verbose

# Sync with remote reference from GitHub
./bin/gomodsync sync -target ./go.mod -reference https://raw.githubusercontent.com/user/repo/main/go.mod -dry-run
```

**Output (default):**
```
Dry-run mode: 14 change(s) identified but not applied.
```

**Output (with -verbose):**
```
Changes to be made:

  github.com/RoaringBitmap/roaring: v1.9.4 -> v1.10.0
  github.com/blugelabs/bluge: v0.2.2 -> v0.2.3
  golang.org/x/crypto: v0.47.0 -> v0.50.0
  ...

Dry-run mode: 14 change(s) identified but not applied.
```

**Output (successful sync):**
```
✓ Successfully updated ./project/go.mod (14 change(s) applied)
```

#### check - Check version differences

Compares dependency versions and reports mismatches. Useful for CI/CD pipelines.

```bash
./bin/gomodsync check -target <target-go.mod> -reference <reference-go.mod> [-strict] [-verbose]
```

**Options:**
- `-target`: Path to the target go.mod file to check (required)
- `-reference`: Path or URL to the reference go.mod file with desired versions (required)
- `-strict`: Fail if target has dependencies not in reference (optional)
- `-verbose`: Show detailed list of all mismatches (optional)

**Exit codes:**
- `0`: All versions match (or in non-strict mode, common dependencies match)
- `1`: Version mismatches found

**Example:**
```bash
# Check common dependencies (minimal output)
./bin/gomodsync check -target ./project/go.mod -reference ./reference/go.mod

# Check with detailed mismatch list
./bin/gomodsync check -target ./project/go.mod -reference ./reference/go.mod -verbose

# Strict mode - fail if target has extra dependencies
./bin/gomodsync check -target ./project/go.mod -reference ./reference/go.mod -strict -verbose

# Check against remote reference
./bin/gomodsync check -target ./go.mod -reference https://raw.githubusercontent.com/user/repo/main/go.mod -verbose
```

**Output (on mismatch, default):**
```
✗ Version check failed: 3 mismatch(es) found
Run with -verbose to see details
```

**Output (on mismatch, with -verbose):**
```
✗ Found 3 version mismatch(es):

  github.com/pkg/errors: v0.9.1 != v0.9.2
  golang.org/x/crypto: v0.47.0 != v0.50.0
  github.com/extra/dep: v1.0.0 (not in reference)
```

**Output (on success):**
```
✓ All dependency versions match!
```

## How It Works

### sync command
1. Parses both the target and reference `go.mod` files
2. Creates a version map from the reference file
3. Compares each dependency in the target file with the reference
4. Updates mismatched versions in the target file
5. Preserves file permissions, structure, and formatting
6. Writes the updated content back to the target file (unless in dry-run mode)

### check command
1. Parses both the target and reference `go.mod` files
2. Compares versions of common dependencies
3. In strict mode, also checks for dependencies only in target
4. Reports all mismatches and exits with appropriate code

## Using Remote References

The reference file can be either a local file path or a URL. This is useful for:
- Syncing with a canonical go.mod from a central repository
- Checking against production versions without local files
- CI/CD pipelines that reference remote standards

### Examples with URLs

**GitHub raw content:**
```bash
# Sync with a specific branch
./bin/gomodsync sync -target ./go.mod \
  -reference https://raw.githubusercontent.com/user/repo/main/go.mod \
  -dry-run -verbose

# Check against a tagged version
./bin/gomodsync check -target ./go.mod \
  -reference https://raw.githubusercontent.com/user/repo/v1.0.0/go.mod \
  -strict
```

**GitLab raw content:**
```bash
./bin/gomodsync sync -target ./go.mod \
  -reference https://gitlab.com/user/repo/-/raw/main/go.mod \
  -dry-run
```

**Any HTTP(S) URL:**
```bash
./bin/gomodsync check -target ./go.mod \
  -reference https://example.com/standards/go.mod \
  -verbose
```

## Use Cases

### CI/CD Pipeline
```bash
# Ensure dependencies match before deployment (show details on failure)
./bin/gomodsync check -target ./go.mod -reference ./prod/go.mod -strict -verbose
if [ $? -ne 0 ]; then
  echo "Dependency mismatch detected!"
  exit 1
fi
```

### Development Workflow
```bash
# Sync local project with team's standard versions (show what changed)
./bin/gomodsync sync -target ./go.mod -reference ./standards/go.mod -verbose
```

### Version Auditing
```bash
# Quick check without details
./bin/gomodsync check -target ./service/go.mod -reference ./security-baseline/go.mod

# Detailed check for auditing
./bin/gomodsync check -target ./service/go.mod -reference ./security-baseline/go.mod -verbose > audit-report.txt
```

### Automated Updates
```bash
# Preview changes silently in scripts
if ./bin/gomodsync sync -target ./go.mod -reference ./reference/go.mod -dry-run > /dev/null 2>&1; then
  echo "No updates needed"
else
  # Show details only when changes exist
  ./bin/gomodsync sync -target ./go.mod -reference ./reference/go.mod -dry-run -verbose
fi
```

### Sync with Remote Repository
```bash
# Keep multiple services in sync with a central go.mod
./bin/gomodsync sync -target ./service-a/go.mod \
  -reference https://raw.githubusercontent.com/company/standards/main/go.mod \
  -verbose

./bin/gomodsync sync -target ./service-b/go.mod \
  -reference https://raw.githubusercontent.com/company/standards/main/go.mod \
  -verbose
```

## Notes

- **sync**: Only dependencies that exist in both files will be updated
- **sync**: Dependencies unique to the target remain unchanged
- **check (non-strict)**: Only reports version mismatches for common dependencies
- **check (strict)**: Also reports dependencies that exist only in target
- File permissions are preserved when syncing
- Original structure and comments are maintained
