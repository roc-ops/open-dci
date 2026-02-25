package main

import (
	"fmt"
	"math/big"
	"net"
	"strings"
)

// SnmpVarbind represents a decoded SNMP MIB varbind.
type SnmpVarbind struct {
	OID   string `json:"oid"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

// ASN.1 tag constants for SNMP value types.
const (
	tagInteger     = 0x02
	tagOctetString = 0x04
	tagNull        = 0x05
	tagOID         = 0x06
	tagIPAddress   = 0x40
	tagCounter32   = 0x41
	tagGauge32     = 0x42 // Also Unsigned32
	tagTimeTicks   = 0x43
	tagOpaque      = 0x44
	tagCounter64   = 0x46
)

// DecodeSnmpVarbind parses a BER-encoded SNMP varbind and returns the
// decoded OID, type name, and value as strings.
func DecodeSnmpVarbind(data []byte) (*SnmpVarbind, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("varbind too short: %d bytes", len(data))
	}

	// Outer SEQUENCE (tag 0x30)
	if data[0] != 0x30 {
		return nil, fmt.Errorf("expected SEQUENCE tag 0x30, got 0x%02X", data[0])
	}

	seqLen, seqHdrLen, err := berLength(data[1:])
	if err != nil {
		return nil, fmt.Errorf("parsing SEQUENCE length: %w", err)
	}

	inner := data[1+seqHdrLen : 1+seqHdrLen+seqLen]

	// First element: OID (tag 0x06)
	if len(inner) < 2 {
		return nil, fmt.Errorf("inner SEQUENCE too short")
	}

	if inner[0] != tagOID {
		return nil, fmt.Errorf("expected OID tag 0x06, got 0x%02X", inner[0])
	}

	oidLen, oidHdrLen, err := berLength(inner[1:])
	if err != nil {
		return nil, fmt.Errorf("parsing OID length: %w", err)
	}

	oidBytes := inner[1+oidHdrLen : 1+oidHdrLen+oidLen]
	oid := decodeOIDBytes(oidBytes)

	// Second element: value
	valStart := 1 + oidHdrLen + oidLen
	if valStart >= len(inner) {
		return nil, fmt.Errorf("no value element in varbind")
	}

	valTag := inner[valStart]
	valLen, valHdrLen, err := berLength(inner[valStart+1:])
	if err != nil {
		return nil, fmt.Errorf("parsing value length: %w", err)
	}

	valBytes := inner[valStart+1+valHdrLen : valStart+1+valHdrLen+valLen]

	typeName, value, err := decodeSnmpValue(valTag, valBytes)
	if err != nil {
		return nil, fmt.Errorf("decoding value: %w", err)
	}

	return &SnmpVarbind{
		OID:   oid,
		Type:  typeName,
		Value: value,
	}, nil
}

// berLength parses a BER length field. Returns the length value and the
// number of bytes consumed by the length encoding.
func berLength(data []byte) (int, int, error) {
	if len(data) == 0 {
		return 0, 0, fmt.Errorf("empty length field")
	}

	if data[0] < 0x80 {
		// Short form: single byte
		return int(data[0]), 1, nil
	}

	if data[0] == 0x80 {
		return 0, 0, fmt.Errorf("indefinite length not supported")
	}

	// Long form: first byte indicates number of subsequent length bytes
	numBytes := int(data[0] & 0x7F)
	if numBytes > 4 {
		return 0, 0, fmt.Errorf("length field too long: %d bytes", numBytes)
	}
	if len(data) < 1+numBytes {
		return 0, 0, fmt.Errorf("truncated length field")
	}

	length := 0
	for i := 0; i < numBytes; i++ {
		length = (length << 8) | int(data[1+i])
	}

	return length, 1 + numBytes, nil
}

// decodeSnmpValue converts a BER-encoded value into a type name and string representation.
func decodeSnmpValue(tag byte, data []byte) (string, string, error) {
	switch tag {
	case tagInteger:
		return "Integer", decodeSignedInteger(data), nil

	case tagOctetString:
		if isPrintable(data) {
			return "String", string(data), nil
		}
		return "HexString", fmt.Sprintf("%X", data), nil

	case tagNull:
		return "Null", "", nil

	case tagOID:
		return "OID", decodeOIDBytes(data), nil

	case tagIPAddress:
		if len(data) != 4 {
			return "IPAddress", fmt.Sprintf("%X", data), nil
		}
		return "IPAddress", net.IP(data).String(), nil

	case tagCounter32:
		return "Counter32", decodeUnsignedInteger(data), nil

	case tagGauge32:
		return "Gauge32", decodeUnsignedInteger(data), nil

	case tagTimeTicks:
		return "TimeTicks", decodeUnsignedInteger(data), nil

	case tagOpaque:
		return "HexString", fmt.Sprintf("%X", data), nil

	case tagCounter64:
		return "Counter64", decodeUnsignedInteger(data), nil

	default:
		return fmt.Sprintf("Unknown(0x%02X)", tag), fmt.Sprintf("%X", data), nil
	}
}

// decodeSignedInteger decodes a BER-encoded signed integer.
func decodeSignedInteger(data []byte) string {
	if len(data) == 0 {
		return "0"
	}

	// BER integers are big-endian, sign-extended
	val := big.NewInt(0)
	if data[0]&0x80 != 0 {
		// Negative: sign extend
		val.SetInt64(-1)
	}
	for _, b := range data {
		val.Lsh(val, 8)
		val.Or(val, big.NewInt(int64(b)))
	}

	return val.String()
}

// decodeUnsignedInteger decodes a BER-encoded unsigned integer.
func decodeUnsignedInteger(data []byte) string {
	if len(data) == 0 {
		return "0"
	}

	val := new(big.Int)
	val.SetBytes(data)
	return val.String()
}

// isPrintable returns true if all bytes are printable ASCII.
func isPrintable(data []byte) bool {
	for _, b := range data {
		if b < 0x20 || b > 0x7E {
			return false
		}
	}
	return len(data) > 0
}

// formatOIDFromDotted formats an OID string (no conversion needed, used for consistency).
func formatOIDFromDotted(parts []int) string {
	strs := make([]string, len(parts))
	for i, p := range parts {
		strs[i] = fmt.Sprintf("%d", p)
	}
	return strings.Join(strs, ".")
}
