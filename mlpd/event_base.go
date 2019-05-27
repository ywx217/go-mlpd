package mlpd

// EventBase event base header info
type EventBase struct {
	// extended info: upper 4 bits] [type: lower 4 bits]
	eventType uint8
	// time diff: uleb128] nanoseconds since last timing
}

// ExtendedInfo get extended info from upper 4 bits
func (e *EventBase) ExtendedInfo() uint8 {
	return (0xf0 & e.eventType) << 4
}

// Type get type from lower 4 bits
func (e *EventBase) Type() uint8 {
	return 0x0f & e.eventType
}
