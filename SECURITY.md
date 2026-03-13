# Security Policy

## Supported Versions

The following versions of Mire are currently supported with security updates:

| Version | Supported          |
| ------- | ------------------ |
| v0.1.x  | :white_check_mark: |
| v0.0.x  | :x:                |

## Reporting a Vulnerability

If you discover a security vulnerability within Mire, please send an email to the maintainers. All security vulnerabilities will be promptly addressed.

Please include the following information:

- Type of vulnerability
- Full paths of source file(s) related to the vulnerability
- Location of the affected source code (tag/branch/commit or direct URL)
- Any special configuration required to reproduce the issue
- Step-by-step instructions to reproduce the issue
- Proof-of-concept or exploit code (if possible)
- Impact of the issue, including how an attacker might exploit it

## Response Timeline

- **Initial Response**: Within 48 hours
- **Severity Assessment**: Within 7 days
- **Fix Timeline**: Based on severity (critical: 7 days, high: 14 days, medium: 30 days)

## Security Best Practices

When using Mire in your applications:

1. **Input Validation**: Always validate log input data, especially when logging user-provided data
2. **Sensitive Data**: Avoid logging sensitive information (passwords, API keys, tokens, etc.)
3. **Access Control**: Ensure log outputs are accessible only to authorized personnel
4. **Rotation**: Configure appropriate log rotation to prevent disk exhaustion
