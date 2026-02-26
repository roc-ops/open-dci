// Command jtd2jsonschema converts a JSON Type Definition (RFC 8927) schema
// to JSON Schema 2020-12. It reads a JTD file from the given path (or stdin)
// and writes the equivalent JSON Schema to stdout.
//
// Usage:
//
//	jtd2jsonschema [--vendors-dir <dir>] <input.jtd.json> [> output.schema.json]
//	cat input.jtd.json | jtd2jsonschema
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func main() {
	var input []byte
	var err error
	var vendorsDir string

	// Parse arguments: support --vendors-dir flag before the positional input path.
	args := os.Args[1:]
	var positionalArgs []string
	for i := 0; i < len(args); i++ {
		if args[i] == "--vendors-dir" && i+1 < len(args) {
			vendorsDir = args[i+1]
			i++ // skip the value
		} else {
			positionalArgs = append(positionalArgs, args[i])
		}
	}

	if len(positionalArgs) > 0 {
		input, err = os.ReadFile(positionalArgs[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading file %s: %v\n", positionalArgs[0], err)
			os.Exit(1)
		}
	} else {
		input, err = io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading stdin: %v\n", err)
			os.Exit(1)
		}
	}

	var jtd map[string]interface{}
	if err := json.Unmarshal(input, &jtd); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing JTD: %v\n", err)
		os.Exit(1)
	}

	jsonSchema := ConvertJTDToJSONSchema(jtd)

	// If --vendors-dir is specified, merge vendor schemas into the output.
	if vendorsDir != "" {
		if err := MergeVendorSchemas(jsonSchema, vendorsDir); err != nil {
			fmt.Fprintf(os.Stderr, "error merging vendor schemas: %v\n", err)
			os.Exit(1)
		}
	}

	output, err := json.MarshalIndent(jsonSchema, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshaling JSON Schema: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(output))
}
