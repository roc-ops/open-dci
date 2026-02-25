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
	)

	flag.StringVar(&inputFile, "i", "", "Input binary config file (default: stdin)")
	flag.StringVar(&inputFile, "input", "", "Input binary config file (default: stdin)")
	flag.StringVar(&outputFile, "o", "", "Output JSONC file (default: stdout)")
	flag.StringVar(&outputFile, "output", "", "Output JSONC file (default: stdout)")
	flag.StringVar(&cmtsSecret, "cmts-secret", "", "CMTS shared secret for MIC verification")
	flag.StringVar(&schemaPath, "schema", "", "Path to JTD schema file (default: ../schemas/docsis-config.jtd.json)")
	flag.StringVar(&mibsDir, "mibs-dir", "", "Path to mibs/ root directory (default: ../mibs/)")
	flag.StringVar(&withMibs, "with-mibs", "", "Comma-separated version-specific MIBs (e.g. DOCS-IF3-MIB@2024-07-05)")
	flag.BoolVar(&noMibs, "no-mibs", false, "Disable MIB resolution")
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
