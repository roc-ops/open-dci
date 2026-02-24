---
id: REQ-010
title: Vendor-specific TLV extension system
status: completed
created_at: 2026-02-23T10:00:00Z
claimed_at: 2026-02-23T10:30:00Z
route: C
completed_at: 2026-02-23T11:30:00Z
user_request: UR-009
---

# Vendor-Specific TLV Extension System

## What
Design and implement a vendor-specific TLV extension system that allows generic representation of unknown vendor TLVs (hex-encoded) and an extension mechanism for vendors to define their own schemas without modifying the main schema.

## Detailed Requirements

### Generic Vendor-Specific Representation
- When no vendor schema exists, vendor-specific sub-TLVs are represented as `{ "type": <number>, "value": "<hex-encoded-bytes>" }`
- Hex encoding is required because binary decode cannot determine data types from raw bytes alone
- Optional `dataType` hint for encoding convenience (e.g., `{ "type": 5, "value": 12345, "dataType": "uint32" }`), but decode always produces hex

### Vendor Schema Directory
- Separate JTD schema files in `schemas/vendors/` keyed by OUI (e.g., `001018.jtd.json` for Broadcom)
- Each file is standalone, defining sub-TLVs for that vendor's OUI
- Main schema is never modified by vendors

### OUI-Gated Extension Points Only
- Vendor schemas only apply where DOCSIS uses OUI-based namespacing:
  - TLV 43 (DOCSIS Extension Field) — starts with sub-TLV 8 (Vendor ID/OUI)
  - TLV 43.5.43 (Vendor Specific L2VPN Subtype)
  - x.43 (Vendor Specific QoS/Classifier Parameters in service flows and classifiers)
  - 202.43 (eRouter Vendor Specific)
  - 219.43 (eTEA Vendor Specific)
- eSAFE containers (201, 216, 217, 220, 221) are NOT vendor-specific — they're spec-defined opaque blobs and don't use this mechanism

### Vendor Schema File Format
- Each vendor file declares OUI, vendor name, and which extension points it covers
- Uses `extensionPoints` metadata field (e.g., `["43", "x.43", "202.43"]`)
- Standard JTD format with `optionalProperties` defining the vendor's sub-TLVs
- Each property has `metadata.tlvType` and `metadata.spec`

### Runtime Behavior
- Tools load main schema + scan `schemas/vendors/` directory
- If OUI matches a vendor schema, use typed encoding/decoding
- Otherwise fall back to generic hex representation
- Graceful degradation — unknown OUI just means hex encoding

### Deliverables
- Update main JTD schema to support the generic `{ type, value }` pattern at all OUI-gated extension points
- Create `schemas/vendors/` directory with a template/example vendor schema
- Update docs/TLVs.md to document the vendor extension system
- Update JSON Schema generation to handle vendor schemas (or document the merge approach)
- Add/update examples showing vendor-specific TLV usage (both generic hex and with vendor schema)
- Update any relevant spec/design documentation

## Design Rationale
- When decoding a binary config file, vendor-specific sub-TLVs are just raw bytes — the decoder has the type number and length but cannot determine whether a 4-byte value is a uint32, an IPv4 address, or opaque data. This is why hex encoding is the mandatory default representation, not a preference.
- Vendor-specific TLVs are distinct from "encapsulated" compound TLVs — encapsulated TLVs have known sub-TLV structure defined in the spec, while vendor-specific TLVs are opaque without the vendor's schema.

## Builder Guidance
- Certainty level: Firm on the OUI-gating rule and hex-as-default. The exact JSON representation format and vendor schema file structure have latitude for refinement during implementation.
- The optional `dataType` hint for encoding is a firm requirement — users who know their vendor's spec should be able to provide typed values for encoding convenience, while decode always produces hex.
- The key principle: differentiation between vendor implementations requires OUI namespacing — no OUI, no vendor schema
- eSAFE exclusion is explicit: user confirmed these are spec-defined, not vendor-specific

## Full Context
See [user-requests/UR-009/input.md](./user-requests/UR-009/input.md) for complete verbatim input.

---
*Source: See UR-009/input.md for full verbatim input*

---

## Triage

**Route: C** - Complex

**Reasoning:** New architectural feature spanning schema definitions, 9 extension points, vendor directory, examples, tests, docs, and JSON Schema generation. Multiple design decisions needed.

**Planning:** Required

## Plan

### Implementation Steps

1. **Add definitions**: Create `VendorSpecificTlvEntry` (type + hex value + optional dataType hint) and `VendorSpecificContainer` (VendorId + SubTlvs array) in the JTD schema
2. **Restructure TLV 43**: Remove `additionalProperties: true` from `DocsisExtensionFieldEntry`, add `VendorSubTlvs` optional property
3. **Update 9 extension points**: Change all vendor-specific properties from flat `"type": "string"` to `"ref": "VendorSpecificContainer"` (x.43 in classifiers/service flows, 26.43, 43.5.43, 202.43, 219.43)
4. **Create vendor directory**: `schemas/vendors/` with README.md, template, and example Broadcom schema
5. **Update examples**: Rewrite vendor-specific-config.jsonc to use new VendorSubTlvs pattern
6. **Update tests**: Fix additionalProperties assertion, add vendor-specific structure tests
7. **Regenerate JSON Schema**
8. **Update documentation**: TLVs.md vendor extension section, schemas/README.md

### Key Design Decisions
- `value` field is always a hex string (avoids JTD oneOf limitation)
- `dataType` is an enum hint for encoding convenience
- TLV 43 keeps General Extension named properties alongside new `VendorSubTlvs`
- Vendor schema merging into JSON Schema is deferred (optional post-generation step)

*Generated by Plan agent*

## Implementation Summary

- Added `VendorSpecificTlvEntry` definition (type + hex value + optional dataType enum hint)
- Added `VendorSpecificContainer` definition (VendorId + SubTlvs array)
- Restructured `DocsisExtensionFieldEntry`: removed `additionalProperties: true`, added `VendorSubTlvs`
- Changed 9 vendor-specific extension points from flat hexstring to `VendorSpecificContainer` ref
- Created `schemas/vendors/` directory with README.md and example `001018.jtd.json` (Broadcom)
- Updated `schemas/examples/vendor-specific-config.jsonc` with new VendorSubTlvs pattern
- Updated converter tests for new schema structure
- Regenerated JSON Schema
- Added Section 5 to docs/TLVs.md documenting the vendor extension system
- Updated schemas/README.md with vendor extension architecture

*Completed by work action (Route C)*

## Testing

**Tests run:** `go test ./tools/jtd2jsonschema/ -v`
**Result:** All tests passing

**Tests updated:**
- `TestConvertDocsisSchema` — updated additionalProperties assertion, added VendorSpecificTlvEntry/VendorSpecificContainer checks

**Existing tests verified:**
- `TestAllPropertiesHaveSpec` — still passing with new definitions
- All 24+ tests passing

*Verified by work action*
