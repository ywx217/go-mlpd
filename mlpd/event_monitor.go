package mlpd

// EventMonitor alloc event
type EventMonitor struct {
	base *EventBase
	// tp MonoProfilerMonitorEvent enum value
	tp MonoProfilerMonitorEvent
	// object the lock object as a difference from obj_base
	object int64
	// If exinfo == TYPE_MONITOR_BT, a backtrace follows.
	bt *Backtrace
}

// IsEventMonitor find out if its an EventMonitor
func IsEventMonitor(base *EventBase) bool {
	return base.Type() == TypeMonitor
}

// ReadEventMonitor reads EventMonitor from reader
func ReadEventMonitor(r *MlpdReader, base *EventBase) (*EventMonitor, error) {
	ev := &EventMonitor{
		base:   base,
		tp:     MonoProfilerMonitorEvent(r.readByte()),
		object: r.readLEB128(),
	}
	if base.ExtendedInfo() == TypeMonitorBT {
		if bt, err := ReadBacktrace(r); err == nil {
			ev.bt = bt
		} else {
			return nil, err
		}
	}
	return ev, nil
}
