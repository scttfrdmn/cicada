# Getting Started: Setting Up Your Credentials

**For researchers and lab managers with minimal technical experience**

This guide will help you set up Cicada to publish DOIs. No technical background required - just follow the steps!

## What You Need

To publish DOIs with Cicada, you need credentials (like a password) from either:
- **Zenodo** (free for everyone)
- **DataCite** (requires institutional membership)

**We recommend Zenodo** for most users because it's completely free and easy to set up.

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

## For DataCite Users

If your institution has a DataCite membership instead of using Zenodo:

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
