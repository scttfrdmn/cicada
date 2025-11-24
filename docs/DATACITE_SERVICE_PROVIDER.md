# DataCite Integration Strategy for Cicada

**Date:** January 24, 2025
**Status:** Research Complete - Awaiting Decision
**Related Issues:** #27 (DataCite API Client), #36 (Integration Tests), #37 (Provider Documentation)

## Executive Summary

After researching DataCite integration options, **becoming a DataCite Registered Service Provider** is the recommended path for Cicada. This approach:

- ✅ **Zero cost** to Cicada as a project
- ✅ **Scales** to unlimited users (each brings their own credentials)
- ✅ **Official certification** from DataCite
- ✅ **Best user experience** for institutional users

**Alternative path for testing:** Use Lyrasis consortium membership ($1,625/year) for sandbox/early development, then transition to Registered Service Provider for production.

## Understanding DataCite Membership Models

### Model 1: Direct Institutional Membership

**Who it's for:** Individual institutions (universities, research centers)

**Cost:**
- €2,000/year base fee
- Plus per-DOI fees (volume-based pricing)
- Total: €2,500 - €5,000+/year depending on usage

**For Cicada:** ❌ **Not Recommended**
- Expensive for a single open-source project
- Doesn't scale to multiple users
- Requires Cicada to manage all DOIs

### Model 2: Consortium Membership

**Who it's for:** Regional organizations (e.g., US institutions via Lyrasis)

**Cost (Lyrasis example):**
- $1,625/year for 1-1,999 DOIs ($1/DOI after)
- $3,600/year flat rate for 2,000-10,000 DOIs
- Other consortia: British Library (UK), EUDAT (EU)

**For Cicada:** ⚠️ **Potential for Early Development**
- Lower cost than direct membership
- Could be used for sandbox testing
- Still requires annual fees
- Temporary solution, not long-term strategy

### Model 3: Registered Service Provider (RECOMMENDED)

**Who it's for:** Software applications that integrate DataCite API

**Cost:**
- **FREE** for the software provider (Cicada)
- Users bring their own DataCite credentials
- No DOI limits, no per-user fees

**For Cicada:** ✅ **RECOMMENDED**
- No cost to Cicada project
- Users use their institution's existing DataCite membership
- Official DataCite certification and listing
- Scales to unlimited users
- Best fit for open-source CLI tool

## What is a DataCite Registered Service Provider?

### Definition

A **Registered Service Provider** is software that:
1. Integrates the DataCite REST API
2. Allows DataCite members to register DOIs using their own credentials
3. Meets DataCite's technical and security requirements
4. Is officially certified and listed by DataCite

### How It Works

```
┌─────────────────────────────────────────────────────────┐
│                      Cicada User                        │
│                  (Research Institution)                 │
│                                                          │
│  Has DataCite membership through:                       │
│  - Direct membership (€2,000/year)                      │
│  - Consortium (e.g., Lyrasis $1,625/year)               │
│  - Institutional affiliation                            │
└─────────────────────────────────────────────────────────┘
                           │
                           │ Provides credentials:
                           │ - Repository ID (prefix)
                           │ - Password
                           ▼
┌─────────────────────────────────────────────────────────┐
│                    Cicada (v0.3.0+)                     │
│              Registered Service Provider                │
│                                                          │
│  $ cicada doi publish sample.fastq \                    │
│      --provider datacite \                              │
│      --datacite-repository-id <user's ID> \             │
│      --datacite-password <user's password>              │
└─────────────────────────────────────────────────────────┘
                           │
                           │ API Request with
                           │ user's credentials
                           ▼
┌─────────────────────────────────────────────────────────┐
│              DataCite REST API                          │
│          https://api.datacite.org                       │
│                                                          │
│  Validates credentials, mints DOI using                 │
│  user's allocation                                      │
└─────────────────────────────────────────────────────────┘
                           │
                           │ DOI: 10.5072/xxxxx
                           ▼
┌─────────────────────────────────────────────────────────┐
│                  User's DOI Record                      │
│              Registered under their prefix              │
└─────────────────────────────────────────────────────────┘
```

