package mlpd

import "errors"

// EventRuntime runtime event
type EventRuntime struct {
	base *EventBase
	// tp MonoProfilerCodeBufferType enum value
	tp MonoProfilerCodeBufferType
	// bufferAddress pointer to the native code as a diff from ptr_base
	bufferAddress int64
	// bufferSize size of the generated code
	bufferSize uint64
	// name buffer description name if type == MONO_PROFILER_CODE_BUFFER_SPECIFIC_TRAMPOLINE
	name string
}

// IsEventRuntime find out if its an EventRuntime
func IsEventRuntime(base *EventBase) bool {
	return base.Type() == TypeRuntime
}

// ReadEventRuntime reads EventRuntime from reader
func ReadEventRuntime(r *MlpdReader, base *EventBase) (*EventRuntime, error) {
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
