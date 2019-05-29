package mlpd

import (
	"fmt"
)

// EventSample alloc event
type EventSample struct {
	// if exinfo == TYPE_SAMPLE_HIT
	// 	[thread: sleb128] thread id as difference from ptr_base
	thread int64
	// 	[count: uleb128] number of following instruction addresses
	// 	[ip: sleb128]* instruction pointer as difference from ptr_base
	ip []int64
	// [mbt_count: uleb128] number of managed backtrace frames
	// [method: sleb128]* MonoMethod* as a pointer difference from the last such
	// 	pointer or the buffer method_base (the first such method can be also indentified by ip, but this is not neccessarily true)
	method []int64
	// if exinfo == TYPE_SAMPLE_USYM
	// 	[address: sleb128] symbol address as a difference from ptr_base
	// 	[size: uleb128] symbol size (may be 0 if unknown)
	// 	[name: string] symbol name
	// if exinfo == TYPE_SAMPLE_UBIN
	// 	[address: sleb128] address where binary has been loaded as a difference from ptr_base
	// 	[offset: uleb128] file offset of mapping (the same file can be mapped multiple times)
	// 	[size: uleb128] memory size
	// 	[name: string] binary name
	address int64
	offset  uint64
	size    uint64
	name    string
	// if exinfo == TYPE_SAMPLE_COUNTERS_DESC
	// 	[len: uleb128] number of counters
	// 	for i = 0 to len
	// 		[section: uleb128] section of counter
	// 		if section == MONO_COUNTER_PERFCOUNTERS:
	// 			[section_name: string] section name of counter
	// 		[name: string] name of counter
	// 		[type: uleb128] type of counter
	// 		[unit: uleb128] unit of counter
	// 		[variance: uleb128] variance of counter
	// 		[index: uleb128] unique index of counter
	sections []SampleCounterDesc
	// if exinfo == TYPE_SAMPLE_COUNTERS
	// 	while true:
	// 		[index: uleb128] unique index of counter
	// 		if index == 0:
	// 			break
	// 		[type: uleb128] type of counter value
	// 		if type == string:
	// 			if value == null:
	// 				[0: byte] 0 -> value is null
	// 			else:
	// 				[1: byte] 1 -> value is not null
	// 				[value: string] counter value
	// 		else:
	// 			[value: uleb128/sleb128/double] counter value, can be sleb128, uleb128 or double (determined by using type)
	values []interface{}
}

// SampleCounterDesc sample counter section
type SampleCounterDesc struct {
	section     uint64
	sectionName string
	name        string
	tp          uint64
	unit        uint64
	variance    uint64
	index       uint64
}

// Name name of the event
func (ev *EventSample) Name() string {
	return "EventSample"
}

// ReadEventSample reads EventSample from reader
func ReadEventSample(r *MlpdReader, base *Event) (*EventSample, error) {
	exInfo := base.ExtendedInfo()
	ev := &EventSample{}

	switch exInfo {
	case TypeSampleHit:
		ev.thread = r.readLEB128()
		count := r.readULEB128()
		ip := make([]int64, count)
		for i := uint64(0); i < count; i++ {
			ip[i] = r.readLEB128()
		}
		ev.ip = ip
		methodCount := r.readULEB128()
		method := make([]int64, methodCount)
		for i := uint64(0); i < methodCount; i++ {
			method[i] = r.readLEB128()
		}
		ev.method = method
	case TypeSampleUsym:
		ev.address = r.readLEB128()
		ev.size = r.readULEB128()
		ev.name = r.readCString()
	case TypeSampleUbin:
		ev.address = r.readLEB128()
		ev.offset = r.readULEB128()
		ev.size = r.readULEB128()
		ev.name = r.readCString()
	case TypeSampleCountersDesc:
		descCount := r.readULEB128()
		desc := make([]SampleCounterDesc, descCount)
		for i := uint64(0); i < descCount; i++ {
			d := &desc[i]
			d.section = r.readULEB128()
			if d.section == MonoCounterPerfcounters {
				d.sectionName = r.readCString()
			}
			d.name = r.readCString()
			d.tp = r.readULEB128()
			d.unit = r.readULEB128()
			d.variance = r.readULEB128()
			d.index = r.readULEB128()
		}
		ev.sections = desc
	case TypeSampleCounters:
		values := make([]interface{}, 0)
		for {
			index := r.readULEB128()
			if index == 0 {
				break
			}
			tp := r.readULEB128()
			switch tp {
			case MonoCounterString:
				firstByte := r.readByte()
				value := ""
				if firstByte != 0 {
					value = r.readCString()
				}
				values = append(values, &value)
			case MonoCounterInt, MonoCounterWord, MonoCounterLong:
				value := r.readLEB128()
				values = append(values, &value)
			case MonoCounterUint, MonoCounterUlong:
				value := r.readULEB128()
				values = append(values, &value)
			case MonoCounterDouble:
				value := r.readFloat64()
				values = append(values, &value)
			default:
				return nil, fmt.Errorf("invalid type for sample counter: %v", tp)
			}
		}
		ev.values = values
	}

	return ev, nil
}