### Key Benefits for Cicada

1. **Zero Cost:** Cicada doesn't need a DataCite membership
2. **No User Limits:** Support unlimited institutions
3. **Credential Security:** Users keep their own credentials
4. **Official Status:** DataCite certification badge
5. **Discoverability:** Listed on DataCite's website
6. **User Trust:** Official certification provides credibility

### Key Benefits for Users

1. **Use Existing Membership:** No additional DataCite costs
2. **Institutional Control:** DOIs registered under their prefix
3. **Quota Management:** Uses their existing allocation
4. **Compliance:** Meets institutional requirements
5. **Certified Software:** Validated by DataCite

## Requirements for Becoming a Registered Service Provider

### 1. Technical Requirements

#### API Integration
- ✅ Implement DataCite REST API v2
- ✅ Support all CRUD operations (Create, Read, Update, Delete)
- ✅ Handle authentication (HTTP Basic Auth)
- ✅ Support both sandbox and production environments
- ✅ Implement proper error handling and retry logic

**Status for Cicada:** Planned in Issue #27 (DataCite API Client)

#### Metadata Schema Support
- ✅ Support DataCite Metadata Schema v4.5 (current)
- ✅ Commit to updating when schema changes
- ✅ Validate metadata before submission

**Status for Cicada:** ✅ Already implemented in v0.2.0

#### Demonstration Requirement
- ✅ Have registered DOIs in production (findable, not just draft)
- ✅ Demonstrate working integration on verification call

