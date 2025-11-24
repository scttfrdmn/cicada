# Cicada Credential Management Specification

**Version:** 0.3.0
**Status:** Design Specification
**Related Issue:** #26 (Provider Configuration System)

## Overview

Cicada must handle sensitive credentials (API tokens, passwords) securely while providing a flexible, user-friendly configuration experience. This document specifies the credential management system for v0.3.0.

## Design Principles

1. **Security First** - Never log, expose, or leak credentials
2. **Defense in Depth** - Multiple layers of protection
3. **Fail Secure** - Reject insecure configurations with clear errors
4. **User Choice** - Support multiple configuration methods
5. **Zero Surprises** - Clear precedence rules and validation

## Configuration Methods

Cicada supports **four configuration methods** with clear precedence:

### Precedence Order (Highest to Lowest)

```
1. Command-line flags         (--zenodo-token, --datacite-password)
2. Environment variables       (CICADA_ZENODO_TOKEN, CICADA_DATACITE_PASSWORD)
3. Cicada config file         (~/.config/cicada/config.yaml)
4. Project .env file          (./.env in current directory)
```

**Rule:** Higher precedence always overrides lower precedence. No merging.

**Example:**
```bash
# .env has ZENODO_TOKEN=old_token
# Command specifies --zenodo-token=new_token
# Result: Uses new_token (command-line wins)
```

## Method 1: Command-Line Flags

### Usage

```bash
cicada doi publish sample.fastq \
  --provider zenodo \
  --zenodo-token "xyzabc123..."
```

```bash
cicada doi publish sample.fastq \
  --provider datacite \
  --datacite-repository-id "10.5072/FK2" \
  --datacite-password "secret123"
```

### Flag Names

| Provider | Flags |
|----------|-------|
| Zenodo | `--zenodo-token` |
| Zenodo Sandbox | `--zenodo-token` + `--zenodo-sandbox` |
| DataCite | `--datacite-repository-id`, `--datacite-password` |
| DataCite Sandbox | Same + `--datacite-sandbox` |

### Security Considerations

**Pros:**
- ✅ Explicit and clear
- ✅ Works for one-off commands
- ✅ No persistent storage

**Cons:**
- ⚠️ **Visible in shell history** (use `history -d` or `HISTCONTROL`)
- ⚠️ **Visible in process list** (`ps aux` shows args)
- ⚠️ **Visible in CI/CD logs** if not masked

**Best Practice:**
```bash
# Read from stdin to avoid history
read -s ZENODO_TOKEN
export ZENODO_TOKEN
cicada doi publish sample.fastq --provider zenodo

# Or use environment variables (Method 2)
```

**Recommendation:** Use command-line flags only for:
- One-off commands
- When other methods aren't available
- Testing/debugging

### Implementation Notes

- Flags are parsed but **never logged**
- Help text shows flag names but not values
- Error messages show `--zenodo-token=***` not actual value

## Method 2: Environment Variables

### Usage

```bash
# Zenodo
export CICADA_ZENODO_TOKEN="xyzabc123..."
cicada doi publish sample.fastq --provider zenodo

# DataCite
export CICADA_DATACITE_REPOSITORY_ID="10.5072/FK2"
export CICADA_DATACITE_PASSWORD="secret123"
cicada doi publish sample.fastq --provider datacite

# Sandbox environments
export CICADA_ZENODO_SANDBOX=true
export CICADA_DATACITE_SANDBOX=true
```

### Variable Names

| Provider | Variables |
|----------|-----------|
| Zenodo Production | `CICADA_ZENODO_TOKEN` |
| Zenodo Sandbox | `CICADA_ZENODO_TOKEN` + `CICADA_ZENODO_SANDBOX=true` |
| DataCite Production | `CICADA_DATACITE_REPOSITORY_ID`, `CICADA_DATACITE_PASSWORD` |
| DataCite Sandbox | Same + `CICADA_DATACITE_SANDBOX=true` |

### Alternative: Provider-Agnostic Names

For compatibility with existing tools:

```bash
# Also supported (without CICADA_ prefix)
export ZENODO_TOKEN="..."
export DATACITE_REPOSITORY_ID="..."
export DATACITE_PASSWORD="..."
```

