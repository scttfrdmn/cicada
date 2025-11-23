# Changelog

All notable changes to Cicada will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2025-11-23

Initial release of Cicada - Fast, reliable file sync for S3.

### Added

- **Sync Engine**: rsync-like sync between local filesystem and S3
  - Concurrent transfers with configurable concurrency
  - MD5/ETag comparison for change detection
  - Dry-run mode for previewing changes
  - Delete mode for mirror syncing
  - Exclude patterns support

- **File Watching**: Auto-sync directories on file changes
  - fsnotify integration for real-time monitoring
  - Debouncing to prevent sync storms
  - Min-age filtering for partial file writes
  - Recursive directory watching
  - Watch persistence via configuration

- **Configuration Management**: YAML-based configuration
  - AWS profile and region configuration
  - Default sync options
  - Watch persistence
  - Global settings
  - CLI commands: init, set, get, list

- **S3 Backend**: Full S3 integration
  - List, read, write, stat, delete operations
  - AWS SDK v2 integration
  - Multi-region support
  - Profile-based authentication

- **Local Backend**: Filesystem operations
  - MD5 checksum calculation
  - Recursive directory listing
  - Efficient file I/O

- **CLI Commands**:
  - `cicada sync`: Manual sync operations
  - `cicada watch`: File watching management
  - `cicada config`: Configuration management
  - `cicada version`: Version information

- **Testing**:
  - Comprehensive unit tests (80%+ coverage)
  - AWS S3 integration tests
  - Makefile targets for local testing

- **Documentation**:
  - Comprehensive README with usage examples
  - Integration test setup guide
  - Development instructions
  - Contributing guidelines

- **CI/CD**:
  - GitHub Actions workflow for build and test
  - goreleaser configuration for releases
  - Homebrew tap integration
  - Scoop bucket integration

### Known Limitations

- No multipart upload support (files >5GB may be slow)
- No bandwidth throttling
- Sync engine path handling needs improvement for some use cases
- No background daemon (runs in foreground)
- No resume capability for interrupted transfers

[Unreleased]: https://github.com/scttfrdmn/cicada/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/scttfrdmn/cicada/releases/tag/v0.1.0
