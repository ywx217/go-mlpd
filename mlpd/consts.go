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
	TypeHeapStart = iota
	TypeHeapEnd
	TypeHeapObject
	TypeHeapRoot
	TypeHeapRootRegister
	TypeHeapRootUnregister
)

/* extended type for TYPE_METADATA */
const (
	TypeEndLoad   = 2
	TypeEndUnload = 4
)

/* extended type for TYPE_GC */
const (
	TypeGCEvent = iota + 1
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
	TypeLeave = iota + 1
	TypeEnter
	TypeExcLeave
	TypeJIT
)

/* extended type for TYPE_EXCEPTION */
const (
	TypeThrowNoBT = 0 << 3
	TypeThrowBT   = 1 << 3
	TypeClause    = 1
)

/* extended type for TYPE_ALLOC */
const (
	TypeAllocNoBT = iota
	TypeAllocBT
)

/* extended type for TYPE_MONITOR */
const (
	TypeMonitorNoBT = iota << 3
	TypeMonitorBT
)

/* extended type for TYPE_SAMPLE */
const (
	TypeSampleHit = iota
	TypeSampleUsym
	TypeSampleUbin
	TypeSampleCountersDesc
	TypeSampleCounters
)

/* extended type for TYPE_RUNTIME */
const (
	TypeJITHelper = 1
)

/* extended type for TYPE_COVERAGE */
const (
	TypeCoverageAssembly = iota
	TypeCoverageMethod
	TypeCoverageStatement
	TypeCoverageClass
)

/* extended type for TYPE_META */
const (
	TypeSyncPoint = 0
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

const (
	/* Counter type, bits 0-7. */
	MonoCounterInt   = iota /* 32 bit int */
	MonoCounterUint         /* 32 bit uint */
	MonoCounterWord         /* pointer-sized int */
	MonoCounterLong         /* 64 bit int */
	MonoCounterUlong        /* 64 bit uint */
	MonoCounterDouble
	MonoCounterString       /* char* */
	MonoCounterTimeInterval /* 64 bits signed int holding usecs. */
	MonoCounterTypeMask     = 0xf
	MonoCounterCallback     = 128 /* ORed with the other values */
	MonoCounterSectionMask  = 0x00ffff00

	/* Sections, bits 8-23 (16 bits) */
	MonoCounterJIT = 1 << (iota + 8)
	MonoCounterGC
	MonoCounterMetadata
	MonoCounterGenerics
	MonoCounterSecurity
	MonoCounterRuntime
	MonoCounterSystem
	MonoCounterPerfcounters
	MonoCounterProfiler
	MonoCounterLastSection = (1 << 16) + 1

	/* Unit, bits 24-27 (4 bits) */
	MonoCounterUnitShift  = 24
	MonoCounterUnitMask   = 0x0f << MonoCounterUnitShift
	MonoCounterRaw        = iota << 24 /* Raw value */
	MonoCounterBytes                   /* Quantity of bytes. RSS, active heap, etc */
	MonoCounterTime                    /* Time interval in 100ns units. Minor pause, JIT compilation*/
	MonoCounterCount                   /*  Number of things (threads, queued jobs) or Number of events triggered (Major collections, Compiled methods).*/
	MonoCounterPercentage              /* [0-1] Fraction Percentage of something. Load average. */

	/* Monotonicity bits 28-31 (4 bits) */
	MonoCounterVarianceShift = 28
	MonoCounterVarianceMask  = 0x0f << MonoCounterVarianceShift
	MonoCounterMonotonic     = 1 << (iota + 28) /* This counter value always increase/decreases over time. Reported by --stat. */
	MonoCounterConstant                         /* Fixed value. Used by configuration data. */
	MonoCounterVariable                         /* This counter value can be anything on each sampling. Only interesting when sampling. */
)

// MonoProfilerSyncPointType enum
type MonoProfilerSyncPointType byte

// enum values of MonoProfilerSyncPointType
const (
	SyncPointPeriodic MonoProfilerSyncPointType = iota
	SyncPointWorldStop
	SyncPointWorldStart
)
