package main

import (
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

	// Compute CMTS MIC: HMAC-MD5 over all TLVs except TLV 7 and TLV 255
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
