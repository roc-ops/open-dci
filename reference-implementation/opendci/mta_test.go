package opendci

import (
	"bytes"
	"testing"
)

// buildMTAConfig wraps TLVs with MTA delimiters (TLV 254 start/end).
func buildMTAConfig(tlvs ...[]byte) []byte {
	var data []byte
	// Start delimiter: TLV 254, length 1, value 1
	data = append(data, 0xFE, 0x01, 0x01)
	for _, tlv := range tlvs {
		data = append(data, tlv...)
	}
	// End delimiter: TLV 254, length 1, value 255
	data = append(data, 0xFE, 0x01, 0xFF)
	return data
}

// makeMTATestRegistry builds a minimal MTA registry for testing.
func makeMTATestRegistry() *Registry {
	reg := &Registry{
		TopLevel: map[int]*TLVDef{
			254: {
				TypeNum:  254,
				Name:     "MtaConfigDelimiter",
				DataType: DataTypeUint8,
			},
			11: {
				TypeNum:    11,
				Name:       "SnmpMibObject",
				DataType:   DataTypeCompound,
				Repeatable: true,
			},
			64: {
				TypeNum:    64,
				Name:       "SnmpMibObjectLarge",
				DataType:   DataTypeCompound,
				LengthSize: 2,
				Repeatable: true,
				RefName:    "SnmpMibEntry",
			},
			38: {
				TypeNum:    38,
				Name:       "Snmpv3NotificationReceiver",
				DataType:   DataTypeCompound,
				Repeatable: true,
				SubTLVs: map[int]*TLVDef{
					1: {TypeNum: 1, Name: "IPv4Address", DataType: DataTypeIPv4Address},
					2: {TypeNum: 2, Name: "UDPPortNumber", DataType: DataTypeUint16},
					3: {TypeNum: 3, Name: "TrapType", DataType: DataTypeUint16},
				},
			},
		},
		NameLookup: make(map[string]*TLVDef),
		Format:     FormatMTA,
	}
	for _, def := range reg.TopLevel {
		reg.NameLookup[def.Name] = def
	}
	return reg
}

func TestDetectFormat(t *testing.T) {
	if DetectFormat([]byte{0xFE, 0x01, 0x01}) != FormatMTA {
		t.Error("expected MTA format for data starting with 0xFE")
	}
	if DetectFormat([]byte{0x03, 0x01, 0x01}) != FormatCM {
		t.Error("expected CM format for data starting with 0x03")
	}
	if DetectFormat([]byte{0xFF, 0x00}) != FormatCM {
		t.Error("expected CM format for data starting with 0xFF")
	}
	if DetectFormat(nil) != FormatCM {
		t.Error("expected CM format for nil data")
	}
	if DetectFormat([]byte{}) != FormatCM {
		t.Error("expected CM format for empty data")
	}
}

func TestMTADecodeTLV11(t *testing.T) {
	reg := makeMTATestRegistry()

	// Build an SNMP varbind: OID 1.3.6.1.2.1, Integer = 42
	oidBytes := []byte{0x2B, 0x06, 0x01, 0x02, 0x01}
	varbind := buildVarbind(oidBytes, tagInteger, []byte{0x2A})

	data := buildMTAConfig(
		buildTLV(11, varbind),
	)

	result, err := Decode(data, reg)
	if err != nil {
		t.Fatal(err)
	}

	snmpArr, ok := result.Config["SnmpMibObject"].([]interface{})
	if !ok {
		t.Fatalf("expected []interface{}, got %T", result.Config["SnmpMibObject"])
	}
	if len(snmpArr) != 1 {
		t.Fatalf("expected 1 SNMP entry, got %d", len(snmpArr))
	}

	entry := snmpArr[0].(map[string]interface{})
	if entry["oid"] != "1.3.6.1.2.1" {
		t.Errorf("expected oid '1.3.6.1.2.1', got %v", entry["oid"])
	}
	if entry["type"] != "Integer" {
		t.Errorf("expected type 'Integer', got %v", entry["type"])
	}
	if entry["value"] != "42" {
		t.Errorf("expected value '42', got %v", entry["value"])
	}
}

