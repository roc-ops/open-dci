package opendci

// Config file format identifiers.
const (
	FormatCM  = "cm"
	FormatMTA = "mta"
)

// DetectFormat examines the first byte of binary config data to determine
// the format. MTA configs start with TLV 254 (0xFE); CM configs do not.
func DetectFormat(data []byte) string {
	if len(data) > 0 && data[0] == 0xFE {
		return FormatMTA
	}
	return FormatCM
}
