# Security Policy

## Supported Versions

We release patches for security vulnerabilities for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| latest  | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security issue, please follow these steps:

### 1. Do Not Disclose Publicly

Please do not create a public GitHub issue for security vulnerabilities.

### 2. Report via GitHub Security Advisory

The preferred method is to use [GitHub Security Advisories](https://github.com/yourusername/gomodsync/security/advisories/new):

1. Go to the Security tab
2. Click "Report a vulnerability"
3. Provide detailed information about the vulnerability

### 3. Alternative: Email Report

If you cannot use GitHub Security Advisories, email the maintainers with:
- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)

### What to Include

A good security report includes:
- Type of vulnerability (e.g., XSS, SSRF, etc.)
- Affected versions
- Step-by-step reproduction instructions
- Proof of concept or example code
- Potential impact and severity
- Any suggested mitigations

### Response Timeline

- **Initial Response**: Within 48 hours
- **Status Update**: Within 7 days
- **Fix Timeline**: Depends on severity
  - Critical: Within 7 days
  - High: Within 14 days
  - Medium: Within 30 days
  - Low: Next release cycle

## Security Considerations

### URL Fetching

This tool fetches go.mod files from URLs. While we validate HTTP responses:
- Only fetch from trusted sources
- Be cautious with URLs from untrusted input
- The tool runs with your user permissions

### File System Access

The tool:
- Modifies local go.mod files (sync command)
- Reads file system paths provided by the user
- Preserves file permissions

Always review changes with `--dry-run` before applying.

## Security Best Practices for Users

1. **Verify URLs**: Only use URLs from trusted sources
2. **Review Changes**: Use `--dry-run` to preview changes before applying
3. **Check Sources**: Verify the integrity of remote go.mod files
4. **Keep Updated**: Use the latest version of the tool
5. **Review Permissions**: The tool respects file permissions

## Known Security Considerations

- **G107 (gosec)**: The tool makes HTTP requests with user-provided URLs. This is intentional functionality for a CLI tool where the user explicitly provides the URL.

## Acknowledgments

We appreciate responsible disclosure and will acknowledge security researchers who report valid vulnerabilities (with permission).

## Questions?

For security-related questions that are not vulnerabilities, please open a regular GitHub issue or discussion.
