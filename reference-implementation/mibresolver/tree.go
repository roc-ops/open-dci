package mibresolver

import (
	"fmt"
	"sort"
	"strings"
	"unicode"

	"github.com/sleepinggenius2/gosmi"
	"github.com/sleepinggenius2/gosmi/types"
)

// EnumValue represents a named integer value in a MIB object's SYNTAX clause.
type EnumValue struct {
	Value int64  `json:"value"`
	Label string `json:"label"`
}

// IndexObject represents an index column from a SNMP table's INDEX clause.
type IndexObject struct {
	Name        string `json:"name"`
	OID         string `json:"oid"`
	Module      string `json:"module"`
	Syntax      string `json:"syntax,omitempty"`
	Description string `json:"description,omitempty"`
}

// MIBTreeNode represents a single node in the MIB object tree.
type MIBTreeNode struct {
	OID         string         `json:"oid"`
	Name        string         `json:"name"`
	Module      string         `json:"module"`
	Description string         `json:"description,omitempty"`
	Syntax      string         `json:"syntax,omitempty"`
	Access      string         `json:"access,omitempty"`
	NodeType    string         `json:"nodeType"`
	Indexes     []IndexObject  `json:"indexes,omitempty"`
	Enums       []EnumValue    `json:"enums,omitempty"`
	Children    []*MIBTreeNode `json:"children,omitempty"`
}

// QueryTree returns the full MIB tree as a hierarchical structure rooted at
// the ISO root (OID "1"). Each loaded MIB node is represented as a
// MIBTreeNode with its children nested recursively.
func (r *Resolver) QueryTree() (*MIBTreeNode, error) {
	// Collect all nodes from all loaded modules. We iterate modules and
	// call GetNodes() on each, which returns every node defined in that
	// module. This is more reliable than GetSubtree() from a root node,
	// which may return an empty list for well-known placeholder roots.
	modules := gosmi.GetLoadedModules()
	if len(modules) == 0 {
		return nil, fmt.Errorf("no MIB modules loaded")
	}

	nodeMap := make(map[string]*MIBTreeNode)

	for _, m := range modules {
		nodes := m.GetNodes()
		for _, sn := range nodes {
			tn := smiNodeToTreeNode(sn)
			// Deduplicate: keep the first occurrence (modules may share nodes).
			if _, exists := nodeMap[tn.OID]; !exists {
				nodeMap[tn.OID] = tn
			}
		}
	}

	if len(nodeMap) == 0 {
		return nil, fmt.Errorf("no MIB nodes found in loaded modules")
	}

	// Create a synthetic root for OID "1" if not already present.
	root, ok := nodeMap["1"]
	if !ok {
		root = &MIBTreeNode{
			OID:      "1",
			Name:     "iso",
			Module:   "<well-known>",
			NodeType: "node",
		}
		nodeMap["1"] = root
	}

	// Ensure ancestor nodes exist so the tree is fully connected.
	// For each node, walk up the OID path and create synthetic ancestors
	// as needed (e.g., "1.3" if only "1.3.6" and "1" exist).
	allOIDs := make([]string, 0, len(nodeMap))
	for oid := range nodeMap {
		allOIDs = append(allOIDs, oid)
	}
	for _, oid := range allOIDs {
		ensureAncestors(oid, nodeMap)
	}

	// Build the tree by mapping each node to its parent.
	// Parent OID is derived by trimming the last component.
	for oid, tn := range nodeMap {
		if oid == "1" {
			continue // root has no parent
		}
		parentOID := parentOIDOf(oid)
		if parentOID == "" {
			continue
		}
		if parent, ok2 := nodeMap[parentOID]; ok2 {
			parent.Children = append(parent.Children, tn)
		}
	}

	// Sort children at each level by OID for deterministic output.
	sortChildren(root)

	return root, nil
}

// ensureAncestors creates synthetic ancestor nodes for the given OID so the
// tree is fully connected back to the root "1".
func ensureAncestors(oid string, nodeMap map[string]*MIBTreeNode) {
	for {
		parent := parentOIDOf(oid)
		if parent == "" {
			break
		}
		if _, exists := nodeMap[parent]; exists {
			break // ancestor already exists, all further ancestors do too
		}
		// Try to look up the node in gosmi for proper metadata.
		parsedOID, err := types.OidFromString(parent)
		if err == nil {
			if sn, err := gosmi.GetNodeByOID(parsedOID); err == nil {
				nodeMap[parent] = smiNodeToTreeNode(sn)
				oid = parent
				continue
			}
		}
		// Fallback: create a synthetic node.
		nodeMap[parent] = &MIBTreeNode{
			OID:      parent,
			Name:     parent,
			Module:   "<well-known>",
			NodeType: "node",
		}
		oid = parent
	}
}

