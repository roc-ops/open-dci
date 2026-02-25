package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// TLVDef describes a single TLV type from the schema.
type TLVDef struct {
	TypeNum     int               // Wire type number (last segment of dotted notation)
	Name        string            // JSON property name
	DataType    DataType          // Wire data type
	Repeatable  bool              // Whether this TLV can appear multiple times
	Chunked     bool              // Whether this TLV uses 254-byte chunked encoding (e.g. CVC certificates)
	SubTLVs     map[int]*TLVDef   // For compound types: sub-TLV definitions keyed by type number
	RefName     string            // Name of the referenced definition (for compound types)
	ValidValues map[string]string // Human-readable labels for enum-like integer values
}

// Registry holds the complete TLV definition hierarchy.
type Registry struct {
	TopLevel   map[int]*TLVDef
	NameLookup map[string]*TLVDef // Reverse lookup: name → definition
}

// jtdSchema represents the top-level JTD schema structure.
type jtdSchema struct {
	Metadata           map[string]interface{}      `json:"metadata"`
	Definitions        map[string]json.RawMessage   `json:"definitions"`
	OptionalProperties map[string]json.RawMessage   `json:"optionalProperties"`
}

// jtdProperty represents a single JTD property.
type jtdProperty struct {
	Metadata           map[string]interface{}      `json:"metadata"`
	Type               string                      `json:"type"`
	Ref                string                      `json:"ref"`
	Enum               []string                    `json:"enum"`
	Elements           *jtdRef                     `json:"elements"`
	Properties         map[string]json.RawMessage   `json:"properties"`
	OptionalProperties map[string]json.RawMessage   `json:"optionalProperties"`
}

type jtdRef struct {
	Ref string `json:"ref"`
}

// LoadRegistry parses the JTD schema file and builds the TLV registry.
func LoadRegistry(schemaPath string) (*Registry, error) {
	data, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("reading schema: %w", err)
	}

	return LoadRegistryFromBytes(data)
}

// LoadRegistryFromBytes parses a JTD schema from raw JSON bytes and builds the
// TLV registry. This is the core parsing logic used by both file-based
// LoadRegistry and in-memory callers such as the WASM entry point.
func LoadRegistryFromBytes(data []byte) (*Registry, error) {
	var schema jtdSchema
	if err := json.Unmarshal(data, &schema); err != nil {
		return nil, fmt.Errorf("parsing schema: %w", err)
	}

	// Parse all definitions into jtdProperty structs for reference resolution.
	defs := make(map[string]*jtdProperty)
	for name, raw := range schema.Definitions {
		var prop jtdProperty
		if err := json.Unmarshal(raw, &prop); err != nil {
			return nil, fmt.Errorf("parsing definition %s: %w", name, err)
		}
		defs[name] = &prop
	}

	reg := &Registry{
		TopLevel:   make(map[int]*TLVDef),
		NameLookup: make(map[string]*TLVDef),
	}

	// Process top-level optionalProperties.
	for name, raw := range schema.OptionalProperties {
		var prop jtdProperty
		if err := json.Unmarshal(raw, &prop); err != nil {
			return nil, fmt.Errorf("parsing property %s: %w", name, err)
		}

		meta := prop.Metadata
		if meta == nil {
			continue
		}

		tlvTypeRaw, ok := meta["tlvType"]
		if !ok {
			continue
		}

		typeNum, err := extractTypeNum(tlvTypeRaw)
		if err != nil {
			continue // Skip non-numeric tlvType
		}

		dt := DataType(getStringMeta(meta, "dataType"))
		repeatable := getBoolMeta(meta, "repeatable")
		chunked := getBoolMeta(meta, "chunked")

		def := &TLVDef{
			TypeNum:    typeNum,
			Name:       name,
			DataType:   dt,
			Repeatable: repeatable,
			Chunked:    chunked,
		}
		def.ValidValues = getValidValuesMeta(meta)

		// For compound types, resolve sub-TLVs from references.
		if dt == DataTypeCompound {
			def.SubTLVs = make(map[int]*TLVDef)
			refName := resolveRefName(&prop)
			if refName != "" {
				def.RefName = refName
				if defProp, ok := defs[refName]; ok {
					populateSubTLVs(def, defProp, defs)
				}
			}
		}

		reg.TopLevel[typeNum] = def
		reg.NameLookup[name] = def
	}

	return reg, nil
}

// TopLevelByName returns the top-level TLVDef for the given property name, or nil.
func (r *Registry) TopLevelByName(name string) *TLVDef {
	return r.NameLookup[name]
}

