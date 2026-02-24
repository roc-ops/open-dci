# Vendor Extension Schemas

This directory contains optional vendor-specific TLV schema definitions. These schemas provide structured type information for vendor sub-TLVs that would otherwise be represented as generic `{ type, value }` hex pairs.

## File Naming Convention

Files are named using the 3-byte IEEE OUI in lowercase hex with no delimiters:

```
<OUI>.jtd.json
```

Examples:
- `001018.jtd.json` -- Broadcom (OUI 00:10:18)
- `00d09e.jtd.json` -- Cisco/Scientific Atlanta (OUI 00:D0:9E)

## Required Metadata

Every vendor schema file must include the following top-level metadata:

| Field | Type | Description |
|-------|------|-------------|
| `oui` | string | 3-byte OUI as a 6-character hex string (e.g., `"001018"`) |
| `vendorName` | string | Human-readable vendor name (e.g., `"Broadcom"`) |
| `extensionPoints` | array of strings | TLV paths where this vendor's sub-TLVs apply (e.g., `["43", "24.43"]`) |

## Definition Naming Convention

Definitions within a vendor schema should be prefixed with the vendor name to avoid collisions:

```
Broadcom_PowerSavingMode
Broadcom_SpectrumManagement
Cisco_CpeManagementMode
```

## Runtime Loading Behavior

Vendor schemas are loaded at runtime by tools that decode or encode DOCSIS configuration files. When a `VendorSpecificContainer` or `VendorSubTlvs` array is encountered:

1. The tool reads the `VendorId` (OUI) from the container or parent TLV 43 entry.
2. If a file matching `<OUI>.jtd.json` exists in this directory, the tool uses that schema to decode sub-TLV type numbers into named properties with typed values.
3. If no matching file exists, sub-TLVs remain in generic `{ type, value }` hex form.

## Important Notes

- The main schema (`docsis-config.jtd.json`) is never modified by vendor extensions. Vendor schemas are supplementary.
- Vendor schemas are optional. The system works without them -- all vendor sub-TLVs will simply appear as generic hex-encoded type/value pairs.
- Vendor schemas in this directory may be illustrative examples. Check each file's metadata for an `illustrative` flag.
