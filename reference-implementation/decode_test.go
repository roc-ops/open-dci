package main

import (
	"encoding/json"
	"testing"
)

// buildTLV constructs a binary TLV: [type][length][value].
func buildTLV(t int, value []byte) []byte {
	result := []byte{byte(t), byte(len(value))}
	result = append(result, value...)
	return result
}

// buildConfig concatenates TLVs and appends the end-of-data marker (TLV 255).
func buildConfig(tlvs ...[]byte) []byte {
	var data []byte
	for _, tlv := range tlvs {
		data = append(data, tlv...)
	}
	// End-of-data marker: type 255, length 0
	data = append(data, 0xFF, 0x00)
	return data
}

// makeTestRegistry builds a minimal registry for testing.
func makeTestRegistry() *Registry {
	reg := &Registry{
		TopLevel:   map[int]*TLVDef{
			1: {
				TypeNum:  1,
				Name:     "DownstreamFrequency",
				DataType: DataTypeUint32,
			},
			2: {
				TypeNum:  2,
				Name:     "UpstreamChannelId",
				DataType: DataTypeUint8,
			},
			3: {
				TypeNum:  3,
				Name:     "NetworkAccess",
				DataType: DataTypeUint8,
			},
			6: {
				TypeNum:  6,
				Name:     "CmMic",
				DataType: DataTypeHexString,
			},
			7: {
				TypeNum:  7,
				Name:     "CmtsMic",
				DataType: DataTypeHexString,
			},
			9: {
				TypeNum:  9,
				Name:     "SwUpgradeFilename",
				DataType: DataTypeString,
			},
			11: {
				TypeNum:    11,
				Name:       "SnmpMibObject",
				DataType:   DataTypeCompound,
				Repeatable: true,
			},
			14: {
				TypeNum:  14,
				Name:     "CpeEthernetMacAddress",
				DataType: DataTypeMacAddress,
			},
			18: {
				TypeNum:  18,
				Name:     "MaxNumCpes",
				DataType: DataTypeUint8,
			},
			20: {
				TypeNum:  20,
				Name:     "TftpServerProvisionedModemIpv4Address",
				DataType: DataTypeIPv4Address,
			},
			24: {
				TypeNum:    24,
				Name:       "UpstreamServiceFlow",
				DataType:   DataTypeCompound,
				Repeatable: true,
				SubTLVs: map[int]*TLVDef{
					1: {TypeNum: 1, Name: "ServiceFlowReference", DataType: DataTypeUint16},
					4: {TypeNum: 4, Name: "ServiceClassName", DataType: DataTypeString},
					6: {TypeNum: 6, Name: "QosParamSetType", DataType: DataTypeUint8},
					8: {TypeNum: 8, Name: "MaxSustainedTrafficRate", DataType: DataTypeUint32},
				},
			},
			43: {
				TypeNum:    43,
				Name:       "DocsisExtensionField",
				DataType:   DataTypeCompound,
				Repeatable: true,
				SubTLVs: map[int]*TLVDef{
					1:  {TypeNum: 1, Name: "CmLoadBalancingPolicyId", DataType: DataTypeUint32},
					8:  {TypeNum: 8, Name: "VendorId", DataType: DataTypeHexString},
					10: {TypeNum: 10, Name: "ServiceTypeIdentifier", DataType: DataTypeString},
				},
			},
		},
		NameLookup: make(map[string]*TLVDef),
	}
	for _, def := range reg.TopLevel {
		reg.NameLookup[def.Name] = def
	}
	return reg
}