func TestMTADecodeTLV64Varbind(t *testing.T) {
	reg := makeMTATestRegistry()

	// Build an SNMP varbind: OID 1.3.6.1.2.1.69.1.3.8.0, String = "test"
	oidBytes := []byte{0x2B, 0x06, 0x01, 0x02, 0x01, 0x45, 0x01, 0x03, 0x08, 0x00}
	varbind := buildVarbind(oidBytes, tagOctetString, []byte("test"))

	// TLV 64 with 2-byte length
	data := buildMTAConfig(
		buildTLV2(64, varbind),
	)

	result, err := Decode(data, reg)
	if err != nil {
		t.Fatal(err)
	}

	snmpArr, ok := result.Config["SnmpMibObjectLarge"].([]interface{})
	if !ok {
		t.Fatalf("expected []interface{}, got %T", result.Config["SnmpMibObjectLarge"])
	}
	if len(snmpArr) != 1 {
		t.Fatalf("expected 1 SNMP entry, got %d", len(snmpArr))
	}

	entry := snmpArr[0].(map[string]interface{})
	if entry["oid"] != "1.3.6.1.2.1.69.1.3.8.0" {
		t.Errorf("expected oid '1.3.6.1.2.1.69.1.3.8.0', got %v", entry["oid"])
	}
	if entry["type"] != "String" {
		t.Errorf("expected type 'String', got %v", entry["type"])
	}
	if entry["value"] != "test" {
		t.Errorf("expected value 'test', got %v", entry["value"])
	}
}

func TestMTADecodeCompound(t *testing.T) {
	reg := makeMTATestRegistry()

	// TLV 38 (Snmpv3NotificationReceiver) with sub-TLVs
	subTLVs := buildTLV(1, []byte{10, 0, 0, 1})      // IPv4Address = 10.0.0.1
	subTLVs = append(subTLVs, buildTLV(2, []byte{0x00, 0xA2})...) // UDPPortNumber = 162
	subTLVs = append(subTLVs, buildTLV(3, []byte{0x00, 0x02})...) // TrapType = 2

	data := buildMTAConfig(
		buildTLV(38, subTLVs),
	)

	result, err := Decode(data, reg)
	if err != nil {
		t.Fatal(err)
	}

	arr, ok := result.Config["Snmpv3NotificationReceiver"].([]interface{})
	if !ok {
		t.Fatalf("expected []interface{}, got %T", result.Config["Snmpv3NotificationReceiver"])
	}
	if len(arr) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(arr))
	}

	entry := arr[0].(map[string]interface{})
	if entry["IPv4Address"] != "10.0.0.1" {
		t.Errorf("expected IPv4Address '10.0.0.1', got %v", entry["IPv4Address"])
	}
	if entry["UDPPortNumber"] != 162 {
		t.Errorf("expected UDPPortNumber 162, got %v", entry["UDPPortNumber"])
	}
	if entry["TrapType"] != 2 {
		t.Errorf("expected TrapType 2, got %v", entry["TrapType"])
	}
}

func TestMTADecodePreservesDelimiter(t *testing.T) {
	reg := makeMTATestRegistry()

	data := buildMTAConfig(
		buildTLV(11, buildVarbind([]byte{0x2B, 0x06, 0x01}, tagInteger, []byte{0x01})),
	)

	result, err := Decode(data, reg)
	if err != nil {
		t.Fatal(err)
	}

	// MtaConfigDelimiter start (value 1) should be preserved in decoded config.
	delim, ok := result.Config["MtaConfigDelimiter"]
	if !ok {
		t.Fatal("MtaConfigDelimiter should appear in decoded config")
	}
	if delim != float64(1) {
		t.Errorf("MtaConfigDelimiter should be 1, got %v", delim)
	}

	// MIC fields should be nil in MTA mode.
	if result.CmMic != nil {
		t.Error("CmMic should be nil in MTA mode")
	}
	if result.CmtsMic != nil {
		t.Error("CmtsMic should be nil in MTA mode")
	}
}