**Precedence:** `CICADA_*` variables override non-prefixed variables.

```bash
export ZENODO_TOKEN="old_token"
export CICADA_ZENODO_TOKEN="new_token"
# Uses: new_token
```

### Security Considerations

**Pros:**
- ✅ Standard Unix pattern
- ✅ Works with CI/CD secrets
- ✅ Session-scoped (not persistent)
- ✅ Doesn't appear in process args

**Cons:**
- ⚠️ Can leak in error messages if not careful
- ⚠️ Visible to child processes
- ⚠️ Can leak in system logs (`/proc` on Linux)

**Best Practice:**
```bash
# In ~/.bashrc or ~/.zshrc (for persistent use)
export CICADA_ZENODO_TOKEN="xyzabc123..."

# Or in a separate file (source when needed)
# ~/.cicada/env.sh
export CICADA_ZENODO_TOKEN="xyzabc123..."
export CICADA_DATACITE_REPOSITORY_ID="10.5072/FK2"
export CICADA_DATACITE_PASSWORD="secret123"

# Usage:
source ~/.cicada/env.sh
cicada doi publish sample.fastq --provider zenodo
```

**Recommendation:** Use environment variables for:
- CI/CD pipelines (GitHub Actions, GitLab CI)
- Docker containers
- Temporary credentials
- Development environments

### Implementation Notes

- Check both `CICADA_*` and unprefixed variants
- Never log environment variable values
- Warn if credentials found in both prefixed and unprefixed

## Method 3: Cicada Config File

### Location

**Primary:**
```
~/.config/cicada/config.yaml
```

**Fallback (if XDG_CONFIG_HOME not set or Windows):**
```
~/.cicada/config.yaml
```

**Discovery order:**
1. `$XDG_CONFIG_HOME/cicada/config.yaml` (if `XDG_CONFIG_HOME` set)
2. `~/.config/cicada/config.yaml` (Linux/macOS default)
3. `~/.cicada/config.yaml` (fallback)
4. `%APPDATA%\cicada\config.yaml` (Windows)

### Format

```yaml
# ~/.config/cicada/config.yaml

# Provider credentials
providers:
  zenodo:
    token: "xyzabc123..."
    environment: production  # or "sandbox"

  datacite:
    repository_id: "10.5072/FK2"
    password: "secret123"
    environment: production  # or "sandbox"

# Optional: Default provider
default_provider: zenodo

# Optional: Security settings
security:
  check_permissions: true  # Verify file permissions
  warn_insecure: true      # Warn about insecure configurations

# Optional: Other settings (future use)
metadata:
  default_publisher: "My Lab"
  default_rights: "CC-BY-4.0"
```

### Minimal Example

```yaml
providers:
  zenodo:
    token: "xyzabc123..."
```

### Security Considerations

**Pros:**
- ✅ Persistent configuration
- ✅ Structured format (YAML)
- ✅ Can include multiple providers
- ✅ Can be encrypted (future enhancement)
- ✅ Not visible in process list or history

**Cons:**
- ⚠️ **Persistent file - must be secured**
- ⚠️ Can be accidentally committed to git
- ⚠️ Can be copied/backed up insecurely

**Best Practice:**
```bash
# Create config directory
mkdir -p ~/.config/cicada

# Create config file
cat > ~/.config/cicada/config.yaml <<'EOF'
providers:
  zenodo:
    token: "your-token-here"
EOF

# Secure permissions (REQUIRED)
chmod 600 ~/.config/cicada/config.yaml

# Verify permissions
ls -la ~/.config/cicada/config.yaml
# Should show: -rw------- (600)
```

**Recommendation:** Use config file for:
- Persistent credentials (daily use)
- Multiple provider configurations
- Shared settings across commands
- Production use

### Security Requirements

#### 1. File Permissions

**On Unix-like systems (Linux, macOS):**

Cicada **MUST** enforce `600` permissions (owner read/write only):

```bash
# Required permissions
-rw-------  1 user  group  123 Jan 24 10:00 config.yaml
```

**Behavior:**
```go
// Check permissions
fileInfo, _ := os.Stat(configPath)
mode := fileInfo.Mode()

if mode.Perm() & 0077 != 0 {  // Check group/other bits
    return fmt.Errorf("insecure permissions on %s: %o (must be 600)",
                      configPath, mode.Perm())
}
```

