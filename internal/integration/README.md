# Cicada Integration Tests

This directory contains integration tests that interact with real external services:
- **AWS S3** - Cloud storage backend tests
- **DOI Providers** - DataCite and Zenodo sandbox API tests

## Table of Contents

1. [S3 Integration Tests](#s3-integration-tests)
2. [DOI Provider Integration Tests](#doi-provider-integration-tests)
3. [Running Tests](#running-integration-tests)
4. [CI/CD Integration](#cicd-integration)

---

## S3 Integration Tests

### Prerequisites

#### 1. AWS Account and Credentials

You need an AWS account with S3 access. Configure your credentials:

```bash
# Using AWS CLI
aws configure --profile aws

# Or manually edit ~/.aws/credentials
[aws]
aws_access_key_id = YOUR_ACCESS_KEY
aws_secret_access_key = YOUR_SECRET_KEY
```

#### 2. S3 Test Bucket

Create a dedicated S3 bucket for integration testing:

```bash
aws s3 mb s3://cicada-integration-test --region us-west-2 --profile aws
```

**Important**: The tests expect:
- **Profile**: `aws` (configured in `~/.aws/credentials`)
- **Region**: `us-west-2`
- **Bucket**: `cicada-integration-test`

You can modify these constants in `s3_test.go` if needed:
```go
const (
    testProfile = "aws"
    testRegion  = "us-west-2"
    testBucket  = "cicada-integration-test"
)
```

#### 3. IAM Permissions

The AWS profile needs the following S3 permissions on the test bucket:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:ListBucket",
        "s3:GetObject",
        "s3:PutObject",
        "s3:DeleteObject",
        "s3:HeadBucket",
        "s3:PutObjectTagging",
        "s3:GetObjectTagging"
      ],
      "Resource": [
        "arn:aws:s3:::cicada-integration-test",
        "arn:aws:s3:::cicada-integration-test/*"
      ]
    }
  ]
}
```

**Note:** The `s3:PutObjectTagging` and `s3:GetObjectTagging` permissions are required for S3 object tagging integration tests. These permissions allow Cicada to:
- Write metadata as S3 object tags during upload (`s3:PutObjectTagging`)
- Read metadata from S3 object tags (`s3:GetObjectTagging`)
- Tag existing objects with metadata (`s3:PutObjectTagging`)

### Test Coverage

#### S3 Backend Operations
- ✅ Write: Upload files to S3
- ✅ Read: Download files from S3
- ✅ Stat: Get file metadata (size, ETag, ModTime)
- ✅ List: List objects with prefix
- ✅ Delete: Remove objects

#### Sync Operations
- ✅ Local → S3: Upload local directory to S3
- ✅ S3 → Local: Download S3 prefix to local directory
- ✅ Nested directories: Handle subdirectories correctly
- ✅ Content verification: Ensure data integrity

#### Test Isolation

Each test run:
1. Creates a unique prefix with timestamp: `test-run-<unix-timestamp>/`
2. Performs operations within that prefix
3. Cleans up all objects at the end (via defer)

This ensures tests can run concurrently without conflicts.

#### Costs

Integration tests create minimal AWS costs:
- Small test files (< 1KB each)
- Short-lived objects (cleaned up immediately)
- Minimal API calls

Expected cost: < $0.01 per test run.

#### Troubleshooting

##### "Test bucket does not exist"

```
Failed to load AWS config: operation error S3: HeadBucket...
```

**Solution**: Create the test bucket:
```bash
aws s3 mb s3://cicada-integration-test --region us-west-2 --profile aws
```

##### "Failed to load AWS config"

```
Failed to load AWS config: no EC2 IMDS role found
```

**Solution**: Configure AWS credentials:
```bash
aws configure --profile aws
```

##### "Access Denied"

```
operation error S3: PutObject, https response error StatusCode: 403
```

**Solution**: Verify your IAM permissions include S3 write access.

##### Tests Hanging

If tests hang, they may be waiting for AWS responses. Check:
1. Network connectivity
2. AWS region configuration
3. S3 bucket region matches test region (us-west-2)

## CI/CD Integration

To run integration tests in CI/CD:

```yaml
# .github/workflows/integration.yml
name: Integration Tests
on: [push, pull_request]

