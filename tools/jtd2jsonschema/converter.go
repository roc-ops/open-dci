// Package main provides a converter from JSON Type Definition (RFC 8927)
// to JSON Schema 2020-12.
package main

import (
	"fmt"
	"sort"
)

// ConvertJTDToJSONSchema converts a top-level JTD schema (with optional
// definitions, metadata, properties, optionalProperties, etc.) into a
// JSON Schema 2020-12 document.
func ConvertJTDToJSONSchema(jtd map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{
		"$schema": "https://json-schema.org/draft/2020-12/schema",
	}

	// Convert top-level metadata to JSON Schema metadata.
	if meta, ok := jtd["metadata"].(map[string]interface{}); ok {
		applyMetadata(result, meta)
	}

	// Convert definitions to $defs.
	if defs, ok := jtd["definitions"].(map[string]interface{}); ok {
		jsonDefs := make(map[string]interface{})
		for name, def := range defs {
			if defMap, ok := def.(map[string]interface{}); ok {
				jsonDefs[name] = convertNode(defMap)
			}
		}
		result["$defs"] = jsonDefs
	}

	// Convert the root schema body (properties, optionalProperties, type, etc.)
	// by merging conversion results into the top-level object.
	rootBody := convertNode(jtd)
	for k, v := range rootBody {
		if k == "$schema" {
			continue
		}
		// Don't overwrite $defs we already set.
		if k == "$defs" {
			continue
		}
		result[k] = v
	}

	return result
}

// convertNode converts a single JTD node to its JSON Schema equivalent.
func convertNode(jtd map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	// Handle metadata.
	if meta, ok := jtd["metadata"].(map[string]interface{}); ok {
		applyMetadata(result, meta)
	}

	// Handle nullable wrapping at the end if needed.
	nullable := false
	if n, ok := jtd["nullable"].(bool); ok && n {
		nullable = true
	}

	// Determine which JTD form this node uses.
	switch {
	case hasKey(jtd, "ref"):
		// Reference form.
		refName, _ := jtd["ref"].(string)
		result["$ref"] = "#/$defs/" + refName

	case hasKey(jtd, "type"):
		// Type form.
		typeName, _ := jtd["type"].(string)
		applyType(result, typeName)

	case hasKey(jtd, "enum"):
		// Enum form.
		result["type"] = "string"
		result["enum"] = jtd["enum"]

	case hasKey(jtd, "elements"):
		// Elements form (array).
		result["type"] = "array"
		if elemNode, ok := jtd["elements"].(map[string]interface{}); ok {
			result["items"] = convertNode(elemNode)
		}

	case hasKey(jtd, "values"):
		// Values form (string-keyed map).
		result["type"] = "object"
		if valNode, ok := jtd["values"].(map[string]interface{}); ok {
			result["additionalProperties"] = convertNode(valNode)
		}

	case hasKey(jtd, "discriminator"):
		// Discriminator form — map to oneOf with discriminator property.
		disc, _ := jtd["discriminator"].(string)
		mapping, _ := jtd["mapping"].(map[string]interface{})
		var oneOf []interface{}
		for tag, schema := range mapping {
			if schemaMap, ok := schema.(map[string]interface{}); ok {
				converted := convertNode(schemaMap)
				// Ensure the discriminator property is present.
				if props, ok := converted["properties"].(map[string]interface{}); ok {
					props[disc] = map[string]interface{}{
						"type":  "string",
						"const": tag,
					}
				} else {
					converted["properties"] = map[string]interface{}{
						disc: map[string]interface{}{
							"type":  "string",
							"const": tag,
						},
					}
					if req, ok := converted["required"].([]interface{}); ok {
						converted["required"] = append(req, disc)
					} else {
						converted["required"] = []interface{}{disc}
					}
				}
				oneOf = append(oneOf, converted)
			}
		}
		// Sort oneOf for deterministic output.
		sort.Slice(oneOf, func(i, j int) bool {
			iMap := oneOf[i].(map[string]interface{})
			jMap := oneOf[j].(map[string]interface{})
			iProps := iMap["properties"].(map[string]interface{})
			jProps := jMap["properties"].(map[string]interface{})
			iConst := iProps[disc].(map[string]interface{})["const"].(string)
			jConst := jProps[disc].(map[string]interface{})["const"].(string)
			return iConst < jConst
		})
		result["oneOf"] = oneOf

	case hasKey(jtd, "properties") || hasKey(jtd, "optionalProperties"):
		// Properties form (object).
		convertProperties(result, jtd)

	default:
		// Empty form — accepts anything.
	}

	// Handle additionalProperties at the JTD level.
	if ap, ok := jtd["additionalProperties"].(bool); ok && ap {
		// Only set if we have an object with properties.
		if _, hasProps := result["properties"]; hasProps {
			result["additionalProperties"] = true
		} else if result["type"] == nil {
			// If this is a properties form that we just built, it's already set.
			result["additionalProperties"] = true
		}
	}

	if nullable {
		return wrapNullable(result)
	}
	return result
}

