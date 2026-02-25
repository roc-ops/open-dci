package main

import (
	"bytes"
	"encoding/hex"
	"strings"
	"testing"
)

func TestEncodeSimpleTLVs(t *testing.T) {
	reg := makeTestRegistry()

	result := &DecodeResult{
		Config: map[string]interface{}{
			"NetworkAccess":       1,
			"DownstreamFrequency": 360000000,
			"UpstreamChannelId":   5,
		},
		TLVOrder: []string{"NetworkAccess", "DownstreamFrequency", "UpstreamChannelId"},
	}

	encoded, err := Encode(result, reg)
	if err != nil {
		t.Fatal(err)
	}

	expected := buildConfig(
		buildTLV(3, []byte{1}),
		buildTLV(1, []byte{0x15, 0x75, 0x2A, 0x00}),
		buildTLV(2, []byte{5}),
	)

	if !bytes.Equal(encoded, expected) {
		t.Errorf("encoded mismatch:\n  got:  %X\n  want: %X", encoded, expected)
	}
}

func TestEncodeStringTLV(t *testing.T) {
	reg := makeTestRegistry()

	result := &DecodeResult{
		Config: map[string]interface{}{
			"SwUpgradeFilename": "firmware.bin",
		},
		TLVOrder: []string{"SwUpgradeFilename"},
	}

	encoded, err := Encode(result, reg)
	if err != nil {
		t.Fatal(err)
	}

	expected := buildConfig(
		buildTLV(9, []byte("firmware.bin\x00")),
	)

	if !bytes.Equal(encoded, expected) {
		t.Errorf("encoded mismatch:\n  got:  %X\n  want: %X", encoded, expected)
	}
}

func TestEncodeMacAddress(t *testing.T) {
	reg := makeTestRegistry()

	result := &DecodeResult{
		Config: map[string]interface{}{
			"CpeEthernetMacAddress": "001A2B3C4D5E",
		},
		TLVOrder: []string{"CpeEthernetMacAddress"},
	}

	encoded, err := Encode(result, reg)
	if err != nil {
		t.Fatal(err)
	}

	expected := buildConfig(
		buildTLV(14, []byte{0x00, 0x1A, 0x2B, 0x3C, 0x4D, 0x5E}),
	)

	if !bytes.Equal(encoded, expected) {
		t.Errorf("encoded mismatch:\n  got:  %X\n  want: %X", encoded, expected)
	}
}

func TestEncodeIPv4Address(t *testing.T) {
	reg := makeTestRegistry()

	result := &DecodeResult{
		Config: map[string]interface{}{
			"TftpServerProvisionedModemIpv4Address": "10.1.2.3",
		},
		TLVOrder: []string{"TftpServerProvisionedModemIpv4Address"},
	}

	encoded, err := Encode(result, reg)
	if err != nil {
		t.Fatal(err)
	}

	expected := buildConfig(
		buildTLV(20, []byte{10, 1, 2, 3}),
	)

	if !bytes.Equal(encoded, expected) {
		t.Errorf("encoded mismatch:\n  got:  %X\n  want: %X", encoded, expected)
	}
}

func TestEncodeCompoundTLV(t *testing.T) {
	reg := makeTestRegistry()

	result := &DecodeResult{
		Config: map[string]interface{}{
			"UpstreamServiceFlow": []interface{}{
				map[string]interface{}{
					"ServiceFlowReference":    1,
					"QosParamSetType":         7,
					"MaxSustainedTrafficRate": 10000,
					"_tlvOrder":               []string{"ServiceFlowReference", "QosParamSetType", "MaxSustainedTrafficRate"},
				},
			},
		},
		TLVOrder: []string{"UpstreamServiceFlow"},
	}

	encoded, err := Encode(result, reg)
	if err != nil {
		t.Fatal(err)
	}

	subTLVs := append(
		buildTLV(1, []byte{0x00, 0x01}),
		buildTLV(6, []byte{7})...,
	)
	subTLVs = append(subTLVs,
		buildTLV(8, []byte{0x00, 0x00, 0x27, 0x10})...,
	)

	expected := buildConfig(
		buildTLV(24, subTLVs),
	)

	if !bytes.Equal(encoded, expected) {
		t.Errorf("encoded mismatch:\n  got:  %X\n  want: %X", encoded, expected)
	}
}

