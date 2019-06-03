package mlpd

import (
	"github.com/emirpasic/gods/trees/redblacktree"
	"github.com/emirpasic/gods/utils"
)

// MethodTree high effeciency for range lookup
type MethodTree struct {
	tree *redblacktree.Tree
}

// NewMethodTree creates a new MethodTree
func NewMethodTree() *MethodTree {
	return &MethodTree{
		tree: redblacktree.NewWith(utils.Int64Comparator),
	}
}

// Add adds an executable range into the tree
func (t *MethodTree) Add(addr int64, size uint64, name string) {
	t.tree.Put(addr, &MethodNode{
		code:   addr,
		length: size,
		name:   name,
	})
}

// Lookup looks up the method from the tree, using the given ip (instruction pointer)
func (t *MethodTree) Lookup(ip int64) *MethodNode {
	if node, found := t.tree.Floor(ip); found {
		if node, ok := node.Value.(*MethodNode); ok {
			return node
		}
	}
	return nil
}

// LookupByIP same with Lookup
func (t *MethodTree) LookupByIP(ip int64) *MethodNode {
	return t.Lookup(ip)
}
