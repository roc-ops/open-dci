package mibresolver

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// mibsDir returns the absolute path to the mibs/ directory at the repository
// root. It uses the source file location so tests work from any working
// directory.
func mibsDir(t *testing.T) string {
	t.Helper()

	// This file lives at reference-implementation/mibresolver/resolver_test.go.
	// The mibs/ directory is at ../../mibs relative to this file.
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot determine test file location")
	}
	dir := filepath.Join(filepath.Dir(filename), "..", "..", "mibs")
	abs, err := filepath.Abs(dir)
	if err != nil {
		t.Fatalf("cannot resolve mibs dir: %v", err)
	}
	if _, err := os.Stat(abs); err != nil {
		t.Skipf("mibs directory not found at %s: %v", abs, err)
	}
	return abs
}

func TestNewResolver(t *testing.T) {
	r, err := New(mibsDir(t))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer r.Close()
}

func TestResolveOID_sysDescr(t *testing.T) {
	r, err := New(mibsDir(t))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer r.Close()

	// sysDescr.0 is 1.3.6.1.2.1.1.1.0
	got := r.ResolveOID("1.3.6.1.2.1.1.1.0")
	want := "SNMPv2-MIB::sysDescr.0"
	if got != want {
		t.Errorf("ResolveOID(sysDescr.0) = %q, want %q", got, want)
	}
}

func TestResolveOID_ifAdminStatus(t *testing.T) {
	r, err := New(mibsDir(t))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer r.Close()

	// ifAdminStatus.1 is 1.3.6.1.2.1.2.2.1.7.1
	got := r.ResolveOID("1.3.6.1.2.1.2.2.1.7.1")
	want := "IF-MIB::ifAdminStatus.1"
	if got != want {
		t.Errorf("ResolveOID(ifAdminStatus.1) = %q, want %q", got, want)
	}
}

func TestResolveOID_noIndex(t *testing.T) {
	r, err := New(mibsDir(t))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer r.Close()

	// ifAdminStatus without index is 1.3.6.1.2.1.2.2.1.7
	got := r.ResolveOID("1.3.6.1.2.1.2.2.1.7")
	want := "IF-MIB::ifAdminStatus"
	if got != want {
		t.Errorf("ResolveOID(ifAdminStatus) = %q, want %q", got, want)
	}
}

func TestResolveOID_unknownOID(t *testing.T) {
	r, err := New(mibsDir(t))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer r.Close()

	// An OID that only matches the well-known root (iso = 1) returns empty
	// because the resolver filters out <well-known> module nodes.
	got := r.ResolveOID("1.99.99.99.99")
	if got != "" {
		t.Errorf("ResolveOID(unknown under iso) = %q, want empty string", got)
	}
}

func TestResolveOID_partialResolve(t *testing.T) {
	r, err := New(mibsDir(t))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer r.Close()

	// 1.3.6.1 = internet (from SNMPv2-SMI). An OID like 1.3.6.1.99... will
	// resolve to the deepest known node with the remainder as suffix.
	got := r.ResolveOID("1.3.6.1.99.99.99")
	if got == "" {
		t.Error("ResolveOID(partial) = empty, want partial resolution")
	}
	// It should start with the module name and contain the object name.
	if !strings.Contains(got, "::") {
		t.Errorf("ResolveOID(partial) = %q, want MODULE::name format", got)
	}
}

func TestResolveOID_emptyString(t *testing.T) {
	r, err := New(mibsDir(t))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer r.Close()

	got := r.ResolveOID("")
	if got != "" {
		t.Errorf("ResolveOID(\"\") = %q, want empty string", got)
	}
}

func TestResolveOID_leadingDot(t *testing.T) {
	r, err := New(mibsDir(t))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer r.Close()

	// Leading dot should be handled.
	got := r.ResolveOID(".1.3.6.1.2.1.1.1.0")
	want := "SNMPv2-MIB::sysDescr.0"
	if got != want {
		t.Errorf("ResolveOID(.1.3.6.1.2.1.1.1.0) = %q, want %q", got, want)
	}
}

func TestResolveEnum_ifAdminStatus(t *testing.T) {
	r, err := New(mibsDir(t))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer r.Close()

	// ifAdminStatus is 1.3.6.1.2.1.2.2.1.7
	// value 1 = up, 2 = down, 3 = testing
	tests := []struct {
		value int64
		want  string
	}{
		{1, "up(1)"},
		{2, "down(2)"},
		{3, "testing(3)"},
	}

	for _, tt := range tests {
		got := r.ResolveEnum("1.3.6.1.2.1.2.2.1.7", tt.value)
		if got != tt.want {
			t.Errorf("ResolveEnum(ifAdminStatus, %d) = %q, want %q", tt.value, got, tt.want)
		}
	}
}