func TestEncodeRepeatableTLV(t *testing.T) {
	reg := makeTestRegistry()

	result := &DecodeResult{
		Config: map[string]interface{}{
			"UpstreamServiceFlow": []interface{}{
				map[string]interface{}{
					"ServiceFlowReference": 1,
					"_tlvOrder":            []string{"ServiceFlowReference"},
				},
				map[string]interface{}{
					"ServiceFlowReference": 2,
					"_tlvOrder":            []string{"ServiceFlowReference"},
				},
			},
		},
		TLVOrder: []string{"UpstreamServiceFlow"},
	}

	encoded, err := Encode(result, reg)
	if err != nil {
		t.Fatal(err)
	}

	expected := buildConfig(
		buildTLV(24, buildTLV(1, []byte{0x00, 0x01})),
		buildTLV(24, buildTLV(1, []byte{0x00, 0x02})),
	)

	if !bytes.Equal(encoded, expected) {
		t.Errorf("encoded mismatch:\n  got:  %X\n  want: %X", encoded, expected)
	}
}

func TestEncodeTLV10SnmpWriteAccess(t *testing.T) {
	reg := makeTestRegistry()

	result := &DecodeResult{
		Config: map[string]interface{}{
			"SnmpWriteAccessControl": []interface{}{
				map[string]interface{}{
					"oid":    "1.3.6.1",
					"access": 1,
				},
			},
		},
		TLVOrder: []string{"SnmpWriteAccessControl"},
	}

	encoded, err := Encode(result, reg)
	if err != nil {
		t.Fatal(err)
	}

	// Decode the encoded result to verify round-trip.
	decoded, err := Decode(encoded, reg)
	if err != nil {
		t.Fatal(err)
	}

	arr := decoded.Config["SnmpWriteAccessControl"].([]interface{})
	entry := arr[0].(map[string]interface{})
	if entry["oid"] != "1.3.6.1" {
		t.Errorf("expected oid '1.3.6.1', got %v", entry["oid"])
	}
	if entry["access"] != 1 {
		t.Errorf("expected access 1, got %v", entry["access"])
	}
}

func TestEncodeTLV11Snmp(t *testing.T) {
	reg := makeTestRegistry()

	result := &DecodeResult{
		Config: map[string]interface{}{
			"SnmpMibObject": []interface{}{
				map[string]interface{}{
					"oid":   "1.3.6.1.2.1",
					"type":  "Integer",
					"value": "42",
				},
			},
		},
		TLVOrder: []string{"SnmpMibObject"},
	}

	encoded, err := Encode(result, reg)
	if err != nil {
		t.Fatal(err)
	}

	// Decode the encoded result to verify.
	decoded, err := Decode(encoded, reg)
	if err != nil {
		t.Fatal(err)
	}

	snmpArr := decoded.Config["SnmpMibObject"].([]interface{})
	entry := snmpArr[0].(map[string]interface{})
	if entry["oid"] != "1.3.6.1.2.1" {
		t.Errorf("expected oid '1.3.6.1.2.1', got %v", entry["oid"])
	}
	if entry["type"] != "Integer" {
		t.Errorf("expected type 'Integer', got %v", entry["type"])
	}
	if entry["value"] != "42" {
		t.Errorf("expected value '42', got %v", entry["value"])
	}
}

func TestEncodeTLV43GeneralExtension(t *testing.T) {
	reg := makeTestRegistry()

	result := &DecodeResult{
		Config: map[string]interface{}{
			"DocsisExtensionField": []interface{}{
				map[string]interface{}{
					"VendorId":                "FFFFFF",
					"CmLoadBalancingPolicyId": 100,
					"_tlvOrder":               []string{"CmLoadBalancingPolicyId"},
				},
			},
		},
		TLVOrder: []string{"DocsisExtensionField"},
	}

	encoded, err := Encode(result, reg)
	if err != nil {
		t.Fatal(err)
	}

	// VendorId sub-TLV 8 = FF:FF:FF
	vendorId := buildTLV(8, []byte{0xFF, 0xFF, 0xFF})
	// CmLoadBalancingPolicyId sub-TLV 1 = 100
	lbPolicy := buildTLV(1, []byte{0x00, 0x00, 0x00, 0x64})
	tlv43Value := append(vendorId, lbPolicy...)
	expected := buildConfig(buildTLV(43, tlv43Value))

	if !bytes.Equal(encoded, expected) {
		t.Errorf("encoded mismatch:\n  got:  %X\n  want: %X", encoded, expected)
	}
}

