package opendci

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestFormatJSONC_NoComments(t *testing.T) {
	config := map[string]interface{}{
		"NetworkAccess": 1,
		"MaxNumCpes":    16,
	}

	out, err := FormatJSONC(config, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Output has a header comment; stripping comments should produce valid JSON.
	stripped := StripJSONCComments(out)
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(stripped), &parsed); err != nil {
		t.Fatalf("expected valid JSON after stripping comments, got error: %v\noutput:\n%s", err, out)
	}

	if int(parsed["NetworkAccess"].(float64)) != 1 {
		t.Errorf("expected NetworkAccess=1, got %v", parsed["NetworkAccess"])
	}
	if int(parsed["MaxNumCpes"].(float64)) != 16 {
		t.Errorf("expected MaxNumCpes=16, got %v", parsed["MaxNumCpes"])
	}

	// Header comment should be present.
	if !strings.Contains(out, "// OpenDCI v"+Version) {
		t.Errorf("expected header comment in output, got:\n%s", out)
	}
}

func TestFormatJSONC_WithMICComments(t *testing.T) {
	config := map[string]interface{}{
		"NetworkAccess": 1,
	}

	comments := []string{
		"// CM MIC: INVALID (expected AABB, computed CCDD)",
		"// \"CmMic\": \"AABB\",",
		"// \"CmtsMic\": \"EEFF\",",
		"// CMTS MIC: SKIPPED (no --cmts-secret provided)",
	}

	out, err := FormatJSONC(config, comments, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Verify schema property is present as regular JSON.
	if !strings.Contains(out, `"NetworkAccess": 1`) {
		t.Errorf("expected NetworkAccess as regular property, got:\n%s", out)
	}

	// Verify each comment line is present (indented with 2 spaces).
	for _, c := range comments {
		indented := "  " + c
		if !strings.Contains(out, indented) {
			t.Errorf("expected comment %q in output, got:\n%s", indented, out)
		}
	}

	// Comments should appear before the closing brace.
	closingIdx := strings.LastIndex(out, "}")
	for _, c := range comments {
		idx := strings.Index(out, c)
		if idx >= closingIdx {
			t.Errorf("comment %q should appear before closing brace", c)
		}
	}

	// The output should start with the header comment and end with }
	trimmed := strings.TrimSpace(out)
	if !strings.HasPrefix(trimmed, "// OpenDCI") || !strings.HasSuffix(trimmed, "}") {
		t.Errorf("expected output to start with header comment and end with }, got:\n%s", out)
	}
}

func TestFormatJSONC_WithUnknownTlvs(t *testing.T) {
	config := map[string]interface{}{
		"NetworkAccess": 1,
	}

	comments := []string{
		"// \"UnknownTlvs\": [",
		"//   {",
		"//     \"type\": 47,",
		"//     \"value\": \"0602000007020016\"",
		"//   }",
		"// ]",
	}

	out, err := FormatJSONC(config, comments, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Verify UnknownTlvs appear as comments, not as properties.
	for _, c := range comments {
		indented := "  " + c
		if !strings.Contains(out, indented) {
			t.Errorf("expected comment %q in output, got:\n%s", indented, out)
		}
	}

	t.Logf("JSONC output:\n%s", out)
}

func TestFormatJSONC_EmptyConfigWithComments(t *testing.T) {
	config := map[string]interface{}{}

	comments := []string{
		"// CM MIC: VALID",
		"// \"CmMic\": \"AABBCCDD\",",
	}

	out, err := FormatJSONC(config, comments, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Should still produce valid structure with header and braces.
	trimmed := strings.TrimSpace(out)
	if !strings.HasPrefix(trimmed, "// OpenDCI") || !strings.HasSuffix(trimmed, "}") {
		t.Errorf("expected output to start with header and end with }, got:\n%s", out)
	}

	// Comments should be present.
	for _, c := range comments {
		if !strings.Contains(out, c) {
			t.Errorf("expected comment %q in output, got:\n%s", c, out)
		}
	}

	t.Logf("JSONC output:\n%s", out)
}

func TestFormatJSONC_SchemaPropertiesRemain(t *testing.T) {
	config := map[string]interface{}{
		"NetworkAccess": 1,
		"MaxNumCpes":    16,
		"DownstreamServiceFlow": []interface{}{
			map[string]interface{}{
				"MaxSustainedTrafficRate": 500000000,
				"QosParamSetType":         7,
				"ServiceFlowReference":    1,
			},
		},
	}

	comments := []string{
		"// \"CmMic\": \"AABB\",",
	}

	out, err := FormatJSONC(config, comments, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	// All schema properties should be in the output as regular JSON (not commented).
	if !strings.Contains(out, `"NetworkAccess"`) {
		t.Error("missing NetworkAccess property")
	}
	if !strings.Contains(out, `"MaxNumCpes"`) {
		t.Error("missing MaxNumCpes property")
	}
	if !strings.Contains(out, `"DownstreamServiceFlow"`) {
		t.Error("missing DownstreamServiceFlow property")
	}
	if !strings.Contains(out, `"MaxSustainedTrafficRate"`) {
		t.Error("missing MaxSustainedTrafficRate property")
	}

	// CmMic should only appear as a comment.
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.Contains(trimmedLine, "CmMic") && !strings.HasPrefix(trimmedLine, "//") {
			t.Errorf("CmMic should be a comment, found as property: %s", line)
		}
	}

	t.Logf("JSONC output:\n%s", out)
}

func TestFormatJSONC_InlineComments(t *testing.T) {
	config := map[string]interface{}{
		"NetworkAccess": 1,
	}
	vv := map[string]map[string]string{
		"NetworkAccess": {"0": "disabled", "1": "enabled"},
	}

	out, err := FormatJSONC(config, nil, vv, nil)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(out, `"NetworkAccess": 1 // enabled`) {
		t.Errorf("expected inline comment '// enabled' on NetworkAccess line, got:\n%s", out)
	}

	t.Logf("JSONC output:\n%s", out)
}

func TestFormatJSONC_InlineCommentNoMatch(t *testing.T) {
	config := map[string]interface{}{
		"NetworkAccess": 99,
	}
	vv := map[string]map[string]string{
		"NetworkAccess": {"0": "disabled", "1": "enabled"},
	}

	out, err := FormatJSONC(config, nil, vv, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Should have the header comment but no inline enum comment.
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		// Skip header comment lines.
		if strings.HasPrefix(trimmedLine, "// OpenDCI") || strings.HasPrefix(trimmedLine, "// https://") {
			continue
		}
		if strings.Contains(trimmedLine, "//") {
			t.Errorf("expected no inline comment for value 99, got line: %s", line)
		}
	}

	t.Logf("JSONC output:\n%s", out)
}

func TestFormatJSONC_InlineAndBlockComments(t *testing.T) {
	config := map[string]interface{}{
		"NetworkAccess": 1,
	}
	comments := []string{
		"// CM MIC: VALID",
	}
	vv := map[string]map[string]string{
		"NetworkAccess": {"0": "disabled", "1": "enabled"},
	}

	out, err := FormatJSONC(config, comments, vv, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Inline comment on the property line.
	if !strings.Contains(out, `"NetworkAccess": 1 // enabled`) {
		t.Errorf("expected inline comment on NetworkAccess, got:\n%s", out)
	}

	// Block comment before closing brace.
	if !strings.Contains(out, "  // CM MIC: VALID") {
		t.Errorf("expected block comment in output, got:\n%s", out)
	}

	t.Logf("JSONC output:\n%s", out)
}

func TestFormatJSONC_NestedInlineComments(t *testing.T) {
	config := map[string]interface{}{
		"UpstreamServiceFlow": []interface{}{
			map[string]interface{}{
				"QosParamSetType":      7,
				"ServiceFlowReference": 1,
			},
		},
	}
	vv := map[string]map[string]string{
		"QosParamSetType": {"0": "reserved", "7": "provisioned, admitted, and active set"},
	}

	out, err := FormatJSONC(config, nil, vv, nil)
	if err != nil {
		t.Fatal(err)
	}

	// The property may have a trailing comma if it's not the last in the object.
	if !strings.Contains(out, `"QosParamSetType": 7, // provisioned, admitted, and active set`) &&
		!strings.Contains(out, `"QosParamSetType": 7 // provisioned, admitted, and active set`) {
		t.Errorf("expected inline comment on nested QosParamSetType, got:\n%s", out)
	}

	t.Logf("JSONC output:\n%s", out)
}

func TestStripJSONCComments_InlineComment(t *testing.T) {
	input := `{
  "NetworkAccess": 1 // enabled
}`
	expected := `{
  "NetworkAccess": 1
}`
	result := StripJSONCComments(input)
	if result != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}

	// The result should be parseable as JSON.
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("result should be valid JSON: %v", err)
	}
}

func TestStripJSONCComments_FullLineComment(t *testing.T) {
	input := `{
  "NetworkAccess": 1,
  // CM MIC: VALID
  // "CmMic": "AABBCCDD",
}`
	result := StripJSONCComments(input)

	// Full-line comments should be removed entirely.
	if strings.Contains(result, "CM MIC") {
		t.Errorf("full-line comment should be stripped, got:\n%s", result)
	}
	if strings.Contains(result, "CmMic") {
		t.Errorf("commented property should be stripped, got:\n%s", result)
	}
}

func TestStripJSONCComments_QuotedSlashes(t *testing.T) {
	input := `{
  "url": "http://example.com"
}`
	result := StripJSONCComments(input)
	if !strings.Contains(result, "http://example.com") {
		t.Errorf("// inside string should be preserved, got:\n%s", result)
	}
}

func TestStripJSONCComments_NoComments(t *testing.T) {
	input := `{
  "NetworkAccess": 1
}`
	result := StripJSONCComments(input)
	if result != input {
		t.Errorf("expected no change, got:\n%s", result)
	}
}
