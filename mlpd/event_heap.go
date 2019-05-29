package mlpd

import "errors"

// EventHeapObject heap event if exinfo == TYPE_HEAP_OBJECT
type EventHeapObject struct {
	// [object: sleb128] the object as a difference from obj_base
	object int64
	// [vtable: sleb128] MonoVTable* as a pointer difference from ptr_base
	vtable int64
	// [size: uleb128] size of the object on the heap
	size uint64
	// [numRefs: uleb128] number of object references
	numRefs uint64
	// each referenced objref is preceded by a uleb128 encoded offset: the
	// first offset is from the object address and each next offset is relative
	// to the previous one
	// [objrefs: sleb128]+ object referenced as a difference from obj_base
	// The same object can appear multiple times, but only the first time
	// with size != 0: in the other cases this data will only be used to
	// provide additional referenced objects.
	objRefs []int64
}

// EventHeapRoot heap event if exinfo == TYPE_HEAP_ROOT
type EventHeapRoot struct {
	// 	[numRoots: uleb128] number of root references
	numRoots uint64
	// 	for i = 0 to num_roots
	// 		[address: sleb128] the root address as a difference from ptr_base
	// 		[object: sleb128] the object address as a difference from obj_base
	address []int64
	object  []int64
}

// EventHeapRootRegister heap event if exinfo == TYPE_HEAP_ROOT_REGISTER
type EventHeapRootRegister struct {
	// [start: sleb128] start address as a difference from ptr_base
	start int64
	// [size: uleb] size of the root region
	size uint64
	// [source: byte] MonoGCRootSource enum value
	source MonoGCRootSource
	// [key: sleb128] root key, meaning dependent on type, value as a difference from ptr_base
	key int64
	// [desc: string] description of the root
	desc string
}

// EventHeapRootUnregister heap event if exinfo == TYPE_HEAP_ROOT_UNREGISTER
type EventHeapRootUnregister struct {
	// [start: sleb128] start address as a difference from ptr_base
	start int64
}

// Name name of the event
func (ev *EventHeapObject) Name() string {
	return "EventHeapObject"
}

// Name name of the event
func (ev *EventHeapRoot) Name() string {
	return "EventHeapRoot"
}

// Name name of the event
func (ev *EventHeapRootRegister) Name() string {
	return "EventHeapRootRegister"
}

// Name name of the event
func (ev *EventHeapRootUnregister) Name() string {
	return "EventHeapRootUnregister"
}

// ReadEventHeap reads EventHeap from reader
func ReadEventHeap(r *MlpdReader, base *Event) (EventData, error) {
	exInfo := base.ExtendedInfo()
	switch exInfo {
	case TypeHeapObject:
		ev := &EventHeapObject{
			object:  r.readLEB128(),
			vtable:  r.readLEB128(),
			size:    r.readULEB128(),
			numRefs: r.readULEB128(),
		}
		objRefs := make([]int64, ev.numRefs)

		for i := uint64(0); i < ev.numRefs; i++ {
			objRefs[i] = r.readLEB128()
		}
		ev.objRefs = objRefs
		return ev, nil
	case TypeHeapRoot:
		numRoots := r.readULEB128()
		address := make([]int64, numRoots)
		object := make([]int64, numRoots)
		for i := uint64(0); i < numRoots; i++ {
			address[i] = r.readLEB128()
			object[i] = r.readLEB128()
		}
		return &EventHeapRoot{
			numRoots: numRoots,
			address:  address,
			object:   object,
		}, nil
	case TypeHeapRootRegister:
		return &EventHeapRootRegister{
			start:  r.readLEB128(),
			size:   r.readULEB128(),
			source: MonoGCRootSource(r.readByte()),
			key:    r.readLEB128(),
			desc:   r.readCString(),
		}, nil
	case TypeHeapRootUnregister:
		return &EventHeapRootUnregister{
			start: r.readLEB128(),
		}, nil
	}
	return nil, errors.New("not implemented")
}
