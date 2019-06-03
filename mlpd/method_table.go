package mlpd

// MethodNode node for the method
type MethodNode struct {
	code   int64
	length uint64
	name   string
}

// MethodTable high effeciency method lookup table
type MethodTable struct {
	methodMap map[int64]*MethodNode
}

// NewMethodTable creates an empty MethodTable
func NewMethodTable() *MethodTable {
	return &MethodTable{
		methodMap: make(map[int64]*MethodNode, 0),
	}
}

// Add adds a method into the table
func (mt *MethodTable) Add(methodID, codeAddress int64, codeSize uint64, name string) {
	mt.methodMap[methodID] = &MethodNode{
		code:   codeAddress,
		length: codeSize,
		name:   name,
	}
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
	for _, v := range mt.methodMap {
		if ip >= v.code && ip < v.code+int64(v.length) {
			return v
		}
	}
	return nil
}