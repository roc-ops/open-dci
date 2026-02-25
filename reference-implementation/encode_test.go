package main

import (
	"bytes"
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
