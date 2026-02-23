# CableLabs Specification References

This document lists the CableLabs specifications relevant to the OpenDCI project
and explains how to obtain them.

## About `docs/external/`

The `docs/external/` directory is a git-ignored folder used to store local copies
of CableLabs specification PDFs. These PDFs are not committed to the repository
due to their size and licensing considerations. Each developer should download
the specs they need into this folder.

## Freely Available Specs

The following spec is freely downloadable without a CableLabs account:

| Spec ID | Version | Download URL |
|---------|---------|--------------|
| CM-SP-OSSIv3.0 | I05 (Dec 2007) | [Direct PDF](https://www.cablelabs.com/wp-content/uploads/2015/08/CM-SP-OSSIv3.0-I05-071206.pdf) |

Alternate CDN link:
https://www-res.cablelabs.com/wp-content/uploads/2019/02/28093824/CM-SP-OSSIv3.0-I05-071206.pdf

To download into `docs/external/`:

```sh
curl -L -o docs/external/CM-SP-OSSIv3.0-I05-071206.pdf \
  "https://www.cablelabs.com/wp-content/uploads/2015/08/CM-SP-OSSIv3.0-I05-071206.pdf"
```

## Account-Gated Specs (Free Registration Required)

The following current-version specifications require a free CableLabs account.
Register at: https://register.cablelabs.com

Once registered, you can download PDFs from the specification landing pages below.

| Spec ID | Version | Landing Page |
|---------|---------|--------------|
| CM-SP-MULPIv3.1 | I21 (Oct 2020) | https://www.cablelabs.com/specifications/CM-SP-MULPIv3.1 |
| CM-SP-MULPIv3.0 | I30 | https://www.cablelabs.com/specifications/CM-SP-MULPIv3.0 |
| CL-SP-CANN | ~I24 (Mar 2025) | https://www.cablelabs.com/specifications/CL-SP-CANN |
| CM-SP-CM-OSSIv3.1 | I20 (Oct 2020) | https://www.cablelabs.com/specifications/CM-SP-CM-OSSIv3.1 |

### How to Download Account-Gated Specs

1. Create a free account at https://register.cablelabs.com
2. Log in at https://www.cablelabs.com
3. Visit each specification landing page listed above
4. Download the PDF and save it into `docs/external/`

### Spec Descriptions

- **CM-SP-MULPIv3.1** -- DOCSIS 3.1 MAC and Upper Layer Protocols Interface.
  Defines the MAC-layer and upper-layer protocol requirements for DOCSIS 3.1
  cable modems and CMTSs.

- **CM-SP-MULPIv3.0** -- DOCSIS 3.0 MAC and Upper Layer Protocols Interface.
  The DOCSIS 3.0 predecessor to MULPIv3.1.

- **CL-SP-CANN** -- CableLabs Common Annex. Contains common definitions,
  encoding rules, and TLV tables shared across multiple DOCSIS specifications.
  This is a key reference for TLV parsing and configuration file formats.

- **CM-SP-CM-OSSIv3.1** -- DOCSIS 3.1 CM Operations Support System Interface.
  Defines the configuration file format, MIBs, and management interfaces for
  DOCSIS 3.1 cable modems.

- **CM-SP-OSSIv3.0** -- DOCSIS 3.0 Operations Support System Interface.
  The DOCSIS 3.0 predecessor to CM-OSSIv3.1. Freely available (see above).

## Expected Files in `docs/external/`

After downloading all specs, the `docs/external/` folder should contain:

```
docs/external/
  CM-SP-OSSIv3.0-I05-071206.pdf          (freely available)
  CM-SP-MULPIv3.1-*.pdf                  (requires account)
  CM-SP-MULPIv3.0-*.pdf                  (requires account)
  CL-SP-CANN-*.pdf                       (requires account)
  CM-SP-CM-OSSIv3.1-*.pdf                (requires account)
```

Note: The exact filenames for account-gated specs will depend on the version
downloaded from CableLabs. The `*` above represents the version/date suffix
that CableLabs includes in the filename.
