# v0.3.0 Milestone 1: Provider Integration Foundation

**Status:** Not Started
**Target:** 3-4 weeks (February 2025)
**Priority:** CRITICAL

## Overview

Milestone 1 implements live DataCite and Zenodo API integrations, enabling users to mint real DOIs from Cicada. This is the most critical feature for v0.3.0 as it completes the DOI workflow that was prepared in v0.2.0.

## Goals

- [ ] Users can mint DOIs in DataCite sandbox
- [ ] Users can mint DOIs in DataCite production
- [ ] Users can upload files and mint DOIs in Zenodo sandbox
- [ ] Users can upload files and mint DOIs in Zenodo production
- [ ] Error handling covers common failure scenarios
- [ ] Configuration management is straightforward
- [ ] Comprehensive documentation for setup and usage

## Task Breakdown

### Phase 1: Foundation (Week 1)

#### Task 1.1: Provider Configuration System

**Complexity:** LOW
**Estimated Time:** 2 days
**Dependencies:** None

**Description:** Create configuration system for managing provider credentials and settings.

**Implementation:**

1. **Config Structure** (`internal/doi/provider_config.go`):
   ```go
   type ProviderConfig struct {
       Active   string                    // Active provider name
       Providers map[string]ProviderCreds // Provider-specific credentials
   }

   type ProviderCreds struct {
       Type       string            // "datacite", "zenodo", etc.
       Sandbox    bool              // Use sandbox environment
       BaseURL    string            // API base URL
       Credentials map[string]string // Provider-specific credentials
   }
   ```

2. **Configuration Commands**:
   - `cicada config set provider.active zenodo`
   - `cicada config set provider.zenodo.token TOKEN`
   - `cicada config set provider.datacite.repository_id REPO_ID`
   - `cicada config set provider.datacite.username USER`
   - `cicada config set provider.datacite.password PASS`

3. **Security**:
   - Encrypt passwords in config file
   - Support environment variables for secrets
   - Warn if secrets in plaintext

**Deliverables:**
- [ ] `internal/doi/provider_config.go` (config structures)
- [ ] Config get/set methods
- [ ] Environment variable support
- [ ] Unit tests for config management

**Acceptance Criteria:**
- [ ] Can configure multiple providers
- [ ] Credentials are securely stored
- [ ] Environment variables override config file
- [ ] Config validation catches errors

---

#### Task 1.2: DataCite API Client - Core Infrastructure

**Complexity:** MEDIUM
**Estimated Time:** 3 days
**Dependencies:** Task 1.1

**Description:** Implement core DataCite API client with authentication and basic DOI operations.

**Implementation:**

1. **Client Structure** (`internal/doi/datacite_client.go`):
   ```go
   type DataCiteClient struct {
       baseURL      string
       repositoryID string
       username     string
       password     string
       httpClient   *http.Client
   }

   func NewDataCiteClient(config *ProviderCreds) (*DataCiteClient, error)
   func (c *DataCiteClient) CreateDOI(metadata *Dataset) (*DOI, error)
   func (c *DataCiteClient) UpdateDOI(doi string, metadata *Dataset) error
   func (c *DataCiteClient) PublishDOI(doi string) error
   func (c *DataCiteClient) GetDOI(doi string) (*DOI, error)
   func (c *DataCiteClient) ListDOIs(filter *DOIFilter) ([]*DOI, error)
   ```

2. **Authentication**:
   - HTTP Basic Auth (username:password)
   - Set User-Agent header
   - Handle 401/403 responses

3. **API Endpoints**:
   - `POST /dois` - Create draft DOI
   - `PUT /dois/{id}` - Update DOI
   - `PUT /dois/{id}/actions/publish` - Publish DOI
   - `GET /dois/{id}` - Get DOI details
   - `GET /dois` - List DOIs

4. **Request/Response Handling**:
   - JSON serialization/deserialization
   - Error response parsing
   - API version headers

**Deliverables:**
- [ ] `internal/doi/datacite_client.go` (client implementation)
- [ ] `internal/doi/datacite_types.go` (API request/response types)
- [ ] Unit tests with mocked HTTP responses

