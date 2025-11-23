# Cicada 0.1.0 Release Checklist

**Target**: Minimal viable sync tool - local ↔ S3 with file watching

## Core Features

### Sync Engine (Completed ✅)
- [x] #1 - S3 sync engine with backend abstraction
- [x] Local filesystem backend with MD5 checksums
- [x] AWS S3 backend with SDK v2
- [x] Concurrent transfers with progress reporting
- [x] Dry-run and delete modes
- [x] CLI sync command

### File Watching (Critical for 0.1.0)
- [ ] #2 - Implement file watching system with fsnotify
- [ ] Watch local directories for changes
- [ ] Debounce rapid file changes
- [ ] Auto-sync on file create/modify/delete
- [ ] CLI watch command with start/stop/status
- [ ] Watch configuration persistence

### Configuration Management
- [ ] Configuration file format (YAML/TOML)
- [ ] Config command: init, set, get, list
- [ ] Default config locations (~/.cicada/config.yaml)
- [ ] Per-directory .cicada config support
- [ ] AWS credentials integration
- [ ] Sync path mappings (local → S3)

## CLI Commands

### Essential Commands
- [x] `cicada sync <source> <destination>` - Manual sync
- [ ] `cicada watch <path>` - Start watching directory
- [ ] `cicada watch stop` - Stop watching
- [ ] `cicada watch status` - Show watch status
- [ ] `cicada config init` - Initialize configuration
- [ ] `cicada config set <key> <value>` - Set config value
- [ ] `cicada version` - Show version info

### Flags & Options
- [x] `--dry-run` - Preview changes
- [x] `--delete` - Delete extraneous files
- [x] `--verbose` - Detailed output
- [ ] `--config <file>` - Use specific config file
- [ ] `--profile <name>` - Use AWS profile

## Testing

### Unit Tests
- [x] Sync engine tests (68% coverage)
- [x] Local backend tests
- [x] S3 URI parsing tests
- [ ] File watching tests
- [ ] Configuration tests

### Integration Tests
- [ ] #13 - AWS integration testing setup
- [ ] Local to S3 sync test
- [ ] S3 to local sync test
- [ ] Bidirectional sync test
- [ ] Watch and auto-sync test
- [ ] Large file handling test

### Manual Testing Checklist
- [ ] Fresh install on clean system
- [ ] Config initialization
- [ ] First sync to empty S3 bucket
- [ ] Incremental sync (changed files only)
- [ ] Delete mode (cleanup)
- [ ] Watch mode for 1 hour stability
- [ ] Network interruption recovery
- [ ] AWS credential errors handled gracefully

## Documentation

### User Documentation
- [ ] README.md: Installation instructions
- [ ] README.md: Quick start guide
- [ ] README.md: Basic usage examples
- [ ] README.md: Configuration guide
- [ ] CHANGELOG.md: 0.1.0 release notes

### Developer Documentation
- [x] CLAUDE.md: Project management
- [x] CLAUDE.md: Go coding standards
- [x] CONTRIBUTING.md: Setup and contribution guide
- [x] SECURITY.md: Security reporting
- [x] CODE_OF_CONDUCT.md: Community guidelines

### Examples
- [ ] example-config.yaml: Sample configuration
- [ ] examples/basic-sync.sh: Simple sync script
- [ ] examples/watch-mode.sh: Watch setup script
- [ ] examples/multi-directory.sh: Multiple paths

## Infrastructure

### Build & Release
- [ ] #16 - Configure goreleaser with existing repos
- [ ] Makefile targets: build, install, release
- [ ] Cross-platform builds (Linux, macOS, Windows)
- [ ] Binary release artifacts
- [ ] Homebrew formula (existing tap repo)
- [ ] Scoop manifest (existing bucket repo)
- [ ] deb/rpm packages
- [ ] `cicada update` command (self-update)
- [ ] Installation script (curl | bash)

### Quality Assurance
- [x] #8 - Go linting (A+ grade)
- [x] All linters passing
- [x] #12 - Test coverage standards (80%+)
- [ ] #11 - Minimal CI/CD (GitHub Actions)
- [ ] Automated testing on push
- [ ] Release workflow

## Compliance & Tags

### AWS Resource Management
- [ ] #15 - Implement AWS resource tagging
- [ ] Tag schema definition
- [ ] S3 bucket tagging
- [ ] Installation ID generation
- [ ] Tag documentation

### Security
- [ ] AWS credentials handling best practices
- [ ] No secrets in config files
- [ ] Secure credential storage guidance
- [ ] Security audit of dependencies

## Known Limitations (Document for 0.1.0)

These are acceptable for 0.1.0 but should be documented:

- [ ] No multipart upload for large files (>5GB)
- [ ] No bandwidth throttling
- [ ] No encryption options (uses AWS defaults)
- [ ] No resume capability for interrupted transfers
- [ ] No conflict resolution (last write wins)
- [ ] No version history
- [ ] Single installation only (no multi-tenant)

## Release Criteria

### Must Have (Blocking)
- [ ] All "Core Features" completed
- [ ] All "Essential Commands" implemented
- [ ] Integration tests passing
- [ ] Manual testing checklist completed
- [ ] README.md with installation and quick start
- [ ] CHANGELOG.md with 0.1.0 notes
- [ ] Builds successfully on Linux and macOS
- [ ] No critical bugs

### Should Have (High Priority)
- [ ] 80%+ test coverage
- [ ] All linters passing
- [ ] GitHub Actions CI
- [ ] Cross-platform binaries
- [ ] Installation script

### Nice to Have (Optional)
- [ ] Homebrew formula
- [ ] Windows support
- [ ] AWS resource tagging
- [ ] Performance benchmarks

## Pre-Release Tasks

- [ ] Version bump to 0.1.0 in code
- [ ] Update CHANGELOG.md
- [ ] Create GitHub release with notes
- [ ] Upload binary artifacts
- [ ] Tag release: `v0.1.0`
- [ ] Announce in project discussions

## Post-Release

- [ ] Monitor GitHub issues for bugs
- [ ] Gather user feedback
- [ ] Plan 0.2.0 features
- [ ] Update project board

---

**Target Date**: TBD
**Owner**: scttfrdmn
**Status**: In Progress (35% complete)
