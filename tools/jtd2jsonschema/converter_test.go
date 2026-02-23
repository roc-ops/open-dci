package main

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

// parseJSON is a test helper that parses a JSON string into a map.
func parseJSON(t *testing.T, s string) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("failed to parse JSON: %v\nInput: %s", err, s)
	}
	return m
}

// assertJSONEqual compares two interface{} values as JSON-serialized strings
// for deep equality, providing better error output than reflect.DeepEqual.
func assertJSONEqual(t *testing.T, label string, expected, actual interface{}) {
	t.Helper()
	expJSON, _ := json.MarshalIndent(expected, "", "  ")
	actJSON, _ := json.MarshalIndent(actual, "", "  ")
	if string(expJSON) != string(actJSON) {
		t.Errorf("%s mismatch:\nexpected:\n%s\nactual:\n%s", label, string(expJSON), string(actJSON))
	}
}

// --- Type Form Tests ---

func TestConvertTypeUint8(t *testing.T) {
	jtd := parseJSON(t, `{"type": "uint8"}`)
	result := convertNode(jtd)
	if result["type"] != "integer" {
		t.Errorf("expected type=integer, got %v", result["type"])
	}
	if result["minimum"] != float64(0) {
		t.Errorf("expected minimum=0, got %v", result["minimum"])
	}
	if result["maximum"] != float64(255) {
		t.Errorf("expected maximum=255, got %v", result["maximum"])
	}
}

func TestConvertTypeUint16(t *testing.T) {
	jtd := parseJSON(t, `{"type": "uint16"}`)
	result := convertNode(jtd)
	if result["type"] != "integer" {
		t.Errorf("expected type=integer, got %v", result["type"])
	}
	if result["minimum"] != float64(0) {
		t.Errorf("expected minimum=0, got %v", result["minimum"])
	}
	if result["maximum"] != float64(65535) {
		t.Errorf("expected maximum=65535, got %v", result["maximum"])
	}
}

func TestConvertTypeUint32(t *testing.T) {
	jtd := parseJSON(t, `{"type": "uint32"}`)
	result := convertNode(jtd)
	if result["type"] != "integer" {
		t.Errorf("expected type=integer, got %v", result["type"])
	}
	if result["minimum"] != float64(0) {
		t.Errorf("expected minimum=0, got %v", result["minimum"])
	}
	if result["maximum"] != float64(4294967295) {
		t.Errorf("expected maximum=4294967295, got %v", result["maximum"])
	}
}

func TestConvertTypeInt8(t *testing.T) {
	jtd := parseJSON(t, `{"type": "int8"}`)
	result := convertNode(jtd)
	if result["minimum"] != float64(-128) {
		t.Errorf("expected minimum=-128, got %v", result["minimum"])
	}
	if result["maximum"] != float64(127) {
		t.Errorf("expected maximum=127, got %v", result["maximum"])
	}
}

func TestConvertTypeInt16(t *testing.T) {
	jtd := parseJSON(t, `{"type": "int16"}`)
	result := convertNode(jtd)
	if result["minimum"] != float64(-32768) {
		t.Errorf("expected minimum=-32768, got %v", result["minimum"])
	}
	if result["maximum"] != float64(32767) {
		t.Errorf("expected maximum=32767, got %v", result["maximum"])
	}
}

func TestConvertTypeInt32(t *testing.T) {
	jtd := parseJSON(t, `{"type": "int32"}`)
	result := convertNode(jtd)
	if result["minimum"] != float64(-2147483648) {
		t.Errorf("expected minimum=-2147483648, got %v", result["minimum"])
	}
	if result["maximum"] != float64(2147483647) {
		t.Errorf("expected maximum=2147483647, got %v", result["maximum"])
	}
}

func TestConvertTypeString(t *testing.T) {
	jtd := parseJSON(t, `{"type": "string"}`)
	result := convertNode(jtd)
	if result["type"] != "string" {
		t.Errorf("expected type=string, got %v", result["type"])
	}
}

func TestConvertTypeBoolean(t *testing.T) {
	jtd := parseJSON(t, `{"type": "boolean"}`)
	result := convertNode(jtd)
	if result["type"] != "boolean" {
		t.Errorf("expected type=boolean, got %v", result["type"])
	}
}

func TestConvertTypeFloat32(t *testing.T) {
	jtd := parseJSON(t, `{"type": "float32"}`)
	result := convertNode(jtd)
	if result["type"] != "number" {
		t.Errorf("expected type=number, got %v", result["type"])
	}
}