**Status for Cicada:** ⏳ Will be ready after Milestone 1 (Issue #36)

### 2. Security Requirements

- ✅ **Secure credential handling:**
  - Never log credentials
  - Support environment variables
  - Support config files with proper permissions
  - Clear documentation on credential security

- ✅ **HTTPS only:** All API calls over TLS
- ✅ **Error handling:** Don't expose credentials in error messages
- ✅ **Documentation:** Security best practices for users

**Status for Cicada:** Planned in Issue #26 (Provider Configuration)

### 3. Documentation Requirements

- ✅ User guide for DataCite integration
- ✅ API usage documentation
- ✅ Security best practices
- ✅ Troubleshooting guide

**Status for Cicada:** Planned in Issue #37 (Provider Documentation)

### 4. Best Practices

- ✅ Follow RESTful principles
- ✅ Implement rate limiting
- ✅ Use appropriate HTTP methods and status codes
- ✅ Provide clear error messages to users
- ✅ Monitor API deprecation notices

**Status for Cicada:** Planned in Issue #35 (Error Handling & Retry Logic)

## The Registration Process

### Step 1: Prepare Integration (Weeks 1-4)

**Goal:** Implement DataCite API client meeting all requirements

**Tasks:**
- ✅ Issue #26: Provider Configuration System
- ✅ Issue #27: DataCite API Client
- ✅ Issue #30: DataCite Metadata Mapping
- ✅ Issue #32: CLI `doi publish` Command
- ✅ Issue #35: Error Handling & Retry Logic
- ✅ Issue #36: Integration Tests (with sandbox)

**Deliverable:** Working integration with DataCite sandbox

**Timeline:** Milestone 1 (4 weeks)

### Step 2: Contact DataCite (Week 4-5)

**Action:** Email support@datacite.org

**Email Template:**
```
Subject: Application to Become DataCite Registered Service Provider - Cicada

Dear DataCite Team,

I am writing to apply for Cicada to become a DataCite Registered Service Provider.

About Cicada:
- Open-source CLI tool for scientific data management
- Written in Go, available on GitHub: https://github.com/scttfrdmn/cicada
- Target users: Research labs, bioinformatics facilities, microscopy cores
- Version 0.3.0 will include full DataCite REST API v2 integration

Current Implementation Status:
✅ DataCite Metadata Schema v4.5 support
✅ REST API v2 client implementation
✅ Sandbox environment testing
✅ Security best practices (credential handling, HTTPS)
✅ Comprehensive documentation

We have successfully:
- Registered test DOIs in the sandbox environment
- Implemented all CRUD operations
- Validated metadata compliance
- Created user documentation for DataCite integration

We would like to schedule a verification call to demonstrate our integration
and complete the registration process.

Best regards,
[Your Name]
Cicada Project Lead
```

**Expected Response:** DataCite will review and schedule verification call

### Step 3: Verification Call (Week 5-6)

**What to Demonstrate:**
1. Live demo of Cicada registering DOI via DataCite sandbox
2. Show metadata validation
3. Show error handling
4. Show credential security
5. Walk through documentation

**Preparation:**
- Have sandbox credentials ready
- Prepare test data file
- Have documentation open
- Rehearse demo flow

**Demo Script:**
```bash
# 1. Show version and provider support
cicada version
cicada doi provider list

# 2. Show metadata extraction
cicada metadata extract sample.fastq --preset illumina-novaseq

# 3. Show DOI preparation (without publishing)
cicada doi prepare sample.fastq \
  --enrich metadata.yaml \
  --publisher "Test Lab" \
  --output doi-metadata.json

# 4. Show validation
cat doi-metadata.json | jq '.'

# 5. Publish to sandbox
cicada doi publish sample.fastq \
  --enrich metadata.yaml \
  --provider datacite \
  --publisher "Test Lab" \
  --datacite-repository-id <SANDBOX_ID> \
  --datacite-password <SANDBOX_PASSWORD> \
  --datacite-sandbox

# 6. Show registered DOI
cicada doi status <DOI>
```

### Step 4: Complete Registration Form (Week 6)

**Information Needed:**
- Application name: **Cicada**
- Version: **0.3.0**
- Website: https://github.com/scttfrdmn/cicada
- Documentation URL: https://github.com/scttfrdmn/cicada/blob/main/docs/DOI_PUBLISHING.md
- Contact email: [Your email]
- API version: **REST API v2**
- Supported schema: **DataCite Metadata Schema v4.5**
- License: **MIT** (or whatever Cicada uses)
- Description: *"Open-source CLI tool for scientific data management and DOI registration. Supports automated metadata extraction from scientific file formats and DOI minting through multiple providers including DataCite and Zenodo."*

### Step 5: Receive Certification (Week 7)

**What You Get:**
- ✅ Official listing on DataCite website
- ✅ "DataCite Registered Service Provider" badge
- ✅ Logo usage permission
- ✅ Direct support contact at DataCite
- ✅ Early notice of API changes

**What to Do:**
- Add badge to README.md
- Update documentation with official status
- Announce in release notes
- Add to project website

## Cost Comparison: Three Scenarios

### Scenario A: Registered Service Provider (RECOMMENDED)

**Cost to Cicada:** $0

**Cost to Users:** Their existing DataCite membership
- Direct: €2,000/year + per-DOI fees
- Consortium (Lyrasis): $1,625-$3,600/year
- Already have it: $0 additional

**Pros:**
- ✅ No cost to Cicada
- ✅ Scales to unlimited users
- ✅ Users control their own DOIs
- ✅ Official DataCite certification
- ✅ Best for open-source model

**Cons:**
- ⚠️ Users must have DataCite membership
- ⚠️ More complex credential management
- ⚠️ Can't provide "free" DataCite DOIs to users

**Best for:**
- Open-source projects
- Tools for institutional users
- Multi-tenant scenarios

### Scenario B: Direct Institutional Membership

**Cost to Cicada:** €2,000-€5,000/year

**Cost to Users:** $0 (Cicada provides DOIs)

**Pros:**
- ✅ Can provide "free" DOIs to users
- ✅ Simple credential management
- ✅ Full control

**Cons:**
- ❌ Expensive for open-source project
- ❌ Doesn't scale (volume pricing)
- ❌ Ongoing maintenance cost
- ❌ Single point of failure
- ❌ All DOIs under Cicada's prefix (user ownership unclear)

**Best for:**
- Commercial services
- Institutional tools with dedicated budget
- Low-volume use cases

### Scenario C: Consortium Membership (TEMPORARY)

**Cost to Cicada:** $1,625/year (Lyrasis, up to 2,000 DOIs)

**Cost to Users:** $0 (Cicada provides DOIs)

**Pros:**
- ✅ Lower cost than direct membership
- ✅ Good for early testing/development
- ✅ Can provide DOIs to beta testers

**Cons:**
- ⚠️ Still has annual cost
- ⚠️ Volume limits
- ⚠️ Not sustainable long-term
- ⚠️ US-specific (Lyrasis)

**Best for:**
- Early development phase
- Beta testing with real DOIs
- Temporary solution before becoming Service Provider

## Recommended Strategy for Cicada

### Phase 1: Sandbox Development (Weeks 1-4) - FREE

**Goal:** Build and test DataCite integration

**Approach:**
1. Use DataCite sandbox (no membership required)
2. Implement all API features (Issues #26, #27, #30)
3. Write integration tests (Issue #36)
4. Document provider setup (Issue #37)

**Cost:** $0

### Phase 2: Verification (Weeks 5-7) - FREE

**Goal:** Become Registered Service Provider

**Approach:**
1. Contact DataCite (support@datacite.org)
2. Schedule verification call
3. Demonstrate integration
4. Complete registration form
5. Receive certification

**Cost:** $0

### Phase 3: Production Use (v0.3.0 Release) - FREE

**Goal:** Launch with official DataCite support

**Approach:**
1. Document credential requirements for users
2. Provide setup guides for different institution types
3. Include Service Provider badge in documentation
4. Support users with their own credentials

**Cost:** $0 (users bring their own DataCite membership)

### Optional: Early Testing with Real DOIs (Week 4-6)

**If you need production DOIs before becoming Service Provider:**

**Option A:** Personal/Institutional Account
- If you have access to a DataCite account (university, etc.)
- Use for early testing only
- Cost: $0 (already have it)

**Option B:** Lyrasis Consortium Membership
- Sign up at https://www.lyrasis.org/programs/Pages/DataCite-US-Community-Membership.aspx
- $1,625/year (1-1,999 DOIs)
- Use for beta testing, then transition to Service Provider
- Cost: $1,625 (one-time or temporary)

**Recommendation:** Not necessary - sandbox is sufficient for development and certification

## User Documentation Strategy

### For Institutional Users (Most Common)

**Title:** Using Cicada with Your Institution's DataCite Membership

**Content:**
```markdown
# Using Cicada with DataCite

Cicada is a DataCite Registered Service Provider. If your institution has a
DataCite membership, you can use Cicada to mint DOIs using your credentials.

## Prerequisites

- DataCite repository ID (prefix)
- DataCite password
- Your institution's DataCite membership (direct or consortium)

## Getting Your Credentials

Contact your institution's library, IT department, or research data office.
They can provide:
- Your repository ID (format: XXXX.YYYY)
- Your password
- Information about your DOI allocation

## Configuration

### Option 1: Environment Variables (Recommended)

export DATACITE_REPOSITORY_ID="your-repo-id"
export DATACITE_PASSWORD="your-password"

cicada doi publish sample.fastq --provider datacite

### Option 2: Command Line Flags

cicada doi publish sample.fastq \
  --provider datacite \
  --datacite-repository-id your-repo-id \
  --datacite-password your-password

### Option 3: Config File

Create ~/.cicada/config.yaml:

providers:
  datacite:
    repository_id: your-repo-id
    password: your-password
    environment: production

cicada doi publish sample.fastq --provider datacite

## Security Best Practices

- Never commit credentials to version control
- Use environment variables or config files with proper permissions
- Rotate passwords regularly
- Use separate credentials for production and testing
```

### For Users Without DataCite

**Title:** Getting Started with DOI Registration

**Content:**
```markdown
# DOI Registration Options

## Option 1: Use Zenodo (FREE, Recommended for Individuals)

Zenodo provides free DOI registration for all users:

cicada doi publish sample.fastq --provider zenodo

See docs/ZENODO_SETUP.md for details.

## Option 2: DataCite via Your Institution

If you're affiliated with a university or research institution:

1. Check if your institution has a DataCite membership
2. Contact your library or IT department
3. Request DataCite credentials
4. Use Cicada with your institutional credentials

See docs/DATACITE_SETUP.md for details.

## Option 3: DataCite Consortium Membership

For independent researchers or small labs:

- US: Lyrasis ($1,625/year for up to 2,000 DOIs)
  https://www.lyrasis.org/programs/Pages/DataCite-US-Community-Membership.aspx

- UK: British Library (~£1,100/year)
  https://www.bl.uk/britishlibrary/~/media/bl/global/services/collection%20metadata/pdfs/datacite-membership-application-form.pdf

- Europe: EUDAT
  https://www.eudat.eu/catalogue/B2SHARE

- Australia: ANDS/ARDC
  https://ardc.edu.au/services/identifier/
```

## Next Steps

### Immediate Actions (This Week)

- [x] ✅ Research DataCite options (COMPLETE)
- [x] ✅ Create this decision document (COMPLETE)
- [ ] 👤 **User Decision:** Confirm Registered Service Provider approach
- [ ] 👤 **User Action:** Review Zenodo account setup (already created)

### Before Starting Development (Week 1)

- [ ] Create DataCite sandbox account (free)
  - URL: https://support.datacite.org/docs/testing-guide
  - Sign up at: https://doi.test.datacite.org/
- [ ] Store sandbox credentials securely
- [ ] Review DataCite REST API documentation
  - API Docs: https://support.datacite.org/docs/api
  - Schema: https://schema.datacite.org/

### During Development (Weeks 1-4)

- [ ] Implement Issue #26: Provider Configuration System
- [ ] Implement Issue #27: DataCite API Client
- [ ] Implement Issue #30: DataCite Metadata Mapping
- [ ] Implement Issue #32: `doi publish` Command
- [ ] Implement Issue #35: Error Handling & Retry Logic
- [ ] Implement Issue #36: Integration Tests (sandbox)
- [ ] Write Issue #37: Provider Documentation

### After Development (Weeks 5-7)

- [ ] Email DataCite: support@datacite.org
- [ ] Prepare verification demo
- [ ] Schedule verification call
- [ ] Complete registration form
- [ ] Receive certification
- [ ] Update documentation with badge
- [ ] Announce in v0.3.0 release notes

### Optional: Early Production Testing

- [ ] Decide if needed (probably not necessary)
- [ ] If yes, evaluate Lyrasis vs institutional account
- [ ] Register test account
- [ ] Use for beta testing only
- [ ] Transition to Service Provider model for release

## Decision Required

**Question:** Confirm the Registered Service Provider approach for Cicada v0.3.0?

**Recommendation:** ✅ **YES** - Proceed with Registered Service Provider path

**Rationale:**
1. Zero cost to Cicada project
2. Scales to unlimited institutional users
3. Official DataCite certification
4. Best fit for open-source model
5. Sandbox sufficient for development

**Alternative considered:** Temporary Lyrasis membership for early testing
- **Verdict:** Not necessary - sandbox is adequate

---

**Next Step:** Confirm decision, then begin Milestone 1 (Issue #26: Provider Configuration System)
