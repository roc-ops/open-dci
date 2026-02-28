package opendci

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"os"
	"testing"
)

func TestVerifyCmMicValid(t *testing.T) {
	// Build a config: TLV 3 (NetworkAccess=1), then TLV 6 (CM MIC), TLV 7 (CMTS MIC)
	configWithoutMICs := buildConfig(
		buildTLV(3, []byte{1}),
		buildTLV(18, []byte{16}),
	)

	// Remove end-of-data marker for MIC computation input
	// Actually, the MIC is computed over the TLV stream excluding TLV 6 and 7.
	// Build the raw bytes including all TLVs.
	tlv3 := buildTLV(3, []byte{1})
	tlv18 := buildTLV(18, []byte{16})
	endMarker := []byte{0xFF, 0x00}

	// Compute CM MIC: plain MD5 over TLVs excluding 6, 7, and 255 (per MULPI Annex D)
	bytesForMic := append(tlv3, tlv18...)

	h := md5.New()
	h.Write(bytesForMic)
	cmMicValue := h.Sum(nil)

	// Build the full config with the MIC
	var fullConfig []byte
	fullConfig = append(fullConfig, tlv3...)
	fullConfig = append(fullConfig, tlv18...)
	fullConfig = append(fullConfig, buildTLV(6, cmMicValue)...)
	fullConfig = append(fullConfig, endMarker...)

	result := VerifyCmMic(fullConfig, cmMicValue)
	if !result.Valid {
		t.Errorf("CM MIC should be valid, computed=%X expected=%X", result.Computed, result.Expected)
	}

	_ = configWithoutMICs // Suppress unused warning in a clean way
}

func TestVerifyCmMicInvalid(t *testing.T) {
	tlv3 := buildTLV(3, []byte{1})
	endMarker := []byte{0xFF, 0x00}

	// Use a wrong MIC
	wrongMic := make([]byte, 16)
	for i := range wrongMic {
		wrongMic[i] = 0xFF
	}

	var fullConfig []byte
	fullConfig = append(fullConfig, tlv3...)
	fullConfig = append(fullConfig, buildTLV(6, wrongMic)...)
	fullConfig = append(fullConfig, endMarker...)

	result := VerifyCmMic(fullConfig, wrongMic)
	if result.Valid {
		t.Error("CM MIC should be invalid with wrong value")
	}
}

func TestVerifyCmtsMicValid(t *testing.T) {
	tlv3 := buildTLV(3, []byte{1})
	endMarker := []byte{0xFF, 0x00}
	secret := "mysecret"

	// Compute CM MIC first (plain MD5 per MULPI Annex D, excludes TLV 255)
	h := md5.New()
	h.Write(tlv3)
	cmMicValue := h.Sum(nil)
	cmMicTLV := buildTLV(6, cmMicValue)

	// Compute CMTS MIC: HMAC-MD5 over TLVs in canonical digest order.
	// Digest order: {1,2,3,4,17,43,6,...} — so type 3 before type 6.
	mac2 := hmac.New(md5.New, []byte(secret))
	mac2.Write(tlv3)
	mac2.Write(cmMicTLV)
	cmtsMicValue := mac2.Sum(nil)

	var fullConfig []byte
	fullConfig = append(fullConfig, tlv3...)
	fullConfig = append(fullConfig, cmMicTLV...)
	fullConfig = append(fullConfig, buildTLV(7, cmtsMicValue)...)
	fullConfig = append(fullConfig, endMarker...)

	result := VerifyCmtsMic(fullConfig, cmtsMicValue, secret)
	if !result.Valid {
		t.Errorf("CMTS MIC should be valid, computed=%X expected=%X", result.Computed, result.Expected)
	}
}

func TestVerifyCmtsMicInvalidSecret(t *testing.T) {
	tlv3 := buildTLV(3, []byte{1})
	endMarker := []byte{0xFF, 0x00}

	// Compute with one secret (excludes TLV 255)
	mac := hmac.New(md5.New, []byte("secret1"))
	mac.Write(tlv3)
	cmtsMicValue := mac.Sum(nil)

	var fullConfig []byte
	fullConfig = append(fullConfig, tlv3...)
	fullConfig = append(fullConfig, buildTLV(7, cmtsMicValue)...)
	fullConfig = append(fullConfig, endMarker...)

	// Verify with wrong secret
	result := VerifyCmtsMic(fullConfig, cmtsMicValue, "secret2")
	if result.Valid {
		t.Error("CMTS MIC should be invalid with wrong secret")
	}
}

