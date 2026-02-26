package mibresolver

import (
	"encoding/json"
	"testing"
)

func TestQueryTree(t *testing.T) {
	r, err := New(mibsDir(t))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer r.Close()

	tree, err := r.QueryTree()
	if err != nil {
		t.Fatalf("QueryTree() error: %v", err)
	}

	// Root should be OID "1" (iso).
	if tree.OID != "1" {
		t.Errorf("root OID = %q, want %q", tree.OID, "1")
	}

	// Root must have children (at least the iso.org subtree).
	if len(tree.Children) == 0 {
		t.Fatal("root has no children")
	}

	// Verify the tree is JSON-serializable.
	data, err := json.Marshal(tree)
	if err != nil {
		t.Fatalf("json.Marshal(tree) error: %v", err)
	}
	if len(data) == 0 {
		t.Error("serialized tree is empty")
	}

	t.Logf("tree root: OID=%s Name=%s Children=%d JSON_size=%d",
		tree.OID, tree.Name, len(tree.Children), len(data))
}

func TestQueryTree_CoreOnly(t *testing.T) {
	r, err := NewFromMIBData(nil)
	if err != nil {
		t.Fatalf("NewFromMIBData(nil) error: %v", err)
	}
	defer r.Close()

	tree, err := r.QueryTree()
	if err != nil {
		t.Fatalf("QueryTree() error: %v", err)
	}

	if tree.OID != "1" {
		t.Errorf("root OID = %q, want %q", tree.OID, "1")
	}
	if len(tree.Children) == 0 {
		t.Fatal("root has no children with core MIBs")
	}
}

// TestMIBTreeNodeFields verifies that a known node (sysDescr) has the
// correct name, module, and nodeType fields.
func TestMIBTreeNodeFields(t *testing.T) {
	r, err := New(mibsDir(t))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer r.Close()

	tree, err := r.QueryTree()
	if err != nil {
		t.Fatalf("QueryTree() error: %v", err)
	}

	// Find sysDescr at OID 1.3.6.1.2.1.1.1
	sysDescr := findNodeByOID(tree, "1.3.6.1.2.1.1.1")
	if sysDescr == nil {
		t.Fatal("sysDescr (1.3.6.1.2.1.1.1) not found in tree")
	}

	if sysDescr.Name != "sysDescr" {
		t.Errorf("sysDescr.Name = %q, want %q", sysDescr.Name, "sysDescr")
	}

	if sysDescr.Module != "SNMPv2-MIB" {
		t.Errorf("sysDescr.Module = %q, want %q", sysDescr.Module, "SNMPv2-MIB")
	}

	if sysDescr.NodeType != "scalar" {
		t.Errorf("sysDescr.NodeType = %q, want %q", sysDescr.NodeType, "scalar")
	}

	if sysDescr.Description == "" {
		t.Error("sysDescr.Description is empty, expected DESCRIPTION text")
	}
}

// TestMIBTreeNodeAccess verifies that access fields are populated for nodes
// with known access values.
func TestMIBTreeNodeAccess(t *testing.T) {
	r, err := New(mibsDir(t))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer r.Close()

	tree, err := r.QueryTree()
	if err != nil {
		t.Fatalf("QueryTree() error: %v", err)
	}

	// sysDescr (1.3.6.1.2.1.1.1) should be read-only.
	sysDescr := findNodeByOID(tree, "1.3.6.1.2.1.1.1")
	if sysDescr == nil {
		t.Fatal("sysDescr not found")
	}
	if sysDescr.Access != "read-only" {
		t.Errorf("sysDescr.Access = %q, want %q", sysDescr.Access, "read-only")
	}

	// sysContact (1.3.6.1.2.1.1.4) should be read-write.
	sysContact := findNodeByOID(tree, "1.3.6.1.2.1.1.4")
	if sysContact == nil {
		t.Fatal("sysContact not found")
	}
	if sysContact.Access != "read-write" {
		t.Errorf("sysContact.Access = %q, want %q", sysContact.Access, "read-write")
	}

	// system (1.3.6.1.2.1.1) is a node — should have empty access.
	system := findNodeByOID(tree, "1.3.6.1.2.1.1")
	if system == nil {
		t.Fatal("system node not found")
	}
	if system.Access != "" {
		t.Errorf("system.Access = %q, want empty", system.Access)
	}
}

