package mlpd

// EventAlloc alloc event
type EventAlloc struct {
	// vtable MonoVTable* as a pointer difference from ptr_base
	vtable int64
	// obj object address as a byte difference from obj_base
	obj int64
	// size size of the object in the heap
	size uint64
	// If exinfo == TYPE_ALLOC_BT, a backtrace follows.
	bt *Backtrace
}

// Name name of the event
func (ev *EventAlloc) Name() string {
	return "EventAlloc"
}

// ReadEventAlloc reads EventAlloc from reader
func ReadEventAlloc(r *MlpdReader, base *Event) (*EventAlloc, error) {
	ev := &EventAlloc{
		vtable: r.readLEB128(),
		obj:    r.readLEB128(),
		size:   r.readULEB128(),
	}
	exInfo := base.ExtendedInfo()
	if exInfo == TypeAllocBT {
		bt, err := ReadBacktrace(r)
		if err != nil {
			return ev, err
		}
		ev.bt = bt
	}
	if exInfo == TypeAllocNoBT {
		return ev, nil
	}
	return nil, makeExInfoError("alloc", exInfo)
}
