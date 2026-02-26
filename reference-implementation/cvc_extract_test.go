package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"math/big"
	"strings"
	"testing"
	"time"

	"go.mozilla.org/pkcs7"
)

// createTestCert creates a self-signed or issuer-signed X.509 certificate
// with optional codeSigning EKU and CA basic constraints.
func createTestCert(t *testing.T, cn string, isCA bool, codeSigning bool, issuerCert *x509.Certificate, issuerKey *rsa.PrivateKey) (*x509.Certificate, *rsa.PrivateKey) {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generating RSA key for %s: %v", cn, err)
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		t.Fatalf("generating serial number: %v", err)
	}

	tmpl := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   cn,
			Organization: []string{"Test"},
		},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  isCA,
	}

	if isCA {
		tmpl.KeyUsage |= x509.KeyUsageCertSign
	}

	if codeSigning {
		tmpl.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageCodeSigning}
	}

	// Self-sign if no issuer provided.
	signerCert := tmpl
	signerKey := key
	if issuerCert != nil && issuerKey != nil {
		signerCert = issuerCert
		signerKey = issuerKey
	}

	certDER, err := x509.CreateCertificate(rand.Reader, tmpl, signerCert, &key.PublicKey, signerKey)
	if err != nil {
		t.Fatalf("creating certificate %s: %v", cn, err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		t.Fatalf("parsing certificate %s: %v", cn, err)
	}

	return cert, key
}

// buildTestPKCS7Firmware creates a PKCS#7 SignedData structure with the given
// content and certificates, simulating a signed firmware binary.
func buildTestPKCS7Firmware(t *testing.T, content []byte, signerCert *x509.Certificate, signerKey *rsa.PrivateKey, extraCerts []*x509.Certificate) []byte {
	t.Helper()

	sd, err := pkcs7.NewSignedData(content)
	if err != nil {
		t.Fatalf("creating SignedData: %v", err)
	}

	sd.SetDigestAlgorithm(pkcs7.OIDDigestAlgorithmSHA256)

	if err := sd.AddSigner(signerCert, signerKey, pkcs7.SignerInfoConfig{}); err != nil {
		t.Fatalf("adding signer: %v", err)
	}

	for _, cert := range extraCerts {
		sd.AddCertificate(cert)
	}

	signed, err := sd.Finish()
	if err != nil {
		t.Fatalf("finishing SignedData: %v", err)
	}

	return signed
}

func TestExtractCVCFromFirmware_SingleManufacturerCVC(t *testing.T) {
	// Create a CA cert (no codeSigning) and a CVC cert (with codeSigning).
	caCert, caKey := createTestCert(t, "Test CA", true, false, nil, nil)
	cvcCert, cvcKey := createTestCert(t, "Manufacturer CVC", false, true, caCert, caKey)

	// Build PKCS#7 firmware signed by the CVC.
	firmware := []byte("test firmware content")
	p7Data := buildTestPKCS7Firmware(t, firmware, cvcCert, cvcKey, []*x509.Certificate{caCert})

	result, err := ExtractCVCFromFirmware(p7Data)
	if err != nil {
		t.Fatalf("ExtractCVCFromFirmware error: %v", err)
	}

	// ManufacturerCvc should be present.
	mfgCvc, ok := result["ManufacturerCvc"]
	if !ok || mfgCvc == nil {
		t.Fatal("ManufacturerCvc should be present")
	}

	// Verify the hex matches the CVC's raw DER.
	expectedHex := strings.ToUpper(hex.EncodeToString(cvcCert.Raw))
	if mfgCvc.(string) != expectedHex {
		t.Errorf("ManufacturerCvc mismatch:\n  got:      %s...\n  expected: %s...", mfgCvc.(string)[:40], expectedHex[:40])
	}

	// CoSignerCvc should be nil.
	if result["CoSignerCvc"] != nil {
		t.Error("CoSignerCvc should be nil when only one CVC is present")
	}

	// ManufacturerCvcChain should be present (degenerate PKCS#7).
	if result["ManufacturerCvcChain"] == nil {
		t.Error("ManufacturerCvcChain should be present")
	}

	// CoSignerCvcChain should be nil.
	if result["CoSignerCvcChain"] != nil {
		t.Error("CoSignerCvcChain should be nil when only one CVC is present")
	}
}

