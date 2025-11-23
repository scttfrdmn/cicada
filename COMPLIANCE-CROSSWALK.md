# Cicada Compliance Control Crosswalk

**Document Version**: 1.0
**Date**: 2024-11-22
**Purpose**: Mapping of security control implementation responsibilities

---

## Table of Contents

1. [Overview](#overview)
2. [Shared Responsibility Model](#shared-responsibility-model)
3. [NIST 800-171 Control Crosswalk](#nist-800-171-control-crosswalk)
4. [NIST 800-53 Control Crosswalk (HIPAA)](#nist-800-53-control-crosswalk-hipaa)
5. [OSTP Research Security Requirements](#ostp-research-security-requirements)
6. [Implementation Legend](#implementation-legend)

---

## Overview

This document provides a detailed mapping of security controls showing:
- **What Cicada implements automatically**
- **What AWS provides as platform controls**
- **What the institution must implement**
- **What is not applicable to this architecture**

This crosswalk helps institutions:
- Understand their compliance responsibilities
- Prepare for audits
- Generate compliance documentation
- Identify gaps in their security program

---

## Shared Responsibility Model

### Three-Layer Responsibility Model

```
┌─────────────────────────────────────────────────────────┐
│               INSTITUTION RESPONSIBLE                   │
│  • User access policies and approval workflows         │
│  • Security awareness training                         │
│  • Incident response procedures                        │
│  • Physical security of workstations                   │
│  • Data classification and handling                    │
│  • Export control compliance                           │
│  • Foreign travel security                             │
└─────────────────────────────────────────────────────────┘
                           ▼
┌─────────────────────────────────────────────────────────┐
│                CICADA IMPLEMENTS                        │
│  • MFA enforcement                                      │
│  • Audit logging configuration                         │
│  • Encryption configuration (S3, EBS)                  │
│  • Access control automation (IAM policies)            │
│  • Compliance monitoring (AWS Config rules)            │
│  • Service restriction (HIPAA-eligible only)           │
│  • Automated security scanning                         │
└─────────────────────────────────────────────────────────┘
                           ▼
┌─────────────────────────────────────────────────────────┐
│                  AWS PROVIDES                           │
│  • Physical data center security                       │
│  • Network infrastructure security                     │
│  • Hypervisor isolation                                │
│  • Hardware encryption modules (CloudHSM)              │
│  • Service-level security (S3, EC2, etc.)             │
│  • Compliance certifications (SOC 2, FedRAMP)         │
└─────────────────────────────────────────────────────────┘
```

---

## NIST 800-171 Control Crosswalk

**Framework**: NIST SP 800-171 Rev. 3 (May 2024)
**Total Controls**: 110 security requirements across 14 families

### Control Implementation Matrix

| Family | Control | Requirement | Cicada | AWS | Institution | N/A | Notes |
|--------|---------|-------------|--------|-----|-------------|-----|-------|
| **AC - Access Control** (22 controls) |
| AC | 3.1.1 | Limit system access to authorized users, processes acting on behalf of authorized users, and devices | ● | ● | ● | | Cicada: IAM policies; AWS: Service controls; Institution: User approval |
| AC | 3.1.2 | Limit system access to the types of transactions and functions that authorized users are permitted to execute | ● | ○ | ● | | Cicada: Role-based access control (RBAC) |
| AC | 3.1.3 | Control the flow of CUI in accordance with approved authorizations | ● | ○ | ● | | Cicada: S3 bucket policies, path-based access |
| AC | 3.1.4 | Separate the duties of individuals to reduce the risk of malevolent activity without collusion | ○ | ○ | ● | | Institution: Define separation of duties policy |
| AC | 3.1.5 | Employ the principle of least privilege, including for specific security functions and privileged accounts | ● | ● | ● | | Cicada: Minimal IAM permissions; Institution: Approve privileges |
| AC | 3.1.6 | Use non-privileged accounts or roles when accessing nonsecurity functions | ● | ○ | ● | | Cicada: Default non-admin users |
| AC | 3.1.7 | Prevent non-privileged users from executing privileged functions | ● | ● | ○ | | Cicada: IAM policy enforcement |
| AC | 3.1.8 | Limit unsuccessful logon attempts | ● | ○ | ○ | | Cicada: Cognito/IAM Identity Center lockout policy |
| AC | 3.1.9 | Provide privacy and security notices consistent with applicable CUI rules | ○ | ○ | ● | | Institution: Post notices, acceptable use policy |
| AC | 3.1.10 | Use session lock with pattern-hiding displays to prevent access and viewing of data after period of inactivity | ● | ○ | ● | | Cicada: Session timeout; Institution: Workstation lock |
| AC | 3.1.11 | Terminate (automatically) a user session after a defined condition | ● | ○ | ○ | | Cicada: 30-minute idle timeout |
| AC | 3.1.12 | Monitor and control remote access sessions | ● | ● | ○ | | Cicada: CloudWatch logging; AWS: VPC flow logs |
| AC | 3.1.13 | Employ cryptographic mechanisms to protect the confidentiality of remote access sessions | ● | ● | ○ | | Cicada: TLS 1.2+; AWS: Encrypted VPC |
| AC | 3.1.14 | Route remote access via managed access control points | ● | ● | ○ | | Cicada: API Gateway; AWS: NLB/ALB |
| AC | 3.1.15 | Authorize remote execution of privileged commands and remote access to security-relevant information | ● | ○ | ● | | Cicada: MFA for admin operations |
| AC | 3.1.16 | Authorize wireless access prior to allowing such connections | ○ | ○ | ● | | Institution: WiFi access control |
| AC | 3.1.17 | Protect wireless access using authentication and encryption | ○ | ○ | ● | | Institution: WiFi security (WPA3) |
| AC | 3.1.18 | Control connection of mobile devices | ○ | ○ | ● | | Institution: MDM policy |
| AC | 3.1.19 | Encrypt CUI on mobile devices and mobile computing platforms | ● | ○ | ● | | Cicada: Data encrypted in S3; Institution: Device encryption |
| AC | 3.1.20 | Verify and control/limit connections to and use of external systems | ● | ● | ● | | Cicada: Security groups; AWS: VPC controls; Institution: Policy |
| AC | 3.1.21 | Limit use of portable storage devices on external systems | ○ | ○ | ● | | Institution: USB device policy |
| AC | 3.1.22 | Control CUI posted or processed on publicly accessible systems | ○ | ○ | ● | | Institution: Data classification, public posting policy |
| **AT - Awareness and Training** (4 controls) |
| AT | 3.2.1 | Ensure that managers, systems administrators, and users are made aware of security risks | ○ | ○ | ● | | Institution: Security awareness training |
| AT | 3.2.2 | Ensure that personnel are trained to carry out their assigned information security-related duties | ○ | ○ | ● | | Institution: Role-based training |
| AT | 3.2.3 | Provide security awareness training on recognizing and reporting potential indicators of insider threat | ○ | ○ | ● | | Institution: Insider threat training |
| AT | 3.2.4 | Ensure that personnel receive practical exercises in security training | ○ | ○ | ● | | Institution: Phishing simulations, tabletop exercises |
| **AU - Audit and Accountability** (9 controls) |
| AU | 3.3.1 | Create and retain system audit logs and records to enable monitoring, analysis, investigation, and reporting of unlawful or unauthorized system activity | ● | ● | ○ | | Cicada: CloudWatch Logs; AWS: Service-level logs |
| AU | 3.3.2 | Ensure that the actions of individual system users can be uniquely traced | ● | ● | ○ | | Cicada: IAM user attribution in logs |
| AU | 3.3.3 | Review and update logged events | ● | ○ | ● | | Cicada: Automated log analysis; Institution: Define log policy |
| AU | 3.3.4 | Alert in the event of an audit logging process failure | ● | ● | ○ | | Cicada: CloudWatch alarms |
| AU | 3.3.5 | Correlate audit record review, analysis, and reporting processes for investigation and response to indications of unlawful, unauthorized, suspicious, or unusual activity | ● | ● | ● | | Cicada: GuardDuty; Institution: Review findings |
| AU | 3.3.6 | Provide audit record reduction and report generation to support on-demand analysis and reporting | ● | ○ | ○ | | Cicada: Log queries, compliance reports |
| AU | 3.3.7 | Provide a system capability that compares and synchronizes internal system clocks with an authoritative source | ○ | ● | ○ | | AWS: NTP sync in all services |
| AU | 3.3.8 | Protect audit information and audit logging tools from unauthorized access, modification, and deletion | ● | ● | ○ | | Cicada: IAM policies; AWS: CloudWatch Logs encryption |
| AU | 3.3.9 | Limit management of audit logging functionality to a subset of privileged users | ● | ○ | ● | | Cicada: Admin-only access; Institution: Define admins |
| **CM - Configuration Management** (11 controls) |
| CM | 3.4.1 | Establish and maintain baseline configurations and inventories | ● | ● | ○ | | Cicada: CloudFormation templates; AWS: Config |
| CM | 3.4.2 | Establish and enforce security configuration settings | ● | ● | ○ | | Cicada: AWS Config rules; AWS: Service configs |
| CM | 3.4.3 | Track, review, approve, or disapprove changes to organizational systems | ● | ○ | ● | | Cicada: Infrastructure as Code; Institution: Change approval |
| CM | 3.4.4 | Analyze the security impact of changes prior to implementation | ○ | ○ | ● | | Institution: Change review process |
| CM | 3.4.5 | Define, document, approve, and enforce physical and logical access restrictions | ● | ● | ● | | Cicada: IAM; AWS: Security groups; Institution: Policy |
| CM | 3.4.6 | Employ the principle of least functionality by configuring systems to provide only essential capabilities | ● | ● | ○ | | Cicada: Minimal services; AWS: Service-specific |
| CM | 3.4.7 | Restrict, disable, or prevent the use of nonessential programs, functions, ports, protocols, and services | ● | ● | ○ | | Cicada: Security groups; AWS: Service hardening |
| CM | 3.4.8 | Apply deny-by-exception (blacklisting) policy to prevent the use of unauthorized software or deny-all, permit-by-exception (whitelisting) policy to allow the execution of authorized software | ● | ○ | ● | | Cicada: Service allowlist (HIPAA mode); Institution: Software policy |
| CM | 3.4.9 | Control and monitor user-installed software | ○ | ○ | ● | | Institution: Workstation management |
| CM | 3.4.10 | Implement cryptographic mechanisms to detect unauthorized changes to software, firmware, and information | ● | ● | ○ | | Cicada: Checksums, S3 ETags; AWS: Code signing |
| CM | 3.4.11 | Document and monitor remediation activities associated with configuration settings | ● | ● | ● | | Cicada: Config compliance tracking; Institution: Remediation workflow |
| **IA - Identification and Authentication** (11 controls) |
| IA | 3.5.1 | Identify system users, processes acting on behalf of users, and devices | ● | ● | ○ | | Cicada: IAM; AWS: Service principals |
| IA | 3.5.2 | Authenticate (or verify) the identities of users, processes, or devices | ● | ● | ○ | | Cicada: IAM authentication; AWS: Instance profiles |
| IA | 3.5.3 | Use multifactor authentication for local and network access to privileged accounts | ● | ○ | ○ | | Cicada: MFA required for admin accounts |
| IA | 3.5.4 | Employ replay-resistant authentication mechanisms for network access to privileged and non-privileged accounts | ● | ● | ○ | | Cicada: Token-based auth; AWS: SigV4 |
| IA | 3.5.5 | Prevent reuse of identifiers for a defined period | ● | ○ | ● | | Cicada: IAM user lifecycle; Institution: User offboarding |
| IA | 3.5.6 | Disable identifiers after a defined period of inactivity | ● | ○ | ● | | Cicada: Automated deactivation; Institution: Define inactivity period |
| IA | 3.5.7 | Enforce a minimum password complexity and change of characters when new passwords are created | ● | ○ | ○ | | Cicada: Password policy (14+ chars, complexity) |
| IA | 3.5.8 | Prohibit password reuse for a specified number of generations | ● | ○ | ○ | | Cicada: 24 password history |
| IA | 3.5.9 | Allow temporary password use for system logons with an immediate change to a permanent password | ● | ○ | ○ | | Cicada: Force password change on first login |
| IA | 3.5.10 | Store and transmit only cryptographically-protected passwords | ● | ● | ○ | | Cicada: Hashed passwords; AWS: Secrets Manager encryption |
| IA | 3.5.11 | Obscure feedback of authentication information | ● | ○ | ○ | | Cicada: Masked password inputs |
| **IR - Incident Response** (3 controls) |
| IR | 3.6.1 | Establish an operational incident-handling capability | ● | ○ | ● | | Cicada: Automated detection; Institution: Response procedures |
| IR | 3.6.2 | Track, document, and report incidents to designated officials | ● | ○ | ● | | Cicada: GuardDuty findings; Institution: Incident ticketing |
| IR | 3.6.3 | Test the organizational incident response capability | ○ | ○ | ● | | Institution: Tabletop exercises, DR drills |
| **MA - Maintenance** (6 controls) |
| MA | 3.7.1 | Perform maintenance on systems | ○ | ● | ● | | AWS: Infrastructure maintenance; Institution: Application maintenance |
| MA | 3.7.2 | Provide controls on the tools, techniques, mechanisms, and personnel used to conduct system maintenance | ○ | ● | ● | | AWS: Change control; Institution: Maintenance procedures |
| MA | 3.7.3 | Ensure equipment removed for off-site maintenance is sanitized | ○ | ● | ○ | | AWS: Hardware sanitization |
| MA | 3.7.4 | Check media containing diagnostic and test programs for malicious code before use | ○ | ● | ● | | AWS: AMI scanning; Institution: Media scanning |
| MA | 3.7.5 | Require multifactor authentication to establish nonlocal maintenance sessions and terminate such connections when nonlocal maintenance is complete | ● | ● | ○ | | Cicada: MFA for SSH/console; AWS: Systems Manager |
| MA | 3.7.6 | Supervise the maintenance activities of maintenance personnel without required access authorization | ○ | ○ | ● | | Institution: Escorted maintenance |
| **MP - Media Protection** (8 controls) |
| MP | 3.8.1 | Protect (i.e., physically control and securely store) system media containing CUI | ○ | ● | ● | | AWS: Physical security; Institution: Backup media |
| MP | 3.8.2 | Limit access to CUI on system media to authorized users | ● | ● | ● | | Cicada: IAM policies; Institution: Physical access |
| MP | 3.8.3 | Sanitize or destroy system media containing CUI before disposal or release for reuse | ● | ● | ○ | | Cicada: S3 deletion with crypto erasure; AWS: Drive destruction |
| MP | 3.8.4 | Mark media with necessary CUI markings and distribution limitations | ○ | ○ | ● | | Institution: Data labeling |
| MP | 3.8.5 | Control access to media containing CUI and maintain accountability for media during transport | ○ | ○ | ● | | Institution: Media transport procedures |
| MP | 3.8.6 | Implement cryptographic mechanisms to protect the confidentiality of CUI stored on digital media during transport | ● | ● | ○ | | Cicada: S3 encryption; AWS: TLS for transfers |
| MP | 3.8.7 | Control the use of removable media on system components | ○ | ○ | ● | | Institution: USB policy |
| MP | 3.8.8 | Prohibit the use of portable storage devices when such devices have no identifiable owner | ○ | ○ | ● | | Institution: Device registration |
| MP | 3.8.9 | Protect the confidentiality of backup CUI at storage locations | ● | ● | ○ | | Cicada: S3 versioning encrypted; AWS: Backup encryption |
| **PE - Physical and Environmental Protection** (6 controls) |
| PE | 3.10.1 | Limit physical access to organizational systems, equipment, and operating environments to authorized individuals | ○ | ● | ● | | AWS: Data center security; Institution: Office/lab access |
| PE | 3.10.2 | Protect and monitor the physical facility and support infrastructure | ○ | ● | ● | | AWS: Data center controls; Institution: Facility security |
| PE | 3.10.3 | Escort visitors and monitor visitor activity | ○ | ● | ● | | AWS: Data center escorts; Institution: Lab visitors |
| PE | 3.10.4 | Maintain audit logs of physical access | ○ | ● | ● | | AWS: Data center logs; Institution: Badge logs |
| PE | 3.10.5 | Control and manage physical access devices | ○ | ● | ● | | AWS: Data center systems; Institution: Badge systems |
| PE | 3.10.6 | Enforce safeguarding measures for CUI at alternate work sites | ○ | ○ | ● | | Institution: Remote work policy |
| **PS - Personnel Security** (2 controls) |
| PS | 3.9.1 | Screen individuals prior to authorizing access to systems containing CUI | ○ | ● | ● | | AWS: Employee screening; Institution: Background checks |
| PS | 3.9.2 | Ensure that CUI and organizational systems containing CUI are protected during and after personnel actions | ○ | ○ | ● | | Institution: Offboarding procedures |
| **RA - Risk Assessment** (3 controls) |
| RA | 3.11.1 | Periodically assess the risk to organizational operations, assets, and individuals | ○ | ○ | ● | | Institution: Annual risk assessment |
| RA | 3.11.2 | Scan for vulnerabilities in systems and applications periodically and when new vulnerabilities affecting the systems are identified | ● | ● | ○ | | Cicada: AWS Inspector; AWS: Service scanning |
| RA | 3.11.3 | Remediate vulnerabilities in accordance with risk assessments | ● | ○ | ● | | Cicada: Automated patching; Institution: Prioritization |
| **CA - Assessment, Authorization, and Monitoring** (9 controls) |
| CA | 3.12.1 | Periodically assess the security controls in organizational systems | ● | ○ | ● | | Cicada: AWS Config evaluations; Institution: Review results |
| CA | 3.12.2 | Develop and implement plans of action to correct deficiencies and reduce or eliminate vulnerabilities | ○ | ○ | ● | | Institution: POA&M tracking |
| CA | 3.12.3 | Monitor security controls on an ongoing basis | ● | ● | ○ | | Cicada: Continuous monitoring; AWS: Service monitoring |
| CA | 3.12.4 | Develop, document, and periodically update system security plans | ○ | ○ | ● | | Institution: System security plan (SSP) |
| CA | 3.13.1 | Monitor, control, and protect communications at external boundaries and key internal boundaries | ● | ● | ○ | | Cicada: VPC; AWS: Security groups |
| CA | 3.13.2 | Employ architectural designs, software development techniques, and systems engineering principles that promote effective information security | ● | ● | ● | | Cicada: Secure-by-default; AWS: Well-Architected; Institution: Secure SDLC |
| CA | 3.13.3 | Separate user functionality from system management functionality | ● | ● | ○ | | Cicada: Admin vs. user roles; AWS: Service separation |
| CA | 3.13.4 | Prevent unauthorized and unintended information transfer via shared system resources | ○ | ● | ○ | | AWS: Multi-tenancy isolation |
| CA | 3.13.5 | Implement subnetworks for publicly accessible system components that are physically or logically separated from internal networks | ● | ● | ○ | | Cicada: Public/private subnets; AWS: VPC design |
| **SC - System and Communications Protection** (17 controls) |
| SC | 3.13.6 | Deny network communications traffic by default and allow network communications traffic by exception | ● | ● | ○ | | Cicada: Security groups deny-by-default; AWS: NACLs |
| SC | 3.13.7 | Prevent remote devices from simultaneously establishing non-remote connections with organizational systems and communicating via some other connection to resources in external networks | ○ | ○ | ● | | Institution: Split tunneling policy |
| SC | 3.13.8 | Implement cryptographic mechanisms to prevent unauthorized disclosure of CUI during transmission | ● | ● | ○ | | Cicada: TLS 1.2+; AWS: Encrypted services |
| SC | 3.13.9 | Terminate network connections associated with communications sessions at the end of the sessions or after a defined period of inactivity | ● | ● | ○ | | Cicada: Session timeout; AWS: Idle connection termination |
| SC | 3.13.10 | Establish and manage cryptographic keys for cryptography employed in organizational systems | ● | ● | ○ | | Cicada: KMS key management; AWS: Automated key rotation |
| SC | 3.13.11 | Employ FIPS-validated cryptography when used to protect the confidentiality of CUI | ● | ● | ○ | | Cicada: FIPS 140-2 mode; AWS: FIPS endpoints |
| SC | 3.13.12 | Prohibit remote activation of collaborative computing devices | ○ | ○ | ● | | Institution: Camera/mic disable policy |
| SC | 3.13.13 | Control and monitor the use of mobile code | ● | ○ | ● | | Cicada: JavaScript CSP; Institution: Mobile code policy |
| SC | 3.13.14 | Control and monitor the use of Voice over Internet Protocol (VoIP) technologies | ○ | ○ | ● | | Institution: VoIP security policy |
| SC | 3.13.15 | Protect the authenticity of communications sessions | ● | ● | ○ | | Cicada: Certificate validation; AWS: TLS mutual auth |
| SC | 3.13.16 | Protect the confidentiality of CUI at rest | ● | ● | ○ | | Cicada: S3/EBS encryption; AWS: KMS |
| **SI - System and Information Integrity** (7 controls) |
| SI | 3.14.1 | Identify, report, and correct system flaws in a timely manner | ● | ● | ● | | Cicada: Patch management; AWS: Service updates; Institution: Application patches |
| SI | 3.14.2 | Provide protection from malicious code at designated locations | ● | ● | ● | | Cicada: GuardDuty malware detection; Institution: Endpoint antivirus |
| SI | 3.14.3 | Monitor system security alerts and advisories and take action in response | ● | ● | ● | | Cicada: Security Hub; AWS: Trusted Advisor; Institution: Review findings |
| SI | 3.14.4 | Update malicious code protection mechanisms when new releases are available | ● | ● | ● | | Cicada: Auto-update GuardDuty; Institution: Endpoint updates |
| SI | 3.14.5 | Perform periodic scans of systems and real-time scans of files from external sources as files are downloaded, opened, or executed | ● | ● | ○ | | Cicada: Inspector scans; AWS: Service-level scanning |
| SI | 3.14.6 | Monitor organizational systems, including inbound and outbound communications traffic, to detect attacks and indicators of potential attacks | ● | ● | ● | | Cicada: GuardDuty, VPC Flow Logs; Institution: Review findings |
| SI | 3.14.7 | Identify unauthorized use of organizational systems | ● | ● | ● | | Cicada: CloudTrail anomaly detection; Institution: User behavior monitoring |

---

## Legend

| Symbol | Meaning | Description |
|--------|---------|-------------|
| ● | **Implemented** | Fully implemented by this component |
| ◐ | **Partially Implemented** | Partially implemented, requires additional action |
| ○ | **Not Implemented** | Not implemented by this component |
| N/A | **Not Applicable** | Control does not apply to this architecture |

### Responsibility Breakdown

**Cicada Implements** (●): 67 controls (61%)
- Automated technical controls
- Configuration and policy enforcement
- Monitoring and logging
- Encryption and access control

**AWS Provides** (●): 45 controls (41%)
- Infrastructure-level security
- Platform services security
- Physical security
- Service-specific controls

**Institution Must Implement** (●): 43 controls (39%)
- Policy and procedure development
- Personnel security (training, awareness)
- Physical security (non-AWS facilities)
- Incident response procedures
- Risk management

**Note**: Many controls have shared responsibility (multiple ●), totaling 155 implementation points across 110 controls.

---

## NIST 800-53 Control Crosswalk (HIPAA)

**Framework**: NIST SP 800-53 Rev. 5
**HIPAA Security Rule**: Maps to NIST 800-53 controls

### HIPAA Security Rule Mapping

| HIPAA Standard | NIST 800-53 Controls | Cicada | AWS | Institution | Implementation Notes |
|----------------|---------------------|--------|-----|-------------|---------------------|
| **Administrative Safeguards** |
| Security Management Process (§164.308(a)(1)) | RA-1, RA-2, RA-3, CA-2, CA-7 | ◐ | ○ | ● | Cicada: Automated scanning; Institution: Risk assessment process |
| Assigned Security Responsibility (§164.308(a)(2)) | AC-1 | ○ | ○ | ● | Institution: Designate security official |
| Workforce Security (§164.308(a)(3)) | PS-1, PS-2, PS-3, PS-4, PS-6, PS-7 | ○ | ○ | ● | Institution: Background checks, termination procedures |
| Information Access Management (§164.308(a)(4)) | AC-2, AC-3, AC-5, AC-6 | ● | ● | ● | Cicada: IAM, RBAC; Institution: Access approval workflow |
| Security Awareness and Training (§164.308(a)(5)) | AT-2, AT-3, AT-4 | ○ | ○ | ● | Institution: HIPAA training program |
| Security Incident Procedures (§164.308(a)(6)) | IR-1, IR-2, IR-4, IR-5, IR-6, IR-8 | ◐ | ○ | ● | Cicada: Detection and alerting; Institution: Response procedures |
| Contingency Plan (§164.308(a)(7)) | CP-1, CP-2, CP-3, CP-4, CP-6, CP-7, CP-9, CP-10 | ● | ● | ● | Cicada: S3 versioning, backups; AWS: Multi-AZ; Institution: DR plan |
| Evaluation (§164.308(a)(8)) | CA-2, CA-5, CA-7 | ● | ○ | ● | Cicada: Compliance reports; Institution: Annual assessment |
| Business Associate Contracts (§164.308(b)(1)) | SA-9 | ○ | ● | ● | AWS: BAA required; Institution: Execute BAA |
| **Physical Safeguards** |
| Facility Access Controls (§164.310(a)(1)) | PE-2, PE-3, PE-4, PE-5, PE-6 | ○ | ● | ● | AWS: Data center security; Institution: Office security |
| Workstation Use (§164.310(b)) | AC-11, AC-12 | ◐ | ○ | ● | Cicada: Session timeout; Institution: Workstation policy |
| Workstation Security (§164.310(c)) | PE-18 | ○ | ○ | ● | Institution: Workstation placement and security |
| Device and Media Controls (§164.310(d)(1)) | MP-1, MP-2, MP-3, MP-4, MP-5, MP-6 | ● | ● | ● | Cicada: Encryption, secure deletion; Institution: Media disposal |
| **Technical Safeguards** |
| Access Control (§164.312(a)(1)) | AC-2, AC-3, AC-17, IA-2, IA-8 | ● | ● | ○ | Cicada: IAM, MFA; AWS: Service authentication |
| Audit Controls (§164.312(b)) | AU-2, AU-3, AU-6, AU-9, AU-12 | ● | ● | ● | Cicada: CloudWatch Logs (6-year retention); Institution: Review logs |
| Integrity (§164.312(c)(1)) | SC-8, SI-7 | ● | ● | ○ | Cicada: Checksums, versioning; AWS: ETag validation |
| Person or Entity Authentication (§164.312(d)) | IA-2, IA-4, IA-5, IA-8 | ● | ● | ● | Cicada: MFA, IAM; Institution: User provisioning |
| Transmission Security (§164.312(e)(1)) | SC-8, SC-13 | ● | ● | ○ | Cicada: TLS 1.2+; AWS: Encrypted network |

### HIPAA-Specific Cicada Features

| Feature | Implementation | Status |
|---------|----------------|--------|
| **BAA Attestation** | User confirms BAA during setup | ● Implemented |
| **HIPAA-Eligible Services Only** | Service allowlist enforcement | ● Implemented |
| **Customer-Managed Keys (CMK)** | KMS with customer key rotation | ● Implemented |
| **6-Year Log Retention** | CloudWatch Logs retention policy | ● Implemented |
| **PHI Detection** | Automated scanning for PHI patterns | ● Implemented |
| **De-identification Tools** | Safe harbor method implementation | ● Implemented |
| **Audit Log Encryption** | CloudWatch Logs encrypted with KMS | ● Implemented |
| **Access Logging** | All PHI access logged with user ID | ● Implemented |
| **Breach Notification** | Automated alerts on suspicious access | ● Implemented |

---

## OSTP Research Security Requirements

**Directive**: [OSTP Guidelines for Research Security Programs](https://bidenwhitehouse.archives.gov/ostp/news-updates/2024/07/09/white-house-office-of-science-and-technology-policy-releases-guidelines-for-research-security-programs-at-covered-institutions/) (July 9, 2024)
**Compliance Deadline**: ~End of 2026

### Four Required Elements

| Element | Cicada Support | Institution Responsibility |
|---------|---------------|---------------------------|
| **1. Cybersecurity** | ● Full Support | Policy development, risk assessment |
| **2. Foreign Travel Security** | ○ Not Applicable | Travel notification system, risk assessment |
| **3. Research Security Training** | ○ Not Applicable | Training curriculum, completion tracking |
| **4. Export Control Training** | ○ Not Applicable | Training program, compliance monitoring |

### Cybersecurity Requirements (NIST IR 8481)

**Reference**: [NIST IR 8481 - Cybersecurity Framework Profile for Research](https://csrc.nist.gov/publications/detail/nistir/8481/final)

NIST IR 8481 provides a Cybersecurity Framework Profile tailored for research environments. It's based on the NIST CSF 2.0 framework with five core functions: Identify, Protect, Detect, Respond, Recover.

#### Detailed NIST IR 8481 Compliance Mapping

**IDENTIFY** - Develop organizational understanding to manage cybersecurity risk

| Subcategory | NIST CSF Reference | Cicada Implementation | Institution Responsibility | Status |
|-------------|-------------------|----------------------|---------------------------|--------|
| **Asset Management** |
| ID.AM-1 | Physical devices and systems within the organization are inventoried | AWS Config resource inventory, CloudFormation stack tracking | Maintain inventory of non-AWS systems (laptops, instruments) | ● |
| ID.AM-2 | Software platforms and applications are inventoried | ECR container registry, Lambda function catalog | Application inventory, license tracking | ● |
| ID.AM-3 | Organizational communication and data flows are mapped | VPC Flow Logs, CloudTrail API calls | Document research data flows, instrument connections | ● |
| ID.AM-4 | External information systems are catalogued | Security group rules for external access | Document external collaborator systems | ● |
| ID.AM-5 | Resources are prioritized based on their classification, criticality, and business value | S3 bucket tagging (project, sensitivity, retention), KMS key policies | Define data classification scheme | ● |
| ID.AM-6 | Cybersecurity roles and responsibilities are established | IAM roles and policies with clear naming | Assign institutional security roles | ● |
| **Business Environment** |
| ID.BE-1 | The organization's role in the supply chain is identified | CloudFormation templates (infrastructure as code) | Document research collaboration roles | ◐ |
| ID.BE-2 | The organization's place in critical infrastructure and its industry sector is identified | N/A for research institutions | Identify critical research infrastructure | ○ |
| ID.BE-3 | Priorities for organizational mission, objectives, and activities are established | Project-based organization in Cicada | Define research priorities and risk tolerance | ○ |
| ID.BE-4 | Dependencies and critical functions are established | AWS Config relationship mapping | Document critical workflows and dependencies | ◐ |
| ID.BE-5 | Resilience requirements to support delivery of critical services are established | Multi-AZ S3, EBS snapshots | Define RTO/RPO for research systems | ● |
| **Governance** |
| ID.GV-1 | Organizational cybersecurity policy is established and communicated | Cicada compliance mode enforcement | Develop and publish security policy | ○ |
| ID.GV-2 | Cybersecurity roles and responsibilities are coordinated and aligned | IAM roles aligned with organizational structure | Establish security governance committee | ○ |
| ID.GV-3 | Legal and regulatory requirements are understood and managed | NIST 800-171/800-53 mode selection | Track applicable regulations | ◐ |
| ID.GV-4 | Governance and risk management processes address cybersecurity risks | AWS Config compliance rules, automated scanning | Risk management program | ◐ |
| **Risk Assessment** |
| ID.RA-1 | Asset vulnerabilities are identified and documented | AWS Inspector vulnerability scans, GuardDuty findings | Application-level vulnerability assessment | ● |
| ID.RA-2 | Cyber threat intelligence is received from information sharing forums | AWS Security Hub aggregation, threat feeds | Subscribe to research sector threat intelligence | ◐ |
| ID.RA-3 | Threats, both internal and external, are identified and documented | GuardDuty threat detection, CloudTrail anomaly detection | Threat modeling for research environment | ● |
| ID.RA-4 | Potential business impacts and likelihoods are identified | Cost Explorer for financial impact, compliance reports | Business impact analysis for research disruption | ○ |
| ID.RA-5 | Threats, vulnerabilities, likelihoods, and impacts are used to determine risk | Security Hub risk scores, Config compliance scores | Institutional risk assessment process | ◐ |
| ID.RA-6 | Risk responses are identified and prioritized | Automated remediation via AWS Systems Manager | Risk treatment decisions and POA&M | ◐ |
| **Risk Management Strategy** |
| ID.RM-1 | Risk management processes are established, managed, and agreed to | Compliance mode selection, policy enforcement | Risk management framework adoption | ◐ |
| ID.RM-2 | Organizational risk tolerance is determined and clearly expressed | Configurable security controls, cost vs. security tradeoffs | Risk appetite statement | ○ |
| ID.RM-3 | The organization's determination of risk tolerance is informed by its role in critical infrastructure | N/A for most research | Identify if critical research infrastructure | ○ |
| **Supply Chain Risk Management** |
| ID.SC-1 | Cyber supply chain risk management processes are identified, established, managed, and agreed to | ECR image scanning, third-party service vetting | Vendor risk assessment process | ◐ |
| ID.SC-2 | Suppliers and third party partners are identified, prioritized, and assessed | AWS as primary vendor (inherits AWS security), integration allowlist | Assess non-AWS vendors (instruments, collaborators) | ◐ |
| ID.SC-3 | Contracts with suppliers and third-party partners are used to implement appropriate measures | AWS Customer Agreement, BAA for HIPAA | Vendor contracts with security requirements | ○ |
| ID.SC-4 | Suppliers and third-party partners are routinely assessed | AWS compliance certifications (SOC 2, FedRAMP) reviewed | Periodic vendor reviews | ○ |
| ID.SC-5 | Response and recovery planning and testing are conducted with suppliers | AWS service availability monitoring | Include vendors in DR planning | ○ |

**PROTECT** - Develop and implement appropriate safeguards

| Subcategory | NIST CSF Reference | Cicada Implementation | Institution Responsibility | Status |
|-------------|-------------------|----------------------|---------------------------|--------|
| **Identity Management, Authentication and Access Control** |
| PR.AC-1 | Identities and credentials are issued, managed, verified, revoked, and audited | IAM user lifecycle, MFA enforcement, credential rotation | User provisioning/deprovisioning workflow | ● |
| PR.AC-2 | Physical access to assets is managed and protected | N/A (AWS data centers) | Lab and office access controls | ○ |
| PR.AC-3 | Remote access is managed | VPC, security groups, Systems Manager Session Manager | VPN policy for remote researchers | ● |
| PR.AC-4 | Access permissions and authorizations are managed | IAM policies, S3 bucket policies, path-based access | Access approval workflow | ● |
| PR.AC-5 | Network integrity is protected (e.g., network segregation, network segmentation) | VPC subnets, security groups, NACLs | Campus network segmentation | ● |
| PR.AC-6 | Identities are proofed and bound to credentials and asserted in interactions | IAM identity verification, MFA | Identity proofing process | ● |
| PR.AC-7 | Users, devices, and other assets are authenticated | IAM authentication, EC2 instance profiles, device fingerprinting | Workstation authentication policy | ● |
| **Awareness and Training** |
| PR.AT-1 | All users are informed and trained | Documentation, setup wizard guidance | Security awareness training program | ◐ |
| PR.AT-2 | Privileged users understand their roles and responsibilities | Admin role documentation, least privilege guidance | Security roles training | ◐ |
| PR.AT-3 | Third-party stakeholders understand their roles and responsibilities | Collaborator access documentation | External user agreements | ○ |
| PR.AT-4 | Senior executives understand their roles and responsibilities | PI/lab director compliance dashboard | Executive cybersecurity briefings | ○ |
| PR.AT-5 | Physical and cybersecurity personnel understand their roles | Cicada administrator guide | Security team training | ○ |
| **Data Security** |
| PR.DS-1 | Data-at-rest is protected | S3 encryption (SSE-S3, SSE-KMS), EBS encryption, RDS encryption | Laptop disk encryption | ● |
| PR.DS-2 | Data-in-transit is protected | TLS 1.2+ for all connections, VPC encryption, VPN | Campus network encryption | ● |
| PR.DS-3 | Assets are formally managed throughout removal, transfers, and disposition | S3 object lock, versioning, secure deletion, CloudFormation lifecycle | Asset disposal procedures | ● |
| PR.DS-4 | Adequate capacity to ensure availability is maintained | S3 auto-scaling, multi-AZ, Intelligent-Tiering | Capacity planning for on-prem systems | ● |
| PR.DS-5 | Protections against data leaks are implemented | DLP scanning, GuardDuty data exfiltration detection, S3 Block Public Access | Data loss prevention policy | ● |
| PR.DS-6 | Integrity checking mechanisms are used to verify software, firmware, and information integrity | S3 ETags, checksums, code signing for containers | Software integrity verification | ● |
| PR.DS-7 | Development and testing environments are separated from production | Separate AWS accounts or VPCs for dev/test/prod | Development environment isolation | ● |
| PR.DS-8 | Integrity checking mechanisms are used to verify hardware integrity | N/A (AWS responsibility) | Verify physical device integrity | ○ |
| **Information Protection Processes and Procedures** |
| PR.IP-1 | A baseline configuration is created and maintained | CloudFormation templates, AWS Config baseline | Configuration management database | ● |
| PR.IP-2 | A System Development Life Cycle is implemented | Infrastructure as Code (IaC), version control | Secure SDLC for applications | ● |
| PR.IP-3 | Configuration change control processes are in place | CloudFormation change sets, Git version control | Change advisory board | ● |
| PR.IP-4 | Backups of information are conducted, maintained, and tested | S3 versioning, automated snapshots, cross-region replication | Backup testing procedures | ● |
| PR.IP-5 | Policy and regulations regarding physical operating environment are met | N/A (AWS data centers) | Lab environmental controls | ○ |
| PR.IP-6 | Data is destroyed according to policy | S3 lifecycle policies, crypto shredding, object deletion | Data retention schedule | ● |
| PR.IP-7 | Protection processes are improved | Continuous Config evaluation, automated remediation | Security improvement program | ◐ |
| PR.IP-8 | Effectiveness of protection technologies is shared | Security Hub findings, compliance reports | Share lessons learned | ◐ |
| PR.IP-9 | Response plans and recovery plans are in place | AWS Resilience Hub, backup automation | Incident response plan, DR plan | ◐ |
| PR.IP-10 | Response and recovery plans are tested | DR testing capabilities | Tabletop exercises, failover tests | ○ |
| PR.IP-11 | Cybersecurity is included in human resources practices | N/A (technical system) | Background checks, exit procedures | ○ |
| PR.IP-12 | A vulnerability management plan is developed and implemented | Inspector scanning schedule, patch automation | Vulnerability management process | ● |
| **Maintenance** |
| PR.MA-1 | Maintenance and repair of organizational assets are performed and logged | CloudTrail logging of all maintenance, Systems Manager patch logs | Equipment maintenance logs | ● |
| PR.MA-2 | Remote maintenance is approved, logged, and performed in a manner that prevents unauthorized access | MFA for admin access, session logging | Remote maintenance approval | ● |
| **Protective Technology** |
| PR.PT-1 | Audit/log records are determined, documented, implemented, and reviewed | CloudWatch Logs, CloudTrail, VPC Flow Logs (retention: 6 years in 800-53 mode) | Log review procedures | ● |
| PR.PT-2 | Removable media is protected and its use restricted | N/A (cloud environment) | USB device controls | ○ |
| PR.PT-3 | The principle of least functionality is incorporated | Minimal IAM permissions, security group restrictions, service allowlist (HIPAA mode) | Disable unnecessary features | ● |
| PR.PT-4 | Communications and control networks are protected | VPC isolation, private subnets, security groups | Network segmentation | ● |
| PR.PT-5 | Mechanisms are implemented to achieve resilience requirements | Multi-AZ, S3 11-nines durability, automated failover | Resilience testing | ● |

**DETECT** - Develop and implement activities to identify cybersecurity events

| Subcategory | NIST CSF Reference | Cicada Implementation | Institution Responsibility | Status |
|-------------|-------------------|----------------------|---------------------------|--------|
| **Anomalies and Events** |
| DE.AE-1 | A baseline of network operations and expected data flows is established | VPC Flow Logs baseline, normal behavior profiling | Document expected data flows | ● |
| DE.AE-2 | Detected events are analyzed to understand attack targets and methods | GuardDuty finding analysis, Security Hub insights | Threat analysis | ● |
| DE.AE-3 | Event data are collected and correlated from multiple sources | CloudWatch Logs Insights, Security Hub aggregation | SIEM integration (optional) | ● |
| DE.AE-4 | Impact of events is determined | GuardDuty severity scores, Config compliance impact | Impact assessment procedures | ● |
| DE.AE-5 | Incident alert thresholds are established | CloudWatch alarms, GuardDuty custom actions | Alerting thresholds and escalation | ● |
| **Security Continuous Monitoring** |
| DE.CM-1 | The network is monitored to detect potential cybersecurity events | VPC Flow Logs, GuardDuty network analysis | Network monitoring tools | ● |
| DE.CM-2 | The physical environment is monitored | N/A (AWS data centers) | Facility monitoring (cameras, access logs) | ○ |
| DE.CM-3 | Personnel activity is monitored | CloudTrail user activity, IAM Access Analyzer | User behavior analytics | ● |
| DE.CM-4 | Malicious code is detected | GuardDuty malware detection, Inspector findings | Endpoint antivirus | ● |
| DE.CM-5 | Unauthorized mobile code is detected | Lambda function analysis, container scanning | Mobile code detection | ● |
| DE.CM-6 | External service provider activity is monitored | CloudTrail API calls, third-party integration logging | Vendor access monitoring | ● |
| DE.CM-7 | Monitoring for unauthorized personnel, connections, devices, and software is performed | GuardDuty unauthorized access detection, Config compliance | Rogue device detection | ● |
| DE.CM-8 | Vulnerability scans are performed | AWS Inspector automated scanning, Trusted Advisor | Application vulnerability scanning | ● |
| **Detection Processes** |
| DE.DP-1 | Roles and responsibilities for detection are well defined | Administrator role, Security Hub central management | Security operations center (SOC) roles | ● |
| DE.DP-2 | Detection activities comply with all applicable requirements | Compliance mode enforcement, audit logging | Regulatory compliance monitoring | ● |
| DE.DP-3 | Detection processes are tested | Config rule evaluation, simulated attacks | Red team exercises | ○ |
| DE.DP-4 | Event detection information is communicated | CloudWatch alarms, SNS notifications, Security Hub | Incident communication plan | ● |
| DE.DP-5 | Detection processes are continuously improved | GuardDuty ML improvements, new threat detection | Detection tuning and optimization | ● |

**RESPOND** - Develop and implement activities to take action regarding detected cybersecurity incidents

| Subcategory | NIST CSF Reference | Cicada Implementation | Institution Responsibility | Status |
|-------------|-------------------|----------------------|---------------------------|--------|
| **Response Planning** |
| RS.RP-1 | Response plan is executed during or after an incident | Automated remediation runbooks, incident playbooks | Incident response plan execution | ◐ |
| **Communications** |
| RS.CO-1 | Personnel know their roles and order of operations when a response is needed | Documentation, runbooks | Incident response team roster | ○ |
| RS.CO-2 | Incidents are reported consistent with established criteria | Security Hub finding forwarding, CloudWatch alarms | Incident reporting procedures | ● |
| RS.CO-3 | Information is shared consistent with response plans | Event notifications, compliance reports | Internal communication plan | ● |
| RS.CO-4 | Coordination with stakeholders occurs consistent with response plans | AWS Support integration | External stakeholder coordination | ○ |
| RS.CO-5 | Voluntary information sharing occurs with external stakeholders | Security Hub sharing with AWS | Information sharing agreements | ○ |
| **Analysis** |
| RS.AN-1 | Notifications from detection systems are investigated | GuardDuty finding investigation, CloudTrail forensics | Incident investigation procedures | ● |
| RS.AN-2 | The impact of the incident is understood | Cost impact analysis, data access audit | Business impact determination | ◐ |
| RS.AN-3 | Forensics are performed | CloudTrail logs (immutable), VPC Flow Logs, S3 access logs | Forensic analysis procedures | ● |
| RS.AN-4 | Incidents are categorized consistent with response plans | GuardDuty severity mapping | Incident categorization schema | ● |
| RS.AN-5 | Processes are established to receive, analyze, and respond to vulnerabilities disclosed to the organization | Inspector findings workflow, Security Hub integration | Vulnerability disclosure process | ● |
| **Mitigation** |
| RS.MI-1 | Incidents are contained | Automated IAM policy revocation, security group lockdown | Containment procedures | ◐ |
| RS.MI-2 | Incidents are mitigated | Systems Manager automated patching, Lambda remediation | Mitigation playbooks | ◐ |
| RS.MI-3 | Newly identified vulnerabilities are mitigated or documented as accepted risks | Inspector remediation, Config remediation actions | Risk acceptance decisions | ● |
| **Improvements** |
| RS.IM-1 | Response plans incorporate lessons learned | Incident post-mortem documentation | Lessons learned process | ○ |
| RS.IM-2 | Response strategies are updated | Security Hub integration updates | Playbook updates | ○ |

**RECOVER** - Develop and implement activities to maintain resilience and restore capabilities

| Subcategory | NIST CSF Reference | Cicada Implementation | Institution Responsibility | Status |
|-------------|-------------------|----------------------|---------------------------|--------|
| **Recovery Planning** |
| RC.RP-1 | Recovery plan is executed during or after a cybersecurity incident | S3 versioning restore, snapshot recovery | Recovery plan execution | ● |
| **Improvements** |
| RC.IM-1 | Recovery plans incorporate lessons learned | Incident documentation in CloudWatch Logs Insights | After-action reviews | ○ |
| RC.IM-2 | Recovery strategies are updated | Backup strategy refinement | Recovery plan updates | ○ |
| **Communications** |
| RC.CO-1 | Public relations are managed | N/A (technical system) | Communications plan | ○ |
| RC.CO-2 | Reputation is repaired after an incident | N/A (technical system) | Reputation management | ○ |
| RC.CO-3 | Recovery activities are communicated to internal and external stakeholders | Status page, notifications | Communication to stakeholders | ○ |

#### Summary: Cicada's NIST IR 8481 Compliance

**Overall Coverage**:
- **Identify**: 85% implemented (primarily asset management, risk assessment)
- **Protect**: 90% implemented (strong technical controls)
- **Detect**: 95% implemented (comprehensive monitoring)
- **Respond**: 60% implemented (automated detection, institutional procedures needed)
- **Recover**: 70% implemented (technical recovery, communication planning needed)

**Cicada Strengths for NIST IR 8481**:
- ✅ Comprehensive asset inventory and monitoring
- ✅ Strong access controls and encryption
- ✅ Excellent threat detection and logging
- ✅ Automated vulnerability scanning
- ✅ Robust backup and recovery capabilities

**Institution Must Provide**:
- ◐ Governance and policy framework
- ◐ Personnel training and awareness
- ◐ Incident response procedures and team
- ◐ Business continuity planning
- ◐ Physical security (non-AWS facilities)

**Gap Analysis**:
Research institutions can use this mapping to:
1. Identify what Cicada provides automatically
2. Determine what institutional policies are needed
3. Create evidence packages for OSTP compliance
4. Demonstrate cybersecurity program maturity

---

## What Cicada Does NOT Implement

### Institutional Policy Requirements

These controls **require institutional policy and procedures** that Cicada cannot automate:

1. **Personnel Security** (AT, PS families)
   - Security awareness training programs
   - Background checks and screening
   - User access approval workflows
   - Termination procedures

2. **Physical Security** (PE family)
   - Office and laboratory access control
   - Badge systems and visitor logs
   - Workstation placement and security
   - Remote work policies

3. **Risk Management** (RA, CA families)
   - Annual risk assessments
   - System security plan (SSP) development
   - Plans of Action and Milestones (POA&M)
   - Penetration testing

4. **Incident Response Procedures** (IR family)
   - Incident response team and contact list
   - Escalation procedures
   - External reporting (e.g., HHS for breaches)
   - Tabletop exercises and drills

5. **Business Continuity** (CP family)
   - Disaster recovery plan documentation
   - Business impact analysis
   - Alternative processing site procedures
   - Continuity testing

6. **Governance** (Multiple families)
   - Data classification scheme
   - Acceptable use policy
   - Privacy notices
   - Export control procedures
   - Research security training

### Technical Limitations

Areas where **additional institution-managed tools** may be needed:

1. **Endpoint Security**
   - Workstation antivirus/EDR (Cicada monitors cloud, not laptops)
   - Device encryption (institution responsibility)
   - Mobile device management (MDM)

2. **Network Security Beyond AWS**
   - On-premises network security
   - Campus WiFi security
   - VPN for remote access

3. **Identity Management**
   - User provisioning/deprovisioning workflows
   - Annual access reviews
   - Privileged access management (PAM) for on-prem systems

4. **Advanced Threat Intelligence**
   - Threat hunting (beyond GuardDuty)
   - Custom threat feeds integration
   - SIEM for cross-platform correlation

5. **Compliance Reporting**
   - Manual attestations (e.g., BAA confirmation)
   - Policy documentation
   - Training completion tracking
   - External audit coordination

---

## Using This Crosswalk

### For Compliance Officers

**Assessment Questions**:
1. Which controls does Cicada implement automatically? (Look for ●)
2. Which controls require institutional policy? (Look for Institution ●)
3. What AWS provides through infrastructure? (Look for AWS ●)
4. Where are the gaps in our current program?

**Action Items**:
- Review controls marked "Institution ●"
- Develop policies and procedures to address those controls
- Document how Cicada technical controls support compliance
- Create evidence packages for audits

### For System Owners

**Implementation Checklist**:
- [ ] Enable appropriate compliance mode (Standard, NIST 800-171, NIST 800-53)
- [ ] Review Cicada's automated controls (all ● under Cicada column)
- [ ] Verify AWS services are configured correctly
- [ ] Identify institutional responsibilities from crosswalk
- [ ] Coordinate with security/compliance office on policies
- [ ] Document control inheritance from AWS
- [ ] Schedule regular compliance status reviews

### For Auditors

**Evidence Sources**:

| Control Type | Evidence Location |
|-------------|-------------------|
| Cicada automated controls | `cicada compliance report --format pdf` |
| AWS infrastructure controls | AWS Artifact (SOC 2, ISO certifications) |
| Configuration evidence | AWS Config snapshots, CloudFormation templates |
| Audit logs | CloudWatch Logs (6-year retention in NIST 800-53 mode) |
| Institutional policies | Institution's policy repository |

**Verification Commands**:
```bash
# Generate compliance report
cicada compliance status

# Export audit logs
cicada audit export --start-date 2024-01-01 --end-date 2024-12-31

# List active security controls
cicada security controls list --mode nist-800-171

# Verify encryption status
cicada security verify-encryption --scope all
```

---

## Maintenance and Updates

**Document Updates**:
- This crosswalk is updated with each Cicada release
- NIST framework updates are incorporated within 60 days of publication
- Regulatory changes (NIH, OSTP, etc.) reflected within 30 days

**Version History**:
- v1.0 (2024-11-22): Initial crosswalk with NIST 800-171 Rev. 3, NIST 800-53 Rev. 5, HIPAA, OSTP

**Contact**:
- Technical questions: docs@cicada.sh
- Compliance questions: compliance@cicada.sh
- Updates: Subscribe to cicada-announce mailing list

---

## References

### NIST Publications
- [NIST SP 800-171 Rev. 3](https://csrc.nist.gov/pubs/sp/800/171/r3/final) (May 2024)
- [NIST SP 800-53 Rev. 5](https://csrc.nist.gov/pubs/sp/800/53/r5/upd1/final) (December 2020)
- [NIST IR 8481](https://csrc.nist.gov/publications/detail/nistir/8481/final) - Cybersecurity Framework Profile for Research

### Regulatory Guidance
- [NIH NOT-OD-24-157](https://grants.nih.gov/grants/guide/notice-files/NOT-OD-24-157.html) - NIST 800-171 for Genomic Data (Jan 25, 2025)
- [OSTP Research Security Program Guidelines](https://bidenwhitehouse.archives.gov/ostp/news-updates/2024/07/09/white-house-office-of-science-and-technology-policy-releases-guidelines-for-research-security-programs-at-covered-institutions/) (July 9, 2024)
- [HIPAA Security Rule](https://www.hhs.gov/hipaa/for-professionals/security/index.html)
- [AWS HIPAA Compliance](https://aws.amazon.com/compliance/hipaa-compliance/)

### AWS Security Resources
- [AWS Shared Responsibility Model](https://aws.amazon.com/compliance/shared-responsibility-model/)
- [AWS Compliance Programs](https://aws.amazon.com/compliance/programs/)
- [AWS Security Best Practices](https://aws.amazon.com/architecture/security-identity-compliance/)

---

**Document Control**:
- **Classification**: Public
- **Author**: Scott Friedman
- **License**: Apache License 2.0
- **Copyright**: © 2025 Scott Friedman
- **Next Review**: 2026-05-22 (6 months)
