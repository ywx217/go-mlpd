package mlpd

// EventGCResize if exinfo == TYPE_GC_RESIZE
type EventGCResize struct {
	// heapSize new heap size
	heapSize uint64
}

// EventGCEvent if exinfo == TYPE_GC_EVENT
type EventGCEvent struct {
	// eventType GC event (MONO_GC_EVENT_* from profiler.h) eventType byte
	eventType byte
	// generation GC generation event refers to generation byte
	generation byte
}

// EventGCMove if exinfo == TYPE_GC_MOVE
type EventGCMove struct {
	// numObjects number of object moves that follow
	numObjects uint64
	// objAddrs num_objects object pointer differences from obj_base
	//    num is always an even number: the even items are the old addresses, the odd numbers are the respective new object addresses
	objAddrs []int64
}

// EventGCHandleCreated if exinfo == TYPE_GC_HANDLE_CREATED[_BT]
type EventGCHandleCreated struct {
	// handleType MonoGCHandleType enum value
	//     upper bits reserved as flags
	handleType uint64
	// handle GC handle value
	handle uint64
	// objAddr object pointer differences from obj_base
	objAddr int64
	// If exinfo == TYPE_GC_HANDLE_CREATED_BT, a backtrace follows.
	bt *Backtrace
}

// EventGCHandleDestroyed if exinfo == TYPE_GC_HANDLE_DESTROYED[_BT]
type EventGCHandleDestroyed struct {
	// handleType MonoGCHandleType enum value, upper bits reserved as flags
	handleType uint64
	// handle GC handle value
	handle uint64
	// If exinfo == TYPE_GC_HANDLE_DESTROYED_BT, a backtrace follows.
	bt *Backtrace
}

// EventGCFinalizeObjectStart if exinfo == TYPE_GC_FINALIZE_OBJECT_START
type EventGCFinalizeObjectStart struct {
	// object the object as a difference from obj_base
	object int64
}

// EventGCFinalizeObjectEnd if exinfo == TYPE_GC_FINALIZE_OBJECT_END
type EventGCFinalizeObjectEnd struct {
	// object the object as a difference from obj_base
	object int64
}

// Name name of the event
func (ev *EventGCResize) Name() string {
	return "EventGCResize"
}

// Name name of the event
func (ev *EventGCEvent) Name() string {
	return "EventGCEvent"
}

// Name name of the event
func (ev *EventGCMove) Name() string {
	return "EventGCMove"
}

// Name name of the event
func (ev *EventGCHandleCreated) Name() string {
	return "EventGCHandleCreated"
}

// Name name of the event
func (ev *EventGCHandleDestroyed) Name() string {
	return "EventGCHandleDestroyed"
}

// Name name of the event
func (ev *EventGCFinalizeObjectStart) Name() string {
	return "EventGCFinalizeObjectStart"
}

// Name name of the event
func (ev *EventGCFinalizeObjectEnd) Name() string {
	return "EventGCFinalizeObjectEnd"
}

// ReadEventGC reads EventGC from reader
func ReadEventGC(r *MlpdReader, base *Event) (EventData, error) {
	extInfo := base.ExtendedInfo()
	switch extInfo {
	case TypeGCResize:
		return &EventGCResize{
			heapSize: r.readULEB128(),
		}, nil
	case TypeGCEvent:
		return &EventGCEvent{
			eventType:  r.readByte(),
			generation: r.readByte(),
		}, nil
	case TypeGCMove:
		numObjects := r.readULEB128()
		objAddrs := make([]int64, numObjects)
		for i := uint64(0); i < numObjects; i++ {
			objAddrs[i] = r.readLEB128()
		}
		return &EventGCMove{
			objAddrs: objAddrs,
		}, nil
	case TypeGCHandleCreated, TypeGCHandleCreatedBT:
		ev := &EventGCHandleCreated{
			handleType: r.readULEB128(),
			handle:     r.readULEB128(),
			objAddr:    r.readLEB128(),
		}
		if extInfo == TypeGCHandleCreatedBT {
			if bt, err := ReadBacktrace(r); err == nil {
				ev.bt = bt
			} else {
				return nil, err
			}
		}
		return ev, nil
	case TypeGCHandleDestroyed, TypeGCHandleDestroyedBT:
		ev := &EventGCHandleDestroyed{
			handleType: r.readULEB128(),
			handle:     r.readULEB128(),
		}
		if extInfo == TypeGCHandleDestroyedBT {
			if bt, err := ReadBacktrace(r); err == nil {
				ev.bt = bt
			} else {
				return nil, err
			}
		}
		return ev, nil
	case TypeGCFinalizeObjectStart:
		return &EventGCFinalizeObjectStart{
			object: r.readLEB128(),
		}, nil
	case TypeGCFinalizeObjectEnd:
		return &EventGCFinalizeObjectEnd{
			object: r.readLEB128(),
		}, nil
	}
	return nil, makeExInfoError("gc", extInfo)
}
