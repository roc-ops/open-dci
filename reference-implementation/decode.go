package main

import (
	"fmt"
)

// UnknownTLV represents a TLV that is not in the registry.
type UnknownTLV struct {
	Type  int    `json:"type"`
	Value string `json:"value"`
}

// DecodeResult holds the decoded config and raw data for MIC verification.
type DecodeResult struct {
	Config   map[string]interface{}
	RawData  []byte   // Original binary data for MIC verification
	CmMic    []byte   // Extracted CM MIC value (TLV 6)
	CmtsMic  []byte   // Extracted CMTS MIC value (TLV 7)
	TLVOrder []string // Insertion order of top-level TLV property names
}

// Decode reads the binary DOCSIS config and returns the decoded JSON structure.
func Decode(data []byte, reg *Registry) (*DecodeResult, error) {
	result := &DecodeResult{
		Config:  make(map[string]interface{}),
		RawData: data,
	}

	var unknowns []interface{}
	offset := 0

	for offset < len(data) {
		if offset+1 >= len(data) {
			break
		}

		tlvType := int(data[offset])

		// TLV 0 = pad byte
		if tlvType == 0 {
			offset++
			continue
		}

		// TLV 255 = end-of-data
		if tlvType == 255 {
			break
		}

		// Determine the length field size: look up the registry first,
		// then read 1 or 2 bytes for the length accordingly.
		def, ok := reg.TopLevel[tlvType]
		lengthSize := 1
		if ok && def.LengthSize == 2 {
			lengthSize = 2
		}

		if offset+1+lengthSize > len(data) {
			return nil, fmt.Errorf("truncated TLV at offset %d", offset)
		}

		var tlvLen int
		if lengthSize == 2 {
			tlvLen = int(data[offset+1])<<8 | int(data[offset+2])
		} else {
			tlvLen = int(data[offset+1])
		}
		valueStart := offset + 1 + lengthSize
		valueEnd := valueStart + tlvLen

		if valueEnd > len(data) {
			return nil, fmt.Errorf("TLV type %d at offset %d: length %d exceeds data", tlvType, offset, tlvLen)
		}

		value := data[valueStart:valueEnd]

		// Extract MIC values for later verification.
		if tlvType == 6 {
			result.CmMic = make([]byte, len(value))
			copy(result.CmMic, value)
		}
		if tlvType == 7 {
			result.CmtsMic = make([]byte, len(value))
			copy(result.CmtsMic, value)
		}
		if !ok {
			// Unknown TLV
			hexVal, _ := DecodeValue(value, DataTypeHexString)
			unknowns = append(unknowns, map[string]interface{}{
				"type":  tlvType,
				"value": hexVal,
			})
			offset = valueEnd
			continue
		}

		// Track top-level ordering (first occurrence only for repeatables).
		result.appendTopOrder(def.Name)

		// Chunked TLV (e.g. CVC certificates): concatenate consecutive same-type
		// instances into a single value, then decode as one.
		if def.Chunked {
			var assembled []byte
			assembled = append(assembled, value...)
			scanOffset := valueEnd
			for scanOffset+1+lengthSize <= len(data) {
				nextType := int(data[scanOffset])
				if nextType != tlvType {
					break
				}
				var nextLen int
				if lengthSize == 2 {
					nextLen = int(data[scanOffset+1])<<8 | int(data[scanOffset+2])
				} else {
					nextLen = int(data[scanOffset+1])
				}
				nextEnd := scanOffset + 1 + lengthSize + nextLen
				if nextEnd > len(data) {
					break
				}
				assembled = append(assembled, data[scanOffset+1+lengthSize:nextEnd]...)
				scanOffset = nextEnd
			}
			decoded, err := DecodeValue(assembled, def.DataType)
			if err != nil {
				return nil, fmt.Errorf("decoding chunked TLV %d at offset %d: %w", tlvType, offset, err)
			}
			result.Config[def.Name] = decoded
			offset = scanOffset
			continue
		}

		// Special case: TLV 10 (SNMP Write Access Control) - BER OID + access flag
		if tlvType == 10 {
			err := decodeTLV10(result.Config, def, value)
			if err != nil {
				return nil, fmt.Errorf("decoding TLV 10 at offset %d: %w", offset, err)
			}
			offset = valueEnd
			continue
		}

		// Special case: TLV 11 (SNMP MIB Object) - BER-encoded varbind
		if tlvType == 11 {
			err := decodeTLV11(result.Config, def, value)
			if err != nil {
				return nil, fmt.Errorf("decoding TLV 11 at offset %d: %w", offset, err)
			}
			offset = valueEnd
			continue
		}

		// Special case: TLV 43 (DOCSIS Extension Field) - vendor-gated compound
		if tlvType == 43 {
			err := decodeTLV43(result.Config, def, value, reg)
			if err != nil {
				return nil, fmt.Errorf("decoding TLV 43 at offset %d: %w", offset, err)
			}
			offset = valueEnd
			continue
		}

		// Compound TLV - recurse into sub-TLVs.
		if def.DataType == DataTypeCompound {
			subResult, err := decodeCompound(value, def, reg)
			if err != nil {
				return nil, fmt.Errorf("decoding compound TLV %d at offset %d: %w", tlvType, offset, err)
			}
			addToResult(result.Config, def, subResult)
			offset = valueEnd
			continue
		}

		// Simple TLV
		decoded, err := DecodeValue(value, def.DataType)
		if err != nil {
			return nil, fmt.Errorf("decoding TLV %d at offset %d: %w", tlvType, offset, err)
		}
		addToResult(result.Config, def, decoded)

		offset = valueEnd
	}

	if len(unknowns) > 0 {
		result.Config["UnknownTlvs"] = unknowns
	}

	return result, nil
}