func TestDecodeSimpleTLVs(t *testing.T) {
	reg := makeTestRegistry()

	data := buildConfig(
		buildTLV(3, []byte{1}),                       // NetworkAccess = 1
		buildTLV(1, []byte{0x15, 0x75, 0x2A, 0x00}),  // DownstreamFrequency = 360000000
		buildTLV(2, []byte{5}),                        // UpstreamChannelId = 5
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
	if result.Config["UpstreamChannelId"] != 5 {
		t.Errorf("expected UpstreamChannelId=5, got %v", result.Config["UpstreamChannelId"])
	}
}

func TestDecodeStringTLV(t *testing.T) {
	reg := makeTestRegistry()

	data := buildConfig(
		buildTLV(9, []byte("firmware.bin\x00")),
	)

	result, err := Decode(data, reg)
	if err != nil {
		t.Fatal(err)
	}

	if result.Config["SwUpgradeFilename"] != "firmware.bin" {
		t.Errorf("expected 'firmware.bin', got %v", result.Config["SwUpgradeFilename"])
	}
}

func TestDecodeMacAddress(t *testing.T) {
	reg := makeTestRegistry()

	data := buildConfig(
		buildTLV(14, []byte{0x00, 0x1A, 0x2B, 0x3C, 0x4D, 0x5E}),
	)

	result, err := Decode(data, reg)
	if err != nil {
		t.Fatal(err)
	}

	if result.Config["CpeEthernetMacAddress"] != "001A2B3C4D5E" {
		t.Errorf("expected '001A2B3C4D5E', got %v", result.Config["CpeEthernetMacAddress"])
	}
}

func TestDecodeIPv4Address(t *testing.T) {
	reg := makeTestRegistry()

	data := buildConfig(
		buildTLV(20, []byte{10, 1, 2, 3}),
	)

	result, err := Decode(data, reg)
	if err != nil {
		t.Fatal(err)
	}

	if result.Config["TftpServerProvisionedModemIpv4Address"] != "10.1.2.3" {
		t.Errorf("expected '10.1.2.3', got %v", result.Config["TftpServerProvisionedModemIpv4Address"])
	}
}

func TestDecodeCompoundTLV(t *testing.T) {
	reg := makeTestRegistry()

	// Build a TLV 24 (UpstreamServiceFlow) with sub-TLVs
	subTLVs := append(
		buildTLV(1, []byte{0x00, 0x01}),              // ServiceFlowReference = 1
		buildTLV(6, []byte{7})...,                     // QosParamSetType = 7
	)
	subTLVs = append(subTLVs,
		buildTLV(8, []byte{0x00, 0x00, 0x27, 0x10})..., // MaxSustainedTrafficRate = 10000
	)

	data := buildConfig(
		buildTLV(24, subTLVs),
	)

	result, err := Decode(data, reg)
	if err != nil {
		t.Fatal(err)
	}

	// UpstreamServiceFlow should be an array (repeatable)
	usfs, ok := result.Config["UpstreamServiceFlow"].([]interface{})
	if !ok {
		t.Fatalf("expected []interface{}, got %T", result.Config["UpstreamServiceFlow"])
	}
	if len(usfs) != 1 {
		t.Fatalf("expected 1 USF, got %d", len(usfs))
	}

	usf := usfs[0].(map[string]interface{})
	if usf["ServiceFlowReference"] != 1 {
		t.Errorf("expected ServiceFlowReference=1, got %v", usf["ServiceFlowReference"])
	}
	if usf["QosParamSetType"] != 7 {
		t.Errorf("expected QosParamSetType=7, got %v", usf["QosParamSetType"])
	}
	if usf["MaxSustainedTrafficRate"] != 10000 {
		t.Errorf("expected MaxSustainedTrafficRate=10000, got %v", usf["MaxSustainedTrafficRate"])
	}
}

func TestDecodeRepeatableTLV(t *testing.T) {
	reg := makeTestRegistry()

	// Two UpstreamServiceFlow TLVs
	usf1 := buildTLV(1, []byte{0x00, 0x01}) // ServiceFlowReference = 1
	usf2 := buildTLV(1, []byte{0x00, 0x02}) // ServiceFlowReference = 2

	data := buildConfig(
		buildTLV(24, usf1),
		buildTLV(24, usf2),
	)

	result, err := Decode(data, reg)
	if err != nil {
		t.Fatal(err)
	}

	usfs, ok := result.Config["UpstreamServiceFlow"].([]interface{})
	if !ok {
		t.Fatalf("expected []interface{}, got %T", result.Config["UpstreamServiceFlow"])
	}
	if len(usfs) != 2 {
		t.Errorf("expected 2 USFs, got %d", len(usfs))
	}
}

func TestDecodeUnknownTLV(t *testing.T) {
	reg := makeTestRegistry()

	data := buildConfig(
		buildTLV(3, []byte{1}),
		buildTLV(200, []byte{0xAB, 0xCD}), // Unknown TLV
	)

	result, err := Decode(data, reg)
	if err != nil {
		t.Fatal(err)
	}

	unknowns, ok := result.Config["UnknownTlvs"].([]interface{})
	if !ok {
		t.Fatal("expected UnknownTlvs array")
	}
	if len(unknowns) != 1 {
		t.Fatalf("expected 1 unknown, got %d", len(unknowns))
	}

	unk := unknowns[0].(map[string]interface{})
	if unk["type"] != 200 {
		t.Errorf("expected type 200, got %v", unk["type"])
	}
	if unk["value"] != "ABCD" {
		t.Errorf("expected value 'ABCD', got %v", unk["value"])
	}
}

func TestDecodeTLV11Snmp(t *testing.T) {
	reg := makeTestRegistry()

	// Build an SNMP varbind: OID 1.3.6.1.2.1, Integer = 42
	oidBytes := []byte{0x2B, 0x06, 0x01, 0x02, 0x01}
	varbind := buildVarbind(oidBytes, tagInteger, []byte{0x2A})

	data := buildConfig(
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

func TestDecodeTLV43GeneralExtension(t *testing.T) {
	reg := makeTestRegistry()

	// VendorId sub-TLV 8 = FF:FF:FF (General Extension)
	vendorId := buildTLV(8, []byte{0xFF, 0xFF, 0xFF})
	// CmLoadBalancingPolicyId sub-TLV 1 = 100
	lbPolicy := buildTLV(1, []byte{0x00, 0x00, 0x00, 0x64})

	tlv43Value := append(vendorId, lbPolicy...)

	data := buildConfig(
		buildTLV(43, tlv43Value),
	)

	result, err := Decode(data, reg)
	if err != nil {
		t.Fatal(err)
	}

	defs, ok := result.Config["DocsisExtensionField"].([]interface{})
	if !ok {
		t.Fatalf("expected []interface{}, got %T", result.Config["DocsisExtensionField"])
	}
	if len(defs) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(defs))
	}

	entry := defs[0].(map[string]interface{})
	if entry["VendorId"] != "FFFFFF" {
		t.Errorf("expected VendorId 'FFFFFF', got %v", entry["VendorId"])
	}
	if entry["CmLoadBalancingPolicyId"] != 100 {
		t.Errorf("expected CmLoadBalancingPolicyId=100, got %v", entry["CmLoadBalancingPolicyId"])
	}
}

func TestDecodeTLV43VendorSpecific(t *testing.T) {
	reg := makeTestRegistry()

	// VendorId sub-TLV 8 = 00:11:22 (vendor-specific)
	vendorId := buildTLV(8, []byte{0x00, 0x11, 0x22})
	// Vendor sub-TLV type 1 with value
	vendorSub := buildTLV(1, []byte{0xAA, 0xBB})

	tlv43Value := append(vendorId, vendorSub...)

	data := buildConfig(
		buildTLV(43, tlv43Value),
	)

	result, err := Decode(data, reg)
	if err != nil {
		t.Fatal(err)
	}

	defs, ok := result.Config["DocsisExtensionField"].([]interface{})
	if !ok {
		t.Fatalf("expected []interface{}, got %T", result.Config["DocsisExtensionField"])
	}

	entry := defs[0].(map[string]interface{})
	if entry["VendorId"] != "001122" {
		t.Errorf("expected VendorId '001122', got %v", entry["VendorId"])
	}

	vendorSubTlvs, ok := entry["VendorSubTlvs"].([]interface{})
	if !ok {
		t.Fatal("expected VendorSubTlvs array")
	}
	if len(vendorSubTlvs) != 1 {
		t.Fatalf("expected 1 vendor sub-TLV, got %d", len(vendorSubTlvs))
	}

	sub := vendorSubTlvs[0].(map[string]interface{})
	if sub["type"] != 1 {
		t.Errorf("expected type 1, got %v", sub["type"])
	}
	if sub["value"] != "AABB" {
		t.Errorf("expected value 'AABB', got %v", sub["value"])
	}
}

func TestDecodeMICExtraction(t *testing.T) {
	reg := makeTestRegistry()

	cmMic := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10}
	cmtsMic := []byte{0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20}

	data := buildConfig(
		buildTLV(3, []byte{1}),
		buildTLV(6, cmMic),
		buildTLV(7, cmtsMic),
	)

	result, err := Decode(data, reg)
	if err != nil {
		t.Fatal(err)
	}

	if len(result.CmMic) != 16 {
		t.Errorf("expected 16 byte CM MIC, got %d", len(result.CmMic))
	}
	if len(result.CmtsMic) != 16 {
		t.Errorf("expected 16 byte CMTS MIC, got %d", len(result.CmtsMic))
	}
}