// smiNodeToTreeNode converts a gosmi SmiNode to a MIBTreeNode.
func smiNodeToTreeNode(sn gosmi.SmiNode) *MIBTreeNode {
	tn := &MIBTreeNode{
		OID:      sn.Oid.String(),
		Name:     sn.Name,
		Module:   sn.GetModule().Name,
		NodeType: nodeKindToString(sn.Kind),
	}

	if sn.Description != "" {
		tn.Description = sn.Description
	}

	if sn.Type != nil {
		if sn.Type.Name != "" {
			tn.Syntax = sn.Type.Name
		} else {
			s := sn.Type.BaseType.String()
			if s != "" && s != "Unknown" {
				tn.Syntax = s
			}
		}
		if sn.Type.Enum != nil && len(sn.Type.Enum.Values) > 0 {
			enums := make([]EnumValue, 0, len(sn.Type.Enum.Values))
			for _, v := range sn.Type.Enum.Values {
				enums = append(enums, EnumValue{Value: v.Value, Label: v.Name})
			}
			sort.Slice(enums, func(i, j int) bool {
				return enums[i].Value < enums[j].Value
			})
			tn.Enums = enums
		}
	}

	accessStr := accessToString(sn.Access)
	if accessStr != "" {
		tn.Access = accessStr
	}

	// Populate indexes for table row (entry) nodes.
	if sn.Kind == types.NodeRow {
		indexNodes := sn.GetIndex()
		if len(indexNodes) > 0 {
			indexes := make([]IndexObject, 0, len(indexNodes))
			for _, idx := range indexNodes {
				io := IndexObject{
					Name:   idx.Name,
					OID:    idx.RenderNumeric(),
					Module: idx.GetModule().Name,
				}
				if idx.Type != nil {
					io.Syntax = idx.Type.Name
				}
				if idx.Description != "" {
					io.Description = idx.Description
				}
				indexes = append(indexes, io)
			}
			tn.Indexes = indexes
		}
	}

	return tn
}

// nodeKindToString converts a gosmi NodeKind to a lowercase string.
func nodeKindToString(k types.NodeKind) string {
	return strings.ToLower(k.String())
}

// accessToString converts a gosmi Access value to a lowercase-hyphenated
// string matching SNMP conventions (e.g., "read-only", "read-write",
// "not-accessible"). Returns empty string for Unknown or NotImplemented.
func accessToString(a types.Access) string {
	switch a {
	case types.AccessUnknown, types.AccessNotImplemented:
		return ""
	default:
		// Convert CamelCase like "ReadOnly" to "read-only".
		return camelToHyphen(a.String())
	}
}

// camelToHyphen converts a CamelCase string to lowercase-hyphenated.
// Example: "ReadOnly" -> "read-only", "NotAccessible" -> "not-accessible"
func camelToHyphen(s string) string {
	var b strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			b.WriteByte('-')
		}
		b.WriteRune(unicode.ToLower(r))
	}
	return b.String()
}

// parentOIDOf returns the parent OID by trimming the last dotted component.
// Example: "1.3.6.1.2" -> "1.3.6.1"
func parentOIDOf(oid string) string {
	idx := strings.LastIndex(oid, ".")
	if idx < 0 {
		return ""
	}
	return oid[:idx]
}

// sortChildren recursively sorts children of each node by OID using numeric
// comparison of OID components.
func sortChildren(n *MIBTreeNode) {
	if len(n.Children) == 0 {
		return
	}
	sort.Slice(n.Children, func(i, j int) bool {
		return compareOIDs(n.Children[i].OID, n.Children[j].OID) < 0
	})
	for _, child := range n.Children {
		sortChildren(child)
	}
}

// compareOIDs compares two dotted-decimal OID strings numerically.
// Returns negative if a < b, 0 if equal, positive if a > b.
func compareOIDs(a, b string) int {
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")

	minLen := len(aParts)
	if len(bParts) < minLen {
		minLen = len(bParts)
	}

	for i := 0; i < minLen; i++ {
		av := oidComponent(aParts[i])
		bv := oidComponent(bParts[i])
		if av != bv {
			return av - bv
		}
	}

	return len(aParts) - len(bParts)
}

// oidComponent parses a single OID component string to int.
func oidComponent(s string) int {
	n := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		}
	}
	return n
}
