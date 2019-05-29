package mlpd

import "fmt"

// Event event definition
type Event struct {
	// eventType [extended info: upper 4 bits] [type: lower 4 bits]
	eventType byte
	// timeDiff: uleb128] nanoseconds since last timing
	timeDiff uint64
	data     EventData
}

// EventData event data interface
type EventData interface {
	Name() string
}

// ReadEvent reads Event from reader
func ReadEvent(r *MlpdReader) (*Event, error) {
	var data EventData
	var err error

	bs, err := r.data.Peek(4)
	if len(bs) == 0 || err != nil {
		return nil, &EventEOFError{}
	}
	ev := &Event{
		eventType: r.readByte(),
		timeDiff:  r.readULEB128(),
	}
	tp := ev.Type()
	switch tp {
	case TypeAlloc:
		data, err = ReadEventAlloc(r, ev)
	case TypeGC:
		data, err = ReadEventGC(r, ev)
	case TypeMetadata:
		data, err = ReadEventMetadata(r, ev)
	case TypeMethod:
		data, err = ReadEventMethod(r, ev)
	case TypeException:
		data, err = ReadEventException(r, ev)
	case TypeRuntime:
		data, err = ReadEventRuntime(r, ev)
	case TypeMonitor:
		data, err = ReadEventMonitor(r, ev)
	case TypeHeap:
		data, err = ReadEventHeap(r, ev)
	case TypeSample:
		data, err = ReadEventSample(r, ev)
	case TypeCoverage:
		data, err = ReadEventCoverage(r, ev)
	case TypeMeta:
		data, err = ReadEventMeta(r, ev)
	default:
		return nil, fmt.Errorf("Unsupported event type %v", tp)
	}

	if err != nil {
		return nil, err
	}
	ev.data = data
	return ev, nil
}

// ExtendedInfo get extended info from upper 4 bits
func (e *Event) ExtendedInfo() uint8 {
	return (0xf0 & e.eventType) >> 4
}

// Type get type from lower 4 bits
func (e *Event) Type() uint8 {
	return 0x0f & e.eventType
}

// EventEOFError end of file when reading an event header
type EventEOFError struct{}

// Error end of file when reading an event header
func (e *EventEOFError) Error() string {
	return "end of file when reading an event header"
}

// InvalidExInfoError invalid exinfo error
type InvalidExInfoError struct {
	name   string
	exInfo uint8
}

// Error error string of invlid exinfo
func (e *InvalidExInfoError) Error() string {
	return fmt.Sprintf("invalid exinfo for %s: %v", e.name, e.exInfo)
}

func makeExInfoError(name string, exInfo uint8) error {
	return &InvalidExInfoError{
		name:   name,
		exInfo: exInfo,
	}
}
