package mlpd

// EventMonitor alloc event
type EventMonitor struct {
	// tp MonoProfilerMonitorEvent enum value
	tp MonoProfilerMonitorEvent
	// object the lock object as a difference from obj_base
	object int64
	// If exinfo == TYPE_MONITOR_BT, a backtrace follows.
	bt *Backtrace
}

// Name name of the event
func (ev *EventMonitor) Name() string {
	return "EventMonitor"
}

// ReadEventMonitor reads EventMonitor from reader
func ReadEventMonitor(r *MlpdReader, base *Event) (*EventMonitor, error) {
	ev := &EventMonitor{
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
