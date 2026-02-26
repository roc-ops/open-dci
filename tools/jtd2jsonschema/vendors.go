package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// MergeVendorSchemas reads vendor JTD schema files from the given directory
// and merges their definitions and mapped properties into the JSON Schema.
// Vendor definitions are added to $defs, and vendor-mapped properties are
// added as optional properties on VendorSpecificContainer.
func MergeVendorSchemas(jsonSchema map[string]interface{}, vendorsDir string) error {
	entries, err := os.ReadDir(vendorsDir)
	if err != nil {
		return fmt.Errorf("reading vendors directory: %w", err)
	}

	// Collect vendor files sorted by name for deterministic output.
	var vendorFiles []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".jtd.json") {
			vendorFiles = append(vendorFiles, entry.Name())
		}
	}
	sort.Strings(vendorFiles)

	// Get or create the $defs map.
	defs, ok := jsonSchema["$defs"].(map[string]interface{})
	if !ok {
		defs = make(map[string]interface{})
		jsonSchema["$defs"] = defs
	}

	// Get the VendorSpecificContainer definition.
	vsc, ok := defs["VendorSpecificContainer"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("VendorSpecificContainer not found in $defs")
	}

	// Get or create the properties map on VendorSpecificContainer.
	vscProps, ok := vsc["properties"].(map[string]interface{})
	if !ok {
		vscProps = make(map[string]interface{})
		vsc["properties"] = vscProps
	}

	for _, filename := range vendorFiles {
		filePath := filepath.Join(vendorsDir, filename)
		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("reading vendor schema %s: %w", filename, err)
		}

		var vendorSchema map[string]interface{}
		if err := json.Unmarshal(data, &vendorSchema); err != nil {
			return fmt.Errorf("parsing vendor schema %s: %w", filename, err)
		}

		vendorMeta, _ := vendorSchema["metadata"].(map[string]interface{})
		vendorName, _ := vendorMeta["vendorName"].(string)
		vendorOUI, _ := vendorMeta["oui"].(string)

		// Convert and add vendor definitions to $defs.
		vendorDefs, _ := vendorSchema["definitions"].(map[string]interface{})
		for defName, def := range vendorDefs {
			if defMap, ok := def.(map[string]interface{}); ok {
				defs[defName] = convertNode(defMap)
			}
		}

		// Process vendor mappings for TLV 43 (VendorSpecificContainer context).
		mapping, _ := vendorSchema["mapping"].(map[string]interface{})
		tlv43Mapping, ok := mapping["43"].(map[string]interface{})
		if !ok {
			continue
		}

		// Sort mapping keys for deterministic output.
		subTlvTypes := sortedKeys(tlv43Mapping)
		for _, subTlvType := range subTlvTypes {
			entry := tlv43Mapping[subTlvType]
			entryMap, ok := entry.(map[string]interface{})
			if !ok {
				continue
			}

			propName, _ := entryMap["name"].(string)
			refName, _ := entryMap["ref"].(string)
			if propName == "" || refName == "" {
				continue
			}

			// Prefix the property name with the vendor name to avoid collisions.
			qualifiedName := vendorName + "_" + propName

			// Build the property schema as a $ref to the vendor definition.
			propSchema := map[string]interface{}{
				"$ref": "#/$defs/" + refName,
			}

			// Add descriptive metadata.
			if vendorOUI != "" {
				propSchema["description"] = fmt.Sprintf("Vendor-specific: %s (OUI %s) sub-TLV %s", vendorName, vendorOUI, subTlvType)
				propSchema["x-docsis-vendorOUI"] = vendorOUI
			}
			if vendorName != "" {
				propSchema["x-docsis-vendorName"] = vendorName
			}

			vscProps[qualifiedName] = propSchema
		}
	}

	return nil
}