func TestConvertTypeFloat64(t *testing.T) {
	jtd := parseJSON(t, `{"type": "float64"}`)
	result := convertNode(jtd)
	if result["type"] != "number" {
		t.Errorf("expected type=number, got %v", result["type"])
	}
}

func TestConvertTypeTimestamp(t *testing.T) {
	jtd := parseJSON(t, `{"type": "timestamp"}`)
	result := convertNode(jtd)
	if result["type"] != "string" {
		t.Errorf("expected type=string, got %v", result["type"])
	}
	if result["format"] != "date-time" {
		t.Errorf("expected format=date-time, got %v", result["format"])
	}
}

// --- Enum Form Tests ---

func TestConvertEnum(t *testing.T) {
	jtd := parseJSON(t, `{"enum": ["Integer", "String", "IPAddress"]}`)
	result := convertNode(jtd)
	if result["type"] != "string" {
		t.Errorf("expected type=string, got %v", result["type"])
	}
	enumVals, ok := result["enum"].([]interface{})
	if !ok {
		t.Fatalf("expected enum to be a slice, got %T", result["enum"])
	}
	if len(enumVals) != 3 {
		t.Errorf("expected 3 enum values, got %d", len(enumVals))
	}
}

// --- Elements Form Tests ---

func TestConvertElements(t *testing.T) {
	jtd := parseJSON(t, `{"elements": {"type": "string"}}`)
	result := convertNode(jtd)
	if result["type"] != "array" {
		t.Errorf("expected type=array, got %v", result["type"])
	}
	items, ok := result["items"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected items to be a map, got %T", result["items"])
	}
	if items["type"] != "string" {
		t.Errorf("expected items.type=string, got %v", items["type"])
	}
}

func TestConvertElementsWithRef(t *testing.T) {
	jtd := parseJSON(t, `{"elements": {"ref": "Foo"}}`)
	result := convertNode(jtd)
	if result["type"] != "array" {
		t.Errorf("expected type=array, got %v", result["type"])
	}
	items, ok := result["items"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected items to be a map, got %T", result["items"])
	}
	if items["$ref"] != "#/$defs/Foo" {
		t.Errorf("expected $ref=#/$defs/Foo, got %v", items["$ref"])
	}
}

// --- Values Form Tests ---

func TestConvertValues(t *testing.T) {
	jtd := parseJSON(t, `{"values": {"type": "uint32"}}`)
	result := convertNode(jtd)
	if result["type"] != "object" {
		t.Errorf("expected type=object, got %v", result["type"])
	}
	ap, ok := result["additionalProperties"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected additionalProperties to be a map, got %T", result["additionalProperties"])
	}
	if ap["type"] != "integer" {
		t.Errorf("expected additionalProperties.type=integer, got %v", ap["type"])
	}
}

// --- Ref Form Tests ---

func TestConvertRef(t *testing.T) {
	jtd := parseJSON(t, `{"ref": "SnmpMibEntry"}`)
	result := convertNode(jtd)
	if result["$ref"] != "#/$defs/SnmpMibEntry" {
		t.Errorf("expected $ref=#/$defs/SnmpMibEntry, got %v", result["$ref"])
	}
}

// --- Properties Form Tests ---

func TestConvertProperties(t *testing.T) {
	jtd := parseJSON(t, `{
		"properties": {
			"name": {"type": "string"},
			"age": {"type": "uint8"}
		}
	}`)
	result := convertNode(jtd)
	if result["type"] != "object" {
		t.Errorf("expected type=object, got %v", result["type"])
	}

	props, ok := result["properties"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected properties map, got %T", result["properties"])
	}
	if len(props) != 2 {
		t.Errorf("expected 2 properties, got %d", len(props))
	}

	required, ok := result["required"].([]interface{})
	if !ok {
		t.Fatalf("expected required array, got %T", result["required"])
	}
	if len(required) != 2 {
		t.Errorf("expected 2 required, got %d", len(required))
	}

	// additionalProperties should default to false.
	if result["additionalProperties"] != false {
		t.Errorf("expected additionalProperties=false, got %v", result["additionalProperties"])
	}
}

