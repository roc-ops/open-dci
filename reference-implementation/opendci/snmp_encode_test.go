package opendci

import (
	"bytes"
	"testing"
)

func TestEncodeSnmpVarbindInteger(t *testing.T) {
	entry := map[string]interface{}{
		"oid":   "1.3.6.1.2.1",
		"type":  "Integer",
		"value": "42",
	}

	encoded, err := EncodeSnmpVarbind(entry)
	if err != nil {
		t.Fatal(err)
	}

	// Decode it back to verify round-trip.
	vb, err := DecodeSnmpVarbind(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if vb.OID != "1.3.6.1.2.1" {
		t.Errorf("expected OID '1.3.6.1.2.1', got %q", vb.OID)
	}
	if vb.Type != "Integer" {
		t.Errorf("expected type 'Integer', got %q", vb.Type)
	}
	if vb.Value != "42" {
		t.Errorf("expected value '42', got %q", vb.Value)
	}
}

func TestEncodeSnmpVarbindNegativeInteger(t *testing.T) {
	entry := map[string]interface{}{
		"oid":   "1.3.6.1.2.1",
		"type":  "Integer",
		"value": "-1",
	}

	encoded, err := EncodeSnmpVarbind(entry)
	if err != nil {
		t.Fatal(err)
	}

	vb, err := DecodeSnmpVarbind(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if vb.Value != "-1" {
		t.Errorf("expected value '-1', got %q", vb.Value)
	}
}

func TestEncodeSnmpVarbindString(t *testing.T) {
	entry := map[string]interface{}{
		"oid":   "1.3.6.1.2.1",
		"type":  "String",
		"value": "hello",
	}

	encoded, err := EncodeSnmpVarbind(entry)
	if err != nil {
		t.Fatal(err)
	}

	vb, err := DecodeSnmpVarbind(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if vb.Type != "String" {
		t.Errorf("expected type 'String', got %q", vb.Type)
	}
	if vb.Value != "hello" {
		t.Errorf("expected value 'hello', got %q", vb.Value)
	}
}

func TestEncodeSnmpVarbindIPAddress(t *testing.T) {
	entry := map[string]interface{}{
		"oid":   "1.3.6.1.2.1",
		"type":  "IPAddress",
		"value": "10.0.0.1",
	}

	encoded, err := EncodeSnmpVarbind(entry)
	if err != nil {
		t.Fatal(err)
	}

	vb, err := DecodeSnmpVarbind(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if vb.Type != "IPAddress" {
		t.Errorf("expected type 'IPAddress', got %q", vb.Type)
	}
	if vb.Value != "10.0.0.1" {
		t.Errorf("expected value '10.0.0.1', got %q", vb.Value)
	}
}

func TestEncodeSnmpVarbindCounter32(t *testing.T) {
	entry := map[string]interface{}{
		"oid":   "1.3.6.1.2.1",
		"type":  "Counter32",
		"value": "1000",
	}

	encoded, err := EncodeSnmpVarbind(entry)
	if err != nil {
		t.Fatal(err)
	}

	vb, err := DecodeSnmpVarbind(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if vb.Type != "Counter32" {
		t.Errorf("expected type 'Counter32', got %q", vb.Type)
	}
	if vb.Value != "1000" {
		t.Errorf("expected value '1000', got %q", vb.Value)
	}
}

func TestEncodeSnmpVarbindGauge32(t *testing.T) {
	entry := map[string]interface{}{
		"oid":   "1.3.6.1.2.1",
		"type":  "Gauge32",
		"value": "100",
	}

	encoded, err := EncodeSnmpVarbind(entry)
	if err != nil {
		t.Fatal(err)
	}

	vb, err := DecodeSnmpVarbind(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if vb.Type != "Gauge32" {
		t.Errorf("expected type 'Gauge32', got %q", vb.Type)
	}
	if vb.Value != "100" {
		t.Errorf("expected value '100', got %q", vb.Value)
	}
}

func TestEncodeSnmpVarbindTimeTicks(t *testing.T) {
	entry := map[string]interface{}{
		"oid":   "1.3.6.1.2.1",
		"type":  "TimeTicks",
		"value": "86400",
	}

	encoded, err := EncodeSnmpVarbind(entry)
	if err != nil {
		t.Fatal(err)
	}

	vb, err := DecodeSnmpVarbind(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if vb.Type != "TimeTicks" {
		t.Errorf("expected type 'TimeTicks', got %q", vb.Type)
	}
	if vb.Value != "86400" {
		t.Errorf("expected value '86400', got %q", vb.Value)
	}
}

func TestEncodeSnmpVarbindCounter64(t *testing.T) {
	entry := map[string]interface{}{
		"oid":   "1.3.6.1.2.1",
		"type":  "Counter64",
		"value": "4294967296",
	}

	encoded, err := EncodeSnmpVarbind(entry)
	if err != nil {
		t.Fatal(err)
	}

	vb, err := DecodeSnmpVarbind(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if vb.Type != "Counter64" {
		t.Errorf("expected type 'Counter64', got %q", vb.Type)
	}
	if vb.Value != "4294967296" {
		t.Errorf("expected value '4294967296', got %q", vb.Value)
	}
}

func TestEncodeSnmpVarbindOIDValue(t *testing.T) {
	entry := map[string]interface{}{
		"oid":   "1.3.6.1.2.1",
		"type":  "OID",
		"value": "1.3.6.1.4",
	}

	encoded, err := EncodeSnmpVarbind(entry)
	if err != nil {
		t.Fatal(err)
	}

	vb, err := DecodeSnmpVarbind(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if vb.Type != "OID" {
		t.Errorf("expected type 'OID', got %q", vb.Type)
	}
	if vb.Value != "1.3.6.1.4" {
		t.Errorf("expected value '1.3.6.1.4', got %q", vb.Value)
	}
}

func TestEncodeSnmpVarbindNull(t *testing.T) {
	entry := map[string]interface{}{
		"oid":   "1.3.6.1.2.1",
		"type":  "Null",
		"value": "",
	}

	encoded, err := EncodeSnmpVarbind(entry)
	if err != nil {
		t.Fatal(err)
	}

	vb, err := DecodeSnmpVarbind(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if vb.Type != "Null" {
		t.Errorf("expected type 'Null', got %q", vb.Type)
	}
}

// Byte-exact round-trip tests using the test helper from snmp_test.go.

func TestSnmpVarbindRoundTripInteger(t *testing.T) {
	oidBytes := []byte{0x2B, 0x06, 0x01, 0x02, 0x01}
	original := buildVarbind(oidBytes, tagInteger, []byte{0x2A})

	vb, err := DecodeSnmpVarbind(original)
	if err != nil {
		t.Fatal(err)
	}

	entry := map[string]interface{}{
		"oid":   vb.OID,
		"type":  vb.Type,
		"value": vb.Value,
	}

	encoded, err := EncodeSnmpVarbind(entry)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(encoded, original) {
		t.Errorf("round-trip failed:\n  got:  %X\n  want: %X", encoded, original)
	}
}

func TestSnmpVarbindRoundTripNegativeInteger(t *testing.T) {
	oidBytes := []byte{0x2B, 0x06, 0x01, 0x02, 0x01}
	original := buildVarbind(oidBytes, tagInteger, []byte{0xFF})

	vb, err := DecodeSnmpVarbind(original)
	if err != nil {
		t.Fatal(err)
	}

	entry := map[string]interface{}{
		"oid":   vb.OID,
		"type":  vb.Type,
		"value": vb.Value,
	}

	encoded, err := EncodeSnmpVarbind(entry)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(encoded, original) {
		t.Errorf("round-trip failed:\n  got:  %X\n  want: %X", encoded, original)
	}
}

func TestSnmpVarbindRoundTripIPAddress(t *testing.T) {
	oidBytes := []byte{0x2B, 0x06, 0x01, 0x02, 0x01}
	original := buildVarbind(oidBytes, tagIPAddress, []byte{10, 0, 0, 1})

	vb, err := DecodeSnmpVarbind(original)
	if err != nil {
		t.Fatal(err)
	}

	entry := map[string]interface{}{
		"oid":   vb.OID,
		"type":  vb.Type,
		"value": vb.Value,
	}

	encoded, err := EncodeSnmpVarbind(entry)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(encoded, original) {
		t.Errorf("round-trip failed:\n  got:  %X\n  want: %X", encoded, original)
	}
}

func TestSnmpVarbindRoundTripCounter32(t *testing.T) {
	oidBytes := []byte{0x2B, 0x06, 0x01, 0x02, 0x01}
	original := buildVarbind(oidBytes, tagCounter32, []byte{0x03, 0xE8})

	vb, err := DecodeSnmpVarbind(original)
	if err != nil {
		t.Fatal(err)
	}

	entry := map[string]interface{}{
		"oid":   vb.OID,
		"type":  vb.Type,
		"value": vb.Value,
	}

	encoded, err := EncodeSnmpVarbind(entry)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(encoded, original) {
		t.Errorf("round-trip failed:\n  got:  %X\n  want: %X", encoded, original)
	}
}

func TestSnmpVarbindRoundTripGauge32(t *testing.T) {
	oidBytes := []byte{0x2B, 0x06, 0x01, 0x02, 0x01}
	original := buildVarbind(oidBytes, tagGauge32, []byte{0x00, 0x64})

	vb, err := DecodeSnmpVarbind(original)
	if err != nil {
		t.Fatal(err)
	}

	entry := map[string]interface{}{
		"oid":   vb.OID,
		"type":  vb.Type,
		"value": vb.Value,
	}

	encoded, err := EncodeSnmpVarbind(entry)
	if err != nil {
		t.Fatal(err)
	}

	// Note: Gauge32 with leading 0x00 byte — the decoder reads "100",
	// the encoder will output {0x64} (1 byte) not {0x00, 0x64} (2 bytes).
	// This is still semantically correct — the leading zero is redundant.
	vb2, err := DecodeSnmpVarbind(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if vb2.Value != vb.Value {
		t.Errorf("values differ: got %q, want %q", vb2.Value, vb.Value)
	}
}

func TestEncodeSnmpVarbindMissingOID(t *testing.T) {
	entry := map[string]interface{}{
		"type":  "Integer",
		"value": "42",
	}
	_, err := EncodeSnmpVarbind(entry)
	if err == nil {
		t.Fatal("expected error for missing oid")
	}
}

func TestEncodeSnmpVarbindMissingType(t *testing.T) {
	entry := map[string]interface{}{
		"oid":   "1.3.6.1.2.1",
		"value": "42",
	}
	_, err := EncodeSnmpVarbind(entry)
	if err == nil {
		t.Fatal("expected error for missing type")
	}
}

// TLV 10 (SNMP Write Access Control) encode tests.

func TestEncodeSnmpWriteAccess(t *testing.T) {
	entry := map[string]interface{}{
		"oid":    "1.3.6.1",
		"access": 1,
	}
	data, err := EncodeSnmpWriteAccess(entry)
	if err != nil {
		t.Fatal(err)
	}
	// Expected: BER OID (06 03 2B 06 01) + access byte (01)
	expected := []byte{0x06, 0x03, 0x2B, 0x06, 0x01, 0x01}
	if !bytes.Equal(data, expected) {
		t.Errorf("expected %X, got %X", expected, data)
	}
}

func TestEncodeSnmpWriteAccessDeny(t *testing.T) {
	entry := map[string]interface{}{
		"oid":    "1.3.6.1.2.1",
		"access": 0,
	}
	data, err := EncodeSnmpWriteAccess(entry)
	if err != nil {
		t.Fatal(err)
	}
	// Expected: BER OID (06 05 2B 06 01 02 01) + access byte (00)
	expected := []byte{0x06, 0x05, 0x2B, 0x06, 0x01, 0x02, 0x01, 0x00}
	if !bytes.Equal(data, expected) {
		t.Errorf("expected %X, got %X", expected, data)
	}
}

func TestSnmpWriteAccessRoundTrip(t *testing.T) {
	entry := map[string]interface{}{
		"oid":    "1.3.6.1.2.1.69.1.2",
		"access": 1,
	}
	encoded, err := EncodeSnmpWriteAccess(entry)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodeSnmpWriteAccess(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if decoded["oid"] != entry["oid"] {
		t.Errorf("OID mismatch: expected %q, got %q", entry["oid"], decoded["oid"])
	}
	if decoded["access"] != entry["access"] {
		t.Errorf("access mismatch: expected %v, got %v", entry["access"], decoded["access"])
	}
}

func TestEncodeSnmpWriteAccessMissingOID(t *testing.T) {
	entry := map[string]interface{}{
		"access": 1,
	}
	_, err := EncodeSnmpWriteAccess(entry)
	if err == nil {
		t.Fatal("expected error for missing OID")
	}
}

func TestEncodeBERLength(t *testing.T) {
	tests := []struct {
		length   int
		expected []byte
	}{
		{0, []byte{0x00}},
		{5, []byte{0x05}},
		{127, []byte{0x7F}},
		{128, []byte{0x81, 0x80}},
		{255, []byte{0x81, 0xFF}},
		{256, []byte{0x82, 0x01, 0x00}},
	}
	for _, tc := range tests {
		got := encodeBERLength(tc.length)
		if !bytes.Equal(got, tc.expected) {
			t.Errorf("encodeBERLength(%d): got %X, want %X", tc.length, got, tc.expected)
		}
	}
}
