package main

import (
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"strings"

	"go.mozilla.org/pkcs7"
)

// ExtractCVCFromFirmware parses a PKCS#7-signed firmware binary and extracts
// the embedded CVC (Code Verification Certificate) certificates.
//
// Signed cable modem firmware files are DER-encoded PKCS#7 SignedData structures.
// CVCs are identified by the codeSigning extended key usage (OID 1.3.6.1.5.5.7.3.3).
// The first CVC with codeSigning EKU is the manufacturer CVC; the second (if present)
// is the co-signer CVC.
//
// Returns a map with keys:
//   - ManufacturerCvc:      uppercase hex DER of the manufacturer CVC, or nil
//   - CoSignerCvc:          uppercase hex DER of the co-signer CVC, or nil
//   - ManufacturerCvcChain: uppercase hex degenerate PKCS#7 containing the manufacturer CVC chain, or nil
//   - CoSignerCvcChain:     uppercase hex degenerate PKCS#7 containing the co-signer CVC chain, or nil
func ExtractCVCFromFirmware(data []byte) (map[string]interface{}, error) {
	result := map[string]interface{}{
		"ManufacturerCvc":      nil,
		"CoSignerCvc":          nil,
		"ManufacturerCvcChain": nil,
		"CoSignerCvcChain":     nil,
	}

	p7, err := pkcs7.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parsing PKCS#7 firmware: %w", err)
	}

	if len(p7.Certificates) == 0 {
		return result, nil // Valid PKCS#7 but no certificates
	}

	// Separate CVC certificates (codeSigning EKU) from CA certificates.
	var cvcs []*x509.Certificate
	var caCerts []*x509.Certificate
	for _, cert := range p7.Certificates {
		if hasCodeSigningEKU(cert) {
			cvcs = append(cvcs, cert)
		} else {
			caCerts = append(caCerts, cert)
		}
	}

	// Classify CVCs: first = manufacturer, second = co-signer.
	if len(cvcs) >= 1 {
		result["ManufacturerCvc"] = strings.ToUpper(hex.EncodeToString(cvcs[0].Raw))
	}
	if len(cvcs) >= 2 {
		result["CoSignerCvc"] = strings.ToUpper(hex.EncodeToString(cvcs[1].Raw))
	}

	// Build CvcChain (degenerate PKCS#7) for each CVC with its issuing CA chain.
	if len(cvcs) >= 1 {
		chain := buildCertChain(cvcs[0], caCerts)
		if chainDER, err := buildDegeneratePKCS7(chain); err == nil {
			result["ManufacturerCvcChain"] = strings.ToUpper(hex.EncodeToString(chainDER))
		}
	}
	if len(cvcs) >= 2 {
		chain := buildCertChain(cvcs[1], caCerts)
		if chainDER, err := buildDegeneratePKCS7(chain); err == nil {
			result["CoSignerCvcChain"] = strings.ToUpper(hex.EncodeToString(chainDER))
		}
	}

	return result, nil
}

// hasCodeSigningEKU checks if a certificate has the codeSigning extended key usage
// (OID 1.3.6.1.5.5.7.3.3).
func hasCodeSigningEKU(cert *x509.Certificate) bool {
	for _, eku := range cert.ExtKeyUsage {
		if eku == x509.ExtKeyUsageCodeSigning {
			return true
		}
	}
	return false
}

// buildCertChain builds an ordered certificate chain from a leaf cert
// up through available CA certificates by following issuer signatures.
// Returns a slice of DER-encoded certificates starting with the leaf.
// Stops when no issuer is found or a self-signed certificate is reached.
func buildCertChain(leaf *x509.Certificate, caCerts []*x509.Certificate) [][]byte {
	chain := [][]byte{leaf.Raw}
	current := leaf

	// Track certificates already in the chain to prevent infinite loops
	// (e.g., a self-signed root verifying itself).
	seen := map[string]bool{
		current.Subject.CommonName + current.SerialNumber.String(): true,
	}

	for {
		found := false
		for _, ca := range caCerts {
			key := ca.Subject.CommonName + ca.SerialNumber.String()
			if seen[key] {
				continue
			}
			if err := current.CheckSignatureFrom(ca); err == nil {
				chain = append(chain, ca.Raw)
				seen[key] = true
				current = ca
				found = true
				break
			}
		}
		if !found {
			break
		}
	}
	return chain
}

// buildDegeneratePKCS7 creates a degenerate PKCS#7 SignedData structure
// containing only certificates (no content or signatures). This is the
// format used by DOCSIS TLV 81 (ManufacturerCvcChain) and TLV 82 (CoSignerCvcChain).
//
// For a chain of multiple certificates, the DER-encoded certs are concatenated
// before wrapping in the PKCS#7 structure, which is the standard way to encode
// a certificate bag in a degenerate PKCS#7.
func buildDegeneratePKCS7(certDERs [][]byte) ([]byte, error) {
	if len(certDERs) == 0 {
		return nil, fmt.Errorf("no certificates to wrap")
	}

	// Concatenate all DER-encoded certificates into a single byte slice.
	// DegenerateCertificate wraps this in a PKCS#7 SignedData with no signers.
	var combined []byte
	for _, der := range certDERs {
		combined = append(combined, der...)
	}

	return pkcs7.DegenerateCertificate(combined)
}