**Acceptance Criteria:**
- [ ] Client authenticates successfully
- [ ] Can serialize Dataset to DataCite JSON
- [ ] Can deserialize API responses
- [ ] Error handling covers API error responses

---

#### Task 1.3: Zenodo API Client - Core Infrastructure

**Complexity:** MEDIUM
**Estimated Time:** 3 days
**Dependencies:** Task 1.1

**Description:** Implement core Zenodo API client with token authentication and basic operations.

**Implementation:**

1. **Client Structure** (`internal/doi/zenodo_client.go`):
   ```go
   type ZenodoClient struct {
       baseURL    string
       token      string
       httpClient *http.Client
   }

   func NewZenodoClient(config *ProviderCreds) (*ZenodoClient, error)
   func (c *ZenodoClient) CreateDeposition() (*Deposition, error)
   func (c *ZenodoClient) UploadFile(depID string, file io.Reader, filename string) error
   func (c *ZenodoClient) UpdateMetadata(depID string, metadata *Dataset) error
   func (c *ZenodoClient) PublishDeposition(depID string) (*DOI, error)
   func (c *ZenodoClient) GetDeposition(depID string) (*Deposition, error)
   ```

2. **Authentication**:
   - Bearer token authentication
   - Token in query param or header

3. **API Endpoints**:
   - `POST /api/deposit/depositions` - Create deposition
   - `POST /api/deposit/depositions/{id}/files` - Upload file
   - `PUT /api/deposit/depositions/{id}` - Update metadata
   - `POST /api/deposit/depositions/{id}/actions/publish` - Publish
   - `GET /api/deposit/depositions/{id}` - Get deposition

4. **File Upload**:
   - Multipart form data
   - Progress tracking (optional)
   - Checksum validation

**Deliverables:**
- [ ] `internal/doi/zenodo_client.go` (client implementation)
- [ ] `internal/doi/zenodo_types.go` (API request/response types)
- [ ] Unit tests with mocked HTTP responses

**Acceptance Criteria:**
- [ ] Client authenticates with token
- [ ] Can create depositions
- [ ] Can upload files
- [ ] Can update metadata and publish

---

### Phase 2: Integration (Week 2)

#### Task 2.1: Provider Registry Enhancement

**Complexity:** LOW
**Estimated Time:** 1 day
**Dependencies:** Tasks 1.2, 1.3

**Description:** Enhance existing provider registry to support live API clients.

**Implementation:**

1. **Update Provider Interface** (`internal/doi/provider.go`):
   ```go
   type Provider interface {
       Name() string
       Validate(dataset *Dataset) error
       Prepare(dataset *Dataset) (*PrepareResult, error)
       CreateDOI(dataset *Dataset) (*DOI, error)      // NEW
       PublishDOI(doi string) error                    // NEW
       UploadFiles(doi string, files []string) error   // NEW (Zenodo)
       GetDOI(doi string) (*DOI, error)                // NEW
   }
   ```

2. **Update Registry**:
   - Register DataCite provider
   - Register Zenodo provider
   - Provider selection from config

3. **Factory Pattern**:
   ```go
   func (r *ProviderRegistry) GetProvider(name string) (Provider, error)
   func NewDataCiteProvider(config *ProviderCreds) *DataCiteProvider
   func NewZenodoProvider(config *ProviderCreds) *ZenodoProvider
   ```

**Deliverables:**
- [ ] Updated `internal/doi/provider.go`
- [ ] Provider factory implementation
- [ ] Integration tests

**Acceptance Criteria:**
- [ ] Can get DataCite provider from registry
- [ ] Can get Zenodo provider from registry
- [ ] Provider selection respects configuration

---

#### Task 2.2: DataCite Metadata Mapping

**Complexity:** MEDIUM
**Estimated Time:** 2 days
**Dependencies:** Task 2.1

**Description:** Map Cicada Dataset to DataCite JSON format (v4.5 schema).

**Implementation:**

1. **Mapping Function** (`internal/doi/datacite_mapper.go`):
   ```go
   func MapToDataCite(dataset *Dataset) (*DataCiteMetadata, error)
   ```

