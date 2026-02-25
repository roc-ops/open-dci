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

// encodeSingleTLV encodes a single TLV (simple, compound, TLV 11, or TLV 43).
func encodeSingleTLV(def *TLVDef, val interface{}, reg *Registry) ([]byte, error) {
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
		return makeTLV(def.TypeNum, varbind), nil
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
		return makeTLV(def.TypeNum, body), nil
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
		return makeTLV(def.TypeNum, body), nil
	}

	// Simple TLV
	valueBytes, err := EncodeValue(val, def.DataType)
	if err != nil {
		return nil, err
	}
	return makeTLV(def.TypeNum, valueBytes), nil
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
	if subDef.DataType == DataTypeCompound {
		m, ok := val.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("nested compound sub-TLV %d: expected map, got %T", subDef.TypeNum, val)
		}
		body, err := encodeCompound(m, subDef)
		if err != nil {
			return nil, err
		}
		return makeTLV(subDef.TypeNum, body), nil
	}

	valueBytes, err := EncodeValue(val, subDef.DataType)
	if err != nil {
		return nil, err
	}
	return makeTLV(subDef.TypeNum, valueBytes), nil
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

// makeTLV builds a binary TLV: [type_byte][length_byte][value_bytes...].
func makeTLV(typeNum int, value []byte) []byte {
	result := []byte{byte(typeNum), byte(len(value))}
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

