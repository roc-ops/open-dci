package main

import (
	"encoding/hex"
	"fmt"
)

// Encode converts a decoded DOCSIS config (as produced by Decode) back into the
// binary TLV format. It uses TLVOrder / _tlvOrder metadata to preserve the
// original TLV ordering. MIC TLVs (6, 7) are omitted by default.
//
// The result does NOT include an end-of-data marker (TLV 255); the caller
// should append 0xFF 0x00 if needed.
func Encode(result *DecodeResult, reg *Registry) ([]byte, error) {
	var out []byte

	order := result.TLVOrder
	if len(order) == 0 {
		// Fallback: iterate map keys in arbitrary order.
		for k := range result.Config {
			if k == "_tlvOrder" || k == "UnknownTlvs" {
				continue
			}
			order = append(order, k)
		}
	}

	for _, name := range order {
		// Skip MIC TLVs — they are not round-trippable without recomputation.
		if name == "CmMic" || name == "CmtsMic" {
			continue
		}

		val, ok := result.Config[name]
		if !ok {
			continue
		}

		def := reg.TopLevelByName(name)
		if def == nil {
			return nil, fmt.Errorf("unknown top-level property: %q", name)
		}

		// Handle chunked TLVs (e.g. CVC certificates): single value split into
		// consecutive ≤254-byte TLV instances on the wire.
		if def.Chunked {
			ls := defLengthSize(def)
			maxChunk := 254
			if ls == 2 {
				maxChunk = 65535
			}
			valueBytes, err := EncodeValue(val, def.DataType)
			if err != nil {
				return nil, fmt.Errorf("encoding chunked %s: %w", name, err)
			}
			for len(valueBytes) > 0 {
				chunkSize := len(valueBytes)
				if chunkSize > maxChunk {
					chunkSize = maxChunk
				}
				out = append(out, makeTLVn(def.TypeNum, valueBytes[:chunkSize], ls)...)
				valueBytes = valueBytes[chunkSize:]
			}
			continue
		}

		// Handle repeatable TLVs (always stored as []interface{}).
		if def.Repeatable {
			arr, ok := val.([]interface{})
			if !ok {
				return nil, fmt.Errorf("expected array for repeatable TLV %q, got %T", name, val)
			}
			for i, elem := range arr {
				b, err := encodeSingleTLV(def, elem, reg)
				if err != nil {
					return nil, fmt.Errorf("encoding %s[%d]: %w", name, i, err)
				}
				out = append(out, b...)
			}
		} else {
			b, err := encodeSingleTLV(def, val, reg)
			if err != nil {
				return nil, fmt.Errorf("encoding %s: %w", name, err)
			}
			out = append(out, b...)
		}
	}

	// Encode UnknownTlvs if present.
	if unknowns, ok := result.Config["UnknownTlvs"]; ok {
		arr, ok := unknowns.([]interface{})
		if ok {
			for _, u := range arr {
				m, ok := u.(map[string]interface{})
				if !ok {
					continue
				}
				b, err := encodeUnknownTLV(m)
				if err != nil {
					return nil, fmt.Errorf("encoding UnknownTlvs: %w", err)
				}
				out = append(out, b...)
			}
		}
	}

	// Append end-of-data marker.
	out = append(out, 0xFF, 0x00)

	return out, nil
}

// PadToAlignment pads the data with zero bytes to the next n-byte boundary.
// If already aligned, no padding is added. Typically called with n=4 after Encode.
func PadToAlignment(data []byte, n int) []byte {
	if n <= 1 || len(data)%n == 0 {
		return data
	}
	pad := n - (len(data) % n)
	return append(data, make([]byte, pad)...)
}

// encodeSingleTLV encodes a single TLV (simple, compound, TLV 11, or TLV 43).
func encodeSingleTLV(def *TLVDef, val interface{}, reg *Registry) ([]byte, error) {
	ls := defLengthSize(def)

	// Special case: TLV 10 (SNMP Write Access Control)
	if def.TypeNum == 10 {
		m, ok := val.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("TLV 10: expected map, got %T", val)
		}
		payload, err := EncodeSnmpWriteAccess(m)
		if err != nil {
			return nil, err
		}
		return makeTLVn(def.TypeNum, payload, ls), nil
	}

	// Special case: TLV 11 (SNMP MIB Object)
	if def.TypeNum == 11 {
		m, ok := val.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("TLV 11: expected map, got %T", val)
		}
		varbind, err := EncodeSnmpVarbind(m)
		if err != nil {
			return nil, err
		}
		return makeTLVn(def.TypeNum, varbind, ls), nil
	}

	// Special case: TLV 43 (DOCSIS Extension Field)
	if def.TypeNum == 43 {
		m, ok := val.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("TLV 43: expected map, got %T", val)
		}
		body, err := encodeTLV43(m, def)
		if err != nil {
			return nil, err
		}
		return makeTLVn(def.TypeNum, body, ls), nil
	}

	// Compound TLV
	if def.DataType == DataTypeCompound {
		m, ok := val.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("compound TLV %d: expected map, got %T", def.TypeNum, val)
		}
		body, err := encodeCompound(m, def)
		if err != nil {
			return nil, err
		}
		return makeTLVn(def.TypeNum, body, ls), nil
	}

	// Simple TLV
	valueBytes, err := EncodeValue(val, def.DataType)
	if err != nil {
		return nil, err
	}
	return makeTLVn(def.TypeNum, valueBytes, ls), nil
}

