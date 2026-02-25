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

// cmtsMicDigestOrder is the canonical TLV type ordering for CMTS MIC computation
// per CM-SP-MULPIv4.0 Annex D. Only TLV types in this list are included in the
// HMAC-MD5 digest, and they are concatenated in this order regardless of wire order.
var cmtsMicDigestOrder = []int{1, 2, 3, 4, 17, 43, 6, 18, 19, 20, 22, 23, 24, 25, 28, 29, 26, 35, 36, 37, 40}

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
// Per MULPI Annex D, the digest is computed over TLV bytes reordered by the canonical
// digest order, excluding TLV 7 (CMTS MIC) and TLV 255 (end-of-data marker).
func VerifyCmtsMic(configData []byte, cmtsMic []byte, sharedSecret string) *MICResult {
	filtered := filterTLVs(configData, 7, 255)
	reordered := reorderForCmtsMic(filtered)
	mac := hmac.New(md5.New, []byte(sharedSecret))
	mac.Write(reordered)
	computed := mac.Sum(nil)

	return &MICResult{
		Valid:    bytes.Equal(computed, cmtsMic),
		Expected: cmtsMic,
		Computed: computed,
	}
}

// ComputeCmMic computes the CM MIC (TLV 6) for encoded config bytes.
// The input should be the TLV stream WITHOUT TLV 6, TLV 7, and TLV 255.
// Returns the 16-byte MD5 digest.
func ComputeCmMic(configTLVBytes []byte) []byte {
	h := md5.New()
	h.Write(configTLVBytes)
	return h.Sum(nil)
}

// ComputeCmtsMic computes the CMTS MIC (TLV 7) for encoded config bytes.
// The input should be the TLV stream WITHOUT TLV 7 and TLV 255 (but including
// TLV 6 / CM MIC). TLVs are reordered per the canonical CMTS MIC digest order
// (MULPI Annex D) before HMAC computation. Returns the 16-byte HMAC-MD5 digest.
func ComputeCmtsMic(configTLVBytes []byte, sharedSecret string) []byte {
	reordered := reorderForCmtsMic(configTLVBytes)
	mac := hmac.New(md5.New, []byte(sharedSecret))
	mac.Write(reordered)
	return mac.Sum(nil)
}

// reorderForCmtsMic extracts TLV instances from the binary stream and
// concatenates them in the canonical CMTS MIC digest order (MULPI Annex D).
// Only TLV types listed in cmtsMicDigestOrder are included; all others are dropped.
// Multiple instances of the same type are preserved in their original wire order.
//
// Note: This function assumes 1-byte TLV length fields. TLVs with 2-byte length
// fields (e.g. 103, 104) are not in the digest order list and would be skipped,
// but if encountered they may cause misparsing of subsequent TLVs.
func reorderForCmtsMic(data []byte) []byte {
	// Build the set of types we care about for fast lookup.
	digestSet := make(map[int]bool, len(cmtsMicDigestOrder))
	for _, t := range cmtsMicDigestOrder {
		digestSet[t] = true
	}

	// Extract TLV instances grouped by type number.
	tlvsByType := make(map[int][][]byte)
	offset := 0
	for offset < len(data) {
		if offset+1 >= len(data) {
			break
		}

		tlvType := int(data[offset])

		// Pad byte (no length field).
		if tlvType == 0 {
			offset++
			continue
		}

		// End-of-data marker.
		if tlvType == 255 {
			break
		}

		tlvLen := int(data[offset+1])
		end := offset + 2 + tlvLen
		if end > len(data) {
			break
		}

		if digestSet[tlvType] {
			tlvBytes := make([]byte, end-offset)
			copy(tlvBytes, data[offset:end])
			tlvsByType[tlvType] = append(tlvsByType[tlvType], tlvBytes)
		}

		offset = end
	}

	// Concatenate in the canonical digest order.
	var result []byte
	for _, t := range cmtsMicDigestOrder {
		for _, tlvBytes := range tlvsByType[t] {
			result = append(result, tlvBytes...)
		}
	}

	return result
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

// insertMICs computes CM MIC and CMTS MIC and inserts them before the end-of-data
// marker in the encoded binary. The input must end with 0xFF 0x00.
func insertMICs(encoded []byte, cmtsSecret string) ([]byte, error) {
	// Remove end-of-data marker.
	if len(encoded) < 2 || encoded[len(encoded)-2] != 0xFF || encoded[len(encoded)-1] != 0x00 {
		return nil, fmt.Errorf("encoded data does not end with end-of-data marker")
	}
	body := encoded[:len(encoded)-2]

	// Compute CM MIC (TLV 6): MD5 of all TLV bytes excluding TLV 6, 7, 255.
	cmMic := ComputeCmMic(body)
	cmMicTLV := makeTLV(6, cmMic)

	// Build data with CM MIC for CMTS MIC computation.
	withCmMic := append(body, cmMicTLV...)

	// Compute CMTS MIC (TLV 7): HMAC-MD5 of all TLV bytes excluding TLV 7, 255.
	cmtsMic := ComputeCmtsMic(withCmMic, cmtsSecret)
	cmtsMicTLV := makeTLV(7, cmtsMic)

	// Reassemble: body + CM MIC + CMTS MIC + end-of-data.
	var result []byte
	result = append(result, body...)
	result = append(result, cmMicTLV...)
	result = append(result, cmtsMicTLV...)
	result = append(result, 0xFF, 0x00)

	return result, nil
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
