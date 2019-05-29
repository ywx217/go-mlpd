package mlpd

// EventMetadata metadata event
type EventMetadata struct {
	// mType metadata type, one of: TYPE_CLASS, TYPE_IMAGE, TYPE_ASSEMBLY, TYPE_DOMAIN, TYPE_THREAD, TYPE_CONTEXT, TYPE_VTABLE
	mType byte
	// pointer pointer of the metadata type depending on mtype
	pointer int64

	// image MonoImage* as a pointer difference from ptr_base
	//   only for TYPE_CLASS, TYPE_ASSEMBLY
	image int64
	// name full class/file/assembly/domain/thread name
	//   only for TYPE_CLASS, TYPE_IMAGE, TYPE_ASSEMBLY, TYPE_DOMAIN && exinfo == 0, TYPE_THREAD && exinfo == 0)
	name string
	// domain domain id as pointer difference from ptr_base, can be zero for proxy VTables
	//   only for TYPE_CONTEXT, TYPE_VTABLE
	domain int64
	// class MonoClass* as a pointer difference from ptr_base
	//   only for TYPE_VTABLE
	class int64
}

// Name name of the event
func (ev *EventMetadata) Name() string {
	return "EventMetadata"
}

// ReadEventMetadata reads EventMetadata from reader
func ReadEventMetadata(r *MlpdReader, base *Event) (*EventMetadata, error) {
	exInfo := base.ExtendedInfo()
	if exInfo != TypeEndLoad && exInfo != TypeEndUnload && exInfo != 0 {
		return nil, makeExInfoError("metadata", exInfo)
	}
	ev := &EventMetadata{
		mType:   r.readByte(),
		pointer: r.readLEB128(),
	}
	mt := ev.mType
	if mt == MetadataTypeClass || mt == MetadataTypeAssembly {
		ev.image = r.readLEB128()
	}
	if mt == MetadataTypeClass || mt == MetadataTypeImage || mt == MetadataTypeAssembly || (mt == MetadataTypeDomain && exInfo == 0) || (mt == MetadataTypeThread && exInfo == 0) {
		ev.name = r.readCString()
	}
	if mt == MetadataTypeContext || mt == MetadataTypeVtable {
		ev.domain = r.readLEB128()
	}
	if mt == MetadataTypeVtable {
		ev.class = r.readLEB128()
	}
	return ev, nil
}
