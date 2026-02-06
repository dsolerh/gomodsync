# Contributing to gomodsync

Thank you for your interest in contributing to gomodsync! This document provides guidelines and instructions for contributing.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for everyone.

## How to Contribute

### Reporting Bugs

Before creating a bug report, please check existing issues to avoid duplicates.

When filing a bug report, include:
- A clear, descriptive title
- Steps to reproduce the issue
- Expected behavior
- Actual behavior
- Go version (`go version`)
- OS and version
- Any relevant logs or error messages

### Suggesting Features

Feature requests are welcome! Please:
- Use a clear, descriptive title
- Provide a detailed description of the proposed feature
- Explain why this feature would be useful
- Include examples of how it would be used

### Pull Requests

1. **Fork the repository** and create your branch from `main`

2. **Make your changes**
   - Follow the existing code style
   - Write clear, descriptive commit messages
   - Add tests for new functionality
   - Update documentation as needed

3. **Run tests and linting**
   ```bash
   make test
   make lint
   ```

4. **Ensure coverage doesn't decrease**
   ```bash
   make test-coverage
   ```

5. **Submit your pull request**
   - Provide a clear description of the changes
   - Reference any related issues
   - Ensure CI checks pass

## Development Setup

### Prerequisites

- Go 1.21 or higher
- Make (optional, but recommended)
- golangci-lint (for linting)

### Setup Steps

```bash
# Clone the repository
git clone https://github.com/yourusername/gomodsync.git
cd gomodsync

# Install dependencies
go mod download

# Build the project
make build

# Run tests
make test

# Run linter
make lint
```

## Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Run `make fmt` before committing
- Keep functions small and focused
- Write clear comments for exported functions
- Use meaningful variable names

## Testing

- Write tests for all new functionality
- Maintain or improve test coverage
- Use table-driven tests where appropriate
- Include both positive and negative test cases
- Test edge cases

Example test structure:
```go
func TestFeature(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"case1", "input1", "output1"},
        {"case2", "input2", "output2"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Feature(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

## Commit Messages

Follow conventional commits format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Adding or updating tests
- `refactor`: Code refactoring
- `chore`: Maintenance tasks
- `ci`: CI/CD changes

Examples:
```
feat(sync): add support for URL references

fix(check): handle missing go.mod gracefully

docs(readme): update installation instructions
```

## Project Structure

```
.
├── .github/          # GitHub workflows
├── bin/              # Built binaries (gitignored)
├── check.go          # Check command logic
├── check_test.go     # Check tests
├── commands.go       # CLI command handlers
├── fetch.go          # URL/file fetching
├── fetch_test.go     # Fetch tests
├── main.go           # Entry point
├── parser.go         # go.mod parsing
├── parser_test.go    # Parser tests
├── sync.go           # Sync command logic
├── sync_test.go      # Sync tests
├── types.go          # Type definitions
├── Makefile          # Build automation
└── README.md         # Documentation
```

## Release Process

Releases are automated via GitHub Actions:

1. Update version in relevant files
2. Create a git tag: `git tag -a v1.0.0 -m "Release v1.0.0"`
3. Push the tag: `git push origin v1.0.0`
4. GitHub Actions will build and create a release

## Questions?

Feel free to open an issue for questions or clarifications.

Thank you for contributing!