**Error message:**
```
Error: Insecure permissions on /home/user/.config/cicada/config.yaml: 0644

  Your config file contains sensitive credentials but has insecure permissions.
  Other users on this system can read your credentials.

  Fix with: chmod 600 /home/user/.config/cicada/config.yaml

  Current: -rw-r--r-- (0644)
  Required: -rw------- (0600)
```

**On Windows:**

Check ACLs to ensure only owner has access:
```go
// Use Windows ACL APIs to verify only owner can read
// Warn if BUILTIN\Users or Everyone has access
```

#### 2. .gitignore Protection

Cicada should **check parent directories** for `.git` and warn:

```
Warning: Config file is in a git repository

  File: /home/user/project/.cicada/config.yaml
  Repo: /home/user/project

  Your config file contains credentials but is in a git repository.
  Ensure it's in .gitignore to prevent accidental commit.

  Add to .gitignore:
    .cicada/config.yaml

  Or move config to: ~/.config/cicada/config.yaml
```

#### 3. World-Readable Check

Warn if config is in shared/world-readable locations:
```
Warning: Config file in potentially shared location

  File: /tmp/cicada/config.yaml

  This location may be accessible to other users.
  Consider moving to: ~/.config/cicada/config.yaml
```

### Implementation Notes

#### Config Loading

```go
func LoadConfig() (*Config, error) {
    // 1. Find config file
    configPath := findConfigFile()

    // 2. Check security
    if err := checkConfigSecurity(configPath); err != nil {
        return nil, err
    }

    // 3. Parse YAML
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, err
    }

    // 4. Unmarshal (but never log contents)
    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("invalid config format: %w", err)
    }

    // 5. Validate (but don't log credential values)
    if err := config.Validate(); err != nil {
        return nil, err
    }

    return &config, nil
}
```

#### Security Checks

```go
func checkConfigSecurity(path string) error {
    // Check 1: File permissions
    if err := checkPermissions(path); err != nil {
        return err
    }

    // Check 2: Inside git repo?
    if inGitRepo(path) {
        // Warning, not error (user might have .gitignore)
        warnGitRepo(path)
    }

    // Check 3: Shared location?
    if isSharedLocation(path) {
        warnSharedLocation(path)
    }

    return nil
}
```

## Method 4: Project .env File

### Location

```
./.env
```

**Discovery:** Current working directory only (not parent directories).

### Format

```bash
# .env (dotenv format)

# Zenodo
CICADA_ZENODO_TOKEN=xyzabc123...

# DataCite
CICADA_DATACITE_REPOSITORY_ID=10.5072/FK2
CICADA_DATACITE_PASSWORD=secret123

# Environment
CICADA_ZENODO_SANDBOX=true
CICADA_DATACITE_SANDBOX=false

# Optional: Default provider
CICADA_DEFAULT_PROVIDER=zenodo
```

### Usage

```bash
# Cicada automatically loads .env if present
cd /path/to/project
cicada doi publish sample.fastq --provider zenodo
# Uses credentials from ./.env
```

### Security Considerations

**Pros:**
- ✅ Project-specific credentials
- ✅ Works with docker-compose patterns
- ✅ Familiar to developers

**Cons:**
- ⚠️ **EASY TO ACCIDENTALLY COMMIT TO GIT**
- ⚠️ Must be in `.gitignore`
- ⚠️ Persistent file - must be secured

**Best Practice:**
```bash
# Create .env
cat > .env <<'EOF'
CICADA_ZENODO_TOKEN=your-token-here
EOF

# Secure permissions
chmod 600 .env

# Add to .gitignore (CRITICAL)
echo ".env" >> .gitignore

# Verify not tracked
git status .env
# Should show: "use git add" (not tracked)
```

**Recommendation:** Use .env for:
- Project-specific credentials (sandbox tokens)
- Development environments
- When using docker-compose
- Team projects (with .env.example template)

### Security Requirements

Same as config file:

1. **File permissions must be 600**
2. **Warn if in git repository without .gitignore**
3. **Check for common mistakes:**