// decodeCompound parses a compound TLV's value as a sequence of sub-TLVs.
func decodeCompound(data []byte, parent *TLVDef, reg *Registry) (map[string]interface{}, error) {
	// VendorSpecificContainer: handle vendor-gated sub-TLVs like TLV 43.
	if parent.RefName == "VendorSpecificContainer" && reg != nil {
		return decodeVendorSpecificContainer(data, reg)
	}

	result := make(map[string]interface{})
	var unknowns []interface{}
	offset := 0

	for offset < len(data) {
		if offset+1 >= len(data) {
			break
		}

		subType := int(data[offset])

		// Pad byte
		if subType == 0 {
			offset++
			continue
		}

		// Determine the length field size from sub-TLV definition.
		subDef, ok := parent.SubTLVs[subType]
		subLengthSize := 1
		if ok && subDef.LengthSize == 2 {
			subLengthSize = 2
		}

		if offset+1+subLengthSize > len(data) {
			return nil, fmt.Errorf("truncated sub-TLV at offset %d", offset)
		}

		var subLen int
		if subLengthSize == 2 {
			subLen = int(data[offset+1])<<8 | int(data[offset+2])
		} else {
			subLen = int(data[offset+1])
		}
		valueStart := offset + 1 + subLengthSize
		valueEnd := valueStart + subLen

		if valueEnd > len(data) {
			return nil, fmt.Errorf("sub-TLV type %d at offset %d: length %d exceeds data", subType, offset, subLen)
		}

		value := data[valueStart:valueEnd]
		if !ok {
			// Unknown sub-TLV
			hexVal, _ := DecodeValue(value, DataTypeHexString)
			unknowns = append(unknowns, map[string]interface{}{
				"type":  subType,
				"value": hexVal,
			})
			offset = valueEnd
			continue
		}

		appendOrder(result, subDef.Name)

		// Nested compound sub-TLV
		if subDef.DataType == DataTypeCompound {
			subResult, err := decodeCompound(value, subDef, reg)
			if err != nil {
				return nil, fmt.Errorf("decoding nested compound sub-TLV %d: %w", subType, err)
			}
			addToResult(result, subDef, subResult)
			offset = valueEnd
			continue
		}

		// Simple sub-TLV
		decoded, err := DecodeValue(value, subDef.DataType)
		if err != nil {
			return nil, fmt.Errorf("decoding sub-TLV %d: %w", subType, err)
		}
		addToResult(result, subDef, decoded)

		offset = valueEnd
	}

	if len(unknowns) > 0 {
		result["UnknownSubTlvs"] = unknowns
	}

	return result, nil
}