// SubTLVByName returns the sub-TLV definition matching the given name within a
// compound TLVDef, or nil if not found.
func SubTLVByName(parent *TLVDef, name string) *TLVDef {
	if parent == nil || parent.SubTLVs == nil {
		return nil
	}
	for _, sub := range parent.SubTLVs {
		if sub.Name == name {
			return sub
		}
	}
	return nil
}

// resolveRefName extracts the definition reference name from a property.
func resolveRefName(prop *jtdProperty) string {
	if prop.Ref != "" {
		return prop.Ref
	}
	if prop.Elements != nil && prop.Elements.Ref != "" {
		return prop.Elements.Ref
	}
	return ""
}

// populateSubTLVs fills the SubTLVs map of a compound TLVDef from a definition.
func populateSubTLVs(parent *TLVDef, defProp *jtdProperty, defs map[string]*jtdProperty) {
	// Process required properties.
	for name, raw := range defProp.Properties {
		processSubTLVProperty(parent, name, raw, defs)
	}
	// Process optional properties.
	for name, raw := range defProp.OptionalProperties {
		processSubTLVProperty(parent, name, raw, defs)
	}
}

// processSubTLVProperty processes a single sub-TLV property and adds it to the parent.
func processSubTLVProperty(parent *TLVDef, name string, raw json.RawMessage, defs map[string]*jtdProperty) {
	var prop jtdProperty
	if err := json.Unmarshal(raw, &prop); err != nil {
		return
	}

	meta := prop.Metadata
	if meta == nil {
		return
	}

	tlvTypeRaw, ok := meta["tlvType"]
	if !ok {
		return // Skip properties without tlvType (like VendorSubTlvs)
	}

	// Extract the last segment of the dotted notation as the sub-TLV type number.
	typeNum, err := extractLastSegment(tlvTypeRaw)
	if err != nil {
		return
	}

	dt := DataType(getStringMeta(meta, "dataType"))
	repeatable := getBoolMeta(meta, "repeatable")

	subDef := &TLVDef{
		TypeNum:    typeNum,
		Name:       name,
		DataType:   dt,
		Repeatable: repeatable,
	}
	subDef.ValidValues = getValidValuesMeta(meta)

	// Recursively resolve compound sub-TLVs.
	if dt == DataTypeCompound {
		subDef.SubTLVs = make(map[int]*TLVDef)
		refName := resolveRefName(&prop)
		if refName != "" {
			subDef.RefName = refName
			if defProp, ok := defs[refName]; ok {
				populateSubTLVs(subDef, defProp, defs)
			}
		}
	}

	parent.SubTLVs[typeNum] = subDef
}

// extractTypeNum extracts a type number from a tlvType metadata value.
// Handles both integer (top-level) and string (dotted notation) formats.
func extractTypeNum(v interface{}) (int, error) {
	switch val := v.(type) {
	case float64:
		return int(val), nil
	case string:
		// For dotted notation at top level, use the first segment.
		parts := strings.Split(val, ".")
		return strconv.Atoi(parts[0])
	default:
		return 0, fmt.Errorf("unsupported tlvType type: %T", v)
	}
}

// extractLastSegment extracts the last numeric segment from a dotted tlvType string.
func extractLastSegment(v interface{}) (int, error) {
	switch val := v.(type) {
	case float64:
		return int(val), nil
	case string:
		parts := strings.Split(val, ".")
		last := parts[len(parts)-1]
		return strconv.Atoi(last)
	default:
		return 0, fmt.Errorf("unsupported tlvType type: %T", v)
	}
}

// getStringMeta safely extracts a string value from metadata.
func getStringMeta(meta map[string]interface{}, key string) string {
	if v, ok := meta[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// getBoolMeta safely extracts a boolean value from metadata.
func getBoolMeta(meta map[string]interface{}, key string) bool {
	if v, ok := meta[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

// getValidValuesMeta extracts the validValues map from metadata.
func getValidValuesMeta(meta map[string]interface{}) map[string]string {
	raw, ok := meta["validValues"]
	if !ok {
		return nil
	}
	obj, ok := raw.(map[string]interface{})
	if !ok {
		return nil
	}
	result := make(map[string]string, len(obj))
	for k, v := range obj {
		if s, ok := v.(string); ok {
			result[k] = s
		}
	}
	return result
}

// ValidValuesMap builds a flat lookup map of property name to valid values.
func (r *Registry) ValidValuesMap() map[string]map[string]string {
	result := make(map[string]map[string]string)
	for _, def := range r.TopLevel {
		collectValidValues(def, result)
	}
	return result
}

// collectValidValues recursively collects valid values from a TLVDef and its sub-TLVs.
func collectValidValues(def *TLVDef, result map[string]map[string]string) {
	if def.ValidValues != nil {
		result[def.Name] = def.ValidValues
	}
	for _, sub := range def.SubTLVs {
		collectValidValues(sub, result)
	}
}
