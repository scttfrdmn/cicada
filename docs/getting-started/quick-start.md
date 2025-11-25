# Quick Start

!!! info "Documentation in Progress"
    This page is being updated as part of v0.3.0. Complete quick start guide coming soon.

## Basic Usage

### Sync to S3

```bash
cicada sync /local/data s3://my-bucket/data
```

### Extract Metadata

```bash
cicada metadata extract sample.fastq.gz
```

### Watch for Changes

```bash
cicada watch add /data/microscope s3://lab-data/microscopy
```

## Next Steps

- [User Guide](../user-guide/overview.md)
- [CLI Reference](../reference/cli.md)

For complete current documentation, see the [README](https://github.com/scttfrdmn/cicada#readme).
