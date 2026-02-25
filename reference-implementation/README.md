# OpenDCI Reference Implementation

A Go CLI tool that decodes and encodes binary DOCSIS configuration files using the OpenDCI JSONC format.

## Build

```bash
cd reference-implementation
go build -o opendci-decode .
```

## Usage

### Decode (binary to JSONC)

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

### Encode (JSONC/JSON to binary)

```bash
# Encode a JSONC file to binary
./opendci-decode --encode -i config.jsonc -o config.bin

# Encode with MIC computation
./opendci-decode --encode -i config.jsonc -o config.bin --cmts-secret "mysecret"

# Pipe from stdin
cat config.jsonc | ./opendci-decode --encode -o config.bin
```

### Round-trip

```bash
# Decode, then re-encode
./opendci-decode -i original.bin -o config.jsonc
./opendci-decode --encode -i config.jsonc -o rebuilt.bin
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-i`, `--input` | stdin | Input file (binary for decode, JSONC/JSON for encode) |
| `-o`, `--output` | stdout | Output file (JSONC for decode, binary for encode) |
| `--encode` | `false` | Encode mode: convert JSON/JSONC to binary TLV format |
| `--cmts-secret` | *(none)* | CMTS shared secret (decode: MIC verification, encode: MIC computation) |
| `--schema` | `../schemas/docsis-config.jtd.json` | Path to JTD schema file |
| `--mibs-dir` | `../mibs/` | Path to mibs/ root directory (decode only) |
| `--with-mibs` | *(none)* | Comma-separated version-specific MIBs, e.g. `DOCS-IF3-MIB@2024-07-05` (decode only) |
| `--no-mibs` | `false` | Disable MIB resolution (decode only) |

## JSONC Output

The decode output is JSONC (JSON with Comments). Non-schema items are emitted as `//` comments so the output is round-trip safe -- if fed to the encoder, only schema-conforming properties are processed.

Commented-out items:
- **MIC values** (`CmMic`, `CmtsMic`) -- computed by the encoder, not provisioned
- **MIC validation results** -- diagnostic messages (also printed to stderr)
- **Unknown TLVs** -- preserved for the operator but invisible to the encoder

## Encoding

The encoder reads JSON or JSONC input and produces a binary DOCSIS config file in TLV format. JSONC comments are automatically stripped before parsing.

- **MIC computation**: When `--cmts-secret` is provided, the encoder computes and inserts both CM MIC (TLV 6) and CMTS MIC (TLV 7). Without it, MIC TLVs are omitted.
- **TLV ordering**: The encoder uses schema-defined TLV type numbers to determine output ordering.
- **End-of-data marker**: TLV 255 is always appended at the end.

## MIC Verification

MIC verification results are printed to stderr and included as JSONC comments in the output.

- **CM MIC** (TLV 6) is always verified if present (uses empty HMAC-MD5 key)
- **CMTS MIC** (TLV 7) is only verified if `--cmts-secret` is provided

## How It Works

Both the decoder and encoder are **schema-driven** -- they load TLV definitions from the OpenDCI JTD schema (`docsis-config.jtd.json`) at runtime. This means they automatically support new TLVs when the schema is updated, without code changes.

### Decode

1. **Registry**: Parses the JTD schema to build a hierarchical TLV type-to-property mapping
2. **Decode**: Reads the binary TLV stream, looks up each TLV in the registry, converts wire bytes to JSON values
3. **Compound TLVs**: Recursively decoded using sub-TLV definitions from the schema
4. **TLV 11 (SNMP)**: BER-encoded varbinds are parsed into `{oid, type, value}` objects
5. **TLV 43 (Vendor Specific)**: Routes to general extension sub-TLVs (OUI=FFFFFF) or collects as generic hex pairs
6. **Unknown TLVs**: Passed through as `{type, value}` hex pairs in JSONC comments -- nothing is silently dropped

### Encode

1. **Registry**: Uses reverse lookup (property name to TLV type code) from the same JTD schema
2. **Encode**: Iterates JSON properties, converts values back to wire bytes using the registry's data type definitions
3. **Compound TLVs**: Recursively encoded with sub-TLV nesting
4. **TLV 11 (SNMP)**: `{oid, type, value}` objects are re-encoded as BER varbinds
5. **TLV 43 (Vendor Specific)**: VendorId emitted first, then sub-TLVs (general extension) or raw hex pairs (vendor-specific)
6. **MICs**: Optionally computed and inserted when `--cmts-secret` is provided

## Tests

```bash
go test ./... -v
```
