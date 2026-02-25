package main

import (
	"bytes"
	"crypto/sha1"
	"testing"
)

func TestBuildHashVarbind_NA(t *testing.T) {
	tlv11, hashOffset, err := buildHashVarbind(packetCableHashOIDs[PacketCableNA])
	if err != nil {
		t.Fatalf("buildHashVarbind(NA) error: %v", err)
	}

	// TLV 11 outer: type=11, then 1-byte length.
	if tlv11[0] != 11 {
		t.Errorf("expected TLV type 11, got %d", tlv11[0])
	}

	// SEQUENCE tag should follow the TLV 11 header (type + length = 2 bytes).
	if tlv11[2] != 0x30 {
		t.Errorf("expected SEQUENCE tag 0x30 at offset 2, got 0x%02X", tlv11[2])
	}

	// The hash value (20 zero bytes) should be at the reported offset.
	hashValue := tlv11[hashOffset : hashOffset+20]
	if !bytes.Equal(hashValue, make([]byte, 20)) {
		t.Errorf("expected 20 zero bytes at hash offset %d, got %X", hashOffset, hashValue)
	}

	// Verify the OID is present by encoding it and checking it appears in the varbind.
	oidBytes, err := encodeOIDBytes(packetCableHashOIDs[PacketCableNA])
	if err != nil {
		t.Fatalf("encodeOIDBytes error: %v", err)
	}
	if !bytes.Contains(tlv11, oidBytes) {
		t.Error("TLV 11 does not contain the expected NA OID bytes")
	}
}

func TestBuildHashVarbind_EU(t *testing.T) {
	tlv11, hashOffset, err := buildHashVarbind(packetCableHashOIDs[PacketCableEU])
	if err != nil {
		t.Fatalf("buildHashVarbind(EU) error: %v", err)
	}

	if tlv11[0] != 11 {
		t.Errorf("expected TLV type 11, got %d", tlv11[0])
	}

	hashValue := tlv11[hashOffset : hashOffset+20]
	if !bytes.Equal(hashValue, make([]byte, 20)) {
		t.Errorf("expected 20 zero bytes at hash offset %d, got %X", hashOffset, hashValue)
	}

	oidBytes, err := encodeOIDBytes(packetCableHashOIDs[PacketCableEU])
	if err != nil {
		t.Fatalf("encodeOIDBytes error: %v", err)
	}
	if !bytes.Contains(tlv11, oidBytes) {
		t.Error("TLV 11 does not contain the expected EU OID bytes")
	}
}

func TestBuildHashVarbind_IETF(t *testing.T) {
	tlv11, hashOffset, err := buildHashVarbind(packetCableHashOIDs[PacketCableIETF])
	if err != nil {
		t.Fatalf("buildHashVarbind(IETF) error: %v", err)
	}

	if tlv11[0] != 11 {
		t.Errorf("expected TLV type 11, got %d", tlv11[0])
	}

	hashValue := tlv11[hashOffset : hashOffset+20]
	if !bytes.Equal(hashValue, make([]byte, 20)) {
		t.Errorf("expected 20 zero bytes at hash offset %d, got %X", hashOffset, hashValue)
	}

	oidBytes, err := encodeOIDBytes(packetCableHashOIDs[PacketCableIETF])
	if err != nil {
		t.Fatalf("encodeOIDBytes error: %v", err)
	}
	if !bytes.Contains(tlv11, oidBytes) {
		t.Error("TLV 11 does not contain the expected IETF OID bytes")
	}
}

