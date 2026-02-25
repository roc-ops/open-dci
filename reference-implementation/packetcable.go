package main

import (
	"bytes"
	"crypto/sha1"
	"fmt"
)

// PacketCable hash OID variants.
const (
	PacketCableNA   = "na"
	PacketCableEU   = "eu"
	PacketCableIETF = "ietf"
)

// packetCableHashOIDs maps variant names to their OID strings.
var packetCableHashOIDs = map[string]string{
	PacketCableNA:   "1.3.6.1.4.1.4491.2.2.1.1.2.7.0",
	PacketCableEU:   "1.3.6.1.4.1.7432.1.1.2.9.0",
	PacketCableIETF: "1.3.6.1.2.1.140.1.2.11.0",
}

// buildHashVarbind constructs a TLV 11 SNMP varbind containing the given OID
// and a 20-byte zero placeholder for the SHA-1 hash value. It returns the
// complete TLV 11 bytes and the offset within those bytes where the 20-byte
// hash value starts, so the caller can overwrite it with the computed hash.
func buildHashVarbind(oid string) (tlv11 []byte, hashValueOffset int, err error) {
	// Encode the OID.
	oidBytes, err := encodeOIDBytes(oid)
	if err != nil {
		return nil, 0, fmt.Errorf("encoding hash OID: %w", err)
	}
	oidTLV := encodeBERTLV(tagOID, oidBytes)

	// 20-byte zero placeholder for the SHA-1 hash.
	placeholder := make([]byte, 20)
	valueTLV := encodeBERTLV(tagOctetString, placeholder)

	// SEQUENCE wrapper: OID TLV + OctetString TLV.
	inner := append(oidTLV, valueTLV...)
	sequence := encodeBERTLV(0x30, inner)

	// Outer TLV 11 wrapper (1-byte length field).
	tlv11 = makeTLVn(11, sequence, 1)

	// The hash value (20 zero bytes) is at the end of the TLV 11 structure.
	hashValueOffset = len(tlv11) - 20

	return tlv11, hashValueOffset, nil
}

// insertPacketCableHash computes a circular SHA-1 hash for a PacketCable MTA
// config file and inserts it as a TLV 11 SNMP varbind before the end-of-data
// marker.
//
// The algorithm:
//  1. Build a TLV 11 varbind with the hash OID and 20 zero-byte placeholder.
//  2. Insert the varbind before the end-of-data marker (0xFF 0x00).
//  3. Compute SHA-1 over the entire assembled file (with zeroed placeholder).
//  4. Replace the 20-byte placeholder with the computed SHA-1 hash.
func insertPacketCableHash(encoded []byte, variant string) ([]byte, error) {
	oid, ok := packetCableHashOIDs[variant]
	if !ok {
		return nil, fmt.Errorf("unknown PacketCable hash variant: %q (use na, eu, or ietf)", variant)
	}

	// Build TLV 11 varbind with 20 zero-byte placeholder.
	varbind, varbindHashOffset, err := buildHashVarbind(oid)
	if err != nil {
		return nil, err
	}

	// Verify end-of-data marker.
	if len(encoded) < 2 || encoded[len(encoded)-2] != 0xFF || encoded[len(encoded)-1] != 0x00 {
		return nil, fmt.Errorf("encoded data does not end with end-of-data marker")
	}
	body := encoded[:len(encoded)-2]

	// Insert varbind before end-of-data marker.
	var withHash []byte
	withHash = append(withHash, body...)
	withHash = append(withHash, varbind...)
	withHash = append(withHash, 0xFF, 0x00)

	// Compute SHA-1 of the entire file with zeroed hash placeholder.
	h := sha1.Sum(withHash)

	// The hash value offset in the assembled output is:
	// len(body) bytes of original data + varbindHashOffset within the varbind.
	absoluteOffset := len(body) + varbindHashOffset

	// Verify the placeholder is where we expect it.
	placeholder := make([]byte, 20)
	if !bytes.Equal(withHash[absoluteOffset:absoluteOffset+20], placeholder) {
		return nil, fmt.Errorf("internal error: hash placeholder not found at expected offset")
	}

	// Replace placeholder with the computed SHA-1 hash.
	copy(withHash[absoluteOffset:absoluteOffset+20], h[:])

	return withHash, nil
}