func TestDecodePadByte(t *testing.T) {
	reg := makeTestRegistry()

	// Include pad bytes (TLV 0) between regular TLVs
	var data []byte
	data = append(data, 0x00)                         // pad
	data = append(data, buildTLV(3, []byte{1})...)     // NetworkAccess
	data = append(data, 0x00)                         // pad
	data = append(data, 0xFF, 0x00)                    // end-of-data

	result, err := Decode(data, reg)
	if err != nil {
		t.Fatal(err)
	}

	if result.Config["NetworkAccess"] != 1 {
		t.Errorf("expected NetworkAccess=1, got %v", result.Config["NetworkAccess"])
	}
}

func TestDecodeEndOfData(t *testing.T) {
	reg := makeTestRegistry()

	// Data after end-of-data should be ignored
	data := buildConfig(
		buildTLV(3, []byte{1}),
	)
	// Append extra bytes that should be ignored
	data = append(data, buildTLV(2, []byte{99})...)

	result, err := Decode(data, reg)
	if err != nil {
		t.Fatal(err)
	}

	if result.Config["NetworkAccess"] != 1 {
		t.Errorf("expected NetworkAccess=1, got %v", result.Config["NetworkAccess"])
	}
	if _, ok := result.Config["UpstreamChannelId"]; ok {
		t.Error("data after end-of-data should be ignored")
	}
}

