package mlpd

// EventMeta meta event
type EventMeta struct {
	// type meta format:
	// type: TYPE_META
	// exinfo: one of: TYPE_SYNC_POINT
	// if exinfo == TYPE_SYNC_POINT
	// [type: byte] MonoProfilerSyncPointType enum value
	tp MonoProfilerSyncPointType
}

// Name name of the event
func (ev *EventMeta) Name() string {
	return "EventMeta"
}

// ReadEventMeta reads EventMeta from reader
func ReadEventMeta(r *MlpdReader, base *Event) (*EventMeta, error) {
	exInfo := base.ExtendedInfo()
	ev := &EventMeta{}
	if exInfo == TypeSyncPoint {
		ev.tp = MonoProfilerSyncPointType(r.readByte())
	}
	return ev, nil
}
