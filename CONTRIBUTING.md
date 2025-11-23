# Contributing to Cicada

Thanks for your interest in contributing to Cicada!

## Getting Started

### Prerequisites

- **Go 1.21 or later**
- **Make** (for build automation)
- **golangci-lint** (for linting)
  ```bash
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  ```
- **AWS CLI** (for integration tests, optional)

### Building from Source

```bash
# Clone the repository
git clone https://github.com/scttfrdmn/cicada.git
cd cicada

# Download dependencies
make deps

# Build the binary
make build

# Verify it works
./bin/cicada --version
```

### Running Tests

```bash
# Run unit tests
make test

# Run with coverage
make test-cover

# Run linting
make lint

# Run all checks
make check
```

## Development Workflow

### 1. Find or Create an Issue

Browse open issues or create a new one:
```bash
gh issue list
gh issue create
```

All work should be tied to a GitHub issue.

### 2. Create a Branch

```bash
git checkout -b feature/issue-123-sync-engine
```

Use descriptive branch names: `feature/`, `fix/`, `docs/`, `refactor/`

### 3. Make Your Changes

Write clear, idiomatic Go code:
- Follow standard Go conventions
- Write tests alongside implementation
- Add documentation for exported functions
- Update CHANGELOG.md if user-facing

### 4. Test Your Changes

Before committing:
```bash
# Format code
make fmt

# Run linters
make lint

# Run tests
make test

# Check coverage (aim for 80%+)
make test-coverage
```

All checks must pass.

### 5. Commit Your Changes

Use [Conventional Commits](https://www.conventionalcommits.org/) format:

```bash
git commit -m "feat: implement S3 sync engine (#123)

- Add multipart upload support
- Implement checksum-based delta detection
- Add progress reporting"
```

**Types**:
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation
- `refactor:` - Code restructuring
- `test:` - Test changes
- `chore:` - Build/tooling changes

### 6. Push and Create Pull Request

```bash
git push origin feature/issue-123-sync-engine

gh pr create --title "Implement S3 sync engine" \
  --body "Closes #123"
```

## Code Standards

### Go Report Card A+

Cicada maintains a **Go Report Card grade of A+**. All code must pass:
- gofmt - Code formatting
- go vet - Suspicious constructs
- gocyclo - Cyclomatic complexity <15
- golint/revive - Style checks
- ineffassign - Dead code detection
- misspell - Spelling
- errcheck - Unchecked errors

Run `make lint` to check locally.

### Test Coverage

Target **80%+ test coverage** overall, but be reasonable:

**Must test**:
- Business logic
- Error handling
- Edge cases
- Public APIs

**May skip**:
- Simple getters/setters
- Main/CLI wiring
- Pure I/O operations

**Never skip**:
- Security-sensitive code
- Compliance features
- Financial calculations

### Code Style

- Use `gofmt` (run `make fmt`)
- Keep functions under 100 lines
- Keep cyclomatic complexity under 15
- Write clear, descriptive names
- Add comments for exported identifiers
- Use table-driven tests

## Pull Request Guidelines

### PR Checklist

- [ ] Linked to issue (e.g., "Closes #123")
- [ ] Tests added or updated
- [ ] All tests pass (`make test`)
- [ ] Linting passes (`make lint`)
- [ ] Documentation updated
- [ ] CHANGELOG.md updated (if user-facing)
- [ ] Commit messages follow conventional format

### PR Review Process

1. Automated checks run (build, lint, test)
2. Maintainer reviews code
3. Address feedback
4. Merge when approved

## Integration Testing (Optional)

For integration tests with real AWS:

1. Copy `.env.test.example` to `.env.test`
2. Add your AWS credentials
3. Run `make test-integration-setup`
4. Run `make test-integration`

See [Issue #13](https://github.com/scttfrdmn/cicada/issues/13) for details.

## Project Structure

```
cicada/
├── cmd/cicada/          # CLI entry point
├── internal/            # Private packages
│   ├── sync/           # Sync engine
│   ├── metadata/       # Metadata system
│   ├── daemon/         # Background service
│   └── testutil/       # Test helpers
├── docs/               # User documentation
├── examples/           # Configuration examples
└── planning/           # Design documents (not tracked)
```

## Getting Help

- **Questions**: Open a [GitHub Discussion](https://github.com/scttfrdmn/cicada/discussions)
- **Bugs**: Create an [issue](https://github.com/scttfrdmn/cicada/issues/new?template=bug_report.yml)
- **Features**: Create an [issue](https://github.com/scttfrdmn/cicada/issues/new?template=feature_request.yml)
- **Security**: See [SECURITY.md](SECURITY.md)

## License

By contributing, you agree that your contributions will be licensed under the Apache License 2.0.
