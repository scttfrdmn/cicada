# Cicada Development Guide

This document provides guidance for development work on the Cicada project.

## Project Management Philosophy

**Cicada uses GitHub's native project management tools instead of extensive markdown planning documents.**

### Why GitHub Issues & Projects?

- **Living documentation**: Issues stay current with development progress
- **Collaborative**: Easy for contributors to see what needs work
- **Traceable**: Links commits, PRs, and discussions directly to tasks
- **No duplication**: Single source of truth, not scattered across multiple markdown files
- **Built-in project boards**: Visual progress tracking with milestones

## Project Structure

### Planning Documents (Private)

The `/planning` directory contains initial project specifications and is **not tracked in git**:
- Original design conversations
- Detailed technical specifications
- Implementation roadmap

These documents were used to bootstrap the project but are **not maintained**. The GitHub project is the living source of truth.

### Public Documentation (Repository)

Public-facing documentation that is maintained:
- `README.md` - Project overview
- `COMPLIANCE-CROSSWALK.md` - Security control mapping
- `COMPLIANCE-ESSENTIALS.md` - Compliance guidance
- `docs/` - User and API documentation
- `examples/` - Configuration examples
- `internal/` - Go code specifications (reference implementations)

## Development Workflow

### 1. Check GitHub Milestones

View project phases and their deadlines:
```bash
gh milestone list
```

Current phases:
- **Phase 1**: Core Storage & Sync (Weeks 1-6)
- **Phase 2**: Metadata & FAIR (Weeks 7-10)
- **Phase 3**: Web UI & User Management (Weeks 11-14)
- **Phase 4**: Compute & Workflows (Weeks 15-18)
- **Phase 5**: Workstations & Portal (Weeks 19-22)
- **Phase 6**: Compliance & Polish (Weeks 23-26)

### 2. Browse Issues by Priority

```bash
# Critical issues
gh issue list --label "priority: critical"

# Good first issues for new contributors
gh issue list --label "good first issue"

# Issues by area
gh issue list --label "area: sync"
gh issue list --label "area: metadata"
gh issue list --label "area: web"
```

### 3. Work on an Issue

1. **Assign yourself**: Prevents duplicate work
   ```bash
   gh issue develop <issue-number> --checkout
   ```

2. **Create feature branch**: Use descriptive names
   ```bash
   git checkout -b feature/issue-123-sync-engine
   ```

3. **Reference issue in commits**: Links commits to issues
   ```bash
   git commit -m "feat: implement S3 sync engine (#123)

   - Add multipart upload support
   - Implement checksum-based delta detection
   - Add progress reporting"
   ```

4. **Create pull request**: Auto-links to issue
   ```bash
   gh pr create --title "Implement S3 sync engine" --body "Closes #123"
   ```

### 4. Update Issue Status

Issues automatically move through states:
- **Open**: Ready to work on
- **In Progress**: Assigned and being worked on (manual label)
- **Review**: PR submitted (automatic when PR created)
- **Closed**: Merged or resolved

## Label System

### Priority Labels
- `priority: critical` - Must have for MVP
- `priority: high` - Important feature
- `priority: medium` - Nice to have
- `priority: low` - Future enhancement

### Type Labels
- `type: feature` - New functionality
- `type: bug` - Bug fix
- `type: docs` - Documentation
- `type: refactor` - Code improvement
- `type: test` - Testing

### Area Labels
- `area: sync` - Sync engine and file operations
- `area: metadata` - Metadata system
- `area: web` - Web UI
- `area: compute` - Compute and workflows
- `area: compliance` - Security and compliance
- `area: cli` - Command-line interface
- `area: daemon` - Background service

### Status Labels
- `status: blocked` - Waiting on dependencies
- `status: in-progress` - Actively being worked on
- `status: review` - In code review
- `good first issue` - Good for new contributors

### Adding New Labels

**Labels should be added to the project as needed**, but always check if an existing label can be used first to avoid label sprawl.

**Before creating a new label**:
1. Check existing labels: `gh label list`
2. Consider if an existing label applies (e.g., use `type: refactor` instead of creating `code-cleanup`)
3. Only create a new label if there's a clear, recurring need

**Common additional label categories**:
- `quality: performance` - Performance optimization
- `quality: security` - Security hardening
- `quality: tech-debt` - Technical debt
- `dependencies` - Dependency updates
- `breaking-change` - Breaking API changes

**How to create a new label**:
```bash
gh label create "quality: performance" --description "Performance optimization needed" --color "fef2c0"
```

**Label naming conventions**:
- Use kebab-case: `tech-debt` not `Tech Debt`
- Use prefixes for categories: `area:`, `priority:`, `type:`, `quality:`
- Be descriptive but concise

## Go Coding Standards

