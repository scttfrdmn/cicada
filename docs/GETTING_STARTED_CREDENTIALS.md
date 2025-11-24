# Getting Started with DOI Publishing

**Setting up Cicada to publish DOIs for your research data**

**For researchers and lab managers with minimal technical experience**

This guide will help you set up Cicada to publish DOIs (Digital Object Identifiers) for your research data. No technical background required - just follow the steps!

## What Should Get a DOI?

Before setting up credentials, it's helpful to understand what research outputs typically get DOIs and what your field expects.

### What is a DOI?

A **DOI (Digital Object Identifier)** is a permanent identifier for your research data - like a permanent web address that never breaks. Even if you move your data to a different server or repository, the DOI always points to the correct location.

**Example DOI:** `10.5281/zenodo.123456`
- Anyone can click this and find your data
- It works forever, even if URLs change
- It's citable in publications

### What Typically Gets a DOI

**✅ Should have a DOI:**

**Research Data:**
- Raw data that supports a publication
- Datasets that others might reuse
- Data required by journal or funder policies
- Data that you want to cite in papers

**Examples:**
- Sequencing data (FASTQ, BAM files)
- Microscopy images (CZI, OME-TIFF files)
- Tabular data (CSV, Excel files)
- Simulation outputs
- Survey data

**Software & Code:**
- Research software releases
- Analysis scripts and pipelines
- Computational models

**Other Research Outputs:**
- Protocols
- Supplementary materials
- Technical reports
- Preprints (though use preprint servers like bioRxiv)

### What Typically Doesn't Need a DOI

**❌ Usually doesn't need a DOI:**

- **Intermediate/temporary files** - Processing outputs that aren't final
- **Duplicate data** - Data already published elsewhere with a DOI
- **Preliminary data** - Still collecting or analyzing
- **Personal notes** - Lab notebooks, internal documentation
- **Already published data** - Available in established databases (NCBI, ENA, etc.)
- **Tiny test files** - Sample data for tutorials

### Domain-Specific Practices

Different research fields have different expectations and norms:

#### Genomics & Bioinformatics

**Common practice:**
- Raw sequencing data (FASTQ) → Submit to NCBI SRA/ENA (they provide accessions)
- Assembled genomes → Submit to NCBI GenBank
- Supporting data (metadata, analysis results) → DOI via Zenodo/DataCite
- Software/pipelines → DOI via Zenodo + GitHub release

**When to use Cicada:**
- Small-scale sequencing projects not in public databases
- Supplementary data for publications
- Custom reference datasets
- Analysis pipelines and scripts

**Example:**
```bash
# Raw reads submitted to SRA (gets accession: SRR123456)
# Supplementary metadata gets DOI via Cicada:
cicada doi publish sample_metadata.csv quality_metrics.csv \
  --provider zenodo \
  --enrich metadata.yaml
```

#### Microscopy & Imaging

**Common practice:**
- Publication-quality images → DOI via institutional repository or Zenodo
- Large imaging datasets → Specialized repositories (IDR, BioImage Archive)
- Instrument calibration data → DOI for reproducibility
- Image analysis workflows → DOI for transparency

**When to use Cicada:**
- Original microscopy files (CZI, OME-TIFF) with metadata
- Supporting images for papers
- Method validation datasets
- Training datasets for analysis

**Example:**
```bash
# Publish microscopy dataset with DOI
cicada doi publish experiment_*.czi \
  --preset zeiss-lsm-880 \
  --provider zenodo \
  --enrich metadata.yaml
```

#### Structural Biology & Chemistry

**Common practice:**
- Crystal structures → PDB (Protein Data Bank)
- NMR data → BMRB
- Supporting data → DOI via Zenodo/DataCite
- Protocols → DOI via protocols.io or Zenodo

**When to use Cicada:**
- Supplementary crystallography data
- Raw diffraction images
- Validation datasets
- Custom structure libraries

#### Ecology & Environmental Science

