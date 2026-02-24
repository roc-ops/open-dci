# OpenDCI Schema Documentation

## Overview

OpenDCI uses a **dual-format schema approach** for DOCSIS configuration files:

1. **JTD Canonical Schema** (`docsis-config.jtd.json`) -- The authoritative schema definition, written in JSON Type Definition (RFC 8927). This is the source of truth for all property names, types, metadata, and structure.

2. **JSON Schema (Generated)** (`generated/docsis-config.schema.json`) -- A JSON Schema 2020-12 file automatically generated from the JTD schema. This is provided for tooling compatibility (IDE autocompletion, linters, CI validation) but is never hand-edited.

The JTD schema is the canonical representation. The JSON Schema is a derived artifact.

## Why JTD (RFC 8927)?

JSON Type Definition was chosen over authoring JSON Schema directly for several reasons:

- **Simplicity** -- JTD has a small, fixed set of forms (type, enum, elements, properties, values, ref, discriminator). There are no ambiguous keywords or draft compatibility issues.
- **Code generation** -- JTD is designed for code generation. The `jtd-codegen` tool can produce Go structs, TypeScript interfaces, and other language bindings directly from the schema.
- **Metadata field** -- JTD's `metadata` keyword provides a clean, standardized place to attach DOCSIS-specific information (TLV types, lengths, spec references) without conflicting with validation logic.
- **RFC standard** -- JTD is RFC 8927, a stable IETF standard. It will not change incompatibly between drafts.
- **Round-trip fidelity** -- The JTD-to-JSON-Schema conversion is mechanical and lossless. DOCSIS metadata maps to `x-docsis-*` extension keywords in JSON Schema.

## Schema Structure

### Root

```json
{
  "metadata": { "title": "...", "version": "...", "docsis": "...", "spec": "..." },
  "definitions": { ... },
  "optionalProperties": { ... }
}
```

- All top-level TLVs are in `optionalProperties` because no single TLV is mandatory in every DOCSIS config file.
- Reusable structures (service flow entries, SNMP MIB entries, vendor extension entries) are in `definitions` and referenced via `ref`.

### JTD Forms Used

| JTD Form | DOCSIS Pattern | Example |
|---|---|---|
| `type` (scalar) | Simple TLVs with fixed types | `DownstreamFrequency` (uint32), `NetworkAccess` (uint8) |
| `enum` | Fixed enumeration values | SNMP MIB `type` field (ASN.1 types) |
| `elements` + `ref` | Repeatable compound TLVs | `UpstreamServiceFlow` (array of entries) |
| `properties` | Required sub-TLVs within a compound | `DocsisExtensionFieldEntry.VendorId` |
| `optionalProperties` | Optional sub-TLVs within a compound | Service flow QoS parameters |
| `ref` (to container) | Vendor-specific extension points | VendorSpecificContainer at TLV x.43 |
| `ref` | Shared definitions | `ServiceFlowErrorEntry` used by both US/DS flows |

## JSONC Comments Support

OpenDCI configuration files use the `.jsonc` extension and support JavaScript-style comments:

```jsonc
{
  // Line comment
  "NetworkAccess": 1,
  /* Block comment */
  "MaxNumCpes": 16
}
```

JSONC is JSON with comments. Any standard JSON is also valid JSONC. Tools that process OpenDCI configs should strip comments before validation.

## Metadata Field Reference

Every property in the JTD schema carries a `metadata` object. The `spec` field is **mandatory** on every property. Other fields are included as appropriate.

| Field | Type | Required | Description |
|---|---|---|---|
| `spec` | string | **Yes** | Specification reference (e.g., `CM-SP-MULPIv4.0 C.1.1.1`) |
| `tlvType` | int or string | No | TLV type number (e.g., `1`) or dotted sub-TLV path (e.g., `24.15`) |
| `tlvLength` | int or string | No | Wire-format length in bytes, or range (e.g., `"2-16"`) |
| `dataType` | string | No | Semantic data type: `uint8`, `uint16`, `uint32`, `string`, `macAddress`, `ipv6Address`, `hexstring`, `oid`, `compound` |
| `description` | string | No | Human-readable description of the field |
| `repeatable` | boolean | No | Whether this TLV can appear multiple times |
| `encoding` | string | No | Special encoding notes (e.g., `bitmask`) |
| `validValues` | object | No | Map of valid numeric values to descriptions |
| `context` | string | No | Contextual notes (e.g., `General Extension (OUI=FFFFFF)`) |

### Metadata in Generated JSON Schema

When the JTD schema is converted to JSON Schema, metadata maps as follows:

```
metadata.description  -->  "description" (standard JSON Schema keyword)
metadata.tlvType      -->  "x-docsis-tlvType"
metadata.tlvLength    -->  "x-docsis-tlvLength"
metadata.dataType     -->  "x-docsis-dataType"
metadata.spec         -->  "x-docsis-spec"
metadata.repeatable   -->  "x-docsis-repeatable"
metadata.encoding     -->  "x-docsis-encoding"
metadata.validValues  -->  "x-docsis-validValues"
metadata.context      -->  "x-docsis-context"
```