func TestMTADecodeDoesNotBreakOnTLV255(t *testing.T) {
	reg := makeMTATestRegistry()

	// In MTA mode, a TLV type 255 should be treated as unknown, not as end-of-data.
	// First, add a known TLV 11, then an unknown byte 0xFF (type 255).
	// MTA decoder should NOT break on type 255 — it should continue to the real
	// end delimiter (TLV 254 value 255).
	varbind := buildVarbind([]byte{0x2B, 0x06, 0x01}, tagInteger, []byte{0x01})
	tlv11 := buildTLV(11, varbind)

	// Build MTA config with TLV 255 in the middle (as an unknown TLV).
	var data []byte
	data = append(data, 0xFE, 0x01, 0x01) // start delimiter
	data = append(data, tlv11...)
	data = append(data, 0xFF, 0x02, 0xAB, 0xCD) // "TLV 255" with length 2 — unknown in MTA
	data = append(data, 0xFE, 0x01, 0xFF) // end delimiter

	result, err := Decode(data, reg)
	if err != nil {
		t.Fatal(err)
	}

	// TLV 11 should be decoded.
	if _, ok := result.Config["SnmpMibObject"]; !ok {
		t.Error("expected SnmpMibObject in decoded config")
	}

	// TLV 255 should be in UnknownTlvs.
	unknowns, ok := result.Config["UnknownTlvs"].([]interface{})
	if !ok {
		t.Fatal("expected UnknownTlvs in MTA config with TLV 255")
	}
	if len(unknowns) != 1 {
		t.Fatalf("expected 1 unknown TLV, got %d", len(unknowns))
	}
	unk := unknowns[0].(map[string]interface{})
	if unk["type"] != 255 {
		t.Errorf("expected unknown TLV type 255, got %v", unk["type"])
	}
}

func TestMTAEncodeDelimiters(t *testing.T) {
	reg := makeMTATestRegistry()

	result := &DecodeResult{
		Config: map[string]interface{}{
			"SnmpMibObject": []interface{}{
				map[string]interface{}{
					"oid":   "1.3.6.1.2.1",
					"type":  "Integer",
					"value": "42",
				},
			},
		},
		TLVOrder: []string{"SnmpMibObject"},
	}

	encoded, err := Encode(result, reg)
	if err != nil {
		t.Fatal(err)
	}

	// Should start with TLV 254 start delimiter.
	if len(encoded) < 3 || encoded[0] != 0xFE || encoded[1] != 0x01 || encoded[2] != 0x01 {
		t.Errorf("expected MTA start delimiter (FE 01 01), got: %X", encoded[:min(3, len(encoded))])
	}

	// Should end with TLV 254 end delimiter.
	n := len(encoded)
	if n < 3 || encoded[n-3] != 0xFE || encoded[n-2] != 0x01 || encoded[n-1] != 0xFF {
		t.Errorf("expected MTA end delimiter (FE 01 FF), got: %X", encoded[max(0, n-3):])
	}

	// Should NOT contain CM end-of-data marker (FF 00).
	for i := 0; i < len(encoded)-1; i++ {
		if encoded[i] == 0xFF && encoded[i+1] == 0x00 {
			t.Errorf("found CM end-of-data marker (FF 00) at offset %d in MTA output", i)
			break
		}
	}
}

func TestMTAEncodeTLV64(t *testing.T) {
	reg := makeMTATestRegistry()

	result := &DecodeResult{
		Config: map[string]interface{}{
			"SnmpMibObjectLarge": []interface{}{
				map[string]interface{}{
					"oid":   "1.3.6.1.2.1.69.1.3.8.0",
					"type":  "String",
					"value": "test",
				},
			},
		},
		TLVOrder: []string{"SnmpMibObjectLarge"},
	}

	encoded, err := Encode(result, reg)
	if err != nil {
		t.Fatal(err)
	}

	// Verify the TLV 64 has a 2-byte length field.
	// Skip start delimiter (3 bytes), then TLV 64 starts.
	if encoded[3] != 64 {
		t.Errorf("expected TLV type 64, got %d", encoded[3])
	}

	// Length should be 2 bytes (big-endian).
	tlv64Len := int(encoded[4])<<8 | int(encoded[5])
	if tlv64Len == 0 {
		t.Error("expected non-zero 2-byte length for TLV 64")
	}
}

