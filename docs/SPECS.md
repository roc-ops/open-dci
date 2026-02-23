# CableLabs Specification References

This document lists the CableLabs specifications relevant to the OpenDCI project
and explains how to obtain them.

## About `docs/external/`

The `docs/external/` directory is a git-ignored folder used to store local copies
of CableLabs specification PDFs. These PDFs are not committed to the repository
due to their size and licensing considerations. Each developer should download
the specs they need into this folder.

## Specifications

All current-version specifications require a free CableLabs account.
Register at: https://register.cablelabs.com

Once registered, you can download PDFs from the specification landing pages below.

### DOCSIS 4.0

| Spec ID | Current Version | Landing Page |
|---------|----------------|--------------|
| CM-SP-MULPIv4.0 | I11 (Feb 2026) | https://www.cablelabs.com/specifications/CM-SP-MULPIv4.0 |
| CM-SP-CM-OSSIv4.0 | I12 (Jun 2025) | https://www.cablelabs.com/specifications/CM-SP-CM-OSSIv4.0 |

### DOCSIS 3.1

| Spec ID | Current Version | Landing Page |
|---------|----------------|--------------|
| CM-SP-MULPIv3.1 | I25 (Apr 2023) | https://www.cablelabs.com/specifications/CM-SP-MULPIv3.1 |
| CM-SP-CM-OSSIv3.1 | I27 (Feb 2025) | https://www.cablelabs.com/specifications/CM-SP-CM-OSSIv3.1 |

### DOCSIS 3.0

| Spec ID | Current Version | Landing Page |
|---------|----------------|--------------|
| CM-SP-MULPIv3.0 | C01 (Dec 2017) | https://www.cablelabs.com/specifications/CM-SP-MULPIv3.0 |
| CM-SP-OSSIv3.0 | C01 (Dec 2017) | https://www.cablelabs.com/specifications/CM-SP-OSSIv3.0 |

### Common

| Spec ID | Current Version | Landing Page |
|---------|----------------|--------------|
| CL-SP-CANN | I24 (Mar 2025) | https://www.cablelabs.com/specifications/CL-SP-CANN |

### How to Download

1. Create a free account at https://register.cablelabs.com
2. Log in at https://www.cablelabs.com
3. Visit each specification landing page listed above
4. Download the PDF and save it into `docs/external/`

### Spec Descriptions

- **CM-SP-MULPIv4.0** -- DOCSIS 4.0 MAC and Upper Layer Protocols Interface.
  The latest generation defining config file structure for DOCSIS 4.0.

- **CM-SP-MULPIv3.1** -- DOCSIS 3.1 MAC and Upper Layer Protocols Interface.
  Defines the MAC-layer and upper-layer protocol requirements for DOCSIS 3.1
  cable modems and CMTSs.

- **CM-SP-MULPIv3.0** -- DOCSIS 3.0 MAC and Upper Layer Protocols Interface.
  The DOCSIS 3.0 predecessor to MULPIv3.1.

- **CL-SP-CANN** -- CableLabs Assigned Names and Numbers. Contains the central
  registry for all TLV values and sub-TLVs shared across DOCSIS specifications.
  This is a key reference for TLV parsing and configuration file formats.

- **CM-SP-CM-OSSIv4.0** -- DOCSIS 4.0 CM Operations Support System Interface.
  Defines the configuration file format, MIBs, and management interfaces for
  DOCSIS 4.0 cable modems.

- **CM-SP-CM-OSSIv3.1** -- DOCSIS 3.1 CM Operations Support System Interface.
  Defines the configuration file format, MIBs, and management interfaces for
  DOCSIS 3.1 cable modems.

- **CM-SP-OSSIv3.0** -- DOCSIS 3.0 Operations Support System Interface.
  The DOCSIS 3.0 predecessor to CM-OSSIv3.1.

## Expected Files in `docs/external/`

After downloading all specs, the `docs/external/` folder should contain:

```
docs/external/
  CM-SP-MULPIv4.0-I11-260219.pdf
  CM-SP-MULPIv3.1-I25-230419.pdf
  CM-SP-MULPIv3.0-C0I-171207.pdf
  CM-SP-CM-OSSIv4.0-I12-250611.pdf
  CM-SP-CM-OSSIv3.1-I27-250219.pdf
  CM-SP-OSSIv3.0-C01-171207.pdf
  CM-SP-OSSIv3.0-I05-071206.pdf
  CL-SP-CANN-I24-250320.pdf
```