The `description` field in JSON Schema is synthesized from `tlvType`, `dataType`, and `spec` when a natural description is not provided.

## TLV 43 -- Vendor Specific Design

TLV 43 (DOCSIS Extension Field) is modeled as a repeatable array of `DocsisExtensionFieldEntry` objects. Each entry has:

- **Required**: `VendorId` -- A 3-byte OUI as a hex string
- **Optional**: CANN-registered General Extension sub-TLVs (`CmLoadBalancingPolicyId`, `CmLoadBalancingPriority`, `CmLoadBalancingGroupId`, `ServiceTypeIdentifier`)
- **Optional**: `VendorSubTlvs` -- An array of generic `{ type, value }` sub-TLVs for vendor-specific data when the OUI is not `FFFFFF`

This design supports both use cases:

```jsonc
// General Extension (OUI = FFFFFF) -- well-known sub-TLVs
{ "VendorId": "FFFFFF", "CmLoadBalancingPolicyId": 2 }

// Vendor Specific (custom OUI) -- generic type/value sub-TLVs
{
  "VendorId": "001018",
  "VendorSubTlvs": [
    { "type": 1, "value": "01" },
    { "type": 2, "value": "0080" }
  ]
}
```

In the binary config file, each entry maps to a separate TLV 43 instance with sub-TLV 8 (Vendor ID) followed by the sub-TLVs.

## Vendor-Specific Extension System

All OUI-gated extension points (TLV x.43 patterns) use a `VendorSpecificContainer` structure containing:

- **`VendorId`** -- 3-byte IEEE OUI identifying the vendor
- **`SubTlvs`** -- Optional array of `VendorSpecificTlvEntry` objects (`{ type, value, dataType? }`)

The `schemas/vendors/` directory holds optional vendor-specific schema files that can map sub-TLV type numbers to named, typed properties. See `schemas/vendors/README.md` for details.

## TLV 11 -- SNMP MIB Design

TLV 11 (SNMP MIB Object) carries a BER-encoded SNMP varbind on the wire. In the JSON representation, each MIB object is decomposed into three required fields:

```jsonc
{
  "oid": "1.3.6.1.2.1.69.1.2.1.6.1",
  "type": "Integer",
  "value": "4"
}
```

The `type` field is an enum of ASN.1 types: `Integer`, `String`, `IPAddress`, `Counter32`, `Gauge32`, `TimeTicks`, `OID`, `HexString`, `Counter64`, `Unsigned32`.

The binary encoder is responsible for BER-encoding the decomposed fields back into a single varbind for TLV 11.

## How to Validate Configs

### Using the Generated JSON Schema

Any JSON Schema 2020-12 compatible validator can be used:

```bash
# Using ajv-cli (Node.js)
npx ajv validate -s schemas/generated/docsis-config.schema.json -d config.json

# Using check-jsonschema (Python)
check-jsonschema --schemafile schemas/generated/docsis-config.schema.json config.json
```

For JSONC files, strip comments first (e.g., with `jsonc-parser` or `jq`).

### Using JTD Validation (Go)

The `json-typedef-go` library provides JTD validation:

```go
import jtd "github.com/jsontypedef/json-typedef-go"

schema, _ := jtd.SchemaFromFile("schemas/docsis-config.jtd.json")
errors, _ := jtd.Validate(schema, configData)
```

### Generating the JSON Schema

Use the included converter tool:

```bash
go run ./tools/jtd2jsonschema schemas/docsis-config.jtd.json > schemas/generated/docsis-config.schema.json
```

## Phase Roadmap

### Phase 1: Representative Subset (Current)

The initial schema covers a representative subset of DOCSIS TLVs that demonstrate all major patterns:

- **Simple scalars**: TLV 1, 2, 3, 9, 14, 18, 29, 58
- **SNMP MIB objects**: TLV 11 (array of varbinds)
- **Upstream/Downstream Service Flows**: TLV 24, 25 (compound with sub-TLVs)
- **Vendor Specific / General Extension**: TLV 43 (compound with `VendorSubTlvs`)

### Phase 2: Core DOCSIS TLVs

Expand to cover all commonly provisioned TLVs:

- Packet Classification (TLV 22, 23, 60)
- Payload Header Suppression (TLV 26)
- Baseline Privacy (TLV 17)
- SNMPv3 Kickstart (TLV 34)
- SNMPv1v2c Coexistence (TLV 53)
- Channel Assignment (TLV 56)
- Subscriber Management (TLV 35-37, 61, 63, 67)

### Phase 3: Advanced and Extension TLVs

- Aggregate Service Flows (TLV 70, 71)
- DOCSIS 3.1/4.0 specific encodings (TLV 84, 85, 88, 96)
- eSAFE TLVs (TLV 201, 202, 216-221)
- L2VPN (TLV 43.5, TLV 65)
- Energy Management (TLV 74, 75, 78, 80)

### Phase 4: Full Coverage

- Complete CANN 11.1 registry coverage
- All sub-TLV hierarchies
- Modem Capabilities (TLV 5) for registration response support
- RPHY/Remote PHY related encodings