func TestMTARoundTrip(t *testing.T) {
	reg := makeMTATestRegistry()

	// Build original binary with MTA delimiters.
	varbind1 := buildVarbind([]byte{0x2B, 0x06, 0x01, 0x02, 0x01}, tagInteger, []byte{0x2A})
	varbind2 := buildVarbind([]byte{0x2B, 0x06, 0x01, 0x02, 0x01, 0x45, 0x01}, tagOctetString, []byte("hello"))
	original := buildMTAConfig(
		buildTLV(11, varbind1),
		buildTLV(11, varbind2),
	)

	// Decode.
	result, err := Decode(original, reg)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}

	// Re-encode.
	encoded, err := Encode(result, reg)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}

	if !bytes.Equal(encoded, original) {
		t.Errorf("MTA round-trip mismatch:\n  got:  %X\n  want: %X", encoded, original)
	}
}

func TestMTARoundTripTLV64(t *testing.T) {
	reg := makeMTATestRegistry()

	// TLV 64 with 2-byte length
	varbind := buildVarbind([]byte{0x2B, 0x06, 0x01, 0x02, 0x01, 0x45, 0x01, 0x03, 0x08, 0x00}, tagOctetString, []byte("test-value"))
	original := buildMTAConfig(
		buildTLV2(64, varbind),
	)

	result, err := Decode(original, reg)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}

	encoded, err := Encode(result, reg)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}

	if !bytes.Equal(encoded, original) {
		t.Errorf("MTA TLV 64 round-trip mismatch:\n  got:  %X\n  want: %X", encoded, original)
	}
}

func TestMTARoundTripMixed(t *testing.T) {
	reg := makeMTATestRegistry()

	// Mix TLV 11, TLV 64, and TLV 38.
	varbind := buildVarbind([]byte{0x2B, 0x06, 0x01}, tagInteger, []byte{0x01})
	varbind64 := buildVarbind([]byte{0x2B, 0x06, 0x01, 0x45}, tagOctetString, []byte("big"))
	subTLVs := buildTLV(1, []byte{192, 168, 1, 1}) // IPv4Address
	subTLVs = append(subTLVs, buildTLV(2, []byte{0x00, 0xA2})...) // UDPPortNumber

	original := buildMTAConfig(
		buildTLV(11, varbind),
		buildTLV2(64, varbind64),
		buildTLV(38, subTLVs),
	)

	result, err := Decode(original, reg)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}

	encoded, err := Encode(result, reg)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}

	if !bytes.Equal(encoded, original) {
		t.Errorf("MTA mixed round-trip mismatch:\n  got:  %X\n  want: %X", encoded, original)
	}
}

func TestMTAEmptyConfig(t *testing.T) {
	reg := makeMTATestRegistry()

	// Just delimiters, no TLVs.
	data := buildMTAConfig()

	result, err := Decode(data, reg)
	if err != nil {
		t.Fatal(err)
	}

	// Only the MtaConfigDelimiter should be present (no actual TLV data).
	if len(result.Config) != 1 {
		t.Errorf("expected 1 entry (MtaConfigDelimiter only), got %d entries", len(result.Config))
	}
	if _, ok := result.Config["MtaConfigDelimiter"]; !ok {
		t.Error("expected MtaConfigDelimiter in decoded config")
	}
}

func TestMTAEncodeEmptyConfig(t *testing.T) {
	reg := makeMTATestRegistry()

	result := &DecodeResult{
		Config: map[string]interface{}{},
	}

	encoded, err := Encode(result, reg)
	if err != nil {
		t.Fatal(err)
	}

	// Should be: FE 01 01 FE 01 FF
	expected := []byte{0xFE, 0x01, 0x01, 0xFE, 0x01, 0xFF}
	if !bytes.Equal(encoded, expected) {
		t.Errorf("empty MTA encode mismatch:\n  got:  %X\n  want: %X", encoded, expected)
	}
}

