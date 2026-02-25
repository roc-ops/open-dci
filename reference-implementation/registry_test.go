package main

import (
	"path/filepath"
	"runtime"
	"testing"
)

func schemaPath(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("unable to get caller info")
	}
	return filepath.Join(filepath.Dir(filename), "..", "schemas", "docsis-config.jtd.json")
}

func TestLoadRegistrySuccess(t *testing.T) {
	reg, err := LoadRegistry(schemaPath(t))
	if err != nil {
		t.Fatal(err)
	}

	if reg == nil {
		t.Fatal("registry is nil")
	}

	// The schema has 82 top-level TLV properties.
	// Some may share type numbers (e.g. TLV 43.8 is both VendorId and CableModemAttributeMasks),
	// but top-level should be close to 82.
	if len(reg.TopLevel) < 50 {
		t.Errorf("expected at least 50 top-level TLVs, got %d", len(reg.TopLevel))
	}

	t.Logf("Loaded %d top-level TLV definitions", len(reg.TopLevel))
}

func TestLoadRegistryTLV1(t *testing.T) {
	reg, err := LoadRegistry(schemaPath(t))
	if err != nil {
		t.Fatal(err)
	}

	def, ok := reg.TopLevel[1]
	if !ok {
		t.Fatal("TLV 1 (DownstreamFrequency) not found")
	}

	if def.Name != "DownstreamFrequency" {
		t.Errorf("expected name 'DownstreamFrequency', got %q", def.Name)
	}
	if def.DataType != DataTypeUint32 {
		t.Errorf("expected data type uint32, got %q", def.DataType)
	}
	if def.Repeatable {
		t.Error("TLV 1 should not be repeatable")
	}
}

func TestLoadRegistryTLV3(t *testing.T) {
	reg, err := LoadRegistry(schemaPath(t))
	if err != nil {
		t.Fatal(err)
	}

	def, ok := reg.TopLevel[3]
	if !ok {
		t.Fatal("TLV 3 (NetworkAccess) not found")
	}

	if def.Name != "NetworkAccess" {
		t.Errorf("expected name 'NetworkAccess', got %q", def.Name)
	}
	if def.DataType != DataTypeUint8 {
		t.Errorf("expected data type uint8, got %q", def.DataType)
	}
}

func TestLoadRegistryTLV10(t *testing.T) {
	reg, err := LoadRegistry(schemaPath(t))
	if err != nil {
		t.Fatal(err)
	}

	def, ok := reg.TopLevel[10]
	if !ok {
		t.Fatal("TLV 10 (SnmpWriteAccessControl) not found")
	}

	if def.Name != "SnmpWriteAccessControl" {
		t.Errorf("expected name 'SnmpWriteAccessControl', got %q", def.Name)
	}
	if def.DataType != DataTypeCompound {
		t.Errorf("expected data type compound, got %q", def.DataType)
	}
	if !def.Repeatable {
		t.Error("TLV 10 should be repeatable")
	}
}

func TestLoadRegistryTLV11(t *testing.T) {
	reg, err := LoadRegistry(schemaPath(t))
	if err != nil {
		t.Fatal(err)
	}

	def, ok := reg.TopLevel[11]
	if !ok {
		t.Fatal("TLV 11 (SnmpMibObject) not found")
	}

	if def.Name != "SnmpMibObject" {
		t.Errorf("expected name 'SnmpMibObject', got %q", def.Name)
	}
	if def.DataType != DataTypeCompound {
		t.Errorf("expected data type compound, got %q", def.DataType)
	}
	if !def.Repeatable {
		t.Error("TLV 11 should be repeatable")
	}
}

func TestLoadRegistryTLV24Compound(t *testing.T) {
	reg, err := LoadRegistry(schemaPath(t))
	if err != nil {
		t.Fatal(err)
	}

	def, ok := reg.TopLevel[24]
	if !ok {
		t.Fatal("TLV 24 (UpstreamServiceFlow) not found")
	}

	if def.Name != "UpstreamServiceFlow" {
		t.Errorf("expected name 'UpstreamServiceFlow', got %q", def.Name)
	}
	if def.DataType != DataTypeCompound {
		t.Errorf("expected data type compound, got %q", def.DataType)
	}
	if !def.Repeatable {
		t.Error("TLV 24 should be repeatable")
	}

	// Check sub-TLVs
	if def.SubTLVs == nil {
		t.Fatal("TLV 24 SubTLVs is nil")
	}

	// Sub-TLV 1 = ServiceFlowReference (uint16)
	sub1, ok := def.SubTLVs[1]
	if !ok {
		t.Fatal("TLV 24.1 (ServiceFlowReference) not found")
	}
	if sub1.Name != "ServiceFlowReference" {
		t.Errorf("expected 'ServiceFlowReference', got %q", sub1.Name)
	}
	if sub1.DataType != DataTypeUint16 {
		t.Errorf("expected uint16, got %q", sub1.DataType)
	}

	// Sub-TLV 4 = ServiceClassName (string)
	sub4, ok := def.SubTLVs[4]
	if !ok {
		t.Fatal("TLV 24.4 (ServiceClassName) not found")
	}
	if sub4.Name != "ServiceClassName" {
		t.Errorf("expected 'ServiceClassName', got %q", sub4.Name)
	}
	if sub4.DataType != DataTypeString {
		t.Errorf("expected string, got %q", sub4.DataType)
	}

	// Sub-TLV 6 = QosParamSetType (uint8)
	sub6, ok := def.SubTLVs[6]
	if !ok {
		t.Fatal("TLV 24.6 (QosParamSetType) not found")
	}
	if sub6.DataType != DataTypeUint8 {
		t.Errorf("expected uint8, got %q", sub6.DataType)
	}

	t.Logf("TLV 24 has %d sub-TLVs", len(def.SubTLVs))
}