func TestExtractCVCFromFirmware_DualCVC(t *testing.T) {
	// Create two separate CA + CVC pairs.
	ca1Cert, ca1Key := createTestCert(t, "Manufacturer CA", true, false, nil, nil)
	cvc1Cert, cvc1Key := createTestCert(t, "Manufacturer CVC", false, true, ca1Cert, ca1Key)

	ca2Cert, _ := createTestCert(t, "CoSigner CA", true, false, nil, nil)
	cvc2Cert, _ := createTestCert(t, "CoSigner CVC", false, true, ca2Cert, nil)

	// Build PKCS#7 with manufacturer CVC as signer, plus co-signer CVC and both CAs as extra certs.
	firmware := []byte("test firmware with dual CVC")
	p7Data := buildTestPKCS7Firmware(t, firmware, cvc1Cert, cvc1Key,
		[]*x509.Certificate{ca1Cert, cvc2Cert, ca2Cert})

	result, err := ExtractCVCFromFirmware(p7Data)
	if err != nil {
		t.Fatalf("ExtractCVCFromFirmware error: %v", err)
	}

	// Both CVCs should be present.
	if result["ManufacturerCvc"] == nil {
		t.Error("ManufacturerCvc should be present")
	}
	if result["CoSignerCvc"] == nil {
		t.Error("CoSignerCvc should be present")
	}

	// Both chains should be present.
	if result["ManufacturerCvcChain"] == nil {
		t.Error("ManufacturerCvcChain should be present")
	}
	if result["CoSignerCvcChain"] == nil {
		t.Error("CoSignerCvcChain should be present")
	}

	// Verify they are different certificates.
	if result["ManufacturerCvc"].(string) == result["CoSignerCvc"].(string) {
		t.Error("ManufacturerCvc and CoSignerCvc should be different certificates")
	}
}

func TestExtractCVCFromFirmware_InvalidInput(t *testing.T) {
	// Non-PKCS#7 data should return an error.
	_, err := ExtractCVCFromFirmware([]byte("this is not PKCS#7"))
	if err == nil {
		t.Error("expected error for non-PKCS#7 input")
	}
}

func TestExtractCVCFromFirmware_EmptyInput(t *testing.T) {
	_, err := ExtractCVCFromFirmware([]byte{})
	if err == nil {
		t.Error("expected error for empty input")
	}
}

func TestExtractCVCFromFirmware_NoCodeSigningCerts(t *testing.T) {
	// Create a cert without codeSigning EKU.
	cert, key := createTestCert(t, "Non-CVC Signer", false, false, nil, nil)

	firmware := []byte("unsigned firmware content")
	p7Data := buildTestPKCS7Firmware(t, firmware, cert, key, nil)

	result, err := ExtractCVCFromFirmware(p7Data)
	if err != nil {
		t.Fatalf("ExtractCVCFromFirmware error: %v", err)
	}

	// All fields should be nil since no codeSigning certs exist.
	for _, key := range []string{"ManufacturerCvc", "CoSignerCvc", "ManufacturerCvcChain", "CoSignerCvcChain"} {
		if result[key] != nil {
			t.Errorf("%s should be nil when no codeSigning certs are present", key)
		}
	}
}

