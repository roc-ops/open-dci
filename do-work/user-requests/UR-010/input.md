---
id: UR-010
title: Complete spec section references for all TLVs
created_at: 2026-02-24T08:00:00Z
requests: [REQ-011]
word_count: 215
status: closed
closed_at: 2026-02-24T09:00:00Z
---

# Complete Spec Section References for All TLVs

## Summary

User wants all TLV spec references to include the specific section number in the relevant specification document, not just the spec name. Currently MULPI TLVs have good section references (e.g., "CM-SP-MULPIv4.0 C.1.2.26") but other specs (eDOCSIS, eRouter, L2VPN, TEI, etc.) only list the spec name without a section. This applies to both docs/TLVs.md and the JTD schema's metadata.spec fields. Additionally, wants a spec-to-document mapping table in TLVs.md, leveraging existing info in SPECS.md.

## Full Verbatim Input

in the TLVs document and the JTD spec document, the MULPI stuff has pretty goo Reference Sections but other specs only list the spec not the referenced section. e.g.  "FdxDownstreamUpperBandEdge": {
      "metadata": {
        "description": "FDX Downstream Upper Band Edge frequency",
        "tlvType": 106,
        "tlvLength": 2,
        "dataType": "uint16",
        "spec": "CM-SP-MULPIv4.0 C.1.2.26"
      },
      "type": "uint16"
    },
    "Eps": {
      "metadata": {
        "description": "ePS (eSAFE Power Supply) configuration - opaque eSAFE encoding",
        "tlvType": 201,
        "dataType": "hexstring",
        "spec": "CM-SP-eDOCSIS-I31"
      },
      "type": "string"
    }, other items only have [SPEC NAME] for the referenced section this should be consistent and have the referenced section in the relevant spec for all. If a TLV spans multiple Specs or is a common TLV in multiple sections, all Referenced SPECS and Sections should be listed. The idea is if there is a question on a particular TLV or sub-TLV we should reference the section in the appropriate spec document. In the TLVs document we should have a table where SPEC e.g. MULPI is mapped to one or more Documents that the references come from. This is available in SPECS.md currently.

---
*Captured: 2026-02-24T08:00:00Z*
