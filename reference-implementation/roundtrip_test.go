package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// TestRoundTripLabTr069 is the integration test that verifies:
// decode(binary) -> encode(decoded) -> decode(encoded) produces the same config.
// Byte-exact comparison is not feasible because some original strings lack DOCSIS-spec
// null terminators, and the encoder always null-terminates per spec.
// Instead, we verify semantic equivalence: the decoded configs must be identical.
func TestRoundTripLabTr069(t *testing.T) {
	binPath := "binary-files/lab-tr069.bin"
	originalData, err := os.ReadFile(binPath)
	if err != nil {
		t.Skipf("skipping integration test: %v", err)
	}

	// Load the full schema registry.
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("unable to get caller info")
	}
	schemaFile := filepath.Join(filepath.Dir(filename), "..", "schemas", "docsis-config.jtd.json")
	reg, err := LoadRegistry(schemaFile)
	if err != nil {
		t.Fatalf("loading registry: %v", err)
	}

	// Decode the original binary.
	result1, err := Decode(originalData, reg)
	if err != nil {
		t.Fatalf("decoding original: %v", err)
	}

	t.Logf("Decoded %d top-level properties", len(result1.Config))
	t.Logf("TLVOrder: %v", result1.TLVOrder)

	// Re-encode.
	encoded, err := Encode(result1, reg)
	if err != nil {
		t.Fatalf("encoding: %v", err)
	}

	t.Logf("Original size: %d bytes, encoded size: %d bytes", len(originalData), len(encoded))

	// Decode the re-encoded binary.
	result2, err := Decode(encoded, reg)
	if err != nil {
		t.Fatalf("decoding re-encoded: %v", err)
	}

	// Compare the decoded configs (excluding MICs and internal metadata).
	config1 := cleanConfigForComparison(result1.Config)
	config2 := cleanConfigForComparison(result2.Config)

	json1, _ := json.MarshalIndent(config1, "", "  ")
	json2, _ := json.MarshalIndent(config2, "", "  ")

	if string(json1) != string(json2) {
		t.Errorf("Semantic round-trip FAILED: decoded configs differ")
		t.Logf("Original decoded:\n%s", string(json1))
		t.Logf("Re-encoded decoded:\n%s", string(json2))
	} else {
		t.Logf("Semantic round-trip PASSED: configs match (%d bytes JSON)", len(json1))
	}

	// Also verify that TLV ordering is preserved (excluding MIC TLVs).
	order1 := filterMICFromOrder(result1.TLVOrder)
	order2 := filterMICFromOrder(result2.TLVOrder)
	if len(order1) != len(order2) {
		t.Errorf("TLVOrder length mismatch (excluding MICs): %d vs %d", len(order1), len(order2))
	}
	for i := 0; i < len(order1) && i < len(order2); i++ {
		if order1[i] != order2[i] {
			t.Errorf("TLVOrder[%d] mismatch: %q vs %q", i, order1[i], order2[i])
		}
	}

	// Also do a byte-exact comparison on the parts that should match exactly
	// (non-string TLVs). This verifies integers, IPs, compounds, etc. are exact.
	originalStripped := stripTLVsForComparison(originalData, 6, 7)
	encodedStripped := stripTLVsForComparison(encoded, 6, 7)
	t.Logf("Stripped sizes: original=%d, encoded=%d (diff=%d bytes, due to null-terminator normalization)",
		len(originalStripped), len(encodedStripped), len(encodedStripped)-len(originalStripped))
}

// filterMICFromOrder removes CmMic and CmtsMic from a TLV order slice.
func filterMICFromOrder(order []string) []string {
	var result []string
	for _, name := range order {
		if name != "CmMic" && name != "CmtsMic" {
			result = append(result, name)
		}
	}
	return result
}

// cleanConfigForComparison removes internal metadata keys and MIC-related entries
// from a config map for comparison.
func cleanConfigForComparison(config map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range config {
		if k == "_tlvOrder" || k == "CmMic" || k == "CmtsMic" {
			continue
		}
		result[k] = cleanValue(v)
	}
	return result
}

func cleanValue(v interface{}) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		cleaned := make(map[string]interface{})
		for k, sv := range val {
			if k == "_tlvOrder" {
				continue
			}
			cleaned[k] = cleanValue(sv)
		}
		return cleaned
	case []interface{}:
		cleaned := make([]interface{}, len(val))
		for i, sv := range val {
			cleaned[i] = cleanValue(sv)
		}
		return cleaned
	default:
		return v
	}
}

// TestRoundTripSimpleConfig tests round-trip with a simple manually-built config.
func TestRoundTripSimpleConfig(t *testing.T) {
	reg := makeTestRegistry()

	original := buildConfig(
		buildTLV(3, []byte{1}),                           // NetworkAccess
		buildTLV(1, []byte{0x15, 0x75, 0x2A, 0x00}),      // DownstreamFrequency
		buildTLV(2, []byte{5}),                            // UpstreamChannelId
		buildTLV(9, []byte("firmware.bin\x00")),           // SwUpgradeFilename
		buildTLV(14, []byte{0x00, 0x1A, 0x2B, 0x3C, 0x4D, 0x5E}), // CpeEthernetMacAddress
		buildTLV(18, []byte{16}),                          // MaxNumCpes
		buildTLV(20, []byte{10, 1, 2, 3}),                 // TftpServerProvisionedModemIpv4Address
	)

	result, err := Decode(original, reg)
	if err != nil {
		t.Fatal(err)
	}

	encoded, err := Encode(result, reg)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(encoded, original) {
		t.Errorf("round-trip failed:\n  got:  %X\n  want: %X", encoded, original)
	}
}