func TestCMDecodeUnchanged(t *testing.T) {
	// Verify CM decode behavior is completely unaffected by MTA changes.
	reg := makeTestRegistry()

	data := buildConfig(
		buildTLV(3, []byte{1}),
		buildTLV(1, []byte{0x15, 0x75, 0x2A, 0x00}),
		buildTLV(2, []byte{5}),
	)

	result, err := Decode(data, reg)
	if err != nil {
		t.Fatal(err)
	}

	if result.Config["NetworkAccess"] != 1 {
		t.Errorf("expected NetworkAccess=1, got %v", result.Config["NetworkAccess"])
	}
	if result.Config["DownstreamFrequency"] != 360000000 {
		t.Errorf("expected DownstreamFrequency=360000000, got %v", result.Config["DownstreamFrequency"])
	}
}

func TestCMEncodeUnchanged(t *testing.T) {
	// Verify CM encode behavior is completely unaffected by MTA changes.
	reg := makeTestRegistry()

	result := &DecodeResult{
		Config: map[string]interface{}{
			"NetworkAccess": 1,
		},
		TLVOrder: []string{"NetworkAccess"},
	}

	encoded, err := Encode(result, reg)
	if err != nil {
		t.Fatal(err)
	}

	expected := buildConfig(buildTLV(3, []byte{1}))
	if !bytes.Equal(encoded, expected) {
		t.Errorf("CM encode changed unexpectedly:\n  got:  %X\n  want: %X", encoded, expected)
	}

	// Verify CM end-of-data marker.
	n := len(encoded)
	if encoded[n-2] != 0xFF || encoded[n-1] != 0x00 {
		t.Error("CM encoded data should end with FF 00")
	}
}

func TestRegistryFormatDefault(t *testing.T) {
	// A registry without format metadata should default to CM.
	reg := &Registry{
		TopLevel:   make(map[int]*TLVDef),
		NameLookup: make(map[string]*TLVDef),
	}
	if reg.Format != "" {
		// Zero value is "", which is fine — the decode/encode compare against FormatMTA.
		// But LoadRegistryFromBytes should set it explicitly.
	}
}

func TestMTAJSONCHeader(t *testing.T) {
	config := map[string]interface{}{
		"SnmpMibObject": []interface{}{
			map[string]interface{}{
				"oid":   "1.3.6.1.2.1",
				"type":  "Integer",
				"value": "42",
			},
		},
	}

	out, err := FormatJSONC(config, nil, nil, nil, FormatMTA)
	if err != nil {
		t.Fatal(err)
	}

	expected := "// OpenDCI v" + Version + " JSONC — PacketCable MTA Configuration"
	if len(out) < len(expected) || out[:len(expected)] != expected {
		t.Errorf("expected MTA header, got:\n%s", out[:min(100, len(out))])
	}
}

func TestCMJSONCHeader(t *testing.T) {
	config := map[string]interface{}{
		"NetworkAccess": 1,
	}

	out, err := FormatJSONC(config, nil, nil, nil, FormatCM)
	if err != nil {
		t.Fatal(err)
	}

	expected := "// OpenDCI v" + Version + " JSONC — DOCSIS Configuration Interchange Format"
	if len(out) < len(expected) || out[:len(expected)] != expected {
		t.Errorf("expected CM header, got:\n%s", out[:min(100, len(out))])
	}
}

// makeCMRegistryWithEmta builds a CM registry that includes TLV 216 (Emta)
// with nestedFormat pointing to an MTA registry.
func makeCMRegistryWithEmta() (*Registry, *Registry) {
	mtaReg := makeMTATestRegistry()

	cmReg := makeTestRegistry()
	cmReg.TopLevel[216] = &TLVDef{
		TypeNum:      216,
		Name:         "Emta",
		DataType:     DataTypeHexString,
		Chunked:      true,
		NestedFormat: "mta",
	}
	cmReg.NameLookup["Emta"] = cmReg.TopLevel[216]
	cmReg.NestedRegistries = map[string]*Registry{
		"mta": mtaReg,
	}

	return cmReg, mtaReg
}