```
Error: .env file has insecure permissions: 0644

  Fix with: chmod 600 .env

Error: .env file is tracked by git

  Your .env file is tracked by git and will be committed.
  This will expose your credentials.

  Fix:
    1. Add to .gitignore:  echo ".env" >> .gitignore
    2. Remove from git:    git rm --cached .env
    3. Commit changes:     git commit -m "Remove .env from tracking"
```

### .env.example Template

Cicada should support creating a template:

```bash
# Generate template (no credentials)
cicada config init --env-example

# Creates .env.example:
CICADA_ZENODO_TOKEN=your-zenodo-token-here
CICADA_DATACITE_REPOSITORY_ID=your-datacite-repo-id
CICADA_DATACITE_PASSWORD=your-datacite-password

# Usage: cp .env.example .env, then edit
```

## Credential Resolution Algorithm

### Pseudocode

```python
def get_credential(provider: str, key: str) -> str:
    """
    Get credential with proper precedence.

    Precedence (highest to lowest):
    1. Command-line flags
    2. Environment variables
    3. Config file
    4. .env file
    """

    # 1. Check command-line flags
    flag_value = get_flag_value(provider, key)
    if flag_value is not None:
        log_source("command-line flag", provider, key)
        return flag_value

    # 2. Check environment variables
    env_value = get_env_value(provider, key)
    if env_value is not None:
        log_source("environment variable", provider, key)
        return env_value

    # 3. Check config file
    config_value = get_config_value(provider, key)
    if config_value is not None:
        log_source("config file", provider, key)
        return config_value

    # 4. Check .env file
    dotenv_value = get_dotenv_value(provider, key)
    if dotenv_value is not None:
        log_source(".env file", provider, key)
        return dotenv_value

    # 5. Not found
    return None

def log_source(source: str, provider: str, key: str):
    """Log where credential came from (but NOT the value)."""
    log.debug(f"Using {provider} {key} from {source}")
```

### Go Implementation

```go
type CredentialSource string

const (
    SourceCommandLine CredentialSource = "command-line flag"
    SourceEnvVar      CredentialSource = "environment variable"
    SourceConfigFile  CredentialSource = "config file"
    SourceDotEnv      CredentialSource = ".env file"
)

type Credential struct {
    Value  string
    Source CredentialSource
}

func GetCredential(provider, key string) (*Credential, error) {
    // 1. Command-line flags
    if val := getFromFlags(provider, key); val != "" {
        return &Credential{Value: val, Source: SourceCommandLine}, nil
    }

    // 2. Environment variables
    if val := getFromEnv(provider, key); val != "" {
        return &Credential{Value: val, Source: SourceEnvVar}, nil
    }

    // 3. Config file
    if val := getFromConfig(provider, key); val != "" {
        return &Credential{Value: val, Source: SourceConfigFile}, nil
    }

    // 4. .env file
    if val := getFromDotEnv(provider, key); val != "" {
        return &Credential{Value: val, Source: SourceDotEnv}, nil
    }

    // 5. Not found
    return nil, fmt.Errorf("%s %s not configured", provider, key)
}
```

## Security Best Practices (Implementation)

### 1. Never Log Credentials

```go
// ❌ BAD
log.Debug("Using token: %s", token)
log.Debug("Config: %+v", config)  // config contains credentials

// ✅ GOOD
log.Debug("Using token from %s", source)
log.Debug("Config loaded from %s", configPath)
```

### 2. Redact in Error Messages

```go
// ❌ BAD
return fmt.Errorf("authentication failed with token %s", token)

// ✅ GOOD
return fmt.Errorf("authentication failed (check token)")

// ✅ GOOD (show source but not value)
return fmt.Errorf("authentication failed (token from %s)", cred.Source)
```

### 3. Redact in Help/Debug Output

```go
// When printing config for debugging:
func (c *Config) String() string {
    return fmt.Sprintf("Config{provider=%s, token=*****, env=%s}",
                       c.Provider, c.Environment)
}
```

### 4. Zero Memory After Use

```go
// For sensitive strings
func clearString(s *string) {
    if s != nil && *s != "" {
        for i := range *s {
            (*s)[i] = 0
        }
    }
}

// Usage:
defer clearString(&token)
```