2. **DataCite JSON Schema**:
   ```json
   {
     "data": {
       "type": "dois",
       "attributes": {
         "prefix": "10.12345",
         "titles": [{"title": "Dataset Title"}],
         "creators": [{"name": "Author Name", "nameType": "Personal"}],
         "publisher": "Publisher Name",
         "publicationYear": 2025,
         "types": {"resourceTypeGeneral": "Dataset"},
         "descriptions": [{"description": "...", "descriptionType": "Abstract"}],
         ...
       }
     }
   }
   ```

3. **Handle All Fields**:
   - Required: title, creators, publisher, publicationYear, resourceType
   - Recommended: subjects, contributors, dates, relatedIdentifiers, descriptions, geoLocations, fundingReferences

4. **Validation**:
   - Ensure all required fields present
   - Validate field formats (ORCID, DOI, etc.)
   - Check constraints (year range, etc.)

**Deliverables:**
- [ ] `internal/doi/datacite_mapper.go`
- [ ] DataCite JSON schema types
- [ ] Mapper unit tests with examples
- [ ] Validation tests

**Acceptance Criteria:**
- [ ] Maps all Dataset fields correctly
- [ ] Produces valid DataCite JSON
- [ ] Handles optional fields gracefully
- [ ] Validation catches common errors

---

#### Task 2.3: Zenodo Metadata Mapping

**Complexity:** MEDIUM
**Estimated Time:** 2 days
**Dependencies:** Task 2.1

**Description:** Map Cicada Dataset to Zenodo JSON format.

**Implementation:**

1. **Mapping Function** (`internal/doi/zenodo_mapper.go`):
   ```go
   func MapToZenodo(dataset *Dataset) (*ZenodoMetadata, error)
   ```

2. **Zenodo JSON Schema**:
   ```json
   {
     "metadata": {
       "title": "Dataset Title",
       "upload_type": "dataset",
       "description": "Dataset description",
       "creators": [{"name": "Author Name", "orcid": "0000-0001-2345-6789"}],
       "keywords": ["keyword1", "keyword2"],
       "access_right": "open",
       "license": "cc-by-4.0",
       ...
     }
   }
   ```

3. **Handle All Fields**:
   - Required: title, upload_type, description, creators
   - Optional: keywords, notes, access_right, license, communities, grants

4. **License Mapping**:
   - Map common license IDs to Zenodo format
   - Default to CC-BY-4.0

**Deliverables:**
- [ ] `internal/doi/zenodo_mapper.go`
- [ ] Zenodo JSON schema types
- [ ] Mapper unit tests
- [ ] License mapping table

**Acceptance Criteria:**
- [ ] Maps all Dataset fields correctly
- [ ] Produces valid Zenodo JSON
- [ ] License mapping works
- [ ] Handles optional fields

---

### Phase 3: CLI Integration (Week 3)

#### Task 3.1: DOI Publish Command

**Complexity:** MEDIUM
**Estimated Time:** 2 days
**Dependencies:** Phase 2

**Description:** Implement `cicada doi publish` command for minting DOIs.

**Implementation:**

1. **Command Structure** (`internal/cli/doi.go`):
   ```bash
   cicada doi publish [files...] [flags]

   Flags:
     --metadata FILE        Metadata file (from doi prepare)
     --provider STRING      Provider: datacite, zenodo (default: from config)
     --files FILE[,FILE]    Files to upload (Zenodo only)
     --landing-page URL     Landing page URL (DataCite)
     --dry-run              Show what would be published
   ```

2. **Workflow**:
   - Load metadata file
   - Select provider from config
   - Call provider.CreateDOI()
   - Upload files (if Zenodo)
   - Call provider.PublishDOI()
   - Display DOI and URL

3. **Output**:
   ```
   Creating DOI...
   Uploading files... (if Zenodo)
   ✓ File 1 uploaded (5.2 MB)
   ✓ File 2 uploaded (3.8 MB)
   Publishing DOI...

   ✓ DOI Published!
   DOI: 10.5281/zenodo.123456
   URL: https://zenodo.org/record/123456
   ```