func TestConvertOptionalProperties(t *testing.T) {
	jtd := parseJSON(t, `{
		"optionalProperties": {
			"nickname": {"type": "string"}
		}
	}`)
	result := convertNode(jtd)
	if result["type"] != "object" {
		t.Errorf("expected type=object, got %v", result["type"])
	}

	props, ok := result["properties"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected properties map, got %T", result["properties"])
	}
	if len(props) != 1 {
		t.Errorf("expected 1 property, got %d", len(props))
	}

	// No required array.
	if _, ok := result["required"]; ok {
		t.Errorf("expected no required array for optional-only properties")
	}
}

func TestConvertMixedProperties(t *testing.T) {
	jtd := parseJSON(t, `{
		"properties": {
			"id": {"type": "string"}
		},
		"optionalProperties": {
			"label": {"type": "string"}
		}
	}`)
	result := convertNode(jtd)
	props := result["properties"].(map[string]interface{})
	if len(props) != 2 {
		t.Errorf("expected 2 properties (required + optional), got %d", len(props))
	}
	required := result["required"].([]interface{})
	if len(required) != 1 || required[0] != "id" {
		t.Errorf("expected required=[id], got %v", required)
	}
}

func TestConvertAdditionalPropertiesTrue(t *testing.T) {
	jtd := parseJSON(t, `{
		"properties": {
			"id": {"type": "string"}
		},
		"additionalProperties": true
	}`)
	result := convertNode(jtd)
	if result["additionalProperties"] != true {
		t.Errorf("expected additionalProperties=true, got %v", result["additionalProperties"])
	}
}

// --- Nullable Tests ---

func TestConvertNullable(t *testing.T) {
	jtd := parseJSON(t, `{"type": "string", "nullable": true}`)
	result := convertNode(jtd)
	anyOf, ok := result["anyOf"].([]interface{})
	if !ok {
		t.Fatalf("expected anyOf array, got %T: %v", result["anyOf"], result)
	}
	if len(anyOf) != 2 {
		t.Errorf("expected 2 anyOf entries, got %d", len(anyOf))
	}
	// Second entry should be null.
	nullEntry, ok := anyOf[1].(map[string]interface{})
	if !ok {
		t.Fatalf("expected null entry to be a map")
	}
	if nullEntry["type"] != "null" {
		t.Errorf("expected second anyOf entry type=null, got %v", nullEntry["type"])
	}
}

// --- Metadata Tests ---

func TestMetadataDescription(t *testing.T) {
	jtd := parseJSON(t, `{
		"metadata": {"description": "Test field", "spec": "CM-SP-MULPIv4.0 C.1.1.1"},
		"type": "uint32"
	}`)
	result := convertNode(jtd)
	if result["description"] != "Test field" {
		t.Errorf("expected description='Test field', got %v", result["description"])
	}
	if result["x-docsis-spec"] != "CM-SP-MULPIv4.0 C.1.1.1" {
		t.Errorf("expected x-docsis-spec, got %v", result["x-docsis-spec"])
	}
}

func TestMetadataTlvType(t *testing.T) {
	jtd := parseJSON(t, `{
		"metadata": {"tlvType": 3, "dataType": "uint8", "spec": "CM-SP-MULPIv4.0 C.1.1.3"},
		"type": "uint8"
	}`)
	result := convertNode(jtd)
	if result["x-docsis-tlvType"] != float64(3) {
		t.Errorf("expected x-docsis-tlvType=3, got %v", result["x-docsis-tlvType"])
	}
	if result["x-docsis-dataType"] != "uint8" {
		t.Errorf("expected x-docsis-dataType=uint8, got %v", result["x-docsis-dataType"])
	}
}

func TestMetadataSynthesizedDescription(t *testing.T) {
	jtd := parseJSON(t, `{
		"metadata": {"tlvType": 1, "dataType": "uint32", "spec": "CM-SP-MULPIv4.0 C.1.1.1"},
		"type": "uint32"
	}`)
	result := convertNode(jtd)
	desc := result["description"].(string)
	expected := "TLV 1 - uint32 (CM-SP-MULPIv4.0 C.1.1.1)"
	if desc != expected {
		t.Errorf("expected description=%q, got %q", expected, desc)
	}
}

func TestMetadataValidValues(t *testing.T) {
	jtd := parseJSON(t, `{
		"metadata": {
			"spec": "CM-SP-MULPIv4.0 C.2.2.5.2",
			"validValues": {"2": "best effort", "6": "unsolicited grant service"}
		},
		"type": "uint8"
	}`)
	result := convertNode(jtd)
	vv, ok := result["x-docsis-validValues"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected x-docsis-validValues map, got %T", result["x-docsis-validValues"])
	}
	if vv["2"] != "best effort" {
		t.Errorf("expected validValues[2]='best effort', got %v", vv["2"])
	}
}