func TestFilterTLVs(t *testing.T) {
	tlv3 := buildTLV(3, []byte{1})
	tlv6 := buildTLV(6, make([]byte, 16))
	tlv7 := buildTLV(7, make([]byte, 16))
	endMarker := []byte{0xFF, 0x00}

	var data []byte
	data = append(data, tlv3...)
	data = append(data, tlv6...)
	data = append(data, tlv7...)
	data = append(data, endMarker...)

	// Filter out TLV 6 and 7
	filtered := filterTLVs(data, 6, 7)

	// Should contain only TLV 3 and end marker
	expected := append(tlv3, endMarker...)
	if len(filtered) != len(expected) {
		t.Errorf("expected %d bytes, got %d", len(expected), len(filtered))
	}

	for i := range expected {
		if filtered[i] != expected[i] {
			t.Errorf("byte %d: expected 0x%02X, got 0x%02X", i, expected[i], filtered[i])
			break
		}
	}
}

func TestFilterTLVsOnlyTLV7(t *testing.T) {
	tlv3 := buildTLV(3, []byte{1})
	tlv6 := buildTLV(6, make([]byte, 16))
	tlv7 := buildTLV(7, make([]byte, 16))
	endMarker := []byte{0xFF, 0x00}

	var data []byte
	data = append(data, tlv3...)
	data = append(data, tlv6...)
	data = append(data, tlv7...)
	data = append(data, endMarker...)

	// Filter out only TLV 7
	filtered := filterTLVs(data, 7)

	// Should contain TLV 3, TLV 6, and end marker
	var expected []byte
	expected = append(expected, tlv3...)
	expected = append(expected, tlv6...)
	expected = append(expected, endMarker...)
	if len(filtered) != len(expected) {
		t.Errorf("expected %d bytes, got %d", len(expected), len(filtered))
	}
}

func TestExtractTLVValue(t *testing.T) {
	tlv3 := buildTLV(3, []byte{1})
	mic := make([]byte, 16)
	for i := range mic {
		mic[i] = byte(i)
	}
	tlv6 := buildTLV(6, mic)
	endMarker := []byte{0xFF, 0x00}

	var data []byte
	data = append(data, tlv3...)
	data = append(data, tlv6...)
	data = append(data, endMarker...)

	val, err := ExtractTLVValue(data, 6)
	if err != nil {
		t.Fatal(err)
	}

	if len(val) != 16 {
		t.Errorf("expected 16 bytes, got %d", len(val))
	}
}

func TestExtractTLVValueNotFound(t *testing.T) {
	data := buildConfig(
		buildTLV(3, []byte{1}),
	)

	_, err := ExtractTLVValue(data, 99)
	if err == nil {
		t.Fatal("expected error for missing TLV")
	}
}

func TestReorderForCmtsMic(t *testing.T) {
	tlv1 := buildTLV(1, []byte{0, 0, 1, 0})   // DownstreamFrequency
	tlv3 := buildTLV(3, []byte{1})              // NetworkAccess
	tlv18 := buildTLV(18, []byte{16})           // MaxCPE

	// Wire order: 18, 3, 1 (reversed from digest order)
	var data []byte
	data = append(data, tlv18...)
	data = append(data, tlv3...)
	data = append(data, tlv1...)

	reordered := reorderForCmtsMic(data)

	// Digest order: {1, 2, 3, 4, 17, 43, 6, 18, ...} → type 1, then 3, then 18
	var expected []byte
	expected = append(expected, tlv1...)
	expected = append(expected, tlv3...)
	expected = append(expected, tlv18...)

	if !bytes.Equal(reordered, expected) {
		t.Errorf("reorder mismatch:\n  got:      %X\n  expected: %X", reordered, expected)
	}
}

func TestReorderForCmtsMicExcludesNonDigestTypes(t *testing.T) {
	tlv3 := buildTLV(3, []byte{1})
	tlv99 := buildTLV(99, []byte{0xAA, 0xBB}) // Not in digest order list
	tlv18 := buildTLV(18, []byte{16})

	var data []byte
	data = append(data, tlv3...)
	data = append(data, tlv99...)
	data = append(data, tlv18...)

	reordered := reorderForCmtsMic(data)

	// TLV 99 should be excluded entirely.
	var expected []byte
	expected = append(expected, tlv3...)
	expected = append(expected, tlv18...)

	if !bytes.Equal(reordered, expected) {
		t.Errorf("non-digest TLV not excluded:\n  got:      %X\n  expected: %X", reordered, expected)
	}
}

