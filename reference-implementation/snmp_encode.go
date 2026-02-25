package main

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"net"
	"strconv"
	"strings"
)

// EncodeSnmpVarbind encodes an SNMP varbind map (with "oid", "type", "value" keys)
// into the BER-encoded binary format expected by TLV 11. This is the inverse of
// DecodeSnmpVarbind.
func EncodeSnmpVarbind(entry map[string]interface{}) ([]byte, error) {
	oidStr, ok := entry["oid"].(string)
	if !ok {
		return nil, fmt.Errorf("varbind missing 'oid' string")
	}
	typeName, ok := entry["type"].(string)
	if !ok {
		return nil, fmt.Errorf("varbind missing 'type' string")
	}
	valueStr, _ := entry["value"].(string)

	// Encode the OID.
	oidBytes, err := encodeOIDBytes(oidStr)
	if err != nil {
		return nil, fmt.Errorf("encoding OID: %w", err)
	}
	oidTLV := encodeBERTLV(tagOID, oidBytes)

	// Encode the value based on its type.
	valTag, valBytes, err := encodeSnmpValue(typeName, valueStr)
	if err != nil {
		return nil, fmt.Errorf("encoding value: %w", err)
	}
	valTLV := encodeBERTLV(valTag, valBytes)

	// Wrap in SEQUENCE.
	inner := append(oidTLV, valTLV...)
	return encodeBERTLV(0x30, inner), nil
}

// encodeSnmpValue encodes a typed SNMP value into a BER tag and value bytes.
func encodeSnmpValue(typeName string, value string) (byte, []byte, error) {
	switch typeName {
	case "Integer":
		b, err := encodeSignedIntegerBER(value)
		return tagInteger, b, err

	case "String":
		return tagOctetString, []byte(value), nil

	case "HexString":
		b, err := hex.DecodeString(value)
		if err != nil {
			return 0, nil, fmt.Errorf("HexString: %w", err)
		}
		// HexString decoded from Opaque (0x44) or OctetString (0x04).
		// We use OctetString (0x04) since that's the more common case and
		// the decoder maps both to "HexString".
		return tagOctetString, b, nil

	case "Null":
		return tagNull, nil, nil

	case "OID":
		oidBytes, err := encodeOIDBytes(value)
		if err != nil {
			return 0, nil, err
		}
		return tagOID, oidBytes, nil

	case "IPAddress":
		ip := net.ParseIP(value)
		if ip == nil {
			return 0, nil, fmt.Errorf("invalid IP address: %q", value)
		}
		ip4 := ip.To4()
		if ip4 == nil {
			return 0, nil, fmt.Errorf("not an IPv4 address: %q", value)
		}
		return tagIPAddress, ip4, nil

	case "Counter32":
		b, err := encodeUnsignedIntegerBER(value)
		return tagCounter32, b, err

	case "Gauge32":
		b, err := encodeUnsignedIntegerBER(value)
		return tagGauge32, b, err

	case "TimeTicks":
		b, err := encodeUnsignedIntegerBER(value)
		return tagTimeTicks, b, err

	case "Counter64":
		b, err := encodeUnsignedIntegerBER(value)
		return tagCounter64, b, err

	default:
		// Handle Unknown(0xNN) format.
		if strings.HasPrefix(typeName, "Unknown(0x") && strings.HasSuffix(typeName, ")") {
			tagHex := typeName[len("Unknown(0x") : len(typeName)-1]
			tagVal, err := strconv.ParseUint(tagHex, 16, 8)
			if err != nil {
				return 0, nil, fmt.Errorf("unknown type tag: %w", err)
			}
			b, err := hex.DecodeString(value)
			if err != nil {
				return 0, nil, fmt.Errorf("unknown type value: %w", err)
			}
			return byte(tagVal), b, nil
		}
		return 0, nil, fmt.Errorf("unsupported SNMP type: %q", typeName)
	}
}

