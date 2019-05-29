package mlpd

import "errors"

// EventRuntime runtime event
type EventRuntime struct {
	// tp MonoProfilerCodeBufferType enum value
	tp MonoProfilerCodeBufferType
	// bufferAddress pointer to the native code as a diff from ptr_base
	bufferAddress int64
	// bufferSize size of the generated code
	bufferSize uint64
	// name buffer description name if type == MONO_PROFILER_CODE_BUFFER_SPECIFIC_TRAMPOLINE
	name string
}

// Name name of the event
func (ev *EventRuntime) Name() string {
	return "EventRuntime"
}

// ReadEventRuntime reads EventRuntime from reader
func ReadEventRuntime(r *MlpdReader, base *Event) (*EventRuntime, error) {
	exInfo := base.ExtendedInfo()
	if exInfo != TypeJITHelper {
		return nil, errors.New("Unexpected exinfo for EventRuntime")
	}
	ev := &EventRuntime{
		tp:            MonoProfilerCodeBufferType(r.readByte()),
		bufferAddress: r.readLEB128(),
		bufferSize:    r.readULEB128(),
	}
	if ev.tp == MonoProfilerCodeBufferSpecificTrampoline {
		ev.name = r.readCString()
	}
	return ev, nil
}
