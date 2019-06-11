package flamewriter

// Record flame record with call stack info and statistical values
type Record struct {
	name     string
	value    int
	children map[string]*Record
}

// NewRecord creates a new Record
func NewRecord(name string, value int) *Record {
	return &Record{
		name:     name,
		value:    value,
		children: make(map[string]*Record, 0),
	}
}

// Add adds a callstack
func (r *Record) Add(stack []string, value int) {
	r.value += value
	if len(stack) == 0 {
		return
	}
	if child, ok := r.children[stack[0]]; ok {
		child.Add(stack[1:], value)
	} else {
		child = NewRecord(stack[0], 0)
		r.AddChild(child)
		child.Add(stack[1:], value)
	}
}

// AddChild adds a child record
func (r *Record) AddChild(child *Record) {
	r.children[child.name] = child
}

// ReduceRoot returns the only child if len(children) == 1, otherwise the root itself
func (r *Record) ReduceRoot() *Record {
	if len(r.children) == 1 {
		for _, child := range r.children {
			return child
		}
	}
	return r
}

// FixRootValue fix value using sum of children values
func (r *Record) FixRootValue() *Record {
	r.value = 0
	for _, child := range r.children {
		r.value += child.value
	}
	return r
}

// ValueStatisticInplace count on each node's value in the given node value -> node count map
func (r *Record) ValueStatisticInplace(m map[int]int) *Record {
	m[r.value]++
	for _, child := range r.children {
		child.ValueStatisticInplace(m)
	}
	return r
}

// CutoffInplace cuts off entries less than threshold in-place
func (r *Record) CutoffInplace(threshold int) *Record {
	cutNames := make([]string, 0)
	for name, child := range r.children {
		if child.value < threshold {
			cutNames = append(cutNames, name)
		} else {
			child.CutoffInplace(threshold)
		}
	}
	for _, name := range cutNames {
		delete(r.children, name)
	}
	return r
}
