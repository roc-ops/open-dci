package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

// DataType represents the wire-format data type for a TLV value.
type DataType string

const (
	DataTypeUint8       DataType = "uint8"
	DataTypeUint16      DataType = "uint16"
	DataTypeUint32      DataType = "uint32"
	DataTypeString      DataType = "string"
	DataTypeHexString   DataType = "hexstring"
	DataTypeMacAddress  DataType = "macAddress"
	DataTypeIPv4Address DataType = "ipv4Address"
	DataTypeIPv6Address DataType = "ipv6Address"
	DataTypeOID         DataType = "oid"
	DataTypeCompound    DataType = "compound"
)

// DecodeValue converts raw TLV value bytes into a Go value suitable for JSON
// serialization, based on the specified data type.
func DecodeValue(data []byte, dt DataType) (interface{}, error) {
	switch dt {
	case DataTypeUint8:
		if len(data) != 1 {
			return nil, fmt.Errorf("uint8 requires 1 byte, got %d", len(data))
		}
		return int(data[0]), nil

	case DataTypeUint16:
		if len(data) != 2 {
			return nil, fmt.Errorf("uint16 requires 2 bytes, got %d", len(data))
		}
		return int(binary.BigEndian.Uint16(data)), nil

	case DataTypeUint32:
		if len(data) != 4 {
			return nil, fmt.Errorf("uint32 requires 4 bytes, got %d", len(data))
		}
		return int(binary.BigEndian.Uint32(data)), nil

	case DataTypeString:
		// Strip null terminator if present
		s := string(data)
		s = strings.TrimRight(s, "\x00")
		return s, nil

	case DataTypeHexString:
		return fmt.Sprintf("%X", data), nil

	case DataTypeMacAddress:
		if len(data) != 6 {
			return nil, fmt.Errorf("macAddress requires 6 bytes, got %d", len(data))
		}
		return fmt.Sprintf("%02X%02X%02X%02X%02X%02X", data[0], data[1], data[2], data[3], data[4], data[5]), nil

	case DataTypeIPv4Address:
		if len(data) != 4 {
			return nil, fmt.Errorf("ipv4Address requires 4 bytes, got %d", len(data))
		}
		return net.IP(data).String(), nil

	case DataTypeIPv6Address:
		if len(data) != 16 {
			return nil, fmt.Errorf("ipv6Address requires 16 bytes, got %d", len(data))
		}
		return net.IP(data).String(), nil

	case DataTypeOID:
		return decodeOIDBytes(data), nil

	case DataTypeCompound:
		// Signal that the caller should recurse into sub-TLVs.
		return nil, nil

	default:
		return nil, fmt.Errorf("unknown data type: %s", dt)
	}
}

// decodeOIDBytes decodes a BER-encoded OID value (without tag/length) into
// dotted-decimal notation.
func decodeOIDBytes(data []byte) string {
	if len(data) == 0 {
		return ""
	}

	parts := make([]string, 0, len(data))

	// First byte encodes first two components: X*40 + Y
	first := int(data[0])
	parts = append(parts, fmt.Sprintf("%d", first/40))
	parts = append(parts, fmt.Sprintf("%d", first%40))

	// Remaining bytes use base-128 encoding with high bit as continuation flag
	val := 0
	for i := 1; i < len(data); i++ {
		val = (val << 7) | int(data[i]&0x7F)
		if data[i]&0x80 == 0 {
			parts = append(parts, fmt.Sprintf("%d", val))
			val = 0
		}
	}

	return strings.Join(parts, ".")
}
