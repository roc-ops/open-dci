# OpenDCI Reference Implementation

A Go CLI tool that decodes binary DOCSIS configuration files into OpenDCI JSONC format.

## Build

```bash
cd reference-implementation
go build -o opendci-decode .
```

## Usage

```bash
# Decode a binary config file to stdout
./opendci-decode -i config.bin

# Pipe from stdin, write to file
cat config.bin | ./opendci-decode -o config.jsonc

# Verify CMTS MIC with shared secret
./opendci-decode -i config.bin --cmts-secret "mysecret"

# Specify schema path (auto-detected by default)
./opendci-decode -i config.bin --schema ../schemas/docsis-config.jtd.json
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-i`, `--input` | stdin | Input binary config file |
| `-o`, `--output` | stdout | Output JSONC file (`.jsonc` extension auto-added if none) |
| `--cmts-secret` | *(none)* | CMTS shared secret for MIC verification |
| `--schema` | `../schemas/docsis-config.jtd.json` | Path to JTD schema file |

## JSONC Output

The output is JSONC (JSON with Comments). Non-schema items are emitted as `//` comments so the output is round-trip safe -- if fed to an encoder, only schema-conforming properties are processed.

Commented-out items:
- **MIC values** (`CmMic`, `CmtsMic`) -- computed by the encoder, not provisioned
- **MIC validation results** -- diagnostic messages (also printed to stderr)
- **Unknown TLVs** -- preserved for the operator but invisible to an encoder

## MIC Verification

MIC verification results are printed to stderr and included as JSONC comments in the output.

- **CM MIC** (TLV 6) is always verified if present (uses empty HMAC-MD5 key)
- **CMTS MIC** (TLV 7) is only verified if `--cmts-secret` is provided

## How It Works

The decoder is **schema-driven** -- it loads TLV definitions from the OpenDCI JTD schema (`docsis-config.jtd.json`) at runtime. This means it automatically supports new TLVs when the schema is updated, without code changes.

1. **Registry**: Parses the JTD schema to build a hierarchical TLV type-to-property mapping
2. **Decode**: Reads the binary TLV stream, looks up each TLV in the registry, converts wire bytes to JSON values
3. **Compound TLVs**: Recursively decoded using sub-TLV definitions from the schema
4. **TLV 11 (SNMP)**: BER-encoded varbinds are parsed into `{oid, type, value}` objects
5. **TLV 43 (Vendor Specific)**: Routes to general extension sub-TLVs (OUI=FFFFFF) or collects as generic hex pairs
6. **Unknown TLVs**: Passed through as `{type, value}` hex pairs in JSONC comments -- nothing is silently dropped

## Tests

```bash
go test ./... -v
```