### 5. Validate Before Use

```go
func ValidateZenodoToken(token string) error {
    if token == "" {
        return errors.New("token is empty")
    }
    if len(token) < 20 {
        return errors.New("token too short (expected 40+ chars)")
    }
    if strings.Contains(token, " ") {
        return errors.New("token contains whitespace")
    }
    return nil
}
```

### 6. Prevent Credential Leaks in Version Control

```bash
# Cicada should create .gitignore entries
cicada config init

# Creates/updates .gitignore:
.env
.cicada/config.yaml
*.credentials
```

## User-Facing Documentation

### Quick Start

```markdown
# Credential Configuration

Cicada needs credentials to publish DOIs. Choose the method that works best for you:

## Option 1: Config File (Recommended)

Create `~/.config/cicada/config.yaml`:

```yaml
providers:
  zenodo:
    token: "your-zenodo-token"
```

Secure the file:
```bash
chmod 600 ~/.config/cicada/config.yaml
```

Use it:
```bash
cicada doi publish sample.fastq --provider zenodo
```

## Option 2: Environment Variables

```bash
export CICADA_ZENODO_TOKEN="your-zenodo-token"
cicada doi publish sample.fastq --provider zenodo
```

## Option 3: Project .env File

Create `.env` in your project:
```bash
CICADA_ZENODO_TOKEN=your-zenodo-token
```

**Important:** Add to `.gitignore`:
```bash
echo ".env" >> .gitignore
```

## Security Checklist

- [ ] Config file has 600 permissions (`chmod 600`)
- [ ] .env file in .gitignore
- [ ] Never commit credentials to git
- [ ] Use sandbox tokens for testing
- [ ] Rotate tokens regularly
```

### Troubleshooting

```markdown
# Common Errors

## "Insecure permissions on config.yaml"

**Cause:** Config file readable by others

**Fix:**
```bash
chmod 600 ~/.config/cicada/config.yaml
```

## ".env file is tracked by git"

**Cause:** .env not in .gitignore

**Fix:**
```bash
echo ".env" >> .gitignore
git rm --cached .env
git commit -m "Remove .env from tracking"
```

## "Authentication failed"

**Cause:** Invalid or expired credentials

**Check:**
1. Verify token is correct (copy/paste carefully)
2. Check token hasn't expired
3. Verify using correct environment (sandbox vs production)
4. Test token directly with provider API

**Debug:**
```bash
cicada doi publish --debug sample.fastq --provider zenodo
# Shows: "Using token from config file"
```
```

## CLI Commands for Credential Management

### Initialize Configuration

```bash
# Create config directory and file
cicada config init

# Creates:
# - ~/.config/cicada/config.yaml (empty template)
# - .gitignore entries (if in git repo)
```

### Set Credentials

```bash
# Interactive prompt (secure input)
cicada config set zenodo-token
# Prompts: Enter Zenodo token: [hidden input]

cicada config set datacite-repository-id
cicada config set datacite-password
```

### View Configuration (Redacted)

```bash
cicada config show

# Output:
Providers:
  zenodo:
    token: ****ab123 (from config file)
    environment: production
  datacite:
    repository_id: 10.5072/FK2 (from environment variable)
    password: ***** (from config file)
    environment: sandbox

Config file: /home/user/.config/cicada/config.yaml
Permissions: 600 ✓
```

### Validate Configuration

```bash
# Check security
cicada config validate

# Output:
✓ Config file permissions: 600
✓ Not in git repository
✓ Zenodo token format valid
✓ DataCite repository ID format valid
⚠ Warning: Using sandbox environment
```

### Test Credentials

```bash
# Test authentication without publishing
cicada config test zenodo
# Output: ✓ Zenodo authentication successful

cicada config test datacite
# Output: ✓ DataCite authentication successful
```

## Implementation Checklist (Issue #26)

### Core Implementation

- [ ] Config file loading (`~/.config/cicada/config.yaml`)
- [ ] Environment variable support (`CICADA_*`)
- [ ] .env file loading (current directory)
- [ ] Command-line flag parsing
- [ ] Credential precedence resolution

### Security Implementation