func TestResolveEnum_unknownValue(t *testing.T) {
	r, err := New(mibsDir(t))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer r.Close()

	// ifAdminStatus with value 99 should return empty.
	got := r.ResolveEnum("1.3.6.1.2.1.2.2.1.7", 99)
	if got != "" {
		t.Errorf("ResolveEnum(ifAdminStatus, 99) = %q, want empty string", got)
	}
}

func TestResolveEnum_noEnum(t *testing.T) {
	r, err := New(mibsDir(t))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer r.Close()

	// sysDescr (1.3.6.1.2.1.1.1) is an OctetString, not an enum.
	got := r.ResolveEnum("1.3.6.1.2.1.1.1", 1)
	if got != "" {
		t.Errorf("ResolveEnum(sysDescr, 1) = %q, want empty string", got)
	}
}

func TestResolveEnum_emptyOID(t *testing.T) {
	r, err := New(mibsDir(t))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer r.Close()

	got := r.ResolveEnum("", 1)
	if got != "" {
		t.Errorf("ResolveEnum(\"\", 1) = %q, want empty string", got)
	}
}

func TestResolveNamed(t *testing.T) {
	r, err := New(mibsDir(t))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer r.Close()

	oid, err := r.ResolveNamed("SNMPv2-MIB::sysDescr")
	if err != nil {
		t.Fatalf("ResolveNamed() error: %v", err)
	}
	want := "1.3.6.1.2.1.1.1"
	if oid != want {
		t.Errorf("ResolveNamed(SNMPv2-MIB::sysDescr) = %q, want %q", oid, want)
	}
}

func TestResolveNamed_invalidFormat(t *testing.T) {
	r, err := New(mibsDir(t))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer r.Close()

	_, err = r.ResolveNamed("invalid-no-double-colon")
	if err == nil {
		t.Error("ResolveNamed(invalid) expected error, got nil")
	}
}

func TestResolveNamed_unknownObject(t *testing.T) {
	r, err := New(mibsDir(t))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer r.Close()

	_, err = r.ResolveNamed("SNMPv2-MIB::nonExistentObject")
	if err == nil {
		t.Error("ResolveNamed(nonExistent) expected error, got nil")
	}
}

func TestNewResolver_invalidDir(t *testing.T) {
	_, err := New("/nonexistent/path/to/mibs")
	if err == nil {
		t.Error("New(nonexistent) expected error, got nil")
	}
}

func TestNewFromMIBData_coreOnly(t *testing.T) {
	r, err := NewFromMIBData(nil)
	if err != nil {
		t.Fatalf("NewFromMIBData(nil) error: %v", err)
	}
	defer r.Close()

	// Core MIBs should resolve SNMPv2-MIB OIDs.
	got := r.ResolveOID("1.3.6.1.2.1.1.1.0")
	if got == "" {
		t.Error("ResolveOID(sysDescr.0) returned empty with core MIBs")
	}
}

func TestNewFromMIBData_withUserMIBs(t *testing.T) {
	// Read a real MIB file from the repository to use as test data.
	mibsRoot := mibsDir(t)

	// Read IF-MIB from the ietf directory.
	ifMibPath := mibsRoot + "/ietf/IF-MIB.mib"
	ifMibData, err := os.ReadFile(ifMibPath)
	if err != nil {
		// Try following the symlink.
		resolved, resolveErr := os.Readlink(ifMibPath)
		if resolveErr != nil {
			t.Skipf("IF-MIB.mib not found: %v", err)
		}
		ifMibData, err = os.ReadFile(mibsRoot + "/ietf/" + resolved)
		if err != nil {
			t.Skipf("IF-MIB.mib not readable: %v", err)
		}
	}

	// Also need IANAifType-MIB which IF-MIB imports.
	ianaPath := mibsRoot + "/iana/IANAifType-MIB.mib"
	ianaData, err := os.ReadFile(ianaPath)
	if err != nil {
		resolved, resolveErr := os.Readlink(ianaPath)
		if resolveErr != nil {
			t.Skipf("IANAifType-MIB.mib not found: %v", err)
		}
		ianaData, err = os.ReadFile(mibsRoot + "/iana/" + resolved)
		if err != nil {
			t.Skipf("IANAifType-MIB.mib not readable: %v", err)
		}
	}

	// Also need SNMPv2-MIB for sysDescr.
	snmpv2Path := mibsRoot + "/ietf/SNMPv2-MIB.mib"
	snmpv2Data, err := os.ReadFile(snmpv2Path)
	if err != nil {
		resolved, resolveErr := os.Readlink(snmpv2Path)
		if resolveErr != nil {
			t.Skipf("SNMPv2-MIB.mib not found: %v", err)
		}
		snmpv2Data, err = os.ReadFile(mibsRoot + "/ietf/" + resolved)
		if err != nil {
			t.Skipf("SNMPv2-MIB.mib not readable: %v", err)
		}
	}

	r, err := NewFromMIBData(map[string][]byte{
		"IF-MIB.mib":         ifMibData,
		"IANAifType-MIB.mib": ianaData,
		"SNMPv2-MIB.mib":     snmpv2Data,
	})
	if err != nil {
		t.Fatalf("NewFromMIBData error: %v", err)
	}
	defer r.Close()

	// Should resolve IF-MIB OIDs.
	got := r.ResolveOID("1.3.6.1.2.1.2.2.1.7.1")
	want := "IF-MIB::ifAdminStatus.1"
	if got != want {
		t.Errorf("ResolveOID(ifAdminStatus.1) = %q, want %q", got, want)
	}

	// Should resolve enums.
	enumGot := r.ResolveEnum("1.3.6.1.2.1.2.2.1.7", 1)
	enumWant := "up(1)"
	if enumGot != enumWant {
		t.Errorf("ResolveEnum(ifAdminStatus, 1) = %q, want %q", enumGot, enumWant)
	}
}

