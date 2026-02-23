// Command jtd2jsonschema converts a JSON Type Definition (RFC 8927) schema
// to JSON Schema 2020-12. It reads a JTD file from the given path (or stdin)
// and writes the equivalent JSON Schema to stdout.
//
// Usage:
//
//	jtd2jsonschema <input.jtd.json> [> output.schema.json]
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

	if len(os.Args) > 1 {
		input, err = os.ReadFile(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading file %s: %v\n", os.Args[1], err)
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

	output, err := json.MarshalIndent(jsonSchema, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshaling JSON Schema: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(output))
}
