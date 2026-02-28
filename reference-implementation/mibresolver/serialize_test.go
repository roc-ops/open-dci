package mibresolver

import (
	"bytes"
	"encoding/gob"
	"os"
	"testing"
)

// loadTestMIBFiles reads a small set of real MIB files from the repository
// for use in serialization tests. Returns a map of filename -> content.
func loadTestMIBFiles(t *testing.T) map[string][]byte {
	t.Helper()
	mibsRoot := mibsDir(t)

	files := map[string]string{
		"IF-MIB.mib":         mibsRoot + "/ietf/IF-MIB.mib",
		"IANAifType-MIB.mib": mibsRoot + "/iana/IANAifType-MIB.mib",
		"SNMPv2-MIB.mib":     mibsRoot + "/ietf/SNMPv2-MIB.mib",
	}

	result := make(map[string][]byte, len(files))
	for name, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Skipf("%s not readable: %v", name, err)
		}
		result[name] = data
	}
	return result
}

func TestSerializeMIBState_RoundTrip(t *testing.T) {
	mibFiles := loadTestMIBFiles(t)

	// Create initial resolver and verify it works.
	r1, err := NewFromMIBData(mibFiles)
	if err != nil {
		t.Fatalf("NewFromMIBData error: %v", err)
	}

	// Verify OID resolution works before serialization.
	got1 := r1.ResolveOID("1.3.6.1.2.1.2.2.1.7.1")
	want := "IF-MIB::ifAdminStatus.1"
	if got1 != want {
		t.Fatalf("before serialize: ResolveOID = %q, want %q", got1, want)
	}

	enum1 := r1.ResolveEnum("1.3.6.1.2.1.2.2.1.7", 1)
	wantEnum := "up(1)"
	if enum1 != wantEnum {
		t.Fatalf("before serialize: ResolveEnum = %q, want %q", enum1, wantEnum)
	}

	r1.Close()

	// Serialize the MIB file contents.
	blob, err := SerializeMIBState(mibFiles)
	if err != nil {
		t.Fatalf("SerializeMIBState error: %v", err)
	}

	if len(blob) == 0 {
		t.Fatal("SerializeMIBState returned empty blob")
	}
	t.Logf("serialized blob size: %d bytes", len(blob))

	// Restore from the serialized blob.
	r2, restoredFiles, err := RestoreMIBState(blob)
	if err != nil {
		t.Fatalf("RestoreMIBState error: %v", err)
	}
	defer r2.Close()

	// Verify the restored file map matches.
	if len(restoredFiles) != len(mibFiles) {
		t.Errorf("restored %d files, want %d", len(restoredFiles), len(mibFiles))
	}
	for name := range mibFiles {
		if _, ok := restoredFiles[name]; !ok {
			t.Errorf("restored files missing %q", name)
		}
	}

	// Verify OID resolution works identically after restore.
	got2 := r2.ResolveOID("1.3.6.1.2.1.2.2.1.7.1")
	if got2 != want {
		t.Errorf("after restore: ResolveOID = %q, want %q", got2, want)
	}

	enum2 := r2.ResolveEnum("1.3.6.1.2.1.2.2.1.7", 1)
	if enum2 != wantEnum {
		t.Errorf("after restore: ResolveEnum = %q, want %q", enum2, wantEnum)
	}

	// Verify named resolution works.
	oid, err := r2.ResolveNamed("IF-MIB::ifAdminStatus")
	if err != nil {
		t.Errorf("after restore: ResolveNamed error: %v", err)
	}
	if oid != "1.3.6.1.2.1.2.2.1.7" {
		t.Errorf("after restore: ResolveNamed = %q, want %q", oid, "1.3.6.1.2.1.2.2.1.7")
	}

	// Verify full name resolution.
	fullName, err := r2.ResolveFullName("1.3.6.1.2.1.2.2.1.7")
	if err != nil {
		t.Errorf("after restore: ResolveFullName error: %v", err)
	}
	if fullName == "" {
		t.Error("after restore: ResolveFullName returned empty")
	}

	// Verify tree query works.
	tree, err := r2.QueryTree()
	if err != nil {
		t.Errorf("after restore: QueryTree error: %v", err)
	}
	if tree == nil || tree.OID != "1" {
		t.Error("after restore: QueryTree returned nil or wrong root")
	}
}

