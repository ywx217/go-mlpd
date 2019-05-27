package mlpd

// EventAlloc alloc event
type EventAlloc struct {
	base *EventBase
	// vtable MonoVTable* as a pointer difference from ptr_base
	vtable int64
	// obj object address as a byte difference from obj_base
	obj int64
	// size size of the object in the heap
	size uint64
	// If exinfo == TYPE_ALLOC_BT, a backtrace follows.
	bt *Backtrace
}

// IsEventAlloc find out if its an EventAlloc
func IsEventAlloc(base *EventBase) bool {
	return base.Type() == TypeAlloc
}

// ReadEventAlloc reads EventAlloc from reader
func ReadEventAlloc(r *MlpdReader, base *EventBase) (*EventAlloc, error) {
	ev := &EventAlloc{
		base:   base,
		vtable: r.readLEB128(),
		obj:    r.readLEB128(),
		size:   r.readULEB128(),
	}
	if base.ExtendedInfo() == TypeAllocBT {
		bt, err := ReadBacktrace(r)
		if err != nil {
			return ev, err
		}
		ev.bt = bt
	}
	return ev, nil
}