// decodeTLV11 handles TLV 11 (SNMP MIB Object) which contains a BER-encoded varbind.
func decodeTLV11(config map[string]interface{}, def *TLVDef, data []byte) error {
	varbind, err := DecodeSnmpVarbind(data)
	if err != nil {
		return err
	}

	entry := map[string]interface{}{
		"oid":   varbind.OID,
		"type":  varbind.Type,
		"value": varbind.Value,
	}

	addToResult(config, def, entry)
	return nil
}

// decodeTLV10 handles TLV 10 (SNMP Write Access Control) which contains a
// BER-encoded OID followed by a 1-byte access flag.
func decodeTLV10(config map[string]interface{}, def *TLVDef, data []byte) error {
	entry, err := DecodeSnmpWriteAccess(data)
	if err != nil {
		return err
	}

	addToResult(config, def, entry)
	return nil
}

// decodeTLV43 handles TLV 43 (DOCSIS Extension Field) with vendor-gated sub-TLVs.
func decodeTLV43(config map[string]interface{}, def *TLVDef, data []byte, reg *Registry) error {
	// First, extract VendorId (sub-TLV 8, 3 bytes).
	vendorId, vendorSubTLVs, err := extractVendorId(data)
	if err != nil {
		return err
	}

	result := make(map[string]interface{})
	result["VendorId"] = vendorId

	if vendorId == "FFFFFF" {
		// General Extension: use defined sub-TLVs from the registry.
		offset := 0
		for offset < len(vendorSubTLVs) {
			if offset+1 >= len(vendorSubTLVs) {
				break
			}

			subType := int(vendorSubTLVs[offset])
			if subType == 0 {
				offset++
				continue
			}

			subLen := int(vendorSubTLVs[offset+1])
			valueStart := offset + 2
			valueEnd := valueStart + subLen
			if valueEnd > len(vendorSubTLVs) {
				return fmt.Errorf("TLV 43 sub-TLV %d: length exceeds data", subType)
			}

			value := vendorSubTLVs[valueStart:valueEnd]

			subDef, ok := def.SubTLVs[subType]
			if !ok {
				// Unknown general extension sub-TLV
				hexVal, _ := DecodeValue(value, DataTypeHexString)
				addUnknownSubTLV(result, subType, hexVal.(string))
				offset = valueEnd
				continue
			}

			appendOrder(result, subDef.Name)

			if subDef.DataType == DataTypeCompound {
				subResult, err := decodeCompound(value, subDef, reg)
				if err != nil {
					return fmt.Errorf("decoding TLV 43 compound sub-TLV %d: %w", subType, err)
				}
				addToResult(result, subDef, subResult)
			} else {
				decoded, err := DecodeValue(value, subDef.DataType)
				if err != nil {
					return fmt.Errorf("decoding TLV 43 sub-TLV %d: %w", subType, err)
				}
				addToResult(result, subDef, decoded)
			}

			offset = valueEnd
		}
	} else {
		// Vendor-specific: collect as VendorSubTlvs array.
		var vendorTlvEntries []interface{}
		offset := 0
		for offset < len(vendorSubTLVs) {
			if offset+1 >= len(vendorSubTLVs) {
				break
			}

			subType := int(vendorSubTLVs[offset])
			if subType == 0 {
				offset++
				continue
			}

			subLen := int(vendorSubTLVs[offset+1])
			valueStart := offset + 2
			valueEnd := valueStart + subLen
			if valueEnd > len(vendorSubTLVs) {
				return fmt.Errorf("vendor sub-TLV %d: length exceeds data", subType)
			}

			value := vendorSubTLVs[valueStart:valueEnd]
			hexVal, _ := DecodeValue(value, DataTypeHexString)

			vendorTlvEntries = append(vendorTlvEntries, map[string]interface{}{
				"type":  subType,
				"value": hexVal,
			})

			offset = valueEnd
		}
		if len(vendorTlvEntries) > 0 {
			result["VendorSubTlvs"] = vendorTlvEntries
		}
	}

	addToResult(config, def, result)
	return nil
}

