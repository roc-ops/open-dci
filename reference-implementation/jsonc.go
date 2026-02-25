package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/roc-ops/open-dci/reference-implementation/mibresolver"
)

// propertyLineRe matches JSON property lines with integer values, e.g.:
//
//	"NetworkAccess": 1,
var propertyLineRe = regexp.MustCompile(`^(\s*)"([^"]+)":\s*(\d+)(,?)$`)

// FormatJSONC marshals config as pretty-printed JSON, then appends the given
// comment lines (each prefixed with "// ") just before the closing "}".
// If comments is empty and validValues is nil the output is plain JSON.
// When validValues is provided, matching integer property values get an
// inline "// label" comment appended to the line.
// When resolver is non-nil, SNMP OID lines get "// MIB::name" comments and
// integer value lines within varbinds get "// enumLabel" comments.
func FormatJSONC(config map[string]interface{}, comments []string, validValues map[string]map[string]string, resolver *mibresolver.Resolver) (string, error) {
	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return "", fmt.Errorf("encoding JSON: %w", err)
	}

	s := string(jsonData)
	s = addInlineComments(s, validValues)
	s = addSnmpComments(s, resolver)

	if len(comments) == 0 {
		return s, nil
	}

	// Find the last '}' which closes the top-level object.
	lastBrace := strings.LastIndex(s, "}")
	if lastBrace < 0 {
		return "", fmt.Errorf("unexpected JSON structure: no closing brace")
	}

	// Everything before the closing brace (includes the trailing newline after
	// the last property, if any).
	before := s[:lastBrace]

	// Ensure there is a newline before the comment block. For an empty object
	// json.MarshalIndent produces "{}" with no newline between the braces.
	if !strings.HasSuffix(before, "\n") {
		before += "\n"
	}

	// Build the comment block. Each line is indented with two spaces.
	var commentBlock strings.Builder
	for _, c := range comments {
		commentBlock.WriteString("  " + c + "\n")
	}

	return before + commentBlock.String() + "}", nil
}

// addInlineComments appends "// label" comments to JSON property lines whose
// integer value has a matching entry in validValues.
func addInlineComments(jsonStr string, validValues map[string]map[string]string) string {
	if len(validValues) == 0 {
		return jsonStr
	}
	lines := strings.Split(jsonStr, "\n")
	for i, line := range lines {
		m := propertyLineRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		propName := m[2]
		valueStr := m[3]
		vvMap, ok := validValues[propName]
		if !ok {
			continue
		}
		label, ok := vvMap[valueStr]
		if !ok {
			continue
		}
		lines[i] = line + " // " + label
	}
	return strings.Join(lines, "\n")
}

// oidLineRe matches JSONC lines like:  "oid": "1.3.6.1.2.1...",
var oidLineRe = regexp.MustCompile(`^(\s*)"oid":\s*"([^"]+)"(,?)$`)

// valueLineRe matches JSONC lines like:  "value": "42",
var valueLineRe = regexp.MustCompile(`^(\s*)"value":\s*"([^"]*)"(,?)$`)

// snmpBlockRe matches the beginning of a SnmpMibObject array.
var snmpBlockRe = regexp.MustCompile(`^\s*"SnmpMibObject":\s*\[`)

// addSnmpComments annotates OID and value lines inside SnmpMibObject arrays.
// For "oid" lines it appends "// MODULE::objectName.index".
// For "value" lines containing an integer it attempts enum resolution,
// stripping the last OID component (the SNMP table index) to obtain the
// column OID that carries the enum definition.
func addSnmpComments(jsonStr string, resolver *mibresolver.Resolver) string {
	if resolver == nil {
		return jsonStr
	}

	lines := strings.Split(jsonStr, "\n")
	inSnmpBlock := false
	braceDepth := 0
	var currentOID string

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Detect entry into SnmpMibObject array.
		if snmpBlockRe.MatchString(line) {
			inSnmpBlock = true
			braceDepth = 1 // account for the '[' on this line (skipped by continue)
			currentOID = ""
			continue
		}

		if !inSnmpBlock {
			continue
		}

		// Track brace/bracket depth to know when we leave the array.
		for _, ch := range trimmed {
			switch ch {
			case '{':
				braceDepth++
			case '}':
				braceDepth--
			case '[':
				braceDepth++
			case ']':
				braceDepth--
			}
		}

		// If depth drops to zero or below, we exited the SnmpMibObject array.
		if braceDepth <= 0 {
			inSnmpBlock = false
			currentOID = ""
			continue
		}

		// Match "oid" lines.
		if m := oidLineRe.FindStringSubmatch(line); m != nil {
			oid := m[2]
			currentOID = oid
			name := resolver.ResolveOID(oid)
			if name != "" {
				lines[i] = m[1] + `"oid": "` + oid + `"` + m[3] + " // " + name
			}
			continue
		}

		// Match "value" lines — try enum resolution.
		if m := valueLineRe.FindStringSubmatch(line); m != nil {
			val := m[2]
			if currentOID != "" {
				intVal, err := strconv.ParseInt(val, 10, 64)
				if err == nil {
					// Strip the last OID component (the table index) to get
					// the column OID where the enum is defined.
					baseOID := stripLastOIDComponent(currentOID)
					if baseOID != "" {
						label := resolver.ResolveEnum(baseOID, intVal)
						if label != "" {
							lines[i] = m[1] + `"value": "` + val + `"` + m[3] + " // " + label
						}
					}
				}
			}
			continue
		}
	}

	return strings.Join(lines, "\n")
}

// stripLastOIDComponent removes the last dotted component from an OID string.
// "1.3.6.1.2.1.2.2.1.7.1" -> "1.3.6.1.2.1.2.2.1.7"
func stripLastOIDComponent(oid string) string {
	idx := strings.LastIndex(oid, ".")
	if idx <= 0 {
		return ""
	}
	return oid[:idx]
}

// StripJSONCComments removes single-line // comments from JSONC text, producing
// valid JSON. It handles comments at the end of lines (inline comments) and
// full-line comments. It is careful not to strip // inside quoted strings.
func StripJSONCComments(input string) string {
	lines := strings.Split(input, "\n")
	var result []string

	for _, line := range lines {
		stripped := stripLineComment(line)
		// Skip lines that are entirely comments (now empty or whitespace-only).
		trimmed := strings.TrimSpace(stripped)
		if trimmed == "" {
			continue
		}
		result = append(result, stripped)
	}

	return strings.Join(result, "\n")
}

// stripLineComment removes the // comment portion from a single line,
// respecting quoted strings. Returns the line with the comment removed.
func stripLineComment(line string) string {
	inString := false
	escape := false

	for i, ch := range line {
		if escape {
			escape = false
			continue
		}
		if ch == '\\' && inString {
			escape = true
			continue
		}
		if ch == '"' {
			inString = !inString
			continue
		}
		if !inString && ch == '/' && i+1 < len(line) && line[i+1] == '/' {
			return strings.TrimRight(line[:i], " \t")
		}
	}

	return line
}
