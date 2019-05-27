package mlpd

/* type for header ID */
const (
	LogHeaderID = 0x4D505A01
	BufID       = 0x4D504C01
)

/* type for event */
const (
	TypeAlloc = iota
	TypeGC
	TypeMetadata
	TypeMethod
	TypeException
	TypeMonitor
	TypeHeap
	TypeSample
	TypeRuntime
	TypeCoverage
	TypeMeta
	TypeFileIO
)

/* extended type for TYPE_HEAP */
const (
	TypeHeapStart = (iota + 1) << 4
	TypeHeapEnd
	TypeHeapObject
	TypeHeapRoot
	TypeHeapRootRegister
	TypeHeapRootUnregister
)

/* extended type for TYPE_METADATA */
const (
	TypeEndLoad   = 2 << 4
	TypeEndUnload = 4 << 4
)

/* extended type for TYPE_GC */
const (
	TypeGCEvent = (iota + 1) << 4
	TypeGCResize
	TypeGCMove
	TypeGCHandleCreated
	TypeGCHandleDestroyed
	TypeGCHandleCreatedBT
	TypeGCHandleDestroyedBT
	TypeGCFinalizeStart
	TypeGCFinalizeEnd
	TypeGCFinalizeObjectStart
	TypeGCFinalizeObjectEnd
)

/* extended type for TYPE_METHOD */
const (
	TypeLeave = (iota + 1) << 4
	TypeEnter
	TypeExcLeave
	TypeJIT
)

/* extended type for TYPE_EXCEPTION */
const (
	TypeThrowNoBT = 0 << 7
	TypeThrowBT   = 1 << 7
	TypeClause    = 1 << 4
)

/* extended type for TYPE_ALLOC */
const (
	TypeAllocNoBT = iota << 4
	TypeAllocBT
)

/* extended type for TYPE_MONITOR */
const (
	TypeMonitorNoBT = iota << 7
	TypeMonitorBT
)

/* extended type for TYPE_SAMPLE */
const (
	TypeSampleHit = iota << 4
	TypeSampleUsym
	TypeSampleUbin
	TypeSampleCountersDesc
	TypeSampleCounters
)

/* extended type for TYPE_RUNTIME */
const (
	TypeJITHelper = 1 << 4
)

/* extended type for TYPE_COVERAGE */
const (
	TypeCoverageAssembly = iota << 4
	TypeCoverageMethod
	TypeCoverageStatement
	TypeCoverageClass
)

/* extended type for TYPE_META */
const (
	TypeSyncPoint = 0 << 4
)

/* metadata type byte for TYPE_METADATA */
const (
	MetadataTypeClass = iota + 1
	MetadataTypeImage
	MetadataTypeAssembly
	MetadataTypeDomain
	MetadataTypeThread
	MetadataTypeContext
	MetadataTypeVtable
)

// MonoProfilerCodeBufferType buffer type
type MonoProfilerCodeBufferType byte

/* The \c data parameter is a \c MonoMethod pointer. */
const (
	MonoProfilerCodeBufferMethod MonoProfilerCodeBufferType = iota
	MonoProfilerCodeBufferMethodTrampoline
	MonoProfilerCodeBufferUnboxTrampoline
	MonoProfilerCodeBufferImtTrampoline
	MonoProfilerCodeBufferGenericsTrampoline
	MonoProfilerCodeBufferSpecificTrampoline
	MonoProfilerCodeBufferHelper
	MonoProfilerCodeBufferMonitor
	MonoProfilerCodeBufferDelegateInvoke
	MonoProfilerCodeBufferExceptionHandling
)

// MonoProfilerMonitorEvent monitor event enum
type MonoProfilerMonitorEvent byte

// enum values for MonoProfilerMonitorEvent
const (
	MonoProfilerMonitorContention MonoProfilerMonitorEvent = iota + 1
	MonoProfilerMonitorDone
	MonoProfilerMonitorFail
)

// MonoGCRootSource mono gc root source enum
type MonoGCRootSource byte

const (
	// MonoRootSourceExternal Roots external to Mono.  Embedders may only use this value.
	MonoRootSourceExternal = iota
	// MonoRootSourceStack Thread stack.  Must not be used to register roots.
	MonoRootSourceStack
	// MonoRootSourceFinalizerQueue Roots in the finalizer queue.  Must not be used to register roots.
	MonoRootSourceFinalizerQueue
	// MonoRootSourceStatic Managed static variables.
	MonoRootSourceStatic
	// MonoRootSourceThreadStatic Static variables with ThreadStaticAttribute.
	MonoRootSourceThreadStatic
	// MonoRootSourceContextStatic Static variables with ContextStaticAttribute.
	MonoRootSourceContextStatic
	// MonoRootSourceGCHandle GCHandle structures.
	MonoRootSourceGCHandle
	// MonoRootSourceJIT Roots in the just-in-time compiler.
	MonoRootSourceJIT
	// MonoRootSourceThreading Roots in the threading subsystem.
	MonoRootSourceThreading
	// MonoRootSourceDomain Roots in application domains.
	MonoRootSourceDomain
	// MonoRootSourceReflection Roots in reflection code.
	MonoRootSourceReflection
	// MonoRootSourceMarshal Roots from P/Invoke or other marshaling.
	MonoRootSourceMarshal
	// MonoRootSourceThreadPool Roots in the thread pool data structures.
	MonoRootSourceThreadPool
	// MonoRootSourceDebugger Roots in the debugger agent.
	MonoRootSourceDebugger
	// MonoRootSourceHandle Handle structures, used for object passed to internal functions
	MonoRootSourceHandle
)
