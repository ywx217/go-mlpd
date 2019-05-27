package mlpd

// EventMethod event method
type EventMethod struct {
	base *EventBase
	// method MonoMethod* as a pointer difference from the last such
	method int64
}

// EventMethodJIT event method when exinfo == TYPE_JIT
type EventMethodJIT struct {
	base *EventBase
	// method MonoMethod* as a pointer difference from the last such
	method int64
	// codeAddress pointer to the native code as a diff from ptr_base
	codeAddress int64
	// codeSize size of the generated code
	codeSize uint64
	// name full method name
	name string
}

// IsEventMethod find out if its an EventMethod
func IsEventMethod(base *EventBase) bool {
	return base.Type() == TypeMethod
}

// ReadEventMethod reads EventMethod from reader
func ReadEventMethod(r *MlpdReader, base *EventBase) (interface{}, error) {
	exInfo := base.ExtendedInfo()
	method := r.readLEB128()
	if exInfo == TypeJIT {
		return &EventMethodJIT{
			base:        base,
			method:      method,
			codeAddress: r.readLEB128(),
			codeSize:    r.readULEB128(),
			name:        r.readCString(),
		}, nil
	}
	return &EventMethod{
		base:   base,
		method: method,
	}, nil
}
