//go:build !js

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/roc-ops/open-dci/reference-implementation/mibresolver"
)

func main() {
	var (
		inputFile  string
		outputFile string
		cmtsSecret string
		schemaPath string
		mibsDir    string
		withMibs   string
		noMibs     bool
		encode     bool
	)

	flag.StringVar(&inputFile, "i", "", "Input file (default: stdin)")
	flag.StringVar(&inputFile, "input", "", "Input file (default: stdin)")
	flag.StringVar(&outputFile, "o", "", "Output file (default: stdout)")
	flag.StringVar(&outputFile, "output", "", "Output file (default: stdout)")
	flag.StringVar(&cmtsSecret, "cmts-secret", "", "CMTS shared secret for MIC verification/computation")
	flag.StringVar(&schemaPath, "schema", "", "Path to JTD schema file (default: ../schemas/docsis-config.jtd.json)")
	flag.StringVar(&mibsDir, "mibs-dir", "", "Path to mibs/ root directory (default: ../mibs/)")
	flag.StringVar(&withMibs, "with-mibs", "", "Comma-separated version-specific MIBs (e.g. DOCS-IF3-MIB@2024-07-05)")
	flag.BoolVar(&noMibs, "no-mibs", false, "Disable MIB resolution")
	flag.BoolVar(&encode, "encode", false, "Encode mode: convert JSON/JSONC to binary TLV format")
	flag.Parse()

	// Resolve schema path default relative to the executable.
	if schemaPath == "" {
		exe, err := os.Executable()
		if err == nil {
			schemaPath = filepath.Join(filepath.Dir(exe), "..", "schemas", "docsis-config.jtd.json")
		}
		// Fall back to relative path from current directory.
		if _, err := os.Stat(schemaPath); err != nil {
			schemaPath = filepath.Join("..", "schemas", "docsis-config.jtd.json")
		}
	}

	// Read input.
	var inputData []byte
	var err error
	if inputFile == "" || inputFile == "-" {
		inputData, err = io.ReadAll(os.Stdin)
	} else {
		inputData, err = os.ReadFile(inputFile)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	if len(inputData) == 0 {
		fmt.Fprintf(os.Stderr, "Error: empty input\n")
		os.Exit(1)
	}

	// Load registry from schema.
	reg, err := LoadRegistry(schemaPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading schema: %v\n", err)
		os.Exit(1)
	}

	if encode {
		runEncode(inputData, reg, outputFile, cmtsSecret)
		return
	}

	runDecode(inputData, reg, outputFile, cmtsSecret, mibsDir, withMibs, noMibs)
}

// runEncode handles encode mode: JSON/JSONC input -> binary TLV output.
func runEncode(inputData []byte, reg *Registry, outputFile string, cmtsSecret string) {
	// Strip JSONC comments to produce valid JSON.
	jsonStr := StripJSONCComments(string(inputData))

	// Parse JSON into config map.
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &config); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	// Build a DecodeResult for the encoder.
	result := &DecodeResult{
		Config: config,
	}

	// Encode to binary.
	encoded, err := Encode(result, reg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding config: %v\n", err)
		os.Exit(1)
	}

	// Optionally compute and insert MICs.
	if cmtsSecret != "" {
		encoded, err = insertMICs(encoded, cmtsSecret)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error computing MICs: %v\n", err)
			os.Exit(1)
		}
	}

	// Write output.
	if outputFile == "" || outputFile == "-" {
		os.Stdout.Write(encoded)
	} else {
		if err := os.WriteFile(outputFile, encoded, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Wrote %d bytes to %s\n", len(encoded), outputFile)
	}
}

// insertMICs computes CM MIC and CMTS MIC and inserts them before the end-of-data
// marker in the encoded binary. The input must end with 0xFF 0x00.
func insertMICs(encoded []byte, cmtsSecret string) ([]byte, error) {
	// Remove end-of-data marker.
	if len(encoded) < 2 || encoded[len(encoded)-2] != 0xFF || encoded[len(encoded)-1] != 0x00 {
		return nil, fmt.Errorf("encoded data does not end with end-of-data marker")
	}
	body := encoded[:len(encoded)-2]

	// Compute CM MIC (TLV 6): MD5 of all TLV bytes excluding TLV 6, 7, 255.
	cmMic := ComputeCmMic(body)
	cmMicTLV := makeTLV(6, cmMic)

	// Build data with CM MIC for CMTS MIC computation.
	withCmMic := append(body, cmMicTLV...)

	// Compute CMTS MIC (TLV 7): HMAC-MD5 of all TLV bytes excluding TLV 7, 255.
	cmtsMic := ComputeCmtsMic(withCmMic, cmtsSecret)
	cmtsMicTLV := makeTLV(7, cmtsMic)

	// Reassemble: body + CM MIC + CMTS MIC + end-of-data.
	var result []byte
	result = append(result, body...)
	result = append(result, cmMicTLV...)
	result = append(result, cmtsMicTLV...)
	result = append(result, 0xFF, 0x00)

	return result, nil
}