func TestSerializeMIBState_EmptyFiles(t *testing.T) {
	_, err := SerializeMIBState(nil)
	if err == nil {
		t.Error("SerializeMIBState(nil) expected error, got nil")
	}

	_, err = SerializeMIBState(map[string][]byte{})
	if err == nil {
		t.Error("SerializeMIBState(empty) expected error, got nil")
	}
}

func TestRestoreMIBState_EmptyBlob(t *testing.T) {
	_, _, err := RestoreMIBState(nil)
	if err == nil {
		t.Error("RestoreMIBState(nil) expected error, got nil")
	}

	_, _, err = RestoreMIBState([]byte{})
	if err == nil {
		t.Error("RestoreMIBState(empty) expected error, got nil")
	}
}

func TestRestoreMIBState_CorruptBlob(t *testing.T) {
	_, _, err := RestoreMIBState([]byte("this is not a valid gob blob"))
	if err == nil {
		t.Error("RestoreMIBState(corrupt) expected error, got nil")
	}
}

func TestRestoreMIBState_BadMagic(t *testing.T) {
	mibFiles := loadTestMIBFiles(t)

	// Serialize normally.
	blob, err := SerializeMIBState(mibFiles)
	if err != nil {
		t.Fatalf("SerializeMIBState error: %v", err)
	}

	// Corrupt the blob by flipping bytes (gob encodes the struct fields,
	// so we need to create a blob with wrong magic via direct encoding).
	snapshot := MIBSnapshot{
		Magic:    0xDEADBEEF, // wrong magic
		Version:  MIBSnapshotVersion,
		MIBFiles: mibFiles,
	}

	_ = blob // original blob unused, we create a bad one

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(snapshot); err != nil {
		t.Fatalf("encoding bad snapshot: %v", err)
	}

	_, _, err = RestoreMIBState(buf.Bytes())
	if err == nil {
		t.Error("RestoreMIBState(bad magic) expected error, got nil")
	}
	t.Logf("bad magic error: %v", err)
}

func TestRestoreMIBState_BadVersion(t *testing.T) {
	mibFiles := loadTestMIBFiles(t)

	snapshot := MIBSnapshot{
		Magic:    MIBSnapshotMagic,
		Version:  999, // future version
		MIBFiles: mibFiles,
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(snapshot); err != nil {
		t.Fatalf("encoding bad version snapshot: %v", err)
	}

	_, _, err := RestoreMIBState(buf.Bytes())
	if err == nil {
		t.Error("RestoreMIBState(bad version) expected error, got nil")
	}
	t.Logf("bad version error: %v", err)
}

func TestRestoreMIBState_AugmentAfterRestore(t *testing.T) {
	mibFiles := loadTestMIBFiles(t)

	// Serialize just the core set of MIBs.
	blob, err := SerializeMIBState(mibFiles)
	if err != nil {
		t.Fatalf("SerializeMIBState error: %v", err)
	}

	// Restore the state.
	r, _, err := RestoreMIBState(blob)
	if err != nil {
		t.Fatalf("RestoreMIBState error: %v", err)
	}
	defer r.Close()

	// Verify initial state works.
	got := r.ResolveOID("1.3.6.1.2.1.2.2.1.7.1")
	if got != "IF-MIB::ifAdminStatus.1" {
		t.Errorf("after restore: ResolveOID = %q, want %q", got, "IF-MIB::ifAdminStatus.1")
	}

	// Load additional MIBs on top of the restored state (augment mode).
	// Use the same files as a simple test — the key point is that
	// LoadAdditionalMIBs doesn't error out after a restore.
	mibsRoot := mibsDir(t)
	tcpMibPath := mibsRoot + "/ietf/TCP-MIB.mib"
	tcpData, err := os.ReadFile(tcpMibPath)
	if err != nil {
		t.Skipf("TCP-MIB.mib not readable: %v", err)
	}

	loaded, err := r.LoadAdditionalMIBs(map[string][]byte{
		"TCP-MIB.mib": tcpData,
	})
	if err != nil {
		t.Fatalf("LoadAdditionalMIBs after restore error: %v", err)
	}
	t.Logf("augment after restore: loaded %d additional modules", loaded)

	// The original OID should still resolve.
	got = r.ResolveOID("1.3.6.1.2.1.2.2.1.7.1")
	if got != "IF-MIB::ifAdminStatus.1" {
		t.Errorf("after augment: ResolveOID = %q, want %q", got, "IF-MIB::ifAdminStatus.1")
	}
}