jobs:
  integration:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.23'

      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v2
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: us-west-2

      - name: Run integration tests
        run: go test -v -tags=integration ./internal/integration/...
```

**Required secrets**:
- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`

## Local Development

For local development, you can use LocalStack to mock AWS services:

```bash
# Start LocalStack
docker run -d -p 4566:4566 localstack/localstack

# Set endpoint in tests
export AWS_ENDPOINT_URL=http://localhost:4566
```

Then modify `s3_test.go` to use the endpoint override.

---

## DOI Provider Integration Tests

### Overview

The DOI provider integration tests verify that Cicada correctly interacts with DataCite and Zenodo sandbox APIs. These tests validate:
- DOI minting and creation
- DOI retrieval
- DOI metadata updates
- Listing depositions
- Error handling (authentication, validation)

### Prerequisites

#### DataCite Sandbox Credentials

1. **Obtain Sandbox Credentials**:
   - Contact your institution's DataCite administrator
   - Request sandbox access (separate from production)
   - You'll receive: Repository ID, Username, Password

2. **Set Environment Variables**:
   ```bash
   export DATACITE_SANDBOX_REPOSITORY_ID="YOUR_REPO_ID"
   export DATACITE_SANDBOX_PASSWORD="YOUR_PASSWORD"
   ```

3. **Verify Access**:
   ```bash
   curl -u "$DATACITE_SANDBOX_REPOSITORY_ID:$DATACITE_SANDBOX_PASSWORD" \
     https://api.test.datacite.org/dois | jq '.'
   ```

#### Zenodo Sandbox Token

1. **Create Sandbox Account**:
   - Visit: https://sandbox.zenodo.org
   - Sign up with GitHub, ORCID, or email

2. **Generate API Token**:
   - Settings → Applications → Personal access tokens
   - New token → Name: "Cicada Testing"
   - Scopes: `deposit:write`, `deposit:actions`
   - Copy token (shown once)

3. **Set Environment Variable**:
   ```bash
   export ZENODO_SANDBOX_TOKEN="your_token_here"
   ```

4. **Verify Access**:
   ```bash
   curl "https://sandbox.zenodo.org/api/deposit/depositions" \
     -H "Authorization: Bearer $ZENODO_SANDBOX_TOKEN" | jq '.'
   ```

### Running DOI Integration Tests

#### Run All DOI Tests

```bash
# DataCite tests only
go test -v -tags=integration ./internal/integration/provider_datacite_test.go

# Zenodo tests only
go test -v -tags=integration ./internal/integration/provider_zenodo_test.go

# All DOI tests
go test -v -tags=integration ./internal/integration/provider_*_test.go
```

#### Run Specific DOI Tests

```bash
# DataCite end-to-end workflow
go test -v -tags=integration -run TestDataCite_EndToEnd ./internal/integration/

# Zenodo end-to-end workflow
go test -v -tags=integration -run TestZenodo_EndToEnd ./internal/integration/

# Authentication failure tests (no credentials needed)
go test -v -tags=integration -run TestDataCite_AuthenticationFailure ./internal/integration/
go test -v -tags=integration -run TestZenodo_AuthenticationFailure ./internal/integration/
```

### Test Coverage

#### DataCite Tests (`provider_datacite_test.go`)

- ✅ **TestDataCite_EndToEnd**: Complete workflow (Mint → Get → Update → List)
- ✅ **TestDataCite_CreateDraftDOI**: Create draft DOI
- ✅ **TestDataCite_GetDOI**: Retrieve existing DOI
- ✅ **TestDataCite_AuthenticationFailure**: Invalid credentials handling
- ✅ **TestDataCite_ValidationError**: Invalid metadata handling
- ✅ **TestDataCite_ListDOIs**: List all DOIs

#### Zenodo Tests (`provider_zenodo_test.go`)

- ✅ **TestZenodo_EndToEnd**: Complete workflow (Create → Get → List)
- ✅ **TestZenodo_CreateDeposition**: Create new deposition
- ✅ **TestZenodo_GetDeposition**: Retrieve existing deposition
- ✅ **TestZenodo_AuthenticationFailure**: Invalid token handling
- ✅ **TestZenodo_ListDepositions**: List all depositions
- ✅ **TestZenodo_InvalidMetadata**: Validation error handling
- ✅ **TestZenodo_GetNonExistentDOI**: 404 error handling

