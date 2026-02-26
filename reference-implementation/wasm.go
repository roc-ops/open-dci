//go:build js && wasm

package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"syscall/js"

	"github.com/roc-ops/open-dci/reference-implementation/mibresolver"
)

// registry holds the loaded TLV registry for use across WASM function calls.
var registry *Registry

// resolver holds the MIB resolver for OID annotation in decode output.
var resolver *mibresolver.Resolver

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

	// Collect JSONC comment lines for MIC verification (mirrors CLI behavior).
	var comments []string

	// Extract optional shared secret for MIC verification.
	var cmtsSecret string
	if len(args) >= 2 {
		cmtsSecret = args[1].String()
	}

	if result.CmMic != nil {
		micResult := VerifyCmMic(data, result.CmMic)
		if micResult.Valid {
			comments = append(comments, "// CM MIC: VALID")
		} else {
			comments = append(comments, fmt.Sprintf("// CM MIC: INVALID (expected %X, computed %X)",
				micResult.Expected, micResult.Computed))
		}
		comments = append(comments, fmt.Sprintf("// \"CmMic\": \"%X\",", result.CmMic))
	}

	if result.CmtsMic != nil {
		comments = append(comments, fmt.Sprintf("// \"CmtsMic\": \"%X\",", result.CmtsMic))
		if cmtsSecret == "" {
			comments = append(comments, "// CMTS MIC: SKIPPED (no --cmts-secret provided)")
		} else {
			micResult := VerifyCmtsMic(data, result.CmtsMic, cmtsSecret)
			if micResult.Valid {
				comments = append(comments, "// CMTS MIC: VALID")
			} else {
				comments = append(comments, fmt.Sprintf("// CMTS MIC: INVALID (expected %X, computed %X)",
					micResult.Expected, micResult.Computed))
			}
		}
	}

	// Extract UnknownTlvs from config — they become JSONC comments.
	if unknowns, ok := result.Config["UnknownTlvs"]; ok {
		delete(result.Config, "UnknownTlvs")
		utJSON, jsonErr := json.MarshalIndent(unknowns, "", "  ")
		if jsonErr == nil {
			lines := strings.Split(string(utJSON), "\n")
			comments = append(comments, fmt.Sprintf("// \"UnknownTlvs\": %s", lines[0]))
			for _, line := range lines[1:] {
				comments = append(comments, "// "+line)
			}
		}
	}

	// Remove MIC values from config (they are verification-only data).
	delete(result.Config, "CmMic")
	delete(result.Config, "CmtsMic")

	// Format as JSONC with MIB resolver if available.
	validValues := registry.ValidValuesMap()
	jsoncData, err := FormatJSONC(result.Config, comments, validValues, resolver)
	if err != nil {
		return jsError("formatting output: " + err.Error())
	}

	obj := js.Global().Get("Object").New()
	obj.Set("result", jsoncData)
	return obj
}

// opendciEncode encodes a JSON/JSONC string to binary DOCSIS config (Uint8Array).
// JS signature: opendciEncode(jsoncString: string, secret?: string, pad?: boolean, packetCableHash?: string) -> {result: Uint8Array} | {error: string}
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

	// Optionally compute and insert MICs if a shared secret is provided.
	if len(args) >= 2 {
		secret := args[1].String()
		if secret != "" {
			encoded, err = insertMICs(encoded, secret)
			if err != nil {
				return jsError("computing MICs: " + err.Error())
			}
		}
	}

	// Optionally compute and insert PacketCable hash.
	if len(args) >= 4 && args[3].Type() == js.TypeString {
		hashVariant := args[3].String()
		if hashVariant != "" {
			encoded, err = insertPacketCableHash(encoded, hashVariant)
			if err != nil {
				return jsError("computing PacketCable hash: " + err.Error())
			}
		}
	}

	// Optionally pad to 4-byte alignment.
	if len(args) >= 3 && args[2].Type() == js.TypeBoolean && args[2].Bool() {
		encoded = PadToAlignment(encoded, 4)
	}

	// Copy the encoded bytes into a JS Uint8Array.
	jsArray := js.Global().Get("Uint8Array").New(len(encoded))
	js.CopyBytesToJS(jsArray, encoded)

	obj := js.Global().Get("Object").New()
	obj.Set("result", jsArray)
	return obj
}

// opendciInitMIBs initializes the MIB resolver with embedded core MIBs only.
// Call this before opendciDecode to enable OID name and enum annotations.
// JS signature: opendciInitMIBs() -> {ok: true, loaded: number} | {error: string}
func opendciInitMIBs(_ js.Value, args []js.Value) interface{} {
	// Close existing resolver if any.
	if resolver != nil {
		resolver.Close()
		resolver = nil
	}

	r, err := mibresolver.NewFromMIBData(nil)
	if err != nil {
		return jsError("initializing MIBs: " + err.Error())
	}
	resolver = r

	obj := js.Global().Get("Object").New()
	obj.Set("ok", true)
	return obj
}