func TestEncodeTLV43VendorSpecific(t *testing.T) {
	reg := makeTestRegistry()

	result := &DecodeResult{
		Config: map[string]interface{}{
			"DocsisExtensionField": []interface{}{
				map[string]interface{}{
					"VendorId": "001122",
					"VendorSubTlvs": []interface{}{
						map[string]interface{}{
							"type":  1,
							"value": "AABB",
						},
					},
				},
			},
		},
		TLVOrder: []string{"DocsisExtensionField"},
	}

	encoded, err := Encode(result, reg)
	if err != nil {
		t.Fatal(err)
	}

	vendorId := buildTLV(8, []byte{0x00, 0x11, 0x22})
	vendorSub := buildTLV(1, []byte{0xAA, 0xBB})
	tlv43Value := append(vendorId, vendorSub...)
	expected := buildConfig(buildTLV(43, tlv43Value))

	if !bytes.Equal(encoded, expected) {
		t.Errorf("encoded mismatch:\n  got:  %X\n  want: %X", encoded, expected)
	}
}

func TestEncodeUnknownTLV(t *testing.T) {
	reg := makeTestRegistry()

	result := &DecodeResult{
		Config: map[string]interface{}{
			"NetworkAccess": 1,
			"UnknownTlvs": []interface{}{
				map[string]interface{}{
					"type":  200,
					"value": "ABCD",
				},
			},
		},
		TLVOrder: []string{"NetworkAccess"},
	}

	encoded, err := Encode(result, reg)
	if err != nil {
		t.Fatal(err)
	}

	expected := buildConfig(
		buildTLV(3, []byte{1}),
		buildTLV(200, []byte{0xAB, 0xCD}),
	)

	if !bytes.Equal(encoded, expected) {
		t.Errorf("encoded mismatch:\n  got:  %X\n  want: %X", encoded, expected)
	}
}

func TestEncodeDecodeRoundTrip(t *testing.T) {
	reg := makeTestRegistry()

	// Build original binary.
	original := buildConfig(
		buildTLV(3, []byte{1}),
		buildTLV(1, []byte{0x15, 0x75, 0x2A, 0x00}),
		buildTLV(2, []byte{5}),
		buildTLV(9, []byte("firmware.bin\x00")),
		buildTLV(18, []byte{16}),
	)

	// Decode.
	result, err := Decode(original, reg)
	if err != nil {
		t.Fatal(err)
	}

	// Re-encode.
	encoded, err := Encode(result, reg)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(encoded, original) {
		t.Errorf("round-trip mismatch:\n  got:  %X\n  want: %X", encoded, original)
	}
}

func TestEncodeDecodeRoundTripCompound(t *testing.T) {
	reg := makeTestRegistry()

	subTLVs := append(
		buildTLV(1, []byte{0x00, 0x01}),
		buildTLV(6, []byte{7})...,
	)
	subTLVs = append(subTLVs,
		buildTLV(4, []byte("svc_class\x00"))...,
	)

	original := buildConfig(
		buildTLV(3, []byte{1}),
		buildTLV(24, subTLVs),
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
		t.Errorf("round-trip mismatch:\n  got:  %X\n  want: %X", encoded, original)
	}
}

func TestStripTLVsForComparison(t *testing.T) {
	data := buildConfig(
		buildTLV(3, []byte{1}),
		buildTLV(6, make([]byte, 16)),
		buildTLV(7, make([]byte, 16)),
	)

	stripped := stripTLVsForComparison(data, 6, 7)

	// Should contain only TLV 3 + end-of-data.
	expected := buildConfig(buildTLV(3, []byte{1}))
	if !bytes.Equal(stripped, expected) {
		t.Errorf("strip mismatch:\n  got:  %X\n  want: %X", stripped, expected)
	}
}