// --- Definitions Tests ---

func TestConvertDefinitions(t *testing.T) {
	jtd := parseJSON(t, `{
		"definitions": {
			"Foo": {"type": "string"}
		},
		"optionalProperties": {
			"bar": {"ref": "Foo"}
		}
	}`)
	result := ConvertJTDToJSONSchema(jtd)
	defs, ok := result["$defs"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected $defs map, got %T", result["$defs"])
	}
	foo, ok := defs["Foo"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected Foo definition map, got %T", defs["Foo"])
	}
	if foo["type"] != "string" {
		t.Errorf("expected Foo.type=string, got %v", foo["type"])
	}

	// Check that $schema is set.
	if result["$schema"] != "https://json-schema.org/draft/2020-12/schema" {
		t.Errorf("expected $schema to be set")
	}
}

// --- Top-Level Conversion Test ---

func TestConvertTopLevelSchema(t *testing.T) {
	jtd := parseJSON(t, `{
		"metadata": {
			"title": "Test Schema",
			"version": "1.0.0",
			"spec": "test-spec"
		},
		"definitions": {
			"Entry": {
				"properties": {
					"name": {"type": "string"}
				}
			}
		},
		"optionalProperties": {
			"items": {
				"elements": {"ref": "Entry"}
			}
		}
	}`)
	result := ConvertJTDToJSONSchema(jtd)

	if result["$schema"] != "https://json-schema.org/draft/2020-12/schema" {
		t.Error("missing $schema")
	}
	if result["title"] != "Test Schema" {
		t.Errorf("expected title='Test Schema', got %v", result["title"])
	}
	if _, ok := result["$defs"]; !ok {
		t.Error("missing $defs")
	}
	if _, ok := result["properties"]; !ok {
		t.Error("missing properties")
	}
}

// --- Empty Form Test ---

func TestConvertEmptyForm(t *testing.T) {
	jtd := parseJSON(t, `{}`)
	result := convertNode(jtd)
	// Empty form should produce an empty schema (accepts anything).
	if len(result) != 0 {
		t.Errorf("expected empty result for empty form, got %v", result)
	}
}

// --- Round-Trip Consistency Test ---

func TestRoundTripConsistency(t *testing.T) {
	// Convert the same JTD twice and verify identical output.
	jtdStr := `{
		"definitions": {
			"Item": {
				"properties": {
					"id": {"type": "uint32"},
					"name": {"type": "string"}
				},
				"optionalProperties": {
					"tags": {"elements": {"type": "string"}}
				}
			}
		},
		"optionalProperties": {
			"items": {"elements": {"ref": "Item"}}
		}
	}`
	jtd1 := parseJSON(t, jtdStr)
	jtd2 := parseJSON(t, jtdStr)

	result1 := ConvertJTDToJSONSchema(jtd1)
	result2 := ConvertJTDToJSONSchema(jtd2)

	json1, _ := json.Marshal(result1)
	json2, _ := json.Marshal(result2)

	if string(json1) != string(json2) {
		t.Error("round-trip consistency failed: two conversions of the same input produced different output")
	}
}

// --- Full DOCSIS Schema Conversion Test ---