- [ ] File permission checking (600 required)
- [ ] Git repository detection and warning
- [ ] Credential redaction in logs
- [ ] Credential redaction in error messages
- [ ] Memory zeroing after use
- [ ] Input validation (token format, etc.)

### CLI Commands

- [ ] `cicada config init` - Create config template
- [ ] `cicada config set <key>` - Set credential (secure prompt)
- [ ] `cicada config show` - View config (redacted)
- [ ] `cicada config validate` - Check security
- [ ] `cicada config test <provider>` - Test authentication

### Testing

- [ ] Unit tests for precedence resolution
- [ ] Unit tests for security checks
- [ ] Integration tests with all config methods
- [ ] Test insecure permissions (should fail)
- [ ] Test git repo detection
- [ ] Test credential redaction

### Documentation

- [ ] User guide for configuration
- [ ] Security best practices
- [ ] Troubleshooting guide
- [ ] Examples for each method

## Example: Complete Workflow

### Developer Setup

```bash
# 1. Initialize Cicada config
cicada config init

# 2. Get Zenodo token
# - Go to https://zenodo.org/account/settings/applications/tokens/new/
# - Create token with deposit:write scope
# - Copy token

# 3. Configure Cicada (secure prompt)
cicada config set zenodo-token
# Enter Zenodo token: [paste token, hidden]
# ✓ Token saved to ~/.config/cicada/config.yaml

# 4. Verify setup
cicada config validate
# ✓ Config file permissions: 600
# ✓ Zenodo token format valid

# 5. Test authentication
cicada config test zenodo
# ✓ Zenodo authentication successful

# 6. Publish DOI
cicada doi publish sample.fastq --provider zenodo
# ✓ DOI: 10.5281/zenodo.123456
```

### CI/CD Setup (GitHub Actions)

```yaml
# .github/workflows/publish-doi.yml
name: Publish DOI

on:
  release:
    types: [published]

jobs:
  publish:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Install Cicada
        run: |
          wget https://github.com/scttfrdmn/cicada/releases/download/v0.3.0/cicada_linux_amd64
          chmod +x cicada_linux_amd64
          sudo mv cicada_linux_amd64 /usr/local/bin/cicada

      - name: Publish DOI
        env:
          CICADA_ZENODO_TOKEN: ${{ secrets.ZENODO_TOKEN }}
        run: |
          cicada doi publish data/sample.fastq \
            --provider zenodo \
            --enrich metadata.yaml
```

### Project-Specific Setup

```bash
# 1. Create .env for project (sandbox)
cat > .env <<'EOF'
CICADA_ZENODO_TOKEN=sandbox-token-here
CICADA_ZENODO_SANDBOX=true
EOF

# 2. Secure .env
chmod 600 .env

# 3. Add to .gitignore
echo ".env" >> .gitignore

# 4. Create template for team
cat > .env.example <<'EOF'
# Copy to .env and add your credentials
CICADA_ZENODO_TOKEN=your-token-here
CICADA_ZENODO_SANDBOX=true
EOF

# 5. Commit template (not .env!)
git add .env.example
git commit -m "Add Cicada config template"
```

## Future Enhancements (Post v0.3.0)

### Encrypted Config File

```yaml
# ~/.config/cicada/config.yaml (encrypted)
providers:
  zenodo:
    token: !encrypted |
      AES256:base64-encrypted-data-here
```

Unlock with passphrase or system keychain.

### Keychain Integration

```bash
# Store in system keychain (macOS, GNOME, etc.)
cicada config set zenodo-token --keychain
# Stored in: macOS Keychain / GNOME Keyring / Windows Credential Manager
```

### Credential Helper

```yaml
# ~/.config/cicada/config.yaml
providers:
  zenodo:
    token: !helper "get-zenodo-token.sh"
```

```bash
# get-zenodo-token.sh
#!/bin/bash
# Fetch token from secret manager
aws secretsmanager get-secret-value --secret-id cicada/zenodo-token | jq -r .SecretString
```

### Multi-Profile Support

```bash
# Switch between profiles
cicada config use-profile personal
cicada config use-profile work
cicada config use-profile sandbox
```

---

**Next Steps:**
1. Review this specification
2. Implement in Issue #26 (Provider Configuration System)
3. Write comprehensive tests
4. Document for users