func TestStripTLVsForComparisonWithPads(t *testing.T) {
	var data []byte
	data = append(data, 0x00)                     // pad
	data = append(data, buildTLV(3, []byte{1})...) // NetworkAccess
	data = append(data, 0x00)                     // pad
	data = append(data, 0xFF, 0x00)                // end-of-data

	stripped := stripTLVsForComparison(data)

	expected := buildConfig(buildTLV(3, []byte{1}))
	if !bytes.Equal(stripped, expected) {
		t.Errorf("strip mismatch:\n  got:  %X\n  want: %X", stripped, expected)
	}
}

func TestEncodeChunkedTLV_SmallValue(t *testing.T) {
	reg := makeTestRegistry()

	// A small value (< 254 bytes) should produce a single TLV.
	smallHex := strings.Repeat("AB", 100) // 100 bytes
	result := &DecodeResult{
		Config:   map[string]interface{}{"ManufacturerCvc": smallHex},
		TLVOrder: []string{"ManufacturerCvc"},
	}

	encoded, err := Encode(result, reg)
	if err != nil {
		t.Fatal(err)
	}

	valBytes, _ := hex.DecodeString(smallHex)
	expected := buildConfig(buildTLV(32, valBytes))
	if !bytes.Equal(encoded, expected) {
		t.Errorf("small chunked TLV mismatch:\n  got:  %X\n  want: %X", encoded, expected)
	}
}

func TestEncodeChunkedTLV_LargeValue(t *testing.T) {
	reg := makeTestRegistry()

	// A value of 400 bytes should produce 2 chunks: 254 + 146.
	largeHex := strings.Repeat("CD", 400) // 400 bytes
	result := &DecodeResult{
		Config:   map[string]interface{}{"ManufacturerCvc": largeHex},
		TLVOrder: []string{"ManufacturerCvc"},
	}

	encoded, err := Encode(result, reg)
	if err != nil {
		t.Fatal(err)
	}

	valBytes, _ := hex.DecodeString(largeHex)
	var expected []byte
	expected = append(expected, buildTLV(32, valBytes[:254])...)
	expected = append(expected, buildTLV(32, valBytes[254:])...)
	expected = append(expected, 0xFF, 0x00) // end-of-data

	if !bytes.Equal(encoded, expected) {
		t.Errorf("chunked TLV mismatch:\n  got len=%d, want len=%d", len(encoded), len(expected))
	}
}

func TestDecodeChunkedTLV_Reassembly(t *testing.T) {
	reg := makeTestRegistry()

	// Build binary with 2 consecutive TLV 32 chunks.
	chunk1 := bytes.Repeat([]byte{0xAA}, 254)
	chunk2 := bytes.Repeat([]byte{0xBB}, 150)
	data := buildConfig(
		buildTLV(32, chunk1),
		buildTLV(32, chunk2),
	)

	result, err := Decode(data, reg)
	if err != nil {
		t.Fatal(err)
	}

	// Should be reassembled into a single hex string.
	hexVal, ok := result.Config["ManufacturerCvc"].(string)
	if !ok {
		t.Fatalf("expected string for ManufacturerCvc, got %T", result.Config["ManufacturerCvc"])
	}

	expectedHex := strings.ToUpper(hex.EncodeToString(append(chunk1, chunk2...)))
	if hexVal != expectedHex {
		t.Errorf("reassembled value length: got %d, want %d", len(hexVal), len(expectedHex))
	}
}

func TestChunkedTLV_RoundTrip(t *testing.T) {
	reg := makeTestRegistry()

	// Create a large certificate-like value (600 bytes → 3 chunks: 254 + 254 + 92).
	certHex := strings.Repeat("EF", 600)
	result := &DecodeResult{
		Config:   map[string]interface{}{"ManufacturerCvc": certHex},
		TLVOrder: []string{"ManufacturerCvc"},
	}

	// Encode.
	encoded, err := Encode(result, reg)
	if err != nil {
		t.Fatal(err)
	}

	// Decode.
	decoded, err := Decode(encoded, reg)
	if err != nil {
		t.Fatal(err)
	}

	// Round-trip: decoded value must match original.
	got := decoded.Config["ManufacturerCvc"].(string)
	if got != certHex {
		t.Errorf("round-trip mismatch: got len %d, want len %d", len(got), len(certHex))
	}
}