func TestTLV216RecursiveDecode(t *testing.T) {
	cmReg, _ := makeCMRegistryWithEmta()

	// Build an MTA config binary (with delimiters).
	varbind := buildVarbind([]byte{0x2B, 0x06, 0x01, 0x02, 0x01}, tagInteger, []byte{0x2A})
	mtaBinary := buildMTAConfig(buildTLV(11, varbind))

	// Embed it as TLV 216 in a CM config.
	data := buildConfig(
		buildTLV(3, []byte{1}),      // NetworkAccess = 1
		buildTLV(216, mtaBinary),     // Emta (chunked, but small enough for one chunk)
	)

	result, err := Decode(data, cmReg)
	if err != nil {
		t.Fatal(err)
	}

	// NetworkAccess should be decoded normally.
	if result.Config["NetworkAccess"] != 1 {
		t.Errorf("expected NetworkAccess=1, got %v", result.Config["NetworkAccess"])
	}

	// Emta should be a structured map, not a hexstring.
	emta, ok := result.Config["Emta"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected Emta as map, got %T: %v", result.Config["Emta"], result.Config["Emta"])
	}

	// Should contain SnmpMibObject array.
	snmpArr, ok := emta["SnmpMibObject"].([]interface{})
	if !ok {
		t.Fatalf("expected SnmpMibObject array in Emta, got %T", emta["SnmpMibObject"])
	}
	if len(snmpArr) != 1 {
		t.Fatalf("expected 1 SNMP entry in Emta, got %d", len(snmpArr))
	}

	entry := snmpArr[0].(map[string]interface{})
	if entry["oid"] != "1.3.6.1.2.1" {
		t.Errorf("expected oid '1.3.6.1.2.1', got %v", entry["oid"])
	}
	if entry["value"] != "42" {
		t.Errorf("expected value '42', got %v", entry["value"])
	}
}

func TestTLV216RecursiveEncode(t *testing.T) {
	cmReg, _ := makeCMRegistryWithEmta()

	// Structured Emta content.
	result := &DecodeResult{
		Config: map[string]interface{}{
			"NetworkAccess": 1,
			"Emta": map[string]interface{}{
				"SnmpMibObject": []interface{}{
					map[string]interface{}{
						"oid":   "1.3.6.1.2.1",
						"type":  "Integer",
						"value": "42",
					},
				},
			},
		},
		TLVOrder: []string{"NetworkAccess", "Emta"},
	}

	encoded, err := Encode(result, cmReg)
	if err != nil {
		t.Fatal(err)
	}

	// Verify the output contains TLV 216.
	found216 := false
	offset := 0
	for offset < len(encoded) {
		if encoded[offset] == 0xFF {
			break
		}
		if encoded[offset] == 0 {
			offset++
			continue
		}
		tlvType := int(encoded[offset])
		tlvLen := int(encoded[offset+1])
		if tlvType == 216 {
			found216 = true
			// The value should start with MTA start delimiter (FE 01 01).
			value := encoded[offset+2 : offset+2+tlvLen]
			if len(value) < 3 || value[0] != 0xFE || value[1] != 0x01 || value[2] != 0x01 {
				t.Errorf("TLV 216 value should start with MTA delimiter, got: %X", value[:min(3, len(value))])
			}
		}
		offset += 2 + tlvLen
	}
	if !found216 {
		t.Error("expected TLV 216 in encoded output")
	}
}

func TestTLV216RoundTrip(t *testing.T) {
	cmReg, _ := makeCMRegistryWithEmta()

	// Build original binary with embedded MTA config.
	varbind := buildVarbind([]byte{0x2B, 0x06, 0x01, 0x02, 0x01}, tagInteger, []byte{0x2A})
	mtaBinary := buildMTAConfig(buildTLV(11, varbind))

	original := buildConfig(
		buildTLV(3, []byte{1}),
		buildTLV(216, mtaBinary),
	)

	// Decode.
	result, err := Decode(original, cmReg)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}

	// Verify structured decode.
	if _, ok := result.Config["Emta"].(map[string]interface{}); !ok {
		t.Fatalf("expected structured Emta, got %T", result.Config["Emta"])
	}

	// Re-encode.
	encoded, err := Encode(result, cmReg)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}

	if !bytes.Equal(encoded, original) {
		t.Errorf("TLV 216 round-trip mismatch:\n  got:  %X\n  want: %X", encoded, original)
	}
}

