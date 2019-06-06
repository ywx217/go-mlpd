package mlpd

// MethodNode node for the method
type MethodNode struct {
	code   int64
	length uint64
	name   string
}

// Name returns name of the node
func (node *MethodNode) Name() string {
	return node.name
}

// MethodTable high effeciency method lookup table
type MethodTable struct {
	methodMap map[int64]*MethodNode
	ipTree    *MethodTree
}

// NewMethodTable creates an empty MethodTable
func NewMethodTable() *MethodTable {
	return &MethodTable{
		methodMap: make(map[int64]*MethodNode, 0),
		ipTree:    NewMethodTree(),
	}
}

// Add adds a method into the table
func (mt *MethodTable) Add(methodID, codeAddress int64, codeSize uint64, name string) {
	mt.methodMap[methodID] = &MethodNode{
		code:   codeAddress,
		length: codeSize,
		name:   name,
	}
	mt.ipTree.Add(codeAddress, codeSize, name)
}

// Lookup looks up a method from the table
func (mt *MethodTable) Lookup(methodID int64) *MethodNode {
	if node, ok := mt.methodMap[methodID]; ok {
		return node
	}
	return nil
}

// LookupByIP looks up a method by instruction pointer
func (mt *MethodTable) LookupByIP(ip int64) *MethodNode {
	node := mt.ipTree.Lookup(ip)
	if node == nil {
		return nil
	}
	if ip >= node.code && ip < node.code+int64(node.length) {
		return node
	}
	return nil
}