// convertProperties handles the JTD properties and optionalProperties forms,
// converting them to JSON Schema properties and required arrays.
func convertProperties(result map[string]interface{}, jtd map[string]interface{}) {
	result["type"] = "object"
	props := make(map[string]interface{})
	var required []interface{}

	// Required properties.
	if reqProps, ok := jtd["properties"].(map[string]interface{}); ok {
		names := sortedKeys(reqProps)
		for _, name := range names {
			if propMap, ok := reqProps[name].(map[string]interface{}); ok {
				props[name] = convertNode(propMap)
			}
			required = append(required, name)
		}
	}

	// Optional properties.
	if optProps, ok := jtd["optionalProperties"].(map[string]interface{}); ok {
		names := sortedKeys(optProps)
		for _, name := range names {
			if propMap, ok := optProps[name].(map[string]interface{}); ok {
				props[name] = convertNode(propMap)
			}
			// Not added to required.
		}
	}

	if len(props) > 0 {
		result["properties"] = props
	}
	if len(required) > 0 {
		sort.Slice(required, func(i, j int) bool {
			return required[i].(string) < required[j].(string)
		})
		result["required"] = required
	}

	// additionalProperties defaults to false in JTD (strict), unless explicitly true.
	if ap, ok := jtd["additionalProperties"].(bool); ok && ap {
		result["additionalProperties"] = true
	} else if !hasKey(jtd, "additionalProperties") {
		result["additionalProperties"] = false
	}
}

// applyType maps a JTD type name to JSON Schema type constraints.
func applyType(result map[string]interface{}, typeName string) {
	switch typeName {
	case "boolean":
		result["type"] = "boolean"
	case "string":
		result["type"] = "string"
	case "timestamp":
		result["type"] = "string"
		result["format"] = "date-time"
	case "float32", "float64":
		result["type"] = "number"
	case "int8":
		result["type"] = "integer"
		result["minimum"] = float64(-128)
		result["maximum"] = float64(127)
	case "int16":
		result["type"] = "integer"
		result["minimum"] = float64(-32768)
		result["maximum"] = float64(32767)
	case "int32":
		result["type"] = "integer"
		result["minimum"] = float64(-2147483648)
		result["maximum"] = float64(2147483647)
	case "uint8":
		result["type"] = "integer"
		result["minimum"] = float64(0)
		result["maximum"] = float64(255)
	case "uint16":
		result["type"] = "integer"
		result["minimum"] = float64(0)
		result["maximum"] = float64(65535)
	case "uint32":
		result["type"] = "integer"
		result["minimum"] = float64(0)
		result["maximum"] = float64(4294967295)
	default:
		// Unknown type — pass through as string.
		result["type"] = "string"
	}
}

// applyMetadata converts JTD metadata to JSON Schema description and
// x-docsis-* extension keywords.
func applyMetadata(result map[string]interface{}, meta map[string]interface{}) {
	// Build description from explicit description or synthesized from metadata.
	if desc, ok := meta["description"].(string); ok {
		result["description"] = desc
	} else {
		result["description"] = synthesizeDescription(meta)
	}

	// Map known metadata fields to x-docsis-* extensions.
	extensionFields := []string{
		"tlvType", "tlvLength", "dataType", "spec",
		"repeatable", "encoding", "validValues", "context",
	}
	for _, field := range extensionFields {
		if val, ok := meta[field]; ok {
			result["x-docsis-"+field] = val
		}
	}

	// Map non-standard metadata fields that aren't in our known list.
	knownFields := map[string]bool{
		"description": true, "tlvType": true, "tlvLength": true,
		"dataType": true, "spec": true, "repeatable": true,
		"encoding": true, "validValues": true, "context": true,
		"title": true, "version": true, "docsis": true,
	}
	for k, v := range meta {
		if !knownFields[k] {
			result["x-docsis-"+k] = v
		}
	}

	// Preserve title and version at the top level.
	if title, ok := meta["title"].(string); ok {
		result["title"] = title
	}
	if version, ok := meta["version"].(string); ok {
		result["x-docsis-version"] = version
	}
}

// synthesizeDescription creates a description string from metadata fields
// when no explicit description is provided.
func synthesizeDescription(meta map[string]interface{}) string {
	desc := ""
	if tlvType, ok := meta["tlvType"]; ok {
		desc = fmt.Sprintf("TLV %v", tlvType)
	}
	if dataType, ok := meta["dataType"].(string); ok {
		if desc != "" {
			desc += " - " + dataType
		} else {
			desc = dataType
		}
	}
	if spec, ok := meta["spec"].(string); ok {
		if desc != "" {
			desc += " (" + spec + ")"
		} else {
			desc = spec
		}
	}
	if desc == "" {
		if title, ok := meta["title"].(string); ok {
			return title
		}
	}
	return desc
}

// wrapNullable wraps a JSON Schema node in an anyOf with null.
func wrapNullable(schema map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"anyOf": []interface{}{
			schema,
			map[string]interface{}{"type": "null"},
		},
	}
}

// hasKey returns true if the map contains the given key.
func hasKey(m map[string]interface{}, key string) bool {
	_, ok := m[key]
	return ok
}

// sortedKeys returns the keys of a map in sorted order.
func sortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
