---
id: UR-009
title: Vendor-specific TLV extension system
created_at: 2026-02-23T10:00:00Z
requests: [REQ-010]
word_count: 450
---

# Vendor-Specific TLV Extension System

## Summary

User wants a system for handling vendor-specific TLVs (OUI-namespaced sub-TLVs) in OpenDCI. This includes a generic hex-encoded representation for unknown vendor TLVs, plus an extension mechanism allowing vendors to define their own schemas in separate files without modifying the main schema. Also needs spec/docs and examples updated.

## Full Verbatim Input

[Discussion between user and assistant leading to design decisions]

User: without making changes I want to work somethings with you: when we have Vendor specific encodings besides 0xffffff this spec will not have the names for the TLVs, we need a way to generically describe an unknown TLV, something like type e.g. 5, value-type e.g. uint32 and value, but I am not sure that will work, I am thinking when you decode a binary file, it may not be possible to decode the value-type, so maybe a hex encoded? Thoughts?

[Assistant proposed hex-encoded values as the default since binary decode can't determine data types from raw bytes. User agreed.]

User: that sounds good, so far, but the other part of this is it would be useful for other vendors to define the schema for their Vendor Specifc sub TLVs is there a clean way to do that without modifying the main schema, like ability to impport schemas?

[Assistant proposed a vendor schema directory at schemas/vendors/ keyed by OUI, with JTD schema files that tools load at runtime. Main schema stays clean, vendor schemas are additive.]

User: broader extension points are fine as long as they comply with vendor OUI like Vendor Extension, otherwirse how would we differenctiate between vendor implementations/

[Assistant confirmed OUI is the natural key. Vendor schemas only apply at OUI-gated extension points: TLV 43, 43.5.43, x.43, 202.43, 219.43. eSAFE containers (201, 216, 217, 220, 221) are NOT vendor-specific — they're spec-defined.]

User: lets capture this as a new request, also noting we need to update the spec and examples

---
*Captured: 2026-02-23T10:00:00Z*