**Common practice:**
- Long-term datasets → Domain repositories (LTER, DataONE)
- Species occurrence data → GBIF
- Climate/sensor data → Specialized repositories
- Small studies → DOI via Zenodo/DataCite

**When to use Cicada:**
- Field study data
- Sensor/instrument data
- Observation records
- Derived datasets

#### Social Sciences & Humanities

**Common practice:**
- Survey data → Domain repositories (ICPSR, UK Data Service)
- Qualitative data → Institutional repositories + DOI
- Interview transcripts → Repository + DOI
- Datasets for replication → DOI required

**When to use Cicada:**
- Anonymized survey data
- Supplementary data for publications
- Replication packages
- Coded datasets

### Publisher & Funder Requirements

Many journals and funding agencies now **require** data DOIs:

**Publishers with data policies:**
- Nature journals - "Data availability statement" with DOIs
- PLOS - Data deposition required
- eLife - Data must be accessible with identifiers
- Cell Press - Data availability required

**Funders requiring data sharing:**
- NIH - Data Management and Sharing Policy
- NSF - Data sharing plans required
- Wellcome Trust - Data must be accessible
- European Commission (Horizon) - Open data by default

**Check your requirements:**
1. Journal submission guidelines ("Data Availability")
2. Funder data management plan
3. Institutional policies

### Best Practices

#### When to Create a DOI

**✅ Create a DOI when:**
- Data is finalized and won't change significantly
- You're ready to share publicly (or with embargo)
- Data supports a publication
- Required by journal/funder
- You want others to cite your data

**⏸️ Wait to create a DOI if:**
- Still collecting/analyzing data
- Data quality is uncertain
- Need institutional approval first
- Preparing for submission to domain repository

#### Versioning

**When to create a new version (same DOI):**
- Minor corrections to metadata
- Adding additional files
- Fixing errors in existing data
- Updating documentation

**When to create a new DOI:**
- Major changes to dataset
- Different analysis/processing
- Conceptually different dataset
- Different paper/project

**Example:**
```bash
# Original dataset (DOI: 10.5281/zenodo.123456)
cicada doi publish data_v1.fastq --provider zenodo

# Minor update (new version, same DOI concept)
cicada doi update 10.5281/zenodo.123456 \
  --add-file corrected_metadata.csv

# Major reanalysis (new DOI)
cicada doi publish reanalyzed_data_v2.fastq --provider zenodo
# Gets new DOI: 10.5281/zenodo.789012
```

#### What Metadata to Include

**Minimum (required):**
- Title (descriptive, not just filename)
- Authors/Creators (with affiliations)
- Description (what is the data, how was it collected)
- Keywords (for discoverability)
- License (how can others use it?)

**Recommended:**
- Related publications (DOIs)
- Funding information
- Methods/protocols
- File formats and software requirements
- Related identifiers (ORCID, grant numbers)

**Example metadata file:**
```yaml
title: "RNA-seq data from drought stress experiment in Arabidopsis"
description: |
  Paired-end Illumina RNA sequencing (150bp) from Arabidopsis thaliana
  leaves under drought stress and control conditions. Three biological
  replicates per condition. Sequenced on NovaSeq 6000.
creators:
  - name: "Jane Smith"
    affiliation: "Plant Biology Lab, State University"
    orcid: "0000-0002-1234-5678"
  - name: "John Doe"
    affiliation: "Plant Biology Lab, State University"
keywords:
  - RNA-seq
  - Arabidopsis thaliana
  - drought stress
  - transcriptomics
license: "CC-BY-4.0"
related_publications:
  - doi: "10.1234/journal.2024.123"
    relation: "IsSupplementTo"
funding:
  - funder: "National Science Foundation"
    award: "NSF-1234567"
```

### Quick Decision Guide

**Ask yourself:**

1. **Is this data final?**
   - Yes → Create DOI ✅
   - No → Wait ⏸️