func TestInsertPacketCableHash(t *testing.T) {
	// Build a minimal config: a single TLV (type=1, length=1, value=0x01) + end-of-data.
	minimalConfig := []byte{
		0x01, 0x01, 0x01, // TLV type 1, length 1, value 0x01
		0xFF, 0x00, // end-of-data marker
	}

	result, err := insertPacketCableHash(minimalConfig, PacketCableNA)
	if err != nil {
		t.Fatalf("insertPacketCableHash error: %v", err)
	}

	// Result should be longer than input (has the TLV 11 varbind inserted).
	if len(result) <= len(minimalConfig) {
		t.Errorf("result length %d should be greater than input length %d", len(result), len(minimalConfig))
	}

	// Result should still end with end-of-data marker.
	if result[len(result)-2] != 0xFF || result[len(result)-1] != 0x00 {
		t.Error("result does not end with end-of-data marker (0xFF 0x00)")
	}

	// The 20 bytes at the hash position should NOT all be zeros (hash was computed).
	// Find the hash position by building the varbind to get the offset.
	varbind, varbindHashOffset, err := buildHashVarbind(packetCableHashOIDs[PacketCableNA])
	if err != nil {
		t.Fatalf("buildHashVarbind error: %v", err)
	}
	bodyLen := len(minimalConfig) - 2 // original body without end-of-data
	absoluteOffset := bodyLen + varbindHashOffset
	hashValue := result[absoluteOffset : absoluteOffset+20]
	if bytes.Equal(hashValue, make([]byte, 20)) {
		t.Error("hash value should not be all zeros after insertion")
	}

	// Verify TLV 11 type byte is present at the expected position.
	varbindStart := bodyLen
	if result[varbindStart] != 11 {
		t.Errorf("expected TLV type 11 at offset %d, got %d", varbindStart, result[varbindStart])
	}

	// Verify the varbind length is consistent.
	expectedLen := len(varbind)
	actualVarbindEnd := varbindStart + expectedLen
	if actualVarbindEnd+2 != len(result) {
		t.Errorf("varbind end (%d) + end-of-data (2) should equal result length (%d)",
			actualVarbindEnd, len(result))
	}
}

func TestInsertPacketCableHash_CircularVerification(t *testing.T) {
	// This is the key test: verify the circular hash property.
	// After inserting the hash, zeroing out the hash value and recomputing SHA-1
	// should yield the same hash that was inserted.

	for _, variant := range []string{PacketCableNA, PacketCableEU, PacketCableIETF} {
		t.Run(variant, func(t *testing.T) {
			// Build a config with some content.
			config := []byte{
				0x01, 0x01, 0x01, // TLV type 1, length 1, value 0x01
				0x02, 0x03, 0x0A, 0x0B, 0x0C, // TLV type 2, length 3, value 0x0A0B0C
				0xFF, 0x00, // end-of-data marker
			}

			result, err := insertPacketCableHash(config, variant)
			if err != nil {
				t.Fatalf("insertPacketCableHash error: %v", err)
			}

			// Find the hash value offset in the result.
			varbind, varbindHashOffset, err := buildHashVarbind(packetCableHashOIDs[variant])
			if err != nil {
				t.Fatalf("buildHashVarbind error: %v", err)
			}
			_ = varbind
			bodyLen := len(config) - 2
			absoluteOffset := bodyLen + varbindHashOffset

			// Extract the inserted hash.
			insertedHash := make([]byte, 20)
			copy(insertedHash, result[absoluteOffset:absoluteOffset+20])

			// Zero out the hash value in a copy of the result.
			zeroed := make([]byte, len(result))
			copy(zeroed, result)
			copy(zeroed[absoluteOffset:absoluteOffset+20], make([]byte, 20))

			// Compute SHA-1 of the zeroed result.
			recomputedHash := sha1.Sum(zeroed)

			// The recomputed hash should match the inserted hash.
			if !bytes.Equal(insertedHash, recomputedHash[:]) {
				t.Errorf("circular hash verification failed:\n  inserted:    %X\n  recomputed:  %X",
					insertedHash, recomputedHash)
			}
		})
	}
}

func TestInsertPacketCableHash_InvalidVariant(t *testing.T) {
	config := []byte{0x01, 0x01, 0x01, 0xFF, 0x00}
	_, err := insertPacketCableHash(config, "invalid")
	if err == nil {
		t.Fatal("expected error for invalid variant, got nil")
	}
}

func TestInsertPacketCableHash_NoEndMarker(t *testing.T) {
	// Config without end-of-data marker.
	config := []byte{0x01, 0x01, 0x01}
	_, err := insertPacketCableHash(config, PacketCableNA)
	if err == nil {
		t.Fatal("expected error for missing end-of-data marker, got nil")
	}
}