### Test Behavior

#### Without Credentials

Tests automatically skip when credentials are not available:

```bash
$ go test -v -tags=integration ./internal/integration/provider_datacite_test.go
=== RUN   TestDataCite_EndToEnd
    provider_datacite_test.go:35: DataCite sandbox credentials not available (set DATACITE_SANDBOX_REPOSITORY_ID and DATACITE_SANDBOX_PASSWORD)
--- SKIP: TestDataCite_EndToEnd (0.00s)
```

#### With Credentials

Tests run against sandbox APIs and create/cleanup test DOIs:

```bash
$ go test -v -tags=integration ./internal/integration/provider_zenodo_test.go
=== RUN   TestZenodo_EndToEnd
    provider_zenodo_test.go:54: Creating new deposition...
    provider_zenodo_test.go:67: Successfully created deposition with DOI: 10.5072/zenodo.123456 (state: draft)
    provider_zenodo_test.go:70: Retrieving deposition...
    provider_zenodo_test.go:80: Successfully retrieved deposition: 10.5072/zenodo.123456
    provider_zenodo_test.go:83: Listing depositions...
    provider_zenodo_test.go:104: Successfully listed 1 depositions
--- PASS: TestZenodo_EndToEnd (2.34s)
```

### Test Data

Each test creates unique test datasets with timestamps to avoid conflicts:

```go
// Example test dataset
{
  "title": "Integration Test Dataset - e2e-1732476234",
  "authors": [{"name": "Test Researcher", "affiliation": "Test University"}],
  "publisher": "Cicada Integration Tests",
  "publication_year": 2025,
  "resource_type": "Dataset",
  "keywords": ["integration testing", "automated testing"]
}
```

### Costs

DOI integration tests create minimal costs:
- **DataCite Sandbox**: Free (test DOIs don't count against quotas)
- **Zenodo Sandbox**: Free (unlimited test depositions)
- **Network**: < 1 MB data transfer per test run

### Troubleshooting

#### "DataCite sandbox credentials not available"

**Solution**: Set environment variables:
```bash
export DATACITE_SANDBOX_REPOSITORY_ID="YOUR_REPO_ID"
export DATACITE_SANDBOX_PASSWORD="YOUR_PASSWORD"
```

#### "Zenodo sandbox token not available"

**Solution**: Set environment variable:
```bash
export ZENODO_SANDBOX_TOKEN="your_token_here"
```

#### "API request failed with status 401"

**Cause**: Invalid credentials

**DataCite Solution**:
```bash
# Verify credentials with your institution
# Test manually:
curl -u "$DATACITE_SANDBOX_REPOSITORY_ID:$DATACITE_SANDBOX_PASSWORD" \
  https://api.test.datacite.org/dois
```

**Zenodo Solution**:
```bash
# Regenerate token at https://sandbox.zenodo.org
# Test manually:
curl "https://sandbox.zenodo.org/api/deposit/depositions" \
  -H "Authorization: Bearer $ZENODO_SANDBOX_TOKEN"
```

#### "API request failed with status 422"

**Cause**: Validation error (usually missing required fields)

**Solution**: Check test metadata structure matches provider requirements

---

## Running Integration Tests

### Run All Integration Tests

```bash
# From project root - runs ALL integration tests (S3 + DOI)
go test -v -tags=integration ./internal/integration/...
```

### Run by Category

```bash
# S3 tests only
go test -v -tags=integration ./internal/integration/ -run TestS3

# DOI tests only
go test -v -tags=integration ./internal/integration/ -run Test.*Cite|Test.*Zenodo

# DataCite tests only
go test -v -tags=integration ./internal/integration/ -run TestDataCite

# Zenodo tests only
go test -v -tags=integration ./internal/integration/ -run TestZenodo
```

### Skip Integration Tests

Integration tests are automatically skipped when:

1. Running without the `integration` build tag:
   ```bash
   go test ./...  # Skips ALL integration tests
   ```

2. Running in short mode:
   ```bash
   go test -short -tags=integration ./internal/integration/
   ```

3. Missing required credentials (DOI tests only - they skip gracefully)