// encodeCompound encodes a compound TLV's sub-TLVs into a byte slice.
func encodeCompound(m map[string]interface{}, parent *TLVDef) ([]byte, error) {
	var out []byte

	// Use _tlvOrder if available for deterministic ordering.
	order := getMapOrder(m, parent)

	for _, name := range order {
		if name == "_tlvOrder" || name == "UnknownSubTlvs" {
			continue
		}

		val, ok := m[name]
		if !ok {
			continue
		}

		subDef := SubTLVByName(parent, name)
		if subDef == nil {
			return nil, fmt.Errorf("unknown sub-TLV property %q in compound TLV %d", name, parent.TypeNum)
		}

		if subDef.Repeatable {
			arr, ok := val.([]interface{})
			if !ok {
				return nil, fmt.Errorf("expected array for repeatable sub-TLV %q, got %T", name, val)
			}
			for _, elem := range arr {
				b, err := encodeSingleSubTLV(subDef, elem)
				if err != nil {
					return nil, err
				}
				out = append(out, b...)
			}
		} else {
			b, err := encodeSingleSubTLV(subDef, val)
			if err != nil {
				return nil, err
			}
			out = append(out, b...)
		}
	}

	// Encode UnknownSubTlvs if present.
	if unknowns, ok := m["UnknownSubTlvs"]; ok {
		arr, ok := unknowns.([]interface{})
		if ok {
			for _, u := range arr {
				um, ok := u.(map[string]interface{})
				if !ok {
					continue
				}
				b, err := encodeUnknownTLV(um)
				if err != nil {
					return nil, err
				}
				out = append(out, b...)
			}
		}
	}

	return out, nil
}

// encodeSingleSubTLV encodes a single sub-TLV.
func encodeSingleSubTLV(subDef *TLVDef, val interface{}) ([]byte, error) {
	ls := defLengthSize(subDef)

	if subDef.DataType == DataTypeCompound {
		m, ok := val.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("nested compound sub-TLV %d: expected map, got %T", subDef.TypeNum, val)
		}
		body, err := encodeCompound(m, subDef)
		if err != nil {
			return nil, err
		}
		return makeTLVn(subDef.TypeNum, body, ls), nil
	}

	valueBytes, err := EncodeValue(val, subDef.DataType)
	if err != nil {
		return nil, err
	}
	return makeTLVn(subDef.TypeNum, valueBytes, ls), nil
}

// encodeTLV43 encodes a TLV 43 (DOCSIS Extension Field) value.
// VendorId (sub-TLV 8) is always emitted first.
func encodeTLV43(m map[string]interface{}, def *TLVDef) ([]byte, error) {
	var out []byte

	// Emit VendorId first (sub-TLV 8).
	vendorIdStr, ok := m["VendorId"].(string)
	if !ok {
		return nil, fmt.Errorf("TLV 43: missing VendorId")
	}
	vendorIdBytes, err := hex.DecodeString(vendorIdStr)
	if err != nil {
		return nil, fmt.Errorf("TLV 43: invalid VendorId hex: %w", err)
	}
	out = append(out, makeTLV(8, vendorIdBytes)...)

	if vendorIdStr == "FFFFFF" {
		// General Extension: encode known sub-TLVs via registry.
		order := getMapOrder(m, def)
		for _, name := range order {
			if name == "_tlvOrder" || name == "VendorId" || name == "UnknownSubTlvs" {
				continue
			}

			val, ok := m[name]
			if !ok {
				continue
			}

			subDef := SubTLVByName(def, name)
			if subDef == nil {
				return nil, fmt.Errorf("TLV 43: unknown sub-TLV property %q", name)
			}

			if subDef.Repeatable {
				arr, ok := val.([]interface{})
				if !ok {
					return nil, fmt.Errorf("expected array for repeatable TLV 43 sub-TLV %q", name)
				}
				for _, elem := range arr {
					b, err := encodeSingleSubTLV(subDef, elem)
					if err != nil {
						return nil, err
					}
					out = append(out, b...)
				}
			} else {
				b, err := encodeSingleSubTLV(subDef, val)
				if err != nil {
					return nil, err
				}
				out = append(out, b...)
			}
		}

		// Unknown sub-TLVs in general extension.
		if unknowns, ok := m["UnknownSubTlvs"]; ok {
			arr, ok := unknowns.([]interface{})
			if ok {
				for _, u := range arr {
					um, ok := u.(map[string]interface{})
					if !ok {
						continue
					}
					b, err := encodeUnknownTLV(um)
					if err != nil {
						return nil, err
					}
					out = append(out, b...)
				}
			}
		}
	} else {
		// Vendor-specific: encode VendorSubTlvs as raw type/hex pairs.
		if vendorSubTlvs, ok := m["VendorSubTlvs"]; ok {
			arr, ok := vendorSubTlvs.([]interface{})
			if !ok {
				return nil, fmt.Errorf("TLV 43: VendorSubTlvs expected array, got %T", vendorSubTlvs)
			}
			for _, entry := range arr {
				em, ok := entry.(map[string]interface{})
				if !ok {
					continue
				}
				b, err := encodeUnknownTLV(em)
				if err != nil {
					return nil, err
				}
				out = append(out, b...)
			}
		}
	}

	return out, nil
}