**Deliverables:**
- [ ] Update `internal/cli/doi.go` with publish command
- [ ] Progress display
- [ ] Error handling
- [ ] Integration tests

**Acceptance Criteria:**
- [ ] Can publish DOI to DataCite
- [ ] Can publish DOI to Zenodo with file upload
- [ ] Progress is displayed
- [ ] Errors are handled gracefully

---

#### Task 3.2: DOI Status Command

**Complexity:** LOW
**Estimated Time:** 1 day
**Dependencies:** Task 3.1

**Description:** Implement `cicada doi status` command for checking DOI state.

**Implementation:**

1. **Command Structure**:
   ```bash
   cicada doi status <doi> [flags]

   Flags:
     --provider STRING    Provider (optional, auto-detect from DOI)
     --format STRING      Output format: table, json, yaml
   ```

2. **Workflow**:
   - Parse DOI
   - Detect provider (from DOI prefix)
   - Call provider.GetDOI()
   - Display status

3. **Output**:
   ```
   DOI: 10.5281/zenodo.123456
   Status: Published
   Title: Dataset Title
   Authors: Jane Smith, John Doe
   Published: 2025-01-15
   URL: https://zenodo.org/record/123456
   ```

**Deliverables:**
- [ ] `cicada doi status` command
- [ ] DOI parsing and provider detection
- [ ] Output formatting

**Acceptance Criteria:**
- [ ] Can check DataCite DOI status
- [ ] Can check Zenodo DOI status
- [ ] Auto-detects provider from DOI

---

#### Task 3.3: DOI List Command

**Complexity:** LOW
**Estimated Time:** 1 day
**Dependencies:** Task 3.1

**Description:** Implement `cicada doi list` command for listing user's DOIs.

**Implementation:**

1. **Command Structure**:
   ```bash
   cicada doi list [flags]

   Flags:
     --provider STRING    Provider: datacite, zenodo, all (default: all)
     --limit INT          Number of results (default: 20)
     --format STRING      Output format: table, json, yaml
   ```

2. **Workflow**:
   - Get provider(s) from config
   - Call provider.ListDOIs()
   - Format output

3. **Output**:
   ```
   DOI                        Title                Status      Published
   ──────────────────────────────────────────────────────────────────────
   10.5281/zenodo.123456     Dataset 1            Published   2025-01-15
   10.12345/dataset-001      Dataset 2            Draft       -
   ```

**Deliverables:**
- [ ] `cicada doi list` command
- [ ] Output formatting
- [ ] Pagination support

**Acceptance Criteria:**
- [ ] Lists DOIs from DataCite
- [ ] Lists DOIs from Zenodo
- [ ] Can list from all providers

---

### Phase 4: Error Handling & Polish (Week 4)

#### Task 4.1: Comprehensive Error Handling

**Complexity:** MEDIUM
**Estimated Time:** 2 days
**Dependencies:** All previous tasks

**Description:** Implement robust error handling for all API operations.

**Implementation:**

1. **Error Types** (`internal/doi/errors.go`):
   ```go
   type APIError struct {
       StatusCode int
       Message    string
       Details    map[string]interface{}
   }

   type AuthenticationError struct{ APIError }
   type ValidationError struct{ APIError }
   type RateLimitError struct{ APIError }
   type NetworkError struct{ error }
   ```

2. **Retry Logic**:
   ```go
   func WithRetry(fn func() error, maxRetries int) error
   ```
   - Exponential backoff
   - Retry on transient failures (5xx, network)
   - Don't retry on client errors (4xx)

3. **User-Friendly Messages**:
   - Parse API error responses
   - Provide actionable suggestions
   - Include error codes

4. **Error Handling Patterns**:
   - Authentication failures: suggest checking credentials
   - Validation errors: show which fields are invalid
   - Rate limiting: suggest wait time
   - Network errors: suggest retry

**Deliverables:**
- [ ] `internal/doi/errors.go` (error types)
- [ ] Retry logic implementation
- [ ] Error message improvements
- [ ] Error handling tests

**Acceptance Criteria:**
- [ ] Transient failures are retried
- [ ] Error messages are helpful
- [ ] Retry logic uses backoff
- [ ] All error types are handled