func TestConvertDocsisSchema(t *testing.T) {
	// Load and convert the actual DOCSIS JTD schema.
	schemaPath := "../../schemas/docsis-config.jtd.json"
	data, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Skipf("skipping: cannot read schema file: %v", err)
	}

	var jtd map[string]interface{}
	if err := json.Unmarshal(data, &jtd); err != nil {
		t.Fatalf("failed to parse JTD schema: %v", err)
	}

	result := ConvertJTDToJSONSchema(jtd)

	// Verify basic structure.
	if result["$schema"] != "https://json-schema.org/draft/2020-12/schema" {
		t.Error("missing $schema")
	}

	defs, ok := result["$defs"].(map[string]interface{})
	if !ok {
		t.Fatal("missing $defs")
	}

	// Check that expected definitions exist.
	expectedDefs := []string{
		"SnmpMibEntry",
		"ServiceFlowErrorEntry",
		"UpstreamServiceFlowEntry",
		"DownstreamServiceFlowEntry",
		"DocsisExtensionFieldEntry",
	}
	for _, name := range expectedDefs {
		if _, ok := defs[name]; !ok {
			t.Errorf("missing definition: %s", name)
		}
	}

	// Check that expected top-level properties exist.
	props, ok := result["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("missing properties")
	}
	expectedProps := []string{
		"DownstreamFrequency", "UpstreamChannelId", "NetworkAccess",
		"SwUpgradeFilename", "SnmpMibObject", "CpeEthernetMacAddress",
		"MaxNumCpes", "PrivacyEnable", "UpstreamServiceFlow",
		"DownstreamServiceFlow", "DocsisExtensionField", "SwUpgradeIpv6TftpServer",
	}
	for _, name := range expectedProps {
		if _, ok := props[name]; !ok {
			t.Errorf("missing top-level property: %s", name)
		}
	}

	// Verify DownstreamFrequency is integer with correct range.
	df, ok := props["DownstreamFrequency"].(map[string]interface{})
	if !ok {
		t.Fatal("DownstreamFrequency is not a map")
	}
	if df["type"] != "integer" {
		t.Errorf("DownstreamFrequency type=%v, expected integer", df["type"])
	}
	if df["minimum"] != float64(0) || df["maximum"] != float64(4294967295) {
		t.Errorf("DownstreamFrequency range incorrect: min=%v max=%v", df["minimum"], df["maximum"])
	}

	// Verify SnmpMibObject is an array with $ref items.
	snmp, ok := props["SnmpMibObject"].(map[string]interface{})
	if !ok {
		t.Fatal("SnmpMibObject is not a map")
	}
	if snmp["type"] != "array" {
		t.Errorf("SnmpMibObject type=%v, expected array", snmp["type"])
	}
	items, ok := snmp["items"].(map[string]interface{})
	if !ok {
		t.Fatal("SnmpMibObject.items is not a map")
	}
	if items["$ref"] != "#/$defs/SnmpMibEntry" {
		t.Errorf("SnmpMibObject.items.$ref=%v, expected #/$defs/SnmpMibEntry", items["$ref"])
	}

	// Verify DocsisExtensionFieldEntry has additionalProperties: true.
	def43, ok := defs["DocsisExtensionFieldEntry"].(map[string]interface{})
	if !ok {
		t.Fatal("DocsisExtensionFieldEntry is not a map")
	}
	if def43["additionalProperties"] != true {
		t.Errorf("DocsisExtensionFieldEntry.additionalProperties=%v, expected true", def43["additionalProperties"])
	}

	// Verify x-docsis-spec is present on DownstreamFrequency.
	if df["x-docsis-spec"] == nil {
		t.Error("DownstreamFrequency missing x-docsis-spec")
	}
}

// --- All Metadata Spec Present Test ---

func TestAllPropertiesHaveSpec(t *testing.T) {
	// Verify that every property in the DOCSIS schema has a spec reference
	// in its metadata.
	schemaPath := "../../schemas/docsis-config.jtd.json"
	data, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Skipf("skipping: cannot read schema file: %v", err)
	}

	var jtd map[string]interface{}
	if err := json.Unmarshal(data, &jtd); err != nil {
		t.Fatalf("failed to parse JTD schema: %v", err)
	}

	checkSpecPresent(t, "root", jtd)
}

// checkSpecPresent recursively checks that every property has metadata.spec.
func checkSpecPresent(t *testing.T, path string, node map[string]interface{}) {
	t.Helper()

	// Check properties.
	if props, ok := node["properties"].(map[string]interface{}); ok {
		for name, prop := range props {
			propPath := path + "." + name
			if propMap, ok := prop.(map[string]interface{}); ok {
				assertHasSpec(t, propPath, propMap)
				checkSpecPresent(t, propPath, propMap)
			}
		}
	}

	// Check optionalProperties.
	if props, ok := node["optionalProperties"].(map[string]interface{}); ok {
		for name, prop := range props {
			propPath := path + "." + name
			if propMap, ok := prop.(map[string]interface{}); ok {
				assertHasSpec(t, propPath, propMap)
				checkSpecPresent(t, propPath, propMap)
			}
		}
	}

	// Check definitions.
	if defs, ok := node["definitions"].(map[string]interface{}); ok {
		for name, def := range defs {
			defPath := path + ".definitions." + name
			if defMap, ok := def.(map[string]interface{}); ok {
				// Definitions themselves should have spec in their metadata.
				assertHasSpec(t, defPath, defMap)
				checkSpecPresent(t, defPath, defMap)
			}
		}
	}

	// Check elements (for array types, the elements schema may have a ref but
	// the outer node should have spec).
	if elem, ok := node["elements"].(map[string]interface{}); ok {
		// Elements with ref don't need their own spec — the ref target has it.
		if _, hasRef := elem["ref"]; !hasRef {
			checkSpecPresent(t, path+".elements", elem)
		}
	}
}