// getMapOrder returns the _tlvOrder from a compound map, or falls back to
// iterating the map keys that correspond to known sub-TLV names.
func getMapOrder(m map[string]interface{}, parent *TLVDef) []string {
	if order, ok := m["_tlvOrder"].([]string); ok && len(order) > 0 {
		return order
	}
	// Fallback: use map keys.
	var keys []string
	for k := range m {
		if k == "_tlvOrder" || k == "UnknownSubTlvs" {
			continue
		}
		keys = append(keys, k)
	}
	return keys
}

// encodeUnknownTLV encodes a {"type": N, "value": "HEX"} map into a TLV.
func encodeUnknownTLV(m map[string]interface{}) ([]byte, error) {
	typeNum, err := toInt(m["type"])
	if err != nil {
		return nil, fmt.Errorf("unknown TLV: invalid type: %w", err)
	}
	hexStr, ok := m["value"].(string)
	if !ok {
		return nil, fmt.Errorf("unknown TLV type %d: expected hex string value", typeNum)
	}
	valueBytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, fmt.Errorf("unknown TLV type %d: invalid hex: %w", typeNum, err)
	}
	return makeTLV(typeNum, valueBytes), nil
}

// defLengthSize returns the effective length field size for a TLVDef.
// Returns 1 if the def is nil or LengthSize is not explicitly set.
func defLengthSize(def *TLVDef) int {
	if def != nil && def.LengthSize == 2 {
		return 2
	}
	return 1
}

// makeTLV builds a binary TLV: [type_byte][length_byte][value_bytes...].
// Uses a 1-byte length field (standard).
func makeTLV(typeNum int, value []byte) []byte {
	return makeTLVn(typeNum, value, 1)
}

// makeTLVn builds a binary TLV with a configurable length field size.
// lengthSize=1: [type][1-byte length][value] (standard)
// lengthSize=2: [type][2-byte big-endian length][value]
func makeTLVn(typeNum int, value []byte, lengthSize int) []byte {
	var result []byte
	if lengthSize == 2 {
		l := len(value)
		result = []byte{byte(typeNum), byte(l >> 8), byte(l & 0xFF)}
	} else {
		result = []byte{byte(typeNum), byte(len(value))}
	}
	result = append(result, value...)
	return result
}

// stripTLVsForComparison removes TLVs of specified types from a binary TLV stream,
// and also removes pad bytes (TLV 0). Used for round-trip comparison.
func stripTLVsForComparison(data []byte, excludeTypes ...int) []byte {
	excludeSet := make(map[int]bool)
	for _, t := range excludeTypes {
		excludeSet[t] = true
	}
	// Always exclude pad bytes.
	excludeSet[0] = true

	var result []byte
	offset := 0
	for offset < len(data) {
		if offset+1 >= len(data) {
			break
		}

		tlvType := int(data[offset])

		// TLV 0 = pad byte (single byte, no length).
		if tlvType == 0 {
			offset++
			continue
		}

		// TLV 255 = end-of-data marker (type + length byte 0x00).
		if tlvType == 255 {
			if !excludeSet[255] {
				result = append(result, data[offset])
				if offset+1 < len(data) {
					result = append(result, data[offset+1])
				}
			}
			break
		}

		tlvLen := int(data[offset+1])
		end := offset + 2 + tlvLen
		if end > len(data) {
			break
		}

		if !excludeSet[tlvType] {
			result = append(result, data[offset:end]...)
		}

		offset = end
	}

	return result
}

