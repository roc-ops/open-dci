# OpenDCI

**Open DOCSIS Configuration Interchange** -- A JSON-based data interchange format for DOCSIS cable modem configuration files.

## What is this?

DOCSIS cable modems are configured using binary files built from Type-Length-Value (TLV) encodings defined across multiple CableLabs specifications. These binary files are opaque, vendor-tool-dependent, and difficult to version, diff, or automate.

OpenDCI defines a JSON representation for DOCSIS configuration data. Instead of working with binary TLV blobs, you work with structured, human-readable JSON files that can be validated against a schema, version-controlled with git, and processed with standard tooling.

```jsonc
{
  "NetworkAccess": 1,
  "MaxNumCpes": 16,
  "PrivacyEnable": 1,
  "DownstreamFrequency": 855000000,

  "UpstreamServiceFlow": [
    {
      "ServiceFlowReference": 1,
      "QosParamSetType": 7,
      "MaxSustainedTrafficRate": 5000000,
      "SchedulingType": 2  // best effort
    }
  ],

  "DownstreamServiceFlow": [
    {
      "ServiceFlowReference": 2,
      "QosParamSetType": 7,
      "MaxSustainedTrafficRate": 50000000
    }
  ]
}
```

OpenDCI config files use the `.jsonc` extension and support JavaScript-style comments (`//` and `/* */`).

## Schema

OpenDCI uses a **dual-format schema** approach:

- **JTD (canonical)** -- `schemas/docsis-config.jtd.json` is the authoritative schema, written in [JSON Type Definition (RFC 8927)](https://datatracker.ietf.org/doc/rfc8927/). This is the source of truth for property names, types, and DOCSIS metadata (TLV type numbers, wire lengths, spec references).

- **JSON Schema (generated)** -- `schemas/generated/docsis-config.schema.json` is a JSON Schema 2020-12 file automatically generated from the JTD source. Use this for IDE autocompletion, linting, and CI validation. It is never hand-edited.

JTD was chosen for its simplicity, code-generation support, and stable RFC standard. See [`schemas/README.md`](schemas/README.md) for the full schema documentation including metadata field reference, design decisions for complex TLVs (vendor-specific, SNMP MIBs), and validation examples.

### Validating configs

```bash
# Using ajv-cli (Node.js)
npx ajv validate -s schemas/generated/docsis-config.schema.json -d config.jsonc

# Using check-jsonschema (Python)
check-jsonschema --schemafile schemas/generated/docsis-config.schema.json config.jsonc
```

### Regenerating the JSON Schema

```bash
go run ./tools/jtd2jsonschema schemas/docsis-config.jtd.json > schemas/generated/docsis-config.schema.json
```

## Repository Layout

```
schemas/
  docsis-config.jtd.json              # Canonical JTD schema (source of truth)
  generated/
    docsis-config.schema.json         # Generated JSON Schema 2020-12
  vendors/                            # Optional vendor-specific extension schemas
    001018.jtd.json                   # Example: Broadcom (OUI 00:10:18)
  examples/
    basic-config.jsonc                # Basic CM provisioning example
    vendor-specific-config.jsonc      # Vendor extension example
docs/
  TLVs.md                            # Complete TLV reference (all CANN-registered TLVs)
  SPECS.md                           # Specification download guide
tools/
  jtd2jsonschema/                    # JTD-to-JSON-Schema converter (Go)
```

## Specifications

OpenDCI covers TLV encodings from across the DOCSIS specification family -- MULPI (core TLVs), L2VPN, eRouter, eDOCSIS, PacketCable, DSG, TEI, SYNC, DPoE, and DPoG.

See [`docs/SPECS.md`](docs/SPECS.md) for the full list of specifications, versions, and download instructions. The complete TLV catalog with spec section references is in [`docs/TLVs.md`](docs/TLVs.md).

## License

[MIT](LICENSE)