// TestMIBTreeNodeSyntax verifies that syntax is populated for typed nodes.
func TestMIBTreeNodeSyntax(t *testing.T) {
	r, err := New(mibsDir(t))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer r.Close()

	tree, err := r.QueryTree()
	if err != nil {
		t.Fatalf("QueryTree() error: %v", err)
	}

	// sysDescr (1.3.6.1.2.1.1.1) has syntax DisplayString (or similar).
	sysDescr := findNodeByOID(tree, "1.3.6.1.2.1.1.1")
	if sysDescr == nil {
		t.Fatal("sysDescr not found")
	}
	if sysDescr.Syntax == "" {
		t.Error("sysDescr.Syntax is empty, expected a type name")
	}
	t.Logf("sysDescr.Syntax = %q", sysDescr.Syntax)

	// sysUpTime (1.3.6.1.2.1.1.3) has syntax TimeTicks.
	sysUpTime := findNodeByOID(tree, "1.3.6.1.2.1.1.3")
	if sysUpTime == nil {
		t.Fatal("sysUpTime not found")
	}
	if sysUpTime.Syntax == "" {
		t.Error("sysUpTime.Syntax is empty, expected a type name")
	}
	t.Logf("sysUpTime.Syntax = %q", sysUpTime.Syntax)

	// system (1.3.6.1.2.1.1) is a container node — syntax should be empty.
	system := findNodeByOID(tree, "1.3.6.1.2.1.1")
	if system == nil {
		t.Fatal("system node not found")
	}
	if system.Syntax != "" {
		t.Errorf("system.Syntax = %q, want empty for container node", system.Syntax)
	}
}

// TestMIBTreeNodeSorted verifies that children are sorted by OID.
func TestMIBTreeNodeSorted(t *testing.T) {
	r, err := New(mibsDir(t))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer r.Close()

	tree, err := r.QueryTree()
	if err != nil {
		t.Fatalf("QueryTree() error: %v", err)
	}

	// Check that children at every level are sorted.
	assertSorted(t, tree)
}

// assertSorted recursively checks that children are sorted by OID.
func assertSorted(t *testing.T, n *MIBTreeNode) {
	t.Helper()
	for i := 1; i < len(n.Children); i++ {
		if compareOIDs(n.Children[i-1].OID, n.Children[i].OID) >= 0 {
			t.Errorf("children of %s not sorted: %s >= %s",
				n.OID, n.Children[i-1].OID, n.Children[i].OID)
		}
	}
	for _, child := range n.Children {
		assertSorted(t, child)
	}
}

// findNodeByOID recursively searches the tree for a node with the given OID.
func findNodeByOID(n *MIBTreeNode, oid string) *MIBTreeNode {
	if n.OID == oid {
		return n
	}
	for _, child := range n.Children {
		if found := findNodeByOID(child, oid); found != nil {
			return found
		}
	}
	return nil
}

// TestMIBTreeNodeEnums verifies that enum values are populated for nodes with
// SYNTAX integer enums (e.g., ifAdminStatus: up/down/testing).
func TestMIBTreeNodeEnums(t *testing.T) {
	r, err := New(mibsDir(t))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer r.Close()

	tree, err := r.QueryTree()
	if err != nil {
		t.Fatalf("QueryTree() error: %v", err)
	}

	// ifAdminStatus (1.3.6.1.2.1.2.2.1.7) should have enums: up(1), down(2), testing(3).
	ifAdminStatus := findNodeByOID(tree, "1.3.6.1.2.1.2.2.1.7")
	if ifAdminStatus == nil {
		t.Fatal("ifAdminStatus (1.3.6.1.2.1.2.2.1.7) not found in tree")
	}

	if len(ifAdminStatus.Enums) == 0 {
		t.Fatal("ifAdminStatus.Enums is empty, expected enum values")
	}

	// Verify expected enums are present and sorted by value.
	expected := []EnumValue{
		{Value: 1, Label: "up"},
		{Value: 2, Label: "down"},
		{Value: 3, Label: "testing"},
	}

	if len(ifAdminStatus.Enums) != len(expected) {
		t.Fatalf("ifAdminStatus.Enums length = %d, want %d", len(ifAdminStatus.Enums), len(expected))
	}

	for i, e := range expected {
		got := ifAdminStatus.Enums[i]
		if got.Value != e.Value || got.Label != e.Label {
			t.Errorf("Enums[%d] = {%d, %q}, want {%d, %q}", i, got.Value, got.Label, e.Value, e.Label)
		}
	}

	// Verify a non-enum node has no enums.
	sysDescr := findNodeByOID(tree, "1.3.6.1.2.1.1.1")
	if sysDescr == nil {
		t.Fatal("sysDescr not found")
	}
	if len(sysDescr.Enums) != 0 {
		t.Errorf("sysDescr.Enums = %v, want empty for non-enum node", sysDescr.Enums)
	}
}

// TestCamelToHyphen verifies the camelCase to hyphen-separated conversion.
func TestCamelToHyphen(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"ReadOnly", "read-only"},
		{"ReadWrite", "read-write"},
		{"NotAccessible", "not-accessible"},
		{"NotImplemented", "not-implemented"},
		{"Notify", "notify"},
		{"Install", "install"},
		{"InstallNotify", "install-notify"},
		{"ReportOnly", "report-only"},
		{"EventOnly", "event-only"},
	}

	for _, tt := range tests {
		got := camelToHyphen(tt.input)
		if got != tt.want {
			t.Errorf("camelToHyphen(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
