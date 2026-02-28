package opendci

import (
	"testing"
)

// buildBERLength encodes a BER length value.
func buildBERLength(length int) []byte {
	if length < 0x80 {
		return []byte{byte(length)}
	}
	if length <= 0xFF {
		return []byte{0x81, byte(length)}
	}
	return []byte{0x82, byte(length >> 8), byte(length)}
}

// buildBERTLV builds a BER TLV: tag + length + value.
func buildBERTLV(tag byte, value []byte) []byte {
	result := []byte{tag}
	result = append(result, buildBERLength(len(value))...)
	result = append(result, value...)
	return result
}

// buildVarbind constructs a complete BER-encoded varbind SEQUENCE.
func buildVarbind(oidBytes []byte, valueTag byte, valueBytes []byte) []byte {
	oidTLV := buildBERTLV(tagOID, oidBytes)
	valTLV := buildBERTLV(valueTag, valueBytes)
	inner := append(oidTLV, valTLV...)
	return buildBERTLV(0x30, inner) // SEQUENCE
}

func TestDecodeSnmpVarbindInteger(t *testing.T) {
	// OID 1.3.6.1.2.1 = {0x2B, 0x06, 0x01, 0x02, 0x01}
	oidBytes := []byte{0x2B, 0x06, 0x01, 0x02, 0x01}
	// Integer value 42 = {0x2A}
	varbindData := buildVarbind(oidBytes, tagInteger, []byte{0x2A})

	vb, err := DecodeSnmpVarbind(varbindData)
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

func TestDecodeSnmpVarbindNegativeInteger(t *testing.T) {
	oidBytes := []byte{0x2B, 0x06, 0x01, 0x02, 0x01}
	// -1 in BER = 0xFF
	varbindData := buildVarbind(oidBytes, tagInteger, []byte{0xFF})

	vb, err := DecodeSnmpVarbind(varbindData)
	if err != nil {
		t.Fatal(err)
	}

	if vb.Value != "-1" {
		t.Errorf("expected value '-1', got %q", vb.Value)
	}
}

func TestDecodeSnmpVarbindString(t *testing.T) {
	oidBytes := []byte{0x2B, 0x06, 0x01, 0x02, 0x01}
	varbindData := buildVarbind(oidBytes, tagOctetString, []byte("hello"))

	vb, err := DecodeSnmpVarbind(varbindData)
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

func TestDecodeSnmpVarbindHexString(t *testing.T) {
	oidBytes := []byte{0x2B, 0x06, 0x01, 0x02, 0x01}
	// Non-printable bytes
	varbindData := buildVarbind(oidBytes, tagOctetString, []byte{0x00, 0xFF, 0x01})

	vb, err := DecodeSnmpVarbind(varbindData)
	if err != nil {
		t.Fatal(err)
	}

	if vb.Type != "HexString" {
		t.Errorf("expected type 'HexString', got %q", vb.Type)
	}
	if vb.Value != "00FF01" {
		t.Errorf("expected value '00FF01', got %q", vb.Value)
	}
}

func TestDecodeSnmpVarbindIPAddress(t *testing.T) {
	oidBytes := []byte{0x2B, 0x06, 0x01, 0x02, 0x01}
	varbindData := buildVarbind(oidBytes, tagIPAddress, []byte{10, 0, 0, 1})

	vb, err := DecodeSnmpVarbind(varbindData)
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

func TestDecodeSnmpVarbindCounter32(t *testing.T) {
	oidBytes := []byte{0x2B, 0x06, 0x01, 0x02, 0x01}
	// Counter32 value 1000 = 0x03E8
	varbindData := buildVarbind(oidBytes, tagCounter32, []byte{0x03, 0xE8})

	vb, err := DecodeSnmpVarbind(varbindData)
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

func TestDecodeSnmpVarbindGauge32(t *testing.T) {
	oidBytes := []byte{0x2B, 0x06, 0x01, 0x02, 0x01}
	varbindData := buildVarbind(oidBytes, tagGauge32, []byte{0x00, 0x64})

	vb, err := DecodeSnmpVarbind(varbindData)
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

func TestDecodeSnmpVarbindTimeTicks(t *testing.T) {
	oidBytes := []byte{0x2B, 0x06, 0x01, 0x02, 0x01}
	varbindData := buildVarbind(oidBytes, tagTimeTicks, []byte{0x00, 0x01, 0x51, 0x80})

	vb, err := DecodeSnmpVarbind(varbindData)
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

func TestDecodeSnmpVarbindCounter64(t *testing.T) {
	oidBytes := []byte{0x2B, 0x06, 0x01, 0x02, 0x01}
	varbindData := buildVarbind(oidBytes, tagCounter64, []byte{0x01, 0x00, 0x00, 0x00, 0x00})

	vb, err := DecodeSnmpVarbind(varbindData)
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

func TestDecodeSnmpVarbindOIDValue(t *testing.T) {
	oidBytes := []byte{0x2B, 0x06, 0x01, 0x02, 0x01}
	// Value is another OID: 1.3.6.1.4
	valOID := []byte{0x2B, 0x06, 0x01, 0x04}
	varbindData := buildVarbind(oidBytes, tagOID, valOID)

	vb, err := DecodeSnmpVarbind(varbindData)
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

func TestDecodeSnmpVarbindTooShort(t *testing.T) {
	_, err := DecodeSnmpVarbind([]byte{0x30})
	if err == nil {
		t.Fatal("expected error for truncated data")
	}
}

func TestDecodeSnmpVarbindBadTag(t *testing.T) {
	_, err := DecodeSnmpVarbind([]byte{0x31, 0x00})
	if err == nil {
		t.Fatal("expected error for wrong outer tag")
	}
}

// TLV 10 (SNMP Write Access Control) decode tests.

func TestDecodeSnmpWriteAccess(t *testing.T) {
	// OID 1.3.6.1 with access=1
	// BER OID: 06 03 2B 06 01   (tag=0x06, len=3, value=1.3.6.1)
	// Access flag: 01
	data := []byte{0x06, 0x03, 0x2B, 0x06, 0x01, 0x01}
	entry, err := DecodeSnmpWriteAccess(data)
	if err != nil {
		t.Fatal(err)
	}
	if entry["oid"] != "1.3.6.1" {
		t.Errorf("expected OID '1.3.6.1', got %q", entry["oid"])
	}
	if entry["access"] != 1 {
		t.Errorf("expected access 1, got %v", entry["access"])
	}
}

func TestDecodeSnmpWriteAccessDeny(t *testing.T) {
	// OID 1.3.6.1.2.1 with access=0 (deny)
	// BER OID: 06 05 2B 06 01 02 01
	data := []byte{0x06, 0x05, 0x2B, 0x06, 0x01, 0x02, 0x01, 0x00}
	entry, err := DecodeSnmpWriteAccess(data)
	if err != nil {
		t.Fatal(err)
	}
	if entry["oid"] != "1.3.6.1.2.1" {
		t.Errorf("expected OID '1.3.6.1.2.1', got %q", entry["oid"])
	}
	if entry["access"] != 0 {
		t.Errorf("expected access 0, got %v", entry["access"])
	}
}

func TestDecodeSnmpWriteAccessTooShort(t *testing.T) {
	_, err := DecodeSnmpWriteAccess([]byte{0x06, 0x01})
	if err == nil {
		t.Fatal("expected error for too-short data")
	}
}

func TestDecodeSnmpWriteAccessBadTag(t *testing.T) {
	// Wrong tag (0x30 instead of 0x06)
	_, err := DecodeSnmpWriteAccess([]byte{0x30, 0x03, 0x2B, 0x06, 0x01, 0x01})
	if err == nil {
		t.Fatal("expected error for wrong OID tag")
	}
}

func TestBerLengthShort(t *testing.T) {
	l, consumed, err := berLength([]byte{0x05})
	if err != nil {
		t.Fatal(err)
	}
	if l != 5 {
		t.Errorf("expected length 5, got %d", l)
	}
	if consumed != 1 {
		t.Errorf("expected 1 byte consumed, got %d", consumed)
	}
}

func TestBerLengthLong(t *testing.T) {
	// 0x81 0x80 = 128
	l, consumed, err := berLength([]byte{0x81, 0x80})
	if err != nil {
		t.Fatal(err)
	}
	if l != 128 {
		t.Errorf("expected length 128, got %d", l)
	}
	if consumed != 2 {
		t.Errorf("expected 2 bytes consumed, got %d", consumed)
	}
}

func TestBerLengthLongTwoBytes(t *testing.T) {
	// 0x82 0x01 0x00 = 256
	l, consumed, err := berLength([]byte{0x82, 0x01, 0x00})
	if err != nil {
		t.Fatal(err)
	}
	if l != 256 {
		t.Errorf("expected length 256, got %d", l)
	}
	if consumed != 3 {
		t.Errorf("expected 3 bytes consumed, got %d", consumed)
	}
}