func TestEncodeMICsSkipped(t *testing.T) {
	reg := makeTestRegistry()

	result := &DecodeResult{
		Config: map[string]interface{}{
			"NetworkAccess": 1,
			"CmMic":         "0102030405060708090A0B0C0D0E0F10",
			"CmtsMic":       "1112131415161718191A1B1C1D1E1F20",
		},
		TLVOrder: []string{"NetworkAccess", "CmMic", "CmtsMic"},
	}

	encoded, err := Encode(result, reg)
	if err != nil {
		t.Fatal(err)
	}

	// MIC TLVs should not be in the output.
	expected := buildConfig(buildTLV(3, []byte{1}))
	if !bytes.Equal(encoded, expected) {
		t.Errorf("expected no MIC TLVs in output:\n  got:  %X\n  want: %X", encoded, expected)
	}
}

func TestEncode2ByteLenTopLevel(t *testing.T) {
	reg := makeTestRegistry2ByteLen()

	result := &DecodeResult{
		Config: map[string]interface{}{
			"CmSshServerConfigurationSettings": map[string]interface{}{
				"SshCmCds":            "DEADBEEF",
				"SshCmCdsDownloadUrl": "http://example.com",
				"_tlvOrder":           []string{"SshCmCds", "SshCmCdsDownloadUrl"},
			},
		},
		TLVOrder: []string{"CmSshServerConfigurationSettings"},
	}

	encoded, err := Encode(result, reg)
	if err != nil {
		t.Fatal(err)
	}

	// Build the expected binary manually.
	// Sub-TLV 1 (SshCmCds, 2-byte len): type=1, len=0x00 0x04, value=DEADBEEF
	subTLV := buildTLV2(1, []byte{0xDE, 0xAD, 0xBE, 0xEF})
	// Sub-TLV 2 (SshCmCdsDownloadUrl, 1-byte len)
	subTLV = append(subTLV, buildTLV(2, []byte("http://example.com\x00"))...)
	// TLV 103 (2-byte len) wrapping the sub-TLVs
	var expected []byte
	expected = append(expected, buildTLV2(103, subTLV)...)
	expected = append(expected, 0xFF, 0x00) // end-of-data

	if !bytes.Equal(encoded, expected) {
		t.Errorf("2-byte length encode mismatch:\n  got:  %X\n  want: %X", encoded, expected)
	}
}

func TestEncode2ByteLenLargePayload(t *testing.T) {
	reg := makeTestRegistry2ByteLen()

	// Create a 300-byte hex payload (600 hex chars)
	bigPayload := make([]byte, 300)
	for i := range bigPayload {
		bigPayload[i] = byte(i % 256)
	}
	bigHex := strings.ToUpper(hex.EncodeToString(bigPayload))

	result := &DecodeResult{
		Config: map[string]interface{}{
			"CmSshServerConfigurationSettings": map[string]interface{}{
				"SshCmCds": bigHex,
				"_tlvOrder": []string{"SshCmCds"},
			},
		},
		TLVOrder: []string{"CmSshServerConfigurationSettings"},
	}

	encoded, err := Encode(result, reg)
	if err != nil {
		t.Fatal(err)
	}

	// The TLV 103 header should be: type=103(0x67), length=303 (300 payload + 3 for sub-TLV header)
	// Sub-TLV 1 header: type=1, length_hi=0x01, length_lo=0x2C (300 = 0x012C)
	// TLV 103 total sub-body length = 3 + 300 = 303 = 0x012F
	if encoded[0] != 103 {
		t.Errorf("expected type byte 103, got %d", encoded[0])
	}
	outerLen := int(encoded[1])<<8 | int(encoded[2])
	if outerLen != 303 {
		t.Errorf("expected outer 2-byte length 303, got %d", outerLen)
	}
	// Inner sub-TLV type should be 1
	if encoded[3] != 1 {
		t.Errorf("expected sub-TLV type 1, got %d", encoded[3])
	}
	innerLen := int(encoded[4])<<8 | int(encoded[5])
	if innerLen != 300 {
		t.Errorf("expected inner 2-byte length 300, got %d", innerLen)
	}
}