// opendciLoadMIBs loads additional MIB files into the resolver for OID resolution.
// If the resolver hasn't been initialized, it will be initialized with core MIBs first.
// JS signature: opendciLoadMIBs(mibFiles: {[filename: string]: string}) -> {ok: true, loaded: number} | {error: string}
func opendciLoadMIBs(_ js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return jsError("opendciLoadMIBs requires 1 argument: object mapping filenames to content strings")
	}

	jsObj := args[0]
	if jsObj.Type() != js.TypeObject {
		return jsError("opendciLoadMIBs argument must be an object {filename: content}")
	}

	// Extract {filename: content} pairs from the JS object.
	keys := js.Global().Get("Object").Call("keys", jsObj)
	keysLen := keys.Get("length").Int()

	mibFiles := make(map[string][]byte, keysLen)
	for i := 0; i < keysLen; i++ {
		key := keys.Index(i).String()
		value := jsObj.Get(key).String()
		mibFiles[key] = []byte(value)
	}

	// Initialize resolver if not yet done.
	if resolver == nil {
		r, err := mibresolver.NewFromMIBData(mibFiles)
		if err != nil {
			return jsError("loading MIBs: " + err.Error())
		}
		resolver = r

		obj := js.Global().Get("Object").New()
		obj.Set("ok", true)
		obj.Set("loaded", keysLen)
		return obj
	}

	// Incremental load into existing resolver.
	loaded, err := resolver.LoadAdditionalMIBs(mibFiles)
	if err != nil {
		return jsError("loading additional MIBs: " + err.Error())
	}

	obj := js.Global().Get("Object").New()
	obj.Set("ok", true)
	obj.Set("loaded", loaded)
	return obj
}

// opendciQueryMIBTree returns the entire loaded MIB tree as a JSON string.
// The tree is rooted at OID "1" (iso) with recursive children.
// JS signature: opendciQueryMIBTree() -> {result: string} | {error: string}
func opendciQueryMIBTree(_ js.Value, args []js.Value) interface{} {
	if resolver == nil {
		return jsError("no MIBs loaded: call opendciInitMIBs first")
	}
	tree, err := resolver.QueryTree()
	if err != nil {
		return jsError("querying MIB tree: " + err.Error())
	}
	jsonBytes, err := json.Marshal(tree)
	if err != nil {
		return jsError("serializing MIB tree: " + err.Error())
	}
	obj := js.Global().Get("Object").New()
	obj.Set("result", string(jsonBytes))
	return obj
}

// opendciResolveName resolves a numeric dotted OID to a full named path.
// JS signature: opendciResolveName(numericOid: string) -> {result: string} | {error: string}
func opendciResolveName(_ js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return jsError("opendciResolveName requires 1 argument: numeric OID string")
	}
	if resolver == nil {
		return jsError("no MIBs loaded: call opendciInitMIBs first")
	}
	result, err := resolver.ResolveFullName(args[0].String())
	if err != nil {
		return jsError("resolving name: " + err.Error())
	}
	obj := js.Global().Get("Object").New()
	obj.Set("result", result)
	return obj
}

// opendciResolveOID resolves a name to a numeric dotted OID.
// Accepts "MODULE::objectName", plain object name, or dotted named path.
// JS signature: opendciResolveOID(name: string) -> {result: string} | {error: string}
func opendciResolveOID(_ js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return jsError("opendciResolveOID requires 1 argument: name string")
	}
	if resolver == nil {
		return jsError("no MIBs loaded: call opendciInitMIBs first")
	}
	result, err := resolver.ResolveToNumericOID(args[0].String())
	if err != nil {
		return jsError("resolving OID: " + err.Error())
	}
	obj := js.Global().Get("Object").New()
	obj.Set("result", result)
	return obj
}

// opendciExtractCVC extracts CVC certificates from a PKCS#7-signed CM firmware binary.
// JS signature: opendciExtractCVC(firmwareBinary: Uint8Array) -> {result: {ManufacturerCvc, CoSignerCvc, ManufacturerCvcChain, CoSignerCvcChain}} | {error: string}
func opendciExtractCVC(_ js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return jsError("opendciExtractCVC requires 1 argument: Uint8Array")
	}

	// Copy binary data from JS Uint8Array to Go slice.
	jsArray := args[0]
	length := jsArray.Get("length").Int()
	data := make([]byte, length)
	js.CopyBytesToGo(data, jsArray)

	certs, err := ExtractCVCFromFirmware(data)
	if err != nil {
		return jsError("extracting CVCs: " + err.Error())
	}

	// Build the result JS object with null for absent certs.
	resultObj := js.Global().Get("Object").New()
	for _, key := range []string{"ManufacturerCvc", "CoSignerCvc", "ManufacturerCvcChain", "CoSignerCvcChain"} {
		if v, ok := certs[key]; ok && v != nil {
			resultObj.Set(key, v.(string))
		} else {
			resultObj.Set(key, js.Null())
		}
	}

	obj := js.Global().Get("Object").New()
	obj.Set("result", resultObj)
	return obj
}

func main() {
	js.Global().Set("opendciLoadSchema", js.FuncOf(opendciLoadSchema))
	js.Global().Set("opendciDecode", js.FuncOf(opendciDecode))
	js.Global().Set("opendciEncode", js.FuncOf(opendciEncode))
	js.Global().Set("opendciInitMIBs", js.FuncOf(opendciInitMIBs))
	js.Global().Set("opendciLoadMIBs", js.FuncOf(opendciLoadMIBs))
	js.Global().Set("opendciQueryMIBTree", js.FuncOf(opendciQueryMIBTree))
	js.Global().Set("opendciResolveName", js.FuncOf(opendciResolveName))
	js.Global().Set("opendciResolveOID", js.FuncOf(opendciResolveOID))
	js.Global().Set("opendciExtractCVC", js.FuncOf(opendciExtractCVC))

	// Block forever to keep the Go runtime alive for JS callbacks.
	select {}
}
