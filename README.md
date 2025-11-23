# Cicada: Dormant Data Commons for Academic Research

## Overview

Cicada is a lightweight, cost-effective data commons platform designed for academic research labs (8-10 people) with limited technical expertise and tight budgets. Like its namesake, Cicada lies dormant most of the time, consuming minimal resources, but emerges powerfully when needed.

### Key Principles

- **Dormant by Design**: Resources spin up on-demand, shut down automatically
- **Cost-Conscious**: Optimized for $50-100/month budgets
- **FAIR by Default**: Findable, Accessible, Interoperable, Reusable
- **Zero AWS Knowledge Required**: Abstracts cloud complexity completely
- **Domain-Flexible**: Supports custom metadata schemas for any research domain

### Target Users

- Small academic labs (8-10 people)
- Non-technical researchers (biologists, chemists, physicists, etc.)
- Limited IT support/budget
- Need for: data protection, collaboration, computational workflows, data sharing

### Core Features

1. **Intelligent Data Sync**: rsync-like engine with automatic tiering
2. **File Watching**: Auto-upload from instruments/workstations
3. **Metadata Management**: Flexible, domain-specific schemas
4. **Workflow Execution**: Snakemake, Nextflow, CWL support
5. **Remote Workstations**: On-demand GPU instances for visualization
6. **Data Sharing**: Public portal with DOI minting
7. **Compliance**: NIST 800-171, HIPAA, GDPR support
8. **User Management**: Simple IAM with Globus Auth integration

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    User Interfaces                      │
│  ┌──────────┐  ┌──────────┐  ┌──────────────────────┐ │
│  │   CLI    │  │  Web UI  │  │  Public Data Portal  │ │
│  └────┬─────┘  └────┬─────┘  └──────────┬───────────┘ │
└───────┼─────────────┼────────────────────┼─────────────┘
        │             │                    │
        └─────────────┼────────────────────┘
                      ▼
┌─────────────────────────────────────────────────────────┐
│                  Cicada Daemon (Local)                  │
│  ┌──────────────────────────────────────────────────┐  │
│  │  Core Services                                    │  │
│  │  • Sync Engine        • Metadata Manager          │  │
│  │  • File Watcher       • Schema Validator          │  │
│  │  • Workflow Orchestrator                          │  │
│  │  • Workstation Manager                            │  │
│  └──────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────┐  │
│  │  HTTP Server (Port 7878)                          │  │
│  │  • REST API           • WebSocket (real-time)     │  │
│  │  • Static Files       • SSE (progress streams)    │  │
│  └──────────────────────────────────────────────────┘  │
└────────────────────────┬────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────┐
│                    AWS Infrastructure                   │
│  ┌─────────────┐  ┌──────────────┐  ┌───────────────┐ │
│  │  S3 Storage │  │  AWS Batch   │  │  EC2/Fargate  │ │
│  │  • Intelligent│  │  • Spot      │  │  • Workstations│ │
│  │    Tiering   │  │    Instances │  │  • Gateway    │ │
│  │  • Versioning│  │  • Auto-scale│  │  • On-demand  │ │
│  └─────────────┘  └──────────────┘  └───────────────┘ │
│  ┌─────────────┐  ┌──────────────┐  ┌───────────────┐ │
│  │  IAM/SSO    │  │  CloudWatch  │  │  Lambda       │ │
│  │  • Policies │  │  • Logs      │  │  • Triggers   │ │
│  │  • Roles    │  │  • Alerts    │  │  • Cleanup    │ │
│  └─────────────┘  └──────────────┘  └───────────────┘ │
└─────────────────────────────────────────────────────────┘
```

## Project Structure

```
cicada/
├── cmd/
│   ├── cicada/              # Main CLI entry point
│   ├── cicada-daemon/       # Background daemon service
│   └── cicada-gateway/      # File gateway orchestrator
│
├── internal/
│   ├── sync/               # rsync-like sync engine
│   │   ├── engine.go
│   │   ├── delta.go
│   │   ├── checksum.go
│   │   └── transfer.go
│   │
│   ├── watch/              # File system watcher
│   │   ├── watcher.go
│   │   ├── debounce.go
│   │   └── patterns.go
│   │
│   ├── storage/            # S3 operations
│   │   ├── s3.go
│   │   ├── multipart.go
│   │   ├── lifecycle.go
│   │   └── versioning.go
│   │
│   ├── metadata/           # Metadata management
│   │   ├── schema.go
│   │   ├── extractor.go
│   │   ├── validator.go
│   │   ├── search.go
│   │   └── export.go
│   │
│   ├── workflow/           # Workflow execution
│   │   ├── executor.go
│   │   ├── batch.go
│   │   ├── snakemake.go
│   │   ├── nextflow.go
│   │   └── environment.go
│   │
│   ├── workstation/        # Remote desktop management
│   │   ├── launcher.go
│   │   ├── session.go
│   │   ├── dcv.go
│   │   └── snapshot.go
│   │
│   ├── auth/               # Authentication & authorization
│   │   ├── iam.go
│   │   ├── globus.go
│   │   ├── session.go
│   │   └── policy.go
│   │
│   ├── doi/                # DOI management
│   │   ├── datacite.go
│   │   ├── minter.go
│   │   └── metadata.go
│   │
│   ├── portal/             # Public data portal
│   │   ├── server.go
│   │   ├── templates.go
│   │   ├── search.go
│   │   └── analytics.go
│   │
│   ├── compliance/         # NIST 800-171, HIPAA, etc.
│   │   ├── audit.go
│   │   ├── encryption.go
│   │   ├── scanner.go
│   │   └── reports.go
│   │
│   ├── cost/               # Cost tracking & optimization
│   │   ├── tracker.go
│   │   ├── predictor.go
│   │   └── optimizer.go
│   │
│   ├── config/             # Configuration management
│   │   ├── config.go
│   │   ├── lab.go
│   │   ├── project.go
│   │   └── user.go
│   │
│   └── webui/              # Web UI backend
│       ├── server.go
│       ├── api.go
│       ├── websocket.go
│       └── handlers/
│
├── pkg/
│   ├── instrument/         # Pluggable instrument adapters
│   │   ├── adapter.go
│   │   ├── zeiss.go
│   │   ├── nikon.go
│   │   ├── illumina.go
│   │   └── generic.go
│   │
│   └── schemas/            # Community metadata schemas
│       ├── core/
│       ├── microscopy/
│       ├── sequencing/
│       ├── proteomics/
│       ├── spectroscopy/
│       └── custom/
│
├── web/
│   ├── src/
│   │   ├── components/
│   │   ├── pages/
│   │   ├── stores/
│   │   └── lib/
│   ├── static/
│   └── public/
│
├── templates/              # CloudFormation/Terraform templates
│   ├── lab-init.yaml
│   ├── file-gateway.yaml
│   ├── batch-compute.yaml
│   └── workstation.yaml
│
├── docs/
│   ├── getting-started.md
│   ├── domains/           # Domain-specific guides
│   │   ├── microscopy.md
│   │   ├── sequencing.md
│   │   ├── proteomics.md
│   │   └── ...
│   ├── metadata-schemas.md
│   ├── workflows.md
│   └── compliance.md
│
└── examples/
    ├── schemas/           # Example metadata schemas
    ├── workflows/         # Example Snakemake/Nextflow
    └── configs/           # Example configurations
