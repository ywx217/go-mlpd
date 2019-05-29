package mlpd

// EventException exception event
type EventException struct {
	// object the exception object as a difference from obj_base
	object int64
	// If exinfo == TYPE_THROW_BT, a backtrace follows.
	bt *Backtrace
}

// EventExceptionClause exception event if exinfo == TYPE_CLAUSE
type EventExceptionClause struct {
	// clauseType MonoExceptionEnum enum value
	clauseType byte
	// clauseIndex index of the current clause
	clauseIndex uint64
	// method MonoMethod* as a pointer difference from the last such pointer or the buffer method_base
	method int64
	// object the exception object as a difference from obj_base
	object int64
}

// Name name of the event
func (ev *EventException) Name() string {
	return "EventException"
}

// Name name of the event
func (ev *EventExceptionClause) Name() string {
	return "EventExceptionClause"
}

// ReadEventException reads EventException from reader
func ReadEventException(r *MlpdReader, base *Event) (EventData, error) {
	exInfo := base.ExtendedInfo()
	if exInfo == TypeClause {
		return &EventExceptionClause{
			clauseType:  r.readByte(),
			clauseIndex: r.readULEB128(),
			method:      r.readLEB128(),
			object:      r.readLEB128(),
		}, nil
	}
	if exInfo == TypeThrowBT || exInfo == TypeThrowNoBT {
		ev := &EventException{
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
	return nil, makeExInfoError("exception", exInfo)
}
