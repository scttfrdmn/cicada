#!/bin/bash

# Create v0.3.0 GitHub issues for Milestone 1: Provider Integration Foundation
# Usage: ./scripts/create_v0.3.0_issues.sh

set -e

echo "Creating v0.3.0 GitHub issues..."
echo ""

# Milestone 1: Provider Integration Foundation (4 weeks)

echo "Creating Issue: Provider Configuration System..."
gh issue create \
  --title "[v0.3.0] Provider Configuration System" \
  --label "type: feature,area: doi,priority: critical,milestone: v0.3.0" \
  --body "Implement configuration system for managing provider credentials and settings.

**Estimate**: 2 days | **Phase**: 1 - Foundation (Week 1)

**Implementation**:
- Config structure for provider credentials
- Secure credential storage (encryption)
- Environment variable support
- Configuration validation
- CLI commands: \`cicada config set provider.*\`

**Tasks**:
- [ ] Create \`internal/doi/provider_config.go\`
- [ ] Define \`ProviderConfig\` and \`ProviderCreds\` structs
- [ ] Implement config get/set methods
- [ ] Add environment variable override support
- [ ] Implement credential encryption
- [ ] Add config validation
- [ ] Write unit tests for config management

**Acceptance Criteria**:
- [ ] Can configure multiple providers (DataCite, Zenodo)
- [ ] Credentials are securely stored
- [ ] Environment variables override config file
- [ ] Config validation catches common errors
- [ ] Unit test coverage > 80%

**Dependencies**: None

**Reference**: docs/ROADMAP_v0.3.0_MILESTONE1.md - Task 1.1"

echo "Creating Issue: DataCite API Client..."
gh issue create \
  --title "[v0.3.0] DataCite API Client - Core Infrastructure" \
  --label "type: feature,area: doi,priority: critical,milestone: v0.3.0" \
  --body "Implement core DataCite API client with authentication and basic DOI operations.

**Estimate**: 3 days | **Phase**: 1 - Foundation (Week 1)

**Implementation**:
- HTTP Basic Auth
- DOI CRUD operations (Create, Update, Publish, Get, List)
- JSON serialization/deserialization
- Error response parsing
- API version headers

**Tasks**:
- [ ] Create \`internal/doi/datacite_client.go\`
- [ ] Implement \`DataCiteClient\` struct with HTTP client
- [ ] Implement \`CreateDOI()\` method
- [ ] Implement \`UpdateDOI()\` method
- [ ] Implement \`PublishDOI()\` method
- [ ] Implement \`GetDOI()\` method
- [ ] Implement \`ListDOIs()\` method
- [ ] Create \`internal/doi/datacite_types.go\` (request/response types)
- [ ] Write unit tests with mocked HTTP responses
- [ ] Document all methods

**API Endpoints**:
- \`POST /dois\` - Create draft DOI
- \`PUT /dois/{id}\` - Update DOI
- \`PUT /dois/{id}/actions/publish\` - Publish DOI
- \`GET /dois/{id}\` - Get DOI details
- \`GET /dois\` - List DOIs

**Acceptance Criteria**:
- [ ] Client authenticates successfully with sandbox API
- [ ] Can serialize Dataset to DataCite JSON
- [ ] Can deserialize API responses
- [ ] Error handling covers 4xx/5xx responses
- [ ] Unit test coverage > 80%

**Dependencies**: Issue #1 (Provider Configuration)

**Reference**: docs/ROADMAP_v0.3.0_MILESTONE1.md - Task 1.2
**API Docs**: https://support.datacite.org/docs/api"

echo "Creating Issue: Zenodo API Client..."
gh issue create \
  --title "[v0.3.0] Zenodo API Client - Core Infrastructure" \
  --label "type: feature,area: doi,priority: critical,milestone: v0.3.0" \
  --body "Implement core Zenodo API client with token authentication and file upload.

**Estimate**: 3 days | **Phase**: 1 - Foundation (Week 1)

**Implementation**:
- Bearer token authentication
- Deposition operations (Create, Update, Publish, Get)
- File upload with multipart form data
- Progress tracking (optional)
- Checksum validation

**Tasks**:
- [ ] Create \`internal/doi/zenodo_client.go\`
- [ ] Implement \`ZenodoClient\` struct with HTTP client
- [ ] Implement \`CreateDeposition()\` method
- [ ] Implement \`UploadFile()\` method with progress tracking
- [ ] Implement \`UpdateMetadata()\` method
- [ ] Implement \`PublishDeposition()\` method
- [ ] Implement \`GetDeposition()\` method
- [ ] Create \`internal/doi/zenodo_types.go\` (request/response types)
- [ ] Write unit tests with mocked HTTP responses
- [ ] Document all methods

**API Endpoints**:
- \`POST /api/deposit/depositions\` - Create deposition
- \`POST /api/deposit/depositions/{id}/files\` - Upload file
- \`PUT /api/deposit/depositions/{id}\` - Update metadata
- \`POST /api/deposit/depositions/{id}/actions/publish\` - Publish
- \`GET /api/deposit/depositions/{id}\` - Get deposition

**Acceptance Criteria**:
- [ ] Client authenticates with token
- [ ] Can create depositions
- [ ] Can upload files with progress tracking
- [ ] Can update metadata and publish
- [ ] Error handling covers API errors
- [ ] Unit test coverage > 80%

**Dependencies**: Issue #1 (Provider Configuration)

**Reference**: docs/ROADMAP_v0.3.0_MILESTONE1.md - Task 1.3
**API Docs**: https://developers.zenodo.org/"

echo "Creating Issue: Provider Registry Enhancement..."
gh issue create \
  --title "[v0.3.0] Provider Registry Enhancement for Live APIs" \
  --label "type: feature,area: doi,priority: high,milestone: v0.3.0" \
  --body "Enhance existing provider registry to support live API clients.

**Estimate**: 1 day | **Phase**: 2 - Integration (Week 2)

**Implementation**:
- Update Provider interface with new methods
- Provider factory pattern
- Provider selection from config
- Integration with DataCite and Zenodo clients

**Tasks**:
- [ ] Update \`internal/doi/provider.go\` interface
- [ ] Add \`CreateDOI()\` method to interface
- [ ] Add \`PublishDOI()\` method to interface
- [ ] Add \`UploadFiles()\` method to interface (Zenodo)
- [ ] Add \`GetDOI()\` method to interface
- [ ] Implement provider factory: \`NewDataCiteProvider()\`
- [ ] Implement provider factory: \`NewZenodoProvider()\`
- [ ] Update registry to support live providers
- [ ] Write integration tests

**Acceptance Criteria**:
- [ ] Can get DataCite provider from registry
- [ ] Can get Zenodo provider from registry
- [ ] Provider selection respects configuration
- [ ] Interface is backward compatible
- [ ] Integration tests pass

**Dependencies**: Issues #2 (DataCite Client), #3 (Zenodo Client)

**Reference**: docs/ROADMAP_v0.3.0_MILESTONE1.md - Task 2.1"

echo "Creating Issue: DataCite Metadata Mapping..."
gh issue create \
  --title "[v0.3.0] DataCite Metadata Mapping to Schema v4.5" \
  --label "type: feature,area: doi,priority: high,milestone: v0.3.0" \
  --body "Map Cicada Dataset to DataCite JSON format (Schema v4.5).

**Estimate**: 2 days | **Phase**: 2 - Integration (Week 2)

**Implementation**:
- Complete mapping from Dataset to DataCite JSON
- Handle all required and recommended fields
- Field validation (ORCID, DOI formats, etc.)
- Constraint checking (year ranges, etc.)

**Tasks**:
- [ ] Create \`internal/doi/datacite_mapper.go\`
- [ ] Implement \`MapToDataCite()\` function
- [ ] Define DataCite JSON schema types
- [ ] Map required fields: title, creators, publisher, publicationYear, resourceType
- [ ] Map recommended fields: subjects, contributors, dates, relatedIdentifiers, descriptions
- [ ] Implement field validation
- [ ] Write mapper unit tests with examples
- [ ] Write validation tests
- [ ] Document mapping rules

**DataCite Required Fields**:
- titles, creators, publisher, publicationYear, resourceType

**DataCite Recommended Fields**:
- subjects, contributors, dates, relatedIdentifiers, descriptions, geoLocations, fundingReferences

**Acceptance Criteria**:
- [ ] Maps all Dataset fields correctly
- [ ] Produces valid DataCite JSON Schema v4.5
- [ ] Handles optional fields gracefully
- [ ] Validation catches common errors (invalid ORCID, etc.)
- [ ] Unit test coverage > 90%

**Dependencies**: Issue #4 (Provider Registry)

**Reference**: docs/ROADMAP_v0.3.0_MILESTONE1.md - Task 2.2
**Schema**: https://schema.datacite.org/meta/kernel-4.5/"

echo "Creating Issue: Zenodo Metadata Mapping..."
gh issue create \
  --title "[v0.3.0] Zenodo Metadata Mapping" \
  --label "type: feature,area: doi,priority: high,milestone: v0.3.0" \
  --body "Map Cicada Dataset to Zenodo JSON format.

**Estimate**: 2 days | **Phase**: 2 - Integration (Week 2)

**Implementation**:
- Complete mapping from Dataset to Zenodo JSON
- License ID mapping
- Handle all required and optional fields
- Upload type classification

**Tasks**:
- [ ] Create \`internal/doi/zenodo_mapper.go\`
- [ ] Implement \`MapToZenodo()\` function
- [ ] Define Zenodo JSON schema types
- [ ] Map required fields: title, upload_type, description, creators
- [ ] Map optional fields: keywords, notes, access_right, license, communities, grants
- [ ] Implement license ID mapping table
- [ ] Write mapper unit tests
- [ ] Document mapping rules

**Zenodo Required Fields**:
- title, upload_type, description, creators

**Zenodo Optional Fields**:
- keywords, notes, access_right, license, communities, grants, related_identifiers

**Acceptance Criteria**:
- [ ] Maps all Dataset fields correctly
- [ ] Produces valid Zenodo JSON
- [ ] License mapping works for common licenses
- [ ] Handles optional fields gracefully
- [ ] Unit test coverage > 90%

**Dependencies**: Issue #4 (Provider Registry)

**Reference**: docs/ROADMAP_v0.3.0_MILESTONE1.md - Task 2.3
**API Docs**: https://developers.zenodo.org/"

echo "Creating Issue: DOI Publish Command..."
gh issue create \
  --title "[v0.3.0] CLI: cicada doi publish Command" \
  --label "type: feature,area: cli,priority: critical,milestone: v0.3.0" \
  --body "Implement \`cicada doi publish\` command for minting DOIs.

**Estimate**: 2 days | **Phase**: 3 - CLI Integration (Week 3)

**Command Structure**:
\`\`\`bash
cicada doi publish [files...] [flags]

Flags:
  --metadata FILE        Metadata file (from doi prepare)
  --provider STRING      Provider: datacite, zenodo (default: from config)
  --files FILE[,FILE]    Files to upload (Zenodo only)
  --landing-page URL     Landing page URL (DataCite)
  --dry-run              Show what would be published
\`\`\`

**Workflow**:
1. Load metadata file
2. Select provider from config
3. Call provider.CreateDOI()
4. Upload files (if Zenodo)
5. Call provider.PublishDOI()
6. Display DOI and URL

**Tasks**:
- [ ] Update \`internal/cli/doi.go\` with publish command
- [ ] Implement metadata file loading
- [ ] Implement provider selection logic
- [ ] Implement file upload with progress display
- [ ] Implement DOI creation workflow
- [ ] Implement dry-run mode
- [ ] Add progress indicators
- [ ] Implement error handling
- [ ] Write integration tests
- [ ] Document command usage

**Output Example**:
\`\`\`
Creating DOI...
Uploading files...
✓ File 1 uploaded (5.2 MB)
✓ File 2 uploaded (3.8 MB)
Publishing DOI...

✓ DOI Published!
DOI: 10.5281/zenodo.123456
URL: https://zenodo.org/record/123456
\`\`\`

**Acceptance Criteria**:
- [ ] Can publish DOI to DataCite sandbox
- [ ] Can publish DOI to Zenodo sandbox with file upload
- [ ] Progress is displayed during upload
- [ ] Errors are handled gracefully
- [ ] Dry-run mode works correctly
- [ ] Integration tests pass

**Dependencies**: Issues #5 (DataCite Mapping), #6 (Zenodo Mapping)

**Reference**: docs/ROADMAP_v0.3.0_MILESTONE1.md - Task 3.1"

echo "Creating Issue: DOI Status Command..."
gh issue create \
  --title "[v0.3.0] CLI: cicada doi status Command" \
  --label "type: feature,area: cli,priority: medium,milestone: v0.3.0" \
  --body "Implement \`cicada doi status\` command for checking DOI state.

**Estimate**: 1 day | **Phase**: 3 - CLI Integration (Week 3)

**Command Structure**:
\`\`\`bash
cicada doi status <doi> [flags]

Flags:
  --provider STRING    Provider (optional, auto-detect from DOI)
  --format STRING      Output format: table, json, yaml
\`\`\`

**Workflow**:
1. Parse DOI
2. Detect provider (from DOI prefix)
3. Call provider.GetDOI()
4. Display status

**Tasks**:
- [ ] Implement \`cicada doi status\` command
- [ ] Implement DOI parsing
- [ ] Implement provider auto-detection from DOI prefix
- [ ] Implement status display (table, json, yaml)
- [ ] Write unit tests
- [ ] Write integration tests
- [ ] Document command usage

**Output Example**:
\`\`\`
DOI: 10.5281/zenodo.123456
Status: Published
Title: Dataset Title
Authors: Jane Smith, John Doe
Published: 2025-01-15
URL: https://zenodo.org/record/123456
\`\`\`

**Acceptance Criteria**:
- [ ] Can check DataCite DOI status
- [ ] Can check Zenodo DOI status
- [ ] Auto-detects provider from DOI prefix
- [ ] All output formats work
- [ ] Integration tests pass

**Dependencies**: Issue #7 (DOI Publish Command)

**Reference**: docs/ROADMAP_v0.3.0_MILESTONE1.md - Task 3.2"

echo "Creating Issue: DOI List Command..."
gh issue create \
  --title "[v0.3.0] CLI: cicada doi list Command" \
  --label "type: feature,area: cli,priority: low,milestone: v0.3.0" \
  --body "Implement \`cicada doi list\` command for listing user's DOIs.

**Estimate**: 1 day | **Phase**: 3 - CLI Integration (Week 3)

**Command Structure**:
\`\`\`bash
cicada doi list [flags]

Flags:
  --provider STRING    Provider: datacite, zenodo, all (default: all)
  --limit INT          Number of results (default: 20)
  --format STRING      Output format: table, json, yaml
\`\`\`

**Workflow**:
1. Get provider(s) from config
2. Call provider.ListDOIs()
3. Format output

**Tasks**:
- [ ] Implement \`cicada doi list\` command
- [ ] Implement provider selection (all, datacite, zenodo)
- [ ] Implement output formatting (table, json, yaml)
- [ ] Implement pagination support
- [ ] Write unit tests
- [ ] Write integration tests
- [ ] Document command usage

**Output Example**:
\`\`\`
DOI                        Title                Status      Published
──────────────────────────────────────────────────────────────────────
10.5281/zenodo.123456     Dataset 1            Published   2025-01-15
10.12345/dataset-001      Dataset 2            Draft       -
\`\`\`

**Acceptance Criteria**:
- [ ] Lists DOIs from DataCite
- [ ] Lists DOIs from Zenodo
- [ ] Can list from all providers
- [ ] Pagination works
- [ ] All output formats work
- [ ] Integration tests pass

**Dependencies**: Issue #7 (DOI Publish Command)

**Reference**: docs/ROADMAP_v0.3.0_MILESTONE1.md - Task 3.3"

echo "Creating Issue: Error Handling and Retry Logic..."
gh issue create \
  --title "[v0.3.0] Comprehensive Error Handling and Retry Logic" \
  --label "type: feature,area: doi,priority: high,milestone: v0.3.0" \
  --body "Implement robust error handling for all API operations.

**Estimate**: 2 days | **Phase**: 4 - Polish (Week 4)

**Implementation**:
- Typed error system
- Retry logic with exponential backoff
- User-friendly error messages
- Error handling patterns

**Tasks**:
- [ ] Create \`internal/doi/errors.go\`
- [ ] Define error types: \`APIError\`, \`AuthenticationError\`, \`ValidationError\`, \`RateLimitError\`, \`NetworkError\`
- [ ] Implement \`WithRetry()\` function with exponential backoff
- [ ] Implement retry on transient failures (5xx, network errors)
- [ ] Don't retry on client errors (4xx)
- [ ] Add user-friendly error messages with suggestions
- [ ] Parse API error responses
- [ ] Write error handling tests
- [ ] Document error types

**Error Handling Patterns**:
- Authentication failures → suggest checking credentials
- Validation errors → show which fields are invalid
- Rate limiting → suggest wait time
- Network errors → suggest retry

**Acceptance Criteria**:
- [ ] Transient failures are retried automatically
- [ ] Error messages are helpful and actionable
- [ ] Retry logic uses exponential backoff
- [ ] All error types are handled
- [ ] Unit test coverage > 80%

**Dependencies**: All previous implementation tasks

**Reference**: docs/ROADMAP_v0.3.0_MILESTONE1.md - Task 4.1"

echo "Creating Issue: Integration Tests with Sandbox APIs..."
gh issue create \
  --title "[v0.3.0] Integration Tests with Sandbox APIs" \
  --label "type: test,area: doi,priority: critical,milestone: v0.3.0" \
  --body "Create comprehensive integration tests using sandbox APIs.

**Estimate**: 3 days | **Phase**: 4 - Polish (Week 4)

**Test Coverage**:
- DataCite sandbox API (end-to-end workflow)
- Zenodo sandbox API (end-to-end workflow)
- CLI commands integration
- Error scenarios (auth failures, validation errors)

**Tasks**:
- [ ] Set up test environment with sandbox credentials
- [ ] Create \`internal/integration/provider_datacite_test.go\`
- [ ] Create \`internal/integration/provider_zenodo_test.go\`
- [ ] Create \`internal/integration/cli_doi_test.go\`
- [ ] Implement DataCite integration tests (create, update, publish, get)
- [ ] Implement Zenodo integration tests (create, upload, publish, get)
- [ ] Implement CLI integration tests
- [ ] Implement error scenario tests
- [ ] Add test fixtures and helpers
- [ ] Document how to run integration tests

**DataCite Tests**:
- \`TestDataCite_EndToEnd\`
- \`TestDataCite_CreateDraftDOI\`
- \`TestDataCite_PublishDOI\`
- \`TestDataCite_UpdateDOI\`
- \`TestDataCite_GetDOI\`
- \`TestDataCite_AuthenticationFailure\`
- \`TestDataCite_ValidationError\`

**Zenodo Tests**:
- \`TestZenodo_EndToEnd\`
- \`TestZenodo_CreateDeposition\`
- \`TestZenodo_UploadFile\`
- \`TestZenodo_PublishDeposition\`
- \`TestZenodo_GetDeposition\`
- \`TestZenodo_AuthenticationFailure\`

**CLI Tests**:
- \`TestCLI_DOIPublish_DataCite\`
- \`TestCLI_DOIPublish_Zenodo\`
- \`TestCLI_DOIStatus\`
- \`TestCLI_DOIList\`

**Environment Variables Required**:
- \`DATACITE_SANDBOX_REPOSITORY_ID\`
- \`DATACITE_SANDBOX_USERNAME\`
- \`DATACITE_SANDBOX_PASSWORD\`
- \`ZENODO_SANDBOX_TOKEN\`

**Acceptance Criteria**:
- [ ] All integration tests pass with sandbox APIs
- [ ] Tests clean up resources after completion
- [ ] Tests are reproducible
- [ ] Tests cover success and error scenarios
- [ ] Tests skip gracefully if credentials not available
- [ ] Integration test coverage > 80%

**Dependencies**: All previous tasks

**Reference**: docs/ROADMAP_v0.3.0_MILESTONE1.md - Task 4.2"

echo "Creating Issue: Provider Integration Documentation..."
gh issue create \
  --title "[v0.3.0] Provider Integration Documentation" \
  --label "type: documentation,area: doi,priority: high,milestone: v0.3.0" \
  --body "Create comprehensive documentation for provider integration.

**Estimate**: 2 days | **Phase**: 4 - Polish (Week 4)

**Documentation Files**:
1. Update \`docs/PROVIDERS.md\` (provider setup guide)
2. Create \`docs/DOI_PUBLISHING.md\` (publishing workflow guide)
3. Update API reference documentation

**Tasks**:
- [ ] Update \`docs/PROVIDERS.md\` with:
  - [ ] DataCite sandbox setup (step-by-step)
  - [ ] DataCite production setup
  - [ ] Zenodo sandbox setup
  - [ ] Zenodo production setup
  - [ ] Configuration examples
  - [ ] Troubleshooting guide
- [ ] Create \`docs/DOI_PUBLISHING.md\` with:
  - [ ] End-to-end publishing workflow
  - [ ] Provider comparison (when to use each)
  - [ ] Command examples
  - [ ] Best practices
  - [ ] Common issues and solutions
- [ ] Update API reference:
  - [ ] Provider interface documentation
  - [ ] Client method documentation
  - [ ] Error types documentation
- [ ] Add example workflows
- [ ] Add troubleshooting section

**Acceptance Criteria**:
- [ ] User can set up DataCite from documentation (< 10 minutes)
- [ ] User can set up Zenodo from documentation (< 5 minutes)
- [ ] All commands are documented with examples
- [ ] Troubleshooting covers common issues
- [ ] Documentation is reviewed for clarity

**Dependencies**: All previous tasks

**Reference**: docs/ROADMAP_v0.3.0_MILESTONE1.md - Task 4.3"

echo ""
echo "✅ Created 12 issues for v0.3.0 Milestone 1: Provider Integration Foundation"
echo ""
echo "Next steps:"
echo "1. Review issues at: https://github.com/scttfrdmn/cicada/issues"
echo "2. Set up GitHub project board for v0.3.0"
echo "3. Create milestone 'v0.3.0: Provider Integration' if not exists"
echo "4. Assign issues to milestone"
echo "5. Begin development with Issue #1"