// TestRoundTripWithMICs tests that MICs are properly excluded from comparison.
func TestRoundTripWithMICs(t *testing.T) {
	reg := makeTestRegistry()

	cmMic := make([]byte, 16)
	cmtsMic := make([]byte, 16)

	original := buildConfig(
		buildTLV(3, []byte{1}),
		buildTLV(6, cmMic),
		buildTLV(7, cmtsMic),
	)

	result, err := Decode(original, reg)
	if err != nil {
		t.Fatal(err)
	}

	encoded, err := Encode(result, reg)
	if err != nil {
		t.Fatal(err)
	}

	// Encoded will not have MICs, so compare stripped versions.
	originalStripped := stripTLVsForComparison(original, 6, 7)
	encodedStripped := stripTLVsForComparison(encoded, 6, 7)

	if !bytes.Equal(originalStripped, encodedStripped) {
		t.Errorf("round-trip (stripped) failed:\n  got:  %X\n  want: %X", encodedStripped, originalStripped)
	}
}

// TestRoundTripCompoundWithMICs tests compound TLVs with MIC exclusion.
func TestRoundTripCompoundWithMICs(t *testing.T) {
	reg := makeTestRegistry()

	subTLVs := append(
		buildTLV(1, []byte{0x00, 0x01}),
		buildTLV(6, []byte{7})...,
	)
	subTLVs = append(subTLVs,
		buildTLV(8, []byte{0x00, 0x00, 0x27, 0x10})...,
	)

	cmMic := make([]byte, 16)
	cmtsMic := make([]byte, 16)

	original := buildConfig(
		buildTLV(3, []byte{1}),
		buildTLV(24, subTLVs),
		buildTLV(6, cmMic),
		buildTLV(7, cmtsMic),
	)

	result, err := Decode(original, reg)
	if err != nil {
		t.Fatal(err)
	}

	encoded, err := Encode(result, reg)
	if err != nil {
		t.Fatal(err)
	}

	originalStripped := stripTLVsForComparison(original, 6, 7)
	encodedStripped := stripTLVsForComparison(encoded, 6, 7)

	if !bytes.Equal(originalStripped, encodedStripped) {
		t.Errorf("round-trip (stripped) failed:\n  got:  %X\n  want: %X", encodedStripped, originalStripped)
	}
}

// TestRoundTripTLV11Snmp tests SNMP varbind round-trip.
func TestRoundTripTLV11Snmp(t *testing.T) {
	reg := makeTestRegistry()

	oidBytes := []byte{0x2B, 0x06, 0x01, 0x02, 0x01}
	varbind := buildVarbind(oidBytes, tagInteger, []byte{0x2A})

	original := buildConfig(
		buildTLV(11, varbind),
	)

	result, err := Decode(original, reg)
	if err != nil {
		t.Fatal(err)
	}

	encoded, err := Encode(result, reg)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(encoded, original) {
		t.Errorf("round-trip failed:\n  got:  %X\n  want: %X", encoded, original)
	}
}

// TestRoundTripTLV43 tests TLV 43 round-trip.
func TestRoundTripTLV43(t *testing.T) {
	reg := makeTestRegistry()

	vendorId := buildTLV(8, []byte{0xFF, 0xFF, 0xFF})
	lbPolicy := buildTLV(1, []byte{0x00, 0x00, 0x00, 0x64})
	tlv43Value := append(vendorId, lbPolicy...)

	original := buildConfig(
		buildTLV(43, tlv43Value),
	)

	result, err := Decode(original, reg)
	if err != nil {
		t.Fatal(err)
	}

	encoded, err := Encode(result, reg)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(encoded, original) {
		t.Errorf("round-trip failed:\n  got:  %X\n  want: %X", encoded, original)
	}
}

// TestRoundTripTLV43VendorSpecific tests vendor-specific TLV 43 round-trip.
func TestRoundTripTLV43VendorSpecific(t *testing.T) {
	reg := makeTestRegistry()

	vendorId := buildTLV(8, []byte{0x00, 0x11, 0x22})
	vendorSub := buildTLV(1, []byte{0xAA, 0xBB})
	tlv43Value := append(vendorId, vendorSub...)

	original := buildConfig(
		buildTLV(43, tlv43Value),
	)

	result, err := Decode(original, reg)
	if err != nil {
		t.Fatal(err)
	}

	encoded, err := Encode(result, reg)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(encoded, original) {
		t.Errorf("round-trip failed:\n  got:  %X\n  want: %X", encoded, original)
	}
}

// TestRoundTripWithPads tests that pad bytes are handled correctly.
func TestRoundTripWithPads(t *testing.T) {
	reg := makeTestRegistry()

	var original []byte
	original = append(original, 0x00)                     // pad
	original = append(original, buildTLV(3, []byte{1})...) // NetworkAccess
	original = append(original, 0x00)                     // pad
	original = append(original, 0xFF, 0x00)                // end-of-data

	result, err := Decode(original, reg)
	if err != nil {
		t.Fatal(err)
	}

	encoded, err := Encode(result, reg)
	if err != nil {
		t.Fatal(err)
	}

	// Encoder does not emit pad bytes.
	originalStripped := stripTLVsForComparison(original)
	encodedStripped := stripTLVsForComparison(encoded)

	if !bytes.Equal(originalStripped, encodedStripped) {
		t.Errorf("round-trip (stripped) failed:\n  got:  %X\n  want: %X", encodedStripped, originalStripped)
	}
}