// assertHasSpec checks that a JTD node has metadata.spec.
func assertHasSpec(t *testing.T, path string, node map[string]interface{}) {
	t.Helper()
	// Nodes with only a ref don't need their own spec (they inherit from the definition).
	if _, hasRef := node["ref"]; hasRef {
		// But if there is metadata, it should have spec.
		if meta, ok := node["metadata"].(map[string]interface{}); ok {
			if _, hasSpec := meta["spec"]; !hasSpec {
				t.Errorf("%s: has metadata but missing spec", path)
			}
		}
		return
	}
	meta, ok := node["metadata"].(map[string]interface{})
	if !ok {
		t.Errorf("%s: missing metadata", path)
		return
	}
	if _, hasSpec := meta["spec"]; !hasSpec {
		t.Errorf("%s: missing metadata.spec", path)
	}
}

// --- Discriminator Form Test ---

func TestConvertDiscriminator(t *testing.T) {
	jtd := parseJSON(t, `{
		"discriminator": "kind",
		"mapping": {
			"cat": {
				"properties": {
					"purrs": {"type": "boolean"}
				}
			},
			"dog": {
				"properties": {
					"barks": {"type": "boolean"}
				}
			}
		}
	}`)
	result := convertNode(jtd)
	oneOf, ok := result["oneOf"].([]interface{})
	if !ok {
		t.Fatalf("expected oneOf, got %T", result["oneOf"])
	}
	if len(oneOf) != 2 {
		t.Errorf("expected 2 oneOf entries, got %d", len(oneOf))
	}
}

// --- Serialization Determinism Test ---

func TestSerializationDeterminism(t *testing.T) {
	// Run conversion multiple times and verify JSON output is identical.
	jtdStr := `{
		"metadata": {"title": "Test", "spec": "test"},
		"definitions": {
			"A": {"type": "string"},
			"B": {"type": "uint8"}
		},
		"optionalProperties": {
			"x": {"ref": "A"},
			"y": {"ref": "B"},
			"z": {"type": "boolean"}
		}
	}`

	var outputs []string
	for i := 0; i < 10; i++ {
		jtd := parseJSON(t, jtdStr)
		result := ConvertJTDToJSONSchema(jtd)
		out, _ := json.Marshal(result)
		outputs = append(outputs, string(out))
	}

	for i := 1; i < len(outputs); i++ {
		if outputs[i] != outputs[0] {
			t.Errorf("non-deterministic output at iteration %d", i)
		}
	}
}

// --- Nested Compound Test ---

func TestNestedCompound(t *testing.T) {
	jtd := parseJSON(t, `{
		"optionalProperties": {
			"outer": {
				"metadata": {"spec": "test"},
				"optionalProperties": {
					"inner": {
						"metadata": {"spec": "test"},
						"type": "string"
					}
				}
			}
		}
	}`)
	result := convertNode(jtd)
	props := result["properties"].(map[string]interface{})
	outer := props["outer"].(map[string]interface{})
	if outer["type"] != "object" {
		t.Errorf("expected outer.type=object, got %v", outer["type"])
	}
	innerProps := outer["properties"].(map[string]interface{})
	inner := innerProps["inner"].(map[string]interface{})
	if inner["type"] != "string" {
		t.Errorf("expected inner.type=string, got %v", inner["type"])
	}
}

// --- JSON Output Validity Test ---

func TestOutputIsValidJSON(t *testing.T) {
	jtd := parseJSON(t, `{
		"metadata": {"title": "Test", "spec": "test"},
		"optionalProperties": {
			"field": {"type": "uint32", "metadata": {"spec": "test"}}
		}
	}`)
	result := ConvertJTDToJSONSchema(jtd)
	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Verify it parses back.
	var parsed map[string]interface{}
	if err := json.Unmarshal(output, &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	// Verify deep equality after round-trip through JSON.
	if !reflect.DeepEqual(result, parsed) {
		t.Error("JSON round-trip changed the structure")
	}
}