func TestReorderForCmtsMicMultipleInstances(t *testing.T) {
	// Two TLV 24 (UpstreamServiceFlow) instances — repeatable TLV in digest order.
	tlv24a := buildTLV(24, []byte{0x01, 0x02, 0x00, 0x01})
	tlv24b := buildTLV(24, []byte{0x01, 0x02, 0x00, 0x02})
	tlv3 := buildTLV(3, []byte{1})

	// Wire order: 24(a), 3, 24(b)
	var data []byte
	data = append(data, tlv24a...)
	data = append(data, tlv3...)
	data = append(data, tlv24b...)

	reordered := reorderForCmtsMic(data)

	// Digest order: 3 before 24; both 24 instances preserved in wire order.
	var expected []byte
	expected = append(expected, tlv3...)
	expected = append(expected, tlv24a...)
	expected = append(expected, tlv24b...)

	if !bytes.Equal(reordered, expected) {
		t.Errorf("multiple instance reorder mismatch:\n  got:      %X\n  expected: %X", reordered, expected)
	}
}

func TestCmtsMicDigestOrderWireOrderMismatch(t *testing.T) {
	// Verify that CMTS MIC computation uses digest order, not wire order.
	// Wire order: TLV 18, TLV 3, TLV 6 (CM MIC)
	// Digest order: TLV 3, TLV 6, TLV 18
	secret := "testsecret"

	tlv3 := buildTLV(3, []byte{1})
	tlv18 := buildTLV(18, []byte{16})

	// Compute CM MIC (uses wire order, includes both TLVs).
	h := md5.New()
	h.Write(tlv18)
	h.Write(tlv3)
	cmMicValue := h.Sum(nil)
	cmMicTLV := buildTLV(6, cmMicValue)

	// Build input for ComputeCmtsMic: wire order is 18, 3, 6.
	var inputBytes []byte
	inputBytes = append(inputBytes, tlv18...)
	inputBytes = append(inputBytes, tlv3...)
	inputBytes = append(inputBytes, cmMicTLV...)

	computed := ComputeCmtsMic(inputBytes, secret)

	// Manually compute expected: digest order is 3, 6, 18.
	mac := hmac.New(md5.New, []byte(secret))
	mac.Write(tlv3)
	mac.Write(cmMicTLV)
	mac.Write(tlv18)
	expected := mac.Sum(nil)

	if !bytes.Equal(computed, expected) {
		t.Errorf("CMTS MIC should use digest order, not wire order:\n  computed: %X\n  expected: %X", computed, expected)
	}

	// Also verify that wire order would give a DIFFERENT result (proving reorder matters).
	macWire := hmac.New(md5.New, []byte(secret))
	macWire.Write(tlv18)
	macWire.Write(tlv3)
	macWire.Write(cmMicTLV)
	wireOrderResult := macWire.Sum(nil)

	if bytes.Equal(computed, wireOrderResult) {
		t.Error("digest-order and wire-order results should differ when wire order != digest order")
	}
}

func TestVerifyCmtsMicWithReorderedWireData(t *testing.T) {
	// Full end-to-end: build config with non-standard wire order,
	// compute MIC with digest ordering, verify it.
	secret := "e2etest"

	tlv3 := buildTLV(3, []byte{1})
	tlv18 := buildTLV(18, []byte{16})
	tlv1 := buildTLV(1, []byte{0, 0, 1, 0})

	// Wire order: 18, 1, 3 (not digest order)
	var body []byte
	body = append(body, tlv18...)
	body = append(body, tlv1...)
	body = append(body, tlv3...)

	// Compute CM MIC in wire order (CM MIC is not reordered).
	cmMic := ComputeCmMic(body)
	cmMicTLV := buildTLV(6, cmMic)

	// Compute CMTS MIC (should internally reorder to: 1, 3, 6, 18).
	withCmMic := append(body, cmMicTLV...)
	cmtsMic := ComputeCmtsMic(withCmMic, secret)
	cmtsMicTLV := buildTLV(7, cmtsMic)

	// Build full config: body + CM MIC + CMTS MIC + end-of-data
	var fullConfig []byte
	fullConfig = append(fullConfig, body...)
	fullConfig = append(fullConfig, cmMicTLV...)
	fullConfig = append(fullConfig, cmtsMicTLV...)
	fullConfig = append(fullConfig, 0xFF, 0x00)

	result := VerifyCmtsMic(fullConfig, cmtsMic, secret)
	if !result.Valid {
		t.Errorf("CMTS MIC should be valid with reordered wire data, computed=%X expected=%X",
			result.Computed, result.Expected)
	}
}

func TestCmMicLabTr069(t *testing.T) {
	data, err := os.ReadFile("binary-files/lab-tr069.bin")
	if err != nil {
		t.Skip("binary-files/lab-tr069.bin not found, skipping integration test")
	}

	cmMic, err := ExtractTLVValue(data, 6)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Expected CM MIC: %X", cmMic)

	result := VerifyCmMic(data, cmMic)
	if !result.Valid {
		t.Errorf("CM MIC should be valid, expected=%X computed=%X", result.Expected, result.Computed)
	}
	t.Logf("Computed CM MIC: %X (valid=%v)", result.Computed, result.Valid)
}
