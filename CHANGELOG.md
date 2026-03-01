# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.1] - 2026-02-28

### Added

#### PacketCable MTA Config Support

- MTA JTD schema (`schemas/mta-config.jtd.json`) covering PacketCable 1.5 TLV catalog (TLVs 11, 38, 64, 254)
- MTA decode/encode mode with TLV 254 start/end delimiters and format auto-detection
- TLV 216 (eMTA) recursive decode/encode -- embedded MTA configs in CM files are automatically expanded to structured JSON instead of opaque hexstrings
- Generalized SNMP varbind handling for both CM TLV 11 and MTA TLV 64 (2-byte length)

#### CLI

- `--format` flag to explicitly select `cm` or `mta` mode (default: auto-detect from content)
- `--mta-schema` flag to specify MTA schema path (default: sibling of CM schema)
- Auto-detection of MTA binary format (first byte `0xFE`)
- Automatic loading of MTA registry as nested registry for TLV 216 recursive decode

#### WebAssembly Module

- `opendciLoadMtaSchema` function to load MTA registry alongside CM registry (13 exported functions total)
- Auto-detection of MTA format in `opendciDecode` (selects MTA registry when first byte is `0xFE`)
- Optional `format` argument in `opendciEncode` for explicit MTA encoding
- Cross-wiring of nested registries regardless of CM/MTA schema load order

## [0.1.0] - 2026-02-28

Initial public release of OpenDCI -- a JSON-based data interchange format for DOCSIS cable modem configuration files.

### Added

#### Schema

- Canonical JTD schema (RFC 8927) covering the full DOCSIS TLV catalog
- Auto-generated JSON Schema 2020-12 for IDE validation and CI
- Vendor-specific extension system with OUI-gated schema loading
- Enum metadata (`validValues`) for 67+ TLVs
- Compound TLV expansion for nested sub-TLV structures
- Spec section references for all TLV definitions
- DPA vendor schemas for Tibit and Readylinks

#### Reference Implementation (Go CLI)

- Decode binary `.cm`/`.bin` files to JSONC
- Encode JSONC back to binary with round-trip fidelity
- CM/CMTS MIC verification and computation (MULPI Annex D canonical ordering)
- SNMP TLV 10/11 BER encode/decode
- MIB resolution (OID to human-readable names)
- PacketCable MTA config hash computation (NA/EU/IETF variants)
- CVC TLV 254-byte chunking for oversize certificates
- Two-byte TLV length support
- 4-byte alignment padding
- Vendor-specific TLV decoding via extension schemas
- JSONC output with inline comments (enum labels, MIC status, unknown TLVs)
- Version-stamped output header

#### WebAssembly Module

- Full decoder/encoder compiled to WASM for browser use
- 12 exported JS functions: `opendciLoadSchema`, `opendciLoadVendorSchema`, `opendciDecode`, `opendciEncode`, `opendciInitMIBs`, `opendciLoadMIBs`, `opendciQueryMIBTree`, `opendciResolveName`, `opendciResolveOID`, `opendciExtractCVC`, `opendciSerializeMIBState`, `opendciRestoreMIBState`

#### Library

- Extracted importable `opendci` Go package for programmatic use

#### Documentation

- Complete TLV reference (`docs/TLVs.md`)
- Specification download guide (`docs/SPECS.md`)
- Example configs (basic provisioning, vendor-specific extensions)
- Schema documentation with metadata field reference and design notes

#### Tooling

- JTD-to-JSON-Schema converter (`tools/jtd2jsonschema`)

[0.1.1]: https://github.com/roc-ops/open-dci/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/roc-ops/open-dci/releases/tag/v0.1.0
