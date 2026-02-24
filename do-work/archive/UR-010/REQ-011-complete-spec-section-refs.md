---
id: REQ-011
title: Complete spec section references for all TLVs
status: done
created_at: 2026-02-24T08:00:00Z
claimed_at: 2026-02-24T08:30:00Z
route: C
user_request: UR-010
---

# Complete Spec Section References for All TLVs

## What
Add specific section numbers to all TLV spec references that currently only list the spec name. Ensure consistency across both docs/TLVs.md and the JTD schema's `metadata.spec` fields. Add a spec-to-document mapping table in TLVs.md.

## Detailed Requirements

### Consistent Section References
- Every TLV and sub-TLV must have a spec reference that includes the specific section number, not just the spec name
- Good example (current MULPI pattern): `"spec": "CM-SP-MULPIv4.0 C.1.2.26"` — has both spec document and section
- Bad example (current non-MULPI pattern): `"spec": "CM-SP-eDOCSIS-I31"` — missing section reference
- This applies to ALL specs: eDOCSIS, eRouter, L2VPN, TEI, PacketCable, DSG, DPoE-MULPI, SYNC, etc.

### Multi-Spec TLVs
- If a TLV spans multiple specs or is defined/referenced in multiple sections, ALL spec references and sections should be listed
- Example: a TLV defined in both MULPI and DPoE-MULPI should reference both with their sections

### Files to Update
- **docs/TLVs.md** — Spec and Reference columns in all tables
- **schemas/docsis-config.jtd.json** — `metadata.spec` field on every property and definition
- **schemas/generated/docsis-config.schema.json** — regenerated from updated JTD

### Spec-to-Document Mapping Table
- Add a table in docs/TLVs.md mapping spec short names to their full document identifiers
- Example: "MULPI" → "CM-SP-MULPIv4.0-I11-260219"
- This information is already available in docs/SPECS.md — use it as the source
- Helps readers find the actual document when they see a section reference

## Builder Guidance
- Certainty level: Firm — user explicitly wants section-level references for all specs, not just MULPI
- The spec PDFs in `docs/external/` are the authoritative source for section numbers
- SPECS.md has the spec-to-document mapping already
- For specs where the section number is not readily available from existing sources, research the spec PDFs in `docs/external/`
- User's key principle: "if there is a question on a particular TLV or sub-TLV we should reference the section in the appropriate spec document"

## Full Context
See [user-requests/UR-010/input.md](./user-requests/UR-010/input.md) for complete verbatim input.

## Triage

- **Route**: C (Complex) — requires researching 8+ spec PDFs to find section numbers for hundreds of TLV entries
- **Risk**: Low — metadata-only changes to spec fields, no functional schema changes
- **Scope**: JTD schema metadata.spec fields, TLVs.md Spec/Reference columns, spec-to-document mapping table

## Plan

1. Update JTD schema `metadata.spec` fields with section numbers for all non-MULPI specs
2. Update TLVs.md Spec and Reference columns to match
3. Add spec-to-document mapping table at top of TLVs.md
4. Regenerate JSON Schema from updated JTD
5. Run tests to verify

## Implementation Summary

### JTD Schema Updates (schemas/docsis-config.jtd.json)
- **TEI**: All 80 `"CM-SP-TEI-I06"` refs → `"CM-SP-TEI-I06 6.7.1.10.1"` (uniform section)
- **eDOCSIS**: All 2 `"CM-SP-eDOCSIS-I31"` refs → `"CM-SP-eDOCSIS-I31 5.2.8.1"` (uniform section)
- **PacketCable PROV1.5**: 1 ref → `"PKT-SP-PROV1.5-C01 9.1"`
- **DSG**: 1 ref → `"CM-SP-eDOCSIS-I31 5.2.8.1, CM-SP-DSG-I25"` (multi-spec)
- **RST-E-DVA**: 1 ref → `"CM-SP-eDOCSIS-I31 5.2.8.1, PKT-SP-RST-E-DVA-C01"` (multi-spec)
- **eRouter**: 34 refs updated with per-property sections (B.4 through B.4.13)
- **L2VPN**: 40 refs updated with per-property sections (B.1.x, B.2.x, B.3.x, B.5.x, B.6, B.7.x)
- **SYNC**: 2 refs updated (D.2.1, D.2.2)
- **DPoE-MULPIv2.0**: 65 refs updated with per-property sections (C.5.x, C.17, C.8.x, C.9.x)
- **DPoE-MULPIv1.0**: 37 refs updated with per-property sections (C.7.1, C.10)

### TLVs.md Updates (docs/TLVs.md)
- Added Spec-to-Document Mapping table after Overview section
- Updated Reference column in eSAFE types table (201-231)
- Updated SYNC TLVs (98-102) references
- Updated all L2VPN sub-TLV references (43.5.x, 45.x, 65.x)
- Updated all MESP sub-TLV references (72.x)
- Updated all L2CP sub-TLV references (83.x)
- Updated all eRouter sub-TLV references (202.x)
- Updated all eTEA sub-TLV references (219.x)

### Generated Schema
- Regenerated schemas/generated/docsis-config.schema.json from updated JTD

## Testing
- All 35 converter tests pass
- TestAllPropertiesHaveSpec passes (all properties have spec with section)
- Zero bare spec references remain (verified with grep)
- JSON validation passes

---
*Source: See UR-010/input.md for full verbatim input*
