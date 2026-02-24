---
id: REQ-009
title: Make TLVs.md tables consistent with data type column
status: completed
created_at: 2026-02-23T08:00:00Z
claimed_at: 2026-02-23T09:00:00Z
route: B
completed_at: 2026-02-23T09:30:00Z
user_request: UR-008
---

# Make TLVs.md Tables Consistent with Data Type Column

## What
Standardize all tables in docs/TLVs.md to use a consistent column structure, matching the top-level table format (which has Spec and Reference columns). Add a Data Type column to all tables using "encapsulated" for compound/container TLVs.

## Detailed Requirements
- All sub-TLV tables (Sections 3.1-3.13 and Section 4) should have the same columns as the top-level table
- Top-level table (Section 2) has: Type, CANN Name, Length, Cfg File, Spec, Reference — this is the gold standard
- Sub-TLV tables currently vary: some have Type/CANN Name/Spec, others have Type/CANN Name/Reference, some have Applies To
- Add a **Data Type** column to all tables (e.g., uint8, uint16, uint32, string, hexstring, macAddress, ipv4Address, ipv6Address, encapsulated for compound TLVs)
- Use "encapsulated" for TLVs that contain sub-TLVs (compound/container types)
- Maintain consistency across every table in the document

## Builder Guidance
- Certainty level: Firm — user explicitly wants consistency across all tables
- The user said "encapsulated" for compound types — use that term
- The top-level table is the reference pattern to follow
- Sub-TLV tables that have an "Applies To" column should keep it (it's useful for shared sub-TLVs like classification and service flow)

## Full Context
See [user-requests/UR-008/input.md](./user-requests/UR-008/input.md) for complete verbatim input.

---
*Source: See UR-008/input.md for full verbatim input*

---

## Triage

**Route: B** - Medium

**Reasoning:** Clear feature with well-defined outcome (consistent tables + Data Type column). Need to explore existing table patterns and determine data types from the JTD schema.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

Rationale: Clear documentation-only change. The target column structure is well-defined, just need to map data types from the JTD schema to each TLV row.

*Skipped by work action*

## Implementation Summary

- Added **Data Type** column to all 34 tables in docs/TLVs.md
- Added **Reference** column to all sub-TLV tables that were missing it (Sections 3.1, 3.3-3.13)
- Added **Spec** column to tables that only had Reference (Sections 3.2, 4)
- Three consistent column structures across all tables:
  - Top-level: `Type | CANN Name | Data Type | Length | Cfg File | Spec | Reference`
  - Sub-TLV: `Type | CANN Name | Data Type | Spec | Reference`
  - Sub-TLV with shared parent: `Type | CANN Name | Data Type | Applies To | Spec | Reference`
- Data types determined from JTD schema metadata and DOCSIS spec conventions:
  - uint8/uint16/uint32 for scalar integers
  - string, hexstring, macAddress, ipv4Address, ipv6Address for typed values
  - "encapsulated" for compound/container TLVs with sub-TLVs
  - "-" for deprecated/reserved/unassigned entries

*Completed by work action (Route B)*

## Testing

**Tests run:** `go test ./tools/jtd2jsonschema/ -v`
**Result:** All 24 tests passing (cached)

**No new tests needed** — documentation-only change, no code modified.

*Verified by work action*
