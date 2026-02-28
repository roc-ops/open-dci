package opendci

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// EncodeValue converts a decoded JSON value back into raw TLV value bytes,
// based on the specified data type. This is the inverse of DecodeValue.
func EncodeValue(val interface{}, dt DataType) ([]byte, error) {
	switch dt {
	case DataTypeUint8:
		n, err := toInt(val)
		if err != nil {
			return nil, fmt.Errorf("uint8: %w", err)
		}
		if n < 0 || n > 255 {
			return nil, fmt.Errorf("uint8 value out of range: %d", n)
		}
		return []byte{byte(n)}, nil

	case DataTypeUint16:
		n, err := toInt(val)
		if err != nil {
			return nil, fmt.Errorf("uint16: %w", err)
		}
		if n < 0 || n > 65535 {
			return nil, fmt.Errorf("uint16 value out of range: %d", n)
		}
		buf := make([]byte, 2)
		binary.BigEndian.PutUint16(buf, uint16(n))
		return buf, nil

	case DataTypeUint32:
		n, err := toInt(val)
		if err != nil {
			return nil, fmt.Errorf("uint32: %w", err)
		}
		if n < 0 || n > 0xFFFFFFFF {
			return nil, fmt.Errorf("uint32 value out of range: %d", n)
		}
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(n))
		return buf, nil

	case DataTypeString:
		s, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf("string: expected string, got %T", val)
		}
		// Always null-terminate per DOCSIS spec.
		return append([]byte(s), 0x00), nil

	case DataTypeHexString:
		s, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf("hexstring: expected string, got %T", val)
		}
		b, err := hex.DecodeString(s)
		if err != nil {
			return nil, fmt.Errorf("hexstring: %w", err)
		}
		return b, nil

	case DataTypeMacAddress:
		s, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf("macAddress: expected string, got %T", val)
		}
		if len(s) != 12 {
			return nil, fmt.Errorf("macAddress: expected 12 hex chars, got %d", len(s))
		}
		b, err := hex.DecodeString(s)
		if err != nil {
			return nil, fmt.Errorf("macAddress: %w", err)
		}
		return b, nil

	case DataTypeIPv4Address:
		s, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf("ipv4Address: expected string, got %T", val)
		}
		ip := net.ParseIP(s)
		if ip == nil {
			return nil, fmt.Errorf("ipv4Address: invalid IP %q", s)
		}
		ip4 := ip.To4()
		if ip4 == nil {
			return nil, fmt.Errorf("ipv4Address: not an IPv4 address %q", s)
		}
		return ip4, nil

	case DataTypeIPv6Address:
		s, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf("ipv6Address: expected string, got %T", val)
		}
		ip := net.ParseIP(s)
		if ip == nil {
			return nil, fmt.Errorf("ipv6Address: invalid IP %q", s)
		}
		ip6 := ip.To16()
		if ip6 == nil {
			return nil, fmt.Errorf("ipv6Address: not an IPv6 address %q", s)
		}
		return ip6, nil

	case DataTypeOID:
		s, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf("oid: expected string, got %T", val)
		}
		return encodeOIDBytes(s)

	case DataTypeCompound:
		// Compound types are handled at a higher level; value bytes are the
		// concatenation of encoded sub-TLVs.
		return nil, nil

	default:
		return nil, fmt.Errorf("unknown data type: %s", dt)
	}
}

// toInt converts a JSON numeric value (float64 from JSON unmarshal, or int from
// the decoder) to an integer.
func toInt(val interface{}) (int, error) {
	switch v := val.(type) {
	case int:
		return v, nil
	case float64:
		return int(v), nil
	case json.Number:
		n, err := v.Int64()
		if err != nil {
			return 0, err
		}
		return int(n), nil
	default:
		return 0, fmt.Errorf("expected number, got %T", val)
	}
}

// encodeOIDBytes encodes a dotted-decimal OID string into BER OID value bytes
// (without tag or length). This is the inverse of decodeOIDBytes.
func encodeOIDBytes(oid string) ([]byte, error) {
	parts := strings.Split(oid, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("OID must have at least 2 components: %q", oid)
	}

	nums := make([]int, len(parts))
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("invalid OID component %q: %w", p, err)
		}
		nums[i] = n
	}

	var buf []byte

	// First two components encoded as X*40 + Y.
	buf = append(buf, byte(nums[0]*40+nums[1]))

	// Remaining components use base-128 encoding with high bit as continuation.
	for i := 2; i < len(nums); i++ {
		buf = append(buf, encodeBase128(nums[i])...)
	}

	return buf, nil
}

// encodeBase128 encodes an integer in base-128 with high-bit continuation.
func encodeBase128(val int) []byte {
	if val == 0 {
		return []byte{0}
	}

	// Collect 7-bit groups in reverse order.
	var parts []byte
	for v := val; v > 0; v >>= 7 {
		parts = append(parts, byte(v&0x7F))
	}

	// Reverse and set continuation bits.
	result := make([]byte, len(parts))
	for i, b := range parts {
		idx := len(parts) - 1 - i
		if idx < len(result)-1 {
			result[idx] = b | 0x80
		} else {
			result[idx] = b
		}
	}

	return result
}