func TestLoadAdditionalMIBs(t *testing.T) {
	// Start with core MIBs only.
	r, err := NewFromMIBData(nil)
	if err != nil {
		t.Fatalf("NewFromMIBData(nil) error: %v", err)
	}
	defer r.Close()

	// ifAdminStatus should NOT resolve with only core MIBs.
	got := r.ResolveOID("1.3.6.1.2.1.2.2.1.7.1")
	if strings.Contains(got, "ifAdminStatus") {
		t.Skipf("ifAdminStatus already resolved with core MIBs only — test invalid")
	}

	// Load IF-MIB incrementally.
	mibsRoot := mibsDir(t)

	ifMibData, err := os.ReadFile(mibsRoot + "/ietf/IF-MIB.mib")
	if err != nil {
		t.Skipf("IF-MIB.mib not readable: %v", err)
	}
	ianaData, err := os.ReadFile(mibsRoot + "/iana/IANAifType-MIB.mib")
	if err != nil {
		t.Skipf("IANAifType-MIB.mib not readable: %v", err)
	}
	snmpv2Data, err := os.ReadFile(mibsRoot + "/ietf/SNMPv2-MIB.mib")
	if err != nil {
		t.Skipf("SNMPv2-MIB.mib not readable: %v", err)
	}

	loaded, err := r.LoadAdditionalMIBs(map[string][]byte{
		"IF-MIB.mib":         ifMibData,
		"IANAifType-MIB.mib": ianaData,
		"SNMPv2-MIB.mib":     snmpv2Data,
	})
	if err != nil {
		t.Fatalf("LoadAdditionalMIBs error: %v", err)
	}
	t.Logf("Loaded %d additional modules", loaded)

	// Now ifAdminStatus should resolve.
	got = r.ResolveOID("1.3.6.1.2.1.2.2.1.7.1")
	want := "IF-MIB::ifAdminStatus.1"
	if got != want {
		t.Errorf("ResolveOID(ifAdminStatus.1) after incremental load = %q, want %q", got, want)
	}
}

func TestResolveFullName_ifAdminStatus(t *testing.T) {
	r, err := New(mibsDir(t))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer r.Close()

	got, err := r.ResolveFullName("1.3.6.1.2.1.2.2.1.7")
	if err != nil {
		t.Fatalf("ResolveFullName() error: %v", err)
	}

	// Should be a dotted named path ending in ifAdminStatus.
	if !strings.HasSuffix(got, "ifAdminStatus") {
		t.Errorf("ResolveFullName(1.3.6.1.2.1.2.2.1.7) = %q, want suffix 'ifAdminStatus'", got)
	}

	// Should contain known intermediate names.
	for _, want := range []string{"iso", "org", "dod", "internet"} {
		if !strings.Contains(got, want) {
			t.Errorf("ResolveFullName result %q missing expected component %q", got, want)
		}
	}
	t.Logf("ResolveFullName(1.3.6.1.2.1.2.2.1.7) = %q", got)
}

func TestResolveFullName_sysDescr(t *testing.T) {
	r, err := New(mibsDir(t))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer r.Close()

	got, err := r.ResolveFullName("1.3.6.1.2.1.1.1")
	if err != nil {
		t.Fatalf("ResolveFullName() error: %v", err)
	}

	if !strings.HasSuffix(got, "sysDescr") {
		t.Errorf("ResolveFullName(1.3.6.1.2.1.1.1) = %q, want suffix 'sysDescr'", got)
	}
	t.Logf("ResolveFullName(1.3.6.1.2.1.1.1) = %q", got)
}