func TestLoadRegistryTLV43(t *testing.T) {
	reg, err := LoadRegistry(schemaPath(t))
	if err != nil {
		t.Fatal(err)
	}

	def, ok := reg.TopLevel[43]
	if !ok {
		t.Fatal("TLV 43 (DocsisExtensionField) not found")
	}

	if def.Name != "DocsisExtensionField" {
		t.Errorf("expected name 'DocsisExtensionField', got %q", def.Name)
	}
	if def.DataType != DataTypeCompound {
		t.Errorf("expected data type compound, got %q", def.DataType)
	}
	if !def.Repeatable {
		t.Error("TLV 43 should be repeatable")
	}

	// Check that VendorId (sub-TLV 8) is present
	sub8, ok := def.SubTLVs[8]
	if !ok {
		t.Fatal("TLV 43.8 not found in sub-TLVs")
	}
	// It could be VendorId or CableModemAttributeMasks depending on parse order.
	t.Logf("TLV 43.8 name: %s", sub8.Name)

	// Check a general extension sub-TLV
	sub1, ok := def.SubTLVs[1]
	if !ok {
		t.Fatal("TLV 43.1 (CmLoadBalancingPolicyId) not found")
	}
	if sub1.DataType != DataTypeUint32 {
		t.Errorf("expected uint32, got %q", sub1.DataType)
	}

	t.Logf("TLV 43 has %d sub-TLVs", len(def.SubTLVs))
}

func TestLoadRegistryTLV32Chunked(t *testing.T) {
	reg, err := LoadRegistry(schemaPath(t))
	if err != nil {
		t.Fatal(err)
	}

	def, ok := reg.TopLevel[32]
	if !ok {
		t.Fatal("TLV 32 (ManufacturerCvc) not found")
	}

	if def.Name != "ManufacturerCvc" {
		t.Errorf("expected name 'ManufacturerCvc', got %q", def.Name)
	}
	if !def.Chunked {
		t.Error("TLV 32 should be chunked")
	}
	if def.Repeatable {
		t.Error("TLV 32 should not be repeatable (chunked instead)")
	}
}

func TestLoadRegistryIPv4(t *testing.T) {
	reg, err := LoadRegistry(schemaPath(t))
	if err != nil {
		t.Fatal(err)
	}

	def, ok := reg.TopLevel[20]
	if !ok {
		t.Fatal("TLV 20 (TftpServerProvisionedModemIpv4Address) not found")
	}

	if def.DataType != DataTypeIPv4Address {
		t.Errorf("expected ipv4Address, got %q", def.DataType)
	}
}

func TestLoadRegistryIPv6(t *testing.T) {
	reg, err := LoadRegistry(schemaPath(t))
	if err != nil {
		t.Fatal(err)
	}

	def, ok := reg.TopLevel[58]
	if !ok {
		t.Fatal("TLV 58 (SwUpgradeIpv6TftpServer) not found")
	}

	if def.DataType != DataTypeIPv6Address {
		t.Errorf("expected ipv6Address, got %q", def.DataType)
	}
}

func TestLoadRegistryMacAddress(t *testing.T) {
	reg, err := LoadRegistry(schemaPath(t))
	if err != nil {
		t.Fatal(err)
	}

	def, ok := reg.TopLevel[14]
	if !ok {
		t.Fatal("TLV 14 (CpeEthernetMacAddress) not found")
	}

	if def.DataType != DataTypeMacAddress {
		t.Errorf("expected macAddress, got %q", def.DataType)
	}
}