func TestMakeTLVn(t *testing.T) {
	// 1-byte length
	result1 := makeTLVn(10, []byte{0xAA, 0xBB}, 1)
	if !bytes.Equal(result1, []byte{10, 2, 0xAA, 0xBB}) {
		t.Errorf("makeTLVn(1) mismatch: got %X", result1)
	}

	// 2-byte length
	result2 := makeTLVn(10, []byte{0xAA, 0xBB}, 2)
	if !bytes.Equal(result2, []byte{10, 0x00, 0x02, 0xAA, 0xBB}) {
		t.Errorf("makeTLVn(2) mismatch: got %X", result2)
	}

	// 2-byte length with value > 255 bytes
	bigVal := make([]byte, 300)
	result3 := makeTLVn(5, bigVal, 2)
	if result3[0] != 5 {
		t.Errorf("expected type 5, got %d", result3[0])
	}
	gotLen := int(result3[1])<<8 | int(result3[2])
	if gotLen != 300 {
		t.Errorf("expected length 300, got %d", gotLen)
	}
	if len(result3) != 3+300 {
		t.Errorf("expected total length 303, got %d", len(result3))
	}
}

func TestRoundTrip2ByteLen(t *testing.T) {
	reg := makeTestRegistry2ByteLen()

	// Build original binary with 2-byte length TLVs.
	subTLV := buildTLV2(1, []byte{0xDE, 0xAD, 0xBE, 0xEF})
	subTLV = append(subTLV, buildTLV(2, []byte("test\x00"))...)
	var original []byte
	original = append(original, buildTLV(3, []byte{1})...)
	original = append(original, buildTLV2(103, subTLV)...)
	original = append(original, 0xFF, 0x00)

	// Decode
	result, err := Decode(original, reg)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}

	// Re-encode
	encoded, err := Encode(result, reg)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}

	if !bytes.Equal(encoded, original) {
		t.Errorf("2-byte length round-trip failed:\n  got:  %X\n  want: %X", encoded, original)
	}
}

func TestRoundTrip2ByteLenLargePayload(t *testing.T) {
	reg := makeTestRegistry2ByteLen()

	// 300-byte payload (exceeds 1-byte length limit of 255)
	bigPayload := make([]byte, 300)
	for i := range bigPayload {
		bigPayload[i] = byte(i % 256)
	}

	subTLV := buildTLV2(1, bigPayload)
	var original []byte
	original = append(original, buildTLV2(103, subTLV)...)
	original = append(original, 0xFF, 0x00)

	// Decode
	result, err := Decode(original, reg)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}

	// Re-encode
	encoded, err := Encode(result, reg)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}

	if !bytes.Equal(encoded, original) {
		t.Errorf("large payload 2-byte length round-trip failed:\n  got len=%d, want len=%d", len(encoded), len(original))
	}
}

func TestPadToAlignment(t *testing.T) {
	tests := []struct {
		inputLen int
		align    int
		wantLen  int
	}{
		{0, 4, 0},   // empty stays empty
		{1, 4, 4},   // 1 -> pad to 4
		{2, 4, 4},   // 2 -> pad to 4
		{3, 4, 4},   // 3 -> pad to 4
		{4, 4, 4},   // 4 -> already aligned
		{5, 4, 8},   // 5 -> pad to 8
		{7, 4, 8},   // 7 -> pad to 8
		{8, 4, 8},   // 8 -> already aligned
		{10, 1, 10}, // alignment=1 is a no-op
		{10, 0, 10}, // alignment=0 is a no-op
	}
	for _, tc := range tests {
		data := make([]byte, tc.inputLen)
		result := PadToAlignment(data, tc.align)
		if len(result) != tc.wantLen {
			t.Errorf("PadToAlignment(len=%d, align=%d): got len=%d, want %d",
				tc.inputLen, tc.align, len(result), tc.wantLen)
		}
		// Verify pad bytes are zero.
		for i := tc.inputLen; i < len(result); i++ {
			if result[i] != 0 {
				t.Errorf("PadToAlignment: pad byte at index %d is 0x%02X, want 0x00", i, result[i])
			}
		}
	}
}