2. **Will others need to cite or access it?**
   - Yes → Create DOI ✅
   - No → Maybe not needed ❌

3. **Is there a domain-specific repository?**
   - Yes → Check if they provide DOIs (use that first)
   - No → Use Zenodo/DataCite via Cicada ✅

4. **Does my journal/funder require it?**
   - Yes → Create DOI ✅
   - No → Optional but recommended

5. **Is the data already public with a stable ID?**
   - Yes → Use existing identifier ❌
   - No → Create DOI ✅

**Still not sure?** Ask your:
- Institutional library
- Research data office
- Department administrator
- Journal editor
- Colleagues in your field

## What You Need

To publish DOIs with Cicada, you need credentials (like a password) from a DOI registration service.

**Two main options:**
- **Zenodo** (free for everyone)
- **DataCite** (requires institutional membership)

### Understanding Zenodo vs DataCite

**What are these services?**

Both Zenodo and DataCite are services that create permanent DOIs (Digital Object Identifiers) for your research data. A DOI is like a permanent web address that always points to your data, even if it moves to a different server.

**Zenodo:**
- A **free repository and DOI service** run by CERN (the European physics lab)
- Anyone can create an account and start publishing immediately
- Includes free file storage (up to 50 GB per dataset)
- DOIs look like: `10.5281/zenodo.123456`
- Perfect for: Individual researchers, small labs, open science projects
- Website: https://zenodo.org

**DataCite:**
- A **membership-based DOI service** for institutions
- Your university or research institution needs to be a member
- You use your institution's credentials through Cicada
- DOIs include your institution's prefix: `10.12345/yourdata`
- Perfect for: Large institutions, institutional repositories, compliance requirements
- Website: https://datacite.org

### Which One Should You Use?

**Check with your institution's library FIRST!**

Before setting anything up, contact your institutional library or research data office and ask:

> "Does our institution have a DataCite membership or preferred DOI service for research data?"

