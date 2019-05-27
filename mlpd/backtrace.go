package mlpd

import "errors"

// Backtrace backtrace in event
type Backtrace struct {
	// frames MonoMethod* as a pointer difference from the last such pointer or the buffer method_base
	frames []int64
}

// ReadBacktrace reads a backtrace from reader
func ReadBacktrace(r *MlpdReader) (*Backtrace, error) {
	num := r.readULEB128()
	if num > 1024 {
		return nil, errors.New("invalid backtrace size")
	}
	bt := &Backtrace{
		frames: make([]int64, num),
	}
	for i := uint64(0); i < num; i++ {
		bt.frames[i] = r.readLEB128()
	}
	return bt, nil
}
