package mibresolver

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

// MIBSnapshotMagic is a magic number identifying a serialized MIB snapshot blob.
// Chosen to be unlikely to collide with other binary formats.
const MIBSnapshotMagic uint32 = 0x4D494253 // "MIBS" in ASCII

// MIBSnapshotVersion is the current serialization format version.
// Increment this when the snapshot format changes incompatibly.
const MIBSnapshotVersion uint32 = 1

// MIBSnapshot is the serializable representation of loaded MIB state.
// It captures the raw MIB file contents so the resolver can be fully
// reconstructed from a snapshot without access to the original files.
type MIBSnapshot struct {
	Magic    uint32
	Version  uint32
	MIBFiles map[string][]byte // filename -> file content
}

// SerializeMIBState encodes the given MIB file contents into a binary blob
// that can later be passed to RestoreMIBState to reconstruct the resolver.
func SerializeMIBState(mibFiles map[string][]byte) ([]byte, error) {
	if len(mibFiles) == 0 {
		return nil, fmt.Errorf("no MIB files to serialize")
	}

	snapshot := MIBSnapshot{
		Magic:    MIBSnapshotMagic,
		Version:  MIBSnapshotVersion,
		MIBFiles: mibFiles,
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(snapshot); err != nil {
		return nil, fmt.Errorf("encoding MIB snapshot: %w", err)
	}

	return buf.Bytes(), nil
}

// RestoreMIBState decodes a binary blob produced by SerializeMIBState and
// creates a new Resolver initialized with the contained MIB files. The
// caller must Close() any existing resolver before calling this function,
// as gosmi uses global state.
func RestoreMIBState(data []byte) (*Resolver, map[string][]byte, error) {
	if len(data) == 0 {
		return nil, nil, fmt.Errorf("empty snapshot data")
	}

	var snapshot MIBSnapshot
	dec := gob.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(&snapshot); err != nil {
		return nil, nil, fmt.Errorf("corrupt MIB snapshot: %w", err)
	}

	if snapshot.Magic != MIBSnapshotMagic {
		return nil, nil, fmt.Errorf("invalid MIB snapshot: bad magic number 0x%08X (expected 0x%08X)",
			snapshot.Magic, MIBSnapshotMagic)
	}

	if snapshot.Version != MIBSnapshotVersion {
		return nil, nil, fmt.Errorf("incompatible MIB snapshot version %d (expected %d)",
			snapshot.Version, MIBSnapshotVersion)
	}

	if len(snapshot.MIBFiles) == 0 {
		return nil, nil, fmt.Errorf("MIB snapshot contains no MIB files")
	}

	// Reconstruct the resolver from the stored MIB file contents.
	r, err := NewFromMIBData(snapshot.MIBFiles)
	if err != nil {
		return nil, nil, fmt.Errorf("restoring MIB state: %w", err)
	}

	return r, snapshot.MIBFiles, nil
}