func TestLoadRegistryInvalidPath(t *testing.T) {
	_, err := LoadRegistry("/nonexistent/path/schema.json")
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestLoadRegistryValidValues(t *testing.T) {
	reg, err := LoadRegistry(schemaPath(t))
	if err != nil {
		t.Fatal(err)
	}

	// TLV 3 = NetworkAccess — should have validValues.
	def, ok := reg.TopLevel[3]
	if !ok {
		t.Fatal("TLV 3 (NetworkAccess) not found")
	}
	if def.ValidValues == nil {
		t.Fatal("NetworkAccess.ValidValues is nil")
	}
	if def.ValidValues["0"] != "disabled" {
		t.Errorf("expected ValidValues[\"0\"]=\"disabled\", got %q", def.ValidValues["0"])
	}
	if def.ValidValues["1"] != "enabled" {
		t.Errorf("expected ValidValues[\"1\"]=\"enabled\", got %q", def.ValidValues["1"])
	}

	// TLV 1 = DownstreamFrequency — should NOT have validValues.
	df, ok := reg.TopLevel[1]
	if !ok {
		t.Fatal("TLV 1 (DownstreamFrequency) not found")
	}
	if df.ValidValues != nil {
		t.Errorf("DownstreamFrequency should not have validValues, got %v", df.ValidValues)
	}
}

func TestValidValuesMap(t *testing.T) {
	reg, err := LoadRegistry(schemaPath(t))
	if err != nil {
		t.Fatal(err)
	}

	vvm := reg.ValidValuesMap()

	// Check top-level properties.
	for _, name := range []string{"NetworkAccess", "PrivacyEnable", "SnmpCpeAccessControl"} {
		if _, ok := vvm[name]; !ok {
			t.Errorf("expected %s in ValidValuesMap", name)
		}
	}

	// Check sub-TLV properties (from definitions).
	for _, name := range []string{"QosParamSetType", "DataRateUnitSetting", "DutControl"} {
		if _, ok := vvm[name]; !ok {
			t.Errorf("expected %s in ValidValuesMap", name)
		}
	}

	// Verify a specific mapping.
	if vvm["DataRateUnitSetting"]["0"] != "bits per second (bps)" {
		t.Errorf("expected DataRateUnitSetting[\"0\"]=\"bits per second (bps)\", got %q", vvm["DataRateUnitSetting"]["0"])
	}

	t.Logf("ValidValuesMap has %d entries", len(vvm))
}

func TestLoadRegistryTLV103LengthSize(t *testing.T) {
	reg, err := LoadRegistry(schemaPath(t))
	if err != nil {
		t.Fatal(err)
	}

	def, ok := reg.TopLevel[103]
	if !ok {
		t.Fatal("TLV 103 (CmSshServerConfigurationSettings) not found")
	}

	if def.LengthSize != 2 {
		t.Errorf("expected LengthSize=2 for TLV 103, got %d", def.LengthSize)
	}
	if def.DataType != DataTypeCompound {
		t.Errorf("expected compound data type, got %q", def.DataType)
	}
}

func TestLoadRegistryTLV104LengthSize(t *testing.T) {
	reg, err := LoadRegistry(schemaPath(t))
	if err != nil {
		t.Fatal(err)
	}

	def, ok := reg.TopLevel[104]
	if !ok {
		t.Fatal("TLV 104 (SecurityConfigurationSettings) not found")
	}

	if def.LengthSize != 2 {
		t.Errorf("expected LengthSize=2 for TLV 104, got %d", def.LengthSize)
	}
	if def.DataType != DataTypeCompound {
		t.Errorf("expected compound data type, got %q", def.DataType)
	}
}

func TestLoadRegistryTLV1LengthSizeDefault(t *testing.T) {
	reg, err := LoadRegistry(schemaPath(t))
	if err != nil {
		t.Fatal(err)
	}

	def, ok := reg.TopLevel[1]
	if !ok {
		t.Fatal("TLV 1 (DownstreamFrequency) not found")
	}

	if def.LengthSize != 1 {
		t.Errorf("expected LengthSize=1 (default) for TLV 1, got %d", def.LengthSize)
	}
}

func TestLoadRegistrySubTLVLengthSize(t *testing.T) {
	reg, err := LoadRegistry(schemaPath(t))
	if err != nil {
		t.Fatal(err)
	}

	// TLV 103 sub-TLV 3 = SnmpBasedAuthConfiguration (compound),
	// which itself has sub-TLV 1 = SshCmCds with tlvLengthSize=2.
	def, ok := reg.TopLevel[103]
	if !ok {
		t.Fatal("TLV 103 not found")
	}

	subDef3, ok := def.SubTLVs[3]
	if !ok {
		t.Fatal("TLV 103.3 (SnmpBasedAuthConfiguration) not found")
	}
	if subDef3.DataType != DataTypeCompound {
		t.Fatalf("expected compound, got %q", subDef3.DataType)
	}

	sshCmCds, ok := subDef3.SubTLVs[1]
	if !ok {
		t.Fatal("TLV 103.3.1 (SshCmCds) not found")
	}
	if sshCmCds.LengthSize != 2 {
		t.Errorf("expected LengthSize=2 for TLV 103.3.1 (SshCmCds), got %d", sshCmCds.LengthSize)
	}

	// Sub-TLV 103.3.2 = SshCmCdsDownloadUrl should have default LengthSize=1
	sshCmCdsUrl, ok := subDef3.SubTLVs[2]
	if !ok {
		t.Fatal("TLV 103.3.2 (SshCmCdsDownloadUrl) not found")
	}
	if sshCmCdsUrl.LengthSize != 1 {
		t.Errorf("expected LengthSize=1 for TLV 103.3.2, got %d", sshCmCdsUrl.LengthSize)
	}
}