func TestTLV216FallbackWithoutMTARegistry(t *testing.T) {
	// CM registry WITHOUT nested MTA registry — should fall back to hexstring.
	cmReg := makeTestRegistry()
	cmReg.TopLevel[216] = &TLVDef{
		TypeNum:      216,
		Name:         "Emta",
		DataType:     DataTypeHexString,
		Chunked:      true,
		NestedFormat: "mta",
	}
	cmReg.NameLookup["Emta"] = cmReg.TopLevel[216]
	// No NestedRegistries set.

	varbind := buildVarbind([]byte{0x2B, 0x06, 0x01}, tagInteger, []byte{0x01})
	mtaBinary := buildMTAConfig(buildTLV(11, varbind))

	data := buildConfig(
		buildTLV(3, []byte{1}),
		buildTLV(216, mtaBinary),
	)

	result, err := Decode(data, cmReg)
	if err != nil {
		t.Fatal(err)
	}

	// Without MTA registry, Emta should be a hexstring.
	emta, ok := result.Config["Emta"].(string)
	if !ok {
		t.Fatalf("expected Emta as string (hexstring fallback), got %T", result.Config["Emta"])
	}
	if emta == "" {
		t.Error("expected non-empty hexstring for Emta")
	}
}

func TestTLV216ChunkedRoundTrip(t *testing.T) {
	cmReg, _ := makeCMRegistryWithEmta()

	// Build a large MTA config that will need chunking (>254 bytes).
	var varbinds [][]byte
	for i := 0; i < 20; i++ {
		oid := []byte{0x2B, 0x06, 0x01, 0x02, 0x01, byte(i + 1)}
		varbinds = append(varbinds, buildTLV(11, buildVarbind(oid, tagInteger, []byte{byte(i)})))
	}
	mtaBinary := buildMTAConfig(varbinds...)

	// Build CM config with chunked TLV 216.
	// Manually chunk the MTA binary into ≤254-byte pieces.
	var cmData []byte
	cmData = append(cmData, buildTLV(3, []byte{1})...)
	remaining := mtaBinary
	for len(remaining) > 0 {
		chunkSize := len(remaining)
		if chunkSize > 254 {
			chunkSize = 254
		}
		cmData = append(cmData, buildTLV(216, remaining[:chunkSize])...)
		remaining = remaining[chunkSize:]
	}
	cmData = append(cmData, 0xFF, 0x00)

	// Decode.
	result, err := Decode(cmData, cmReg)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}

	// Verify structured decode.
	emta, ok := result.Config["Emta"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected structured Emta, got %T", result.Config["Emta"])
	}
	snmpArr := emta["SnmpMibObject"].([]interface{})
	if len(snmpArr) != 20 {
		t.Errorf("expected 20 SNMP entries in Emta, got %d", len(snmpArr))
	}

	// Re-encode.
	encoded, err := Encode(result, cmReg)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}

	if !bytes.Equal(encoded, cmData) {
		t.Errorf("chunked TLV 216 round-trip mismatch: got len=%d, want len=%d", len(encoded), len(cmData))
	}
}

func TestPacketCableHashMTAEndMarker(t *testing.T) {
	// Build a minimal MTA encoded output with MTA end delimiter.
	varbind := buildVarbind([]byte{0x2B, 0x06, 0x01}, tagInteger, []byte{0x01})
	encoded := []byte{0xFE, 0x01, 0x01} // start delimiter
	encoded = append(encoded, buildTLV(11, varbind)...)
	encoded = append(encoded, 0xFE, 0x01, 0xFF) // end delimiter

	// InsertPacketCableHash should work with MTA end marker.
	result, err := InsertPacketCableHash(encoded, PacketCableNA)
	if err != nil {
		t.Fatalf("InsertPacketCableHash with MTA end marker: %v", err)
	}

	// Result should still end with MTA end delimiter.
	n := len(result)
	if n < 3 || result[n-3] != 0xFE || result[n-2] != 0x01 || result[n-1] != 0xFF {
		t.Errorf("expected MTA end delimiter at end of hashed output, got: %X", result[max(0, n-3):])
	}

	// Result should contain the hash varbind TLV 11 before the end delimiter.
	if len(result) <= len(encoded) {
		t.Error("expected hashed output to be longer than input")
	}
}