func TestExtractCVCFromFirmware_ChainContainsCACert(t *testing.T) {
	// Create a CA and CVC, verify the chain's degenerate PKCS#7 can be parsed
	// and contains both the CVC and the CA.
	caCert, caKey := createTestCert(t, "Test Root CA", true, false, nil, nil)
	cvcCert, cvcKey := createTestCert(t, "Test CVC", false, true, caCert, caKey)

	firmware := []byte("firmware")
	p7Data := buildTestPKCS7Firmware(t, firmware, cvcCert, cvcKey, []*x509.Certificate{caCert})

	result, err := ExtractCVCFromFirmware(p7Data)
	if err != nil {
		t.Fatalf("ExtractCVCFromFirmware error: %v", err)
	}

	// Parse the chain PKCS#7 and verify it contains certificates.
	chainHex := result["ManufacturerCvcChain"].(string)
	chainDER, err := hex.DecodeString(chainHex)
	if err != nil {
		t.Fatalf("decoding chain hex: %v", err)
	}

	chainP7, err := pkcs7.Parse(chainDER)
	if err != nil {
		t.Fatalf("parsing chain PKCS#7: %v", err)
	}

	if len(chainP7.Certificates) == 0 {
		t.Fatal("chain PKCS#7 should contain certificates")
	}

	// At minimum, the CVC itself should be in the chain.
	foundCVC := false
	for _, cert := range chainP7.Certificates {
		if cert.Subject.CommonName == "Test CVC" {
			foundCVC = true
		}
	}
	if !foundCVC {
		t.Error("chain should contain the CVC certificate")
	}
}

func TestHasCodeSigningEKU(t *testing.T) {
	// Certificate with codeSigning EKU.
	csCert, _ := createTestCert(t, "CodeSigning", false, true, nil, nil)
	if !hasCodeSigningEKU(csCert) {
		t.Error("expected codeSigning EKU to be detected")
	}

	// Certificate without codeSigning EKU.
	noCert, _ := createTestCert(t, "NoCodeSigning", false, false, nil, nil)
	if hasCodeSigningEKU(noCert) {
		t.Error("expected codeSigning EKU to NOT be detected")
	}
}

func TestBuildCertChain(t *testing.T) {
	// Create root -> intermediate -> leaf chain.
	rootCert, rootKey := createTestCert(t, "Root", true, false, nil, nil)
	intCert, intKey := createTestCert(t, "Intermediate", true, false, rootCert, rootKey)
	leafCert, _ := createTestCert(t, "Leaf", false, true, intCert, intKey)

	caCerts := []*x509.Certificate{intCert, rootCert}
	chain := buildCertChain(leafCert, caCerts)

	// Chain should be: leaf, intermediate, root.
	if len(chain) != 3 {
		t.Fatalf("expected chain length 3, got %d", len(chain))
	}

	// Verify first cert is the leaf.
	parsed, err := x509.ParseCertificate(chain[0])
	if err != nil {
		t.Fatalf("parsing chain[0]: %v", err)
	}
	if parsed.Subject.CommonName != "Leaf" {
		t.Errorf("chain[0] should be Leaf, got %s", parsed.Subject.CommonName)
	}

	// Verify last cert is root.
	parsedRoot, err := x509.ParseCertificate(chain[2])
	if err != nil {
		t.Fatalf("parsing chain[2]: %v", err)
	}
	if parsedRoot.Subject.CommonName != "Root" {
		t.Errorf("chain[2] should be Root, got %s", parsedRoot.Subject.CommonName)
	}
}

func TestBuildDegeneratePKCS7(t *testing.T) {
	cert, _ := createTestCert(t, "Test", false, false, nil, nil)

	degDER, err := buildDegeneratePKCS7([][]byte{cert.Raw})
	if err != nil {
		t.Fatalf("buildDegeneratePKCS7 error: %v", err)
	}

	// Should be parseable as PKCS#7.
	p7, err := pkcs7.Parse(degDER)
	if err != nil {
		t.Fatalf("parsing degenerate PKCS#7: %v", err)
	}

	if len(p7.Certificates) != 1 {
		t.Errorf("expected 1 certificate, got %d", len(p7.Certificates))
	}
}

func TestBuildDegeneratePKCS7_Empty(t *testing.T) {
	_, err := buildDegeneratePKCS7(nil)
	if err == nil {
		t.Error("expected error for empty cert list")
	}
}
