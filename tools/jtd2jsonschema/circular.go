package main

import "strings"

// BreakCircularRefs detects circular $ref chains in the $defs section of a
// JSON Schema and replaces back-edge references with inline
// {"type": "object", "additionalProperties": true} to break the cycle.
// This prevents browser JSON validators (like Monaco Editor) from crashing
// on infinite $ref recursion.
func BreakCircularRefs(schema map[string]interface{}) {
	defs, ok := schema["$defs"].(map[string]interface{})
	if !ok || len(defs) == 0 {
		return
	}

	// For each definition, simulate $ref resolution with a visited stack.
	// When following a $ref would revisit a definition already in the
	// resolution chain, replace that $ref with an inline object.
	for name := range defs {
		visited := map[string]bool{name: true}
		if defMap, ok := defs[name].(map[string]interface{}); ok {
			breakRefsInNode(defMap, defs, visited)
		}
	}
}

// breakRefsInNode walks a JSON Schema node tree. When it encounters a $ref
// pointing to a definition already in the visited set, it replaces the $ref
// with an inline {"type": "object", "additionalProperties": true}. When a
// $ref is not circular, it follows the reference into the target definition
// to check for deeper cycles.
func breakRefsInNode(node map[string]interface{}, defs map[string]interface{}, visited map[string]bool) {
	for _, val := range node {
		switch child := val.(type) {
		case map[string]interface{}:
			if ref, ok := child["$ref"].(string); ok {
				defName := extractDefName(ref)
				if defName != "" {
					if visited[defName] {
						// Circular reference — replace with inline object.
						delete(child, "$ref")
						child["type"] = "object"
						child["additionalProperties"] = true
					} else if targetDef, ok := defs[defName].(map[string]interface{}); ok {
						// Follow the ref to detect deeper cycles.
						visited[defName] = true
						breakRefsInNode(targetDef, defs, visited)
						delete(visited, defName)
					}
					continue
				}
			}
			breakRefsInNode(child, defs, visited)
		case []interface{}:
			for _, item := range child {
				if itemMap, ok := item.(map[string]interface{}); ok {
					breakRefsInNode(itemMap, defs, visited)
				}
			}
		}
	}
}

// extractDefName extracts the definition name from a $ref like "#/$defs/Foo".
func extractDefName(ref string) string {
	const prefix = "#/$defs/"
	if strings.HasPrefix(ref, prefix) {
		return ref[len(prefix):]
	}
	return ""
}
