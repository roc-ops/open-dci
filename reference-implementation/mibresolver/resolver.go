// Package mibresolver loads MIB definitions using gosmi and resolves OID
// dotted notation to NET-SNMP style names and integer values to enum labels.
// It supports the OpenDCI mibs/ directory layout where versioned files are
// named MIB-NAME@YYYY-MM-DD.mib and symlinks MIB-NAME.mib point to the
// latest version.
package mibresolver

import (
	"embed"
	"fmt"
	"os"
	"strings"

	"github.com/sleepinggenius2/gosmi"
	"github.com/sleepinggenius2/gosmi/types"
)

// coreMIBs embeds correct versions of foundational MIB files (SNMPv2-SMI,
// SNMPv2-TC, SNMPv2-CONF, RFC1155-SMI) that the repository's copies may be
// incomplete excerpts from RFCs. gosmi needs these to parse virtually all
// other MIBs.
//
//go:embed coremib/*.mib
var coreMIBs embed.FS

// Resolver loads MIBs and resolves OIDs to names and enum values.
type Resolver struct {
	tempDir string // temporary directory for core MIBs and version overrides (native only)
}

// config holds resolver configuration set via Option functions.
type config struct {
	versionOverrides map[string]string // MIB-NAME -> "MIB-NAME@YYYY-MM-DD"
}

// Option configures the resolver.
type Option func(*config)

// WithVersionOverrides sets specific MIB version overrides.
// Format: ["DOCS-IF3-MIB@2024-07-05", "DOCS-QOS3-MIB@2023-11-22"]
func WithVersionOverrides(overrides []string) Option {
	return func(c *config) {
		for _, o := range overrides {
			parts := strings.SplitN(o, "@", 2)
			if len(parts) == 2 {
				c.versionOverrides[parts[0]] = o
			}
		}
	}
}

// ResolveOID takes a dotted OID string and returns the NET-SNMP style name.
// Returns empty string if OID not found in any loaded MIB.
// Example: "1.3.6.1.2.1.2.2.1.7.1" -> "IF-MIB::ifAdminStatus.1"
func (r *Resolver) ResolveOID(oid string) string {
	if oid == "" {
		return ""
	}

	// Strip leading dot if present (e.g., ".1.3.6.1..." -> "1.3.6.1...")
	oid = strings.TrimPrefix(oid, ".")

	parsedOID, err := types.OidFromString(oid)
	if err != nil {
		return ""
	}

	node, err := gosmi.GetNodeByOID(parsedOID)
	if err != nil {
		return ""
	}

	moduleName := node.GetModule().Name

	// Skip nodes that resolved only to the well-known root (iso, ccitt, etc.).
	// These are placeholder nodes, not real MIB objects.
	if moduleName == "<well-known>" {
		return ""
	}

	// Get the registered OID for this node.
	registeredOID := node.Oid.String()

	// Determine the index suffix by comparing input OID to registered OID.
	suffix := ""
	if len(oid) > len(registeredOID) && oid[:len(registeredOID)] == registeredOID {
		suffix = oid[len(registeredOID):]
	}

	if suffix != "" {
		return fmt.Sprintf("%s::%s%s", moduleName, node.Name, suffix)
	}
	return fmt.Sprintf("%s::%s", moduleName, node.Name)
}

// ResolveEnum takes a dotted OID (of the table column, not including index)
// and an integer value, returns the enum label.
// Returns empty string if no enum defined.
// Example: "1.3.6.1.2.1.2.2.1.7", 1 -> "up(1)"
func (r *Resolver) ResolveEnum(oid string, value int64) string {
	if oid == "" {
		return ""
	}

	oid = strings.TrimPrefix(oid, ".")

	parsedOID, err := types.OidFromString(oid)
	if err != nil {
		return ""
	}

	node, err := gosmi.GetNodeByOID(parsedOID)
	if err != nil {
		return ""
	}

	if node.Type == nil {
		return ""
	}

	if node.Type.Enum == nil {
		return ""
	}

	for _, v := range node.Type.Enum.Values {
		if v.Value == value {
			return fmt.Sprintf("%s(%d)", v.Name, value)
		}
	}

	return ""
}

// ResolveNamed takes a MIB-NAME@date::objectName spec and resolves to OID.
// For future encoding support.
// Example: "DOCS-IF3-MIB@2024-07-05::docsIf3CmStatusValue" -> "1.3.6.1.4.1.4491...."
func (r *Resolver) ResolveNamed(spec string) (string, error) {
	// Parse "MIB-NAME::objectName" or "MIB-NAME@date::objectName"
	parts := strings.SplitN(spec, "::", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid spec format: expected MODULE::object, got %q", spec)
	}

	objectName := parts[1]
	if objectName == "" {
		return "", fmt.Errorf("empty object name in spec %q", spec)
	}

	// Look up the node by name.
	node, err := gosmi.GetNode(objectName)
	if err != nil {
		return "", fmt.Errorf("node %q not found: %w", objectName, err)
	}

	return node.Oid.String(), nil
}

// ResolveFullName takes a numeric dotted OID and returns the full named path
// by resolving each OID prefix to its object name.
// Example: "1.3.6.1.2.1.2.2.1.7" -> "iso.org.dod.internet.mgmt.mib-2.interfaces.ifTable.ifEntry.ifAdminStatus"
// Components that can't be resolved fall back to their numeric value.
func (r *Resolver) ResolveFullName(oid string) (string, error) {
	if oid == "" {
		return "", fmt.Errorf("empty OID")
	}
	oid = strings.TrimPrefix(oid, ".")

	parts := strings.Split(oid, ".")
	names := make([]string, len(parts))

	for i := range parts {
		prefix := strings.Join(parts[:i+1], ".")
		parsedOID, err := types.OidFromString(prefix)
		if err != nil {
			names[i] = parts[i]
			continue
		}
		node, err := gosmi.GetNodeByOID(parsedOID)
		if err != nil || node.Name == "" {
			names[i] = parts[i]
			continue
		}
		names[i] = node.Name
	}

	return strings.Join(names, "."), nil
}

// ResolveToNumericOID takes a name and resolves it to a numeric dotted OID.
// Accepts multiple formats:
//   - "MODULE::objectName" (e.g., "IF-MIB::ifAdminStatus")
//   - Plain object name (e.g., "ifAdminStatus")
//   - Dotted named path — uses the last component (e.g., "iso.org...ifAdminStatus")
//
// Returns the numeric OID string or an error if the name cannot be resolved.
func (r *Resolver) ResolveToNumericOID(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("empty name")
	}

	// MODULE::objectName format — delegate to existing ResolveNamed.
	if strings.Contains(name, "::") {
		return r.ResolveNamed(name)
	}

	// Dotted named path — extract leaf name.
	lookupName := name
	if strings.Contains(name, ".") {
		parts := strings.Split(name, ".")
		lookupName = parts[len(parts)-1]
	}

	node, err := gosmi.GetNode(lookupName)
	if err != nil {
		return "", fmt.Errorf("resolving %q: %w", name, err)
	}
	return node.Oid.String(), nil
}

// Close cleans up gosmi state and temporary files.
func (r *Resolver) Close() {
	gosmi.Exit()
	if r.tempDir != "" {
		os.RemoveAll(r.tempDir)
	}
}