---

#### Task 4.2: Integration Testing with Sandbox APIs

**Complexity:** MEDIUM
**Estimated Time:** 3 days
**Dependencies:** All previous tasks

**Description:** Create comprehensive integration tests using sandbox APIs.

**Implementation:**

1. **Test Setup** (`internal/integration/provider_test.go`):
   - Sandbox credentials from environment variables
   - Skip tests if credentials not available
   - Clean up resources after tests

2. **DataCite Integration Tests**:
   ```go
   func TestDataCite_EndToEnd(t *testing.T)
   func TestDataCite_CreateDraftDOI(t *testing.T)
   func TestDataCite_PublishDOI(t *testing.T)
   func TestDataCite_UpdateDOI(t *testing.T)
   func TestDataCite_GetDOI(t *testing.T)
   func TestDataCite_AuthenticationFailure(t *testing.T)
   func TestDataCite_ValidationError(t *testing.T)
   ```

3. **Zenodo Integration Tests**:
   ```go
   func TestZenodo_EndToEnd(t *testing.T)
   func TestZenodo_CreateDeposition(t *testing.T)
   func TestZenodo_UploadFile(t *testing.T)
   func TestZenodo_PublishDeposition(t *testing.T)
   func TestZenodo_GetDeposition(t *testing.T)
   func TestZenodo_AuthenticationFailure(t *testing.T)
   ```

4. **CLI Integration Tests**:
   ```go
   func TestCLI_DOIPublish_DataCite(t *testing.T)
   func TestCLI_DOIPublish_Zenodo(t *testing.T)
   func TestCLI_DOIStatus(t *testing.T)
   func TestCLI_DOIList(t *testing.T)
   ```

**Deliverables:**
- [ ] `internal/integration/provider_datacite_test.go`
- [ ] `internal/integration/provider_zenodo_test.go`
- [ ] `internal/integration/cli_doi_test.go`
- [ ] Test fixtures and helpers

**Acceptance Criteria:**
- [ ] All integration tests pass with sandbox APIs
- [ ] Tests clean up resources
- [ ] Tests are reproducible
- [ ] Tests cover success and error scenarios

---

#### Task 4.3: Documentation

**Complexity:** LOW
**Estimated Time:** 2 days
**Dependencies:** All previous tasks

**Description:** Create comprehensive documentation for provider integration.

**Documentation Files:**

1. **Provider Setup Guide** (update `docs/PROVIDERS.md`):
   - DataCite sandbox setup (step-by-step)
   - DataCite production setup
   - Zenodo sandbox setup
   - Zenodo production setup
   - Configuration examples
   - Troubleshooting

2. **DOI Publishing Guide** (new `docs/DOI_PUBLISHING.md`):
   - End-to-end publishing workflow
   - Provider comparison (when to use each)
   - Command examples
   - Best practices
   - Common issues

3. **API Reference** (update):
   - Provider interface documentation
   - Client method documentation
   - Error types documentation

**Deliverables:**
- [ ] Updated `docs/PROVIDERS.md`
- [ ] New `docs/DOI_PUBLISHING.md`
- [ ] API reference documentation
- [ ] Example workflows

**Acceptance Criteria:**
- [ ] User can set up DataCite from documentation
- [ ] User can set up Zenodo from documentation
- [ ] All commands are documented with examples
- [ ] Troubleshooting covers common issues

---

## Testing Strategy

### Unit Tests

- [ ] Provider configuration (Task 1.1)
- [ ] DataCite client (mocked HTTP) (Task 1.2)
- [ ] Zenodo client (mocked HTTP) (Task 1.3)
- [ ] Metadata mappers (Tasks 2.2, 2.3)
- [ ] Error handling (Task 4.1)

**Target Coverage:** > 80%

### Integration Tests

- [ ] DataCite sandbox API (Task 4.2)
- [ ] Zenodo sandbox API (Task 4.2)
- [ ] CLI commands (Task 4.2)

**Environment Variables Required:**
- `DATACITE_SANDBOX_REPOSITORY_ID`
- `DATACITE_SANDBOX_USERNAME`
- `DATACITE_SANDBOX_PASSWORD`
- `ZENODO_SANDBOX_TOKEN`