// extractVendorId extracts the VendorId (sub-TLV 8) from TLV 43 data and returns
// the remaining sub-TLV bytes.
func extractVendorId(data []byte) (string, []byte, error) {
	var remaining []byte
	vendorId := ""
	offset := 0

	for offset < len(data) {
		if offset+1 >= len(data) {
			break
		}

		subType := int(data[offset])
		if subType == 0 {
			offset++
			continue
		}

		subLen := int(data[offset+1])
		valueStart := offset + 2
		valueEnd := valueStart + subLen
		if valueEnd > len(data) {
			return "", nil, fmt.Errorf("sub-TLV %d: length exceeds data", subType)
		}

		if subType == 8 && vendorId == "" {
			// VendorId
			hexVal, _ := DecodeValue(data[valueStart:valueEnd], DataTypeHexString)
			vendorId = hexVal.(string)
		} else {
			remaining = append(remaining, data[offset:valueEnd]...)
		}

		offset = valueEnd
	}

	if vendorId == "" {
		return "", nil, fmt.Errorf("VendorId (sub-TLV 8) not found in TLV 43")
	}

	return vendorId, remaining, nil
}

// decodeVendorSpecificContainer handles VendorSpecificContainer compounds by
// extracting VendorId first, then using the appropriate sub-TLV definitions.
// For FFFFFF (general extension), it uses TLV 43's sub-TLV definitions.
// For other vendor IDs, sub-TLVs are collected as generic VendorSubTlvs.
func decodeVendorSpecificContainer(data []byte, reg *Registry) (map[string]interface{}, error) {
	vendorId, vendorSubTLVs, err := extractVendorId(data)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	result["VendorId"] = vendorId

	if vendorId == "FFFFFF" {
		// General Extension: use TLV 43's sub-TLV definitions.
		extDef := reg.TopLevel[43]

		offset := 0
		for offset < len(vendorSubTLVs) {
			if offset+1 >= len(vendorSubTLVs) {
				break
			}

			subType := int(vendorSubTLVs[offset])
			if subType == 0 {
				offset++
				continue
			}

			subLen := int(vendorSubTLVs[offset+1])
			valueStart := offset + 2
			valueEnd := valueStart + subLen
			if valueEnd > len(vendorSubTLVs) {
				return nil, fmt.Errorf("vendor-specific sub-TLV %d: length exceeds data", subType)
			}

			value := vendorSubTLVs[valueStart:valueEnd]

			// Look up in TLV 43 definitions when available.
			var subDef *TLVDef
			if extDef != nil {
				subDef = extDef.SubTLVs[subType]
			}

			if subDef == nil {
				hexVal, _ := DecodeValue(value, DataTypeHexString)
				addUnknownSubTLV(result, subType, hexVal.(string))
				offset = valueEnd
				continue
			}

			appendOrder(result, subDef.Name)

			if subDef.DataType == DataTypeCompound {
				subResult, err := decodeCompound(value, subDef, reg)
				if err != nil {
					return nil, fmt.Errorf("decoding vendor-specific compound sub-TLV %d: %w", subType, err)
				}
				addToResult(result, subDef, subResult)
			} else {
				decoded, err := DecodeValue(value, subDef.DataType)
				if err != nil {
					return nil, fmt.Errorf("decoding vendor-specific sub-TLV %d: %w", subType, err)
				}
				addToResult(result, subDef, decoded)
			}

			offset = valueEnd
		}
	} else {
		// Vendor-specific: check for vendor schema, otherwise collect as VendorSubTlvs.
		var vendorDefs map[int]*TLVDef
		if reg.VendorSchemas != nil {
			vendorDefs = reg.VendorSchemas[vendorId]
		}

		var vendorTlvEntries []interface{}
		offset := 0
		for offset < len(vendorSubTLVs) {
			if offset+1 >= len(vendorSubTLVs) {
				break
			}

			subType := int(vendorSubTLVs[offset])
			if subType == 0 {
				offset++
				continue
			}

			subLen := int(vendorSubTLVs[offset+1])
			valueStart := offset + 2
			valueEnd := valueStart + subLen
			if valueEnd > len(vendorSubTLVs) {
				return nil, fmt.Errorf("vendor sub-TLV %d: length exceeds data", subType)
			}

			value := vendorSubTLVs[valueStart:valueEnd]

			// Try vendor schema definitions first.
			if vDef := vendorDefs[subType]; vDef != nil {
				appendOrder(result, vDef.Name)
				if vDef.DataType == DataTypeCompound {
					subResult, err := decodeCompound(value, vDef, reg)
					if err != nil {
						return nil, fmt.Errorf("decoding vendor sub-TLV %d: %w", subType, err)
					}
					addToResult(result, vDef, subResult)
				} else {
					decoded, err := DecodeValue(value, vDef.DataType)
					if err != nil {
						return nil, fmt.Errorf("decoding vendor sub-TLV %d: %w", subType, err)
					}
					addToResult(result, vDef, decoded)
				}
			} else {
				// No vendor schema definition — collect as generic entry.
				hexVal, _ := DecodeValue(value, DataTypeHexString)
				vendorTlvEntries = append(vendorTlvEntries, map[string]interface{}{
					"type":  subType,
					"value": hexVal,
				})
			}

			offset = valueEnd
		}
		if len(vendorTlvEntries) > 0 {
			result["VendorSubTlvs"] = vendorTlvEntries
		}
	}

	return result, nil
}