```

## Technology Stack

### Backend (Go)
- **CLI**: cobra
- **Config**: viper
- **AWS SDK**: aws-sdk-go-v2
- **File Watching**: fsnotify
- **HTTP**: chi or gin
- **WebSocket**: gorilla/websocket
- **Database**: SQLite (local) + DynamoDB (cloud metadata index)
- **Encryption**: crypto/aes, aws-encryption-sdk-go

### Frontend (Web UI)
- **Framework**: Svelte or Vue 3 (lightweight)
- **CSS**: Tailwind CSS
- **Charts**: Chart.js
- **File Upload**: uppy
- **Real-time**: WebSocket client

### Infrastructure (AWS)
- **Storage**: S3 (Intelligent-Tiering)
- **Compute**: AWS Batch, EC2, Fargate
- **Auth**: IAM, Cognito (optional), Globus Auth
- **Monitoring**: CloudWatch
- **Functions**: Lambda
- **CDN**: CloudFront (for portal)
- **IaC**: CloudFormation or Terraform

## Development Phases

### Phase 1: Core Storage & Sync (Weeks 1-6)
- S3 sync engine (rsync-like)
- File watching and auto-upload
- Basic CLI commands
- Local daemon
- Cost tracking

### Phase 2: Metadata & FAIR (Weeks 7-10)
- Metadata schema system
- Automatic extraction
- Validation engine
- Search/discovery
- Export formats (DataCite, ISA-Tab, etc.)

### Phase 3: Web UI & User Management (Weeks 11-14)
- Web UI (file browser, upload, etc.)
- User/group/project management
- IAM automation
- Globus Auth integration

### Phase 4: Compute & Workflows (Weeks 15-18)
- Workflow execution (Snakemake, Nextflow)
- AWS Batch integration
- Spot instance management
- Environment capture

### Phase 5: Workstations & Portal (Weeks 19-22)
- Remote workstation launcher
- DCV/noVNC integration
- Public data portal
- DOI minting (DataCite)

### Phase 6: Compliance & Polish (Weeks 23-26)
- NIST 800-171 mode
- Audit logging
- DLP scanning
- Documentation
- Community schemas

## Getting Started (for Developers)

### Prerequisites
```bash
# Install Go 1.21+
go version

# Install Node.js 20+ (for web UI)
node --version

# AWS CLI configured
aws sts get-caller-identity

# Install development tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Build & Run

```bash
# Clone repository
git clone https://github.com/your-org/cicada.git
cd cicada

# Install dependencies
go mod download

# Build CLI
make build

# Run tests
make test

# Run locally
./bin/cicada daemon start --dev

# Build web UI
cd web
npm install
npm run dev
```

### Development Workflow

1. Create feature branch
2. Write tests first (TDD)
3. Implement feature
4. Run linters: `make lint`
5. Run tests: `make test`
6. Manual testing with `--dev` mode
7. Create PR with description

## Documentation

- [Getting Started Guide](docs/getting-started.md)
- [Architecture Deep Dive](docs/architecture.md)
- [Metadata Schema Guide](docs/metadata-schemas.md)
- [Domain-Specific Guides](docs/domains/)
- [API Reference](docs/api.md)
- [Contributing](CONTRIBUTING.md)

## Community

- GitHub: https://github.com/your-org/cicada
- Discussions: https://github.com/your-org/cicada/discussions
- Slack: cicada-data.slack.com
- Twitter: @cicada_data

## License

MIT License - see [LICENSE](LICENSE)

## Funding

This project was developed with support from:
- [Grant agency if applicable]
- Community contributions

## Citation

If you use Cicada in your research, please cite:

```
@software{cicada2024,
  title = {Cicada: Dormant Data Commons for Academic Research},
  author = {Your Name and Contributors},
  year = {2024},
  url = {https://github.com/your-org/cicada}
}
```