// runDecode handles decode mode: binary TLV input -> JSONC output.
func runDecode(inputData []byte, reg *Registry, outputFile string, cmtsSecret string, mibsDir string, withMibs string, noMibs bool) {
	validValues := reg.ValidValuesMap()

	// Initialize MIB resolver (graceful degradation: nil resolver if unavailable).
	var resolver *mibresolver.Resolver
	if !noMibs {
		resolvedMibsDir := mibsDir
		if resolvedMibsDir == "" {
			// Default relative to the executable.
			exe, exeErr := os.Executable()
			if exeErr == nil {
				resolvedMibsDir = filepath.Join(filepath.Dir(exe), "..", "mibs")
			}
			// Fall back to relative path from current directory.
			if resolvedMibsDir == "" {
				resolvedMibsDir = filepath.Join("..", "mibs")
			} else if _, statErr := os.Stat(resolvedMibsDir); statErr != nil {
				resolvedMibsDir = filepath.Join("..", "mibs")
			}
		}

		var opts []mibresolver.Option
		if withMibs != "" {
			overrides := strings.Split(withMibs, ",")
			opts = append(opts, mibresolver.WithVersionOverrides(overrides))
		}

		r, mibErr := mibresolver.New(resolvedMibsDir, opts...)
		if mibErr != nil {
			fmt.Fprintf(os.Stderr, "MIB resolution disabled: %v\n", mibErr)
		} else {
			resolver = r
			defer resolver.Close()
		}
	}

	// Decode the config.
	result, err := Decode(inputData, reg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding config: %v\n", err)
		os.Exit(1)
	}

	// Collect JSONC comment lines for non-schema-conforming items.
	var comments []string

	// Verify MICs and report to stderr (kept as-is) and collect as JSONC comments.
	if result.CmMic != nil {
		micResult := VerifyCmMic(inputData, result.CmMic)
		if micResult.Valid {
			fmt.Fprintf(os.Stderr, "CM MIC: VALID\n")
			comments = append(comments, "// CM MIC: VALID")
		} else {
			msg := fmt.Sprintf("CM MIC: INVALID (expected %X, computed %X)",
				micResult.Expected, micResult.Computed)
			fmt.Fprintf(os.Stderr, "%s\n", msg)
			comments = append(comments, "// "+msg)
		}
		// Add CmMic value as a commented-out property.
		comments = append(comments, fmt.Sprintf("// \"CmMic\": \"%X\",", result.CmMic))
	}

	if result.CmtsMic != nil {
		// Add CmtsMic value as a commented-out property.
		comments = append(comments, fmt.Sprintf("// \"CmtsMic\": \"%X\",", result.CmtsMic))

		if cmtsSecret == "" {
			fmt.Fprintf(os.Stderr, "CMTS MIC: SKIPPED (no --cmts-secret provided)\n")
			comments = append(comments, "// CMTS MIC: SKIPPED (no --cmts-secret provided)")
		} else {
			micResult := VerifyCmtsMic(inputData, result.CmtsMic, cmtsSecret)
			if micResult.Valid {
				fmt.Fprintf(os.Stderr, "CMTS MIC: VALID\n")
				comments = append(comments, "// CMTS MIC: VALID")
			} else {
				msg := fmt.Sprintf("CMTS MIC: INVALID (expected %X, computed %X)",
					micResult.Expected, micResult.Computed)
				fmt.Fprintf(os.Stderr, "%s\n", msg)
				comments = append(comments, "// "+msg)
			}
		}
	}

	// Extract UnknownTlvs from config before marshaling — they become JSONC comments.
	if unknowns, ok := result.Config["UnknownTlvs"]; ok {
		delete(result.Config, "UnknownTlvs")

		// Marshal the UnknownTlvs array to pretty JSON, then comment each line.
		utJSON, err := json.MarshalIndent(unknowns, "", "  ")
		if err == nil {
			lines := strings.Split(string(utJSON), "\n")
			comments = append(comments, fmt.Sprintf("// \"UnknownTlvs\": %s", lines[0]))
			for _, line := range lines[1:] {
				comments = append(comments, "// "+line)
			}
		}
	}

	// Remove CmMic and CmtsMic from config — they were added by the decoder
	// as schema-resolved properties, but we emit them only as comments.
	delete(result.Config, "CmMic")
	delete(result.Config, "CmtsMic")

	// Strip internal ordering metadata before JSON output.
	StripTLVOrder(result.Config)

	// Format as JSONC.
	jsoncData, err := FormatJSONC(result.Config, comments, validValues, resolver)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSONC: %v\n", err)
		os.Exit(1)
	}

	// Write output.
	if outputFile == "" || outputFile == "-" {
		fmt.Println(jsoncData)
	} else {
		// Default extension: .jsonc
		if filepath.Ext(outputFile) == "" {
			outputFile += ".jsonc"
		}
		if err := os.WriteFile(outputFile, []byte(jsoncData+"\n"), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
			os.Exit(1)
		}
	}
}
