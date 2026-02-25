//go:build js && wasm

package main

import (
	"encoding/json"
	"syscall/js"
)

// registry holds the loaded TLV registry for use across WASM function calls.
var registry *Registry

// jsError returns a JS object with an "error" property containing the message.
func jsError(msg string) interface{} {
	obj := js.Global().Get("Object").New()
	obj.Set("error", msg)
	return obj
}

// opendciLoadSchema loads a TLV registry from a JSON schema string.
// JS signature: opendciLoadSchema(schemaJSON: string) -> {ok: true} | {error: string}
func opendciLoadSchema(_ js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return jsError("opendciLoadSchema requires 1 argument: schemaJSON string")
	}

	schemaJSON := args[0].String()
	reg, err := LoadRegistryFromBytes([]byte(schemaJSON))
	if err != nil {
		return jsError("loading schema: " + err.Error())
	}

	registry = reg

	obj := js.Global().Get("Object").New()
	obj.Set("ok", true)
	return obj
}

// opendciDecode decodes a binary DOCSIS config (Uint8Array) to a JSON string.
// JS signature: opendciDecode(binaryData: Uint8Array) -> {result: string} | {error: string}
func opendciDecode(_ js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return jsError("opendciDecode requires 1 argument: Uint8Array")
	}

	if registry == nil {
		return jsError("no schema loaded: call opendciLoadSchema first")
	}

	// Copy binary data from JS Uint8Array to Go slice.
	jsArray := args[0]
	length := jsArray.Get("length").Int()
	data := make([]byte, length)
	js.CopyBytesToGo(data, jsArray)

	// Decode the binary config.
	result, err := Decode(data, registry)
	if err != nil {
		return jsError("decoding: " + err.Error())
	}

	// Strip internal ordering metadata before JSON output.
	StripTLVOrder(result.Config)

	// Remove MIC values from config (they are verification-only data).
	delete(result.Config, "CmMic")
	delete(result.Config, "CmtsMic")

	// Format as JSONC (no MIB resolver in WASM — pass nil).
	validValues := registry.ValidValuesMap()
	jsoncData, err := FormatJSONC(result.Config, nil, validValues, nil)
	if err != nil {
		return jsError("formatting output: " + err.Error())
	}

	obj := js.Global().Get("Object").New()
	obj.Set("result", jsoncData)
	return obj
}

// opendciEncode encodes a JSON/JSONC string to binary DOCSIS config (Uint8Array).
// JS signature: opendciEncode(jsoncString: string) -> {result: Uint8Array} | {error: string}
func opendciEncode(_ js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return jsError("opendciEncode requires 1 argument: JSON/JSONC string")
	}

	if registry == nil {
		return jsError("no schema loaded: call opendciLoadSchema first")
	}

	// Strip JSONC comments to produce valid JSON.
	jsonStr := StripJSONCComments(args[0].String())

	// Parse JSON into config map.
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &config); err != nil {
		return jsError("parsing JSON: " + err.Error())
	}

	// Build a DecodeResult for the encoder.
	result := &DecodeResult{
		Config: config,
	}

	// Encode to binary.
	encoded, err := Encode(result, registry)
	if err != nil {
		return jsError("encoding: " + err.Error())
	}

	// Copy the encoded bytes into a JS Uint8Array.
	jsArray := js.Global().Get("Uint8Array").New(len(encoded))
	js.CopyBytesToJS(jsArray, encoded)

	obj := js.Global().Get("Object").New()
	obj.Set("result", jsArray)
	return obj
}

func main() {
	js.Global().Set("opendciLoadSchema", js.FuncOf(opendciLoadSchema))
	js.Global().Set("opendciDecode", js.FuncOf(opendciDecode))
	js.Global().Set("opendciEncode", js.FuncOf(opendciEncode))

	// Block forever to keep the Go runtime alive for JS callbacks.
	select {}
}
