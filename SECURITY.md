# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 0.x.x   | :white_check_mark: |

## Reporting a Vulnerability

**DO NOT open public GitHub issues for security vulnerabilities.**

Security vulnerabilities should be reported privately to allow time for fixes before public disclosure.

### How to Report

1. Go to https://github.com/scttfrdmn/cicada/security/advisories/new
2. Click "New draft security advisory"
3. Provide detailed information:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if known)

### What to Expect

- **Initial Response**: Within 48 hours
- **Status Updates**: Every 72 hours
- **Fix Timeline**: Depends on severity
  - Critical: 1-7 days
  - High: 1-2 weeks
  - Medium: 2-4 weeks
  - Low: Best effort

### Disclosure Policy

- We will coordinate disclosure timing with you
- Credit will be given to reporters (unless requested otherwise)
- Public disclosure after fix is released

## Security Best Practices

### For Users

**AWS Account Security**:
- Enable MFA on AWS root account
- Use IAM users with least-privilege policies
- Rotate access keys regularly
- Use AWS Organizations for multi-account setup
- Enable CloudTrail logging

**Cicada Configuration**:
- Use strong encryption for sensitive data
- Regularly update to latest version
- Review IAM policies in COMPLIANCE-CROSSWALK.md
- Enable audit logging
- Use AWS KMS for encryption keys

**Data Protection**:
- Classify data before upload (public, internal, confidential)
- Use NIST 800-171 mode for CUI
- Use HIPAA mode for PHI (see COMPLIANCE-ESSENTIALS.md)
- Review access permissions regularly
- Enable versioning on S3 buckets

### For Developers

**Code Security**:
- All code must pass gosec security linter
- Never commit AWS credentials
- Use AWS SDK credential chains
- Validate all user input
- Use prepared statements for database queries
- Follow OWASP Top 10 guidelines

**Dependencies**:
- Keep dependencies updated
- Review go.sum for unexpected changes
- Use `go mod verify` before building
- Monitor security advisories

**Testing**:
- Test with least-privilege IAM policies
- Test error handling paths
- Fuzz test input validation
- Test rate limiting

## Security Features

Cicada implements multiple security layers:

### Encryption
- **At Rest**: S3 server-side encryption (SSE-S3 or SSE-KMS)
- **In Transit**: TLS 1.3 for all AWS API calls
- **Metadata**: Encrypted with same key as data

### Access Control
- **IAM Integration**: Native AWS IAM policies
- **Least Privilege**: Minimal required permissions
- **Audit Logs**: All operations logged to CloudTrail
- **MFA Support**: Can require MFA for sensitive operations

### Compliance
- **NIST 800-171**: For CUI and NIH genomic data
- **NIST 800-53**: For HIPAA/PHI
- **NIST IR 8481**: For OSTP research security

See COMPLIANCE-CROSSWALK.md for detailed control mappings.

## Known Security Considerations

### AWS Credentials
Cicada requires AWS credentials to function. Users are responsible for:
- Securing credentials on local systems
- Using appropriate IAM policies
- Rotating credentials regularly

### Local Data Cache
Cicada caches data locally during sync. Users should:
- Encrypt local disks (FileVault, LUKS, BitLocker)
- Secure workstation physical access
- Clear cache when decommissioning systems

### Network Security
- Cicada uses HTTPS for all AWS communication
- No network services exposed by default
- Web UI (when enabled) binds to localhost only by default

## Security Contacts

- **GitHub Security Advisories**: Preferred method
- **Email**: security@cicada.sh (when project has dedicated email)

## Recognition

We appreciate security researchers who help keep Cicada secure. Contributors will be acknowledged in:
- Security advisories
- CHANGELOG.md
- Project website (when launched)

## References

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [AWS Security Best Practices](https://aws.amazon.com/security/best-practices/)
- [NIST Cybersecurity Framework](https://www.nist.gov/cyberframework)
- [CIS AWS Foundations Benchmark](https://www.cisecurity.org/benchmark/amazon_web_services)