func TestDecodeJSONOutput(t *testing.T) {
	reg := makeTestRegistry()

	data := buildConfig(
		buildTLV(3, []byte{1}),
		buildTLV(1, []byte{0x15, 0x75, 0x2A, 0x00}),
		buildTLV(18, []byte{16}),
	)

	result, err := Decode(data, reg)
	if err != nil {
		t.Fatal(err)
	}

	// Verify it can be marshalled to JSON
	jsonBytes, err := json.MarshalIndent(result.Config, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	// Verify it can be unmarshalled back
	var parsed map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &parsed); err != nil {
		t.Fatal(err)
	}

	// JSON numbers are float64
	if int(parsed["NetworkAccess"].(float64)) != 1 {
		t.Errorf("expected NetworkAccess=1 in JSON, got %v", parsed["NetworkAccess"])
	}

	t.Logf("JSON output:\n%s", string(jsonBytes))
}

func TestDecodeTruncatedTLV(t *testing.T) {
	reg := makeTestRegistry()

	// TLV with length exceeding available data
	data := []byte{3, 10, 1} // Type 3, length 10, but only 1 value byte

	_, err := Decode(data, reg)
	if err == nil {
		t.Fatal("expected error for truncated TLV")
	}
}

func TestDecodeEmptyConfig(t *testing.T) {
	reg := makeTestRegistry()

	// Just the end-of-data marker
	data := []byte{0xFF, 0x00}

	result, err := Decode(data, reg)
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Config) != 0 {
		t.Errorf("expected empty config, got %d entries", len(result.Config))
	}
}
