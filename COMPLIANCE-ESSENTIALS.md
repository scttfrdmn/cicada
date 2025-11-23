# Compliance Essentials for Researchers

**What You Really Need to Know and Understand About Compliance**

---

## ⚠️ IMPORTANT DISCLAIMER

**THIS DOCUMENT IS NOT A REPLACEMENT FOR GUIDANCE FROM YOUR INSTITUTION'S CISO, COMPLIANCE OFFICER, OR LEGAL COUNSEL.**

This guide provides a practical overview of compliance requirements relevant to Cicada users. However:
- Your institution may have additional or more restrictive requirements
- Regulatory interpretations vary by institution and context
- When in doubt, **always consult your institution's compliance office**
- This is educational information, not legal or compliance advice

**When to talk to your CISO/Compliance Office**:
- Before handling controlled-access data (NIH genomic data, CUI, PHI)
- When setting up Cicada for the first time with sensitive data
- If you're unsure which compliance mode to use
- Before sharing data externally
- If you experience a potential security incident

---

## Table of Contents

1. [Compliance in Plain English](#compliance-in-plain-english)
2. [Do I Need to Care About Compliance?](#do-i-need-to-care-about-compliance)
3. [Three Common Compliance Scenarios](#three-common-compliance-scenarios)
4. [What Cicada Does (and Doesn't Do)](#what-cicada-does-and-doesnt-do)
5. [Your Responsibilities](#your-responsibilities)
6. [Common Misconceptions](#common-misconceptions)
7. [Quick Decision Tree](#quick-decision-tree)
8. [Practical Dos and Don'ts](#practical-dos-and-donts)
9. [What Happens During an Audit?](#what-happens-during-an-audit)
10. [Who to Contact](#who-to-contact)

---

## Compliance in Plain English

### What is "compliance"?

**Compliance** = Following the rules and standards that apply to your research data.

Different types of data have different rules:
- **Regular research data**: Basic security (passwords, backups)
- **NIH genomic data**: NIST 800-171 controls (110 security requirements)
- **Patient data (PHI)**: HIPAA rules mapped to NIST 800-53 (~900 controls)
- **Government contract data (CUI)**: NIST 800-171 controls

### Why do these rules exist?

1. **Protect sensitive information**: Patient privacy, national security, research integrity
2. **Prevent data breaches**: Avoid costly incidents and loss of trust
3. **Meet funding requirements**: NIH, DOD, NSF mandate compliance for certain data types
4. **Enable collaboration**: Standard security enables multi-institution research

### What does "being compliant" actually mean?

It means you can answer "yes" to these questions:
- ✅ Do you know what type of data you have?
- ✅ Are you following the rules for that data type?
- ✅ Can you prove it (documentation, logs, policies)?
- ✅ Do you have a plan if something goes wrong?

---

## Do I Need to Care About Compliance?

**Short answer**: It depends on your data.

### Flowchart: Does Compliance Apply to Me?

```
Do you work with any of these?
├─ NIH controlled-access genomic data (dbGaP, AnVIL, etc.)
│  └─ YES → You need NIST 800-171 (as of Jan 25, 2025)
│
├─ Patient health information (PHI)
│  └─ YES → You need HIPAA/NIST 800-53
│
├─ Government contract data marked "CUI"
│  └─ YES → You need NIST 800-171
│
├─ Export-controlled research (ITAR, EAR)
│  └─ YES → Special requirements (talk to export control office)
│
└─ None of the above, just regular research data
   └─ You still need basic security, but not strict compliance
```

### If you answered "NO" to everything above:

You still benefit from Cicada's **Standard mode**, which provides:
- ✅ Basic security (encryption, access control, MFA)
- ✅ Backups and versioning
- ✅ Audit logging
- ✅ OSTP research security alignment (due ~2026)

**No special compliance mode required**, but good security is still important!

---

## Three Common Compliance Scenarios

### Scenario 1: "I'm a genomics researcher using NIH dbGaP data"

**What you need**: NIST 800-171 mode

**Why**: NIH policy effective January 25, 2025 requires all users of controlled-access genomic data to comply with NIST 800-171.

**What this means for you**:
- **Before January 25, 2025**: Your institution must attest NIST 800-171 compliance to access NIH data
- **Cicada helps**: Enabling NIST 800-171 mode implements all 110 required controls automatically
- **You must still**: Follow your institution's data access policies, complete training, use MFA

**Action items**:
1. Talk to your institution's dbGaP coordinator
2. Enable Cicada NIST 800-171 mode: `cicada init --compliance-mode nist-800-171`
3. Complete required training
4. Follow data use agreement terms

**Key restriction**: You can **only** use dbGaP data on systems that are NIST 800-171 compliant

---

### Scenario 2: "I'm doing clinical research with patient data"

**What you need**: NIST 800-53 mode (HIPAA)

**Why**: Patient health information (PHI) is regulated by HIPAA. Cicada maps HIPAA requirements to NIST 800-53 controls.

**What this means for you**:
- **Your institution must have a Business Associate Agreement (BAA) with AWS** (this is separate from Cicada)
- **Cicada helps**: NIST 800-53 mode enforces HIPAA-eligible services only, required encryption, 6-year audit logs
- **You must still**: Complete HIPAA training, minimize PHI exposure, report breaches

**Action items**:
1. Confirm your institution has signed a BAA with AWS (talk to CISO/compliance office)
2. Complete HIPAA training (required!)
3. Enable Cicada NIST 800-53 mode: `cicada init --compliance-mode nist-800-53`
4. Only store PHI in Cicada (never in email, Dropbox, personal drives)
5. Report any suspected PHI exposure immediately

**Key restrictions**:
- **Cannot** use non-HIPAA services (Cicada blocks these automatically)
- **Must** use customer-managed encryption keys (Cicada enforces this)
- **Must** report breaches within 60 days (talk to compliance office!)

---

### Scenario 3: "I'm working on a DoD or federal contract"

**What you need**: Probably NIST 800-171, possibly more

**Why**: Federal contracts often require protection of Controlled Unclassified Information (CUI).

**What this means for you**:
- **Check your contract**: Look for CUI requirements, DFARS clauses, or NIST 800-171 mentions
- **Cicada helps**: NIST 800-171 mode implements required controls
- **You must still**: Follow contract terms, mark CUI properly, report incidents

**Action items**:
1. Review your grant/contract for data security requirements
2. Talk to your institution's contracts office and CISO
3. Enable Cicada NIST 800-171 mode if required
4. Mark data appropriately (CUI markings)
5. Understand reporting requirements (incidents must be reported to DCSA)

**Key restrictions**:
- CUI must be protected even on personal devices (use Cicada, not local storage)
- Incident reporting timelines are strict (72 hours for some incidents)
- Foreign national access may be restricted (check contract terms)

---

## What Cicada Does (and Doesn't Do)

### What Cicada DOES Automatically

When you enable a compliance mode, Cicada automatically:

**Technical Controls** ✅:
- Enforces encryption (data at rest and in transit)
- Requires multi-factor authentication (MFA)
- Logs all data access (audit trail)
- Scans for vulnerabilities
- Monitors for security threats
- Blocks unauthorized services (in HIPAA mode)
- Implements access controls (who can see what)
- Backs up data with versioning
- Detects anomalous behavior

**Compliance Reporting** ✅:
- Generates compliance status reports
- Tracks control implementation
- Provides audit evidence (logs, configurations)
- Monitors security posture
- Alerts on policy violations

**Example**: If you enable NIST 800-171 mode, Cicada implements all 110 security requirements automatically. You don't need to manually configure encryption, logging, or access controls.

### What Cicada DOES NOT Do (Your Responsibility)

**Policies and Procedures** ❌:
- Security awareness training for your team
- Data classification decisions
- Access approval workflows (who gets access?)
- Incident response procedures (what to do when things go wrong)
- Business continuity planning
- Privacy notices and user agreements

**Physical and Personnel Security** ❌:
- Background checks for lab personnel
- Physical access control to your lab/office
- Visitor management
- Workstation security (locking your laptop)
- Secure disposal of printed documents

**Governance** ❌:
- Risk assessments (what could go wrong?)
- Policy documentation (acceptable use, data retention)
- Compliance attestations (official sign-offs)
- Third-party vendor agreements
- Export control determinations

**Why the split?**: Cicada is a technical tool that implements technical controls. Your institution must provide the governance, policy, and human elements of compliance.

---

## Your Responsibilities

### As a Researcher Using Cicada

Even with Cicada handling technical controls, you must:

**1. Know Your Data** 🔍
- What type of data do you have?
- Does it require special protections?
- Who is allowed to access it?
- How long must you keep it?

**2. Follow Basic Security Hygiene** 🔒
- Use strong, unique passwords
- Enable MFA on all accounts
- Lock your computer when you step away
- Don't share credentials
- Report suspicious activity

**3. Complete Required Training** 📚
- HIPAA training (if handling PHI)
- Security awareness training (annual)
- Data handling training specific to your data type
- Export control training (if required)

**4. Follow Policies** 📋
- Your institution's acceptable use policy
- Data use agreements (DUAs) for controlled data
- Grant requirements
- Lab-specific data handling procedures

**5. Report Issues Promptly** 🚨
- Lost laptop or device
- Suspected data breach
- Accidental data exposure
- Unusual system behavior
- Compliance questions or concerns

**6. Document Your Work** 📝
- Keep records of data access
- Document analysis workflows
- Track data sharing (who, what, when)
- Note any incidents or near-misses

**7. Plan for Departures** 👋
- Transfer data ownership before leaving
- Document your analysis workflows
- Ensure data is properly stored in Cicada
- Return institutional credentials

---

## Common Misconceptions

### ❌ MYTH: "Cicada handles compliance, so I don't need to do anything"

**✅ REALITY**: Cicada handles technical controls, but you're responsible for following policies, training, and proper data handling.

**Example**: Cicada ensures data is encrypted, but you must decide who gets access and follow the approval process.

---

### ❌ MYTH: "HIPAA only applies to hospitals, not researchers"

**✅ REALITY**: Any research involving patient data (PHI) must comply with HIPAA, including academic research.

**Example**: A genetics study using patient samples requires HIPAA compliance even if conducted at a university.

---

### ❌ MYTH: "I can use standard Dropbox/Google Drive for NIH genomic data if I'm careful"

**✅ REALITY**: As of January 25, 2025, NIH requires NIST 800-171 compliance for genomic data. Consumer services don't meet this standard.

**Example**: dbGaP data must be stored in NIST 800-171 compliant systems like Cicada. Dropbox is not compliant.

---

### ❌ MYTH: "Compliance is just about checking boxes for auditors"

**✅ REALITY**: Compliance protects research participants, prevents data breaches, and enables collaboration.

**Example**: HIPAA de-identification rules protect patient privacy. Audit logs help identify the source of a breach.

---

### ❌ MYTH: "Once I set up compliance mode, I'm done forever"

**✅ REALITY**: Compliance is ongoing. Policies change, risks evolve, and you need to stay current.

**Example**: NIH updated genomic data requirements in 2024. Researchers needed to adapt by January 2025.

---

### ❌ MYTH: "My data isn't important enough for hackers to care about"

**✅ REALITY**: All data has value. Research data can be held for ransom, patient data can be sold, and credentials can be used for other attacks.

**Example**: University researchers are common targets because security is often less mature than enterprises.

---

### ❌ MYTH: "Compliance is expensive and slows down research"

**✅ REALITY**: Good security enables research by preventing incidents that would halt work for months. Cicada is designed to be low-cost and non-intrusive.

**Example**: A data breach investigation can shut down a lab for months. Prevention is cheaper than recovery.

---

## Quick Decision Tree

**"Which Cicada compliance mode should I use?"**

```
START: What type of data are you working with?

┌─ NIH controlled-access genomic data?
│  └─ YES → Use NIST 800-171 mode
│     └─ Command: cicada init --compliance-mode nist-800-171
│
├─ Patient health information (PHI)?
│  └─ YES → Use NIST 800-53 mode (HIPAA)
│     ├─ First: Confirm institution has AWS BAA
│     └─ Command: cicada init --compliance-mode nist-800-53
│
├─ Federal contract data marked "CUI"?
│  └─ YES → Use NIST 800-171 mode
│     └─ Command: cicada init --compliance-mode nist-800-171
│
└─ Regular research data (no special requirements)?
   └─ Use Standard mode (still secure, less restrictive)
      └─ Command: cicada init
```

**Still not sure?** → **Talk to your institution's CISO or compliance office**

---

## Practical Dos and Don'ts

### ✅ DO

1. **DO** enable the appropriate compliance mode when setting up Cicada
2. **DO** use multi-factor authentication (MFA) on all accounts
3. **DO** complete required training before accessing sensitive data
4. **DO** report security incidents immediately (even if you're not sure)
5. **DO** follow data use agreements (DUAs) to the letter
6. **DO** keep your software and systems updated
7. **DO** ask questions when you're unsure
8. **DO** document your data handling practices
9. **DO** use Cicada for sensitive data (not Dropbox, email, USB drives)
10. **DO** involve your CISO early in new projects with sensitive data

### ❌ DON'T

1. **DON'T** share credentials or let others use your account
2. **DON'T** disable security features because they're "inconvenient"
3. **DON'T** store sensitive data on personal devices or consumer services
4. **DON'T** email PHI or controlled data (use secure methods)
5. **DON'T** grant data access without following approval procedures
6. **DON'T** ignore security alerts or warnings
7. **DON'T** assume "just this once" is okay
8. **DON'T** wait to report a suspected incident
9. **DON'T** mix personal and research work on the same systems
10. **DON'T** skip training or attestations

---

## What Happens During an Audit?

### Types of Audits

**1. Institutional Self-Assessment**
- Your compliance office reviews systems periodically
- **You may be asked**: Show how you're protecting data, provide documentation

**2. Federal Agency Audit** (NIH, NSF, DOD)
- For grants/contracts, agencies audit periodically
- **You may be asked**: Demonstrate compliance with grant terms, show security controls

**3. HIPAA Audit** (HHS Office for Civil Rights)
- Random or triggered by complaints
- **You may be asked**: Prove HIPAA training, show BAA, demonstrate PHI protections

**4. Third-Party Assessment** (ISO, SOC 2, etc.)
- For institutional certifications
- **You may be asked**: Participate in interviews, provide evidence

### What Auditors Want to See

**Documentation** 📄:
- Policies and procedures
- Training completion records
- Data use agreements
- Access approval records
- Incident reports (if any)

**Technical Evidence** 💻:
- Audit logs (Cicada provides these)
- Encryption status (Cicada compliance report)
- Access controls (IAM policies)
- Vulnerability scan results (Cicada provides)
- Backup verification (Cicada provides)

**Processes** 🔄:
- How do you grant access?
- What happens when someone leaves?
- How do you handle incidents?
- Who reviews logs?

### How Cicada Helps During Audits

Cicada provides:
- ✅ **Compliance reports**: `cicada compliance report --format pdf`
- ✅ **Audit logs**: Comprehensive activity logs with 6-year retention (HIPAA mode)
- ✅ **Control evidence**: Proof that technical controls are in place
- ✅ **Configuration snapshots**: Documentation of security settings
- ✅ **Gap analysis**: Shows what Cicada handles vs. institutional responsibility

### Your Role in an Audit

**Before the audit**:
- Keep good records (who accessed what, when)
- Complete all required training
- Follow policies consistently
- Document any incidents or exceptions

**During the audit**:
- Be honest and cooperative
- Don't speculate – stick to facts
- Refer technical questions to Cicada compliance reports
- Refer policy questions to your compliance office
- Take notes on auditor feedback

**After the audit**:
- Address any findings promptly
- Implement recommended improvements
- Document corrective actions
- Update training and procedures as needed

---

## Who to Contact

### When Something Goes Wrong

**Security Incident (suspected breach, lost device, unauthorized access)**:
1. **Immediate**: Report to your institution's security incident response team
   - Often: IT Security, CISO office, or campus police
2. **Within hours**: Notify your PI/lab director
3. **As directed**: File official incident report
4. **For PHI breaches**: Institution must report to HHS within 60 days

**Questions About Compliance**:
- **Compliance Office**: Policy interpretation, audit preparation
- **CISO/IT Security**: Technical security requirements
- **Contracts Office**: Grant/contract requirements
- **IRB**: Human subjects research and data use
- **Export Control Office**: International collaboration, ITAR/EAR

**Questions About Cicada**:
- **Technical issues**: Cicada support (docs@cicada.sh)
- **Compliance mode selection**: Start with your CISO, then Cicada
- **Control implementation**: Cicada documentation or support

### Building a Compliance Team

For labs with ongoing compliance needs, establish relationships with:

**Internal**:
- CISO or IT Security representative
- Compliance officer
- Research data management librarian
- IRB coordinator (if applicable)
- Export control officer (if applicable)

**External**:
- NIH dbGaP help desk (for genomic data)
- AWS support (for BAA and infrastructure)
- Cicada support (for technical compliance)

**Pro tip**: Introduce your CISO to Cicada early. They'll appreciate being involved upfront rather than discovering compliance tools after the fact.

---

## Compliance Checklist for New Projects

Use this checklist when starting a new research project:

### Before You Begin

- [ ] **Data Classification**: What type of data will you collect?
- [ ] **Regulatory Review**: Does this require IRB, IACUC, IBC approval?
- [ ] **Compliance Determination**: Which compliance requirements apply?
- [ ] **Approval**: Have you received necessary approvals?
- [ ] **Training**: Have all team members completed required training?
- [ ] **DUAs**: Are data use agreements in place?
- [ ] **Funding**: Does the grant have specific security requirements?

### Setting Up Cicada

- [ ] **Compliance Mode**: Selected appropriate mode (Standard, NIST 800-171, NIST 800-53)
- [ ] **BAA**: If HIPAA, confirmed institution has AWS BAA
- [ ] **Access Control**: Defined who needs access and at what level
- [ ] **Data Organization**: Planned folder structure and naming conventions
- [ ] **Backup Strategy**: Configured retention policies
- [ ] **Monitoring**: Set up alerts and notifications

### Ongoing Operations

- [ ] **Quarterly**: Review access lists (remove departed members)
- [ ] **Quarterly**: Check compliance reports for issues
- [ ] **Annually**: Complete security awareness training
- [ ] **Annually**: Review and update policies
- [ ] **As needed**: Report incidents promptly
- [ ] **As needed**: Update Cicada when new features are released

### Before Project Ends

- [ ] **Data Archival**: Move data to long-term storage if needed
- [ ] **Access Removal**: Revoke access for departed team members
- [ ] **Documentation**: Archive data handling documentation
- [ ] **Public Data**: Share data per grant requirements (with DOI)
- [ ] **Closeout**: Final report to funders (if required)

---

## Additional Resources

### NIST Publications
- **NIST SP 800-171**: [Protecting CUI](https://csrc.nist.gov/pubs/sp/800/171/r3/final)
- **NIST SP 800-53**: [Security Controls](https://csrc.nist.gov/pubs/sp/800/53/r5/upd1/final)
- **NIST IR 8481**: [Research Cybersecurity](https://csrc.nist.gov/publications/detail/nistir/8481/final)

### Federal Guidance
- **NIH Data Sharing**: [NIH Genomic Data Sharing Policy](https://grants.nih.gov/grants/guide/notice-files/NOT-OD-24-157.html)
- **OSTP Research Security**: [Guidelines (July 2024)](https://bidenwhitehouse.archives.gov/ostp/news-updates/2024/07/09/white-house-office-of-science-and-technology-policy-releases-guidelines-for-research-security-programs-at-covered-institutions/)
- **HIPAA**: [HHS HIPAA Homepage](https://www.hhs.gov/hipaa/index.html)

### Cicada Documentation
- **Compliance Crosswalk**: [Detailed control mapping](COMPLIANCE-CROSSWALK.md)
- **Project Summary**: [Architecture and features](PROJECT-SUMMARY.md)
- **CLI Reference**: [Command documentation](docs/cli-reference.md)

### Training Resources
- **HIPAA Training**: Check with your institution's compliance office
- **CITI Program**: Research ethics and compliance training
- **NIST Training**: Free cybersecurity courses at [nist.gov](https://www.nist.gov/itl/applied-cybersecurity/nice/resources/online-learning-content)

---

## Key Takeaways

### The Bottom Line

1. **Compliance = Following rules for your data type**
   - Different data, different rules
   - Know what you have

2. **Cicada handles technical controls automatically**
   - Encryption, logging, monitoring, access control
   - Choose the right compliance mode

3. **You handle policy, training, and procedures**
   - Complete training
   - Follow policies
   - Report incidents

4. **When in doubt, ask your CISO**
   - They're there to help
   - Better to ask than to guess wrong

5. **Good security enables research**
   - Prevents incidents that would halt work
   - Builds trust with collaborators and funders
   - Protects research participants

### Remember

**Compliance is not about perfection** – it's about:
- Understanding requirements
- Making good-faith efforts
- Documenting your practices
- Improving continuously
- Reporting issues promptly

**You're not alone** – your institution has resources:
- CISO and IT security team
- Compliance office
- Research data services
- IRB (for human subjects)
- Export control office

**Cicada is a tool** – it implements technical controls, but:
- Your judgment matters
- Your institution's policies apply
- Your CISO has final say
- You're responsible for using it correctly

---

## Final Reminder

⚠️ **THIS DOCUMENT IS EDUCATIONAL, NOT LEGAL OR COMPLIANCE ADVICE**

Always consult with your institution's compliance office, CISO, or legal counsel for:
- Official compliance determinations
- Policy interpretations
- Incident response procedures
- Audit preparation
- Regulatory questions

When in doubt: **Ask first, act second.**

---

---

**Document Version**: 1.0
**Last Updated**: 2025-11-22
**Author**: Scott Friedman
**License**: Apache License 2.0
**Copyright**: © 2025 Scott Friedman

**Questions?**: Contact your institution's CISO

**Good luck with your research!** 🦗
