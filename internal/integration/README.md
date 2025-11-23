# Cicada Integration Tests

This directory contains integration tests that interact with real AWS services.

## Prerequisites

### 1. AWS Account and Credentials

You need an AWS account with S3 access. Configure your credentials:

```bash
# Using AWS CLI
aws configure --profile aws

# Or manually edit ~/.aws/credentials
[aws]
aws_access_key_id = YOUR_ACCESS_KEY
aws_secret_access_key = YOUR_SECRET_KEY
```

### 2. S3 Test Bucket

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

### 3. IAM Permissions

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
        "s3:HeadBucket"
      ],
      "Resource": [
        "arn:aws:s3:::cicada-integration-test",
        "arn:aws:s3:::cicada-integration-test/*"
      ]
    }
  ]
}
```

## Running Integration Tests

### Run All Integration Tests

```bash
# From project root
go test -v -tags=integration ./internal/integration/...
```

### Run Specific Test

```bash
go test -v -tags=integration ./internal/integration/ -run TestS3Backend_Integration
go test -v -tags=integration ./internal/integration/ -run TestLocalToS3Sync_Integration
go test -v -tags=integration ./internal/integration/ -run TestS3ToLocalSync_Integration
```

### Skip Integration Tests

Integration tests are automatically skipped when:

1. Running without the `integration` build tag:
   ```bash
   go test ./...  # Skips integration tests
   ```

2. Running in short mode:
   ```bash
   go test -short -tags=integration ./internal/integration/
   ```

## Test Coverage

The integration tests cover:

### S3 Backend Operations
- ✅ Write: Upload files to S3
- ✅ Read: Download files from S3
- ✅ Stat: Get file metadata (size, ETag, ModTime)
- ✅ List: List objects with prefix
- ✅ Delete: Remove objects

### Sync Operations
- ✅ Local → S3: Upload local directory to S3
- ✅ S3 → Local: Download S3 prefix to local directory
- ✅ Nested directories: Handle subdirectories correctly
- ✅ Content verification: Ensure data integrity

## Test Isolation

Each test run:
1. Creates a unique prefix with timestamp: `test-run-<unix-timestamp>/`
2. Performs operations within that prefix
3. Cleans up all objects at the end (via defer)

This ensures tests can run concurrently without conflicts.

## Costs

Integration tests create minimal AWS costs:
- Small test files (< 1KB each)
- Short-lived objects (cleaned up immediately)
- Minimal API calls

Expected cost: < $0.01 per test run.

## Troubleshooting

### "Test bucket does not exist"

```
Failed to load AWS config: operation error S3: HeadBucket...
```

**Solution**: Create the test bucket:
```bash
aws s3 mb s3://cicada-integration-test --region us-west-2 --profile aws
```

### "Failed to load AWS config"

```
Failed to load AWS config: no EC2 IMDS role found
```

**Solution**: Configure AWS credentials:
```bash
aws configure --profile aws
```

### "Access Denied"

```
operation error S3: PutObject, https response error StatusCode: 403
```

**Solution**: Verify your IAM permissions include S3 write access.

### Tests Hanging

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
