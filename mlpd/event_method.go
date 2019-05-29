package mlpd

// EventMethod event method
type EventMethod struct {
	// method MonoMethod* as a pointer difference from the last such
	method int64
}

// EventMethodJIT event method when exinfo == TYPE_JIT
type EventMethodJIT struct {
	// method MonoMethod* as a pointer difference from the last such
	method int64
	// codeAddress pointer to the native code as a diff from ptr_base
	codeAddress int64
	// codeSize size of the generated code
	codeSize uint64
	// name full method name
	name string
}

// Name name of the event
func (ev *EventMethodJIT) Name() string {
	return "EventMethodJIT"
}

// Name name of the event
func (ev *EventMethod) Name() string {
	return "EventMethod"
}

// ReadEventMethod reads EventMethod from reader
func ReadEventMethod(r *MlpdReader, base *Event) (EventData, error) {
	exInfo := base.ExtendedInfo()
	method := r.readLEB128()
	if exInfo == TypeJIT {
		return &EventMethodJIT{
			method:      method,
			codeAddress: r.readLEB128(),
			codeSize:    r.readULEB128(),
			name:        r.readCString(),
		}, nil
	}
	if exInfo == TypeLeave || exInfo == TypeEnter || exInfo == TypeExcLeave {
		return &EventMethod{
			method: method,
		}, nil
	}
	return nil, makeExInfoError("method", exInfo)
}
