package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"fmt"
)

// MICResult holds the result of a MIC verification.
type MICResult struct {
	Valid    bool
	Expected []byte
	Computed []byte
}

// VerifyCmMic verifies the CM MIC (TLV 6) per MULPI Annex D.
// The digest is a plain MD5 hash over all configuration setting TLV bytes in order,
// excluding TLV 6 (CM MIC), TLV 7 (CMTS MIC), and TLV 255 (end-of-data marker).
func VerifyCmMic(configData []byte, cmMic []byte) *MICResult {
	filtered := filterTLVs(configData, 6, 7, 255)
	h := md5.New()
	h.Write(filtered)
	computed := h.Sum(nil)

	return &MICResult{
		Valid:    bytes.Equal(computed, cmMic),
		Expected: cmMic,
		Computed: computed,
	}
}

// VerifyCmtsMic verifies the CMTS MIC (TLV 7) using HMAC-MD5 with the shared secret.
// The digest is computed over all configuration setting TLV bytes in order,
// excluding TLV 7 (CMTS MIC) and TLV 255 (end-of-data marker).
func VerifyCmtsMic(configData []byte, cmtsMic []byte, sharedSecret string) *MICResult {
	filtered := filterTLVs(configData, 7, 255)
	mac := hmac.New(md5.New, []byte(sharedSecret))
	mac.Write(filtered)
	computed := mac.Sum(nil)

	return &MICResult{
		Valid:    bytes.Equal(computed, cmtsMic),
		Expected: cmtsMic,
		Computed: computed,
	}
}

// filterTLVs returns a copy of the config data with the specified TLV types removed.
// It walks the TLV stream and copies all TLVs except the excluded types.
func filterTLVs(data []byte, excludeTypes ...int) []byte {
	excludeSet := make(map[int]bool)
	for _, t := range excludeTypes {
		excludeSet[t] = true
	}

	var result []byte
	offset := 0
	for offset < len(data) {
		if offset+1 >= len(data) {
			break
		}

		tlvType := int(data[offset])

		// TLV 0 = pad byte (no length)
		if tlvType == 0 {
			if !excludeSet[0] {
				result = append(result, data[offset])
			}
			offset++
			continue
		}

		// End-of-data marker
		if tlvType == 255 {
			if !excludeSet[255] {
				result = append(result, data[offset])
			}
			offset++
			// TLV 255 has a length byte of 0
			if offset < len(data) {
				if !excludeSet[255] {
					result = append(result, data[offset])
				}
				offset++
			}
			break
		}

		if offset+1 >= len(data) {
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

// ExtractTLVValue finds the first TLV of the given type and returns its value bytes.
func ExtractTLVValue(data []byte, tlvType int) ([]byte, error) {
	offset := 0
	for offset < len(data) {
		if offset+1 >= len(data) {
			break
		}

		t := int(data[offset])

		if t == 0 {
			offset++
			continue
		}

		if t == 255 {
			break
		}

		if offset+1 >= len(data) {
			break
		}

		l := int(data[offset+1])
		end := offset + 2 + l
		if end > len(data) {
			break
		}

		if t == tlvType {
			return data[offset+2 : end], nil
		}

		offset = end
	}

	return nil, fmt.Errorf("TLV type %d not found", tlvType)
}
