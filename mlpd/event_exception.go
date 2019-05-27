package mlpd

// EventException exception event
type EventException struct {
	base *EventBase
	// object the exception object as a difference from obj_base
	object int64
	// If exinfo == TYPE_THROW_BT, a backtrace follows.
	bt *Backtrace
}

// EventExceptionClause exception event if exinfo == TYPE_CLAUSE
type EventExceptionClause struct {
	base *EventBase
	// clauseType MonoExceptionEnum enum value
	clauseType byte
	// clauseIndex index of the current clause
	clauseIndex uint64
	// method MonoMethod* as a pointer difference from the last such pointer or the buffer method_base
	method int64
	// object the exception object as a difference from obj_base
	object int64
}

// IsEventException find out if its an EventException
func IsEventException(base *EventBase) bool {
	return base.Type() == TypeException
}

// ReadEventException reads EventException from reader
func ReadEventException(r *MlpdReader, base *EventBase) (interface{}, error) {
	exInfo := base.ExtendedInfo()
	if exInfo == TypeClause {
		return &EventExceptionClause{
			base:        base,
			clauseType:  r.readByte(),
			clauseIndex: r.readULEB128(),
			method:      r.readLEB128(),
			object:      r.readLEB128(),
		}, nil
	} else {
		ev := &EventException{
			base:   base,
			object: r.readLEB128(),
		}
		if exInfo == TypeThrowBT {
			if bt, err := ReadBacktrace(r); err == nil {
				ev.bt = bt
			} else {
				return nil, err
			}
		}
		return ev, nil
	}
}
