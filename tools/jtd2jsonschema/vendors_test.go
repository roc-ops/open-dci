package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestMergeVendorSchemas(t *testing.T) {
	// Create a temporary vendor schema directory.
	tmpDir := t.TempDir()

	vendorSchema := map[string]interface{}{
		"metadata": map[string]interface{}{
			"oui":        "001018",
			"vendorName": "Broadcom",
		},
		"definitions": map[string]interface{}{
			"Broadcom_PowerSavingMode": map[string]interface{}{
				"metadata": map[string]interface{}{
					"description": "Power saving mode",
					"subTlvType":  1,
					"dataType":    "uint8",
					"spec":        "vendor-defined",
				},
				"type": "uint8",
			},
		},
		"mapping": map[string]interface{}{
			"43": map[string]interface{}{
				"1": map[string]interface{}{
					"name": "PowerSavingMode",
					"ref":  "Broadcom_PowerSavingMode",
				},
			},
		},
	}

	data, _ := json.Marshal(vendorSchema)
	os.WriteFile(filepath.Join(tmpDir, "001018.jtd.json"), data, 0644)

	// Build a minimal JSON Schema with VendorSpecificContainer.
	jsonSchema := map[string]interface{}{
		"$schema": "https://json-schema.org/draft/2020-12/schema",
		"$defs": map[string]interface{}{
			"VendorSpecificContainer": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"VendorId": map[string]interface{}{
						"type": "string",
					},
				},
				"required":            []interface{}{"VendorId"},
				"additionalProperties": false,
			},
		},
	}

	err := MergeVendorSchemas(jsonSchema, tmpDir)
	if err != nil {
		t.Fatalf("MergeVendorSchemas failed: %v", err)
	}

	defs := jsonSchema["$defs"].(map[string]interface{})

	// Verify vendor definition was added.
	bpsm, ok := defs["Broadcom_PowerSavingMode"].(map[string]interface{})
	if !ok {
		t.Fatal("expected Broadcom_PowerSavingMode in $defs")
	}
	if bpsm["type"] != "integer" {
		t.Errorf("expected Broadcom_PowerSavingMode.type=integer, got %v", bpsm["type"])
	}

	// Verify vendor property was added to VendorSpecificContainer.
	vsc := defs["VendorSpecificContainer"].(map[string]interface{})
	props := vsc["properties"].(map[string]interface{})

	prop, ok := props["Broadcom_PowerSavingMode"].(map[string]interface{})
	if !ok {
		t.Fatal("expected Broadcom_PowerSavingMode property on VendorSpecificContainer")
	}
	if prop["$ref"] != "#/$defs/Broadcom_PowerSavingMode" {
		t.Errorf("expected $ref=#/$defs/Broadcom_PowerSavingMode, got %v", prop["$ref"])
	}
	if prop["x-docsis-vendorOUI"] != "001018" {
		t.Errorf("expected x-docsis-vendorOUI=001018, got %v", prop["x-docsis-vendorOUI"])
	}
}

func TestMergeVendorSchemasNoVendorDir(t *testing.T) {
	jsonSchema := map[string]interface{}{
		"$defs": map[string]interface{}{
			"VendorSpecificContainer": map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
	}

	err := MergeVendorSchemas(jsonSchema, "/nonexistent/path")
	if err == nil {
		t.Error("expected error for nonexistent directory")
	}
}

func TestMergeVendorSchemasEmptyDir(t *testing.T) {
	tmpDir := t.TempDir()
	jsonSchema := map[string]interface{}{
		"$defs": map[string]interface{}{
			"VendorSpecificContainer": map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
	}

	err := MergeVendorSchemas(jsonSchema, tmpDir)
	if err != nil {
		t.Fatalf("expected no error for empty dir, got: %v", err)
	}
}

func TestMergeVendorSchemasWithRealSchema(t *testing.T) {
	// Load the actual DOCSIS schema and convert, then merge with real vendor schemas.
	schemaPath := "../../schemas/docsis-config.jtd.json"
	vendorsDir := "../../schemas/vendors"

	data, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Skipf("skipping: cannot read schema file: %v", err)
	}

	var jtd map[string]interface{}
	if err := json.Unmarshal(data, &jtd); err != nil {
		t.Fatalf("failed to parse JTD schema: %v", err)
	}

	jsonSchema := ConvertJTDToJSONSchema(jtd)

	// Check vendors directory exists.
	if _, err := os.Stat(vendorsDir); os.IsNotExist(err) {
		t.Skipf("skipping: vendors directory not found: %v", err)
	}

	err = MergeVendorSchemas(jsonSchema, vendorsDir)
	if err != nil {
		t.Fatalf("MergeVendorSchemas failed: %v", err)
	}

	defs := jsonSchema["$defs"].(map[string]interface{})

	// Verify VendorSpecificContainer exists and has expected structure.
	vsc, ok := defs["VendorSpecificContainer"].(map[string]interface{})
	if !ok {
		t.Fatal("VendorSpecificContainer not found in $defs")
	}

	// Should NOT have additionalProperties: true.
	if vsc["additionalProperties"] == true {
		t.Error("VendorSpecificContainer should not have additionalProperties: true")
	}

	props := vsc["properties"].(map[string]interface{})

	// Should have VendorId.
	if _, ok := props["VendorId"]; !ok {
		t.Error("VendorSpecificContainer missing VendorId")
	}

	// Should have FFFFFF sub-TLV properties.
	expectedFFFFFF := []string{
		"CmLoadBalancingPolicyId",
		"CmLoadBalancingPriority",
		"CmLoadBalancingGroupId",
		"CmRangingClassIdExtension",
		"L2vpnEncoding",
		"ExtendedCmtsMicConfigurationSetting",
		"SourceAddressVerification",
		"CableModemAttributeMasks",
		"IpMulticastJoinAuthorization",
		"ServiceTypeIdentifier",
		"DemarcAutoConfiguration",
		"VendorSubTlvs",
	}
	for _, name := range expectedFFFFFF {
		if _, ok := props[name]; !ok {
			t.Errorf("VendorSpecificContainer missing FFFFFF property: %s", name)
		}
	}

	// Should have Broadcom vendor properties.
	if _, ok := props["Broadcom_PowerSavingMode"]; !ok {
		t.Error("VendorSpecificContainer missing Broadcom_PowerSavingMode")
	}

	// Broadcom definitions should be in $defs.
	if _, ok := defs["Broadcom_PowerSavingMode"]; !ok {
		t.Error("Broadcom_PowerSavingMode not found in $defs")
	}

	// L2vpnEncoding should have a $ref (via items since it's an array).
	l2vpn, ok := props["L2vpnEncoding"].(map[string]interface{})
	if !ok {
		t.Fatal("L2vpnEncoding not found")
	}
	items, ok := l2vpn["items"].(map[string]interface{})
	if !ok {
		t.Fatal("L2vpnEncoding.items not found (expected array)")
	}
	if items["$ref"] != "#/$defs/L2vpnEncodingEntry" {
		t.Errorf("L2vpnEncoding.items.$ref=%v, expected #/$defs/L2vpnEncodingEntry", items["$ref"])
	}
}