**If they say YES to DataCite:**
- ✅ Use DataCite (see [Part 2B: DataCite Setup](#for-datacite-users))
- Your DOIs will be under your institution's prefix
- May be required for institutional compliance
- Ask them for your credentials (repository ID and password)

**If they say NO or don't have DOI services:**
- ✅ Use Zenodo (see [Part 1: Zenodo Setup](#part-1-getting-your-zenodo-credentials-5-minutes))
- Free and easy to set up yourself
- No institutional approval needed
- Works great for most research data

**If you're not sure or don't have institutional support:**
- ✅ Start with Zenodo - it's free and takes 5 minutes
- You can always add DataCite later if needed

**Can I use both?**

Yes! Cicada supports both, and you can switch between them with a simple flag:
```bash
cicada doi publish data.fastq --provider zenodo
cicada doi publish data.fastq --provider datacite
```

## Part 1: Getting Your Zenodo Credentials (5 minutes)

### Step 1: Create a Zenodo Account

1. Go to https://zenodo.org
2. Click **"Sign Up"** in the top right corner
3. You can sign up with:
   - Your email address
   - Your GitHub account (if you have one)
   - Your ORCID (if you have one)
4. Follow the instructions to complete sign-up

**That's it! Your Zenodo account is ready.**

### Step 2: Create an Access Token

An access token is like a password that Cicada uses to publish DOIs on your behalf.

1. Log in to Zenodo: https://zenodo.org
2. Click your **name/profile icon** in the top right
3. Click **"Applications"** in the dropdown menu
4. Click the **"New Token"** button
5. Give your token a name (e.g., "Cicada DOI Publishing")
6. Under **"Scopes"**, check the box for **"deposit:write"**
   - This lets Cicada create DOIs for you
   - Don't worry - it can only create DOIs, not delete anything
7. Click **"Create"**
8. **Important:** Copy the token that appears - you'll need it in the next step
   - It looks like: `AbCdEf123456...` (40+ characters)
   - You can only see it once, so copy it now!

**✓ You now have your Zenodo access token!**

---

## Part 2: Giving Cicada Your Credentials

You have **three options** for giving Cicada your token. Choose the one that sounds easiest to you.

### Option A: Simple Setup (Recommended for Most Users)

This is the easiest method. Cicada will guide you through setup.

**On macOS or Linux:**

1. Open Terminal (don't worry, we'll guide you!)
   - **macOS:** Press `Cmd + Space`, type "Terminal", press Enter
   - **Linux:** Press `Ctrl + Alt + T`

2. Type this command and press Enter:
   ```bash
   cicada config init
   ```

3. Cicada will ask: "Enter your Zenodo token:"
   - Paste the token you copied earlier (right-click → Paste)
   - Don't worry if you don't see the token when you paste - that's for security!
   - Press Enter

4. Cicada will save your token securely

**That's it! Your credentials are set up.**

To test that it worked:
```bash
cicada config test zenodo
```

You should see: `✓ Zenodo authentication successful`

### Option B: Manual Setup (Config File)

If you prefer to create the file yourself:

**On macOS or Linux:**

1. Create a folder for Cicada settings:
   ```bash
   mkdir -p ~/.config/cicada
   ```

2. Create a settings file:
   ```bash
   nano ~/.config/cicada/config.yaml
   ```
   (You can also use your favorite text editor instead of `nano`)

3. Copy and paste this into the file:
   ```yaml
   providers:
     zenodo:
       token: "PASTE_YOUR_TOKEN_HERE"
   ```

4. Replace `PASTE_YOUR_TOKEN_HERE` with the token you copied from Zenodo
   - Keep the quotes around it!
   - Example: `token: "AbCdEf123456..."`

5. Save and close the file:
   - In nano: Press `Ctrl + X`, then `Y`, then Enter

6. Secure your file (important!):
   ```bash
   chmod 600 ~/.config/cicada/config.yaml
   ```
   This makes sure only you can read your credentials.

**On Windows:**

1. Open Notepad

2. Copy and paste this:
   ```yaml
   providers:
     zenodo:
       token: "PASTE_YOUR_TOKEN_HERE"
   ```

3. Replace `PASTE_YOUR_TOKEN_HERE` with your Zenodo token

4. Save the file as:
   ```
   C:\Users\YourUsername\AppData\Roaming\cicada\config.yaml
   ```
   Replace `YourUsername` with your Windows username.

### Option C: Project-Specific Setup (For Advanced Users)

If you're working on a specific project and want credentials just for that project:

1. In your project folder, create a file named `.env`

2. Add this line:
   ```
   CICADA_ZENODO_TOKEN=your-token-here
   ```

3. **Important:** Tell Git to ignore this file (so you don't accidentally share your token):
   ```bash
   echo ".env" >> .gitignore
   ```

---

## Part 3: Testing Your Setup

Let's make sure everything is working:

```bash
cicada config test zenodo
```

**If you see:** `✓ Zenodo authentication successful`
- **Success!** You're all set up and ready to publish DOIs.

**If you see an error**, see the [Troubleshooting](#troubleshooting) section below.

---

## Part 4: Publishing Your First DOI

Now that you're set up, here's how to publish a DOI:

### Simple Example

```bash
cicada doi publish my-data-file.fastq --provider zenodo
```

This will:
1. Extract metadata from your file
2. Upload it to Zenodo
3. Create a DOI for it
4. Print the DOI (something like `10.5281/zenodo.123456`)

### With Additional Information

You probably want to add more details about your data:

1. Create a file called `metadata.yaml` with your information:
   ```yaml
   title: "My Research Data"
   description: "RNA sequencing data from..."
   creators:
     - name: "Jane Smith"
       affiliation: "University Lab"
   keywords:
     - RNA-seq
     - genomics
   ```

2. Publish with the extra information:
   ```bash
   cicada doi publish my-data-file.fastq \
     --enrich metadata.yaml \
     --provider zenodo
   ```

**That's it!** You've published your first DOI.

---

## Part 2B: DataCite Setup (For Institutional Users)

If your institution has a DataCite membership (recommended to check with your library first):

### Step 1: Get Your Credentials from Your Institution

Contact your institution's library, IT department, or research data office and ask:

> "I need DataCite credentials to publish DOIs. Can you provide me with:
> - A DataCite repository ID
> - A DataCite password
> - Information about our DOI allocation"

They should give you:
- **Repository ID**: Looks like `10.12345/INST` or `CLIENT.MEMBER`
- **Password**: A password for that repository

### Step 2: Configure Cicada

**Option 1: Simple setup (recommended):**
```bash
cicada config init
```
When prompted, enter your DataCite repository ID and password.

**Option 2: Manual config file:**

Create `~/.config/cicada/config.yaml`:
```yaml
providers:
  datacite:
    repository_id: "your-repo-id-here"
    password: "your-password-here"
```

Then secure it:
```bash
chmod 600 ~/.config/cicada/config.yaml
```

### Step 3: Test It

```bash
cicada config test datacite
```

Should show: `✓ DataCite authentication successful`

### Step 4: Publish DOIs

```bash
cicada doi publish my-data-file.fastq --provider datacite
```

---

## Troubleshooting

### "Command not found: cicada"

**Problem:** Your computer doesn't know where Cicada is installed.

**Solution:**
1. Make sure you installed Cicada (see main README)
2. Try the full path: `/usr/local/bin/cicada` instead of just `cicada`

### "Authentication failed"

**Problem:** Your credentials aren't working.

**Check these:**

1. **Did you copy the token correctly?**
   - Go back to Zenodo and create a new token
   - Make sure you copy the entire token (no spaces before/after)
   - Try pasting it again

2. **Did you create the token with the right permissions?**
   - The token needs "deposit:write" scope
   - Create a new token if you're not sure

3. **Are you using sandbox vs production?**
   - If you're testing, make sure you use a sandbox token
   - See "Testing with Sandbox" below

### "Insecure permissions on config.yaml"

**Problem:** Your credentials file can be read by other users on your computer.

**Solution:**
```bash
chmod 600 ~/.config/cicada/config.yaml
```

This command makes the file readable only by you.

### "File not found: ~/.config/cicada/config.yaml"

**Problem:** The config file doesn't exist yet.

**Solution:** Run the setup command:
```bash
cicada config init
```

Or create the file manually (see Option B above).

### ".env file is tracked by git"

**Problem:** Your credentials are about to be shared on GitHub!

**Solution:**
```bash
echo ".env" >> .gitignore
git rm --cached .env
git commit -m "Remove credentials from git"
```

This removes the credentials from git and prevents it from happening again.

### "Can't connect to Zenodo/DataCite"

**Problem:** Network or service issue.

**Check:**
1. Do you have internet connection?
2. Can you access https://zenodo.org in your browser?
3. Is your institution's firewall blocking API access?

**Try:**
- Wait a few minutes and try again
- Contact your IT department if the problem persists

---

## Testing with Sandbox (Optional)

Before publishing real DOIs, you might want to test with a sandbox (a test environment).

### Zenodo Sandbox

1. Create a sandbox account: https://sandbox.zenodo.org
   - This is separate from your main Zenodo account
2. Create a token the same way as before
3. Tell Cicada to use sandbox:
   ```bash
   cicada doi publish test-file.fastq \
     --provider zenodo \
     --zenodo-sandbox
   ```

Or in your config file:
```yaml
providers:
  zenodo:
    token: "your-sandbox-token"
    environment: sandbox
```

### DataCite Sandbox

1. Get sandbox credentials from your institution (or create test account)
2. Use the `--datacite-sandbox` flag:
   ```bash
   cicada doi publish test-file.fastq \
     --provider datacite \
     --datacite-sandbox
   ```

**Sandbox DOIs are NOT REAL** - they're for testing only and will be deleted eventually.

---

## Security Tips (Please Read!)

Your credentials are like passwords - keep them safe:

### ✅ DO:
- Keep credentials in the config file (`~/.config/cicada/config.yaml`)
- Make sure the config file has secure permissions (`chmod 600`)
- Create different tokens for different purposes (sandbox vs production)
- Delete old tokens when you're not using them anymore

### ❌ DON'T:
- Share your tokens with anyone
- Commit tokens to git/GitHub
- Email tokens
- Put tokens in filenames or comments
- Use production tokens for testing

### If Your Token is Exposed

If you accidentally shared your token:

1. Go to Zenodo → Applications
2. Find your token
3. Click "Revoke" to delete it
4. Create a new token
5. Update your Cicada config with the new token

---

## Quick Reference

### Setup Commands

```bash
# Initialize config (guided setup)
cicada config init

# Test credentials
cicada config test zenodo
cicada config test datacite

# View current config (tokens hidden)
cicada config show

# Validate security
cicada config validate
```

### Publishing Commands

```bash
# Basic publish
cicada doi publish file.fastq --provider zenodo

# With metadata
cicada doi publish file.fastq \
  --enrich metadata.yaml \
  --provider zenodo

# Using DataCite instead
cicada doi publish file.fastq --provider datacite

# Test with sandbox
cicada doi publish file.fastq \
  --provider zenodo \
  --zenodo-sandbox
```

### Getting Help

```bash
# General help
cicada --help

# Help for specific command
cicada doi publish --help
cicada config --help
```

---

## Example: Complete Workflow

Here's a complete example from start to finish:

```bash
# 1. Set up credentials (one time)
cicada config init
# Enter your Zenodo token when prompted

# 2. Test that it works
cicada config test zenodo
# ✓ Zenodo authentication successful

# 3. Create metadata file
cat > metadata.yaml <<EOF
title: "My RNA-seq Data"
description: "Sequencing data from experiment 123"
creators:
  - name: "Jane Smith"
    affiliation: "University Lab"
    orcid: "0000-0002-1234-5678"
keywords:
  - RNA-seq
  - gene expression
rights: "Creative Commons Attribution 4.0 (CC-BY-4.0)"
EOF

# 4. Publish your data
cicada doi publish my-data.fastq \
  --enrich metadata.yaml \
  --provider zenodo

# Output:
# ✓ Metadata extracted from my-data.fastq
# ✓ Uploading to Zenodo...
# ✓ DOI created: 10.5281/zenodo.123456
# ✓ URL: https://zenodo.org/record/123456

# 5. Your data is now published with a DOI!
```

---

## Getting More Help

### Documentation
- Main README: https://github.com/scttfrdmn/cicada
- Detailed guides: https://github.com/scttfrdmn/cicada/docs

### Support
- GitHub Issues: https://github.com/scttfrdmn/cicada/issues
- Ask a question or report a problem

### Zenodo Help
- Zenodo FAQ: https://help.zenodo.org
- Zenodo Support: info@zenodo.org

### DataCite Help
- DataCite Support Guide: https://support.datacite.org
- Contact: support@datacite.org

---

## Glossary

**DOI (Digital Object Identifier)**: A permanent identifier for your data (like `10.5281/zenodo.123456`)

**Token**: A secret code that lets Cicada access your Zenodo account

**Credentials**: Your login information (tokens, passwords, etc.)

**Sandbox**: A test environment where you can practice without creating real DOIs

**Config file**: A file where Cicada stores your settings and credentials

**Repository ID**: Your institution's DataCite identifier

**Metadata**: Information about your data (title, authors, description, etc.)

**Provider**: The service that creates your DOI (Zenodo or DataCite)

---

**Questions?** Create an issue on GitHub: https://github.com/scttfrdmn/cicada/issues

**Ready to start?** Go back to [Part 1](#part-1-getting-your-zenodo-credentials-5-minutes) and follow the steps!