**Write idiomatic Go code that passes all linters.** Cicada maintains a Go Report Card grade of A+.

### Key Idioms to Follow

1. **Always handle errors explicitly**
   ```go
   // ❌ BAD: Silently ignoring errors
   defer file.Close()

   // ✅ GOOD: Explicitly ignoring errors
   defer func() { _ = file.Close() }()
   ```

2. **Use short variable names in limited scopes**
   ```go
   // ✅ GOOD: Short names in small scopes
   for i, f := range files {
       if err := processFile(f); err != nil {
           return fmt.Errorf("file %d: %w", i, err)
       }
   }
   ```

3. **Return early, avoid else**
   ```go
   // ❌ BAD: Unnecessary else
   if err != nil {
       return err
   } else {
       doWork()
   }

   // ✅ GOOD: Early return
   if err != nil {
       return err
   }
   doWork()
   ```

4. **Wrap errors with context**
   ```go
   // ❌ BAD: Losing error context
   return err

   // ✅ GOOD: Adding context
   return fmt.Errorf("sync file %s: %w", path, err)
   ```

5. **Use defer for cleanup**
   ```go
   f, err := os.Open(path)
   if err != nil {
       return err
   }
   defer func() { _ = f.Close() }()

   // Work with file...
   ```

### Linting Requirements

Before committing, ensure code passes all linters:

```bash
make lint    # Must show "0 issues"
```

Common linter errors to avoid:
- **errcheck**: All errors must be checked or explicitly ignored
- **gocyclo**: Keep cyclomatic complexity < 15
- **gosec**: No security vulnerabilities
- **govet**: No suspicious constructs
- **misspell**: Check spelling in comments

### Test Coverage

- Target 80%+ overall coverage
- Be reasonable: don't test getters, generated code, or pure I/O
- Write table-driven tests for complex logic
- Use testify/assert for clearer test assertions

```bash
make test-cover    # Generate coverage report
```

## Creating New Issues

When you identify new work, create an issue instead of editing markdown:

```bash
gh issue create \
  --title "Short descriptive title" \
  --body "Detailed description with requirements and acceptance criteria" \
  --milestone "Phase X: Name" \
  --label "type: feature" \
  --label "area: sync" \
  --label "priority: medium"
```

### Good Issue Template

```markdown
## Context
Brief explanation of what needs to be done and why.

## Requirements
- [ ] Specific requirement 1
- [ ] Specific requirement 2
- [ ] Specific requirement 3

## Acceptance Criteria
Clear, testable criteria for completion.

## Technical Notes
- Implementation hints
- Related files/modules
- Dependencies on other issues

## References
- Links to docs, specs, or related issues
```

## Finding Work

### For Core Contributors

1. Check current milestone: `gh issue list --milestone "Phase X"`
2. Filter by priority: Start with `priority: critical`
3. Pick issues in your area of expertise

### For New Contributors

1. Start with: `gh issue list --label "good first issue"`
2. Read `CONTRIBUTING.md` for setup instructions
3. Ask questions in issue comments before starting work

### For Specific Features

1. Search issues: `gh issue list --search "keyword"`
2. Check related area label: `gh issue list --label "area: sync"`
3. Review milestone description for context

## Tracking Progress

### Project Board

View visual progress on GitHub:
```
https://github.com/scttfrdmn/cicada/projects
```

### Milestone Progress

```bash
# View specific milestone
gh api repos/scttfrdmn/cicada/milestones/1 | jq '{
  title: .title,
  open_issues: .open_issues,
  closed_issues: .closed_issues,
  progress: (.closed_issues / (.open_issues + .closed_issues) * 100)
}'
```

### Your Assigned Issues

```bash
gh issue list --assignee @me
```

## Avoiding Markdown Planning Hell

❌ **Don't**:
- Create extensive markdown TODO lists
- Maintain separate planning documents
- Track progress in multiple places
- Write long specification docs that get outdated

✅ **Do**:
- Create GitHub issues for all work
- Use issue comments for discussion
- Link related issues and PRs
- Update issue descriptions as requirements evolve
- Close issues when work is complete

## Architecture and Design Decisions

For major architectural decisions, use **GitHub Discussions** or **Architecture Decision Records (ADRs)** in the repository:

```
docs/adr/
  001-use-aws-s3-as-backend.md
  002-metadata-stored-in-files.md
  003-go-for-implementation.md
```

This keeps architectural context discoverable and version-controlled without becoming stale planning documents.

## Questions?

- **General questions**: Open a GitHub Discussion
- **Bug reports**: Create an issue with `type: bug`
- **Feature requests**: Create an issue with `type: feature`
- **Security issues**: See `SECURITY.md`

---

**Remember**: GitHub issues are the single source of truth for what needs to be done. Planning documents are reference material only and are not maintained.