// addToResult adds a decoded value to the result map, handling repeatable TLVs as arrays.
func addToResult(result map[string]interface{}, def *TLVDef, value interface{}) {
	if def.Repeatable {
		existing, ok := result[def.Name]
		if ok {
			arr, isArr := existing.([]interface{})
			if isArr {
				result[def.Name] = append(arr, value)
			} else {
				result[def.Name] = []interface{}{existing, value}
			}
		} else {
			result[def.Name] = []interface{}{value}
		}
	} else {
		result[def.Name] = value
	}
}

// appendTopOrder adds a name to the top-level TLVOrder, avoiding duplicates.
func (r *DecodeResult) appendTopOrder(name string) {
	for _, n := range r.TLVOrder {
		if n == name {
			return
		}
	}
	r.TLVOrder = append(r.TLVOrder, name)
}

// appendOrder adds a name to the _tlvOrder slice in a map, avoiding duplicates
// for repeatable TLVs (the first occurrence records the position).
func appendOrder(m map[string]interface{}, name string) {
	order, _ := m["_tlvOrder"].([]string)
	for _, n := range order {
		if n == name {
			return
		}
	}
	m["_tlvOrder"] = append(order, name)
}

// StripTLVOrder removes "_tlvOrder" keys from a config map and all nested maps,
// preparing it for clean JSON output. The top-level TLVOrder on DecodeResult is
// not affected.
func StripTLVOrder(config map[string]interface{}) {
	delete(config, "_tlvOrder")
	for _, v := range config {
		switch val := v.(type) {
		case map[string]interface{}:
			StripTLVOrder(val)
		case []interface{}:
			for _, elem := range val {
				if m, ok := elem.(map[string]interface{}); ok {
					StripTLVOrder(m)
				}
			}
		}
	}
}

// addUnknownSubTLV appends an unknown sub-TLV to the UnknownSubTlvs array.
func addUnknownSubTLV(result map[string]interface{}, subType int, hexValue string) {
	entry := map[string]interface{}{
		"type":  subType,
		"value": hexValue,
	}
	existing, ok := result["UnknownSubTlvs"]
	if ok {
		arr := existing.([]interface{})
		result["UnknownSubTlvs"] = append(arr, entry)
	} else {
		result["UnknownSubTlvs"] = []interface{}{entry}
	}
}