func TestResolveFullName_leadingDot(t *testing.T) {
	r, err := New(mibsDir(t))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer r.Close()

	got, err := r.ResolveFullName(".1.3.6.1.2.1.1.1")
	if err != nil {
		t.Fatalf("ResolveFullName() error: %v", err)
	}

	if !strings.HasSuffix(got, "sysDescr") {
		t.Errorf("ResolveFullName with leading dot = %q, want suffix 'sysDescr'", got)
	}
}

func TestResolveFullName_empty(t *testing.T) {
	r, err := New(mibsDir(t))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer r.Close()

	_, err = r.ResolveFullName("")
	if err == nil {
		t.Error("ResolveFullName(\"\") expected error, got nil")
	}
}

func TestResolveToNumericOID_plainName(t *testing.T) {
	r, err := New(mibsDir(t))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer r.Close()

	got, err := r.ResolveToNumericOID("ifAdminStatus")
	if err != nil {
		t.Fatalf("ResolveToNumericOID() error: %v", err)
	}
	want := "1.3.6.1.2.1.2.2.1.7"
	if got != want {
		t.Errorf("ResolveToNumericOID(ifAdminStatus) = %q, want %q", got, want)
	}
}

func TestResolveToNumericOID_qualified(t *testing.T) {
	r, err := New(mibsDir(t))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer r.Close()

	got, err := r.ResolveToNumericOID("IF-MIB::ifAdminStatus")
	if err != nil {
		t.Fatalf("ResolveToNumericOID() error: %v", err)
	}
	want := "1.3.6.1.2.1.2.2.1.7"
	if got != want {
		t.Errorf("ResolveToNumericOID(IF-MIB::ifAdminStatus) = %q, want %q", got, want)
	}
}

func TestResolveToNumericOID_dottedPath(t *testing.T) {
	r, err := New(mibsDir(t))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer r.Close()

	// Dotted named path — uses leaf name for lookup.
	got, err := r.ResolveToNumericOID("iso.org.dod.internet.mgmt.mib-2.interfaces.ifTable.ifEntry.ifAdminStatus")
	if err != nil {
		t.Fatalf("ResolveToNumericOID(dotted path) error: %v", err)
	}
	want := "1.3.6.1.2.1.2.2.1.7"
	if got != want {
		t.Errorf("ResolveToNumericOID(dotted path) = %q, want %q", got, want)
	}
}

func TestResolveToNumericOID_unknown(t *testing.T) {
	r, err := New(mibsDir(t))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer r.Close()

	_, err = r.ResolveToNumericOID("nonExistentObject12345")
	if err == nil {
		t.Error("ResolveToNumericOID(nonExistent) expected error, got nil")
	}
}

func TestResolveToNumericOID_empty(t *testing.T) {
	r, err := New(mibsDir(t))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer r.Close()

	_, err = r.ResolveToNumericOID("")
	if err == nil {
		t.Error("ResolveToNumericOID(\"\") expected error, got nil")
	}
}

func TestResolveRoundTrip(t *testing.T) {
	r, err := New(mibsDir(t))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer r.Close()

	// numeric → named → numeric round-trip
	numericOID := "1.3.6.1.2.1.2.2.1.7"
	namedPath, err := r.ResolveFullName(numericOID)
	if err != nil {
		t.Fatalf("ResolveFullName() error: %v", err)
	}

	gotOID, err := r.ResolveToNumericOID(namedPath)
	if err != nil {
		t.Fatalf("ResolveToNumericOID() error: %v", err)
	}

	if gotOID != numericOID {
		t.Errorf("round-trip: %q → %q → %q, want %q", numericOID, namedPath, gotOID, numericOID)
	}
}

func TestWithVersionOverrides(t *testing.T) {
	// Test that the option function correctly parses overrides.
	cfg := &config{versionOverrides: make(map[string]string)}
	opt := WithVersionOverrides([]string{"DOCS-IF3-MIB@2024-07-05", "DOCS-QOS3-MIB@2023-11-22"})
	opt(cfg)

	if v, ok := cfg.versionOverrides["DOCS-IF3-MIB"]; !ok || v != "DOCS-IF3-MIB@2024-07-05" {
		t.Errorf("versionOverrides[DOCS-IF3-MIB] = %q, want %q", v, "DOCS-IF3-MIB@2024-07-05")
	}
	if v, ok := cfg.versionOverrides["DOCS-QOS3-MIB"]; !ok || v != "DOCS-QOS3-MIB@2023-11-22" {
		t.Errorf("versionOverrides[DOCS-QOS3-MIB] = %q, want %q", v, "DOCS-QOS3-MIB@2023-11-22")
	}
}
