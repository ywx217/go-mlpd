package mlpd

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"time"

	"ekyu.moe/leb128"
)

// MlpdReader mono profile output file (mlpd) reader
type MlpdReader struct {
	data   *bufio.Reader
	header *Header
}

// Header file header for mlpd file
type Header struct {
	// id constant value: LOG_HEADER_ID
	id int32
	// major and minor version of the log profiler
	major, minor byte
	// format version of the data format for the rest of the file
	format byte
	// pstrSize ptrsize size in bytes of a pointer in the profiled program
	pstrSize byte
	// startupTime startup time time in milliseconds since the unix epoch when the program started
	startupTime time.Time
	// timerOverhead timer overhead approximate overhead in nanoseconds of the timer
	timerOverhead uint32
	// flags file format flags, should be 0 for now
	flags uint32
	// pid pid of the profiled process
	pid uint32
	// port tcp port for server if != 0
	port uint16
	// args arguments passed to the profiler
	args string
	// arch architecture the profiler is running on
	arch string
	// os operating system the profiler is running on
	os string
}

// BufferHeader buffer header
type BufferHeader struct {
	// id constant value: BUF_ID
	id int32
	// length size of the data following the buffer header
	length uint32
	// timeBase time base in nanoseconds since an unspecified epoch
	timeBase time.Time
	// ptrBase base value for pointers
	ptrBase uint64
	// objBase base value for object addresses
	objBase uint64
	// threadID system-specific thread ID (pthread_t for example)
	threadID uint64
	// methodBase base value for MonoMethod pointers
	methodBase uint64
}

func (r *MlpdReader) readBytes(size int) []byte {
	b := make([]byte, size)
	n, err := r.data.Read(b)
	if err != nil {
		return nil
	}
	if n == size {
		return b
	}
	return b[:n]
}

func (r *MlpdReader) readByte() byte {
	b, _ := r.data.ReadByte()
	return b
}

func (r *MlpdReader) readInt16() int16 {
	return int16(r.readUint16())
}

func (r *MlpdReader) readUint16() uint16 {
	return binary.LittleEndian.Uint16(r.readBytes(2))
}

func (r *MlpdReader) readInt32() int32 {
	return int32(r.readUint32())
}

func (r *MlpdReader) readUint32() uint32 {
	return binary.LittleEndian.Uint32(r.readBytes(4))
}

func (r *MlpdReader) readInt64() int64 {
	return int64(r.readUint64())
}

func (r *MlpdReader) readUint64() uint64 {
	return binary.LittleEndian.Uint64(r.readBytes(8))
}

func (r *MlpdReader) readSizedString() string {
	size := r.readUint32()
	b := r.readBytes(int(size))
	if len(b) > 0 && b[len(b)-1] == 0 {
		return string(b[:len(b)-1])
	}
	return string(b)
}

func (r *MlpdReader) readCString() string {
	b, err := r.data.ReadBytes(0)
	if err != nil {
		return ""
	}
	return string(b)
}

func (r *MlpdReader) readTimeInMillis() time.Time {
	ts := r.readUint64()
	return time.Unix(int64(ts/1000), int64(ts%1000))
}

func (r *MlpdReader) readTimeInNanos() time.Time {
	ts := r.readUint64()
	return time.Unix(int64(ts/1000000), int64(ts%1000000))
}

func (r *MlpdReader) readLEB128() int64 {
	if bs, _ := r.data.Peek(10); len(bs) > 0 {
		num, advance := leb128.DecodeSleb128(bs)
		if advance > 0 {
			r.readBytes(int(advance))
		}
		return num
	}
	return 0
}

func (r *MlpdReader) readULEB128() uint64 {
	if bs, _ := r.data.Peek(10); len(bs) > 0 {
		num, advance := leb128.DecodeUleb128(bs)
		if advance > 0 {
			r.readBytes(int(advance))
		}
		return num
	}
	return 0
}

// ReadHeader reads file header
func (r *MlpdReader) ReadHeader() (*Header, error) {
	if r.header != nil {
		return r.header, nil
	}
	headerID := r.readInt32()
	if headerID != LogHeaderID {
		return nil, fmt.Errorf("Invalid log header id: 0x%x", headerID)
	}
	header := &Header{
		id:            headerID,
		major:         r.readByte(),
		minor:         r.readByte(),
		format:        r.readByte(),
		pstrSize:      r.readByte(),
		startupTime:   r.readTimeInMillis(),
		timerOverhead: r.readUint32(),
		flags:         r.readUint32(),
		pid:           r.readUint32(),
		port:          r.readUint16(),
		args:          r.readSizedString(),
		arch:          r.readSizedString(),
		os:            r.readSizedString(),
	}
	return header, nil
}

// ReadBufferHeader reads buffer header of mlpd file
func (r *MlpdReader) ReadBufferHeader() (*BufferHeader, error) {
	if _, err := r.ReadHeader(); err != nil {
		return nil, err
	}
	headerID := r.readInt32()
	if headerID != BufID {
		return nil, fmt.Errorf("Invalid buffer header id: 0x%x", headerID)
	}
	bh := &BufferHeader{
		id:         headerID,
		length:     r.readUint32(),
		timeBase:   r.readTimeInNanos(),
		ptrBase:    r.readUint64(),
		objBase:    r.readUint64(),
		threadID:   r.readUint64(),
		methodBase: r.readUint64(),
	}
	return bh, nil
}
