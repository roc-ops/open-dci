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

### DOCSIS Extensions

| Spec ID | Current Version | TLV Area | Landing Page |
|---------|----------------|----------|--------------|
| CM-SP-L2VPN | I17 (Oct 2025) | TLV 43.5, 45, 65 | https://www.cablelabs.com/specifications/business-services-over-docsis-layer-2-virtual-private-networks |
| CM-SP-eRouter | I22 (May 2024) | TLV 202 | https://www.cablelabs.com/specifications/CM-SP-eRouter |
| CM-SP-DSG | I25 (Sep 2017) | TLV 217 | https://www.cablelabs.com/specifications/CM-SP-DSG |
| CM-SP-eDOCSIS | I31 (Aug 2022) | TLV 201-231 | https://www.cablelabs.com/specifications/CM-SP-eDOCSIS |
| CM-SP-TEI | I06 (Jun 2010) | TLV 219 | https://www.cablelabs.com/specifications/CM-SP-TEI |
| CM-SP-SYNC | I03 (Jul 2022) | TLV 98-102 | https://www.cablelabs.com/specifications/CM-SP-SYNC |

### PacketCable

| Spec ID | Current Version | TLV Area | Landing Page |
|---------|----------------|----------|--------------|
| PKT-SP-PROV1.5 | C01 (Nov 2019) | TLV 216 (eMTA) | https://www.cablelabs.com/specifications/packetcable-mta-device-provisioning-specification |
| PKT-SP-RST-E-DVA | C01 (Mar 2014) | TLV 220 (eDVA) | https://www.cablelabs.com/specifications/PKT-SP-RST-E-DVA |

### DPoE

| Spec ID | Current Version | TLV Area | Landing Page |
|---------|----------------|----------|--------------|
| DPoE-SP-MULPIv1.0 | C01 (Aug 2016) | Various | https://www.cablelabs.com/specifications/DPoE-SP-MULPIv1.0 |
| DPoE-SP-MULPIv2.0 | I14 (Mar 2023) | TLV 43.5, 72, 83 | https://www.cablelabs.com/specifications/DPoE-SP-MULPIv2.0 |

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

- **CM-SP-L2VPN** -- Business Services over DOCSIS Layer 2 Virtual Private
  Networks. Defines L2VPN TLV encodings (TLV 43.5, 45, 65) for provisioning
  layer 2 VPN services over DOCSIS.

- **CM-SP-eRouter** -- eRouter Specification. Defines embedded router
  sub-TLVs (TLV 202) for TR-069, SNMP, and IPv4/IPv6 configuration.

- **CM-SP-DSG** -- DOCSIS Set-top Gateway Interface Specification. Defines
  DSG TLV encodings (TLV 217) for embedded set-top box (eSTB) provisioning.

- **CM-SP-eDOCSIS** -- Embedded DOCSIS Specification. Defines the eCM eSAFE
  configuration file TLVs (TLV 201-231) for embedded DOCSIS devices.

- **CM-SP-TEI** -- Business Services over DOCSIS TDM Emulation Interface.
  Defines TDM emulation TLV encodings (TLV 219) for eTEA provisioning.

- **CM-SP-SYNC** -- Synchronization Techniques for DOCSIS Technology. Defines
  DOCSIS timing synchronization TLVs (TLV 98-102).

- **PKT-SP-PROV1.5** -- PacketCable 1.5 MTA Device Provisioning. Defines
  eMTA configuration TLVs (TLV 216) for PacketCable voice services.

- **PKT-SP-RST-E-DVA** -- PacketCable 2.0 RST E-DVA Provisioning. Defines
  eDVA configuration TLVs (TLV 220) for PacketCable 2.0 devices.

- **DPoE-SP-MULPIv1.0** -- DPoE MAC and Upper Layer Protocols Interface v1.0.
  Defines TLV encodings for DOCSIS Provisioning of EPON (DPoE) v1.0 devices.

- **DPoE-SP-MULPIv2.0** -- DPoE MAC and Upper Layer Protocols Interface v2.0.
  Defines TLV encodings for DOCSIS Provisioning of EPON (DPoE) v2.0 devices.

## Expected Files in `docs/external/`

After downloading all specs, the `docs/external/` folder should contain:

```
docs/external/
  # Core DOCSIS
  CM-SP-MULPIv4.0-I11-260219.pdf
  CM-SP-MULPIv3.1-I25-230419.pdf
  CM-SP-MULPIv3.0-C0I-171207.pdf
  CM-SP-CM-OSSIv4.0-I12-250611.pdf
  CM-SP-CM-OSSIv3.1-I27-250219.pdf
  CM-SP-OSSIv3.0-C01-171207.pdf
  CM-SP-OSSIv3.0-I05-071206.pdf
  CL-SP-CANN-I24-250320.pdf

  # DOCSIS Extensions
  CM-SP-L2VPN-I17-251013.pdf
  CM-SP-eRouter-I22-240503.pdf
  CM-SP-DSG-I25-170906.pdf
  CM-SP-eDOCSIS-I31-220831.pdf
  CM-SP-TEI-I06-100611.pdf
  CM-SP-SYNC-I03-220715.pdf

  # PacketCable
  PKT-SP-PROV1.5-C01-191120.pdf
  PKT-SP-RST-E-DVA-C01-140314.pdf

  # DPoE
  DPoE-SP-MULPIv1.0-C01-160830.pdf
  DPoE-SP-MULPIv2.0-I14-230322.pdf
```
