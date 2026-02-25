package main

import (
	"bytes"
	"testing"
)

func TestEncodeValueUint8(t *testing.T) {
	b, err := EncodeValue(42, DataTypeUint8)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, []byte{42}) {
		t.Errorf("expected [42], got %v", b)
	}
}

func TestEncodeValueUint8Float64(t *testing.T) {
	// JSON unmarshal produces float64
	b, err := EncodeValue(float64(42), DataTypeUint8)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, []byte{42}) {
		t.Errorf("expected [42], got %v", b)
	}
}

func TestEncodeValueUint8OutOfRange(t *testing.T) {
	_, err := EncodeValue(256, DataTypeUint8)
	if err == nil {
		t.Fatal("expected error for out of range")
	}
}

func TestEncodeValueUint16(t *testing.T) {
	b, err := EncodeValue(256, DataTypeUint16)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, []byte{0x01, 0x00}) {
		t.Errorf("expected [0x01, 0x00], got %X", b)
	}
}

func TestEncodeValueUint32(t *testing.T) {
	b, err := EncodeValue(1000000, DataTypeUint32)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, []byte{0x00, 0x0F, 0x42, 0x40}) {
		t.Errorf("expected [00 0F 42 40], got %X", b)
	}
}

func TestEncodeValueString(t *testing.T) {
	b, err := EncodeValue("hello", DataTypeString)
	if err != nil {
		t.Fatal(err)
	}
	expected := []byte("hello\x00")
	if !bytes.Equal(b, expected) {
		t.Errorf("expected %q, got %q", expected, b)
	}
}

func TestEncodeValueHexString(t *testing.T) {
	b, err := EncodeValue("DEADBEEF", DataTypeHexString)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, []byte{0xDE, 0xAD, 0xBE, 0xEF}) {
		t.Errorf("expected [DEADBEEF], got %X", b)
	}
}

func TestEncodeValueMacAddress(t *testing.T) {
	b, err := EncodeValue("001A2B3C4D5E", DataTypeMacAddress)
	if err != nil {
		t.Fatal(err)
	}
	expected := []byte{0x00, 0x1A, 0x2B, 0x3C, 0x4D, 0x5E}
	if !bytes.Equal(b, expected) {
		t.Errorf("expected %X, got %X", expected, b)
	}
}

func TestEncodeValueIPv4Address(t *testing.T) {
	b, err := EncodeValue("192.168.1.1", DataTypeIPv4Address)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, []byte{192, 168, 1, 1}) {
		t.Errorf("expected [192 168 1 1], got %v", b)
	}
}

func TestEncodeValueIPv6Address(t *testing.T) {
	b, err := EncodeValue("2001:db8::1", DataTypeIPv6Address)
	if err != nil {
		t.Fatal(err)
	}
	expected := []byte{
		0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
	}
	if !bytes.Equal(b, expected) {
		t.Errorf("expected %X, got %X", expected, b)
	}
}

func TestEncodeValueOID(t *testing.T) {
	b, err := EncodeValue("1.3.6.1.2.1", DataTypeOID)
	if err != nil {
		t.Fatal(err)
	}
	expected := []byte{0x2B, 0x06, 0x01, 0x02, 0x01}
	if !bytes.Equal(b, expected) {
		t.Errorf("expected %X, got %X", expected, b)
	}
}

func TestEncodeValueOIDLargeComponent(t *testing.T) {
	b, err := EncodeValue("1.3.6.1.2.1.200", DataTypeOID)
	if err != nil {
		t.Fatal(err)
	}
	expected := []byte{0x2B, 0x06, 0x01, 0x02, 0x01, 0x81, 0x48}
	if !bytes.Equal(b, expected) {
		t.Errorf("expected %X, got %X", expected, b)
	}
}

// Round-trip tests: EncodeValue(DecodeValue(bytes)) == bytes for all data types.

func TestRoundTripUint8(t *testing.T) {
	original := []byte{42}
	val, err := DecodeValue(original, DataTypeUint8)
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := EncodeValue(val, DataTypeUint8)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded, original) {
		t.Errorf("round-trip failed: got %X, want %X", encoded, original)
	}
}

func TestRoundTripUint16(t *testing.T) {
	original := []byte{0x01, 0x00}
	val, err := DecodeValue(original, DataTypeUint16)
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := EncodeValue(val, DataTypeUint16)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded, original) {
		t.Errorf("round-trip failed: got %X, want %X", encoded, original)
	}
}

func TestRoundTripUint32(t *testing.T) {
	original := []byte{0x15, 0x75, 0x2A, 0x00}
	val, err := DecodeValue(original, DataTypeUint32)
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := EncodeValue(val, DataTypeUint32)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded, original) {
		t.Errorf("round-trip failed: got %X, want %X", encoded, original)
	}
}

func TestRoundTripString(t *testing.T) {
	original := []byte("firmware.bin\x00")
	val, err := DecodeValue(original, DataTypeString)
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := EncodeValue(val, DataTypeString)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded, original) {
		t.Errorf("round-trip failed: got %q, want %q", encoded, original)
	}
}

func TestRoundTripHexString(t *testing.T) {
	original := []byte{0xDE, 0xAD, 0xBE, 0xEF}
	val, err := DecodeValue(original, DataTypeHexString)
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := EncodeValue(val, DataTypeHexString)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded, original) {
		t.Errorf("round-trip failed: got %X, want %X", encoded, original)
	}
}

func TestRoundTripMacAddress(t *testing.T) {
	original := []byte{0x00, 0x1A, 0x2B, 0x3C, 0x4D, 0x5E}
	val, err := DecodeValue(original, DataTypeMacAddress)
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := EncodeValue(val, DataTypeMacAddress)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded, original) {
		t.Errorf("round-trip failed: got %X, want %X", encoded, original)
	}
}

func TestRoundTripIPv4Address(t *testing.T) {
	original := []byte{10, 1, 2, 3}
	val, err := DecodeValue(original, DataTypeIPv4Address)
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := EncodeValue(val, DataTypeIPv4Address)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded, original) {
		t.Errorf("round-trip failed: got %X, want %X", encoded, original)
	}
}

func TestRoundTripIPv6Address(t *testing.T) {
	original := []byte{
		0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
	}
	val, err := DecodeValue(original, DataTypeIPv6Address)
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := EncodeValue(val, DataTypeIPv6Address)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded, original) {
		t.Errorf("round-trip failed: got %X, want %X", encoded, original)
	}
}

func TestRoundTripOID(t *testing.T) {
	original := []byte{0x2B, 0x06, 0x01, 0x02, 0x01}
	val, err := DecodeValue(original, DataTypeOID)
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := EncodeValue(val, DataTypeOID)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded, original) {
		t.Errorf("round-trip failed: got %X, want %X", encoded, original)
	}
}

func TestRoundTripOIDLargeComponent(t *testing.T) {
	original := []byte{0x2B, 0x06, 0x01, 0x02, 0x01, 0x81, 0x48}
	val, err := DecodeValue(original, DataTypeOID)
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := EncodeValue(val, DataTypeOID)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded, original) {
		t.Errorf("round-trip failed: got %X, want %X", encoded, original)
	}
}

func TestEncodeValueUnknownType(t *testing.T) {
	_, err := EncodeValue("foo", DataType("foobar"))
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
}