// encodeSignedIntegerBER encodes a decimal string as a BER signed integer value
// using the minimal two's complement encoding.
func encodeSignedIntegerBER(s string) ([]byte, error) {
	val := new(big.Int)
	_, ok := val.SetString(s, 10)
	if !ok {
		return nil, fmt.Errorf("invalid integer: %q", s)
	}

	if val.Sign() == 0 {
		return []byte{0}, nil
	}

	if val.Sign() > 0 {
		b := val.Bytes()
		// If the high bit is set, prepend a zero byte to keep it positive.
		if b[0]&0x80 != 0 {
			b = append([]byte{0}, b...)
		}
		return b, nil
	}

	// Negative: determine how many bytes are needed for two's complement.
	// The bit length of the absolute value tells us the minimum bits needed.
	pos := new(big.Int).Neg(val)
	bitLen := pos.BitLen()

	// Number of bytes needed: ceil((bitLen+1) / 8) for sign bit,
	// but if pos is an exact power of 2, the value -pos fits in bitLen bits
	// (e.g., -128 fits in 1 byte as 0x80).
	numBytes := (bitLen + 8) / 8
	if pos.Cmp(new(big.Int).Lsh(big.NewInt(1), uint(bitLen-1))) == 0 {
		// Exact power of 2: -128 = 0x80 (1 byte), -32768 = 0x8000 (2 bytes)
		numBytes = (bitLen + 7) / 8
	}

	// Compute two's complement: 2^(numBytes*8) + val
	twoComp := new(big.Int).Lsh(big.NewInt(1), uint(numBytes*8))
	twoComp.Add(twoComp, val)
	tcBytes := twoComp.Bytes()

	// Pad to numBytes if needed.
	result := make([]byte, numBytes)
	for i := range result {
		result[i] = 0xFF
	}
	copy(result[numBytes-len(tcBytes):], tcBytes)
	return result, nil
}

// encodeUnsignedIntegerBER encodes a decimal string as BER unsigned integer value bytes.
func encodeUnsignedIntegerBER(s string) ([]byte, error) {
	val := new(big.Int)
	_, ok := val.SetString(s, 10)
	if !ok {
		return nil, fmt.Errorf("invalid unsigned integer: %q", s)
	}
	b := val.Bytes()
	if len(b) == 0 {
		return []byte{0}, nil
	}
	return b, nil
}

// EncodeSnmpWriteAccess encodes a TLV 10 (SNMP Write Access Control) entry.
// The entry map must have "oid" (string) and "access" (int) keys.
// Returns the BER-encoded OID followed by a single access flag byte.
func EncodeSnmpWriteAccess(entry map[string]interface{}) ([]byte, error) {
	oidStr, ok := entry["oid"].(string)
	if !ok {
		return nil, fmt.Errorf("write-access entry missing 'oid' string")
	}
	accessVal, err := toInt(entry["access"])
	if err != nil {
		return nil, fmt.Errorf("write-access entry missing 'access' integer: %w", err)
	}

	oidBytes, err := encodeOIDBytes(oidStr)
	if err != nil {
		return nil, fmt.Errorf("encoding OID: %w", err)
	}
	oidTLV := encodeBERTLV(tagOID, oidBytes)

	// Concatenate BER-encoded OID + single access byte (no SEQUENCE wrapper).
	result := append(oidTLV, byte(accessVal))
	return result, nil
}

// encodeBERLength encodes a length value in BER format.
func encodeBERLength(length int) []byte {
	if length < 0x80 {
		return []byte{byte(length)}
	}
	if length <= 0xFF {
		return []byte{0x81, byte(length)}
	}
	if length <= 0xFFFF {
		return []byte{0x82, byte(length >> 8), byte(length)}
	}
	// 3-byte length (should be sufficient for DOCSIS config TLVs).
	return []byte{0x83, byte(length >> 16), byte(length >> 8), byte(length)}
}

// encodeBERTLV wraps a value in a BER tag-length-value structure.
func encodeBERTLV(tag byte, value []byte) []byte {
	result := []byte{tag}
	result = append(result, encodeBERLength(len(value))...)
	result = append(result, value...)
	return result
}