### Manual Tests

- [ ] Test with production DataCite (if available)
- [ ] Test with production Zenodo
- [ ] Test large file upload to Zenodo (> 100 MB)
- [ ] Test error scenarios (bad credentials, network failure)

## Dependencies & Prerequisites

### Required Accounts

1. **DataCite Sandbox** (free):
   - Sign up at https://support.datacite.org/docs/testing-guide
   - Get repository ID, username, password

2. **Zenodo Sandbox** (free):
   - Sign up at https://sandbox.zenodo.org
   - Generate API token

3. **DataCite Production** (institutional):
   - Requires institutional membership
   - Contact library/IT for credentials

4. **Zenodo Production** (free):
   - Sign up at https://zenodo.org
   - Generate API token

### External Libraries

- `net/http` - HTTP client (standard library)
- `encoding/json` - JSON handling (standard library)
- Consider: `hashicorp/go-retryablehttp` for retry logic

### Documentation References

- DataCite REST API: https://support.datacite.org/docs/api
- DataCite Metadata Schema v4.5: https://schema.datacite.org/meta/kernel-4.5/
- Zenodo API: https://developers.zenodo.org/
- Zenodo Depositions: https://developers.zenodo.org/#depositions

## Success Criteria

### Functional

- [ ] Can create draft DOI in DataCite sandbox
- [ ] Can publish DOI in DataCite sandbox
- [ ] Can create deposition in Zenodo sandbox
- [ ] Can upload files to Zenodo sandbox
- [ ] Can publish DOI in Zenodo sandbox
- [ ] Can retrieve DOI status from both providers
- [ ] Can list DOIs from both providers

### Quality

- [ ] Unit test coverage > 80%
- [ ] All integration tests pass
- [ ] Error messages are helpful
- [ ] Documentation is complete

### User Experience

- [ ] Setup takes < 10 minutes with documentation
- [ ] Publishing DOI takes < 30 seconds (excluding upload time)
- [ ] Errors provide clear next steps
- [ ] Progress is visible during long operations

## Risk Mitigation

### Risk 1: API Changes

**Risk:** Provider APIs might change during development
**Probability:** LOW
**Impact:** MEDIUM
**Mitigation:**
- Use API versioning in URLs
- Monitor provider changelog
- Integration tests catch breaking changes

### Risk 2: Authentication Issues

**Risk:** Authentication might be more complex than expected
**Probability:** MEDIUM
**Impact:** MEDIUM
**Mitigation:**
- Start with sandbox environment
- Comprehensive error messages
- Fallback to manual token entry

### Risk 3: File Upload Performance

**Risk:** Large file uploads to Zenodo might be slow/unreliable
**Probability:** MEDIUM
**Impact:** LOW
**Mitigation:**
- Implement progress tracking
- Add resume capability
- Document file size limits

### Risk 4: Rate Limiting

**Risk:** APIs might have rate limits we hit during testing
**Probability:** LOW
**Impact:** LOW
**Mitigation:**
- Implement backoff logic
- Cache API responses where possible
- Document rate limits

## Timeline Summary

| Week | Phase | Key Tasks | Deliverables |
|------|-------|-----------|--------------|
| 1 | Foundation | Config, DataCite client, Zenodo client | API clients with unit tests |
| 2 | Integration | Provider registry, metadata mapping | Working provider implementations |
| 3 | CLI | Publish, status, list commands | CLI commands with integration |
| 4 | Polish | Error handling, integration tests, docs | Production-ready milestone |

**Total Duration:** 4 weeks
**Buffer:** None (aggressive timeline)
**Contingency:** Week 5 if needed

## Next Steps

1. **Set up sandbox accounts** (Day 1)
2. **Implement Task 1.1** (Provider configuration) (Days 1-2)
3. **Implement Task 1.2** (DataCite client) (Days 3-5)
4. **Implement Task 1.3** (Zenodo client) (Days 6-8)
5. **Continue with Phase 2** (Week 2)

**Decision Point:** End of Week 2 - Evaluate progress and adjust timeline if needed.
