package main

import (
	"testing"
)

func TestDecodeValueUint8(t *testing.T) {
	val, err := DecodeValue([]byte{42}, DataTypeUint8)
	if err != nil {
		t.Fatal(err)
	}
	if val != 42 {
		t.Errorf("expected 42, got %v", val)
	}
}

func TestDecodeValueUint8WrongLength(t *testing.T) {
	_, err := DecodeValue([]byte{1, 2}, DataTypeUint8)
	if err == nil {
		t.Fatal("expected error for wrong length")
	}
}

func TestDecodeValueUint16(t *testing.T) {
	val, err := DecodeValue([]byte{0x01, 0x00}, DataTypeUint16)
	if err != nil {
		t.Fatal(err)
	}
	if val != 256 {
		t.Errorf("expected 256, got %v", val)
	}
}

func TestDecodeValueUint16WrongLength(t *testing.T) {
	_, err := DecodeValue([]byte{1}, DataTypeUint16)
	if err == nil {
		t.Fatal("expected error for wrong length")
	}
}

func TestDecodeValueUint32(t *testing.T) {
	val, err := DecodeValue([]byte{0x00, 0x0F, 0x42, 0x40}, DataTypeUint32)
	if err != nil {
		t.Fatal(err)
	}
	if val != 1000000 {
		t.Errorf("expected 1000000, got %v", val)
	}
}

func TestDecodeValueUint32WrongLength(t *testing.T) {
	_, err := DecodeValue([]byte{1, 2, 3}, DataTypeUint32)
	if err == nil {
		t.Fatal("expected error for wrong length")
	}
}

func TestDecodeValueString(t *testing.T) {
	val, err := DecodeValue([]byte("hello\x00"), DataTypeString)
	if err != nil {
		t.Fatal(err)
	}
	if val != "hello" {
		t.Errorf("expected 'hello', got %q", val)
	}
}

func TestDecodeValueStringNoNull(t *testing.T) {
	val, err := DecodeValue([]byte("hello"), DataTypeString)
	if err != nil {
		t.Fatal(err)
	}
	if val != "hello" {
		t.Errorf("expected 'hello', got %q", val)
	}
}

func TestDecodeValueHexString(t *testing.T) {
	val, err := DecodeValue([]byte{0xDE, 0xAD, 0xBE, 0xEF}, DataTypeHexString)
	if err != nil {
		t.Fatal(err)
	}
	if val != "DEADBEEF" {
		t.Errorf("expected 'DEADBEEF', got %v", val)
	}
}

func TestDecodeValueMacAddress(t *testing.T) {
	val, err := DecodeValue([]byte{0x00, 0x1A, 0x2B, 0x3C, 0x4D, 0x5E}, DataTypeMacAddress)
	if err != nil {
		t.Fatal(err)
	}
	if val != "001A2B3C4D5E" {
		t.Errorf("expected '001A2B3C4D5E', got %v", val)
	}
}

func TestDecodeValueMacAddressWrongLength(t *testing.T) {
	_, err := DecodeValue([]byte{0x00, 0x01, 0x02}, DataTypeMacAddress)
	if err == nil {
		t.Fatal("expected error for wrong length")
	}
}

func TestDecodeValueIPv4Address(t *testing.T) {
	val, err := DecodeValue([]byte{192, 168, 1, 1}, DataTypeIPv4Address)
	if err != nil {
		t.Fatal(err)
	}
	if val != "192.168.1.1" {
		t.Errorf("expected '192.168.1.1', got %v", val)
	}
}

func TestDecodeValueIPv4AddressWrongLength(t *testing.T) {
	_, err := DecodeValue([]byte{1, 2, 3}, DataTypeIPv4Address)
	if err == nil {
		t.Fatal("expected error for wrong length")
	}
}

func TestDecodeValueIPv6Address(t *testing.T) {
	ipv6 := []byte{
		0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
	}
	val, err := DecodeValue(ipv6, DataTypeIPv6Address)
	if err != nil {
		t.Fatal(err)
	}
	if val != "2001:db8::1" {
		t.Errorf("expected '2001:db8::1', got %v", val)
	}
}

func TestDecodeValueIPv6AddressWrongLength(t *testing.T) {
	_, err := DecodeValue([]byte{1, 2, 3, 4}, DataTypeIPv6Address)
	if err == nil {
		t.Fatal("expected error for wrong length")
	}
}

func TestDecodeValueOID(t *testing.T) {
	// OID 1.3.6.1.2.1 encoded as BER OID bytes (no tag/length)
	// 1.3 -> 43, 6 -> 6, 1 -> 1, 2 -> 2, 1 -> 1
	oidBytes := []byte{0x2B, 0x06, 0x01, 0x02, 0x01}
	val, err := DecodeValue(oidBytes, DataTypeOID)
	if err != nil {
		t.Fatal(err)
	}
	if val != "1.3.6.1.2.1" {
		t.Errorf("expected '1.3.6.1.2.1', got %v", val)
	}
}

func TestDecodeValueOIDLargeComponent(t *testing.T) {
	// OID 1.3.6.1.2.1.1 with a multi-byte component: 200 = 0x81 0x48
	oidBytes := []byte{0x2B, 0x06, 0x01, 0x02, 0x01, 0x81, 0x48}
	val, err := DecodeValue(oidBytes, DataTypeOID)
	if err != nil {
		t.Fatal(err)
	}
	expected := "1.3.6.1.2.1.200"
	if val != expected {
		t.Errorf("expected %q, got %v", expected, val)
	}
}

func TestDecodeValueCompound(t *testing.T) {
	val, err := DecodeValue([]byte{1, 2, 3}, DataTypeCompound)
	if err != nil {
		t.Fatal(err)
	}
	if val != nil {
		t.Errorf("expected nil for compound, got %v", val)
	}
}

func TestDecodeValueUnknownType(t *testing.T) {
	// Unknown data types should gracefully fall back to hex encoding.
	val, err := DecodeValue([]byte{0xDE, 0xAD, 0xBE, 0xEF}, DataType("foobar"))
	if err != nil {
		t.Fatalf("unexpected error for unknown type: %v", err)
	}
	expected := "DEADBEEF"
	if val != expected {
		t.Errorf("expected %q, got %v", expected, val)
	}
}

func TestDecodeValueEmptyType(t *testing.T) {
	// Empty data type (from missing schema metadata) should fall back to hex.
	val, err := DecodeValue([]byte{0x03, 0x0A, 0xFF}, DataType(""))
	if err != nil {
		t.Fatalf("unexpected error for empty type: %v", err)
	}
	expected := "030AFF"
	if val != expected {
		t.Errorf("expected %q, got %v", expected, val)
	}
}
