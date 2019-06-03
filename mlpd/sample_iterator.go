package mlpd

import (
	"time"
)

// SampleData mono sample profile data
type SampleData struct {
	time     time.Time
	threadID uint64
	methods  []*MethodNode
}

// SampleIter sample iterator function
type SampleIter func(d *SampleData) error

// MakeEventIterFromSampleIter makes an EventIter from a given SampleIter
func MakeEventIterFromSampleIter(it SampleIter) EventIter {
	var lastHeader *BufferHeader
	var timeBase time.Time
	var methodBase, ptrBase int64
	var threadID uint64
	methodTable := NewMethodTable()
	unmanagedSymbols := NewMethodTree()

	return func(bh *BufferHeader, ev *Event, ver byte) error {
		if lastHeader != bh {
			// header changed
			timeBase = bh.timeBase
			methodBase = int64(bh.methodBase)
			ptrBase = int64(bh.ptrBase)
			threadID = bh.threadID
			lastHeader = bh
		}
		timeBase = timeBase.Add(time.Nanosecond * time.Duration(ev.timeDiff))
		switch d := ev.data.(type) {
		case *EventMethodJIT:
			methodBase += d.method
			methodTable.Add(methodBase, ptrBase+d.codeAddress, d.codeSize, d.name)
		case *EventMethod:
			methodBase += d.method
		case *EventExceptionClause:
			methodBase += d.method
		case *EventSample:
			switch ev.ExtendedInfo() {
			case TypeSampleUsym:
				unmanagedSymbols.Add(ptrBase+d.address, d.size, d.name)
			case TypeSampleHit:
				var tid uint64
				if ver > 10 {
					tid = uint64(ptrBase + d.thread)
				} else {
					tid = threadID
				}
				methods := make([]*MethodNode, 0)
				for _, ip := range d.ip {
					if node := methodTable.LookupByIP(ip); node != nil {
						methods = append(methods, node)
					} else if node := unmanagedSymbols.Lookup(ip); node != nil {
						methods = append(methods, node)
					} else {
						methods = append(methods, nil)
					}
				}
				if d.method != nil {
					for _, methodID := range d.method {
						methods = append(methods, methodTable.Lookup(methodID))
					}
				}
				err := it(&SampleData{
					time:     timeBase,
					threadID: tid,
					methods:  methods,
				})
				if err != nil {
					return nil
				}
			}
		}
		return nil
	}
}